//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func (q *quickbarOwner) prepareExpandedWindows(t *testing.T) {
	t.Helper()
	for row := 0; row < 4; row++ {
		bar := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1048196+uintptr(256*row)), 64)
		legacy.PortTestBookQuickbarInit(unsafe.Pointer(&bar[0]), 229, 378-60*row)
		parent := (*gui.Window)(unsafe.Pointer(uintptr(bar[52])))
		for i := 0; i < 5; i++ {
			w := q.c.GUI.NewWindowRaw(parent, gui.StatusFlags(8), 10+36*i, 1, 10, 10, nil)
			bar[58+i] = uint32(uintptr(w.C()))
		}
	}
	button := q.c.GUI.NewWindowRaw(nil, gui.StatusFlags(8), 500, 400, 10, 10, nil)
	*q.quickWords["dword_5d4594_1049512"] = uint32(uintptr(button.C()))
}
func TestQuickbarExpandedRows(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for selected := 0; selected < 5; selected++ {
		q.reset(t)
		q.prepareExpandedWindows(t)
		for i := 0; i < 25; i++ {
			q.bar[2*i] = uint32(i + 1)
			q.bar[2*i+1] = 0xabcdef00 + uint32(i&1)
		}
		q.bar[50] = uint32(selected)
		q.bar[51] = uint32(uintptr(unsafe.Pointer(&q.bar[10*selected])))
		q.call("sub_460920")
		q.check(t, byte(q.bar[50]) == 4 && memmap.Uint32(0x5D4594, 1049476) == 1, "expanded view selects bottom row")
		for row := 0; row < 4; row++ {
			for slot := 0; slot < 5; slot++ {
				p := memmap.PtrUint32(0x5D4594, 1048196+uintptr(256*row+8*slot))
				q.check(t, *p == uint32(row*5+slot+1), "expanded row copies spell id")
				q.check(t, *(*byte)(unsafe.Add(unsafe.Pointer(p), 4)) == byte(slot+row*5)&1, "expanded row copies flag byte")
			}
		}
		rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-expanded", selected), 0))
		// Change one expanded slot to prove closing writes it back.
		*memmap.PtrUint32(0x5D4594, 1048196+8) = 75
		*memmap.PtrUint8(0x5D4594, 1048196+12) = 0x80
		ret := q.call("nox_xxx_quickBarClose_4606B0")
		q.check(t, ret == 1 && memmap.Uint32(0x5D4594, 1049476) == 0, "close exits expanded view")
		q.check(t, q.bar[2] == 75 && q.bar[3] == 0xabcdef80, "close copies low flag byte without clobbering upper bytes")
		// Closing restores the row saved by the matching expansion.
		q.check(t, byte(q.bar[50]) == byte(selected), "closing restores the selected row")
		rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-closed", selected), ret))
	}
	spellbookCapture(t, "quickbar-expanded-rows", rows, "4a5c57ea1d6a2ee64eee652785ae96e93737efbba9b8fa0c21281ba8396ac323")
}

func TestQuickbarRowNavigation(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for selected := uint32(0); selected < 5; selected++ {
		for _, op := range []string{"nox_client_spellSetNext_4604F0", "nox_client_spellSetPrev_460540"} {
			for guard := 0; guard < 3; guard++ {
				q.reset(t)
				q.bar[50] = selected
				q.bar[51] = uint32(uintptr(unsafe.Pointer(&q.bar[10*selected])))
				if guard == 1 {
					*memmap.PtrUint32(0x5D4594, 1049476) = 1
				}
				if guard == 2 {
					*q.quickWords["dword_5d4594_1049496"] = 1
				}
				ret := q.call(op)
				want := selected
				if guard == 0 {
					if op == "nox_client_spellSetNext_4604F0" {
						want = (selected + 1) % 5
					} else {
						want = (selected + 4) % 5
					}
				}
				q.check(t, byte(q.bar[50]) == byte(want) && q.bar[51] == uint32(uintptr(unsafe.Pointer(&q.bar[10*want]))), "navigation row and pointer")
				rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-op%s-guard%d", selected, op, guard), ret))
			}
		}
	}
	for selected := uint32(0); selected < 3; selected++ {
		for _, op := range []string{"nox_client_trapSetNext_4603A0", "nox_client_trapSetPrev_4603F0"} {
			q.reset(t)
			q.call("nox_client_trapSetSelect_4604B0", selected)
			q.call(op)
			want := (selected + 1) % 3
			if op == "nox_client_trapSetPrev_4603F0" {
				want = (selected + 2) % 3
			}
			q.check(t, q.call("sub_4604E0") == want, "trap row wraps across three rows")
			q.check(t, memmap.Uint32(0x5D4594, 1048144) == uint32(uintptr(memmap.PtrOff(0x5D4594, 1047940+uintptr(40*want)))), "trap row pointer")
			rows = append(rows, q.snapshot(fmt.Sprintf("trap-selected%d-op%s", selected, op), 0))
		}
	}
	spellbookCapture(t, "quickbar-row-navigation", rows, "b11e94d59665bbf330ddd50e30839b4f301d4deddfdfacf073be9fb26d08d2f8")
}
