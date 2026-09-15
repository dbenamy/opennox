//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func tradeUIAddRecord(data []byte, side byte, typ, code uint16, value uint32) {
	clear(data)
	data[0], data[1], data[2] = 201, 4, side
	binary.LittleEndian.PutUint16(data[3:], typ)
	binary.LittleEndian.PutUint16(data[5:], code)
	binary.LittleEndian.PutUint32(data[7:], value)
	for i := 11; i < 15; i++ {
		data[i] = 255
	}
}
func TestClientTradeUIQuantityAllocationFailure(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, open := range []bool{false, true} {
		o.reset(t)
		o.construct(t)
		o.initAmount(t)
		if open {
			legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
		}
		oldItem := *o.windowWords["nox_gui_itemAmount_item_1319256"]
		oldState := *o.windowWords["dword_5d4594_1319268"]
		dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
		oldPos := dialog.Offs()
		oldCount := tradeQuantityText(dialog, 3601)
		pool := o.c.Objs.Alloc.PortTestCounts()
		o.exhausted(t, func() {
			if got := legacy.PortTestTradeUI(7, o.quantityArgs(400, 300, 16)...); got != 0 {
				t.Fatalf("failed allocation returned%d", got)
			}
		})
		if *o.windowWords["nox_gui_itemAmount_item_1319256"] != oldItem || *o.windowWords["dword_5d4594_1319268"] != oldState || dialog.Offs() != oldPos || tradeQuantityText(dialog, 3601) != oldCount || o.c.Objs.Alloc.PortTestCounts() != pool {
			t.Fatal("failed allocation changed existing quantity dialog")
		}
		rows = append(rows, o.tradeCapture(t, len(rows), 20))
	}
	tradeUICapture(t, "quantity-allocation", rows, "381635116ce6fecc85a97e47a18d63ccd1600721bd378348e5aae9c313e55a36")
}
func TestClientTradeUIMissingResources(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, amount := range []bool{false, true} {
		o.reset(t)
		o.construct(t)
		before := len(o.objects)
		op := 34
		if amount {
			o.missingAmount = true
			op = 3
		} else {
			o.missingTrade = true
		}
		r := o.tradeCapture(t, len(rows), op)
		rows = append(rows, r)
		if r.Return != 0 || len(o.objects) != before {
			t.Fatal("missing resource created dialog objects")
		}
		if amount {
			if *o.windowWords["nox_gui_itemAmount_dialog_1319228"] != 0 {
				t.Fatal("missing quantity resource returned a window")
			}
		} else if o.tradeWindow() != nil {
			t.Fatal("missing trade resource returned a window")
		}
	}
	tradeUICapture(t, "missing-resources", rows, "2d913258e781726c492bde598c9d14ccb5fc7eb4d99aeb204f107f58e38363c8")
}
func TestClientTradeUIAddAllocationFailure(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data, free := alloc.Make([]byte{}, 15)
	defer free()
	for _, side := range []byte{1, 0} {
		o.reset(t)
		o.construct(t)
		o.constructTrade(t)
		*o.windowWords["dword_5d4594_1320964"] = 1
		tradeUIAddRecord(data, side, uint16(o.c.Things.TypeByID("RedApple").Index()), 100, 7)
		before := len(o.objects)
		pool := o.c.Objs.Alloc.PortTestCounts()
		o.exhausted(t, func() {
			if got := legacy.PortTestTradeUI(27, txptr(unsafe.Pointer(&data[0]))); got != 0 {
				t.Fatalf("failed add returned%d", got)
			}
		})
		for _, cells := range legacy.PortTestTradeUICells() {
			for _, cell := range cells {
				if cell.Drawable != nil || cell.Count != 0 || cell.Value != 0 {
					t.Fatal("failed allocation published trade item")
				}
			}
		}
		if len(o.objects) != before || o.c.Objs.Alloc.PortTestCounts() != pool {
			t.Fatal("failed add changed live objects")
		}
		rows = append(rows, o.tradeCapture(t, len(rows), 20))
	}
	tradeUICapture(t, "add-allocation", rows, "72326ab9c4441df6d80958344fa972d2a38e2eacf8e360910cf90f999b8e98f6")
}
func TestClientTradeUIMissingRemoval(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	data, free := alloc.Make([]byte{}, 4)
	defer free()
	data[0], data[1] = 201, 5
	binary.LittleEndian.PutUint16(data[2:], 999)
	for _, active := range []uint32{0, 1, 2} {
		o.reset(t)
		o.construct(t)
		o.constructTrade(t)
		*o.windowWords["dword_5d4594_1320964"] = active
		before := len(o.events)
		left := append([]byte(nil), o.tradeRegions[0]...)
		right := append([]byte(nil), o.tradeRegions[1]...)
		r := o.tradeCapture(t, len(rows), 35, txptr(unsafe.Pointer(&data[0])))
		rows = append(rows, r)
		if r.Return != 0 || len(o.events) != before || !bytes.Equal(left, o.tradeRegions[0]) || !bytes.Equal(right, o.tradeRegions[1]) {
			t.Fatal("missing-item removal changed trade ownership/grid")
		}
		if active != 0 {
			if len(o.console) != 1 || !strings.Contains(o.console[0], "Trade item missing") {
				t.Fatalf("missing item was not reported: %v", o.console)
			}
		} else if len(o.console) != 0 {
			t.Fatal("inactive removal reported a missing item")
		}
	}
	tradeUICapture(t, "missing-removal", rows, "7bf2cbf7991c34652246f47976345a1f3e3c540174699dfc3a243816482d384c")
}
