package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_506870
func sub_506870(kind C.int, player C.int, name *C.wchar2_t) {
	voteCast(uint32(kind), (*server.Object)(unsafe.Pointer(uintptr(uint32(player)))), (*uint16)(unsafe.Pointer(name)))
}

//export sub_506C90
func sub_506C90(kind C.int, player C.int, name *C.wchar2_t) {
	voteWithdraw(uint32(kind), (*server.Object)(unsafe.Pointer(uintptr(uint32(player)))), (*uint16)(unsafe.Pointer(name)))
}

//export sub_48CAD0
func sub_48CAD0() C.int { return C.int(voteGUIHide()) }

//export sub_48CB10
func sub_48CB10(topic C.int) { voteGUIShow(uint32(topic)) }

//export sub_48D4A0
func sub_48D4A0() C.int { return C.int(voteGUIReset()) }
