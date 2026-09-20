package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

//export sub_455A00
func sub_455A00(show C.int) C.int { return C.int(teamUIHUDShow(false, int(show))) }

func sub_455A50(count C.char) C.char { return C.char(teamUICTFOpen(byte(count))) }

func sub_455C10() C.int { return C.int(teamUIHUDHide(false)) }

func sub_455D80(id C.uchar, state C.char) { teamUICTFTooltip(byte(id), byte(state)) }

//export sub_455F10
func sub_455F10(show C.int) C.int { return C.int(teamUIHUDShow(true, int(show))) }

func sub_455F60() C.int { return C.int(teamUIBallOpen()) }

func sub_456050() C.int { return C.int(teamUIHUDHide(true)) }

//export nox_xxx_guiServerPlayersLoad_456270
func nox_xxx_guiServerPlayersLoad_456270(parent C.int) C.int {
	return C.int(teamUIPlayersConstruct((*gui.Window)(unsafe.Pointer(uintptr(uint32(parent))))))
}

//export sub_456D60
func sub_456D60(destroy C.int) *C.int {
	return (*C.int)(unsafe.Pointer(teamUIPlayersDestroy(destroy != 0)))
}

func sub_456EA0(name *C.wchar2_t) C.int {
	return C.int(teamUITeamRemove(alloc.GoString16((*uint16)(unsafe.Pointer(name)))))
}

func sub_456FA0() C.int { return C.int(teamUITeamClear()) }

func sub_4571A0(code, id C.int) C.int { return C.int(teamUIPlayerTeam(int(code), int(id))) }

func sub_457230(name *C.wchar2_t) *C.char {
	return (*C.char)(unsafe.Pointer(uintptr(teamUITeamAdd(alloc.GoString16((*uint16)(unsafe.Pointer(name)))))))
}

func sub_4573A0() C.int { return C.int(teamUIPlayersRefreshOpen()) }
