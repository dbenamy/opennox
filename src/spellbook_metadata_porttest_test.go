//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookGuideAndAbilityLists(t *testing.T) {
	o := newSpellbookOwner(t)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	guides := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 740076)), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	titles := make([]string, 41)
	names := make([][]uint16, 41)
	for id := 1; id <= 40; id++ {
		titles[id] = fmt.Sprintf("%s %d", []string{"zeta", "Alpha", "ALPHA", "beta"}[id%4], id%3)
		names[id] = tooltipWide(t, []uint16(titlesToWords(titles[id])))
	}
	var rows []spellbookResult
	for _, guide := range []bool{false, true} {
		sizes := []int{0, 1, 2, int(server.AbilityMax) - 1}
		if guide {
			sizes = []int{0, 1, 2, 17, 18, 19, 35, 36, 37, 38, 39, 40}
		}
		for _, size := range sizes {
			for mode := 0; mode < 3; mode++ {
				for pattern := 0; pattern < 3; pattern++ {
					o.resetBook(t)
					if o.bookCall("nox_xxx_bookInit_45B9D0") != 1 {
						t.Fatal("book setup")
					}
					clear(guides)
					clear(o.c.srv.abilities.defs[:])
					if guide {
						*o.words["dword_5d4594_1046868"] = 1
					}
					if mode > 0 {
						noxflags.SetGame(noxflags.GameFlag(0x2000))
					}
					if mode == 2 {
						noxflags.SetGame(noxflags.GameModeQuest)
					}
					var want []int
					for id := 1; id <= size; id++ {
						valid := pattern != 2 || id%4 != 0
						known := pattern == 0 || id%3 != 0
						if guide {
							guides[7*id] = uint32(uintptr(unsafe.Pointer(&names[id][0])))
							if valid {
								guides[7*id+1] = 1
							}
						} else if valid {
							o.c.srv.abilities.defs[id] = AbilityDef{name: titles[id], field24: 1}
						}
						off := 3696 + 4*id
						if guide {
							off = 4244 + 4*id
						}
						if known {
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), off)) = uint32(id%5 + 1)
						}
						if valid && (known || mode == 1) && (!guide || mode == 2 || id < 37) {
							want = append(want, id)
						}
					}
					sort.SliceStable(want, func(i, j int) bool { return strings.ToLower(titles[want[i]]) < strings.ToLower(titles[want[j]]) })
					ret := o.bookCall("nox_xxx_guiSpellSortList_45ADF0", 0)
					got := make([]int, len(want))
					for i := range got {
						got[i] = int(memmap.Uint32(0x5D4594, 1046960+4*uintptr(i)))
					}
					if (ret != 0) != (len(want) > 0) || memmap.Uint32(0x5D4594, 1047508) != uint32(len(want)) || memmap.Uint32(0x5D4594, 1046940) != uint32(len(want)/18+1) || *o.words["dword_5d4594_1047512"] != 0 || (len(want) > 0 && !reflect.DeepEqual(got, want)) {
						t.Fatalf("guide=%v size=%d mode=%d pattern=%d: %v want %v", guide, size, mode, pattern, got, want)
					}
					rows = append(rows, o.bookSnapshot(fmt.Sprintf("guide%v-size%d-mode%d-pattern%d", guide, size, mode, pattern), ret))
				}
			}
		}
	}
	spellbookCapture(t, "metadata-lists", rows, "")
}

// ASCII titles deliberately include equal case-folded names while retaining
// the actual C UTF16 comparator and actual metadata owners.
func titlesToWords(s string) []uint16 {
	r := make([]uint16, len(s))
	for i := range s {
		r[i] = uint16(s[i])
	}
	return r
}
