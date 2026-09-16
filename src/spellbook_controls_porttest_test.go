//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestSpellbookPositionAndControls(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	o.c.srv.abilities.defs[1] = AbilityDef{name: "Berserk", field24: 1}
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	guides[8] = 1
	point, freePoint := alloc.New([2]int32{})
	defer freePoint()
	var rows []spellbookResult
	for _, kind := range []uint32{1, 2, 3, 4} {
		for _, id := range []uint32{1, 999} {
			for _, pos := range []image.Point{image.Pt(5, 160), image.Pt(-7, 53)} {
				prepareSpellbookAddition(t, o, bar, 640)
				class := byte(1)
				if kind == 3 {
					class = 0
				} else if kind == 4 {
					class = 2
					*o.words["dword_5d4594_1046868"] = 1
					*o.words["dword_5d4594_1046872"] = 1
				}
				*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = class
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = 1
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4248)) = 1
				*point = [2]int32{int32(pos.X), int32(pos.Y)}
				ret := o.bookCall("nox_xxx_bookSetForward_45D200", kind, id, uint32(uintptr(unsafe.Pointer(point))))
				if o.bookWindow().Off != pos {
					t.Fatal("explicit book position")
				}
				if kind >= 2 && id == 1 && (*o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 0 || *o.words["dword_5d4594_1046932"] != 0) {
					t.Fatal("explicit known entry selection")
				}
				rows = append(rows, o.bookSnapshot(fmt.Sprintf("kind%d-id%d-pos%v", kind, id, pos), ret))
				*point = [2]int32{99, 100}
				ret = o.bookCall("sub_45D550", uint32(uintptr(unsafe.Pointer(point))))
				icon := (*gui.Window)(unsafe.Pointer(uintptr(*o.words["dword_5d4594_1046952"])))
				actual := icon.GlobalPos()
				if ret != 0 || *point != [2]int32{int32(actual.X), int32(actual.Y)} {
					t.Fatal("book icon position")
				}
				if o.bookCall("sub_45D550", 0) != 0 || o.bookCall("nox_xxx_book_45BD30", 0, 0) != 1 {
					t.Fatal("nullable position and empty draw callback")
				}
			}
		}
	}
	for _, active := range []uint32{0, 1, 2} {
		for _, shown := range []bool{false, true} {
			for _, show := range []uint32{0, 1} {
				prepareSpellbookAddition(t, o, bar, 640)
				o.bookWindow().SetHidden(!shown)
				*o.words["dword_5d4594_1046864"] = active
				ret := o.bookCall("sub_45D500", show)
				want := shown
				if active != 0 {
					want = show != 0
				}
				if (!o.bookWindow().GetFlags().IsHidden()) != want {
					t.Fatal("temporary book visibility")
				}
				rows = append(rows, o.bookSnapshot(fmt.Sprintf("active%d-shown%v-show%d", active, shown, show), ret))
			}
		}
	}
	for _, active := range []uint32{0, 1, 2} {
		prepareSpellbookAddition(t, o, bar, 640)
		*o.words["dword_5d4594_1047520"] = active
		if o.bookCall("sub_45D9B0") != active {
			t.Fatal("addition state getter")
		}
		o.bookCall("sub_45D810")
		if o.bookCall("sub_45D9B0") != 0 {
			t.Fatal("addition stop clears state")
		}
		rows = append(rows, o.bookSnapshot(fmt.Sprintf("stop%d", active), 0))
	}
	spellbookCapture(t, "position-controls", rows, "84c6d01b3ff8723634c1f29cfc82fc7625ad49e57f10f3c106c6780c7419309f")
}
