//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerPanelsSpellToggleRestrictions(t *testing.T) {
	type row struct {
		Flags          uint32
		Mode           uint16
		ID             uint
		Checked        bool
		Mask           uint32
		AfterChecked   bool
		Dialog, Focus  int
		Title, Message string
		Dirty          uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 64} {
		for _, mode := range []uint16{0, 0x40} {
			for _, id := range []uint{1120, 1121} {
				for _, checked := range []bool{false, true} {
					t.Run(fmt.Sprintf("%x-%x-%d-%t", flags, mode, id, checked), func(t *testing.T) {
						o := newServerOptionsOwner(t)
						o.installSubpanels(t, true)
						t.Cleanup(legacy.PortTestServerOptionsRuleState())
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						binary.LittleEndian.PutUint16(o.settings[52:], mode)
						oldSpells := o.c.srv.Spells
						t.Cleanup(func() { o.c.srv.Spells = oldSpells })
						s := server.PortTestSpellClassServer([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x1000000, Valid: true}, {Index: 132, Flags: 0x2000000, Valid: true}})
						o.c.srv.Spells = s.Spells
						o.c.srv.Spells.DefByInd(spell.ID(1)).Title = "Ordinary spell"
						o.c.srv.Spells.DefByInd(spell.ID(132)).Title = "Restricted spell"
						configure, restore := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "spelllst.c:Notice", Vals: []strman.Variant{{Str: "Notice"}}}, strman.Entry{ID: "plyrspel.c:Illegal", Vals: []strman.Variant{{Str: "Unavailable in this mode"}}})
						t.Cleanup(restore)
						configure(0)
						oldStrings := strMan
						strMan = o.c.srv.Strings()
						t.Cleanup(func() { strMan = oldStrings })
						raw := legacy.PortTestServerPanelsConstruct("spell", o.options, unsafe.Pointer(&o.settings[0]))
						w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
						if w == nil {
							t.Fatal("spell restriction root")
						}
						child := w.ChildByID(id)
						child.DrawData().Field0 &^= 4
						if checked {
							child.DrawData().Field0 |= 4
						}
						index := 1
						if id == 1121 {
							index = 132
						}
						masks := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045488), 5)
						for i := range masks {
							masks[i] = 0x55555555
						}
						oldDialog, oldFocus := legacy.Nox_xxx_dialogMsgBoxCreate_449A10, legacy.Sub_44A360
						t.Cleanup(func() { legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = oldDialog; legacy.Sub_44A360 = oldFocus })
						r := row{Flags: flags, Mode: mode, ID: id, Checked: checked}
						legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(parent *gui.Window, title, text string, df gui.DialogFlags, ok, cancel func()) {
							if parent != w || ok != nil || cancel != nil {
								t.Error("spell dialog owner")
							}
							r.Dialog = int(df)
							r.Title = title
							r.Message = text
						}
						legacy.Sub_44A360 = func(value int) { r.Focus = value }
						*o.optionWords["dirty"] = 0
						if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(child.C()), 0))) != 0 {
							t.Fatal("spell toggle result")
						}
						restricted := index == 132 && (flags&64 != 0 || mode&0x40 != 0)
						want := uint32(0x55555555)
						wantChecked := checked
						wantDirty := uint32(1)
						bit := uint32(1) << uint(index%32)
						if restricted {
							wantChecked = !checked
							wantDirty = 0
							if r.Dialog != 33 || r.Focus != 1 || r.Title != "Notice" || r.Message != "Unavailable in this mode" {
								t.Fatalf("restriction dialog %+v", r)
							}
						} else {
							if checked {
								want &^= bit
							} else {
								want |= bit
							}
							if r.Dialog != 0 || r.Focus != 0 {
								t.Fatal("unexpected restriction")
							}
						}
						if masks[index/32] != want || (child.DrawData().Field0&4 != 0) != wantChecked || *o.optionWords["dirty"] != wantDirty {
							t.Fatal("spell toggle state")
						}
						for i, v := range masks {
							if i != index/32 && v != 0x55555555 {
								t.Fatal("spell toggle unrelated word")
							}
						}
						r.Mask = masks[index/32]
						r.AfterChecked = child.DrawData().Field0&4 != 0
						r.Dirty = *o.optionWords["dirty"]
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-panels-spell-toggle-restrictions", rows, "5fdda6fb43e55a415191cec31c5e4171148670ecd1144343a348b705657d0a99")
}
