//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookSpellLists(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	var rows []spellbookResult
	for _, size := range []int{0, 1, 2, 17, 18, 19, 35, 36, 37, 38, 136} {
		for _, class := range []int{1, 2} {
			for mode := 0; mode < 3; mode++ {
				for special := 0; special < 3; special++ {
					o.resetBook(t)
					if got := o.bookCall("nox_xxx_bookInit_45B9D0"); got != 1 {
						t.Fatal("book setup")
					}
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
					if mode > 0 {
						noxflags.SetGame(noxflags.GameFlag(0x2000))
					}
					if mode == 2 {
						noxflags.SetGame(noxflags.GameModeQuest)
					}
					var defs []server.PortTestSpellClassDef
					var eligible []int
					titles := map[int]string{}
					specialCount := 0
					for id := 1; id <= size; id++ {
						flags := uint32(things.SpellClassAny)
						isSpecial := special != 0 && id > size-3
						if isSpecial {
							flags |= []uint32{0x1000, 0x4000, 0x10000}[id%3]
						}
						hidden := special == 2 && id%11 == 0 && !isSpecial
						if hidden {
							flags |= 0x2000
						}
						defs = append(defs, server.PortTestSpellClassDef{Index: uint32(id), Flags: flags, Valid: true})
						known := id%3 != 0
						if known {
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696+4*id)) = uint32(id%5 + 1)
						}
						titles[id] = fmt.Sprintf("%s %03d", []string{"zeta", "ALPHA", "beta", "Alpha"}[id%4], id%9)
						if id == 34 || (!known && mode != 1) {
							continue
						}
						if isSpecial {
							specialCount++
						}
						if !hidden {
							eligible = append(eligible, id)
						}
					}
					configure(defs)
					for id, title := range titles {
						o.c.srv.Spells.DefByInd(spell.ID(id)).Title = title
					}
					// The canonical layout keeps a special-entry suffix and sorts the normal
					// prefix. A full normal page still reserves a following page in the C UI.
					normal := len(eligible) - specialCount
					sort.SliceStable(eligible[:normal], func(i, j int) bool {
						return strings.ToLower(titles[eligible[i]]) < strings.ToLower(titles[eligible[j]])
					})
					got := o.bookCall("nox_xxx_guiSpellSortList_45ADF0", uint32(class))
					if (got != 0) != (len(eligible) != 0) {
						t.Fatal("book list availability")
					}
					if *o.words["dword_5d4594_1046656"] != 13 {
						t.Fatal("actual font line height")
					}
					if memmap.Uint32(0x5D4594, 1047508) != uint32(len(eligible)) || *o.words["dword_5d4594_1047512"] != uint32(specialCount) || memmap.Uint32(0x5D4594, 1046940) != uint32(normal/18+1) {
						t.Fatal("book list counts")
					}
					actual := make([]int, len(eligible))
					for i := range actual {
						actual[i] = int(memmap.Uint32(0x5D4594, 1046960+4*uintptr(i)))
					}
					if !reflect.DeepEqual(actual, eligible) && len(eligible) != 0 {
						t.Fatalf("book sort %v != %v", actual, eligible)
					}
					rows = append(rows, o.bookSnapshot(fmt.Sprintf("size%d-class%d-mode%d-special%d", size, class, mode, special), got))
				}
			}
		}
	}
	spellbookCapture(t, "spell-lists", rows, "")
}
