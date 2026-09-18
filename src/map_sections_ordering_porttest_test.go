//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapSectionsMetadataCountWidth(t *testing.T) {
	var cases []legacy.PortTestMapSectionSpec
	var counts []uint16
	for _, kind := range []string{"windows", "breakable", "secret"} {
		for _, count := range []uint16{0, 1, 32768, 65535} {
			version := uint16(2)
			if kind == "breakable" {
				version = 1
			}
			wire := mapSectionMetadataWire(kind, version, [][2]uint32{{10, 12}})
			binary.LittleEndian.PutUint16(wire[2:], count)
			cases = append(cases, legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("%s/count%d", kind, count), Flags: 0x400000, ScratchWalls: []legacy.PortTestPaintWall{{X: 10, Y: 12}}, Paint: legacy.PortTestPaintSpec{Seed: 73, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: kind, Read: true, Data: wire}}})
			counts = append(counts, count)
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			p := int64(len(r.IO[0].Data))
			if counts[i] == 0 {
				p = 4
			}
			if r.IO[0].Return != 1 || r.IO[0].Position != p {
				t.Fatalf("count consumption %d want%d", r.IO[0].Position, p)
			}
			s := r.Paint.Steps[0]
			for _, key := range []string{"section-breakable-index", "section-secret-index"} {
				want := uint32(0)
				if counts[i] != 0 && key == "section-"+cases[i].IO[0].Function+"-index" {
					want = 1
				}
				if s.Globals[key] != want {
					t.Fatalf("counter %s=%d want%d", key, s.Globals[key], want)
				}
			}
		})
	}
	spellbookCapture(t, "map-sections-metadata-count", out, "f92926a92ffb8c20b223a2eadc5f26071a07338a5a8baa7a23a55faefa38a6c0")
}

func TestMapSectionsFloorOrdering(t *testing.T) {
	nodes := [][4]uint32{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}}
	points := [][3]int{{2, 3, 1}, {10, 1, 2}, {1, 10, 3}}
	sp := legacy.PortTestMapSectionSpec{Name: "row-order-and-bound-width", Paint: legacy.PortTestPaintSpec{Seed: 79, Actions: []legacy.PortTestPaintAction{{Op: 0}}}, IO: []legacy.PortTestMapSectionIO{{Function: "floor"}}}
	for i, p := range points {
		w := map[int]uint32{0: uint32(p[2])}
		for half := 0; half < 2; half++ {
			for j, v := range nodes[i] {
				w[4*(1+5*half+j)] = v
			}
		}
		sp.Paint.Cells = append(sp.Paint.Cells, legacy.PortTestPaintCell{X: p[0], Y: p[1], Words: w})
	}
	// Upper flags bytes do not make a tile present for either bound scan.
	sp.Paint.Cells = append(sp.Paint.Cells, legacy.PortTestPaintCell{X: 0, Y: 5, Words: map[int]uint32{0: 0x100}})
	wire := []byte{4, 0}
	for _, v := range []uint32{1, 1, 10, 10} {
		wire = binary.LittleEndian.AppendUint32(wire, v)
	}
	for _, i := range []int{1, 0, 2} {
		p := points[i]
		coord := uint16(p[0]<<8 | p[1])
		if p[2]&1 != 0 {
			coord |= 0x8000
		}
		if p[2]&2 != 0 {
			coord |= 0x80
		}
		wire = binary.LittleEndian.AppendUint16(wire, coord)
		for half := 0; half < 2; half++ {
			if p[2]&(1<<half) != 0 {
				wire = append(wire, mapSectionTileBytes([][4]uint32{nodes[i]})...)
			}
		}
	}
	wire = append(wire, 255, 255)
	out := mapSectionsRun(t, []legacy.PortTestMapSectionSpec{sp})
	got := out[0].IO[0]
	if got.Return != 1 || !bytes.Equal(got.Data, wire) {
		t.Fatalf("floor ordering bytes %x want%x", got.Data, wire)
	}
	spellbookCapture(t, "map-sections-floor-ordering", out, "fcae64f6894e2611956fc0d708ae1cc0addbc0d73f83a0a474f7f2740ac38879")
}
