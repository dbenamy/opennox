//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func TestClientInventoryTransactionsVisiblePriority(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	coord, free := alloc.Make([]int32{}, 2)
	defer free()
	mods, freeMods := alloc.Make([]uint32{}, 6)
	defer freeMods()
	for _, full := range []bool{false, true} {
		o.reset(t)
		dr := o.item(t, "RedApple", 123)
		for i := range legacy.PortTestInventoryCells() {
			cell := &legacy.PortTestInventoryCells()[i]
			cell.Drawable, cell.Count, cell.Codes[0] = dr, 32, uint32(i+1)
		}
		for _, idx := range []int{20, 41, 62, 83} {
			legacy.PortTestInventoryCells()[idx].Count = 0
		}
		if !full {
			legacy.PortTestInventoryCells()[63].Count = 0
			legacy.PortTestInventoryCells()[1].Count = 0
		}
		coord[0], coord[1] = -99, -99
		ret := o.call(3, 555, uintptr(dr.TypeIDVal), txptr(unsafe.Pointer(&mods[0])), txptr(unsafe.Pointer(&coord[0])))
		if full {
			if ret != 0 || coord[0] != -99 || coord[1] != -99 {
				t.Fatal("hidden row must not expand pickup capacity")
			}
		} else if ret != 1 || coord[0] != 3 || coord[1] != 0 {
			t.Fatalf("pickup priority coordinates%v", coord)
		}
	}
	o.reset(t)
	dr := o.item(t, "RedApple", 0)
	for _, idx := range []int{1, 63, 20} {
		cell := &legacy.PortTestInventoryCells()[idx]
		cell.Drawable, cell.Count, cell.Codes[0] = dr, 1, uint32(idx)
	}
	ret := o.call(4, 999, uintptr(dr.TypeIDVal), 0, 0)
	if ret != o.cell(63) || legacy.PortTestInventoryCells()[63].Codes[1] != 999 {
		t.Fatal("stack append must prioritize row before column")
	}
	o.c.Things.TypeByID("RedApple").ObjClass = 0x4000000
	if ret := o.call(4, 1000, uintptr(dr.TypeIDVal), 0, 0); ret != 0 {
		t.Fatal("type-level no-stack flag ignored")
	}
	o.c.Things.TypeByID("RedApple").ObjClass = 0x10
}

func TestClientInventoryTransactionsErrorRing(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	for i := 0; i < 5; i++ {
		o.call(6, 999, 0, 0, 0)
		slot := (i + 1) % 3
		text := legacy.GoWStringP(memmap.PtrOff(0x5D4594, 823804+uintptr(644*slot)))
		if text != "DroppedNotFound" {
			t.Fatalf("centered error %q", text)
		}
		if *o.txwords[11] != uint32(slot) {
			t.Fatal("centered message ring wrap")
		}
		if got := memmap.Uint32(0x5D4594, 824440+uintptr(644*slot)); got != o.c.srv.Frame()+5*uint32(o.c.srv.TickRate()) {
			t.Fatalf("error expiry frame%d", got)
		}
		if len(o.console) != i+1 || !strings.HasSuffix(o.console[i], "System: DroppedNotFound") {
			t.Fatalf("actual formatted console%v", o.console)
		}
	}
}

func TestClientInventoryTransactionsEquipmentClasses(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	cases := [][3]uint32{{0, 0, 9}, {0x1000, 0, 7}, {0x1000000, 0, 7}, {0x1000000, 2, 0}, {0x1000000, 4, 8}, {0x2000000, 1, 1}, {0x2000000, 4, 2}, {0x2000000, 0x100, 2}, {0x2000000, 0x10, 3}, {0x2000000, 0x80, 3}, {0x2000000, 0x20, 4}, {0x2000000, 2, 8}, {0x2000000, 8, 5}, {0x2000000, 0, 9}, {0x3000000, 3, 0}, {0x3000000, 0x145, 1}}
	for _, c := range cases {
		o.reset(t)
		dr := o.item(t, "RedApple", 123)
		dr.ObjClass, dr.ObjSubClass = object.Class(c[0]), object.SubClass(c[1])
		if got := o.call(11, txptr(dr.C()), 0, 0, 0); got != c[2] {
			t.Fatalf("class%x sub%x slot%d want%d", c[0], c[1], got, c[2])
		}
	}
}
