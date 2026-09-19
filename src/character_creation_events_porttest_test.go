//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestCharacterCreationPaletteEvents(t *testing.T) {
	o := newEntryOwner(t)
	root := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 640, 480, nil)
	controls := map[int]*gui.Window{}
	for _, id := range []int{711, 712, 713, 714, 720, 721, 722, 723, 724, 725, 726, 727, 728, 729, 731, 732, 733, 734, 999} {
		w := o.c.GUI.NewWindowRaw(root, 8, 0, 0, 1, 1, nil)
		w.SetID(uint(id))
		controls[id] = w
	}
	// The shipped palette is detached: lookup also traverses previous root windows.
	menu := o.c.GUI.NewWindowRaw(nil, 8, 20, 30, 175, 74, nil)
	menu.SetID(760)
	var swatches []*gui.Window
	for i := 0; i < 32; i++ {
		w := o.c.GUI.NewWindowRaw(menu, 8, 0, 0, 1, 1, nil)
		w.SetID(uint(761 + i))
		swatches = append(swatches, w)
	}
	defer legacy.PortTestCharacterPaletteOwner(root.C(), menu.C())()
	serverConfigOwnBytes(t, 0x5D4594, 1307788, 4)
	type row struct {
		Event, ID, Index, Return int
		Position                 image.Point
		Hidden                   bool
		Bank, Packed, Flags      uint32
	}
	var rows []row
	for id := 720; id <= 729; id++ {
		for _, index := range []int{0, 31} {
			menu.SetHidden(true)
			*characterWord(controls[id], 32) = 0xa5a5a5a5
			ret := legacy.PortTestCharacterColorEvent(root.C(), controls[id].C(), 16391, uint32(546)|uint32(88)<<16)
			bank := uint32(0)
			if id == 720 {
				bank = 2
			} else if id <= 724 {
				bank = 1
			}
			if ret != 1 || menu.GetFlags().IsHidden() || menu.Offs() != image.Pt(371, 51) || memmap.Uint32(0x5D4594, 1307788) != bank {
				t.Fatal("palette open", id, ret, menu.Offs())
			}
			for i, w := range swatches {
				if *characterWord(w, 32) != bank|uint32(i)<<16 {
					t.Fatal("event swatch", i)
				}
			}
			ret = legacy.PortTestCharacterColorEvent(menu.C(), swatches[index].C(), 16391, 0)
			packed := *characterWord(controls[id], 32)
			if ret != 1 || !menu.GetFlags().IsHidden() || packed != bank|uint32(index)<<16 {
				t.Fatal("palette selection", id, index, ret, packed)
			}
			rows = append(rows, row{16391, id, index, ret, menu.Offs(), true, bank, packed, 0})
		}
	}
	for id := 731; id <= 734; id++ {
		for _, enabled := range []bool{false, true} {
			w := controls[id-20]
			w.Flags = 0x100
			if enabled {
				w.Flags |= 8
			}
			ret := legacy.PortTestCharacterColorEvent(root.C(), controls[id].C(), 16391, 0)
			want := gui.StatusFlags(0x108)
			if enabled {
				want = 0x100
			}
			if ret != 1 || w.Flags != want {
				t.Fatal("override toggle", id, enabled)
			}
			rows = append(rows, row{Event: 16391, ID: id, Return: ret, Flags: uint32(w.Flags)})
		}
	}
	for _, event := range []int{0, 5, 16389, 16391} {
		ret := legacy.PortTestCharacterColorEvent(root.C(), controls[999].C(), event, 0)
		want := 0
		if event == 16389 || event == 16391 {
			want = 1
		}
		if ret != want {
			t.Fatal("other color event", event, ret)
		}
		rows = append(rows, row{Event: event, ID: 999, Return: ret})
	}
	spellbookCapture(t, "character-creation-palette-events", rows, "c9ab688f0a6a3b23c032ffdadc355057dd34c6b47066d73796e28eeceb07c9ee")
}
func TestCharacterCreationPaletteOutside(t *testing.T) {
	o := newEntryOwner(t)
	menu := o.c.GUI.NewWindowRaw(nil, 8, 20, 30, 175, 74, nil)
	defer legacy.PortTestCharacterPaletteOwner(nil, menu.C())()
	type row struct {
		Event  int
		Point  image.Point
		Return int
		Hidden bool
	}
	var rows []row
	for _, event := range []int{0, 4, 5, 6} {
		for _, p := range []image.Point{image.Pt(20, 30), image.Pt(195, 104), image.Pt(19, 30), image.Pt(20, 29), image.Pt(196, 104), image.Pt(195, 105), image.Pt(100, 70), image.Pt(65535, 65535)} {
			menu.SetHidden(false)
			ret := legacy.PortTestCharacterPaletteOutside(menu.C(), event, uint32(p.X)|uint32(p.Y)<<16)
			outside := p.X < 20 || p.X > 195 || p.Y < 30 || p.Y > 104
			want := 0
			if event == 5 && outside {
				want = 1
			}
			if ret != want || menu.GetFlags().IsHidden() != (want == 1) {
				t.Fatal("outside click", event, p, ret)
			}
			rows = append(rows, row{event, p, ret, menu.GetFlags().IsHidden()})
		}
	}
	spellbookCapture(t, "character-creation-palette-outside", rows, "2252743928e0d2856b28410e64af0f94732abe4f9569b75af1382410adc30166")
}
