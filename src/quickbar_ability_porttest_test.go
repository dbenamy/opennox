//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
)

func TestQuickbarAbilityState(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for _, id := range []uint32{0, 1, 2, 3, 4, 5, 6, 31} {
		for _, enabled := range []uint32{0, 1, 2, 0xffffffff} {
			q.reset(t)
			for n := 1; n <= 5; n++ {
				*memmap.PtrUint32(0x5D4594, 1047764+uintptr(24*n)+12) = 0xa5a5a5a5
				*memmap.PtrUint32(0x5D4594, 1047764+uintptr(24*n)+20) = 77
			}
			before := append([]uint32(nil), q.raw...)
			ret := q.call("sub_461120", id, enabled)
			for n := 1; n <= 5; n++ {
				off := 1047764 + 24*n + 12
				want := uint32(0xa5a5a5a5)
				if uint32(n) == id {
					if enabled != 0 {
						want |= 1 << id
					} else {
						want &^= 1 << id
					}
				}
				q.check(t, memmap.Uint32(0x5D4594, uintptr(off)) == want, "ability flag changes only requested bit")
			}
			available := q.call("sub_461160", id)
			if id >= 1 && id <= 5 {
				q.check(t, available == uint32Bool(enabled != 0), "ability availability matches changed flag")
			} else {
				q.check(t, available == 0, "unknown ability is unavailable")
			}
			// Only matching record flag words may differ.
			for i, v := range q.raw {
				off := 1047548 + 4*i
				if id >= 1 && id <= 5 && off == 1047764+24*int(id)+12 {
					continue
				}
				q.check(t, v == before[i], "ability flag update preserves other state")
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("flag-id%d-enabled%d", id, enabled), ret))
			ret = q.call("sub_461090", id, enabled)
			for n := 1; n <= 5; n++ {
				off := uintptr(1047764 + 24*n)
				if uint32(n) == id {
					q.check(t, memmap.Uint32(0x5D4594, off+8) == enabled, "ability active state")
					if enabled != 0 {
						q.check(t, memmap.Uint32(0x5D4594, off+20) == 0, "enabled ability clears cooldown timestamp")
					}
				} else {
					q.check(t, memmap.Uint32(0x5D4594, off+20) == 77, "other ability timestamps unchanged")
				}
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("state-id%d-enabled%d", id, enabled), ret))
		}
	}
	for id := uint32(1); id <= 6; id++ {
		q.reset(t)
		ret := q.call("sub_4610D0", id)
		for n := uint32(1); n <= 5; n++ {
			want := uint32(0)
			if n == id || id == 6 {
				want = 1
			}
			q.check(t, memmap.Uint32(0x5D4594, 1047764+uintptr(24*n)+8) == want, "single/all ability reset")
		}
		rows = append(rows, q.snapshot(fmt.Sprintf("reset-id%d", id), ret))
	}
	spellbookCapture(t, "quickbar-ability-state", rows, "d4e2e826f1bf4948aaea4ea4554199d98ee69773c3525ae2634e4147426211e7")
}

func uint32Bool(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
