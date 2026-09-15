//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestClientInventoryTransactionsPickupTransitions(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	var rows []inventoryTransactionResult
	mods, free := alloc.Make([]uint32{}, 6)
	defer free()
	coords, freeCoords := alloc.Make([]int32{}, 2)
	defer freeCoords()
	mask := (*byte)(unsafe.Add(unsafe.Pointer(o.weapon), 62))
	oldMask := *mask
	defer func() { *mask = oldMask }()
	for _, class := range []byte{0, 1, 2, 7, 8, 31, 32, 255} {
		for _, allowed := range []byte{0, 1, 2, 4, 128, 255} {
			for _, existingAlt := range []bool{false, true} {
				o.reset(t)
				*mask = allowed
				*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = class
				eq := o.item(t, "Bow", 100)
				o.equipment[7] = uint32(uintptr(eq.C()))
				if existingAlt {
					dr := o.item(t, "Bow", 111)
					cell := &legacy.PortTestInventoryCells()[21]
					cell.Drawable, cell.Codes[0], cell.Count, cell.Alternate = dr, 111, 1, 1
					*o.words[2] = o.cell(21)
				}
				ret := o.call(3, 222, uintptr(eq.TypeIDVal), txptr(unsafe.Pointer(&mods[0])), txptr(unsafe.Pointer(&coords[0])))
				selected := *o.words[2]
				want := uint32(0)
				if existingAlt {
					want = o.cell(21)
				} else if byte(uint32(1)<<(class&31))&allowed != 0 {
					want = o.cell(0)
				}
				if ret != 1 || selected != want {
					t.Fatalf("class%d allowed%x existing%t selected%x want%x", class, allowed, existingAlt, selected, want)
				}
				rows = append(rows, o.snapshot(t, len(rows), 3, ret))
			}
		}
	}
	// Follow a dequip notification through the new equip notification and restore
	// the former weapon as secondary, using the actual pending-code fields.
	o.reset(t)
	first := o.item(t, "Bow", 100)
	second := o.item(t, "Bow", 200)
	a, b := &legacy.PortTestInventoryCells()[0], &legacy.PortTestInventoryCells()[21]
	a.Drawable, a.Codes[0], a.Count = first, 100, 1
	b.Drawable, b.Codes[0], b.Count = second, 200, 1
	o.call(10, 100, 0, 0, 0)
	o.call(0, uintptr(o.cell(21)), 0, 0, 0)
	equipped := o.equipment[8]
	if equipped == 0 {
		t.Fatal("bow was not equipped in slot8")
	}
	*o.txwords[1] = equipped
	o.call(14, 100, 0, 0, 0)
	rows = append(rows, o.snapshot(t, len(rows), 14, 0))
	if *o.txwords[2] != 100 {
		t.Fatal("previous weapon code not queued")
	}
	o.call(10, 200, 0, 0, 0)
	rows = append(rows, o.snapshot(t, len(rows), 10, 0))
	if *o.txwords[2] != 0 || *o.words[2] != o.cell(0) || a.Alternate != 1 || b.Alternate != 0 {
		t.Fatal("secondary weapon follow-up did not complete")
	}
	// A selected stack moved during compaction keeps its selected identity.
	o.reset(t)
	dr := o.item(t, "Bow", 321)
	cell := &legacy.PortTestInventoryCells()[67]
	cell.Drawable, cell.Count, cell.Codes[0], cell.Alternate = dr, 1, 321, 1
	*o.words[2] = o.cell(67)
	ret := o.call(7, 0, 0, 0, 0)
	delta := ret - o.cell(0)
	rows = append(rows, o.snapshot(t, len(rows), 7, 0xea800000+delta))
	if *o.words[2] != o.cell(0) || legacy.PortTestInventoryCells()[0].Codes[0] != 321 {
		t.Fatal("compaction lost selected cell identity")
	}
	// A completely occupied visible grid and an empty grid exercise both return paths.
	for _, full := range []bool{false, true} {
		o.reset(t)
		if full {
			dr := o.item(t, "RedApple", 123)
			for i := range legacy.PortTestInventoryCells() {
				cell := &legacy.PortTestInventoryCells()[i]
				cell.Drawable, cell.Count, cell.Codes[0] = dr, 1, uint32(i+1)
			}
		}
		ret := o.call(7, 0, 0, 0, 0)
		delta := ret - o.cell(0)
		want := uint32(15380)
		if full {
			want = 3096
		}
		if delta != want {
			t.Fatalf("full%t compaction return%d want%d", full, delta, want)
		}
		rows = append(rows, o.snapshot(t, len(rows), 7, 0xea800000+delta))
	}
	// Cursor clearing and the actual secondary-change sound are observable.
	o.reset(t)
	player := o.item(t, "RedApple", 1)
	alt := o.item(t, "Bow", 2)
	o.c.dragndropItem = alt
	*memmap.PtrUint32(0x852978, 8) = uint32(uintptr(player.C()))
	cell = &legacy.PortTestInventoryCells()[0]
	cell.Drawable, cell.Count, cell.Codes[0] = alt, 1, 2
	*o.words[2] = o.cell(0)
	rect := unsafe.Slice((*int32)(memmap.PtrOff(0x587000, 136336)), 4)
	*memmap.PtrInt32(0x5D4594, 1062572) = rect[0]
	*memmap.PtrInt32(0x5D4594, 1062576) = rect[1]
	o.call(24, 0, 0, 0, 0)
	if o.c.dragndropItem != nil || !reflect.DeepEqual(o.sounds, [][2]int{{895, 100}}) {
		t.Fatalf("secondary cursor/sound %p %v", o.c.dragndropItem, o.sounds)
	}
	rows = append(rows, o.snapshot(t, len(rows), 24, 0))
	inventoryTransactionCapture(t, "transitions", rows, "e58201bbcf8e6d2a8f7fca122d2fe47f9bdbf74957b66eb90d0c0fb99b22493d")
}
