//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUIBuyQuantity(t *testing.T) {
	o := newShopUIOwner(t)
	o.installMessageDialog(t)
	var rows []shopUIResult
	p, free := alloc.New([2]int32{})
	defer free()
	for _, count := range []uint32{1, 2, 32} {
		for _, price := range []uint32{0, 1, 7, 0x80000000, 0xffffffff} {
			for _, gold := range []uint32{0, 6, 15, 0xffffffff} {
				sub_44A400()
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				o.initAmount(t)
				c := &legacy.PortTestShopUICells()[0]
				c.Drawable = o.item(t, "RedApple", 123)
				c.Count = count
				c.Value = price
				for i := range c.Codes {
					c.Codes[i] = uint32(123 + i)
				}
				*o.windowWords["dword_5d4594_1062552"] = gold
				pos := o.shopWindow().ChildByID(3806).GlobalPos()
				*p = [2]int32{int32(pos.X + 25), int32(pos.Y + 25)}
				rows = append(rows, o.shopCapture(t, len(rows), 32, txptr(unsafe.Pointer(p))))
				want := count
				if price != 0 && uint64(gold)/uint64(price) < uint64(want) {
					want = uint32(uint64(gold) / uint64(price))
				}
				if want == 0 {
					if *o.windowWords["dword_5d4594_1319268"] != 0 || nox_gui_curDialog_830224 == nil {
						t.Fatalf("count%d price%#x gold%#x unaffordable flow", count, price, gold)
					}
					continue
				}
				if nox_gui_curDialog_830224 != nil || *o.windowWords["dword_5d4594_1319268"] != 1 || *o.windowWords["dword_5d4594_1319248"] != want {
					t.Fatalf("count%d price%#x gold%#x max%d want%d", count, price, gold, *o.windowWords["dword_5d4594_1319248"], want)
				}
				w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
				chosen := uint32(1)
				if want > 1 {
					legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3602).C()), 0)
					chosen = 2
				}
				legacy.PortTestTradeUI(4, txptr(w.C()))
				rows = append(rows, o.shopCapture(t, len(rows), 0))
				legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3604).C()), 0)
				r := o.shopCapture(t, len(rows), 0)
				rows = append(rows, r)
				code := c.Codes[count-1]
				msg := []byte{201, 22, byte(code), byte(code >> 8)}
				if chosen > 1 {
					typ := c.Drawable.TypeIDVal
					msg = []byte{201, 23, byte(typ), byte(typ >> 8), byte(chosen)}
				}
				got := shopTradeMessages(r)
				if len(got) != 1 || !bytes.Equal(got[0], msg) {
					t.Fatalf("buy quantity request %v want%v", got, msg)
				}
				if *o.windowWords["dword_5d4594_1319268"] != 0 {
					t.Fatal("buy accept retained quantity dialog")
				}
			}
		}
	}
	shopUICapture(t, "buy-quantity", rows, "ac58373a27af486283e577ed8308f8cc03f689a5dc619568c260169a4f59f097")
}
func TestClientShopUISellRepairQuantity(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, name := range []string{"RedApple", "Bow"} {
		for _, count := range []int{1, 2, 32} {
			for _, op := range []int{40, 41} {
				for _, accept := range []bool{false, true} {
					o.reset(t)
					o.construct(t)
					o.constructShop(t)
					o.initAmount(t)
					dr := o.stack(t, 0, 0, count, name, 65535)
					typ := dr.TypeIDVal
					if name == "Bow" {
						*txword(dr, 432) = uint32(txptr(o.mods[0].C()))
					}
					*o.windowWords["dword_5d4594_1062552"] = 999
					rows = append(rows, o.shopCapture(t, len(rows), op, 65535, 7))
					pending := "dword_5d4594_1098616"
					maximum := count
					if op == 41 {
						pending = "dword_5d4594_1098620"
						maximum = 1
					}
					if *o.shopWords[pending] != 1 || *o.windowWords["dword_5d4594_1319268"] != 1 || int(*o.windowWords["dword_5d4594_1319248"]) != maximum {
						t.Fatalf("%s op%d count%d quantity did not open", name, op, count)
					}
					item := unsafe.Pointer(uintptr(*o.windowWords["nox_gui_itemAmount_item_1319256"]))
					if name == "Bow" && *(*uint32)(unsafe.Add(item, 432)) != uint32(txptr(o.mods[0].C())) {
						t.Fatal("shop quantity did not copy actual modifier")
					}
					// A duplicate reply while pending must preserve the same quantity drawable.
					legacy.PortTestShopUI(op, 65535, 7)
					if *o.windowWords["nox_gui_itemAmount_item_1319256"] != uint32(txptr(item)) {
						t.Fatal("pending duplicate replaced quantity item")
					}
					w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
					chosen := 1
					if maximum > 1 {
						legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3602).C()), 0)
						chosen = 2
					}
					button := uint(3605)
					if accept {
						button = 3604
					}
					legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(button).C()), 0)
					r := o.shopCapture(t, len(rows), 0)
					rows = append(rows, r)
					got := shopTradeMessages(r)
					var want []byte
					if accept {
						want = []byte{201, 24, 255, 255}
						if op == 41 {
							want[1] = 26
						} else if chosen > 1 {
							want = []byte{201, 25, byte(typ), byte(typ >> 8), byte(chosen)}
						}
					}
					if want == nil {
						if len(got) != 0 {
							t.Fatalf("cancel sent%v", got)
						}
					} else if len(got) != 1 || !bytes.Equal(got[0], want) {
						t.Fatalf("op%d quantity request%v want%v", op, got, want)
					}
					if *o.shopWords[pending] != 0 || *o.windowWords["dword_5d4594_1319268"] != 0 {
						t.Fatal("quantity completion retained flags")
					}
				}
			}
		}
	}
	shopUICapture(t, "sell-repair-quantity", rows, "e872a354a1574158235f0a6fc0b45caac1e24e33857f72da18e03053c4481a5f")
}
