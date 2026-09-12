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

func controlBool(v bool) C.int {
	if v {
		return 1
	}
	return 0
}

//export sub_4E5F40
func sub_4E5F40(a C.int) C.int { return C.int(controlRemoveGlyphs(objectFromInt(a))) }

//export sub_4E5FC0
func sub_4E5FC0(a C.int) { controlRemoveCreatures(objectFromInt(a)) }

//export sub_4E6150
func sub_4E6150(a *C.nox_playerInfo) *C.nox_object_t {
	return (*C.nox_object_t)(controlNextObserver(unsafe.Pointer(a)).CObj())
}

//export sub_4E6230
func sub_4E6230() C.int { return inventoryInt(controlFindBall()) }

//export nox_xxx_playerObserverFindGoodSlave0_4E6280
func nox_xxx_playerObserverFindGoodSlave0_4E6280(a *C.nox_playerInfo) *C.nox_object_t {
	return (*C.nox_object_t)(controlObserverSlave(unsafe.Pointer(a)).CObj())
}

//export nox_xxx_playerLeaveObserver_0_4E6AA0
func nox_xxx_playerLeaveObserver_0_4E6AA0(a *C.nox_playerInfo) {
	controlLeaveObserver(unsafe.Pointer(a))
}

//export nox_xxx_playerObserverFindGoodSlave2_4EC3E0
func nox_xxx_playerObserverFindGoodSlave2_4EC3E0(a C.int) C.int {
	return inventoryInt(controlSlave(objectFromInt(a), false))
}

//export nox_xxx_playerObserverFindGoodSlave_4EC420
func nox_xxx_playerObserverFindGoodSlave_4EC420(a C.int) C.int {
	return inventoryInt(controlSlave(objectFromInt(a), true))
}

//export nox_xxx_unitRemoveChild_4EC470
func nox_xxx_unitRemoveChild_4EC470(a *C.nox_object_t) { controlRemoveChildren(asObjectS(a)) }

//export nox_xxx_unitTransferSlaves_4EC4B0
func nox_xxx_unitTransferSlaves_4EC4B0(a *C.nox_object_t) { controlTransferChildren(asObjectS(a)) }

//export nox_xxx_abilGivePlayerAll_4EED40
func nox_xxx_abilGivePlayerAll_4EED40(a C.int, b C.char, c C.int) {
	controlGiveAbilities(objectFromInt(a), int8(b), int32(c))
}

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

//export nox_xxx_getRespawnWeaponFlags_4EF580
func nox_xxx_getRespawnWeaponFlags_4EF580() C.char { return C.char(controlRespawnFlags()) }

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

//export nox_xxx_netSendPlayerRespawn_4EFC30
func nox_xxx_netSendPlayerRespawn_4EFC30(a C.int, b C.char) C.int {
	return C.int(controlRespawnNotify(objectFromInt(a), byte(b)))
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

//export sub_4F7950
func sub_4F7950(a *C.nox_object_t) { controlClearWaypoints(asObjectS(a)) }

//export nox_xxx_playerSetCustomWP_4F79A0
func nox_xxx_playerSetCustomWP_4F79A0(a, b, c C.int) {
	controlSetWaypoint(objectFromInt(a), uint32(b), uint32(c))
}

//export nox_xxx_playerConfusedGetDirection_4F7A40
func nox_xxx_playerConfusedGetDirection_4F7A40(a *C.nox_object_t) C.int {
	return C.int(controlConfusedDirection(asObjectS(a)))
}

//export nox_xxx_mapFindPlayerStart_4F7AB0
func nox_xxx_mapFindPlayerStart_4F7AB0(a *C.float2, b *C.nox_object_t) {
	controlFindStart((*types.Pointf)(unsafe.Pointer(a)), asObjectS(b))
}

//export sub_4F7CE0
func sub_4F7CE0(a, b C.int) C.int {
	return controlBool(controlStartEligible(objectFromInt(a), int32(b)))
}

//export nox_xxx_playerSubStamina_4F7D30
func nox_xxx_playerSubStamina_4F7D30(a *C.nox_object_t, b C.int) C.int {
	return C.int(controlSubStamina(asObjectS(a), int32(b)))
}

//export sub_4F7DB0
func sub_4F7DB0(a C.int, b C.char) { controlAdjustStamina(objectFromInt(a), int8(b)) }

//export nox_xxx_checkWinkFlags_4F7DF0
func nox_xxx_checkWinkFlags_4F7DF0(a *C.nox_object_t) C.int {
	return C.int(controlDropBall(asObjectS(a)))
}

//export nox_xxx_weaponGetStaminaByType_4F7E80
func nox_xxx_weaponGetStaminaByType_4F7E80(a C.int) C.int {
	return C.int(controlWeaponStamina(uint32(a)))
}

//export nox_xxx_playerRespawn_4F7EF0
func nox_xxx_playerRespawn_4F7EF0(a *C.nox_object_t) C.short {
	return C.short(controlRespawn(asObjectS(a)))
}

//export sub_4F80C0
func sub_4F80C0(a C.int, b *C.float2) C.int {
	return C.int(controlNearStart(objectFromInt(a), (*types.Pointf)(unsafe.Pointer(b))))
}

//export sub_4F9A80
func sub_4F9A80(a *C.nox_object_t) C.int { return controlBool(controlHasWaypoint(asObjectS(a))) }

//export sub_4F9AB0
func sub_4F9AB0(a *C.nox_object_t) C.int { return C.int(controlWalkWaypoint(asObjectS(a))) }

//export nox_xxx_playerCanMove_4F9BC0
func nox_xxx_playerCanMove_4F9BC0(a *C.nox_object_t) C.int {
	return controlBool(controlCanMove(asObjectS(a)))
}

//export nox_xxx_playerCanAttack_4F9C40
func nox_xxx_playerCanAttack_4F9C40(a *C.nox_object_t) C.int {
	return controlBool(controlCanAttack(asObjectS(a)))
}

//export nox_xxx_playerInputAttack_4F9C70
func nox_xxx_playerInputAttack_4F9C70(a *C.nox_object_t) { controlInputAttack(asObjectS(a)) }

//export nox_xxx_playerAimsAtEnemy_4F9DC0
func nox_xxx_playerAimsAtEnemy_4F9DC0(a C.int) C.int {
	return controlBool(controlAimsAtEnemy(objectFromInt(a)))
}

//export sub_4F9E10
func sub_4F9E10(a *C.nox_object_t) C.int { return C.int(controlFollowEnemy(asObjectS(a))) }

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

//export nox_xxx_playerBotCreate_4FA700
func nox_xxx_playerBotCreate_4FA700(a *C.nox_object_t) C.int {
	return C.int(controlBotCreate(asObjectS(a)))
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

//export nox_xxx_monsterActionToPlrState_4FABC0
func nox_xxx_monsterActionToPlrState_4FABC0(a C.int) C.char {
	return C.char(controlBotState(objectFromInt(a)))
}

//export nox_xxx_respawnPlayerBot_4FAC70
func nox_xxx_respawnPlayerBot_4FAC70(a C.int) C.int {
	return C.int(controlRespawnBot(objectFromInt(a)))
}

//export nox_xxx_netSendRewardNotify_4FAD50
func nox_xxx_netSendRewardNotify_4FAD50(a, b, c C.int, d C.char) C.int {
	return C.int(controlRewardNotify(objectFromInt(a), int32(b), objectFromInt(c), byte(d)))
}

//export sub_4FADD0
func sub_4FADD0(a C.int, b *C.char, c C.char) { controlLockedDoor(objectFromInt(a), b, byte(c)) }

//export sub_4FB000
func sub_4FB000(a, b C.int) C.int {
	return C.int(controlGuideLevel(objectFromInt(a), objectFromInt(b)))
}

//export sub_4FB050
func sub_4FB050(a, b C.int, c *C.int) C.int {
	return C.int(controlGuideDamage(objectFromInt(a), objectFromInt(b), (*int32)(unsafe.Pointer(c))))
}

//export nox_xxx_playerDoSchedSpell_4FB0E0
func nox_xxx_playerDoSchedSpell_4FB0E0(a, b *C.nox_object_t) C.int {
	return C.int(controlScheduledSpell(asObjectS(a), asObjectS(b), false))
}

//export nox_xxx_playerDoSchedSpellQueue_4FB1D0
func nox_xxx_playerDoSchedSpellQueue_4FB1D0(a, b *C.nox_object_t) C.int {
	return C.int(controlScheduledSpell(asObjectS(a), asObjectS(b), true))
}
