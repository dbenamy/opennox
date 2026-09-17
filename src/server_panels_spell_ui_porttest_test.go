//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerPanelsSpellList(t *testing.T) {
	type row struct {
		Flags                    uint32
		Empty                    bool
		Names, Classes           []string
		Before, Cleared, Enabled [5]uint32
		Hidden, Checked          []bool
		Dirty                    uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 0x4001, 0x8001} {
		for _, empty := range []bool{false, true} {
			t.Run(fmt.Sprintf("%x-%t", flags, empty), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				o.installSubpanels(t, true)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				t.Cleanup(legacy.PortTestServerOptionsRuleState())
				oldSpells := o.c.srv.Spells
				t.Cleanup(func() { o.c.srv.Spells = oldSpells })
				var defs []server.PortTestSpellClassDef
				if !empty {
					combinations := []uint32{0, 0x1000000, 0x2000000, 0x4000000, 0x6000000, 0x7000000, 0x1001000, 0x2004000, 0x4010000, 0x80}
					for i := 1; i < 61; i++ {
						defs = append(defs, server.PortTestSpellClassDef{Index: uint32(i), Flags: combinations[(i-1)%len(combinations)], Valid: i%13 != 0})
					}
				}
				s := server.PortTestSpellClassServer(defs)
				o.c.srv.Spells = s.Spells
				var wantNames, wantClasses []string
				var shownIDs []uint32
				for _, d := range defs {
					sp := o.c.srv.Spells.DefByInd(spell.ID(d.Index))
					sp.Title = fmt.Sprintf("Spell %02d é", d.Index)
					sp.Enabled = d.Index%3 == 0
					if !d.Valid || d.Flags&0x15000 != 0 || d.Flags&0x7000000 == 0 {
						continue
					}
					wantNames = append(wantNames, sp.Title)
					shownIDs = append(shownIDs, d.Index)
					class := "Wizard"
					if d.Flags&0x1000000 != 0 || d.Flags&0x6000000 == 0x6000000 {
						class = "Common"
					} else if d.Flags&0x4000000 != 0 {
						class = "Conjurer"
					}
					wantClasses = append(wantClasses, class)
				}
				language, restore := o.c.srv.PortTestMeterStrings(
					strman.Entry{ID: "spelllst.c:Common", Vals: []strman.Variant{{Str: "Common"}}},
					strman.Entry{ID: "spelllst.c:SpellWizard", Vals: []strman.Variant{{Str: "Wizard"}}},
					strman.Entry{ID: "spelllst.c:SpellConjurer", Vals: []strman.Variant{{Str: "Conjurer"}}},
					strman.Entry{ID: "WindowDir:Blank", Vals: []strman.Variant{{Str: ""}}},
				)
				t.Cleanup(restore)
				language(0)
				oldStrings := strMan
				strMan = o.c.srv.Strings()
				t.Cleanup(func() { strMan = oldStrings })
				raw := legacy.PortTestServerPanelsConstruct("spell", o.options, unsafe.Pointer(&o.settings[0]))
				if raw == 0 {
					t.Fatal("spell constructor")
				}
				w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
				names := (*gui.ScrollListBoxData)(w.ChildByID(1110).WidgetData)
				classes := (*gui.ScrollListBoxData)(w.ChildByID(1112).WidgetData)
				var gotNames, gotClasses []string
				for _, it := range unsafe.Slice(names.Items, int(names.Field_11_0)) {
					gotNames = append(gotNames, alloc.GoString16(&it.Text[0]))
				}
				for _, it := range unsafe.Slice(classes.Items, int(classes.Field_11_0)) {
					gotClasses = append(gotClasses, alloc.GoString16(&it.Text[0]))
				}
				if !reflect.DeepEqual(gotNames, wantNames) || !reflect.DeepEqual(gotClasses, wantClasses) {
					t.Fatalf("population names=%v classes=%v want classes=%v", gotNames, gotClasses, wantClasses)
				}
				r := row{Flags: flags, Empty: empty, Names: gotNames, Classes: gotClasses}
				masks := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045488), 5)
				copy(r.Before[:], masks)
				for i := 0; i < 14; i++ {
					child := w.ChildByID(uint(1120 + i))
					hidden := child.Flags.IsHidden()
					checked := child.DrawData().Field0&4 != 0
					wantChecked := false
					if i < len(shownIDs) {
						wantChecked = shownIDs[i]%3 == 0
					}
					if hidden != (i >= len(shownIDs)) || checked != wantChecked {
						t.Fatalf("row %d hidden/checked %t/%t", i, hidden, checked)
					}
					r.Hidden = append(r.Hidden, hidden)
					r.Checked = append(r.Checked, checked)
				}
				enabled := flags&1 != 0 && flags&0xc000 == 0
				if w.ChildByID(1115).Flags.IsEnabled() != enabled {
					t.Fatal("host controls")
				}
				// Dispatch both bulk actions through the installed production callback.
				for _, id := range []uint{1116, 1115} {
					*o.optionWords["dirty"] = 0
					child := w.ChildByID(id)
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(unsafe.Pointer(child)), 0))) != 0 {
						t.Fatal("bulk return")
					}
					for _, index := range shownIDs {
						if (masks[index/32]&(1<<uint(index%32)) != 0) != (id == 1115) {
							t.Fatalf("bulk %d spell %d", id, index)
						}
					}
					if *o.optionWords["dirty"] != 1 {
						t.Fatal("bulk dirty")
					}
					if id == 1116 {
						copy(r.Cleared[:], masks)
					} else {
						copy(r.Enabled[:], masks)
					}
				}

				if len(shownIDs) > 14 {
					down := w.ChildByID(1114)
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16384, uintptr(down.C()), 0))) != 0 {
						t.Fatal("scroll result")
					}
					if names.Field_13_1 == 0 || names.Field_13_1 != classes.Field_13_1 {
						t.Fatal("parallel spell list scroll")
					}
					child := w.ChildByID(1120)
					child.DrawData().Field0 |= 4
					index := shownIDs[1]
					before := masks[index/32]
					if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(child.C()), 0))) != 0 || masks[index/32] != before&^(1<<uint(index%32)) {
						t.Fatal("scrolled spell row toggle")
					}
					up := w.ChildByID(1113)
					w.Func94(gui.AsWindowEvent(16391, uintptr(up.C()), 0))
					if names.Field_13_1 != 0 || classes.Field_13_1 != 0 {
						t.Fatal("spell scroll up")
					}
				}
				// Empty or unmatched rows mark the panel dirty without changing a mask.
				if empty {
					before := append([]uint32(nil), masks...)
					w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(1120).C()), 0))
					if !reflect.DeepEqual(before, masks) {
						t.Fatal("empty spell row changed masks")
					}
				}
				r.Dirty = *o.optionWords["dirty"]
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-panels-spell-list", rows, "1afa85092d476ad6e29a469993de7bd4f24b9ad71cbfee05ab318075efdc3fb2")
}
