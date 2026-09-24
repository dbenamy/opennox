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

func nox_xxx_journalEntryAdd_427490(p *C.nox_playerInfo, name *C.char, flags C.short) *C.nox_playerInfo_journal {
	return (*C.nox_playerInfo_journal)(unsafe.Pointer(journalAdd((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags))))
}

func nox_xxx_journalEntryRemove_427590(p *C.nox_playerInfo, name *C.journal_const_char) C.int {
	return C.int(journalRemove((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

func nox_xxx_journalUpdateEntry_4276B0(p *C.nox_playerInfo, name *C.journal_const_char, flags C.short) C.int {
	return C.int(uintptr(unsafe.Pointer(journalUpdate((*server.Player)(unsafe.Pointer(p)), alloc.GoString((*byte)(unsafe.Pointer(name))), uint16(flags)))))
}

func nox_xxx_cliBuildJournalString_469BC0() { journalMeasure() }
