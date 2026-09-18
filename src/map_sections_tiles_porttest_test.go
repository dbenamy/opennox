//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func mapSectionsRun(t *testing.T, cases []legacy.PortTestMapSectionSpec) []legacy.PortTestMapSectionResult {
	t.Helper()
	out := legacy.PortTestMapSections(cases, t.TempDir(), func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
	for _, r := range out {
		if !r.Paint.Intact || !r.Paint.ControlOK {
			t.Fatal(r.Name, "fixture ownership/control state")
		}
	}
	return out
}
func mapSectionTileBytes(nodes [][4]uint32) []byte {
	out := make([]byte, 0, 6+5*(len(nodes)-1))
	for i, node := range nodes {
		out = append(out, byte(node[0]), 0, 0, byte(node[2]), byte(node[3]))
		binary.LittleEndian.PutUint16(out[len(out)-4:], uint16(node[1]))
		if i == 0 {
			out = append(out, byte(len(nodes)-1))
		}
	}
	return out
}
func mapSectionTileNarrow(n [4]uint32) [4]uint32 {
	return [4]uint32{uint32(byte(n[0])), uint32(int32(int16(n[1]))), uint32(byte(n[2])), uint32(byte(n[3]))}
}
func TestMapSectionsTileRecords(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	var inputs [][][4]uint32
	for _, read := range []bool{false, true} {
		for _, version := range []int{1, 3, 4} {
			for _, count := range []int{0, 1, 2, 9, 10, 11, 254, 255, 256} {
				for _, value := range []uint32{0, 127, 128, 255, 32768, 0xffffffff} {
					nodes := make([][4]uint32, count+1)
					for i := range nodes {
						v := value + uint32(i*73)
						nodes[i] = [4]uint32{v, v + 0x12340000, v + 31, v + 129}
					}
					raw := mapSectionTileBytes(nodes)
					name := fmt.Sprintf("read%v/v%d/count%d/value%x", read, version, count, value)
					sp := legacy.PortTestMapSectionSpec{Name: name, Paint: legacy.PortTestPaintSpec{Seed: 23,
						Records: []legacy.PortTestMapRoomRecord{{Size: 20, Words: map[int]uint32{0: nodes[0][0], 4: nodes[0][1], 8: nodes[0][2], 12: nodes[0][3]}}},
						Actions: []legacy.PortTestPaintAction{{Op: 0, Args: [6]legacy.PortTestMapRoomArg{{Value: uint32(version)}, {Slot: 1}}}},
					}, IO: []legacy.PortTestMapSectionIO{{Function: "tile", Read: read, Data: raw}}}
					if !read {
						sp.ChainSlot = 1
						sp.Chain = nodes[1:]
					}
					cases = append(cases, sp)
					inputs = append(inputs, nodes)
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			nodes := inputs[i]
			got := r.IO[0]
			wire := mapSectionTileBytes(nodes)
			if !bytes.Equal(got.Data, wire) {
				t.Fatalf("tile bytes: got %x want %x", got.Data, wire)
			}
			count := len(nodes) - 1
			read := cases[i].IO[0].Read
			if read {
				count = int(byte(count))
			}
			if got.Position != int64(6+5*count) {
				t.Fatalf("file position %d want %d", got.Position, 6+5*count)
			}
			ret := uint32(0)
			if count != 0 {
				ret = 1
			}
			if got.Return != ret {
				t.Fatalf("return %d want %d", got.Return, ret)
			}
			if len(got.Tile) != count+1 {
				t.Fatalf("chain nodes %d want %d", len(got.Tile), count+1)
			}
			for j, n := range got.Tile {
				want := nodes[j]
				if read || j == 0 {
					want = mapSectionTileNarrow(want)
				}
				if n != want {
					t.Fatalf("node%d: %x want %x", j, n, want)
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-tile-records", out, "c2f5fba52dd7cfa74468c5b35cbe38c4efe9e90a6b218a65b251bc865b730b39")
}
