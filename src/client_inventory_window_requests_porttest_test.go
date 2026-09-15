//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientInventoryWindowTradeSelection(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	pos, free := alloc.Make([]uint32{}, 2)
	defer free()
	for _, count := range []int{0, 1, 2, 32} {
		for _, repair := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			if count > 0 {
				o.stack(t, 0, 0, count, "RedApple", 100)
			}
			pos[0], pos[1] = 339, 38
			var r inventoryWindowResult
			if repair {
				*o.windowWords["dword_5d4594_1049864"] = 6
				r = o.capture(t, len(rows), 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0)
			} else {
				r = o.capture(t, len(rows), 0, 9, txptr(unsafe.Pointer(&pos[0])), 0, 0, 0)
			}
			rows = append(rows, r)
			var msgs [][]byte
			for _, m := range r.Inventory.Messages {
				if len(m) == 4 && m[0] == 201 {
					msgs = append(msgs, m)
				}
			}
			if count == 0 {
				if len(msgs) != 0 {
					t.Fatal("empty cell sent trade")
				}
				continue
			}
			if len(msgs) != 1 {
				t.Fatalf("trade sent%d messages", len(msgs))
			}
			code, op := uint16(100+count-1), byte(28)
			if repair {
				code, op = 100, 30
			}
			if msgs[0][1] != op || binary.LittleEndian.Uint16(msgs[0][2:]) != code {
				t.Fatalf("trade request %v want op%d code%d", msgs[0], op, code)
			}
		}
	}
	inventoryWindowCapture(t, "trade-selection", rows, "7af0e9ba97bd18f6b8665f5b7252085fe110298249be2da5167d25f80638f58e")
}

func TestClientInventoryWindowCancelFallback(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, equipped := range []bool{false, true} {
		o.reset(t)
		o.construct(t)
		o.openInput()
		o.stack(t, 0, 0, 1, "Bow", 100)
		rows = append(rows, o.capture(t, len(rows), 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0))
		o.stack(t, 0, 0, 32, "Meat", 500)
		if equipped {
			dr := o.item(t, "Bow", 100)
			o.equipment[7] = uint32(txptr(dr.C()))
		}
		rows = append(rows, o.capture(t, len(rows), 1, 33, 0, 0, 0, 0))
		cells := legacy.PortTestInventoryCells()
		if cells[0].Count != 32 || cells[0].Codes[0] != 500 {
			t.Fatal("fallback changed full source cell")
		}
		if *memmap.PtrUint32(0x5D4594, 1049848) != 0 || o.c.dragndropItem != nil {
			t.Fatal("fallback retained drag")
		}
		found := false
		for _, cell := range cells {
			if cell.Count == 1 && cell.Codes[0] == 100 {
				found = true
				if cell.Drawable == nil {
					t.Fatal("restored code has no drawable")
				}
			}
		}
		if !found {
			t.Fatal("fallback lost dragged item")
		}
	}
	inventoryWindowCapture(t, "cancel-fallback", rows, "83117e86628d2a3dfcf6091ab15579236d6ebf1403da861d4950f7ffaf87e58d")
}

func TestClientInventoryWindowDropBlockedAndMissing(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	point, free := alloc.Make([]int32{}, 2)
	defer free()
	point[0], point[1] = 123, 234
	for _, blocked := range []uint32{0, 1, 2} {
		for _, present := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			*o.windowWords["dword_5d4594_1320964"] = blocked
			if present {
				o.stack(t, 0, 0, 2, "RedApple", 100)
			}
			r := o.capture(t, len(rows), 0, 13, txptr(unsafe.Pointer(&point[0])), 100, 0, 2)
			rows = append(rows, r)
			want := 0
			if blocked == 0 && present {
				want = 2
			}
			if len(inventoryWindowDropMessages(r)) != want {
				t.Fatalf("blocked%d present%v drops%d want%d", blocked, present, len(inventoryWindowDropMessages(r)), want)
			}
			if *memmap.PtrUint32(0x5D4594, 1049848) != 0 {
				t.Fatal("quantity callback retained transient drag")
			}
		}
	}
	inventoryWindowCapture(t, "drop-blocked", rows, "7655ff1d6a780f46e879d385ac99b58b17432cafc38a86d988b571abe2fc8776")
}

func TestClientInventoryWindowClickDragThreshold(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, name := range []string{"RedApple", "Bow"} {
		for _, delta := range []int{0, 9, 10, 11} {
			for _, elapsed := range []uint32{0, 9, 10, 11, 30} {
				o.reset(t)
				o.construct(t)
				o.openInput()
				o.stack(t, 0, 0, 2, name, 100)
				o.c.srv.SetFrame(120)
				rows = append(rows, o.capture(t, len(rows), 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0))
				o.c.srv.SetFrame(120 + elapsed)
				rows = append(rows, o.capture(t, len(rows), 1, 1, txptr(o.mainWindow().C()), 6, inventoryWindowPoint(339+delta, 38), 0))
				if *memmap.PtrUint32(0x5D4594, 1049848) != 0 || o.c.dragndropItem != nil {
					t.Fatal("click/drag threshold retained drag")
				}
			}
		}
	}
	inventoryWindowCapture(t, "click-drag-threshold", rows, "5914efec30d81d9bc57370e75a6607e7dff119958b022364ff44c7b963df41e9")
}
