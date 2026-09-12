package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_playerLeaveObserver_0_4E6AA0
func nox_xxx_playerLeaveObserver_0_4E6AA0(a *C.nox_playerInfo) {
	controlLeaveObserver(unsafe.Pointer(a))
}

//export nox_xxx_unitRemoveChild_4EC470
func nox_xxx_unitRemoveChild_4EC470(a *C.nox_object_t) { controlRemoveChildren(asObjectS(a)) }

//export nox_xxx_unitTransferSlaves_4EC4B0
func nox_xxx_unitTransferSlaves_4EC4B0(a *C.nox_object_t) { controlTransferChildren(asObjectS(a)) }

//export nox_xxx_plrReadVals_4EEDC0
func nox_xxx_plrReadVals_4EEDC0(a *C.nox_object_t, b C.int) C.int {
	return C.int(controlReadStats(asObjectS(a), int32(b)))
}

//export sub_4EF140
func sub_4EF140(a C.int) C.int { return C.int(controlLevelFromXP(objectFromInt(a))) }

//export nox_xxx_calcBoltDamage_4EF1E0
func nox_xxx_calcBoltDamage_4EF1E0(a, b C.int) C.double {
	return C.double(controlBoltDamage(int32(a), unsafe.Pointer(uintptr(b))))
}

//export sub_4EF410
func sub_4EF410(a C.int, b C.uchar) { controlSetLevel(objectFromInt(a), byte(b)) }

//export sub_4EF6F0
func sub_4EF6F0(a C.int) C.int { return C.int(controlGlyphCount(objectFromInt(a))) }

//export nox_xxx_playerRespawnItem_4EF750
func nox_xxx_playerRespawnItem_4EF750(a *C.nox_object_t, b *C.char, c *C.int, d, e C.int) *C.nox_object_t {
	return (*C.nox_object_t)(controlRespawnItem(asObjectS(a), C.GoString(b), unsafe.Pointer(c), int32(d), int32(e)).CObj())
}

//export nox_xxx_playerMakeDefItems_4EF7D0
func nox_xxx_playerMakeDefItems_4EF7D0(a, b, c C.int) C.char {
	return C.char(controlDefaultItems(objectFromInt(a), int32(b), int32(c)))
}

//export nox_xxx_unitInitPlayer_4EFE80
func nox_xxx_unitInitPlayer_4EFE80(a *C.nox_object_t) C.char {
	return C.char(controlInitPlayer(asObjectS(a)))
}

//export sub_4EFF10
func sub_4EFF10(a C.int) C.int { return C.int(controlResetPlayer(objectFromInt(a))) }

//export nox_xxx_equipedItemByCode_4F7920
func nox_xxx_equipedItemByCode_4F7920(a, b C.int) C.int {
	return inventoryInt(controlEquippedByCode(objectFromInt(a), uint32(b)))
}

//export nox_xxx_playerSetCustomWP_4F79A0
func nox_xxx_playerSetCustomWP_4F79A0(a, b, c C.int) {
	controlSetWaypoint(objectFromInt(a), uint32(b), uint32(c))
}

//export nox_xxx_mapFindPlayerStart_4F7AB0
func nox_xxx_mapFindPlayerStart_4F7AB0(a *C.float2, b *C.nox_object_t) {
	controlFindStart((*types.Pointf)(unsafe.Pointer(a)), asObjectS(b))
}

//export nox_xxx_weaponGetStaminaByType_4F7E80
func nox_xxx_weaponGetStaminaByType_4F7E80(a C.int) C.int {
	return C.int(controlWeaponStamina(uint32(a)))
}

//export nox_xxx_playerRespawn_4F7EF0
func nox_xxx_playerRespawn_4F7EF0(a *C.nox_object_t) C.short {
	return C.short(controlRespawn(asObjectS(a)))
}

//export sub_4FA280
func sub_4FA280(a C.int) C.int { return C.int(controlWeaponAnimation(uint32(a))) }

//export nox_common_mapPlrActionToStateId_4FA2B0
func nox_common_mapPlrActionToStateId_4FA2B0(a *C.nox_object_t) C.int {
	return C.int(controlActionState(asObjectS(a)))
}

//export nox_xxx_checkInversionEffect_4FA4F0
func nox_xxx_checkInversionEffect_4FA4F0(a, b C.int) C.int {
	return C.int(controlInversion(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_mobMorphFromPlayer_4FAAC0
func nox_xxx_mobMorphFromPlayer_4FAAC0(a *C.uint32_t) C.char {
	return C.char(controlMorphFromPlayer((*server.Object)(unsafe.Pointer(a))))
}

//export nox_xxx_mobMorphToPlayer_4FAAF0
func nox_xxx_mobMorphToPlayer_4FAAF0(a *C.uint32_t) C.char {
	return C.char(controlMorphToPlayer((*server.Object)(unsafe.Pointer(a))))
}

//export nox_xxx_updatePlayerMonsterBot_4FAB20
func nox_xxx_updatePlayerMonsterBot_4FAB20(a *C.uint32_t) C.int {
	return C.int(controlBotUpdate((*server.Object)(unsafe.Pointer(a))))
}

//export nox_xxx_netSendRewardNotify_4FAD50
func nox_xxx_netSendRewardNotify_4FAD50(a, b, c C.int, d C.char) C.int {
	return C.int(controlRewardNotify(objectFromInt(a), int32(b), objectFromInt(c), byte(d)))
}

//export sub_4FADD0
func sub_4FADD0(a C.int, b *C.char, c C.char) { controlLockedDoor(objectFromInt(a), b, byte(c)) }

//export sub_4FB050
func sub_4FB050(a, b C.int, c *C.int) C.int {
	return C.int(controlGuideDamage(objectFromInt(a), objectFromInt(b), (*int32)(unsafe.Pointer(c))))
}
