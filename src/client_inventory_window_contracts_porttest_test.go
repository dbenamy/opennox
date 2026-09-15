//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientInventoryWindowEquipmentRegions(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	point, free := alloc.Make([]uint32{}, 2)
	defer free()
	// Every valid integer paperdoll position is covered even when slot0 is empty.
	for y := 15; y <= 214; y++ {
		for x := 11; x <= 210; x++ {
			point[0], point[1] = uint32(x), uint32(y)
			got := int32(o.invoke(11, txptr(unsafe.Pointer(&point[0])), 0, 0, 0))
			if got < 1 || got > 8 {
				t.Fatalf("uncovered paperdoll point %d,%d: slot%d", x, y, got)
			}
		}
	}
	for _, p := range [][2]uint32{{0xffffffff, 0xffffffff}, {1000, 1000}, {0, 0}} {
		copy(point, p[:])
		if got := int32(o.invoke(11, txptr(unsafe.Pointer(&point[0])), 0, 0, 0)); got != -1 {
			t.Fatalf("outside point %v: slot%d", p, got)
		}
	}
}

func TestClientInventoryWindowTradeIdentifyPredicate(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	prior := o.item(t, "RedApple", 11)
	selected := o.item(t, "Meat", 22)
	*o.windowWords["dword_5d4594_1098624"] = 1
	*o.windowWords["dword_5d4594_1098628"] = 2
	*o.windowWords["dword_5d4594_1049864"] = 5
	*memmap.PtrUint8(0x5D4594, 1049868) = 2
	// Main window has absolute position(4,6); the trade predicate takes local coordinates.
	rect := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 1098380)), 4)
	copy(rect, []uint32{700, 400, 799, 499})
	cell := unsafe.Slice((*uint32)(memmap.PtrOff(0x5D4594, 1098636)), 35)
	cell[0], cell[1], cell[2] = uint32(uintptr(selected.C())), 1, 99
	for _, inside := range []bool{false, true, false} {
		*o.windowWords["dword_5d4594_1063116"] = uint32(uintptr(prior.C()))
		x, y := 699, 399
		if inside {
			x, y = 700, 400
		}
		if got := o.invoke(1, txptr(o.parent.C()), 5, inventoryWindowPoint(x+4, y+6), 0); got != 1 {
			t.Fatalf("identify event return%d", got)
		}
		want := uint32(uintptr(prior.C()))
		if inside {
			want = uint32(uintptr(selected.C()))
		}
		if got := *o.windowWords["dword_5d4594_1063116"]; got != want {
			t.Fatalf("inside=%v selected %#x want %#x", inside, got, want)
		}
		if inside && selected.NetCode32 != 99 {
			t.Fatal("trade selection omitted stack code")
		}
	}
}

func TestClientInventoryWindowEquipmentPriority(t *testing.T) {
	o := newInventoryWindowOwner(t)
	point, free := alloc.Make([]uint32{}, 2)
	defer free()
	for _, slotZero := range []bool{false, true} {
		o.reset(t)
		if slotZero {
			o.equipment[0] = uint32(txptr(o.item(t, "Quiver", 100).C()))
		}
		point[0], point[1] = 80, 50
		want := uint32(7)
		if slotZero {
			want = 0
		}
		if got := o.invoke(11, txptr(unsafe.Pointer(&point[0])), 0, 0, 0); got != want {
			t.Fatalf("slot0=%v got%d want%d", slotZero, got, want)
		}
	}
	for _, shield := range []bool{false, true} {
		o.reset(t)
		first := o.item(t, "RedApple", 100)
		second := o.item(t, "Meat", 101)
		o.equipment[8] = uint32(txptr(first.C()))
		*txword(first, 368) = uint32(txptr(second.C()))
		if shield {
			*txword(second, 112) = 0x2000000
			*txword(second, 116) = 2
		}
		point[0], point[1] = 150, 65
		want := uint32(5)
		if shield {
			want = 8
		}
		if got := o.invoke(11, txptr(unsafe.Pointer(&point[0])), 0, 0, 0); got != want {
			t.Fatalf("shield=%v got%d want%d", shield, got, want)
		}
	}
}
