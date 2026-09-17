//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "defs.h"
#include "GAME4_1.h"
extern void* nox_alloc_pendingOwn_2386916;
extern uint32_t dword_5d4594_2386920;
extern void* dword_5d4594_1599540;
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
	oldPool, oldHead := C.nox_alloc_pendingOwn_2386916, C.dword_5d4594_2386920
	if C.nox_xxx_allocPendingOwnsArray_516EE0() == 0 {
		panic("pending ownership allocation")
	}
	read := func() (rows [][2]uint32) {
		for p := unsafe.Pointer(uintptr(C.dword_5d4594_2386920)); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 8)) {
			rows = append(rows, *(*[2]uint32)(p))
		}
		return
	}
	return read, func() { C.sub_516F10(); C.nox_alloc_pendingOwn_2386916 = oldPool; C.dword_5d4594_2386920 = oldHead }
}

// The editor list owns its libc list nodes, while the fixture owns its objects.
func PortTestObjectXferEditorList() (func() []*server.Object, func()) {
	old := C.dword_5d4594_1599540
	C.dword_5d4594_1599540 = nil
	read := func() (out []*server.Object) {
		for p := C.dword_5d4594_1599540; p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 4)) {
			out = append(out, *(**server.Object)(p))
		}
		return
	}
	return read, func() {
		for p := C.dword_5d4594_1599540; p != nil; {
			next := *(*unsafe.Pointer)(unsafe.Add(p, 4))
			C.free(p)
			p = next
		}
		C.dword_5d4594_1599540 = old
	}
}
