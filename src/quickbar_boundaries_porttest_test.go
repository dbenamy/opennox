//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
)

func TestQuickbarBoundaryRows(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for _, selected := range []uint32{0, 4, 5, 6, 127, 254, 255} {
		for _, next := range []bool{false, true} {
			q.reset(t)
			// The current row remains valid even when the selected byte is out of range.
			q.bar[50] = 0xaabbcc00 | selected
			previous := q.bar[51]
			op := "nox_client_spellSetPrev_460540"
			target := int(selected) - 1
			if next {
				op = "nox_client_spellSetNext_4604F0"
				target = int(selected) + 1
				if target > 4 {
					target = 0
				}
			} else if target < 0 {
				target = 4
			}
			ret := q.call(op)
			if target >= 0 && target < 5 {
				q.check(t, q.bar[50] == 0xaabbcc00|uint32(target), "valid navigation changes only the selected byte")
				q.check(t, q.bar[51] == uint32(uintptr(unsafe.Pointer(&q.bar[10*target]))), "valid navigation updates its row address")
			} else {
				q.check(t, ret == uint32(target) && q.bar[50] == 0xaabbcc00|selected && q.bar[51] == previous, "invalid previous row is rejected without selecting row zero")
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-next%v", selected, next), ret))
		}
	}
	spellbookCapture(t, "quickbar-boundary-rows", rows, "676ce19d44b078ad596f7444eb122948710e2f068bdccb5a6c62e7e579c44b83")
}

func TestQuickbarBoundaryTray(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for class := byte(0); class < 3; class++ {
		for _, target := range []uint32{0, 479, 480, 481, 0x7fffffff, 0x80000000, 0xffffffff} {
			q.reset(t)
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = class
			*q.words["nox_win_height"] = 480
			q.call("nox_xxx_quickBarCreate_45E190")
			right := (*gui.Window)(unsafe.Pointer(uintptr(*q.quickWords["dword_5d4594_1049504"])))
			// Set the raw coordinate independently of GUI rectangle normalization.
			right.Off.Y = int(int32(target))
			*q.quickWords["dword_5d4594_1049536"] = target
			ret := q.call("nox_xxx_quickbarDrawFn_460000")
			want := target
			if target > 480 {
				want = 406
			}
			q.check(t, ret == 1 && *q.quickWords["dword_5d4594_1049536"] == want, "arrived tray compares target as unsigned before resetting its destination")
			q.check(t, uint32(right.Off.Y) == target, "arrival changes destination without moving the panel in that frame")
			rows = append(rows, q.snapshot(fmt.Sprintf("class%d-target%08x", class, target), ret))
		}
	}
	spellbookCapture(t, "quickbar-boundary-tray", rows, "043e2f970d63df24510c97ca416d07586524b782aac184ee5dcb27525cd0c42c")
}

func TestQuickbarBoundaryTrapRows(t *testing.T) {
	q := newQuickbarOwner(t)
	type result struct {
		Selected, CurrentOffset, Padding uint32
		Case                             string
		Sounds                           any
	}
	var rows []result
	for _, selected := range []byte{0, 2, 3, 5, 254, 255} {
		for _, next := range []bool{false, true} {
			q.reset(t)
			*memmap.PtrUint32(0x5D4594, 1048140) = 0xaabbcc00 | uint32(selected)
			op := "nox_client_trapSetPrev_4603F0"
			want := selected
			if next {
				op = "nox_client_trapSetNext_4603A0"
				if want == 2 {
					want = 0
				} else {
					want++
				}
			} else if want == 0 {
				want = 2
			} else {
				want--
			}
			q.call(op)
			word := memmap.Uint32(0x5D4594, 1048140)
			offset := memmap.Uint32(0x5D4594, 1048144) - uint32(uintptr(memmap.PtrOff(0x5D4594, 1047940)))
			q.check(t, word == 0xaabbcc00|uint32(want) && offset == 40*uint32(want), "trap navigation computes the byte-selected address without accessing the row")
			rows = append(rows, result{uint32(byte(word)), offset, word & 0xffffff00, fmt.Sprintf("selected%d-next%v", selected, next), q.sounds})
		}
	}
	spellbookCapture(t, "quickbar-boundary-trap-rows", rows, "9a6b72b2234e6c9dba02e66f99afff706eb61f2a4df182507b2d649f616fa3f6")
}
