//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSpellbookShowAndToggle(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
	var rows []spellbookResult
	for _, present := range []bool{false, true} {
		for _, known := range []uint32{0, 1} {
			for _, cursor := range []uint32{0, 1, 5} {
				for _, toggle := range []bool{false, true} {
					o.resetBook(t)
					if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
						t.Fatal("book setup")
					}
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = 1
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = known
					if !present {
						*o.words["dword_8531A0_2576"] = 0
					}
					*memmap.PtrUint32(0x5D4594, 1096672) = cursor
					op := "nox_xxx_bookShowMB_45AD70"
					if toggle {
						op = "nox_client_toggleSpellbook_45AC70"
					}
					ret := o.bookCall(op, 0)
					shown := cursor == 0 && (!present || known != 0)
					if (o.bookCall("sub_45CFC0") != 0) != shown {
						t.Fatalf("show player%v known%d cursor%d", present, known, cursor)
					}
					if cursor == 0 && present && known == 0 {
						slot := *o.txwords[11]
						if legacy.GoWStringP(memmap.PtrOff(0x5D4594, 823804+644*uintptr(slot))) != "No known entries" {
							t.Fatal("empty book feedback")
						}
					}
					label := fmt.Sprintf("player%v-known%d-cursor%d-toggle%v", present, known, cursor, toggle)
					rows = append(rows, o.bookSnapshot(label, ret))
					if shown {
						o.bookCall("nox_client_toggleSpellbook_45AC70")
						if o.bookCall("sub_45CFC0") != 0 {
							t.Fatal("toggle hides shown book")
						}
						rows = append(rows, o.bookSnapshot(label+"-hide", 0))
					}
				}
			}
		}
	}
	spellbookCapture(t, "show-toggle", rows, "c8ac97be08d31ab4c7480b33c16f166d8a0f6d1c111c18e31b8c86a546422e49")
}
