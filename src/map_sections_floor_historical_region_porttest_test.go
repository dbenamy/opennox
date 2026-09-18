//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestMapSectionsFloorHistoricalRegionRead(t *testing.T) {
	node := [4]uint32{37, 0xffff8001, 19, 23}
	var cases []legacy.PortTestMapSectionSpec
	type expected struct {
		x, y, half int
		generation bool
	}
	var wants []expected
	for _, generation := range []bool{false, true} {
		for _, delta := range [][2]int{{0, 0}, {2, 4}, {-2, -4}} {
			for _, point := range [][2]int{{10, 10}, {11, 11}, {12, 14}, {13, 15}} {
				targetX, targetY := 10+delta[0], 10+delta[1]
				quad := [8]uint32{uint32((targetX + 2) * 23), uint32(targetY * 23), uint32(targetX * 23), uint32((targetY + 2) * 23), uint32((targetX + 4) * 23), uint32((targetY + 2) * 23), uint32((targetX + 2) * 23), uint32((targetY + 4) * 23)}
				words := map[int]uint32{}
				for i, v := range quad {
					words[4*i] = v
				}
				wire := []byte{3, 0}
				for _, v := range []uint32{10, 10, 5, 6} {
					wire = binary.LittleEndian.AppendUint32(wire, v)
				}
				for y := 10; y < 16; y++ {
					for x := 10; x < 15; x++ {
						if x == point[0] && y == point[1] {
							wire = append(wire, 1)
							for _, v := range node {
								wire = binary.LittleEndian.AppendUint32(wire, v)
							}
							wire = append(wire, 0)
						} else {
							wire = append(wire, 0)
						}
					}
				}
				flags := uint32(0)
				if generation {
					flags = 0x400000
				}
				sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("generation%v/delta%v/point%v", generation, delta, point), Flags: flags,
					Paint: legacy.PortTestPaintSpec{Seed: 37, Records: []legacy.PortTestMapRoomRecord{{Size: 32, Words: words}}, Actions: []legacy.PortTestPaintAction{{Op: 0, Args: [6]legacy.PortTestMapRoomArg{{Slot: 1}}}}},
					IO:    []legacy.PortTestMapSectionIO{{Function: "floor", Read: true, Data: wire}},
				}
				x, y := point[0]+delta[0], point[1]+delta[1]
				half := 2
				if y&1 != 0 {
					half = 1
					x = (x - 1) / 2
					y = (y + 1) / 2
				} else {
					x /= 2
					y /= 2
				}
				cases = append(cases, sp)
				wants = append(wants, expected{x, y, half, generation})
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			want := wants[i]
			wire := r.IO[0]
			state := r.Paint.Steps[0]
			if wire.Return != 1 || wire.Position != int64(len(wire.Data)) {
				t.Fatal("region return/consumption")
			}
			n := node
			if !want.generation {
				if len(state.Cells) != 1 {
					t.Fatalf("cells %+v", state.Cells)
				}
				c := state.Cells[0]
				base := 1
				if want.half == 2 {
					base = 6
				}
				if c.X != want.x || c.Y != want.y || c.Words[0] != uint32(want.half) {
					t.Fatal("relocated cell", c, want)
				}
				for j, v := range n {
					if c.Words[base+j] != v {
						t.Fatal("region tile fields", c, n)
					}
				}
			} else {
				if len(state.Cells) != 0 {
					t.Fatal("generation wrote live grid")
				}
				dataCount, nodeCount := 0, 0
				for _, record := range state.Records {
					if !record.Alive {
						continue
					}
					if record.Kind == "input" && len(record.Words) == 5 {
						dataCount++
						for j, v := range n {
							if record.Words[j] != v {
								t.Fatal("generation tile fields", record, n)
							}
						}
					}
					if record.Kind == "section-allocation" && len(record.Words) == 6 {
						nodeCount++
						x, y := float32(want.x*46), float32(want.y*46)
						if want.half == 1 {
							x += 23
						} else {
							y += 23
						}
						if record.Words[1] != math.Float32bits(x) || record.Words[2] != math.Float32bits(y) || record.Words[3] != uint32(want.half) {
							t.Fatal("generation tile position/orientation", record, want)
						}
					}
				}
				if dataCount != 1 || nodeCount != 1 {
					t.Fatal("generation allocation owners", dataCount, nodeCount)
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-floor-historical-region", out, "cbf18d8bf7739ae2ff2e6be3bf769c04c596b5f392e6b3fc221a6d34f99b84b2")
}
