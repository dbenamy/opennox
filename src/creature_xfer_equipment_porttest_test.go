//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestCreatureXferEquipmentOrder(t *testing.T) {
	s := newCreatureXferOwner(t)
	var rows []struct {
		Case  string
		Flags []uint32
	}
	defer func() { spellbookCapture(t, "creature-xfer-equipment-order", rows, "") }()
	kinds := []struct {
		class, subclass uint32
		group           int
		equipped        bool
	}{
		{0x1000000, 4, 1, true}, {0x1000, 0x400, 1, true}, {0x2000000, 2, 2, true}, {0x1000000, 1, 0, true}, {0x1000000, 4, 1, false},
	}
	for length, total := 0, 1; length <= 4; length, total = length+1, total*len(kinds) {
		for sequence := 0; sequence < total; sequence++ {
			t.Run(fmt.Sprintf("len%d-seq%d", length, sequence), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "NPC")
				var items []*server.Object
				var before [][193]uint32
				var expected []uint32
				first, code := 0, sequence
				for i := 0; i < length; i++ {
					kind := kinds[code%len(kinds)]
					code /= len(kinds)
					it := newItemXferObject(t, s, "Gold")
					it.ObjClass = object.Class(kind.class)
					it.ObjSubClass = object.SubClass(kind.subclass)
					it.ObjFlags = 0x400000
					if kind.equipped {
						it.ObjFlags |= 0x100
					}
					flags := uint32(it.ObjFlags)
					if kind.equipped && kind.group != 0 {
						if first == 0 {
							first = kind.group
						} else if first != kind.group {
							flags &^= 0x100
						}
					}
					items = append(items, it)
					expected = append(expected, flags)
				}
				for i, it := range items {
					if i+1 < len(items) {
						it.InvNextItem = items[i+1]
					}
					before = append(before, *(*[193]uint32)(unsafe.Pointer(it)))
				}
				if len(items) != 0 {
					u.InvFirstItem = items[0]
				}
				t.Cleanup(func() {
					u.InvFirstItem = nil
					for _, it := range items {
						it.InvNextItem = nil
					}
				})
				if ret := legacy.PortTestCreatureXferHelper(5, u, nil, 0); ret != 0 {
					t.Fatalf("result=%x", ret)
				}
				for i, it := range items {
					got := *(*[193]uint32)(unsafe.Pointer(it))
					want := before[i]
					want[4] = expected[i]
					if got != want {
						t.Fatalf("item%d flags or unrelated state", i)
					}
				}
				rows = append(rows, struct {
					Case  string
					Flags []uint32
				}{t.Name(), expected})
			})
		}
	}
	if ret := legacy.PortTestCreatureXferHelper(5, nil, nil, 0); ret != 0 {
		t.Fatal("nil inventory owner")
	}
}
