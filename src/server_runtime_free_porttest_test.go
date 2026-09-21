//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerRuntimeRejectedListFree(t *testing.T) {
	type row struct {
		Count, Order  int
		ReturnedFirst bool
		Releases      []uint32
	}
	var rows []row
	for _, count := range []int{0, 1, 2, 3, 17, 256, 1024} {
		for order := 0; order < 3; order++ {
			head := legacy.PortTestResourceAllocate(32)
			if head == nil {
				t.Fatal("list head allocation")
			}
			words := unsafe.Slice((*uint32)(head), 8)
			headAddr := uint32(uintptr(head))
			words[0] = headAddr
			words[1] = headAddr
			words[2] = headAddr
			for i := 3; i < 8; i++ {
				words[i] = 0xa5a5a5a5
			}
			var nodes []unsafe.Pointer
			for i := 0; i < count; i++ {
				p := legacy.PortTestResourceAllocate(0x20c)
				if p == nil {
					t.Fatal("rejected line allocation")
				}
				nodes = append(nodes, p)
			}
			var traversal []int
			for i := 0; i < count; i++ {
				traversal = append(traversal, i)
			}
			if order == 1 {
				for a, b := 0, count-1; a < b; a, b = a+1, b-1 {
					traversal[a], traversal[b] = traversal[b], traversal[a]
				}
			}
			if order == 2 && count > 1 {
				traversal = append(traversal[count/2:], traversal[:count/2]...)
			}
			previous := head
			for _, index := range traversal {
				p := nodes[index]
				w := unsafe.Slice((*uint32)(p), 3)
				w[0] = headAddr
				w[1] = uint32(uintptr(previous))
				w[2] = uint32(index * 13)
				*(*uint32)(previous) = uint32(uintptr(p))
				words[1] = uint32(uintptr(p))
				previous = p
			}
			first := uintptr(words[0])
			if count == 0 {
				first = 0
			}
			var returned uintptr
			releases := legacy.PortTestObserveResourceFrees(nodes, func() { returned = legacy.PortTestRuntimeRejectedClear(head) })
			if returned != first || len(releases) != count {
				t.Fatal("rejected list clear result", count, order, returned, first, len(releases))
			}
			for i, index := range traversal {
				if releases[i] != uint32(index+1) {
					t.Fatal("rejected line release order", count, order, i)
				}
			}
			if words[0] != headAddr || words[1] != headAddr || words[2] != headAddr {
				t.Fatal("list sentinel was not preserved")
			}
			for _, v := range words[3:] {
				if v != 0xa5a5a5a5 {
					t.Fatal("list head guard")
				}
			}
			if legacy.PortTestRuntimeRejectedClear(head) != 0 {
				t.Fatal("repeated empty list clear")
			}
			rows = append(rows, row{count, order, returned == first, releases})
			legacy.PortTestResourceRelease(head)
		}
	}
	interactionCapture(t, "server-runtime-rejected-list-free", rows)
}
