package legacy

/*
#include "GAME1_1.h"
#include "client__gui__guijourn.h"
typedef const char journal_const_char;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_journalEntryAdd_427490
func nox_xxx_journalEntryAdd_427490(p *C.nox_playerInfo, name *C.char, flags C.short) *C.nox_playerInfo_journal {
	return (*C.nox_playerInfo_journal)(unsafe.Pointer(journalAdd((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags))))
}

//export nox_xxx_comJournalEntryAdd_427500
func nox_xxx_comJournalEntryAdd_427500(u *C.nox_object_t, name *C.char, flags C.short) {
	journalUnitAdd((*server.Object)(unsafe.Pointer(u)), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags))
}

//export nox_xxx_journalEntryRemove_427590
func nox_xxx_journalEntryRemove_427590(p *C.nox_playerInfo, name *C.journal_const_char) C.int {
	return C.int(journalRemove((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

//export nox_xxx_comJournalEntryRemove_427630
func nox_xxx_comJournalEntryRemove_427630(u C.int, name *C.journal_const_char) {
	journalUnitRemove((*server.Object)(unsafe.Pointer(uintptr(uint32(u)))), alloc.GoString((*byte)(unsafe.Pointer(name))))
}

//export nox_xxx_comRemoveEntryAll_427680
func nox_xxx_comRemoveEntryAll_427680(name *C.journal_const_char) C.int {
	return C.int(journalRemoveAll(alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

//export nox_xxx_journalUpdateEntry_4276B0
func nox_xxx_journalUpdateEntry_4276B0(p *C.nox_playerInfo, name *C.journal_const_char, flags C.short) C.int {
	return C.int(uintptr(unsafe.Pointer(journalUpdate((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags)))))
}

//export nox_xxx_comJournalEntryUpdate_427720
func nox_xxx_comJournalEntryUpdate_427720(u C.int, name *C.journal_const_char, flags C.short) C.int {
	return C.int(journalUnitUpdate((*server.Object)(unsafe.Pointer(uintptr(uint32(u)))), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags)))
}

//export nox_xxx_comUpdateEntryAll_427770
func nox_xxx_comUpdateEntryAll_427770(name *C.journal_const_char, flags C.short) C.int {
	return C.int(journalUpdateAll(alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags)))
}

//export sub_4277B0
func sub_4277B0(u *C.nox_object_t, mask C.ushort) C.int {
	return C.int(journalRemoveMask((*server.Object)(unsafe.Pointer(u)), uint16(mask)))
}

//export nox_xxx_cliBuildJournalString_469BC0
func nox_xxx_cliBuildJournalString_469BC0() { journalMeasure() }
