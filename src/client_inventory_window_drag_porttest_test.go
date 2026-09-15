//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (o *inventoryWindowOwner) stack(t *testing.T, col, row, count int, name string, code uint32) *client.Drawable {
	dr := o.item(t, name, code)
	cell := &legacy.PortTestInventoryCells()[col*21+row]
	cell.Drawable = dr
	cell.Count = byte(count)
	for i := 0; i < count; i++ {
		cell.Codes[i] = code + uint32(i)
	}
	return dr
}
func (o *inventoryWindowOwner) openInput() {
	*memmap.PtrUint8(0x5D4594, 1049868) = 2
	*o.windowWords["dword_587000_136184"] = 0
}
func inventoryWindowDropMessages(r inventoryWindowResult) [][]byte {
	var out [][]byte
	for _, msg := range r.Inventory.Messages {
		if len(msg) == 7 && msg[0] == 114 {
			out = append(out, msg)
		}
	}
	return out
}

func TestClientInventoryWindowDragDrop(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, count := range []int{1, 2, 3, 32} {
		for _, where := range []string{"same", "empty", "occupied", "outside", "cancel"} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			o.initAmount(t)
			o.stack(t, 0, 0, count, "RedApple", 100)
			if where == "occupied" {
				o.stack(t, 1, 0, 1, "Meat", 900)
			}
			rows = append(rows, o.capture(t, len(rows), 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0))
			if *memmap.PtrUint32(0x5D4594, 1049848) == 0 {
				t.Fatal("press did not pick up inventory item")
			}
			cell := &legacy.PortTestInventoryCells()[0]
			if cell.Count != byte(count-1) || (count > 1 && cell.Codes[0] != 101) {
				t.Fatal("press did not remove the first stack code")
			}
			if dragged := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1049848)); dragged.NetCode32 != 100 {
				t.Fatal("drag copy has the wrong stack code")
			}
			if where == "cancel" {
				rows = append(rows, o.capture(t, len(rows), 1, 33, 0, 0, 0, 0))
			} else {
				x, y := 339, 38
				if where == "empty" || where == "occupied" {
					x = 389
				}
				if where == "outside" {
					x, y = 600, 400
				}
				rows = append(rows, o.capture(t, len(rows), 1, 1, txptr(o.mainWindow().C()), 6, inventoryWindowPoint(x, y), 0))
			}
			if *memmap.PtrUint32(0x5D4594, 1049848) != 0 || o.c.dragndropItem != nil {
				t.Fatalf("%s left an active drag", where)
			}
			if where == "outside" && count > 1 {
				if cell.Count != byte(count) || cell.Codes[count-1] != 100 || cell.Codes[0] != 101 {
					t.Fatal("quantity prompt did not restore dragged code at stack end")
				}
				dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
				if *o.windowWords["dword_5d4594_1319268"] != 1 {
					t.Fatal("stack drop did not open actual amount dialog")
				}
				up := dialog.ChildByID(3602)
				dialog.Func94(&gui.RawEvent{Event: 16391, Arg1: txptr(up.C())})
				// Capture the external dialog transition through a no-op inventory operation.
				rows = append(rows, o.capture(t, len(rows), 2, 7, 0, 0, 0, 0))
				yes := dialog.ChildByID(3604)
				dialog.Func94(&gui.RawEvent{Event: 16391, Arg1: txptr(yes.C())})
				r := o.capture(t, len(rows), 3, 7, 0, 0, 0, 0)
				rows = append(rows, r)
				msgs := inventoryWindowDropMessages(r)
				if len(msgs) != 2 {
					t.Fatalf("amount acceptance sent %d drops", len(msgs))
				}
				for i, m := range msgs {
					if code := binary.LittleEndian.Uint16(m[1:]); code != uint16(100+(i+1)%count) {
						t.Fatalf("drop%d code%d", i, code)
					}
				}
				if *o.windowWords["dword_5d4594_1319268"] != 0 {
					t.Fatal("accepted amount dialog stayed open")
				}
			}
		}
	}
	inventoryWindowCapture(t, "drag-drop", rows, "41ca0fa4c4a5a72d1b8b4f4a795d10c396922d940873da843adbd6a4768012f3")
}
func TestClientInventoryWindowAmountCancel(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, count := range []int{2, 32} {
		o.reset(t)
		o.construct(t)
		o.openInput()
		o.initAmount(t)
		o.stack(t, 0, 0, count, "RedApple", 100)
		rows = append(rows, o.capture(t, count, 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0))
		rows = append(rows, o.capture(t, count, 1, 1, txptr(o.mainWindow().C()), 6, inventoryWindowPoint(600, 400), 0))
		dialog := legacy.Get_nox_gui_itemAmount_dialog_1319228()
		no := dialog.ChildByID(3605)
		dialog.Func94(&gui.RawEvent{Event: 16391, Arg1: txptr(no.C())})
		r := o.capture(t, count, 2, 7, 0, 0, 0, 0)
		rows = append(rows, r)
		if len(inventoryWindowDropMessages(r)) != 0 || legacy.PortTestInventoryCells()[0].Count != byte(count) {
			t.Fatal("cancel dropped or lost stack items")
		}
	}
	inventoryWindowCapture(t, "amount-cancel", rows, "56063ddc15745b1644cd4f1321eaa0a476661b197ce8a728e25de2a7d3d54e32")
}
func TestClientInventoryWindowDropQuantity(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	for _, count := range []int{1, 2, 32} {
		for _, amount := range []int{0, 1, count} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			dr := o.stack(t, 0, 0, count, "RedApple", 200)
			pos[0], pos[1] = 123, 234
			r := o.capture(t, count*100+amount, 0, 13, txptr(unsafe.Pointer(&pos[0])), 200, uintptr(dr.TypeIDVal), uintptr(amount))
			rows = append(rows, r)
			if len(inventoryWindowDropMessages(r)) != amount {
				t.Fatalf("requested%d drops, got%d", amount, len(inventoryWindowDropMessages(r)))
			}
		}
	}
	inventoryWindowCapture(t, "drop-quantity", rows, "cba7af4e3f468b441803c49b41579840f6a9783ec075c50f465c4163ebe980c0")
}
func TestClientInventoryWindowMissingIdentifyResource(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.missingIdentify = true
	r := o.capture(t, 0, 0, 15, 0, 0, 0, 0)
	if r.Return != 0 || *o.windowWords["dword_5d4594_1062476"] != 0 {
		t.Fatal("missing resource did not stop inventory construction")
	}
	if len(o.objects) != 0 {
		t.Fatal("failed constructor created hidden-row items")
	}
	inventoryWindowCapture(t, "missing-resource", []inventoryWindowResult{r}, "91e9fc3bb7b0fbc9f5959ed25450c0ef1806185488d0c0b3fb02003754adb27b")
}

func TestClientInventoryWindowDragAllocationFailure(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.construct(t)
	o.openInput()
	dr := o.stack(t, 0, 0, 3, "RedApple", 100)
	o.exhausted(t, func() {
		if got := o.invoke(1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0); got != 1 {
			t.Fatalf("press return%d", got)
		}
	})
	cell := &legacy.PortTestInventoryCells()[0]
	if cell.Drawable != dr || cell.Count != 3 || cell.Codes[0] != 100 || *memmap.PtrUint32(0x5D4594, 1049848) != 0 || o.c.dragndropItem != nil {
		t.Fatal("failed allocation changed stack or began drag")
	}
	if len(o.console) == 0 {
		t.Fatal("failed drag allocation was not reported")
	}
	inventoryWindowCapture(t, "drag-allocation", []inventoryWindowResult{o.capture(t, 0, 0, 7, 0, 0, 0, 0)}, "a6b60e6bfad294d4a4c7118b3f9de00b17b3a1c65aec77608afb2830c679ad4d")
}
