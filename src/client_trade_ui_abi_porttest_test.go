//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTradeUICQuantityCallbackGC(t *testing.T) {
	o := newTradeUIOwner(t)
	var calls [][7]uint32
	_, _, restore := legacy.PortTestTradeUIObserve(func(v [7]uint32) { runtime.GC(); calls = append(calls, v) })
	t.Cleanup(restore)
	for cycle := 0; cycle < 100; cycle++ {
		o.reset(t)
		o.construct(t)
		o.initAmount(t)
		args := o.quantityArgs(200, 100, 32)
		args[3] = 0xffffffff
		args[7] = 0x80000001
		legacy.PortTestTradeUIBoundary(7, args...)
		w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
		legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3604).C()), 0)
		want := [7]uint32{0, 200, 100, 0xffffffff, uint32(o.c.Things.TypeByID("RedApple").Index()), 1, 0x80000001}
		if len(calls) != cycle+1 || calls[cycle] != want {
			t.Fatalf("callback cycle%d: %v", cycle, calls)
		}
		if *o.windowWords["dword_5d4594_1319268"] != 0 || *o.windowWords["nox_gui_itemAmount_item_1319256"] != 0 {
			t.Fatal("callback did not release quantity drawable")
		}
	}
}
func TestClientTradeUICPriceAndTooltip(t *testing.T) {
	o := newTradeUIOwner(t)
	data := o.record(t, 15)
	for side, reportSide := range []byte{1, 0} {
		for _, price := range []uint32{0, 1, 31, 4096, 0x40000000, 0x80000000, 0xffffffff} {
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			*o.windowWords["dword_5d4594_1320964"] = 1
			tradeUIAddRecord(data, reportSide, uint16(o.c.Things.TypeByID("RedApple").Index()), 123, price)
			if got := legacy.PortTestTradeUIBoundary(27, txptr(unsafe.Pointer(&data[0]))); got != price {
				t.Fatalf("C price return %#x want %#x", got, price)
			}
			runtime.GC()
			cell := &legacy.PortTestTradeUICells()[side][0]
			if cell.Value != price || cell.Count != 1 {
				t.Fatal("C price boundary changed metadata")
			}
			cell.Drawable.NetCode32 = 999
			w := o.tradeWindow().ChildByID(uint(3704 + side))
			p := w.GlobalPos()
			w.TooltipFunc(inventoryWindowPoint(p.X+25, p.Y+25))
			if cell.Drawable.NetCode32 != 123 {
				t.Fatal("raw C tooltip callback did not select the first item")
			}
		}
	}
}
