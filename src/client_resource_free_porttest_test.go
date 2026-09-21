//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientResourceFreeGraph(t *testing.T) {
	type row struct {
		Kind, Pattern, Allocations int
		Order                      []uint32
	}
	var rows []row
	for _, kind := range []int{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
		for pattern := 0; pattern < 6; pattern++ {
			t.Run(fmt.Sprintf("kind%d-pattern%d", kind, pattern), func(t *testing.T) {
				var pointers []unsafe.Pointer
				allocate := func(size uintptr) unsafe.Pointer {
					p := legacy.PortTestResourceAllocate(size)
					if p == nil {
						t.Fatal("resource allocation")
					}
					pointers = append(pointers, p)
					return p
				}
				store := func(p unsafe.Pointer, off uintptr, v uint32) { *(*uint32)(unsafe.Add(p, off)) = v }
				include := func(index int) bool {
					switch pattern {
					case 0:
						return false
					case 1:
						return true
					case 2:
						return index%2 == 0
					case 3:
						return index%2 == 1
					case 4:
						return index%7 == 0
					default:
						return index%17 == 16
					}
				}
				// A direction record has eight owned frame lists and a nonpointer middle word.
				vector := func(p unsafe.Pointer, index int) {
					store(p, 16, 0x5a5a5a5a)
					offsets := []uintptr{0, 4, 8, 12, 20, 24, 28, 32}
					for direction, off := range offsets {
						if include(index*8 + direction) {
							store(p, off, uint32(uintptr(allocate(16))))
						}
					}
				}
				rootSize := uintptr(64)
				switch kind {
				case 6:
					rootSize = 55*264 + 52
				case 7:
					rootSize = 16*48 + 8
				case 8:
					rootSize = 3*48 + 8
				}
				root := allocate(rootSize)
				switch kind {
				case 2, 3:
					if include(0) {
						store(root, 4, uint32(uintptr(allocate(16))))
					}
				case 4:
					for i := 0; i < 5; i++ {
						if include(i) {
							store(root, 4+4*uintptr(i), uint32(uintptr(allocate(16))))
						}
					}
				case 5:
					vector(unsafe.Add(root, 4), 0)
				case 6:
					for state := 0; state < 55; state++ {
						for slot := 0; slot < 54; slot++ {
							index := state*54 + slot
							if include(index) {
								group := allocate(48)
								store(root, 52+uintptr(state)*264+uintptr(slot)*4, uint32(uintptr(group)))
								vector(unsafe.Add(group, 4), index)
							}
						}
					}
				case 7, 8:
					groups := 16
					if kind == 8 {
						groups = 3
					}
					for i := 0; i < groups; i++ {
						vector(unsafe.Add(root, 8+48*i), i)
					}
				}
				order := legacy.PortTestObserveResourceFrees(pointers, func() { legacy.Nox_xxx_draw_44C650_free_kind(root, kind) })
				seen := make(map[uint32]bool)
				for _, id := range order {
					if id == 0 || int(id) > len(pointers) || seen[id] {
						t.Fatal("unknown or duplicate release", id)
					}
					seen[id] = true
				}
				for i, p := range pointers {
					if !seen[uint32(i+1)] {
						legacy.PortTestResourceRelease(p)
					}
				}
				if len(order) != len(pointers) {
					t.Fatalf("released %d of %d allocations", len(order), len(pointers))
				}
				if order[len(order)-1] != 1 {
					t.Fatal("root freed before its children")
				}
				rows = append(rows, row{kind, pattern, len(pointers), order})
			})
		}
	}
	interactionCapture(t, "client-resource-free-graph", rows)
}
