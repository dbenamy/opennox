package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func sub_506870(kind C.int, player C.int, name *C.wchar2_t) {
	voteCast(uint32(kind), (*server.Object)(unsafe.Pointer(uintptr(uint32(player)))), (*uint16)(unsafe.Pointer(name)))
}

func sub_506C90(kind C.int, player C.int, name *C.wchar2_t) {
	voteWithdraw(uint32(kind), (*server.Object)(unsafe.Pointer(uintptr(uint32(player)))), (*uint16)(unsafe.Pointer(name)))
}

func sub_48D4A0() C.int { return C.int(voteGUIReset()) }
