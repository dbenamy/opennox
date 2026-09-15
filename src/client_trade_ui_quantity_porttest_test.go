//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (o *tradeUIOwner) quantityArgs(x, y int, maximum uint32) []uintptr {
	return []uintptr{txptr(unsafe.Pointer(alloc.InternCString16("Choose quantity"))), uintptr(x), uintptr(y), 123, uintptr(o.c.Things.TypeByID("RedApple").Index()), 0, uintptr(maximum), 77, txptr(o.accept), txptr(o.cancel)}
}
func tradeQuantityText(w *gui.Window, id uint) string {
	return alloc.GoString16((*gui.StaticTextData)(w.ChildByID(id).WidgetData).Text)
}
func TestClientTradeUIQuantityPriceBounds(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, maximum := range []uint32{0, 1, 2, 32} {
		for _, priced := range []uintptr{0, 1, 2} {
			for _, unit := range []uint32{0, 7, 0xffffffff} {
				o.reset(t)
				o.construct(t)
				o.initAmount(t)
				legacy.PortTestTradeUI(9, priced, uintptr(unit))
				rows = append(rows, o.tradeCapture(t, len(rows), 7, o.quantityArgs(200, 100, maximum)...))
				dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
				quantity := uint32(1)
				for _, id := range []uint{3603, 3602, 3602, 3603, 3603, 3602} {
					rows = append(rows, o.tradeCapture(t, len(rows), 5, txptr(dialog.C()), 16391, txptr(dialog.ChildByID(id).C()), 0))
					if id == 3602 && quantity+1 <= maximum {
						quantity++
					}
					if id == 3603 && quantity > 1 {
						quantity--
					}
					if got := tradeQuantityText(dialog, 3601); got != fmt.Sprint(quantity) {
						t.Fatalf("maximum%d button%d quantity%q want%d", maximum, id, got, quantity)
					}
					if priced != 0 {
						if got := tradeQuantityText(dialog, 3607); got != fmt.Sprint(int32(unit*quantity)) {
							t.Fatalf("price %q want%d", got, int32(unit*quantity))
						}
					}
				}
				rows = append(rows, o.tradeCapture(t, len(rows), 4, txptr(dialog.C())))
			}
		}
	}
	tradeUICapture(t, "quantity-price-bounds", rows, "b3ced2334402b26ecb68b25f88a513909e79f09eeb589473736cb038684faa18")
}

func TestClientTradeUIQuantityCallbacks(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	cases := []struct {
		text string
		want uint32
	}{{"0", 0}, {"1", 1}, {"2", 2}, {"31", 31}, {"32", 32}, {"33", 32}, {"-1", 32}, {"+2", 2}, {" 7tail", 7}, {"word", 0}, {"2147483648", 32}, {"4294967295", 32}, {"-2147483649", 32}}
	for _, cancel := range []bool{false, true} {
		for _, present := range []bool{false, true} {
			for _, tc := range cases {
				o.reset(t)
				o.construct(t)
				o.initAmount(t)
				args := o.quantityArgs(200, 100, 32)
				if !present {
					args[8], args[9] = 0, 0
				}
				legacy.PortTestTradeUI(7, args...)
				dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
				dialog.ChildByID(3601).Func94(&gui.StaticTextSetText{Str: tc.text})
				if got := tradeQuantityText(dialog, 3601); got != tc.text {
					t.Fatalf("quantity fixture text %q want %q", got, tc.text)
				}
				op := 2
				var values []uintptr
				if !cancel {
					op = 5
					values = []uintptr{txptr(dialog.C()), 16391, txptr(dialog.ChildByID(3604).C()), 0}
				}
				rows = append(rows, o.tradeCapture(t, len(rows), op, values...))
				wantCalls := 0
				if present {
					wantCalls = 1
				}
				if len(o.callbacks) != wantCalls {
					t.Fatalf("cancel%v present%v callbacks %v", cancel, present, o.callbacks)
				}
				if present {
					which := uint32(0)
					if cancel {
						which = 1
					}
					want := [7]uint32{which, 200, 100, 123, uint32(o.c.Things.TypeByID("RedApple").Index()), tc.want, 77}
					if o.callbacks[0] != want {
						t.Fatalf("text%q callback %v want%v", tc.text, o.callbacks[0], want)
					}
				}
				if *o.windowWords["dword_5d4594_1319268"] != 0 {
					t.Fatal("completed amount dialog stayed active")
				}
			}
		}
	}
	tradeUICapture(t, "quantity-callbacks", rows, "0523105ae3eec930bab6b3ed86bc5c744717ec17aea848d4101180542831eabe")
}

func TestClientTradeUIQuantityPositions(t *testing.T) {
	o := newTradeUIOwner(t)
	var rows []tradeUIResult
	for _, x := range []int{-20, 0, 102, 103, 104, 563, 564, 640, 1000} {
		for _, y := range []int{-20, 0, 37, 38, 39, 418, 419, 480, 1000} {
			o.reset(t)
			o.construct(t)
			o.initAmount(t)
			legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
			rows = append(rows, o.tradeCapture(t, len(rows), 8, uintptr(x), uintptr(y)))
			dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
			wantX, wantY := x-103, y-38
			if wantX < 0 {
				wantX = 0
			}
			if wantY < 0 {
				wantY = 0
			}
			if wantX+180 >= 640 {
				wantX = 460
			}
			if wantY+100 >= 480 {
				wantY = 380
			}
			if got := dialog.Offs(); got.X != wantX || got.Y != wantY {
				t.Fatalf("input%d,%d position%v want%d,%d", x, y, got, wantX, wantY)
			}
		}
	}
	tradeUICapture(t, "quantity-positions", rows, "3b1af4affdd12697f07781a87f851609844a52b5185cda14d197d419e15ba00e")
}
