//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientTradeUIQuantityMouse(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for event := 0; event <= 16; event++ {
		for _, inside := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			o.initAmount(t)
			legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
			w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
			x, y := 500, 400
			if inside {
				x, y = 200, 100
			}
			r := o.tradeCapture(t, len(rows), 1, txptr(w.C()), uintptr(event), inventoryWindowPoint(x, y))
			rows = append(rows, r)
			want := uint32(0)
			switch event {
			case 5, 6, 7, 9, 10, 11, 13, 14, 15:
				want = 1
			}
			cancel := !inside && (event == 5 || event == 9 || event == 13)
			if r.Return != want || (*o.windowWords["dword_5d4594_1319268"] == 0) != cancel || (len(o.callbacks) == 1) != cancel {
				t.Fatalf("event%d inside%v: return%d callbacks%v", event, inside, r.Return, o.callbacks)
			}
		}
	}
	tradeUICapture(t, "quantity-mouse", rows, "3feda7aec486a9b7a2f1cd9a57884014f4e924639d98708a3911ae9b1f44540a")
}

func TestClientTradeUIQuantityDrawButtons(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, priced := range []uintptr{0, 1} {
		for mask := 0; mask < 16; mask++ {
			o.reset(t)
			o.construct(t)
			o.initAmount(t)
			legacy.PortTestTradeUI(9, priced, 7)
			legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
			w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
			for i, id := range []uint{3602, 3603, 3604, 3605} {
				if mask&(1<<i) != 0 {
					w.ChildByID(id).DrawData().Field0 |= 4
				}
			}
			rows = append(rows, o.tradeCapture(t, len(rows), 4, txptr(w.C())))
			if len(o.callbacks) != 0 || *o.windowWords["dword_5d4594_1319268"] != 1 {
				t.Fatal("drawing completed the dialog")
			}
		}
	}
	tradeUICapture(t, "quantity-draw-buttons", rows, "4e7045671ec7351fe77ef10a57fcc728dde93cbdf903df196c61c34dc7ac4005")
}

func TestClientTradeUIPopulatedDrawHover(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, offset := range []image.Point{{0, 0}, {20, 30}} {
		for mask := 0; mask < 32; mask++ {
			o.reset(t)
			o.construct(t)
			o.constructTrade(t)
			o.tradeWindow().SetPos(offset)
			*o.tradeWords["dword_5d4594_1320944"] = uint32(mask & 1)
			*o.tradeWords["dword_5d4594_1320948"] = uint32(mask & 2)
			if mask&4 != 0 {
				o.tradeWindow().ChildByID(3708).DrawData().Field0 |= 4
			}
			if mask&8 != 0 {
				o.tradeWindow().ChildByID(3710).DrawData().Field0 |= 4
			}
			o.tradeRegions[10][0] = byte(mask & 16)
			for side := 0; side < 2; side++ {
				for i := 0; i < 4; i++ {
					o.tradeCell(t, side, i, []int{1, 2, 31, 32}[i], "RedApple", uint32(100+side*100+i*32))
				}
			}
			rows = append(rows, o.tradeCapture(t, len(rows), 16))
			for side, cells := range legacy.PortTestTradeUICells() {
				origin := o.tradeWindow().ChildByID(uint(3704 + side)).GlobalPos()
				for i, cell := range cells {
					x, y := origin.X+25+(i/2)*50, origin.Y+25+(i%2)*50
					if *txword(cell.Drawable, 12) != uint32(x) || *txword(cell.Drawable, 16) != uint32(y) {
						t.Fatal("trade drawable position")
					}
					// Hover is independent of stack size and publishes the first item code.
					rows = append(rows, o.tradeCapture(t, len(rows), 17, txptr(o.tradeWindow().C()), 0, inventoryWindowPoint(x, y)))
					if *txword(cell.Drawable, 128) != cell.Codes[0] {
						t.Fatal("trade hover selected wrong item code")
					}
				}
			}
		}
	}
	tradeUICapture(t, "populated-draw-hover", rows, "b3d0db51342fbb90d8a173e292747b9d11ee65273d4643c6494d725553924438")
}

func TestClientTradeUICodeHelpers(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructTrade(t)
	cell := &legacy.PortTestTradeUICells()[0][0]
	ptr := txptr(unsafe.Pointer(cell))
	var rows []tradeUIResult
	for _, count := range []uint32{0, 1, 31, 32} {
		for which := -1; which < 32; which++ {
			cell.Count = count
			cell.Value = 1234
			for i := range cell.Codes {
				cell.Codes[i] = uint32(100 + i)
			}
			code := uintptr(100 + which)
			r := o.tradeCapture(t, len(rows), 26, ptr, code)
			rows = append(rows, r)
			want := uint32(0)
			if count != 0 && which >= 0 {
				want = 1
			}
			if r.Return != want {
				t.Fatalf("find count%d slot%d return%d", count, which, r.Return)
			}
			before := cell.Codes
			r = o.tradeCapture(t, len(rows), 25, ptr, code)
			rows = append(rows, r)
			result := uint32(32)
			if which >= 0 {
				copy(before[which:], before[which+1:])
				before[31] = 0
				result = 131
				if which == 31 {
					result = 31
				}
			}
			if cell.Codes != before || cell.Count != count || cell.Value != 1234 || r.Return != result {
				t.Fatalf("remove count%d slot%d return%d want%d", count, which, r.Return, result)
			}
		}
	}
	tradeUICapture(t, "code-helpers", rows, "646b72e0bae8943a5dab09014090e8527d8811ab9805b639418eb891dcd0610e")
}
