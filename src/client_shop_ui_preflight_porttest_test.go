//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUIConstruction(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	for id := uint(3801); id <= 3810; id++ {
		if o.shopWindow().ChildByID(id) == nil {
			t.Errorf("missing shop widget %d", id)
		}
	}
	if !o.shopWindow().Flags.IsHidden() {
		t.Error("new shop must start hidden")
	}
	for i, c := range legacy.PortTestShopUICells() {
		if c.Drawable != nil || c.Count != 0 || c.Value != 0 || c.Codes != [32]uint32{} {
			t.Errorf("cell %d not empty", i)
		}
	}
	shopUICapture(t, "construction", []shopUIResult{o.shopCapture(t, 0, 0)}, "5edc5d6e24302ef67b7304f51458b25e95515d6e3fcc0a52c72ccd63b15dcf16")
}
func TestClientShopUIFullStackSelection(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	cells := legacy.PortTestShopUICells()
	cells[0].Drawable = o.item(t, "RedApple", 100)
	cells[0].Count = 32
	typ := uintptr(cells[0].Drawable.TypeIDVal)
	if got, want := legacy.PortTestShopUI(17, typ), uint32(txptr(unsafe.Pointer(&cells[10]))); got != want {
		t.Errorf("full stack selected: got %#x want first empty %#x", got, want)
	}
	cells[10].Drawable = o.item(t, "RedApple", 101)
	cells[10].Count = 31
	if got, want := legacy.PortTestShopUI(17, typ), uint32(txptr(unsafe.Pointer(&cells[10]))); got != want {
		t.Errorf("full first stack blocks partial second: got %#x want %#x", got, want)
	}
	for i := range cells {
		if cells[i].Drawable == nil {
			cells[i].Drawable = o.item(t, "RedApple", uint32(200+i))
		}
		cells[i].Count = 32
	}
	if got := legacy.PortTestShopUI(17, typ); got != 0 {
		t.Errorf("all-full shop selected cell %#x", got)
	}
}
func TestClientShopUICodeLookup(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	cells := legacy.PortTestShopUICells()
	for i := range cells {
		cells[i].Drawable = o.item(t, "RedApple", uint32(100+i))
		cells[i].Count = 1
		for j := range cells[i].Codes {
			cells[i].Codes[j] = uint32(1 + i*32 + j)
		}
	}
	for i := range cells {
		for _, slot := range []int{0, 1, 15, 31} {
			code := cells[i].Codes[slot]
			if got, want := legacy.PortTestShopUI(3, uintptr(code)), uint32(txptr(unsafe.Pointer(&cells[i]))); got != want {
				t.Fatalf("cell %d code slot %d: got %#x want %#x", i, slot, got, want)
			}
		}
	}
	if got := legacy.PortTestShopUI(3, 99999); got != 0 {
		t.Fatalf("missing item returned %#x", got)
	}
}
func TestClientShopUIHitRegions(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	p, free := alloc.New([2]int32{})
	defer free()
	w := o.shopWindow().ChildByID(3806)
	pos, sz := w.GlobalPos(), w.Size()
	cells := legacy.PortTestShopUICells()
	for i := range cells {
		cells[i].Drawable = o.item(t, "RedApple", uint32(100+i))
		cells[i].Count = 1
		cells[i].Codes[0] = uint32(900 + i)
	}
	count := 0
	for _, scroll := range []uint32{0, 1, 49, 50, 125, 299, 300} {
		*o.windowWords["dword_5d4594_1107036"] = scroll
		for x := pos.X - 2; x <= pos.X+sz.X+2; x += 1 {
			for y := pos.Y - 2; y <= pos.Y+sz.Y+2; y += 1 {
				*p = [2]int32{int32(x), int32(y)}
				inside := x >= pos.X && x <= pos.X+sz.X && y >= pos.Y && y <= pos.Y+sz.Y
				got := legacy.PortTestShopUI(30, txptr(unsafe.Pointer(p)))
				if (got != 0) != inside {
					t.Fatalf("inside at %d,%d: got%d want%v", x, y, got, inside)
				}
				ret := legacy.PortTestShopUI(31, txptr(unsafe.Pointer(p)))
				var want uint32
				if inside {
					col, row := (x-pos.X)/50, (y-pos.Y+int(scroll))/50
					if col > 5 {
						col = 5
					}
					if row > 9 {
						row = 9
					}
					c := &cells[col*10+row]
					want = uint32(txptr(c.Drawable.C()))
					if c.Drawable.NetCode32 != c.Codes[0] {
						t.Fatal("hit did not publish first code")
					}
				}
				if ret != want {
					t.Fatalf("hit at %d,%d scroll%d: got%#x want%#x", x, y, scroll, ret, want)
				}
				count++
			}
		}
	}
	t.Logf("%d independent inside/item coordinate contracts", count)
}
