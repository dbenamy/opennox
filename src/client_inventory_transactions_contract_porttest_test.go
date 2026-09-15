//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInventoryTransactionsPickupCoordinates(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	mods, free := alloc.Make([]uint32{}, 6)
	defer free()
	for col := 0; col < 4; col++ {
		for row := 0; row < 20; row++ {
			o.reset(t)
			dr := o.item(t, "Bow", 123)
			cell := &legacy.PortTestInventoryCells()[col*21+row]
			cell.Drawable, cell.Codes[0], cell.Count = dr, 123, 1
			if got := o.call(2, 124, uintptr(dr.TypeIDVal), txptr(unsafe.Pointer(&mods[0])), 0); got != 1 {
				t.Fatalf("pickup %d,%d returned %d", col, row, got)
			}
			if cell.Count != 2 || cell.Codes[1] != 124 {
				t.Fatal("existing-stack append")
			}
			if got := *memmap.PtrInt32(0x5D4594, 1062580); got != int32(col) {
				t.Fatalf("column %d want %d", got, col)
			}
			if got := *memmap.PtrInt32(0x5D4594, 1062584); got != int32(row) {
				t.Fatalf("row %d want %d", got, row)
			}
			o.snapshot(t, col*20+row, 2, 1)
		}
	}
}

func TestClientInventoryTransactionsRemovalOwnership(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	lookup, free := alloc.New(legacy.PortTestInventoryLookup{})
	defer free()
	for _, n := range []int{1, 2, 3, 16, 31, 32} {
		for index := 0; index < n; index++ {
			o.reset(t)
			dr := o.item(t, "RedPotion", 77)
			cell := &legacy.PortTestInventoryCells()[45]
			cell.Drawable, cell.Count = dr, byte(n)
			want := make([]uint32, 0, n-1)
			for i := 0; i < n; i++ {
				cell.Codes[i] = uint32(100 + i)
				if i != index {
					want = append(want, uint32(100+i))
				}
			}
			lookup.Cell, lookup.Index = cell, uint32(index)
			o.call(8, txptr(unsafe.Pointer(lookup)), 0, 0, 0)
			if int(cell.Count) != n-1 || !reflect.DeepEqual(cell.Codes[:n-1], want) {
				t.Fatalf("remove n%d index%d count%d codes%v want%v", n, index, cell.Count, cell.Codes[:n-1], want)
			}
			if (cell.Drawable == nil) != (n == 1) {
				t.Fatal("drawable lifetime must follow last stack code")
			}
			o.snapshot(t, n*32+index, 8, 0)
		}
	}
}

func TestClientInventoryTransactionsPlaceCopies(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	for col := 0; col < 4; col++ {
		for row := 0; row < 21; row++ {
			o.reset(t)
			dr := o.item(t, "Bow", 0x1234)
			// Modifiers are actual owned modifier pointers; charge bytes are values.
			for i := 0; i < 4; i++ {
				*txword(dr, 432+uintptr(i*4)) = uint32(uintptr(o.mods[i].C()))
			}
			*txword(dr, 448) = 0x11223344
			*txword(dr, 452) = 0x55667788
			*txword(dr, 292) = 0x12340056
			if got := o.call(15, txptr(dr.C()), uintptr(col), uintptr(row), 0); got != 1 {
				t.Fatalf("place %d,%d ret%d", col, row, got)
			}
			cell := &legacy.PortTestInventoryCells()[col*21+row]
			if cell.Drawable == dr || cell.Drawable == nil || cell.Count != 1 || cell.Codes[0] != 0x1234 {
				t.Fatal("place must create independent stack drawable")
			}
			if !bytes.Equal(unsafe.Slice((*byte)(unsafe.Add(cell.Drawable.C(), 432)), 24), unsafe.Slice((*byte)(unsafe.Add(dr.C(), 432)), 24)) || *txword(cell.Drawable, 292) != *txword(dr, 292) {
				t.Fatal("place metadata/durability copy")
			}
			o.snapshot(t, col*21+row, 15, 1)
		}
	}
}

func TestClientInventoryTransactionsRequests(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	for _, code := range []uint32{0, 1, 32767, 32768, 65535, 0xffffffff} {
		for _, class := range []uint32{0, 0x400000, 0x20000000} {
			for _, op := range []int{1, 12, 17, 18, 19, 21, 23} {
				o.reset(t)
				dr := o.item(t, "RedApple", code)
				dr.ObjClass = object.Class(class)
				binary.LittleEndian.PutUint32(o.regions[1], uint32(uintptr(dr.C())))
				pos[0], pos[1] = -40000, 70000
				a := txptr(dr.C())
				if op == 1 || op == 18 || op == 19 {
					a = uintptr(code)
				}
				if op == 21 {
					a = txptr(unsafe.Pointer(&pos[0]))
				}
				o.call(op, a, 0, 0, 0)
				wire := uint16(code)
				if op == 12 || op == 17 || op == 21 || op == 23 {
					if code >= 0x8000 {
						wire = 0
					} else if class&0x20400000 != 0 {
						wire |= 0x8000
					}
				}
				opcode := map[int]byte{1: 241, 12: 117, 17: 118, 18: 201, 19: 201, 21: 114, 23: 116}[op]
				want := []byte{opcode, byte(wire), byte(wire >> 8)}
				if op == 18 || op == 19 {
					sub := byte(30)
					if op == 19 {
						sub = 28
					}
					want = []byte{201, sub, byte(wire), byte(wire >> 8)}
				}
				if op == 21 {
					want = append(want, byte(uint32(pos[0])), byte(uint32(pos[0])>>8), byte(uint32(pos[1])), byte(uint32(pos[1])>>8))
				}
				var messages [][]byte
				o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { messages = append(messages, append([]byte(nil), b...)); return false })
				if len(messages) != 1 || !bytes.Equal(messages[0], want) {
					t.Fatalf("op%d code%x class%x messages%x want%x", op, code, class, messages, want)
				}
				o.snapshot(t, op, op, 0)
			}
		}
	}
}

func TestClientInventoryTransactionsAllocationFailure(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	mods, free := alloc.Make([]uint32{}, 6)
	defer free()
	coord, freeCoord := alloc.Make([]int32{}, 2)
	defer freeCoord()
	for _, op := range []int{2, 3, 10, 15} {
		o.reset(t)
		dr := o.item(t, "Bow", 123)
		cell := &legacy.PortTestInventoryCells()[21]
		cell.Drawable, cell.Codes[0], cell.Count = dr, 123, 1
		if op == 2 {
			cell.Count = 32
		}
		o.exhausted(t, func() {
			switch op {
			case 2:
				o.call(op, 124, uintptr(dr.TypeIDVal), txptr(unsafe.Pointer(&mods[0])), 0)
			case 3:
				o.call(op, 124, uintptr(dr.TypeIDVal), txptr(unsafe.Pointer(&mods[0])), txptr(unsafe.Pointer(&coord[0])))
			case 10:
				o.call(op, 123, 0, 0, 0)
			case 15:
				o.call(op, txptr(dr.C()), 0, 0, 0)
			}
		})
		if len(o.console) == 0 {
			t.Fatalf("op%d exhausted pool did not report an error", op)
		}
		if legacy.PortTestInventoryCells()[0].Count != 0 {
			t.Fatal("allocation failure occupied a cell")
		}
		o.snapshot(t, op, op, 0)
	}
}
