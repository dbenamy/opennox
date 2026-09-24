package legacy

/*
#include "GAME1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func nox_xxx_netNeedTimestampStatus_4174F0(pl *C.nox_playerInfo, mask C.int) C.int {
	return C.int(playerStateAddStatus((*server.Player)(unsafe.Pointer(pl)), uint32(mask)))
}

func nox_xxx_playerUnsetStatus_417530(pl *C.nox_playerInfo, mask C.int) C.char {
	return C.char(playerStateRemoveStatus((*server.Player)(unsafe.Pointer(pl)), uint32(mask)))
}

func nox_xxx_cliPlayerRespawn_417680(pl C.int, mask C.char) {
	playerStateRespawn((*server.Player)(unsafe.Pointer(uintptr(uint32(pl)))), byte(mask))
}
