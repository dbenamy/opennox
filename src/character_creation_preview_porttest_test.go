//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestCharacterCreationPreview(t *testing.T) {
	o := newEntryOwner(t)
	colors := characterPalette(t)
	var controls []*gui.Window
	var ptrs []unsafe.Pointer
	for i := 0; i < 15; i++ {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 1, 1, nil)
		controls = append(controls, w)
		ptrs = append(ptrs, w.C())
		*characterWord(w, 32) = uint32(i%3) | uint32((i*3+2)%32)<<16
	}
	player, free := alloc.Calloc(1, 128)
	defer free()
	defer legacy.PortTestCharacterAppearanceOwner(player, ptrs)()
	var mods []unsafe.Pointer
	for k := 0; k < 3; k++ {
		p, done := alloc.Calloc(1, 64)
		defer done()
		b := unsafe.Slice((*byte)(p), 64)
		for i := 15; i < 33; i++ {
			b[i] = byte(17*k + 5*i)
		}
		binary.LittleEndian.PutUint32(b[36:], 3)
		binary.LittleEndian.PutUint32(b[40:], 2)
		binary.LittleEndian.PutUint32(b[44:], 4)
		mods = append(mods, p)
	}
	defer legacy.PortTestCharacterPreviewOwner(mods[0], mods[1], mods[2])()
	table := unsafe.Slice(memmap.PtrUint32(0x973A20, 16), 56)
	old := append([]uint32(nil), table...)
	defer copy(table, old)
	for i := range table {
		table[i] = uint32(uintptr(o.images[i%len(o.images)].C()))
	}
	window := o.c.GUI.NewWindowRaw(nil, 8, 10, 11, 40, 40, nil)
	type row struct {
		Sex, Mask int
		Position  image.Point
		Pixels    string
		Materials []byte
	}
	var rows []row
	for sex := 0; sex < 2; sex++ {
		for mask := 0; mask < 16; mask++ {
			for _, pos := range []image.Point{image.Pt(10, 11), image.Pt(-4, -3), image.Pt(91, 92)} {
				*(*byte)(unsafe.Add(player, 67)) = byte(sex)
				for i := 0; i < 4; i++ {
					controls[10+i].Flags = 0
					if mask&(1<<i) != 0 {
						controls[10+i].Flags = 8
					}
				}
				window.SetPos(pos)
				for i := range o.pix.Pix {
					o.pix.Pix[i] = 0x1234
				}
				if legacy.PortTestCharacterPreview(window.C()) != 1 {
					t.Fatal("preview return")
				}
				mat := unsafe.Slice((*byte)(unsafe.Add(o.c.r.Data().C(), 264)), 16*48)
				// Shirt is the final layer. Its two selected materials replace slots2 and4.
				for j, control := range []int{6, 7} {
					i := 32*(control%3) + (control*3+2)%32
					rgb := colors[3*i:][:3]
					word := uint32(rgb[0]/8)<<10 | uint32(rgb[1]/8)<<5 | uint32(rgb[2]/8)
					want := word | word<<16
					got := binary.LittleEndian.Uint32(mat[(2+2*j)*48+40:])
					if got != want {
						t.Fatalf("shirt material%d %x want%x", j, got, want)
					}
				}
				rows = append(rows, row{sex, mask, pos, effectsPixelHash(o.pix), append([]byte(nil), mat...)})
			}
		}
	}
	spellbookCapture(t, "character-creation-preview", rows, "7e4e4ff4483deb28183c7d6caf95ac408461cf0edbe31d4e90c1c1b5bbacba9c")
}
