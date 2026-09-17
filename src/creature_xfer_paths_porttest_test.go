//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferPaths(t *testing.T) {
	s := newCreatureXferOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	path := filepath.Join(t.TempDir(), "paths.bin")
	type row struct {
		Case              string
		Wire, Points      []byte
		Waypoints         []uint32
		ReadCRC, WriteCRC uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-paths", rows, "83ac69db2230cc65c5f5cef23c335266dff56079b6582e8b268f2a8f18f95b8c")
	}()
	for _, points := range []int{0, 1, 2, 32} {
		for _, waypoints := range []int{0, 1, 2, 16} {
			for _, known := range []bool{false, true} {
				t.Run(fmt.Sprintf("points%d-waypoints%d-known%v", points, waypoints, known), func(t *testing.T) {
					u := newCreatureXferObject(t, s, "Monster")
					creatureXferSetTimers(u, 1)
					objectXferSetWord(u.UpdateData, 8, uint32(points))
					objectXferSetWord(u.UpdateData, 296, uint32(waypoints))
					pointData := unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, 12)), 8*points)
					for i := range pointData {
						pointData[i] = byte(17*i + 3)
					}
					var list []*server.Waypoint
					var head *server.Waypoint
					ids := new(mapDrawableStream)
					for i := 0; i < waypoints; i++ {
						wp, free := alloc.New(server.Waypoint{})
						t.Cleanup(free)
						wp.Index = uint32(100 + i)
						wp.WpNext = head
						head = wp
						list = append(list, wp)
						objectXferSetWord(u.UpdateData, 300+4*i, uint32(uintptr(unsafe.Pointer(wp))))
						ids.u32(wp.Index)
					}
					base := creatureXferEmptyActionRecord(u, 4, s.Frame()).Bytes()
					p := new(mapDrawableStream)
					p.Write(base[:15])
					p.Write(pointData)
					p.Write(base[15:36])
					p.Write(ids.Bytes())
					p.Write(base[36:])
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					if ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0); ret != 1 {
						t.Fatalf("write=%d", ret)
					}
					writeCRC := cryptfile.Global().PortTestChecksum()
					cryptfile.Close()
					got, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(got, p.Bytes()) {
						t.Fatalf("path writer wire: got=%x want=%x", got, p.Bytes())
					}
					old := s.WPs.Pending
					if known {
						s.WPs.Pending = head
					} else {
						s.WPs.Pending = nil
					}
					defer func() { s.WPs.Pending = old }()
					v := newCreatureXferObject(t, s, "Monster")
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					if ret := legacy.PortTestCreatureXferHelper(0, v, nil, 0); ret != 1 {
						t.Fatalf("read=%d", ret)
					}
					pos, err := cryptfile.Global().File.Seek(0, 1)
					data := unsafe.Slice((*byte)(unsafe.Add(v.UpdateData, 12)), 8*points)
					if err != nil || pos != int64(p.Len()) || !bytes.Equal(data, pointData) || objectXferGetWord(v.UpdateData, 8) != uint32(points) || objectXferGetWord(v.UpdateData, 296) != uint32(waypoints) {
						t.Fatal("path counts/points/position")
					}
					normalized := make([]uint32, waypoints)
					for i, wp := range list {
						want := uint32(0)
						if known {
							want = uint32(uintptr(unsafe.Pointer(wp)))
							normalized[i] = wp.Index
						}
						if objectXferGetWord(v.UpdateData, 300+4*i) != want {
							t.Fatalf("waypoint identity %d", i)
						}
					}
					rows = append(rows, row{t.Name(), got, append([]byte(nil), data...), normalized, cryptfile.Global().PortTestChecksum(), writeCRC})
				})
			}
		}
	}
}
