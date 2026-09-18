//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsWallCurrentWrite(t *testing.T) {
	newMapPolygonsOwner(t)
	var cases []legacy.PortTestMapSectionSpec
	var wants [][]byte
	for _, flags := range []byte{0, 128} {
		for _, direction := range []byte{0, 1, 2, 11} {
			for _, material := range []byte{0, 1, 2} {
				for _, variation := range []byte{0, 3, 4, 255} {
					for _, gameFlags := range []uint32{0, 0x200000} {
						sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("flags%x/dir%d/tile%d/var%d/game%x", flags, direction, material, variation, gameFlags), Flags: gameFlags,
							Paint: legacy.PortTestPaintSpec{Seed: 41, Globals: map[string]legacy.PortTestMapRoomArg{"section-magic-wall": {Value: 2}},
								Walls:   []legacy.PortTestPaintWall{{X: 8, Y: 10, Words: map[int]uint32{0: uint32(direction) | uint32(material)<<8 | uint32(variation)<<16 | 0xa5000000, 4: uint32(flags) | 8<<8 | 10<<16 | 77<<24, 8: 37 | 0x55<<8 | 17<<16, 12: 0xdeadbeef}}},
								Actions: []legacy.PortTestPaintAction{{Op: 0}},
							}, IO: []legacy.PortTestMapSectionIO{{Function: "walls"}},
						}
						wire := []byte{7, 0}
						for _, v := range []uint32{8, 10, 1, 1} {
							wire = binary.LittleEndian.AppendUint32(wire, v)
						}
						if material != 2 {
							tag := byte(0xef)
							if gameFlags != 0 {
								tag = 0
							}
							wire = append(wire, 8, 10, direction|(flags&128), material, variation, 100, tag)
						}
						wire = append(wire, 255)
						cases = append(cases, sp)
						wants = append(wants, wire)
					}
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			got := r.IO[0]
			if got.Return != 1 || got.Position != int64(len(wants[i])) {
				t.Fatalf("wall result/position %d/%d", got.Return, got.Position)
			}
			if !bytes.Equal(got.Data, wants[i]) {
				t.Fatalf("wall record bytes %x want%x", got.Data, wants[i])
			}
		})
	}
	spellbookCapture(t, "map-sections-wall-write", out, "5bdafe2aab5895f49ea57d9443df66bdc5f0345cdd06593337da74ebd0314bc0")
}
