//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func inventoryCancelDragged() *client.Drawable {
	return (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1049848))
}
func (o *inventoryWindowOwner) assertCancelReleased(t *testing.T, drag *client.Drawable, borrowed bool, before int) {
	t.Helper()
	if inventoryCancelDragged() != nil || o.c.dragndropItem != nil || *o.windowWords["dword_5d4594_1049856"] != 0 {
		t.Fatal("cancellation retained drag state")
	}
	if o.c.GUI.Captured() != nil {
		t.Fatal("cancellation retained mouse capture")
	}
	idx, ok := o.identities[drag]
	if !ok || o.objects[idx].Live != borrowed {
		t.Fatalf("drag ownership live=%v borrowed=%v", ok && o.objects[idx].Live, borrowed)
	}
	deletes := 0
	for _, e := range o.events[before:] {
		if e[0] == 2 && e[2] == o.objects[idx].Ref {
			deletes++
		}
	}
	want := 1
	if borrowed {
		want = 0
	}
	if deletes != want {
		t.Fatalf("drag deleted %d times, want %d", deletes, want)
	}
	before = len(o.events)
	if ret := o.invoke(33, 0, 0, 0, 0); ret != 0 {
		t.Fatalf("repeated cancellation returned %d", ret)
	}
	if len(o.events) != before {
		t.Fatal("repeated cancellation allocated or deleted")
	}
}

func TestClientInventoryWindowCancelRepeatedOwnership(t *testing.T) {
	for _, count := range []int{1, 2, 32} {
		for _, op := range []int{33, 29} {
			t.Run(fmt.Sprintf("stack%d/op%d", count, op), func(t *testing.T) {
				o := newInventoryWindowOwner(t)
				o.reset(t)
				o.construct(t)
				o.stack(t, 0, 0, count, "RedApple", 100)
				pool := o.c.Objs.Alloc.PortTestCounts()
				live := o.c.Objs.Count
				for cycle := 0; cycle < 100; cycle++ {
					o.openInput()
					o.invoke(1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0)
					drag := inventoryCancelDragged()
					if drag == nil {
						t.Fatal("press did not create drag")
					}
					before := len(o.events)
					if o.invoke(op, 0, 0, 0, 0) != 1 {
						t.Fatal("cancellation failed")
					}
					o.assertCancelReleased(t, drag, false, before)
					if got := o.c.Objs.Alloc.PortTestCounts(); got != pool || o.c.Objs.Count != live {
						t.Fatalf("cycle%d pool%v want%v count%d want%d", cycle, got, pool, o.c.Objs.Count, live)
					}
					cell := legacy.PortTestInventoryCells()[0]
					if int(cell.Count) != count || cell.Drawable == nil {
						t.Fatal("cancellation lost stack")
					}
					seen := make(map[uint32]bool)
					for _, code := range cell.Codes[:cell.Count] {
						if code < 100 || code >= 100+uint32(count) || seen[code] {
							t.Fatal("cancellation changed item identities")
						}
						seen[code] = true
					}
				}
			})
		}
	}
}

func TestClientInventoryWindowCancelBorrowedEquipment(t *testing.T) {
	for _, op := range []int{33, 29} {
		t.Run(fmt.Sprintf("op%d", op), func(t *testing.T) {
			o := newInventoryWindowOwner(t)
			o.reset(t)
			o.construct(t)
			dr := o.item(t, "Bow", 100)
			o.equipment[7] = uint32(txptr(dr.C()))
			pool := o.c.Objs.Alloc.PortTestCounts()
			for cycle := 0; cycle < 100; cycle++ {
				o.openInput()
				o.invoke(1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(80, 50), 0)
				if inventoryCancelDragged() != dr || *o.windowWords["dword_5d4594_1049856"] != 1 {
					t.Fatal("equipment press did not borrow drawable")
				}
				before := len(o.events)
				o.invoke(op, 0, 0, 0, 0)
				o.assertCancelReleased(t, dr, true, before)
				if o.equipment[7] != uint32(txptr(dr.C())) || dr.NetCode32 != 100 || o.c.Objs.Alloc.PortTestCounts() != pool {
					t.Fatal("cancellation changed borrowed equipment")
				}
			}
		})
	}
}

func TestClientInventoryWindowCancelRestoreOwnership(t *testing.T) {
	for _, full := range []bool{false, true} {
		t.Run(fmt.Sprintf("full%v", full), func(t *testing.T) {
			o := newInventoryWindowOwner(t)
			o.reset(t)
			o.construct(t)
			o.openInput()
			o.stack(t, 0, 0, 1, "Bow", 100)
			o.invoke(1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0)
			drag := inventoryCancelDragged()
			if drag == nil {
				t.Fatal("missing drag")
			}
			o.stack(t, 0, 0, 32, "Meat", 500)
			if full {
				for col := 0; col < 4; col++ {
					for row := 0; row < 20; row++ {
						if col == 0 && row == 0 {
							continue
						}
						o.stack(t, col, row, 32, "Meat", uint32(1000+(col*20+row)*32))
					}
				}
			}
			pool := o.c.Objs.Alloc.PortTestCounts()
			before := len(o.events)
			o.invoke(33, 0, 0, 0, 0)
			o.assertCancelReleased(t, drag, false, before)
			found := 0
			for _, cell := range legacy.PortTestInventoryCells() {
				for _, code := range cell.Codes[:cell.Count] {
					if code == 100 {
						found++
					}
				}
			}
			want := 1
			delta := 0
			if full {
				want = 0
				delta = -1
			}
			if found != want {
				t.Fatalf("restored instances %d want %d", found, want)
			}
			got := o.c.Objs.Alloc.PortTestCounts()
			if got[0] != pool[0]+delta || got[1] != pool[1]-delta || got[2] != pool[2] {
				t.Fatalf("restoration pool %v before %v expected live delta %d", got, pool, delta)
			}
			if full && len(o.console) == 0 {
				t.Fatal("full inventory was not reported")
			}
		})
	}
}
