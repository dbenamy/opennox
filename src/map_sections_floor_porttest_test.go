//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func mapSectionFloorWire(x, y int, flags byte, a, b [4]uint32, emptyHeader bool) []byte {
	out := []byte{4, 0}
	if flags == 0 && !emptyHeader {
		return out
	}
	for _, n := range []uint32{uint32(x), uint32(y), 1, 1} {
		out = binary.LittleEndian.AppendUint32(out, n)
	}
	if flags&3 != 0 {
		coord := uint16(x<<8 | y)
		if flags&1 != 0 {
			coord |= 0x8000
		}
		if flags&2 != 0 {
			coord |= 0x80
		}
		out = binary.LittleEndian.AppendUint16(out, coord)
		if flags&1 != 0 {
			out = append(out, mapSectionTileBytes([][4]uint32{a})...)
		}
		if flags&2 != 0 {
			out = append(out, mapSectionTileBytes([][4]uint32{b})...)
		}
	}
	return append(out, 255, 255)
}
func TestMapSectionsFloorCurrentRecords(t *testing.T) {
	a, b := [4]uint32{0x123401, 0xffff8001, 0x180, 0x1ff}, [4]uint32{0x345622, 0x8000, 0x17f, 0x200}
	var cases []legacy.PortTestMapSectionSpec
	type expected struct {
		words    [11]uint32
		wire     []byte
		position int64
		x, y     int
	}
	var wants []expected
	for _, read := range []bool{false, true} {
		for _, pos := range [][2]int{{0, 0}, {1, 2}, {2, 1}, {63, 64}, {126, 127}, {127, 126}, {127, 127}} {
			for _, flags := range []byte{0, 1, 2, 3, 4, 0x81, 0x82, 0xff} {
				x, y := pos[0], pos[1]
				wire := mapSectionFloorWire(x, y, flags, a, b, read)
				words := [11]uint32{0x11223300 | uint32(flags)}
				copy(words[1:5], a[:])
				copy(words[6:10], b[:])
				if read {
					words = [11]uint32{0x112233a4, 0xfefdfcfb, 0xfefdfcfb, 0xfefdfcfb, 0xfefdfcfb, 0, 0xfefdfcfb, 0xfefdfcfb, 0xfefdfcfb, 0xfefdfcfb, 0}
				}
				input := map[int]uint32{}
				for i, v := range words {
					input[4*i] = v
				}
				want := expected{words: words, wire: wire, position: int64(len(wire)), x: x, y: y}
				present := flags & 3
				if read {
					want.words[0] &= 0xffffff00
					if x == 127 && y == 127 && present == 3 {
						present = 0
						want.position = 20
					}
					want.words[0] |= uint32(present)
				}
				if present&1 != 0 {
					n := mapSectionTileNarrow(a)
					copy(want.words[1:5], n[:])
				}
				if present&2 != 0 {
					n := mapSectionTileNarrow(b)
					copy(want.words[6:10], n[:])
				}
				sp := legacy.PortTestMapSectionSpec{Name: fmt.Sprintf("read%v/x%d/y%d/flags%x", read, x, y, flags), Paint: legacy.PortTestPaintSpec{Seed: 29,
					Cells: []legacy.PortTestPaintCell{{X: x, Y: y, Words: input}}, Actions: []legacy.PortTestPaintAction{{Op: 0}},
				}, IO: []legacy.PortTestMapSectionIO{{Function: "floor", Read: read, Data: wire}}}
				cases = append(cases, sp)
				wants = append(wants, want)
			}
		}
	}
	out := mapSectionsRun(t, cases)
	for i, r := range out {
		t.Run(r.Name, func(t *testing.T) {
			want := wants[i]
			got := r.IO[0]
			if got.Return != 1 || got.Position != want.position {
				t.Fatalf("return/position %d/%d want1/%d", got.Return, got.Position, want.position)
			}
			if !bytes.Equal(got.Data, want.wire) {
				t.Fatalf("floor bytes got%x want%x", got.Data, want.wire)
			}
			cells := r.Paint.Steps[0].Cells
			if len(cells) != 1 || cells[0].X != want.x || cells[0].Y != want.y || cells[0].Words != want.words {
				t.Fatalf("floor cells %+v want%x", cells, want.words)
			}
		})
	}
	spellbookCapture(t, "map-sections-floor-current", out, "ea601ffce1ae9335d4add46e23c0eec4da09da8a3de4362a9828ba2be8030902")
}
