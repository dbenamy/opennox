//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookPageRendering(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	guides := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 740076)), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	title := tooltipWide(t, titlesToWords("Creature name that wraps over several lines"))
	desc := tooltipWide(t, titlesToWords("A creature description with enough text to wrap across several lines on the right page."))
	var rows []spellbookResult
	for class := 0; class < 3; class++ {
		for _, guide := range []uint32{0, 1} {
			for _, detail := range []uint32{0, 1} {
				for _, length := range []int{0, 1, 3} {
					o.resetBook(t)
					if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
						t.Fatal("book setup")
					}
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
					o.bookWindow().SetPos(image.Pt(40, 100))
					o.bookCall("nox_xxx_bookSetColor_45AC40")
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | 0x12c, Valid: true}})
					sd := o.c.srv.Spells.DefByInd(spell.ID(1))
					sd.Title = "Spell name that wraps over several lines"
					sd.Desc = "A spell description with enough text to wrap across several lines on the right page."
					sd.Icon = unsafe.Pointer(o.images[1])
					sd.IconEnabled = unsafe.Pointer(o.images[2])
					o.c.srv.abilities.defs[1] = AbilityDef{name: "Ability name that wraps over several lines", desc: "An ability description that wraps across several lines on the right page.", field24: 1, icon8: o.images[1], icon12: o.images[2]}
					clear(guides)
					guides[7] = uint32(uintptr(unsafe.Pointer(&title[0])))
					guides[8] = 1
					guides[9] = uint32(uintptr(unsafe.Pointer(&desc[0])))
					guides[11] = uint32(uintptr(o.images[1].C()))
					guides[13] = []uint32{1, 2, 4}[class]
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = uint32(length + 1)
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4232)) = 1
					*o.words["dword_5d4594_1046868"] = guide
					*o.words["dword_5d4594_1046872"] = guide
					*o.words["nox_xxx_aNox_cfg_0_587000_132132"] = 1 - detail
					*memmap.PtrUint32(0x5D4594, 1047508) = uint32(length)
					for i := 0; i < length; i++ {
						*memmap.PtrUint32(0x5D4594, 1046960+4*uintptr(i)) = 1
					}
					if detail != 0 && length == 0 {
						continue
					} // An empty book has no detail entry.
					blank := effectsPixelHash(o.pix)
					ret := o.bookCall("nox_xxx_bookDrawList_45BD40", *o.words["nox_win_unk1"])
					if ret != 1 || effectsPixelHash(o.pix) == blank {
						t.Fatalf("book must draw visible pixels: class=%d guide=%d detail=%d length=%d ret=%d", class, guide, detail, length, ret)
					}
					if length > 0 && len(o.displayText) == 0 {
						t.Fatal("nonempty page must draw text")
					}
					rows = append(rows, o.bookSnapshot(fmt.Sprintf("class%d-guide%d-detail%d-count%d", class, guide, detail, length), ret))
				}
			}
		}
	}
	spellbookCapture(t, "render", rows, "e09d57ec73188ed7cb8c4e81bcbd736b5da2be8268b9df78362241a41b9cd37a")
}
