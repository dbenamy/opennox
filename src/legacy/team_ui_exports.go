package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func sub_455A50(count C.char) C.char { return C.char(teamUICTFOpen(byte(count))) }

func sub_455C10() C.int { return C.int(teamUIHUDHide(false)) }

func sub_455D80(id C.uchar, state C.char) { teamUICTFTooltip(byte(id), byte(state)) }

func sub_455F60() C.int { return C.int(teamUIBallOpen()) }

func sub_456050() C.int { return C.int(teamUIHUDHide(true)) }

func sub_456EA0(name *C.wchar2_t) C.int {
	return C.int(teamUITeamRemove(alloc.GoString16((*uint16)(unsafe.Pointer(name)))))
}

func sub_456FA0() C.int { return C.int(teamUITeamClear()) }

func sub_4571A0(code, id C.int) C.int { return C.int(teamUIPlayerTeam(int(code), int(id))) }

func sub_457230(name *C.wchar2_t) *C.char {
	return (*C.char)(unsafe.Pointer(uintptr(teamUITeamAdd(alloc.GoString16((*uint16)(unsafe.Pointer(name)))))))
}

func sub_4573A0() C.int { return C.int(teamUIPlayersRefreshOpen()) }
