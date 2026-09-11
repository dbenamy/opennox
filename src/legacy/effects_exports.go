package legacy

/*
#include "GAME3_2.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func effectsMod(p C.int) *server.ModifierEff {
	return (*server.ModifierEff)(unsafe.Pointer(uintptr(uint32(p))))
}

//export sub_4DFB50
func sub_4DFB50(a1 C.int, a2 C.int) { effectsEngageFlag(objectFromInt(a2), 8, 75) }

//export sub_4DFB80
func sub_4DFB80(a1 C.int, a2 C.int) { effectsDisengageFlag(objectFromInt(a2), 8, 76) }

//export nox_xxx_enchantItemTestInventory_4DFBB0
func nox_xxx_enchantItemTestInventory_4DFBB0(a1 C.int, a2 C.char) C.int {
	return C.int(effectsInventory(objectFromInt(a1), byte(a2)))
}

//export nox_xxx_effectSpeedEngage_4DFC30
func nox_xxx_effectSpeedEngage_4DFC30(a1 C.int, a2 C.int) {
	effectsSpeed(effectsMod(a1), objectFromInt(a2), true)
}

//export nox_xxx_effectSpeedDisengage_4DFCA0
func nox_xxx_effectSpeedDisengage_4DFCA0(a1 C.int, a2 C.int) {
	effectsSpeed(effectsMod(a1), objectFromInt(a2), false)
}

//export sub_4DFD10
func sub_4DFD10(a1 C.int, a2 C.int) { effectsEngageFlag(objectFromInt(a2), 1, 102) }

//export nox_xxx_modifFireProtection_4DFD40
func nox_xxx_modifFireProtection_4DFD40(a1 C.int, a2 C.int, a3 C.int) {
	if a2 != 0 && a3 != 0 {
		effectsDisengageFlag(objectFromInt(a2), 1, 103)
	}
}

//export nox_xxx_buff_4DFD80
func nox_xxx_buff_4DFD80(a1 C.int, a2 C.int) { effectsEngageFlag(objectFromInt(a2), 4, 106) }

//export sub_4DFDB0
func sub_4DFDB0(a1 C.int, a2 C.int) { effectsDisengageFlag(objectFromInt(a2), 4, 107) }

//export nox_xxx_checkPoisonProtectEnch_4DFDE0
func nox_xxx_checkPoisonProtectEnch_4DFDE0(a1 C.int, a2 C.int) {
	effectsEngageFlag(objectFromInt(a2), 2, 110)
}

//export sub_4DFE10
func sub_4DFE10(a1 C.int, a2 C.int) { effectsDisengageFlag(objectFromInt(a2), 2, 111) }

//export nox_xxx_checkFireProtect_4DFE40
func nox_xxx_checkFireProtect_4DFE40(a1 *C.uint32_t) C.double {
	return C.double(effectsProtection(equipmentObject(a1), C.sub_4DFD10, 17, "FireSpellProtection", .5, .60000002))
}

//export nox_xxx_checkElectrProtect_4DFF40
func nox_xxx_checkElectrProtect_4DFF40(a1 *C.uint32_t) C.double {
	return C.double(effectsProtection(equipmentObject(a1), C.nox_xxx_buff_4DFD80, 20, "ElectricitySpellProtection", .5, .60000002))
}

//export nox_xxx_getPoisonDmg_4E0040
func nox_xxx_getPoisonDmg_4E0040(a1 *C.uint32_t) C.double {
	return C.double(effectsProtection(equipmentObject(a1), C.nox_xxx_checkPoisonProtectEnch_4DFDE0, 18, "PoisonSpellProtection", .69999999, .89999998))
}

//export sub_4E0140
func sub_4E0140(a1 C.int, a2 C.int) { effectsEngageFlag(objectFromInt(a2), 32, 123) }

//export sub_4E0170
func sub_4E0170(a1 C.int, a2 C.int) {
	u := objectFromInt(a2)
	if u != nil && u.ObjClass&4 != 0 {
		effectsDisengageFlag(u, 32, 124)
	}
}

//export nox_xxx_effectRegeneration_4E01D0
func nox_xxx_effectRegeneration_4E01D0(a1 C.int, a2 C.int) {
	effectsRegeneration(effectsMod(a1), objectFromInt(a2))
}

//export nox_xxx_attribContinualReplen_4E02C0
func nox_xxx_attribContinualReplen_4E02C0(a1 C.int, a2 *C.uint32_t) {
	effectsReplenish(effectsMod(a1), equipmentObject(a2))
}

//export sub_4E0370
func sub_4E0370(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 C.int, a6 *C.float) *C.float {
	v := (*float32)(unsafe.Pointer(a6))
	*v = float32(float64(effectsMod(a1).Defend76.Valf) * float64(*v))
	return a6
}

//export sub_4E0380
func sub_4E0380(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 C.int, a6 *C.float) *C.float {
	v := (*float32)(unsafe.Pointer(a6))
	*v = float32((1 - float64(effectsMod(a1).Defend76.Valf) + 1) * float64(*v))
	return a6
}

//export nox_xxx_inversionEffect_4E03D0
func nox_xxx_inversionEffect_4E03D0(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 C.int, a6 *C.int) C.int {
	return C.int(effectsGrip(effectsMod(a1), (*int32)(unsafe.Pointer(a6)), true))
}

//export nox_xxx_unusedCheckGripEffect_4E03F0
func nox_xxx_unusedCheckGripEffect_4E03F0(a1 C.int, a2 C.int, a3 C.int, a4 C.int) C.int {
	return C.int(effectsGripSearch(objectFromInt(a1), objectFromInt(a2), objectFromInt(a3), objectFromInt(a4)))
}

//export nox_xxx_gripEffect_4E0480
func nox_xxx_gripEffect_4E0480(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 C.int, a6 *C.int) C.int {
	return C.int(effectsGrip(effectsMod(a1), (*int32)(unsafe.Pointer(a6)), false))
}

//export nox_xxx_effectDamageMultiplier_4E04C0
func nox_xxx_effectDamageMultiplier_4E04C0(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 *C.float) *C.float {
	v := (*float32)(unsafe.Pointer(a5))
	*v = float32(float64(effectsMod(a1).Attack40.Valf) * float64(*v))
	return a5
}

//export nox_xxx_stunEffect_4E04D0
func nox_xxx_stunEffect_4E04D0(a1 C.int, a2 C.int, a3 C.int, a4 C.int) {
	effectsStatus(effectsMod(a1), objectFromInt(a3), objectFromInt(a4), true)
}

//export nox_xxx_recoilEffect_4E0640
func nox_xxx_recoilEffect_4E0640(a1 C.int, a2 C.int, a3 C.int, a4 C.int) {
	effectsRecoil(effectsMod(a1), objectFromInt(a2), objectFromInt(a4))
}

//export nox_xxx_confuseEffect_4E0670
func nox_xxx_confuseEffect_4E0670(a1 C.int, a2 C.int, a3 C.int, a4 C.int) {
	effectsStatus(effectsMod(a1), objectFromInt(a3), objectFromInt(a4), false)
}

//export nox_xxx_lightngEffect_4E06F0
func nox_xxx_lightngEffect_4E06F0(a1 C.int, a2 C.int, a3 C.int, a4 C.int) {
	effectsLightning(effectsMod(a1), objectFromInt(a2), objectFromInt(a3), objectFromInt(a4))
}

//export nox_xxx_drainMEffect_4E0740
func nox_xxx_drainMEffect_4E0740(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 *C.int) {
	u, target := objectFromInt(a3), objectFromInt(a4)
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	effectsDrainMana(effectsMod(a1), u, target, int32(*a5))
}

//export nox_xxx_vampirismEffect_4E07C0
func nox_xxx_vampirismEffect_4E07C0(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 *C.int) {
	u, target := objectFromInt(a3), objectFromInt(a4)
	if u == nil || target == nil || target.ObjClass&6 == 0 || target.ObjClass&2 != 0 && target.ObjSubClass&0x40 != 0 {
		return
	}
	effectsVampirism(effectsMod(a1), u, target, int32(*a5))
}

//export nox_xxx_poisonEffect_4E0850
func nox_xxx_poisonEffect_4E0850(a1 C.int, a2 C.int, a3 C.int, a4 C.int) {
	effectsPoison(effectsMod(a1), objectFromInt(a3), objectFromInt(a4))
}

//export nox_xxx_sympathyEffect_4E08E0
func nox_xxx_sympathyEffect_4E08E0(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 *C.int) {
	u, target := objectFromInt(a3), objectFromInt(a4)
	if u == nil || target == nil || target.ObjClass&6 == 0 {
		return
	}
	effectsSympathy(effectsMod(a1), u, target, int32(*a5))
}

//export nox_xxx_itemCheckReadinessEffect_4E0960
func nox_xxx_itemCheckReadinessEffect_4E0960(a1 C.int) C.int {
	return C.int(effectsReadiness(objectFromInt(a1)))
}

//export nox_xxx_effectProjectileSpeed_4E09B0
func nox_xxx_effectProjectileSpeed_4E09B0(a1 C.int, a2 C.int, a3 C.int, a4 C.int, a5 C.int) C.int {
	u := objectFromInt(a5)
	u.SpeedCur = float32(float64(effectsMod(a1).Attack40.Valf) * float64(u.SpeedCur))
	return a5
}

//export nox_xxx_rechargeItem_53C520
func nox_xxx_rechargeItem_53C520(a1 C.int, a2 C.int) C.int {
	return C.int(effectsRecharge(objectFromInt(a1), int32(a2)))
}

//export nox_xxx_getRechargeRate_53C940
func nox_xxx_getRechargeRate_53C940(a1 *C.uint32_t) C.int {
	return C.int(effectsRechargeRate(equipmentObject(a1)))
}

//export nox_xxx_useLesserFireballStaff_53F290
func nox_xxx_useLesserFireballStaff_53F290(a1 C.int, a2 *C.uint32_t) C.int {
	return C.int(effectsLesserFireball(objectFromInt(a1), equipmentObject(a2)))
}

//export nox_xxx_wandShot_53F480
func nox_xxx_wandShot_53F480(a1 C.int, a2 C.int, a3 *C.int, a4 *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(effectsWandShot(objectFromInt(a1), int(a2), *(*types.Pointf)(unsafe.Pointer(a3)), uint32(uintptr(unsafe.Pointer(a4)))))))
}

//export nox_xxx_useWandCastSpell_53F4F0
func nox_xxx_useWandCastSpell_53F4F0(a1 C.int, a2 *C.uint32_t) C.int {
	return C.int(effectsWandCast(objectFromInt(a1), equipmentObject(a2)))
}

//export nox_xxx_useFireWand_53F670
func nox_xxx_useFireWand_53F670(a1 C.int, a2 C.int) C.int {
	return C.int(effectsFireWand(objectFromInt(a1), objectFromInt(a2)))
}

//export nox_xxx_useByNetCode_53F8E0
func nox_xxx_useByNetCode_53F8E0(a1 C.int, a2 C.int) C.int {
	return C.int(effectsUse(objectFromInt(a1), objectFromInt(a2)))
}
