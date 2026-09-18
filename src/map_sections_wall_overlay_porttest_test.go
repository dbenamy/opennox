//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsWallOverlay(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	var wants []byte
	composition := blobdata.PortTestMapPaintingTables()[71276]
	for _, version := range []uint16{1, 3, 7} {
		for _, merge := range []uint32{0, 1} {
			for old := byte(0); old < 13; old++ {
				for incoming := byte(0); incoming < 13; incoming++ {
					wire := binary.LittleEndian.AppendUint16(nil, version)
					for _, v := range []uint32{8, 10, 0, 0} {
						wire = binary.LittleEndian.AppendUint32(wire, v)
					}
					if version == 7 {
						wire = append(wire, 8, 10)
					}
					wire = append(wire, incoming)
					if version >= 3 {
						wire = append(wire, 1, 3)
					}
					if version == 7 {
						wire = append(wire, 173, 127, 255)
					}
					sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("v%d/merge%d/old%d/new%d", version, merge, old, incoming), Paint: legacy.PortTestPaintSpec{Seed: 71, Globals: map[string]legacy.PortTestMapRoomArg{"section-magic-wall": {Value: 2}, "section-wall-load-flags": {Value: merge}}, Walls: []legacy.PortTestPaintWall{{X: 8, Y: 10, Words: map[int]uint32{0: uint32(old), 4: 0x000a0880}}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: "walls", Read: true, Data: wire}}}
					want := incoming
					if merge != 0 {
						want = composition[13*int(old)+int(incoming)]
					}
					cases = append(cases, sp)
					wants = append(wants, want)
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			w := r.IO[0]
			s := r.Paint.Steps[0]
			if w.Return != 1 || w.Position != int64(len(w.Data)) || len(s.Walls.ByPos) != 1 {
				t.Fatal("overlay return/consumption/owner")
			}
			p := s.Walls.Records[(s.Walls.ByPos[0][1]-0x20000000)/64]
			if byte(p[0]) != wants[i] || byte(p[1]) != 128 {
				t.Fatalf("overlay direction/retained flag %x want%d", p, wants[i])
			}
		})
	}
	spellbookCapture(t, "map-sections-wall-overlay", out, "34cdfe9f57636c288b91b6c487afd91fd43a4503d35b6284b0f7013abb8acada")
}
