//go:build porttest

package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// PortTestJournal calls the real journal entry owners and renderer.
func PortTestJournal(op int, a, b, c uintptr) uint32 {
	switch op {
	case 0:
		return uint32(uintptr(unsafe.Pointer(journalAdd((*server.Player)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b))), uint16(c)))))
	case 1:
		journalUnitAdd((*server.Object)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b))), uint16(c))
	case 2:
		return uint32(journalRemove((*server.Player)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b)))))
	case 3:
		journalUnitRemove((*server.Object)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b))))
	case 4:
		return uint32(journalRemoveAll(alloc.GoString((*byte)(unsafe.Pointer(b)))))
	case 5:
		return uint32(uintptr(unsafe.Pointer(journalUpdate((*server.Player)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b))), uint16(c)))))
	case 6:
		return journalUnitUpdate((*server.Object)(unsafe.Pointer(a)), alloc.GoString((*byte)(unsafe.Pointer(b))), uint16(c))
	case 7:
		return uint32(journalUpdateAll(alloc.GoString((*byte)(unsafe.Pointer(b))), uint16(c)))
	case 8:
		return uint32(journalRemoveMask((*server.Object)(unsafe.Pointer(a)), uint16(c)))
	case 9:
		journalMeasure()
	case 10:
		journalDraw(int(a), int(b), int(c))
	case 11:
		return uint32(C.sub_41BEC0(unsafe.Pointer(a), nil))
	}
	return 0
}
