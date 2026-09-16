//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSpellbookIconsAndTooltips(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	o.c.srv.abilities.defs[1] = AbilityDef{name: "Berserk", field24: 1, icon8: o.images[1]}
	var rows []spellbookResult
	for class := 0; class < 3; class++ {
		for view := uint32(0); view < 4; view++ {
			for mode := 0; mode < 4; mode++ {
				for _, known := range []uint32{0, 1} {
					for _, special := range []uint32{0, 0x1000} {
						o.resetBook(t)
						if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
							t.Fatal("book setup")
						}
						configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | special, Valid: true}, {Index: 75, Flags: uint32(things.SpellClassAny), Valid: true}})
						for _, id := range []spell.ID{1, 75} {
							o.c.srv.Spells.DefByInd(id).Icon = unsafe.Pointer(o.images[1])
						}
						*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
						*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4232)) = known
						if mode&1 != 0 {
							noxflags.SetGame(noxflags.GameFlag(0x2000))
						}
						if mode&2 != 0 {
							noxflags.SetGame(noxflags.GameModeQuest)
						}
						*o.words["dword_5d4594_1046868"] = view
						*memmap.PtrUint32(0x5D4594, 1046960) = 1
						*memmap.PtrUint32(0x5D4594, 1047508) = 1
						icon := *o.words["dword_5d4594_1046952"]
						blank := effectsPixelHash(o.pix)
						ret := o.bookCall("nox_xxx_bookDrawIconFn_45CB30", icon)
						draw := view == 0 || (view == 1 && class == 2 && (known != 0 || mode == 1))
						if ret != 1 || (effectsPixelHash(o.pix) != blank) != draw {
							t.Fatalf("icon draw class%d view%d mode%d known%d", class, view, mode, known)
						}
						label := fmt.Sprintf("class%d-view%d-mode%d-known%d-special%d", class, view, mode, known, special)
						rows = append(rows, o.bookSnapshot(label+"-draw", ret))
						ret = o.bookCall("nox_xxx_bookWndFn_45CC10", icon, 5, 0)
						denied := view == 1 && (class != 2 || (mode == 0 && known == 0))
						drag := !denied && (view == 0 && (class == 0 || special == 0) || view == 1 && class == 2)
						wantID := uint32(1)
						if view == 1 {
							wantID = 75
						}
						if denied && ret != 0 || !denied && ret != 1 {
							t.Fatal("icon press eligibility")
						}
						if drag {
							if o.c.dragndrapSpell != wantID || o.c.dragndropSpellType != 1 || o.c.GUI.Captured() == nil {
								t.Fatal("icon drag setup")
							}
						} else if o.c.dragndropSpellType != 0 || o.c.GUI.Captured() != nil {
							t.Fatal("unavailable icon must not start drag")
						}
						rows = append(rows, o.bookSnapshot(label+"-press", ret))
					}
				}
			}
		}
	}
	for class := 0; class < 3; class++ {
		o.resetBook(t)
		o.bookCall("nox_xxx_bookInit_45B9D0")
		*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
		o.collect()
		for _, w := range o.windows {
			if w.ID() != 1310 && w.ID() != 1320 {
				continue
			}
			ret := o.bookCall("nox_xxx_book_45CF00", uint32(uintptr(w.C())))
			want := "Spells"
			if w.ID() == 1320 {
				want = "Creatures"
			} else if class == 0 {
				want = "Abilities"
			}
			r := o.bookSnapshot(fmt.Sprintf("class%d-tooltip%d", class, w.ID()), ret)
			if ret != 1 || r.Tooltip != want {
				t.Fatalf("book tooltip got%q want%q", r.Tooltip, want)
			}
			rows = append(rows, r)
		}
	}
	spellbookCapture(t, "icons-tooltips", rows, "4620f33300afa75b363e009048b010db88a549fef9b68d385d51ed473f09c1b8")
}
