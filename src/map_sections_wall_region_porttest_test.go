//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsWallRegionRead(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	type expected struct {
		x, y                           int
		version, load                  uint32
		direction, material, variation byte
		generation                     bool
	}
	var wants []expected
	for _, version := range []uint32{6, 7} {
		for _, generation := range []bool{false, true} {
			for _, load := range []uint32{0, 1} {
				for _, material := range []byte{0, 1, 2} {
					for _, direction := range []byte{0, 11} {
						for _, variation := range []byte{0, 3, 4, 255} {
							for _, point := range [][2]int{{8, 10}} {
								wire := binary.LittleEndian.AppendUint16(nil, uint16(version))
								for _, v := range []uint32{8, 10, 1, 1} {
									wire = binary.LittleEndian.AppendUint32(wire, v)
								}
								wire = append(wire, byte(point[0]), byte(point[1]), direction|128, material, variation, 173)
								if version == 7 {
									wire = append(wire, 127)
								}
								wire = append(wire, 255)
								flags := uint32(0)
								if generation {
									flags = 0x400000
								}
								quad := [8]uint32{22 * 23, 20 * 23, 20 * 23, 22 * 23, 24 * 23, 22 * 23, 22 * 23, 24 * 23}
								words := map[int]uint32{}
								for i, v := range quad {
									words[i*4] = v
								}
								sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("v%d/gen%v/load%d/tile%d/dir%d/var%d/point%v", version, generation, load, material, direction, variation, point), Flags: flags,
									Paint: legacy.PortTestPaintSpec{Seed: 43, Records: []legacy.PortTestMapRoomRecord{{Size: 32, Words: words}}, Globals: map[string]legacy.PortTestMapRoomArg{"section-magic-wall": {Value: 2}, "section-wall-load-flags": {Value: load}, "section-secret-index": {Value: 19}, "section-breakable-index": {Value: 23}}, Actions: []legacy.PortTestPaintAction{{Op: 0, Args: [6]legacy.PortTestMapRoomArg{{Slot: 1}}}}},
									IO:    []legacy.PortTestMapSectionIO{{Function: "walls", Read: true, Data: wire}},
								}
								cases = append(cases, sp)
								wants = append(wants, expected{20, 20, version, load, direction, material, variation, generation})
							}
						}
					}
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			want := wants[i]
			wire := r.IO[0]
			state := r.Paint.Steps[0]
			position := int64(len(wire.Data))
			if want.x == 255 {
				position = 19
			}
			if wire.Return != 1 || wire.Position != position {
				t.Fatalf("wall return/position %d/%d want1/%d", wire.Return, wire.Position, position)
			}
			if state.Globals["section-secret-index"] != 0 || state.Globals["section-breakable-index"] != 0 {
				t.Fatal("section counter reset")
			}
			if want.x == 255 {
				if len(wire.Scratch) != 0 || len(state.Walls.ByPos) != 0 {
					t.Fatal("end marker created wall")
				}
				return
			}
			var wall [9]uint32
			if want.generation {
				if len(wire.Scratch) != 1 || len(state.Walls.ByPos) != 0 {
					t.Fatal("scratch wall owner")
				}
				wall = wire.Scratch[0]
			} else {
				if len(wire.Scratch) != 0 || len(state.Walls.ByPos) != 1 {
					t.Fatal("live wall owner")
				}
				index := (state.Walls.ByPos[0][1] - 0x20000000) / 64
				wall = state.Walls.Records[index]
			}
			variation := want.variation
			if want.load != 0 && variation >= 4 {
				variation = 0
			}
			if byte(wall[0]) != want.direction || byte(wall[0]>>8) != want.material || byte(wall[0]>>16) != variation {
				t.Fatalf("wall material/direction/variation %x", wall)
			}
			if byte(wall[1]) != 128 || byte(wall[1]>>8) != byte(want.x) || byte(wall[1]>>16) != byte(want.y) || byte(wall[1]>>24) != 80 {
				t.Fatalf("wall flag/position/health %x", wall)
			}
			light, tag := uint32(173), uint32(127)
			if want.version == 6 {
				light = 1
				tag = 0
			}
			if byte(wall[2]) != byte(light) || wall[3] != tag {
				t.Fatalf("wall version fields %x", wall)
			}
		})
	}
	spellbookCapture(t, "map-sections-wall-region", out, "c48644df05641db6ab3f1e44c1942987bbeb977f5c82b1bbed8ff51a5cd0bad7")
}
