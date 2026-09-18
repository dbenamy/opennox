//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsFloorHistoricalRead(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	type expected struct {
		x, y  int
		flags byte
		value uint32
		count int
	}
	var wants []expected
	for _, point := range [][2]int{{0, 0}, {1, 2}, {127, 127}} {
		for _, flags := range []byte{0, 1, 2, 3, 4, 0x81, 0xff} {
			for _, value := range []uint32{0, 0x800080ff, 0xffffffff} {
				for _, count := range []int{0, 1, 11} {
					wire := []byte{3, 0}
					for _, n := range []uint32{uint32(point[0]), uint32(point[1]), 1, 1} {
						wire = binary.LittleEndian.AppendUint32(wire, n)
					}
					wire = append(wire, flags)
					for half := 0; half < 2; half++ {
						if flags&(1<<half) == 0 {
							continue
						}
						for node := 0; node <= count; node++ {
							for field := 0; field < 4; field++ {
								wire = binary.LittleEndian.AppendUint32(wire, value+uint32(half*17+node*31+field))
							}
							if node == 0 {
								wire = append(wire, byte(count))
							}
						}
					}
					cases = append(cases, legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("point%v/flags%x/value%x/count%d", point, flags, value, count), Paint: legacy.PortTestPaintSpec{Seed: 53, Cells: []legacy.PortTestPaintCell{{X: point[0], Y: point[1], Words: map[int]uint32{0: 0x11223300}}}, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: "floor", Read: true, Data: wire}}})
					wants = append(wants, expected{point[0], point[1], flags, value, count})
				}
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			w := wants[i]
			got := r.IO[0]
			if got.Return != 1 || got.Position != int64(len(got.Data)) {
				t.Fatalf("historical floor return/position %d/%d", got.Return, got.Position)
			}
			cells := r.Paint.Steps[0].Cells
			if len(cells) != 1 {
				t.Fatalf("cells %d", len(cells))
			}
			cell := cells[0]
			if cell.X != w.x || cell.Y != w.y || cell.Words[0] != 0x11223300|uint32(w.flags) {
				t.Fatalf("cell flags/position %+v", cell)
			}
			for half := 0; half < 2; half++ {
				base := 1 + half*5
				if w.flags&(1<<half) == 0 {
					continue
				}
				for field := 0; field < 4; field++ {
					if cell.Words[base+field] != w.value+uint32(half*17+field) {
						t.Fatalf("historical full-width scalar %x", cell.Words)
					}
				}
				if (cell.Words[base+4] == 0) != (w.count == 0) {
					t.Fatal("historical chain ownership")
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-floor-historical", out, "85af9fb637c0096b22da58c9d79bc7d8c81a9a5fa427bfa19942e89df8125415")
}
