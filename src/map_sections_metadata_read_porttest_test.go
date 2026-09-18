//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsMetadataRead(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	type expected struct {
		kind         string
		version      uint16
		generation   bool
		state, flags byte
		start        uint32
	}
	var wants []expected
	for _, kind := range []string{"windows", "breakable", "secret"} {
		versions := []uint16{0, 1, 2}
		if kind == "breakable" {
			versions = []uint16{0, 1}
		}
		for _, version := range versions {
			for _, generation := range []bool{false, true} {
				for _, state := range []byte{0, 1, 3, 255} {
					for _, flag := range []byte{0, 8} {
						for _, start := range []uint32{0, 65535, 0xffffffff} {
							wire := mapSectionMetadataWire(kind, version, [][2]uint32{{10, 12}})
							if kind == "secret" {
								wire[16] = flag
								if version >= 2 {
									wire[17] = state
									wire[18] = 19
								}
							}
							sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("%s/v%d/gen%v/state%d/flag%x/start%x", kind, version, generation, state, flag, start), Paint: legacy.PortTestPaintSpec{Seed: 61, Globals: map[string]legacy.PortTestMapRoomArg{"section-breakable-index": {Value: start}, "section-secret-index": {Value: start}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: kind, Read: true, Data: wire}}}
							walls := []legacy.PortTestPaintWall{{X: 10, Y: 12, Words: map[int]uint32{0: 0x00070302, 4: 0x000c0a80}}}
							if generation {
								sp.Flags = 0x400000
								sp.ScratchWalls = walls
							} else {
								sp.Paint.Walls = walls
							}
							cases = append(cases, sp)
							wants = append(wants, expected{kind, version, generation, state, flag, start})
						}
					}
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			w := wants[i]
			got := r.IO[0]
			snap := r.Paint.Steps[0]
			if got.Return != 1 || got.Position != int64(len(got.Data)) {
				t.Fatal("metadata read result/position")
			}
			var wall [9]uint32
			if w.generation {
				wall = got.Scratch[0]
			} else {
				wall = snap.Walls.Records[(snap.Walls.ByPos[0][1]-0x20000000)/64]
			}
			bit := byte(0x40)
			if w.kind == "breakable" {
				bit = 8
			}
			if w.kind == "secret" {
				bit = 4
			}
			if byte(wall[1]) != 128|bit {
				t.Fatalf("metadata flags%x", wall[1])
			}
			if w.kind == "windows" {
				variation := byte(7)
				if w.version < 2 {
					variation = 0
				}
				if byte(wall[0]>>16) != variation {
					t.Fatal("window version variation")
				}
				return
			}
			if uint16(wall[2]>>16) != uint16(w.start) || snap.Globals["section-"+w.kind+"-index"] != w.start+1 {
				t.Fatal("metadata index width")
			}
			if w.kind == "breakable" && !w.generation && len(got.Breakable) != 1 {
				t.Fatal("breakable list owner")
			}
			if w.kind == "secret" {
				var data []uint32
				for _, p := range snap.Records {
					if p.ID == wall[7] {
						data = p.Words
					}
				}
				if len(data) != 8 {
					t.Fatalf("secret owner %x", wall[7])
				}
				state, variation, end := w.state, byte(19), uint32(0xfedcba98)
				if w.version < 2 {
					state = 0
					variation = 0
					end = 0
				}
				if state == 0 {
					state = 1
					variation = 0
					end = 0
					if w.flags&8 != 0 {
						state = 3
						variation = 23
						end = 0xffffffff
					}
				}
				if data[4] != 0x12345600 || byte(data[5]) != w.flags || byte(data[5]>>8) != state || byte(data[5]>>16) != variation || data[7] != end {
					t.Fatalf("secret defaults %x", data)
				}
				if binary.LittleEndian.Uint32(got.Data[4:]) != 10 {
					t.Fatal("source coordinate changed")
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-metadata-read", out, "a72fdec825a58414df696a32304bf3d048c22c3fa1b04398dc78eac343ff5f71")
}
