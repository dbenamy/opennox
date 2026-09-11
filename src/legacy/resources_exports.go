package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func resourceReturnPointer(v uint32) *C.uint32_t { return (*C.uint32_t)(unsafe.Pointer(uintptr(v))) }

//export nox_xxx_unitSetHP_4E4560
func nox_xxx_unitSetHP_4E4560(u *nox_object_t, v C.ushort) C.int {
	return C.int(resourceSetHP(asObjectS(u), uint16(v)))
}

//export nox_xxx_unitAdjustHP_4EE460
func nox_xxx_unitAdjustHP_4EE460(u *nox_object_t, v C.int) { resourceAdjustHP(asObjectS(u), int32(v)) }

//export nox_xxx_unitDamageClear_4EE5E0
func nox_xxx_unitDamageClear_4EE5E0(u *nox_object_t, v C.int) { resourceDamage(asObjectS(u), int32(v)) }

//export nox_xxx_unitHPsetOnMax_4EE6F0
func nox_xxx_unitHPsetOnMax_4EE6F0(u C.int) { resourceRestoreHP(objectFromInt(u)) }

//export nox_xxx_playerHP_4EE730
func nox_xxx_playerHP_4EE730(u C.int) { resourceHPHistory(objectFromInt(u)) }

//export nox_xxx_unitGetHP_4EE780
func nox_xxx_unitGetHP_4EE780(u *nox_object_t) C.short { return C.short(resourceGetHP(asObjectS(u))) }

//export nox_xxx_unitGetMaxHP_4EE7A0
func nox_xxx_unitGetMaxHP_4EE7A0(u C.int) C.short { return C.short(resourceGetMaxHP(objectFromInt(u))) }

//export nox_xxx_unitSetMaxHP_4EE7C0
func nox_xxx_unitSetMaxHP_4EE7C0(u C.int, v C.short) C.int {
	return C.int(resourceSetMaxHP(objectFromInt(u), uint16(v)))
}

//export nox_xxx_activatePoison_4EE7E0
func nox_xxx_activatePoison_4EE7E0(u, v, max C.int) C.int {
	return C.int(bool2int(resourcePoison(objectFromInt(u), int32(v), int32(max))))
}

//export nox_xxx_updatePoison_4EE8F0
func nox_xxx_updatePoison_4EE8F0(u *nox_object_t, v C.int) {
	resourceReducePoison(asObjectS(u), int32(v))
}

//export nox_xxx_removePoison_4EE9D0
func nox_xxx_removePoison_4EE9D0(u *nox_object_t) { resourceRemovePoison(asObjectS(u)) }

//export nox_xxx_setSomePoisonData_4EEA90
func nox_xxx_setSomePoisonData_4EEA90(u, v C.int) { resourceSetPoison(objectFromInt(u), int32(v)) }

//export nox_xxx_playerManaAdd_4EEB80
func nox_xxx_playerManaAdd_4EEB80(u *nox_object_t, v C.short) C.ushort {
	return C.ushort(resourceAddMana(asObjectS(u), int16(v)))
}

//export nox_xxx_playerManaSub_4EEBF0
func nox_xxx_playerManaSub_4EEBF0(u, v C.int) *C.uint32_t {
	return resourceReturnPointer(resourceSubMana(objectFromInt(u), int32(v)))
}

//export nox_xxx_unitGetOldMana_4EEC80
func nox_xxx_unitGetOldMana_4EEC80(u C.int) C.short {
	return C.short(resourceGetMana(objectFromInt(u)))
}

//export nox_xxx_playerGetMaxMana_4EECB0
func nox_xxx_playerGetMaxMana_4EECB0(u C.int) C.short {
	return C.short(resourceGetMaxMana(objectFromInt(u)))
}

//export nox_xxx_playerSetMaxMana_4EECD0
func nox_xxx_playerSetMaxMana_4EECD0(u C.int, v C.short) C.int {
	return C.int(resourceSetMaxMana(objectFromInt(u), uint16(v)))
}

//export nox_xxx_playerManaRefresh_4EECF0
func nox_xxx_playerManaRefresh_4EECF0(u C.int) *C.uint32_t {
	return resourceReturnPointer(resourceRefreshMana(objectFromInt(u)))
}

//export nox_xxx_playerAddGold_4FA590
func nox_xxx_playerAddGold_4FA590(u, v C.int) *C.uint32_t {
	return resourceReturnPointer(resourceAddGold(objectFromInt(u), uint32(v)))
}

//export nox_xxx_playerSubGold_4FA5D0
func nox_xxx_playerSubGold_4FA5D0(u C.int, v C.uint) *C.uint32_t {
	return resourceReturnPointer(resourceSubGold(objectFromInt(u), uint32(v)))
}

//export nox_xxx_playerGetGold_4FA6B0
func nox_xxx_playerGetGold_4FA6B0(u C.int) C.int { return C.int(resourceGetGold(objectFromInt(u))) }

//export nox_object_getGold_4FA6D0
func nox_object_getGold_4FA6D0(u *nox_object_t) C.int { return C.int(resourceObjectGold(asObjectS(u))) }

//export nox_xxx_pickupGold_4F3A60_obj_pickup
func nox_xxx_pickupGold_4F3A60_obj_pickup(u, item, flags C.int) C.int {
	return C.int(bool2int(resourceGoldPickup(objectFromInt(u), objectFromInt(item), int(flags))))
}
