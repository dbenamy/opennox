//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUIStockCapacity(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	mods, free := alloc.New([4]byte{})
	defer free()
	*mods = [4]byte{255, 255, 255, 255}
	typ := uintptr(o.c.Things.TypeByID("RedApple").Index())
	var rows []shopUIResult
	cells := legacy.PortTestShopUICells()
	// Fill all 1,920 legitimate slots through the actual add owner. No C call is
	// made beyond capacity until the independent selection prerequisite is fixed.
	for n := 0; n < 60*32; n++ {
		scan := n / 32
		idx := (scan%6)*10 + scan/6
		ret := legacy.PortTestShopUI(16, typ, uintptr(100+n), uintptr(7+n%5), 123, txptr(unsafe.Pointer(mods)))
		if ret != uint32(n%32+1) || cells[idx].Count != ret || cells[idx].Codes[n%32] != uint32(100+n) || cells[idx].Value != uint32(7+n%5) {
			t.Fatalf("add %d cell%d count%v return%d", n, idx, cells[idx].Count, ret)
		}
		if n < 3 || n%32 == 31 {
			rows = append(rows, o.shopCapture(t, len(rows), 0))
		}
	}
	before := append([]legacy.PortTestShopUICell(nil), cells...)
	if got := legacy.PortTestShopUI(16, typ, 99999, 88, 123, txptr(unsafe.Pointer(mods))); got != 0 {
		t.Fatalf("full shop accepted item: %d", got)
	}
	for i := range cells {
		if cells[i] != before[i] {
			t.Fatalf("full shop mutated cell %d", i)
		}
	}
	rows = append(rows, o.shopCapture(t, len(rows), 0))
	shopUICapture(t, "stock-capacity", rows, "ca3d769f1bff64454877f63cc895959f7a4ee67b63e0ce82da401ded7bf16a75")
}
func TestClientShopUIAllocationFailure(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	c := &legacy.PortTestShopUICells()[0]
	c.Value = 77
	c.Codes[5] = 999
	before := *c
	o.exhausted(t, func() {
		if ret := legacy.PortTestShopUI(16, uintptr(o.c.Things.TypeByID("RedApple").Index()), 123, 7, 31, 0); ret != 0 {
			t.Errorf("failed allocation return %d", ret)
		}
	})
	if *c != before {
		t.Fatal("failed stock allocation changed cell")
	}
	shopUICapture(t, "allocation-failure", []shopUIResult{o.shopCapture(t, 0, 0)}, "a6172e37d31db9828c26ccd5b3e6ad9178f3b8a717712065f635e851154e38fd")
}
func TestClientShopUIResetOwnership(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, op := range []int{12, 13} {
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		cells := legacy.PortTestShopUICells()
		for i := range cells {
			cells[i].Drawable = o.item(t, "RedApple", uint32(100+i))
			cells[i].Count = uint32(i%32 + 1)
			cells[i].Value = uint32(71 + i)
			for j := range cells[i].Codes {
				cells[i].Codes[j] = uint32(1000 + i*32 + j)
			}
		}
		*o.windowWords["dword_5d4594_1107036"] = 125
		eventStart := len(o.events)
		refs := make([]uint32, 60)
		for i := range cells {
			refs[i] = o.objects[o.identities[cells[i].Drawable]].Ref
		}
		rows = append(rows, o.shopCapture(t, len(rows), op))
		if op == 13 {
			o.shopReady = false
		}
		for i, c := range cells {
			if c.Drawable != nil || c.Count != 0 || c.Value != 0 || c.Codes != [32]uint32{} {
				t.Fatalf("reset cell%d retained contents", i)
			}
		}
		if *o.windowWords["dword_5d4594_1107036"] != 0 {
			t.Fatal("reset retained scroll")
		}
		if len(o.events)-eventStart != 60 {
			t.Fatalf("reset events %d", len(o.events)-eventStart)
		}
		for scan, e := range o.events[eventStart:] {
			idx := (scan%6)*10 + scan/6
			if e[0] != 2 || e[2] != refs[idx] {
				t.Fatalf("reset deletion %d: %v want ref%#x", scan, e, refs[idx])
			}
		}
		rows = append(rows, o.shopCapture(t, len(rows), 12))
		if len(o.events)-eventStart != 60 {
			t.Fatal("reset deleted twice")
		}
	}
	shopUICapture(t, "reset-ownership", rows, "b1c703af467247bd27b890e742770177450e44c52b24f58da236fc7ed3f015f5")
}
