package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_xxx_unitSparkInit_4F0390
func nox_xxx_unitSparkInit_4F0390(a C.int) *C.uint32_t {
	return (*C.uint32_t)(rewardInitSpark(objectFromInt(a)))
}

//export nox_xxx_initFrog_4F03B0
func nox_xxx_initFrog_4F03B0(a C.int) C.int { return C.int(rewardInitFrog(objectFromInt(a))) }

//export nox_xxx_initChest_4F0400
func nox_xxx_initChest_4F0400(a C.int) *C.int {
	u := objectFromInt(a)
	rewardInitBreakable(u)
	return (*C.int)(u.CObj())
}

//export nox_xxx_unitBoulderInit_4F0420
func nox_xxx_unitBoulderInit_4F0420(a *C.uint32_t) *C.uint32_t {
	u := equipmentObject(a)
	u.Pos39 = u.PosVec
	return a
}

//export sub_4F0450
func sub_4F0450(a C.int) C.int { return C.int(rewardInitDirection(objectFromInt(a), true)) }

//export sub_4F0490
func sub_4F0490(a C.int) C.int { return C.int(rewardInitDirection(objectFromInt(a), false)) }

//export nox_xxx_unitInitGold_4F04B0
func nox_xxx_unitInitGold_4F04B0(a C.int) C.int { return C.int(rewardInitGold(objectFromInt(a))) }

//export nox_xxx_breakInit_4F0570
func nox_xxx_breakInit_4F0570(a C.int) *C.int {
	u := objectFromInt(a)
	rewardInitBreakable(u)
	return (*C.int)(u.CObj())
}

//export nox_xxx_unitInitGenerator_4F0590
func nox_xxx_unitInitGenerator_4F0590(a C.int) C.int {
	return C.int(rewardInitGenerator(objectFromInt(a)))
}

//export nox_server_rewardgen_activateMarker_4F0720
func nox_server_rewardgen_activateMarker_4F0720(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardMarker(objectFromInt(a), uint32(b)).CObj())
}

//export nox_xxx_rewardSpellBook_4F09F0
func nox_xxx_rewardSpellBook_4F09F0(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardBook(objectFromInt(a), uint32(b), 0).CObj())
}

//export nox_server_rewardGen_pickRandomSlots_4F0B60
func nox_server_rewardGen_pickRandomSlots_4F0B60(a C.uint) C.int { return C.int(rewardTier(uint32(a))) }

//export nox_xxx_rewardAbilityBook_4F0C70
func nox_xxx_rewardAbilityBook_4F0C70(a C.int) *C.uint32_t {
	return (*C.uint32_t)(rewardBook(objectFromInt(a), 0, 1).CObj())
}

//export nox_xxx_rewardFieldGuide_4F0D20
func nox_xxx_rewardFieldGuide_4F0D20(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardBook(objectFromInt(a), uint32(b), 2).CObj())
}

//export nox_xxx_rewardMakeArmor_4F0E80
func nox_xxx_rewardMakeArmor_4F0E80(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardEquipment(uint32(b), true).CObj())
}

//export nox_xxx_rewardMakeWeapon_4F14E0
func nox_xxx_rewardMakeWeapon_4F14E0(a C.int, b C.uint) C.int {
	return inventoryInt(rewardEquipment(uint32(b), false))
}

//export nox_xxx_rewardMakePotion_4F1C40
func nox_xxx_rewardMakePotion_4F1C40(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardPotion(uint32(b)).CObj())
}

//export nox_xxx_createGem_4F1D30
func nox_xxx_createGem_4F1D30(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardGem(uint32(b)).CObj())
}

//export nox_xxx_createGem2_4F1F00
func nox_xxx_createGem2_4F1F00(a C.int, b C.uint) *C.uint32_t {
	return (*C.uint32_t)(rewardGem(uint32(b)).CObj())
}

//export sub_4F2110
func sub_4F2110() { rewardPlaceAnkh() }

//export sub_4F2210
func sub_4F2210() C.int { rewardSelectMarkers(); return 0 }
