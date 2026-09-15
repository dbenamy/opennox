//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTradeUIQuantityReopenOwnership(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.initAmount(t)
	legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
	pool := o.c.Objs.Alloc.PortTestCounts()
	for cycle := 0; cycle < 100; cycle++ {
		previous := (*client.Drawable)(unsafe.Pointer(uintptr(*o.windowWords["nox_gui_itemAmount_item_1319256"])))
		if previous == nil {
			t.Fatal("open dialog has no drawable")
		}
		legacy.PortTestTradeUI(7, o.quantityArgs(250, 150, 16)...)
		current := (*client.Drawable)(unsafe.Pointer(uintptr(*o.windowWords["nox_gui_itemAmount_item_1319256"])))
		if *o.windowWords["dword_5d4594_1319268"] != 1 || current == nil {
			t.Fatal("reopening closed the dialog or discarded its new drawable")
		}
		if o.objects[o.identities[previous]].Live {
			t.Fatal("reopening retained previous drawable")
		}
		if !o.objects[o.identities[current]].Live || o.c.Objs.Alloc.PortTestCounts() != pool {
			t.Fatal("reopening changed owned live count")
		}
	}
}

func TestClientTradeUIFullStackSelection(t *testing.T) {
	for side, op := range []int{29, 30} {
		t.Run(fmt.Sprintf("side%d", side), func(t *testing.T) {
			o := newTradeUIOwner(t)
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			cells := legacy.PortTestTradeUICells()[side]
			cells[0].Drawable = o.item(t, "RedApple", 100)
			cells[0].Count = 32
			typ := uintptr(cells[0].Drawable.TypeIDVal)
			if got := legacy.PortTestTradeUI(28, typ, txptr(unsafe.Pointer(&cells[0]))); got != 0 {
				t.Errorf("preferred full cell accepted another item: %d", got)
			}
			if got, want := legacy.PortTestTradeUI(op, typ), uint32(txptr(unsafe.Pointer(&cells[2]))); got != want {
				t.Errorf("full first stack: got%#x want first empty%#x", got, want)
			}
			cells[1].Drawable = o.item(t, "RedApple", 200)
			cells[1].Count = 31
			if got, want := legacy.PortTestTradeUI(op, typ), uint32(txptr(unsafe.Pointer(&cells[1]))); got != want {
				t.Errorf("full first stack: got%#x want compatible partial%#x", got, want)
			}
			cells[1].Count = 32
			for _, i := range []int{2, 3} {
				cells[i].Drawable = o.item(t, "RedApple", uint32(300+i))
				cells[i].Count = 32
			}
			if got := legacy.PortTestTradeUI(op, typ); got != 0 {
				t.Errorf("all stacks full returned%#x", got)
			}
		})
	}
}
