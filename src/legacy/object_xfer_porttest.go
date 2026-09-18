//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "GAME4_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// The C object reader allocates names with libc calloc, outside alloc's tracker.
func PortTestObjectXferFreeName(p unsafe.Pointer) { C.free(p) }

// Preserve globals while using the actual pending-ownership allocator/list.
func PortTestObjectXferPendingOwners() (func() [][2]uint32, func()) {
	oldPool, oldHead := monsterPendingClass, monsterPendingHead
	if monsterPendingInit() == 0 {
		panic("pending ownership allocation")
	}
	read := func() (rows [][2]uint32) {
		for p := unsafe.Pointer(uintptr(monsterPendingHead)); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 8)) {
			rows = append(rows, *(*[2]uint32)(p))
		}
		return
	}
	return read, func() { monsterPendingFree(); monsterPendingClass = oldPool; monsterPendingHead = oldHead }
}

// The editor list owns its libc list nodes, while the fixture owns its objects.
func PortTestObjectXferEditorList() (func() []*server.Object, func()) {
	old := Get_dword_5d4594_1599540()
	Set_dword_5d4594_1599540(nil)
	read := func() (out []*server.Object) {
		for p := Get_dword_5d4594_1599540(); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 4)) {
			out = append(out, *(**server.Object)(p))
		}
		return
	}
	return read, func() {
		for p := Get_dword_5d4594_1599540(); p != nil; {
			next := *(*unsafe.Pointer)(unsafe.Add(p, 4))
			C.free(p)
			p = next
		}
		Set_dword_5d4594_1599540(old)
	}
}
