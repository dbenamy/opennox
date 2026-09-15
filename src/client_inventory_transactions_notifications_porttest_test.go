//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestClientInventoryTransactionsNotificationMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	for _, name := range []string{"Bow", "Quiver", "RedApple"} {
		for _, mode := range []int{0, 1, 2, 3, 4, 5, 6} {
			o.reset(t)
			dr := o.item(t, name, 123)
			if name == "Quiver" {
				dr.ObjSubClass = 2
				o.c.Things.TypeByID(name).ObjSubClass = 2
			}
			if mode != 0 && mode != 2 {
				cell := &legacy.PortTestInventoryCells()[45]
				cell.Drawable, cell.Count, cell.Codes[0] = dr, 2, 123
				cell.Codes[1] = 124
				if mode == 3 {
					cell.Alternate = 1
					*o.words[2] = o.cell(45)
				}
			}
			if mode == 2 {
				*memmap.PtrUint32(0x5D4594, 1049848) = uint32(uintptr(dr.C()))
			}
			if mode == 4 {
				dr.ObjClass = 0x2000000
				dr.ObjSubClass = 0x100
			}
			ret := o.call(10, 123, 0, 0, 0)
			rows = append(rows, o.snapshot(t, len(rows), 10, ret))
			o.checkEquipment(t)
			if mode == 5 || mode == 6 {
				for _, p := range o.equipment {
					if p != 0 {
						*o.txwords[1] = p
						break
					}
				}
				if mode == 6 {
					alt := o.item(t, "Bow", 222)
					cell := &legacy.PortTestInventoryCells()[65]
					cell.Drawable, cell.Count, cell.Codes[0], cell.Alternate = alt, 1, 222, 1
					*o.words[2] = o.cell(65)
				}
			}
			ret = o.call(14, 123, 0, 0, 0)
			// Production drawable deletion returns the remaining object count.
			rows = append(rows, o.snapshot(t, len(rows), 14, ret))
			o.checkEquipment(t)
			ret = o.call(6, 123, 0, 0, 0)
			rows = append(rows, o.snapshot(t, len(rows), 6, ret))
			o.checkEquipment(t)
		}
	}
	inventoryTransactionCapture(t, "notifications", rows, "14ab15f14b260d50914914fc502f872fddcb8cf583b04b172d8f299478c5b50d")
}

func TestClientInventoryTransactionsSecondaryMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	for _, index := range []int{0, 19, 20, 21, 42, 63, 82, 83} {
		for _, previous := range []bool{false, true} {
			o.reset(t)
			dr := o.item(t, "Bow", 123)
			cell := &legacy.PortTestInventoryCells()[index]
			cell.Drawable, cell.Count, cell.Codes[0] = dr, 2, 321
			cell.Codes[1] = 322
			if previous {
				p := o.item(t, "Bow", 500)
				old := &legacy.PortTestInventoryCells()[(index+1)%84]
				old.Drawable, old.Count, old.Codes[0], old.Alternate = p, 1, 501, 1
				*o.words[2] = o.cell((index + 1) % 84)
			}
			ret := o.call(0, uintptr(o.cell(index)), 0, 0, 0)
			rows = append(rows, o.snapshot(t, len(rows), 0, ret))
			if cell.Alternate != 1 || dr.NetCode32 != 321 || *o.words[2] != o.cell(index) {
				t.Fatal("alternate selection must update drawable code and selected cell")
			}
			for i, c := range legacy.PortTestInventoryCells() {
				if i != index && c.Alternate != 0 {
					t.Fatal("old alternate flag survived selection")
				}
			}
			ret = o.call(5, 0, 0, 0, 0)
			delta := ret - o.cell(0)
			if delta != 15532 {
				t.Fatalf("clear-alt return offset%d", delta)
			}
			rows = append(rows, o.snapshot(t, len(rows), 5, 0xea800000+delta))
			ret = o.call(0, 0, 0, 0, 0)
			rows = append(rows, o.snapshot(t, len(rows), 0, ret))
		}
	}
	for mode := 0; mode < 8; mode++ {
		o.reset(t)
		player := o.item(t, "RedApple", 11)
		if mode != 0 {
			*memmap.PtrUint32(0x852978, 8) = uint32(uintptr(player.C()))
		}
		if mode == 2 {
			o.tx.Nox_client_setCursorType(1)
		} else {
			o.tx.Nox_client_setCursorType(0)
		}
		if mode == 3 {
			*memmap.PtrUint32(0x5D4594, 1047764+24) = 1
			*memmap.PtrUint32(0x5D4594, 1047764+24+12) = 2
		}
		if mode == 4 {
			*o.meters.NamedWord("dword_8531A0_2576") = 0
		}
		if mode == 5 {
			*txword(player, 276) = 34
		}
		if mode == 6 || mode == 7 {
			alt := o.item(t, "Bow", 123)
			cell := &legacy.PortTestInventoryCells()[21]
			cell.Drawable, cell.Count, cell.Codes[0], cell.Alternate = alt, 1, 321, 1
			*o.words[2] = o.cell(21)
			if mode == 7 {
				equip := o.item(t, "Bow", 100)
				o.equipment[7] = uint32(uintptr(equip.C()))
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4)) = 4
			}
		}
		ret := o.call(24, 0, 0, 0, 0)
		rows = append(rows, o.snapshot(t, len(rows), 24, ret))
	}
	inventoryTransactionCapture(t, "secondary", rows, "64d40c0dd63c7367db630a9013777e7335667316b354d0fb985bbf84f4ead587")
}

func TestClientInventoryTransactionsDragMatrix(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	for _, index := range []int{0, 19, 20, 21, 42, 63, 82, 83} {
		for _, n := range []byte{1, 2, 31, 32} {
			for _, op := range []int{20, 22} {
				o.reset(t)
				dr := o.item(t, "Bow", 0)
				cell := &legacy.PortTestInventoryCells()[index]
				cell.Drawable, cell.Count = dr, n
				for i := 0; i < int(n); i++ {
					cell.Codes[i] = uint32(123 + i)
				}
				*txword(dr, 292) = 0x43210123
				*txword(dr, 448) = 0x00110008
				*o.txwords[8], *o.txwords[9] = uint32(index/21), uint32(index%21)
				ret := o.call(op, uintptr(index/21), uintptr(index%21), 0, 0)
				rows = append(rows, o.snapshot(t, len(rows), op, ret))
				want := int(n) - 1
				if op == 22 {
					want = int(n)
				}
				if int(cell.Count) != want {
					t.Fatalf("drag/place count%d want%d", cell.Count, want)
				}
			}
		}
	}
	for _, col := range []int{-2, -1, 0, 1, 3, 4, 5, 32767} {
		for _, row := range []int{-2, -1, 0, 19, 20, 21, 22, 32767} {
			o.reset(t)
			ret := o.call(16, uintptr(col), uintptr(row), 0, 0)
			want := uint32(0)
			if col >= 0 && col < 4 && row >= 0 && row < 21 {
				want = 1
			}
			if ret != want {
				t.Fatalf("cell bounds %d,%d got%d", col, row, ret)
			}
			rows = append(rows, o.snapshot(t, len(rows), 16, ret))
		}
	}
	// Non-stackable objects reject a second item even when the type matches.
	o.reset(t)
	dr := o.item(t, "RedApple", 123)
	dr.ObjClass = object.Class(0x4000000)
	cell := &legacy.PortTestInventoryCells()[0]
	cell.Drawable, cell.Count, cell.Codes[0] = dr, 1, 122
	ret := o.call(15, txptr(dr.C()), 0, 0, 0)
	if ret != 0 || cell.Count != 1 {
		t.Fatal("non-stackable placement")
	}
	rows = append(rows, o.snapshot(t, len(rows), 15, ret))
	inventoryTransactionCapture(t, "drag", rows, "a6ff117cb2201610b92b24cd0e7940e4b5fd648344dc3f2547c25578414f4ca3")
}
