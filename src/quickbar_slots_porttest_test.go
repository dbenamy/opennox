//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func TestQuickbarSlotSearch(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for selected := 0; selected < 5; selected++ {
		for empty := -1; empty < 25; empty++ {
			for guard := 0; guard < 3; guard++ {
				q.reset(t)
				q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
				for i := 0; i < 25; i++ {
					q.bar[2*i] = 1
					q.bar[2*i+1] = 0x76543201
				}
				if empty >= 0 {
					q.bar[2*empty] = 0
				}
				q.bar[50] = 0xabcdef00 | uint32(selected)
				q.bar[51] = uint32(uintptr(unsafe.Pointer(&q.bar[10*selected])))
				if guard == 1 {
					*memmap.PtrUint32(0x5D4594, 1049476) = 1
				}
				if guard == 2 {
					*q.quickWords["dword_5d4594_1049496"] = 1
				}
				before := append([]uint32(nil), q.bar[:50]...)
				local := q.call("sub_4612A0")
				wantLocal := uint32(0xffffffff)
				if empty >= 0 && empty/5 == selected {
					wantLocal = uint32(empty % 5)
				}
				q.check(t, local == wantLocal, "current-row search")
				got := q.call("nox_xxx_buttonFindFirstEmptySlot_461250")
				want := uint32(0xffffffff)
				if empty >= 0 {
					want = uint32(empty % 5)
				}
				q.check(t, got == want, "all-row search")
				wantRow := selected
				if empty >= 0 && guard == 0 {
					wantRow = empty / 5
				}
				q.check(t, byte(q.bar[50]) == byte(wantRow), "search row selection respects guards")
				q.check(t, reflect.DeepEqual(before, q.bar[:50]), "search never changes slots")
				q.check(t, q.call("nox_xxx_buttonHaveSpellInBarMB_4612D0", 1) == 1, "existing spell found")
				q.check(t, q.call("nox_xxx_buttonHaveSpellInBarMB_4612D0", 2) == 0, "missing spell absent")
				rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-empty%d-guard%d", selected, empty, guard), got))
			}
		}
	}
	spellbookCapture(t, "quickbar-slot-search", rows, "f67eaccbd3c97f7a0de9205929bd1a559f2a349718d4a9fcd9f63bbbbea236af")
}

func TestQuickbarSlotMutation(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for selected := 0; selected < 5; selected++ {
		for a := 0; a < 5; a++ {
			for b := 0; b < 5; b++ {
				q.reset(t)
				var defs []server.PortTestSpellClassDef
				for id := 1; id <= 25; id++ {
					defs = append(defs, server.PortTestSpellClassDef{Index: uint32(id), Flags: uint32(things.SpellClassAny), Valid: true})
				}
				q.configure(defs)
				for i := 0; i < 25; i++ {
					q.bar[2*i] = uint32(i + 1)
					q.bar[2*i+1] = 0xdead0000 + uint32(i)
				}
				q.bar[50] = uint32(selected)
				q.bar[51] = uint32(uintptr(unsafe.Pointer(&q.bar[10*selected])))
				want := append([]uint32(nil), q.bar[:50]...)
				ia, ib := selected*10+a*2, selected*10+b*2
				want[ia], want[ib] = want[ib], want[ia]
				want[ia+1], want[ib+1] = want[ib+1], want[ia+1]
				q.call("nox_xxx_clientSwapQuickbarKeys_45F300", uint32(uintptr(unsafe.Pointer(&q.bar[0]))), uint32(a), uint32(b))
				q.check(t, reflect.DeepEqual(want, q.bar[:50]), "swaps preserve full flag words")
				rows = append(rows, q.snapshot(fmt.Sprintf("row%d-swap%d-%d", selected, a, b), 0))
			}
		}
	}
	spellbookCapture(t, "quickbar-slot-mutation", rows, "0bc7ff3461e8544461840edae1ce82f71eee2facc995b966d3067c69f382046e")
}
