package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_XFerSpellReward_4F5F30
func nox_xxx_XFerSpellReward_4F5F30(u *C.int) C.int {
	return C.int(itemXferSpellReward(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerAbilityReward_4F6240
func nox_xxx_XFerAbilityReward_4F6240(u *C.int) C.int {
	return C.int(itemXferAbilityReward(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerFieldGuide_4F6390
func nox_xxx_XFerFieldGuide_4F6390(u *C.int) C.int {
	return C.int(itemXferFieldGuide(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerWeapon_4F64A0
func nox_xxx_XFerWeapon_4F64A0(u C.int) C.int { return C.int(itemXferWeapon(objectFromInt(u))) }

//export nox_xxx_XFerArmor_4F6860
func nox_xxx_XFerArmor_4F6860(u C.int) C.int { return C.int(itemXferArmor(objectFromInt(u))) }

//export nox_xxx_XFerAmmo_4F6B20
func nox_xxx_XFerAmmo_4F6B20(u *C.int) C.int {
	return C.int(itemXferAmmo(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerTeam_4F6D20
func nox_xxx_XFerTeam_4F6D20(u *C.int) C.int {
	return C.int(itemXferTeam(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerGold_4F6EC0
func nox_xxx_XFerGold_4F6EC0(u C.int) C.int { return C.int(itemXferGold(objectFromInt(u))) }

//export nox_xxx_XFerObelisk_4F6F60
func nox_xxx_XFerObelisk_4F6F60(u *C.int) C.int {
	return C.int(itemXferObelisk(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerToxicCloud_4F70A0
func nox_xxx_XFerToxicCloud_4F70A0(u C.int) C.int { return C.int(itemXferToxicCloud(objectFromInt(u))) }

//export nox_xxx_XFerMonsterGen_4F7130
func nox_xxx_XFerMonsterGen_4F7130(u *C.int) C.int {
	return C.int(itemXferGenerator(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export nox_xxx_XFerRewardMarker_4F74D0
func nox_xxx_XFerRewardMarker_4F74D0(u *C.int) C.int {
	return C.int(itemXferRewardMarker(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}
