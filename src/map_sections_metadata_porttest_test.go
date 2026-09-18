//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func mapSectionMetadataWire(kind string, version uint16, coords [][2]uint32) []byte {
	out := binary.LittleEndian.AppendUint16(nil, version)
	out = binary.LittleEndian.AppendUint16(out, uint16(len(coords)))
	for i, p := range coords {
		out = binary.LittleEndian.AppendUint32(out, p[0])
		out = binary.LittleEndian.AppendUint32(out, p[1])
		if kind == "secret" {
			out = binary.LittleEndian.AppendUint32(out, uint32(0x12345600+i))
			out = append(out, byte(i*8))
			if version >= 2 {
				out = append(out, 0, 19)
				out = binary.LittleEndian.AppendUint32(out, 0x87654321)
				out = binary.LittleEndian.AppendUint32(out, 0xfedcba98)
			}
		}
	}
	return out
}
func TestMapSectionsMetadataMissingScratchWall(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	for _, kind := range []string{"breakable", "secret"} {
		version := uint16(1)
		if kind == "secret" {
			version = 2
		}
		cases = append(cases, legacy.PortTestMapSectionSpec{Name: kind + "/present-then-missing", Flags: 0x400000,
			ScratchWalls: []legacy.PortTestPaintWall{{X: 10, Y: 12}}, Paint: legacy.PortTestPaintSpec{Seed: 31,
				Globals: map[string]legacy.PortTestMapRoomArg{"section-" + map[string]string{"breakable": "breakable-index", "secret": "secret-index"}[kind]: {Value: 17}},
				Actions: []legacy.PortTestPaintAction{{Op: 0}},
			}, IO: []legacy.PortTestMapSectionIO{{Function: kind, Read: true, Data: mapSectionMetadataWire(kind, version, [][2]uint32{{10, 12}, {11, 12}})}},
		})
	}
	region := [8]uint32{12 * 23, 8 * 23, 10 * 23, 10 * 23, 14 * 23, 10 * 23, 12 * 23, 12 * 23}
	words := map[int]uint32{}
	for i, v := range region {
		words[4*i] = v
	}
	cases = append(cases, legacy.PortTestMapSectionSpec{Name: "windows/missing-with-region", Flags: 0x400000,
		Paint: legacy.PortTestPaintSpec{Seed: 31, Records: []legacy.PortTestMapRoomRecord{{Size: 32, Words: words}},
			Globals: map[string]legacy.PortTestMapRoomArg{"section-map-min-x": {Value: 10}, "section-map-min-y": {Value: 8}},
			Actions: []legacy.PortTestPaintAction{{Op: 0, Args: [6]legacy.PortTestMapRoomArg{{Slot: 1}}}},
		}, IO: []legacy.PortTestMapSectionIO{{Function: "windows", Read: true, Data: mapSectionMetadataWire("windows", 2, [][2]uint32{{11, 12}})}},
	})
	cases = append(cases, legacy.PortTestMapSectionSpec{Name: "secret/missing-only", Flags: 0x400000,
		Paint: legacy.PortTestPaintSpec{Seed: 31, Actions: []legacy.PortTestPaintAction{{Op: 0}}},
		IO:    []legacy.PortTestMapSectionIO{{Function: "secret", Read: true, Data: mapSectionMetadataWire("secret", 2, [][2]uint32{{11, 12}})}},
	})
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			wire := r.IO[0]
			state := r.Paint.Steps[0]
			if wire.Return != 1 || wire.Position != int64(len(wire.Data)) {
				t.Fatal("metadata result/consumption", wire.Return, wire.Position)
			}
			switch i {
			case 0, 1:
				key := "section-breakable-index"
				if i == 1 {
					key = "section-secret-index"
				}
				if state.Globals[key] != 18 || len(wire.Scratch) != 1 || uint16(wire.Scratch[0][2]>>16) != 17 {
					t.Fatalf("missing entry reused previous wall: counter%d walls%v", state.Globals[key], wire.Scratch)
				}
			case 2:
				for _, p := range state.Records {
					if p.Kind == "input" && len(p.Words) == 8 {
						for j, v := range region {
							if p.Words[j] != v {
								t.Fatalf("missing wall mutated region word%d: %x want%x", j, p.Words[j], v)
							}
						}
					}
				}
			case 3:
				for _, p := range state.Records {
					if p.Kind == "section-allocation" && p.Alive && len(p.Words) == 8 {
						t.Fatal("missing wall retained unattached secret record")
					}
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-metadata-missing", out, "de4633140b0c35be625c3f8cde0bd7ea05bd548ac0332cc5668b32ef5f83030b")
}
