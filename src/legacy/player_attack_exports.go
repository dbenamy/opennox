package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_playerPreAttackEffects_538290
func nox_xxx_playerPreAttackEffects_538290(a, b, c, d C.int) C.int {
	return C.int(attackPreEffects(objectFromInt(a), objectFromInt(b), objectFromInt(c), (*attackRecord)(unsafe.Pointer(uintptr(d)))))
}

//export nox_xxx_playerTraceAttack_538330
func nox_xxx_playerTraceAttack_538330(a, b C.int) C.int {
	return C.int(attackTrace(objectFromInt(a), (*attackRecord)(unsafe.Pointer(uintptr(b)))))
}

//export sub_538510
func sub_538510(a, b C.int) { attackHit(objectFromInt(a), (*attackRecord)(unsafe.Pointer(uintptr(b)))) }

//export sub_5386A0
func sub_5386A0(a, b C.int) { attackNearest(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_itemApplyAttackEffect_538840
func nox_xxx_itemApplyAttackEffect_538840(a, b, c C.int) C.int {
	return C.int(attackItemEffects(objectFromInt(a), objectFromInt(b), (*attackRecord)(unsafe.Pointer(uintptr(c)))))
}

//export nox_xxx_playerAttack_538960
func nox_xxx_playerAttack_538960(a *nox_object_t) C.int { return C.int(attackPlayer(asObjectS(a))) }

//export nox_xxx_warcryStunMonsters_539B90
func nox_xxx_warcryStunMonsters_539B90(a, b C.int) C.short {
	return C.short(attackWarcry(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_shootBowCrossbow1_539BD0
func nox_xxx_shootBowCrossbow1_539BD0(a, b C.int) C.int {
	return C.int(attackBow(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_shootBowCrossbow2_539D80
func nox_xxx_shootBowCrossbow2_539D80(a, b, c C.int, d *C.char) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(attackShoot(objectFromInt(a), objectFromInt(b), objectFromInt(c), uint32(uintptr(unsafe.Pointer(d)))))))
}

//export nox_xxx_shootApplyEffects_539F40
func nox_xxx_shootApplyEffects_539F40(a, b, c C.int) C.int {
	return C.int(attackShotEffects(objectFromInt(a), objectFromInt(b), objectFromInt(c)))
}

//export sub_539FB0
func sub_539FB0(a *C.uint32_t) C.int { return C.int(attackReload(equipmentObject(a), 128)) }

//export nox_xxx_playerTryReloadQuiver_539FF0
func nox_xxx_playerTryReloadQuiver_539FF0(a *C.uint32_t) C.int {
	return C.int(attackReload(equipmentObject(a), 2))
}
