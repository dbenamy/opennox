//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "client__gui__guijourn.h"
*/
import "C"
import "unsafe"

// PortTestJournal calls the real journal entry owners and renderer.
func PortTestJournal(op int, a, b, c uintptr) uint32 {
	switch op {
	case 0:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_journalEntryAdd_427490((*C.nox_playerInfo)(unsafe.Pointer(a)), (*C.char)(unsafe.Pointer(b)), C.short(c)))))
	case 1:
		C.nox_xxx_comJournalEntryAdd_427500((*C.nox_object_t)(unsafe.Pointer(a)), (*C.char)(unsafe.Pointer(b)), C.short(c))
	case 2:
		return uint32(C.nox_xxx_journalEntryRemove_427590((*C.nox_playerInfo)(unsafe.Pointer(a)), (*C.char)(unsafe.Pointer(b))))
	case 3:
		C.nox_xxx_comJournalEntryRemove_427630(C.int(a), (*C.char)(unsafe.Pointer(b)))
	case 4:
		return uint32(C.nox_xxx_comRemoveEntryAll_427680((*C.char)(unsafe.Pointer(b))))
	case 5:
		return uint32(C.nox_xxx_journalUpdateEntry_4276B0((*C.nox_playerInfo)(unsafe.Pointer(a)), (*C.char)(unsafe.Pointer(b)), C.short(c)))
	case 6:
		return uint32(C.nox_xxx_comJournalEntryUpdate_427720(C.int(a), (*C.char)(unsafe.Pointer(b)), C.short(c)))
	case 7:
		return uint32(C.nox_xxx_comUpdateEntryAll_427770((*C.char)(unsafe.Pointer(b)), C.short(c)))
	case 8:
		return uint32(C.sub_4277B0((*C.nox_object_t)(unsafe.Pointer(a)), C.ushort(c)))
	case 9:
		C.nox_xxx_cliBuildJournalString_469BC0()
	case 11:
		return uint32(C.sub_41BEC0(unsafe.Pointer(a), nil))
	case 10:
		C.nox_xxx_guiDrawJournal_469D40(C.int(a), C.int(b), C.int(c))
	}
	return 0
}
