//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestChatBubbleOverlapAndShift(t *testing.T) {
	o := newChatBubbleRenderOwner(t)
	a := o.create(t, 1, "A", 1, 0, 1)
	b := o.create(t, 2, "B", 1, 0, 1)
	for _, p := range []unsafe.Pointer{a, b} {
		objectXferSetWord(p, 648, 200)
		objectXferSetWord(p, 652, 200)
		objectXferSetWord(p, 672, 40)
		objectXferSetWord(p, 676, 13)
	}
	type row struct {
		Name     string
		Overlaps int32
		Shift    [2]int32
	}
	var rows []row
	for _, dx := range []int{-63, -62, -61, 0, 61, 62, 63} {
		for _, dy := range []int{-36, -35, -34, 0, 34, 35, 36} {
			name := fmt.Sprintf("overlap/%d,%d", dx, dy)
			t.Run(name, func(t *testing.T) {
				objectXferSetWord(b, 648, uint32(200+dx))
				objectXferSetWord(b, 652, uint32(200+dy))
				v := legacy.PortTestChatBubbleOverlap(a, b)
				want := int32(0)
				if dx > -62 && dx < 62 && dy > -35 && dy < 35 {
					want = 1
				}
				if v != want {
					t.Fatalf("overlap %d want %d", v, want)
				}
				rows = append(rows, row{Name: name, Overlaps: v})
			})
		}
	}
	objectXferSetWord(a, 652, 180)
	objectXferSetWord(b, 648, 300)
	objectXferSetWord(b, 652, 220)
	objectXferSetWord(b, 672, 50)
	objectXferSetWord(b, 676, 20)
	expected := [][2]int32{{238, 185}, {372, 185}, {238, 262}, {372, 262}, {200, 185}, {200, 262}, {238, 180}, {372, 180}}
	for i, want := range expected {
		mask := byte(1 << i)
		t.Run(fmt.Sprintf("shift/%02x", mask), func(t *testing.T) {
			got := legacy.PortTestChatBubbleShift(mask, a, b)
			if got != want {
				t.Fatalf("shift %v want %v", got, want)
			}
			rows = append(rows, row{Name: fmt.Sprintf("shift/%02x", mask), Shift: got})
		})
	}
	spellbookCapture(t, "chat-bubble-overlap-shift", rows, "9aa60d424a458b934e9ce2fd87617c1fe115b0565ee4ae0fb8fb796b099a979d")
}
func TestChatBubbleRegions(t *testing.T) {
	o := newChatBubbleRenderOwner(t)
	_ = o
	type row struct{ X, Y, Region int32 }
	var rows []row
	xs := [][2]int32{{-1, 16}, {0, 16}, {212, 16}, {213, 32}, {425, 32}, {426, 64}, {639, 64}, {640, 64}}
	ys := [][2]int32{{-1, 1}, {0, 1}, {159, 1}, {160, 2}, {319, 2}, {320, 4}, {479, 4}, {480, 4}}
	for _, x := range xs {
		for _, y := range ys {
			t.Run(fmt.Sprintf("%d,%d", x[0], y[0]), func(t *testing.T) {
				got := legacy.PortTestChatBubbleRegion(x[0], y[0])
				if got != x[1]|y[1] {
					t.Fatalf("region %d want %d", got, x[1]|y[1])
				}
				rows = append(rows, row{x[0], y[0], got})
			})
		}
	}
	spellbookCapture(t, "chat-bubble-regions", rows, "70dac7866c424a0c8f34142ae8692904768a01714af84b160f4bcf871ebbb79b")
}
