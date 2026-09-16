//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookSpellKnowledge(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	var rows []spellbookResult
	for _, class := range []uint32{0, 1, 2} {
		for family := 0; family < 4; family++ {
			for _, rank := range []uint32{0, 1, 2, 5, 0xffffffff} {
				for _, present := range []bool{false, true} {
					o.resetBook(t)
					if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
						t.Fatal("book setup")
					}
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
					parent := []uint32{0, 0x1000, 0x4000, 0x10000}[family]
					child := []uint32{0, 0x2000, 0x8000, 0x20000}[family]
					defs := []server.PortTestSpellClassDef{
						{Index: 1, Flags: uint32(things.SpellClassAny) | parent, Valid: true},
						{Index: 2, Flags: uint32(things.SpellClassAny) | child, Valid: true},
						{Index: 3, Flags: uint32(things.SpellClassAny) | child, Valid: false},
						{Index: 4, Flags: uint32(things.SpellClassAny), Valid: true},
					}
					configure(defs)
					levels := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696)), 137)
					for i := 1; i <= 4; i++ {
						levels[i] = 7
					}
					if !present {
						*o.words["dword_8531A0_2576"] = 0
					}
					ret := o.bookCall("nox_xxx_netSpellRewardCli_45CFE0", 1, rank, 0, 0)
					want1, want2 := uint32(7), uint32(7)
					if present && class != 0 {
						want1 = rank
						if family != 0 {
							want2 = rank
						}
					}
					if levels[1] != want1 || levels[2] != want2 || levels[3] != 7 || levels[4] != 7 {
						t.Fatalf("knowledge class%d family%d rank%d player%v: %v", class, family, rank, present, levels[:5])
					}
					rows = append(rows, o.bookSnapshot(fmt.Sprintf("class%d-family%d-rank%d-player%v", class, family, rank, present), ret))
				}
			}
		}
	}
	spellbookCapture(t, "spell-knowledge", rows, "c03cf5fc9649ff34021a66f171c18a90ffea75afb568d222963621fb9b70a2e3")
}
