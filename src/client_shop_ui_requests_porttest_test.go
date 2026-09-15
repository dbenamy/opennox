//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func shopTradeMessages(r shopUIResult) [][]byte {
	var out [][]byte
	for _, b := range r.Window.Inventory.Messages {
		if len(b) >= 2 && b[0] == 201 {
			out = append(out, b)
		}
	}
	return out
}
func TestClientShopUIRequests(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, op := range []int{7, 23, 24, 25, 27, 28, 33, 34} {
		for _, code := range []uintptr{0, 1, 65535, 0x12345} {
			for _, count := range []uintptr{0, 1, 2, 32, 255, 256} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				typ := uintptr(o.c.Things.TypeByID("RedApple").Index())
				var args []uintptr
				var want []byte
				message := func(sub byte, val uintptr, multi bool) {
					want = []byte{201, sub, byte(val), byte(val >> 8)}
					if multi {
						want = append(want, byte(count))
					}
				}
				switch op {
				case 7:
					args = []uintptr{0, code, typ, count}
					if count == 1 {
						message(22, code, false)
					} else if count != 0 {
						message(23, typ, true)
					}
				case 23:
					args = []uintptr{0, code, typ, count}
					if count == 1 {
						message(24, code, false)
					} else if count != 0 {
						message(25, typ, true)
					}
				case 24:
					args = []uintptr{code}
					message(24, code, false)
				case 25:
					args = []uintptr{typ, count}
					message(25, typ, true)
				case 27:
					args = []uintptr{0, code}
					message(26, code, false)
				case 28:
					args = []uintptr{code}
					message(26, code, false)
				case 33:
					args = []uintptr{typ, code}
					message(22, code, false)
				case 34:
					args = []uintptr{typ, count}
					message(23, typ, true)
				}
				*o.shopWords["dword_5d4594_1098616"], *o.shopWords["dword_5d4594_1098620"] = 1, 1
				r := o.shopCapture(t, len(rows), op, args...)
				rows = append(rows, r)
				got := shopTradeMessages(r)
				if want == nil {
					if len(got) != 0 {
						t.Fatalf("operation%d count%d unexpected request%v", op, count, got)
					}
				} else if len(got) != 1 || !bytes.Equal(got[0], want) {
					t.Fatalf("operation%d code%d count%d request%v want%v", op, code, count, got, want)
				}
				if op == 23 && *o.shopWords["dword_5d4594_1098616"] != 0 {
					t.Fatal("sell callback retained pending flag")
				}
				if op == 27 && *o.shopWords["dword_5d4594_1098620"] != 0 {
					t.Fatal("repair callback retained pending flag")
				}
			}
		}
	}
	shopUICapture(t, "requests", rows, "0e850090c6c21a84b4ff4e8d3a4d736c6505b0b955386c57a32bd108087e87bb")
}
func TestClientShopUIBuyInventoryFull(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, op := range []int{7, 33, 34} {
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		for col := 0; col < 4; col++ {
			for row := 0; row < 20; row++ {
				o.stack(t, col, row, 32, "RedApple", uint32(100+col*1000+row*32))
			}
		}
		typ := uintptr(o.c.Things.TypeByID("RedApple").Index())
		args := []uintptr{typ, 123}
		if op == 7 {
			args = []uintptr{0, 123, typ, 2}
		} else if op == 34 {
			args = []uintptr{typ, 2}
		}
		r := o.shopCapture(t, len(rows), op, args...)
		rows = append(rows, r)
		if got := shopTradeMessages(r); len(got) != 0 {
			t.Fatalf("full inventory request %v", got)
		}
		snd := r.Window.Inventory.Sounds
		if len(snd) == 0 || snd[len(snd)-1] != [2]int{925, 100} {
			t.Fatalf("full inventory sound %v", snd)
		}
	}
	shopUICapture(t, "buy-inventory-full", rows, "45e05279e2548674622698aa62ba7d003eb4023305f5ed57a90cfc26834ef8cf")
}
func TestClientShopUIDrawableLookup(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, active := range []uint32{0, 1, 2} {
		for _, value := range []uint32{0, 0xffffffff} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			c := &legacy.PortTestShopUICells()[23]
			c.Drawable = o.item(t, "RedApple", 123)
			c.Count = 1
			c.Codes[0] = 123
			c.Codes[31] = 456
			c.Value = value
			*o.windowWords["dword_5d4594_1098624"] = active
			for _, code := range []uintptr{0, 123, 456, 999} {
				r := o.shopCapture(t, len(rows), 2, code)
				rows = append(rows, r)
				want := uint32(0)
				if active != 0 && code != 999 {
					want = o.norm(uint32(txptr(c.Drawable.C())))
				}
				if r.Return != want {
					t.Fatalf("lookup active%d code%d price%#x got%#x", active, code, value, r.Return)
				}
			}
		}
	}
	shopUICapture(t, "drawable-lookup", rows, "5d0a8fa2b6472b7dc4731461e4856ad6c24bfd224372946213278a5c91a9e30d")
}
func TestClientShopUIRemoveCodes(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, count := range []uint32{1, 2, 3, 32} {
		for _, slot := range []int{0, 1, 15, 31} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			c := &legacy.PortTestShopUICells()[10]
			c.Drawable = o.item(t, "RedApple", 123)
			c.Count = count
			c.Value = 57
			for j := range c.Codes {
				c.Codes[j] = uint32(100 + j)
			}
			before := c.Codes
			code := c.Codes[slot]
			r := o.shopCapture(t, len(rows), 19, uintptr(code))
			rows = append(rows, r)
			for j := 0; j < 31; j++ {
				want := before[j]
				if j >= slot {
					want = before[j+1]
				}
				if c.Codes[j] != want {
					t.Fatalf("count%d slot%d shift at%d", count, slot, j)
				}
			}
			if c.Codes[31] != 0 || c.Count != count-1 {
				t.Fatal("remove count/tail")
			}
			if count == 1 {
				if c.Drawable != nil || c.Value != 0 {
					t.Fatal("last removal retained item/price")
				}
			} else if c.Drawable == nil || c.Value != 57 {
				t.Fatal("partial removal changed unit price/lifetime")
			}
		}
	}
	// Helpers scan all 32 slots even when Count is smaller; preserve that actual
	// storage convention, including first-match priority and missing-code return.
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	c := &legacy.PortTestShopUICells()[0]
	for j := range c.Codes {
		c.Codes[j] = uint32(10 + j)
	}
	before := *c
	r := o.shopCapture(t, len(rows), 20, txptr(unsafe.Pointer(c)), 999)
	rows = append(rows, r)
	if r.Return != 32 || *c != before {
		t.Fatal("missing code helper mutated storage")
	}
	shopUICapture(t, "remove-codes", rows, "bdf4610203c09715c0b36bb614e70cb72ece8fae2767d51a5ceb92d93955d6f2")
}
