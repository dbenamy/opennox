//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Two outgoing links distinguish element indexing from byte-count pointer steps.
// Use the actual pending-waypoint owner and resolver, not a simulated graph.
func TestPrefabRuntimeWaypointLinkLayout(t *testing.T) {
	_, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	oldServer, oldFile := noxServer, cryptfile.Global()
	noxServer = &Server{Server: &server.Server{}}
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile); noxServer = oldServer })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	path := filepath.Join(t.TempDir(), "waypoints.bin")
	for _, version := range []uint16{1, 2, 3, 4} {
		func() {
			var b bytes.Buffer
			binary.Write(&b, binary.LittleEndian, version)
			binary.Write(&b, binary.LittleEndian, uint32(3))
			for i := uint32(0); i < 3; i++ {
				binary.Write(&b, binary.LittleEndian, uint32(100+i))
				x, y := uint32(10+i), uint32(20+i)
				if version >= 4 {
					x = math.Float32bits(float32(x))
					y = math.Float32bits(float32(y))
				}
				binary.Write(&b, binary.LittleEndian, x)
				binary.Write(&b, binary.LittleEndian, y)
				if version >= 3 {
					b.WriteByte(0)
				}
				binary.Write(&b, binary.LittleEndian, uint32(1))
				n := uint32(0)
				if i == 0 {
					n = 2
				}
				if version < 4 {
					binary.Write(&b, binary.LittleEndian, n)
				} else {
					b.WriteByte(byte(n))
				}
				for j := uint32(0); j < n; j++ {
					binary.Write(&b, binary.LittleEndian, uint32(101+j))
					if version >= 2 {
						b.WriteByte(byte(2 + 5*j))
					}
				}
			}
			if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer func() {
				for wp := noxServer.WPs.Pending; wp != nil; {
					next := wp.WpNext
					alloc.Free(wp)
					wp = next
				}
				noxServer.WPs.Pending = nil
				cryptfile.Close()
			}()
			if legacy.PortTestPrefabCall(34, [6]uint32{}) != 1 {
				t.Fatalf("v%d read failed", version)
			}
			var first *server.Waypoint
			for wp := noxServer.WPs.Pending; wp != nil; wp = wp.WpNext {
				if wp.Index == 100 {
					first = wp
				}
			}
			if first == nil {
				t.Fatal("missing source waypoint")
			}
			if first.Field348[0] != 101 || first.Field348[1] != 102 {
				t.Fatalf("v%d saved link IDs %v; want consecutive 101,102", version, first.Field348[:5])
			}
			if noxServer.WPs.Sub_579CA0() != 1 || first.Points[0].Waypoint.Index != 101 || first.Points[1].Waypoint.Index != 102 {
				t.Fatalf("v%d actual link resolution failed", version)
			}
		}()
	}
}
