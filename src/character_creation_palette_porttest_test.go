//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func characterPalette(t *testing.T) []byte {
	t.Helper()
	p := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1307796), 288)
	old := append([]byte(nil), p...)
	t.Cleanup(func() { copy(p, old) })
	for i := 0; i < 96; i++ {
		p[3*i] = byte(i + 1)
		p[3*i+1] = byte(97 + i)
		p[3*i+2] = byte(201 - i)
	}
	return p
}
func characterWord(w *gui.Window, off int) *uint32 { return (*uint32)(unsafe.Add(w.C(), off)) }

func TestCharacterCreationPaletteMatch(t *testing.T) {
	o := newEntryOwner(t)
	colors := characterPalette(t)
	root := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 90, 90, nil)
	root.SetID(700)
	w := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 10, 10, nil)
	w.SetID(721)
	enabled := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 10, 10, nil)
	enabled.SetID(711)
	toggle := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 10, 10, nil)
	toggle.SetID(731)
	t.Cleanup(legacy.PortTestCharacterPaletteOwner(root.C(), nil))
	p, free := alloc.Calloc(1, 3)
	defer free()
	rgb := unsafe.Slice((*byte)(p), 3)
	type row struct {
		Bank, Index             int
		Packed, Enabled, Toggle uint32
		ReturnKind              string
	}
	var rows []row
	for bank := 0; bank < 3; bank++ {
		for index := -1; index < 32; index++ {
			copy(rgb, []byte{0, 0, 0})
			if index >= 0 {
				copy(rgb, colors[96*bank+3*index:][:3])
			}
			*characterWord(w, 32) = 0xa5a5a5a5
			enabled.Flags = 8
			*characterWord(toggle, 36) = 0x100
			ret := legacy.PortTestCharacterPaletteMatch(w.C(), bank, p)
			wantIndex := index
			if index < 0 {
				wantIndex = 32
				if bank == 1 {
					wantIndex = 9
				}
			}
			want := uint32(bank) | uint32(wantIndex)<<16
			if got := *characterWord(w, 32); got != want {
				t.Fatalf("bank%d index%d: %x want%x", bank, index, got, want)
			}
			wantEnabled, wantToggle := uint32(8), uint32(0x100)
			if bank == 1 {
				if index >= 0 {
					wantToggle |= 6
				} else {
					wantEnabled = 0
				}
			}
			if uint32(enabled.Flags) != wantEnabled || *characterWord(toggle, 36) != wantToggle {
				t.Fatalf("bank%d index%d visibility/toggle", bank, index)
			}
			kind := "scalar"
			if ret == uintptr(toggle.C()) {
				kind = "toggle"
			} else if ret >= uintptr(unsafe.Pointer(&colors[0])) && ret <= uintptr(unsafe.Pointer(&colors[0]))+uintptr(len(colors))+1 {
				kind = "palette"
			}
			rows = append(rows, row{bank, index, want, uint32(enabled.Flags), *characterWord(toggle, 36), kind})
		}
	}
	// First matching palette entry wins even with duplicates later in the bank.
	copy(colors[96+3*31:][:3], colors[96:][:3])
	copy(rgb, colors[96:][:3])
	legacy.PortTestCharacterPaletteMatch(w.C(), 1, p)
	if *characterWord(w, 32) != 1 {
		t.Fatal("duplicate must select first")
	}
	spellbookCapture(t, "character-creation-palette-match", rows, "177ecd3c53ebaf6dfc108c08e410c84cc31a797e5dfde810e8169bf3e47a04fa")
}

func TestCharacterCreationPaletteMenu(t *testing.T) {
	o := newEntryOwner(t)
	root := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 90, 90, nil)
	root.SetID(700)
	menu := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 32, 32, nil)
	menu.SetID(760)
	var swatches []*gui.Window
	for i := 0; i < 32; i++ {
		w := o.c.GUI.NewWindowRaw(menu, 8, i%8, i/8, 1, 1, nil)
		w.SetID(uint(761 + i))
		swatches = append(swatches, w)
	}
	target := o.c.GUI.NewWindowRaw(menu, 8, 0, 0, 1, 1, nil)
	target.SetID(721)
	t.Cleanup(legacy.PortTestCharacterPaletteOwner(root.C(), menu.C()))
	bankPtr := memmap.PtrUint32(0x5D4594, 1307788)
	old := *bankPtr
	t.Cleanup(func() { *bankPtr = old })
	type row struct {
		Bank, Index            uint16
		Words                  []uint32
		Target                 uint32
		Hidden, ReturnedTarget bool
	}
	var rows []row
	for _, bank := range []uint16{0, 1, 2, 0xffff} {
		ret := legacy.PortTestCharacterPaletteFill(bank)
		if ret != uintptr(swatches[31].C()) || *bankPtr != uint32(bank) {
			t.Fatal("fill return/bank")
		}
		words := make([]uint32, 32)
		for i, w := range swatches {
			words[i] = *characterWord(w, 32)
			if words[i] != uint32(bank)|uint32(i)<<16 {
				t.Fatal("packed swatch", i)
			}
		}
		for _, index := range []uint16{0, 1, 31, 32, 0xdead, 0xffff} {
			*characterWord(target, 32) = 0xa5a5a5a5
			menu.SetHidden(false)
			ret := legacy.PortTestCharacterPaletteClose(721, index)
			want := uint32(0xa5a5a5a5)
			if index < 32 {
				want = uint32(bank) | uint32(index)<<16
			}
			if *characterWord(target, 32) != want || !menu.GetFlags().IsHidden() {
				t.Fatal("selection", bank, index)
			}
			rows = append(rows, row{bank, index, words, want, menu.GetFlags().IsHidden(), ret == uintptr(target.C())})
		}
	}
	// Unknown target leaves the palette cells intact.
	before := *characterWord(swatches[0], 32)
	if legacy.PortTestCharacterPaletteClose(9999, 0) != 0 || *characterWord(swatches[0], 32) != before {
		t.Fatal("missing selection target")
	}
	spellbookCapture(t, "character-creation-palette-menu", rows, "12bc7a25ca85dcaad6e3dad893a3821811288ff5c6216833f0c5f137e6772225")
}

func TestCharacterCreationSwatchPixels(t *testing.T) {
	o := newEntryOwner(t)
	colors := characterPalette(t)
	w := o.c.GUI.NewWindowRaw(nil, 8, 4, 5, 7, 9, nil)
	type row struct {
		Bank, Index int
		Pixels      string
	}
	var rows []row
	for bank := 0; bank < 3; bank++ {
		for index := 0; index < 32; index++ {
			*characterWord(w, 32) = uint32(bank) | uint32(index)<<16
			for i := range o.pix.Pix {
				o.pix.Pix[i] = 0xffff
			}
			if legacy.PortTestCharacterPaletteDraw(w.C()) != 1 {
				t.Fatal("draw return")
			}
			h := effectsPixelHash(o.pix)
			rgb := colors[3*(32*bank+index):][:3]
			color := uint16(rgb[0]/8)<<10 | uint16(rgb[1]/8)<<5 | uint16(rgb[2]/8)
			for y := 0; y < o.pix.Rect.Dy(); y++ {
				for x := 0; x < o.pix.Rect.Dx(); x++ {
					want := uint16(0xffff)
					if x >= 4 && x < 11 && y >= 5 && y < 14 {
						want = color
					}
					if got := o.pix.Pix[y*o.pix.Stride+x]; got != want {
						t.Fatalf("bank%d index%d pixel(%d,%d)=%x want%x", bank, index, x, y, got, want)
					}
				}
			}
			rows = append(rows, row{bank, index, h})
		}
	}
	spellbookCapture(t, "character-creation-swatch-pixels", rows, "453728fb5ee8e7b81d76343558f5d26cd1e8e6aa313053af1341dbb81c5c6519")
}
