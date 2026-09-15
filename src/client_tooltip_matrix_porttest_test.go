//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTooltipEquipmentMatrix(t *testing.T) {
	o := newTooltipOwner(t)
	var out []tooltipResult
	id := 0
	for _, lang := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		o.language(lang)
		for _, class := range []uint32{0x1000, 0x1000000, 0x2000000, 0x10000000, 0x13001000} {
			for _, sub := range []uint32{0, 0x800000, 0x7800000} {
				for mask := uint32(0); mask < 16; mask++ {
					o.prepare(class, sub, 0, mask)
					out = append(out, o.invoke(t, id, 0, o.dr))
					out = append(out, o.invoke(t, id, 1, o.dr))
					id++
				}
			}
		}
	}
	effectsCapture(t, "tooltip-equipment", out, len(out), "f8a3ade90b2a2ad4ac5d5754eaeda3af50ee3a4257fa337c6a2a0fc2bb78a805")
}
func TestClientTooltipBooksMatrix(t *testing.T) {
	o := newTooltipOwner(t)
	var out []tooltipResult
	id := 0
	for _, lang := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		o.language(lang)
		for sub := uint32(0); sub < 8; sub++ {
			for _, code := range []uint32{0, 1, 0x7fff, 0x8000, 0xffffffff} {
				for _, class := range []uint32{0x100, 0x400100, 0x20000100} {
					sentinel := uint32(0)
					if sub&1 != 0 {
						sentinel = 137
					} else if sub&2 != 0 {
						sentinel = 41
					} else if sub&4 != 0 {
						sentinel = 6
					}
					for _, metadata := range []uint32{0, 1, 2, sentinel} {
						o.prepare(class, sub, metadata, 0)
						o.dr.NetCode32 = code
						out = append(out, o.invoke(t, id, 0, o.dr))
						out = append(out, o.invoke(t, id, 1, o.dr))
						id++
					}
				}
			}
		}
	}
	effectsCapture(t, "tooltip-books", out, len(out), "b76c1936b9bc0f46fbbc899561aceca7a8ed1654745bde1ec350e6544f410830")
}
func TestClientTooltipNamesMatrix(t *testing.T) {
	o := newTooltipOwner(t)
	var out []tooltipResult
	id := 0
	for _, lang := range []int{0, 2, 3, 5, 6} {
		o.language(lang)
		for _, kind := range []int{-1, 0, 7, 8} {
			typ := o.c.Things.TypeByInd(4)
			if kind < 0 {
				typ.PrettyName = nil
			} else {
				typ.PrettyName = &o.text[kind][0]
			}
			for _, class := range []uint32{0, 0x100, 0x200, 0x80000000} {
				o.prepare(class, 0, 0, 0)
				out = append(out, o.invoke(t, id, 0, o.dr))
				out = append(out, o.invoke(t, id, 1, nil))
				id++
			}
		}
		for _, class := range []uint32{0x1000, 0x2000000} {
			for _, missing := range []bool{false, true} {
				o.prepare(class, 0, 0, 15)
				if missing {
					o.weapon.TypeInd, o.armor.TypeInd = 5, 5
				}
				out = append(out, o.invoke(t, id, 0, o.dr))
				o.weapon.TypeInd, o.armor.TypeInd = 4, 4
				id++
			}
		}
		// Preserve raw UTF16 and distinguish absent modifier, nil descriptor,
		// and present-but-empty descriptor. All pointers refer to owned storage.
		for slot := 0; slot < 4; slot++ {
			offset := uintptr(8)
			if slot == 3 {
				offset = 12
			}
			p := (*unsafe.Pointer)(unsafe.Add(o.mods[slot].C(), offset))
			old := *p
			for _, kind := range []int{-1, 7, 8} {
				if kind < 0 {
					*p = nil
				} else {
					*p = unsafe.Pointer(&o.text[kind][0])
				}
				o.prepare(0x1000, 0, 0, 15)
				out = append(out, o.invoke(t, id, 0, o.dr))
				id++
			}
			*p = old
		}
	}
	effectsCapture(t, "tooltip-names", out, len(out), "1395d70fe2dbc419f762fb9b64a80ade504af0e78f63d327c23c062b61bbfcf7")
}
func TestClientTooltipCursorMatrix(t *testing.T) {
	o := newTooltipOwner(t)
	type result struct {
		Length, Step  int
		Cursor, Input []uint16
	}
	var out []result
	for _, n := range []int{0, 1, 2, 63, 254, 255, 256, 257, 511, 1024} {
		text := make([]uint16, n)
		for i := range text {
			text[i] = []uint16{'a', 0xd800, 'Ω', 0xdc01}[i%4]
		}
		p := tooltipWide(t, text)
		o.prepare(0, 0, 0, 0)
		for step := 0; step < 3; step++ {
			input := &p[0]
			if step == 1 {
				input = nil
			}
			if step == 2 {
				input = &p[len(p)-1]
			}
			legacy.PortTestTooltipCursor(input)
			out = append(out, result{n, step, append([]uint16(nil), o.cursor...), append([]uint16(nil), p...)})
		}
	}
	effectsCapture(t, "tooltip-cursor", out, len(out), "2d1dcab558f44301d78a937ed84117ea6e32063daff79b197963da26bfcf8b3f")
}
