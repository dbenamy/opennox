package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

//export nox_xxx_parseDamageTypeByName_4E0A00
func nox_xxx_parseDamageTypeByName_4E0A00(a *C.char) C.int {
	return C.int(damageTypeByName(alloc.GoString((*byte)(unsafe.Pointer(a)))))
}

//export nox_xxx_projectileReflect_4E0A70
func nox_xxx_projectileReflect_4E0A70(a, b C.int) C.int {
	return C.int(damageReflect(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_damageDefaultProc_4E0B30
func nox_xxx_damageDefaultProc_4E0B30(a, b, c, d, e C.int) C.int {
	return C.int(damageDefault(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_gameballOnPlayerDamage_4E1230
func nox_xxx_gameballOnPlayerDamage_4E1230(a, b, c C.int) {
	damageBall(objectFromInt(a), objectFromInt(b), int32(c))
}

//export nox_xxx_itemApplyDefendEffect2_4E1320
func nox_xxx_itemApplyDefendEffect2_4E1320(a, b, c C.int, d *C.int, e C.int) C.int {
	return C.int(damageDefend(objectFromInt(a), objectFromInt(b), objectFromInt(c), (*int32)(unsafe.Pointer(d)), int32(e)))
}

//export nox_xxx_itemApplyPreDamageEffect_4E13B0
func nox_xxx_itemApplyPreDamageEffect_4E13B0(a, b, c, d C.int) C.int {
	return C.int(damagePre(objectFromInt(a), objectFromInt(b), objectFromInt(c), (*int32)(unsafe.Pointer(uintptr(d)))))
}

//export sub_4E1400
func sub_4E1400(a C.int, b *C.uint32_t) C.int {
	return C.int(bool2int(damageMelee(objectFromInt(a), equipmentObject(b))))
}

//export sub_4E1470
func sub_4E1470(a C.int) C.int { return C.int(bool2int(damageFriendlyWeapon(objectFromInt(a)))) }

//export sub_4E14A0
func sub_4E14A0() C.int { return 0 }

//export sub_4E14B0
func sub_4E14B0(a, b, c, d, e C.int) C.int {
	return C.int(damageWeapon(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_damageArmor_4E1500
func nox_xxx_damageArmor_4E1500(a, b, c, d, e C.int) C.int {
	return C.int(damageArmor(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_playerDamageWeapon_4E1560
func nox_xxx_playerDamageWeapon_4E1560(a, b, c, d C.int, e C.float, f C.int) {
	damageDurability(objectFromInt(a), objectFromInt(b), objectFromInt(c), objectFromInt(d), float32(e), int32(f), true)
}

//export nox_xxx_itemDestroyed_4E1650
func nox_xxx_itemDestroyed_4E1650(a C.int, b *C.uint32_t, c, d C.ushort) C.int {
	return C.int(damageItemReport(int32(a), equipmentObject(b), uint16(c), uint16(d)))
}

//export nox_xxx_equipDamage_4E16D0
func nox_xxx_equipDamage_4E16D0(a, b, c, d C.int, e C.float, f C.int) {
	damageDurability(objectFromInt(a), objectFromInt(b), objectFromInt(c), objectFromInt(d), float32(e), int32(f), false)
}

//export nox_server_handler_PlayerDamage_4E17B0
func nox_server_handler_PlayerDamage_4E17B0(a, b, c, d, e C.int) C.int {
	return C.int(damagePlayer(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_playerDecrementHPMana_4E20F0
func nox_xxx_playerDecrementHPMana_4E20F0(a, b C.int, c C.float) {
	damageFraction(objectFromInt(a), (*int32)(unsafe.Pointer(uintptr(b))), float32(c))
}

//export nox_xxx_playerDamageItems_4E2180
func nox_xxx_playerDamageItems_4E2180(a, b, c, d C.int, e C.float) {
	damageInventory(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), float32(e))
}

//export sub_4E2220
func sub_4E2220(a C.int) C.double { return C.double(damageConductivity(objectFromInt(a))) }

//export sub_4E22A0
func sub_4E22A0(a, b, c, d C.int, e C.float, f C.int) C.int {
	return C.int(damageBlockingItem(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), float32(e), int32(f), true))
}

//export sub_4E2330
func sub_4E2330(a, b, c, d C.int, e C.float, f C.int) C.int {
	return C.int(damageBlockingItem(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), float32(e), int32(f), false))
}

//export sub_4E23C0
func sub_4E23C0(a, b, c, d, e C.int) C.int {
	return C.int(damageSkeleton(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export sub_4E24B0
func sub_4E24B0(a, b, c, d, e C.int) C.int {
	return C.int(damageDefault(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export sub_4E24E0
func sub_4E24E0(a, b, c, d, e C.int) C.int {
	if e == 9 || e == 17 {
		d *= 2
	}
	return C.int(damageDefault(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_damageFlammable_4E2520
func nox_xxx_damageFlammable_4E2520(a, b, c, d, e C.int) C.int {
	if e == 1 || e == 12 || e == 7 {
		d = 9999999
	}
	return C.int(damageDefault(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_damageBlackPowder_4E2560
func nox_xxx_damageBlackPowder_4E2560(a, b, c, d, e C.int) C.int {
	if e != 0 && e != 1 && e != 2 && e != 12 {
		return 0
	}
	d = 999999
	return C.int(damageDefault(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}

//export nox_xxx_damageMonsterGen_4E27D0
func nox_xxx_damageMonsterGen_4E27D0(a, b, c, d, e C.int) C.int {
	return C.int(damageGenerator(objectFromInt(a), objectFromInt(b), objectFromInt(c), int32(d), int32(e)))
}
