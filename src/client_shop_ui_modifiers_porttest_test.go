//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUIStockModifiers(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	names := make([]*byte, len(o.mods))
	for i, m := range o.mods {
		names[i] = *(**byte)(m.C())
	}
	t.Cleanup(o.c.srv.Server.PortTestControlsModifiers(o.mods, names))
	ids, free := alloc.New([4]byte{})
	defer free()
	for _, name := range []string{"RedApple", "Bow"} {
		for _, health := range []uintptr{0, 1, 32767, 32768, 65535} {
			for _, set := range [][4]byte{{255, 255, 255, 255}, {1, 2, 3, 4}, {4, 3, 2, 1}, {0, 5, 254, 255}} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				*ids = set
				typ := uintptr(o.c.Things.TypeByID(name).Index())
				r := o.shopCapture(t, len(rows), 16, typ, 123, 0xffffffff, health, txptr(unsafe.Pointer(ids)))
				rows = append(rows, r)
				c := &legacy.PortTestShopUICells()[0]
				if r.Return != 1 || c.Drawable == nil || c.Count != 1 || c.Codes[0] != 123 || c.Value != 0xffffffff {
					t.Fatal("stock metadata")
				}
				if *(*uint16)(unsafe.Add(c.Drawable.C(), 292)) != uint16(health) || *(*uint16)(unsafe.Add(c.Drawable.C(), 294)) != uint16(health) || *txword(c.Drawable, 120)&0x40000000 == 0 {
					t.Fatal("stock drawable health/flags")
				}
				for i, id := range set {
					want := uint32(0)
					if name == "Bow" && id >= 1 && id <= 4 {
						want = uint32(txptr(o.mods[id-1].C()))
					}
					if got := *txword(c.Drawable, 432+4*uintptr(i)); got != want {
						t.Fatalf("%s modifier%d ID%d got%#x want%#x", name, i, id, got, want)
					}
				}
			}
		}
	}
	shopUICapture(t, "stock-modifiers", rows, "d4f53bb932fd27e8c644a33fd43a4945a299d95a33d8ecf0cd957c96f1aed30f")
}
func TestClientShopUIStockSelection(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	ids, free := alloc.New([4]byte{})
	defer free()
	*ids = [4]byte{255, 255, 255, 255}
	for _, mask := range []uint32{0, 0x4000000} {
		for _, count := range []uint32{1, 31, 32} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			cells := legacy.PortTestShopUICells()
			c := &cells[0]
			c.Drawable = o.item(t, "RedApple", 123)
			c.Count = count
			c.Value = 57
			*txword(c.Drawable, 112) |= mask
			typ := uintptr(c.Drawable.TypeIDVal)
			want := 0
			if mask != 0 || count == 32 {
				want = 10
			}
			if got := legacy.PortTestShopUI(17, typ); got != uint32(txptr(unsafe.Pointer(&cells[want]))) {
				t.Fatalf("mask%#x count%d lookup selected%#x", mask, count, got)
			}
			rows = append(rows, o.shopCapture(t, len(rows), 16, typ, 777, 61, 99, txptr(unsafe.Pointer(ids))))
			if cells[want].Codes[cells[want].Count-1] != 777 || cells[want].Value != 61 {
				t.Fatal("stock selection appended into wrong cell")
			}
			if want == 0 && *(*uint16)(unsafe.Add(c.Drawable.C(), 292)) != 0 {
				t.Fatal("append replaced original drawable health")
			}
		}
	}
	shopUICapture(t, "stock-selection", rows, "fec65809865c55ea6ca661fe9bc5f3dc28a0b0fe95921026e62c213588b6081f")
}
