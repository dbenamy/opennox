package legacy

/*
#include "GAME1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_40A770
func sub_40A770() C.int { return C.int(playerStateCompetitors()) }

//export nox_xxx_gamePlayIsAnyPlayers_40A8A0
func nox_xxx_gamePlayIsAnyPlayers_40A8A0() C.int { return C.int(playerStateMultiple()) }

//export sub_40A970
func sub_40A970() { playerStateReset() }

//export sub_40AA70
func sub_40AA70(pl *C.nox_playerInfo) C.int {
	return C.int(playerStateAdmission((*server.Player)(unsafe.Pointer(pl))))
}

//export nox_xxx_playerForceSendLessons_416E50
func nox_xxx_playerForceSendLessons_416E50(send C.int) *C.char {
	playerStateLessons(int32(send))
	return nil
}

//export nox_xxx_netNeedTimestampStatus_4174F0
func nox_xxx_netNeedTimestampStatus_4174F0(pl *C.nox_playerInfo, mask C.int) C.int {
	return C.int(playerStateAddStatus((*server.Player)(unsafe.Pointer(pl)), uint32(mask)))
}

//export nox_xxx_playerUnsetStatus_417530
func nox_xxx_playerUnsetStatus_417530(pl *C.nox_playerInfo, mask C.int) C.char {
	return C.char(playerStateRemoveStatus((*server.Player)(unsafe.Pointer(pl)), uint32(mask)))
}

//export nox_xxx_sendAllClientStatus_4175C0
func nox_xxx_sendAllClientStatus_4175C0(to C.int) *C.char {
	playerStateAllStatus((*server.Player)(unsafe.Pointer(uintptr(uint32(to)))))
	return nil
}

//export nox_xxx_cliPlayerRespawn_417680
func nox_xxx_cliPlayerRespawn_417680(pl C.int, mask C.char) {
	playerStateRespawn((*server.Player)(unsafe.Pointer(uintptr(uint32(pl)))), byte(mask))
}

//export nox_xxx_clientEquipWeaponArmor_417AA0
func nox_xxx_clientEquipWeaponArmor_417AA0(command C.char, code, mask, mods C.int) *C.char {
	return (*C.char)(playerStateEquip(byte(command), uint32(code), uint32(mask), (*[4]byte)(unsafe.Pointer(uintptr(uint32(mods))))))
}

//export sub_417B80
func sub_417B80(command C.char, code, mask C.int) *C.char {
	return (*C.char)(playerStateUnequip(byte(command), uint32(code), uint32(mask)))
}
