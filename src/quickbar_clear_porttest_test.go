//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"
)

func TestQuickbarClearAndRestore(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for selected := 0; selected < 5; selected++ {
		for _, op := range []string{"sub_4602F0", "sub_460380", "sub_461360", "sub_461400"} {
			q.reset(t)
			q.bar[50] = uint32(selected)
			q.bar[51] = uint32(uintptr(unsafe.Pointer(&q.bar[10*selected])))
			for i := 0; i < 25; i++ {
				q.bar[2*i] = uint32(i%3 + 1)
				q.bar[2*i+1] = 0x12345680 + uint32(i)
				*memmap.PtrUint32(0x5D4594, 1047564+uintptr(i*8)) = uint32(25 - i)
				*memmap.PtrUint32(0x5D4594, 1047568+uintptr(i*8)) = 0xfedcba00 + uint32(i)
			}
			for i := 0; i < 5; i++ {
				*memmap.PtrUint32(0x5D4594, 1047764+24+uintptr(24*i)+16) = 7
			}
			ret := q.call(op, 2)
			for i := 0; i < 25; i++ {
				id := uint32(i%3 + 1)
				flag := 0x12345680 + uint32(i)
				switch op {
				case "sub_4602F0", "sub_460380":
					id = 0
					flag &^= 255
				case "sub_461360":
					if id == 2 {
						id = 0
					}
				case "sub_461400":
					id = uint32(25 - i)
					flag = flag&0xffffff00 | uint32(i)
				}
				q.check(t, q.bar[2*i] == id && q.bar[2*i+1] == flag, "slot clearing/restoration preserves exact flag width")
			}
			for i := 0; i < 5; i++ {
				want := uint32(7)
				if op == "sub_460380" {
					want = 0
				}
				q.check(t, memmap.Uint32(0x5D4594, 1047764+24+uintptr(24*i)+16) == want, "only ability reset clears known ranks")
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("selected%d-op%s", selected, op), ret))
		}
	}
	for _, value := range []uint32{0, 1, 2, 0xffffffff} {
		q.reset(t)
		q.check(t, q.call("sub_461440", value) == value && q.call("sub_461450") == value, "restore-pending state keeps full word")
		rows = append(rows, q.snapshot(fmt.Sprintf("restore-pending%d", value), value))
	}
	spellbookCapture(t, "quickbar-clear-restore", rows, "cb7f298704397f2737a8eb5cbcea72dfdb026c003ee06ba7df995bb5a5ccc23a")
}
