//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"
)

func TestQuickbarTimedSelection(t *testing.T) {
	q := newQuickbarOwner(t)
	frame := q.c.srv.Frame()
	t.Cleanup(func() { q.c.srv.SetFrame(frame) })
	var rows []quickbarResult
	for _, gap := range []uint32{0, 1, q.c.srv.TickRate() / 2, q.c.srv.TickRate()/2 + 1} {
		q.reset(t)
		q.c.srv.SetFrame(1000)
		for step := 0; step < 8; step++ {
			if step > 0 {
				q.c.srv.SetFrame(q.c.srv.Frame() + gap)
			}
			q.call("nox_client_spellSetSelect_460590")
			want := []uint32{0, 1, 2, 3, 4, 4, 4, 4}[step]
			if gap == q.c.srv.TickRate()/2 {
				// Saturation does not refresh the timeout: after two half-rate
				// gaps the next press starts a new selection sequence.
				want = []uint32{0, 1, 2, 3, 4, 4, 0, 1}[step]
			}
			if gap > q.c.srv.TickRate()/2 {
				want = 0
			}
			q.check(t, byte(q.bar[50]) == byte(want), "timed selection resets after strict timeout and saturates at last row")
			q.check(t, q.bar[51] == uint32(uintptr(unsafe.Pointer(&q.bar[10*want]))), "reachable saturated selection keeps matching row pointer")
			rows = append(rows, q.snapshot(fmt.Sprintf("gap%d-step%d", gap, step), 0))
		}
	}
	// Boundaries of the byte-sized spell highlight timer store.
	for _, index := range []uint32{0xffffffff, 0, 1, 135, 136, 139, 140} {
		for _, value := range []uint32{0, 1, 127, 128, 255} {
			q.reset(t)
			q.call("sub_460EB0", index, value)
			for i := 0; i < 140; i++ {
				want := byte(0)
				if uint32(i) == index {
					want = byte(value)
				}
				q.check(t, *memmap.PtrUint8(0x5D4594, 1049544+uintptr(i)) == want, "bounded timer store")
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("timer-id%d-value%d", index, value), 0))
		}
	}
	spellbookCapture(t, "quickbar-timing", rows, "c9a7f9c54b4e72e0c0fe8a5437a05d36b6a25f52ad97616118cb8cee12cd06a6")
}
