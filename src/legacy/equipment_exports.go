package legacy

/*
#include "defs.h"
#include <stdbool.h>
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func equipmentObject(u *C.uint32_t) *server.Object { return (*server.Object)(unsafe.Pointer(u)) }

//export nox_xxx_equipWeaponNPC_53A030
func nox_xxx_equipWeaponNPC_53A030(u, it C.int) C.int {
	return C.int(equipmentNPCDequipWeapon(objectFromInt(u), objectFromInt(it)))
}

//export sub_53A0F0
func sub_53A0F0(u, a, b C.int) { equipmentDequipAmmo(objectFromInt(u), int(a), int(b)) }

//export nox_xxx_playerDequipWeapon_53A140
func nox_xxx_playerDequipWeapon_53A140(u *C.uint32_t, it *nox_object_t, a, b C.int) C.int {
	return C.int(equipmentDequipWeapon(equipmentObject(u), asObjectS(it), int(a), int(b)))
}

//export nox_xxx_NPCEquipWeapon_53A2C0
func nox_xxx_NPCEquipWeapon_53A2C0(u C.int, it *nox_object_t) C.int {
	return C.int(equipmentNPCEquipWeapon(objectFromInt(u), asObjectS(it)))
}

//export sub_53A3D0
func sub_53A3D0(u *C.uint32_t) { equipmentRemoveShields(equipmentObject(u)) }

//export nox_xxx_playerEquipWeapon_53A420
func nox_xxx_playerEquipWeapon_53A420(u *C.uint32_t, it *nox_object_t, a, b C.int) C.int {
	return C.int(equipmentEquipWeapon(equipmentObject(u), asObjectS(it), int(a), int(b)))
}

//export sub_53A680
func sub_53A680(u C.int) C.int { return C.int(equipmentEquipBow(objectFromInt(u))) }

//export sub_53A6C0
func sub_53A6C0(u C.int, it *nox_object_t) { equipmentPickupSound(objectFromInt(u), asObjectS(it)) }

//export sub_53AAB0
func sub_53AAB0(it C.int) { equipmentDropSound(objectFromInt(it)) }

//export sub_53AB90
func sub_53AB90(u, it C.int) { equipmentSecondary(objectFromInt(u), objectFromInt(it)) }

//export sub_53E2D0
func sub_53E2D0(it C.int) C.int { return C.int(equipmentArmorMask(objectFromInt(it))) }

//export nox_xxx_recalculateArmorVal_53E300
func nox_xxx_recalculateArmorVal_53E300(u *C.uint32_t) C.int {
	return C.int(equipmentRecalculate(equipmentObject(u)))
}

//export sub_53E3A0
func sub_53E3A0(u C.int, it *nox_object_t) C.int {
	return C.int(equipmentNPCDequipArmor(objectFromInt(u), asObjectS(it)))
}

//export sub_53E430
func sub_53E430(u *C.uint32_t, it *nox_object_t, a, b C.int) C.int {
	return C.int(equipmentDequipArmor(equipmentObject(u), asObjectS(it), int(a), int(b)))
}

//export nox_xxx_NPCEquipArmor_53E520
func nox_xxx_NPCEquipArmor_53E520(u C.int, it *C.uint32_t) C.int {
	return C.int(equipmentNPCEquipArmor(objectFromInt(u), equipmentObject(it)))
}

//export sub_53E600
func sub_53E600(u *C.uint32_t) { equipmentRemoveWeapons(equipmentObject(u)) }

//export nox_xxx_playerEquipArmor_53E650
func nox_xxx_playerEquipArmor_53E650(u *C.uint32_t, it *nox_object_t, a, b C.int) C.int {
	return C.int(equipmentEquipArmor(equipmentObject(u), asObjectS(it), int(a), int(b)))
}

//export nox_xxx_armorHaveSameSubclass_53E7B0
func nox_xxx_armorHaveSameSubclass_53E7B0(u, it C.int) *C.uint32_t {
	return (*C.uint32_t)(equipmentSameArmor(objectFromInt(u), objectFromInt(it)).CObj())
}

//export sub_53EAE0
func sub_53EAE0(it C.int) { equipmentArmorDropSound(objectFromInt(it)) }

//export sub_53EC40
func sub_53EC40() *C.char { equipmentInitDropTable(); return nil }

//export sub_53EC80
func sub_53EC80(it, mask C.int) C.int {
	return C.int(equipmentDropPolicy(objectFromInt(it), int(mask)))
}

//export nox_xxx_npcSetItemEquipFlags_4E4B20
func nox_xxx_npcSetItemEquipFlags_4E4B20(u C.int, it *nox_object_t, v C.int) *C.int {
	return (*C.int)(equipmentNPCSync(objectFromInt(u), asObjectS(it), int(v)))
}

//export nox_xxx_inventoryCountObjects_4E7D30
func nox_xxx_inventoryCountObjects_4E7D30(u, typ C.int) C.int {
	return C.int(equipmentCount(objectFromInt(u), int(typ)))
}

//export sub_4E7EC0
func sub_4E7EC0(u C.int, it *nox_object_t) C.int {
	return C.int(equipmentDuplicate(objectFromInt(u), asObjectS(it)))
}

//export nox_xxx_playerTryEquip_4F2F70
func nox_xxx_playerTryEquip_4F2F70(u, it *nox_object_t) C.int {
	return C.int(equipmentTryEquip(asObjectS(u), asObjectS(it)))
}

//export nox_xxx_playerTryDequip_4F2FB0
func nox_xxx_playerTryDequip_4F2FB0(u, it *nox_object_t) C.int {
	return C.int(equipmentTryDequip(asObjectS(u), asObjectS(it)))
}

//export nox_xxx_itemApplyEngageEffect_4F2FF0
func nox_xxx_itemApplyEngageEffect_4F2FF0(it *nox_object_t, u C.int) C.int {
	return C.int(equipmentEffects(asObjectS(it), objectFromInt(u), true))
}

//export nox_xxx_itemApplyDisengageEffect_4F3030
func nox_xxx_itemApplyDisengageEffect_4F3030(it *nox_object_t, u C.int) C.int {
	return C.int(equipmentEffects(asObjectS(it), objectFromInt(u), false))
}

//export nox_xxx_playerCheckStrength_4F3180
func nox_xxx_playerCheckStrength_4F3180(u, it *nox_object_t) C.bool {
	return C.bool(equipmentCheckStrength(asObjectS(u), asObjectS(it)))
}

//export sub_980523
func sub_980523(u *nox_object_t) { equipmentSaveShield(asObjectS(u)) }

//export sub_9805EB
func sub_9805EB(u *nox_object_t) *nox_object_t { return asObjectC(equipmentFindShield(asObjectS(u))) }

//export nox_xxx_itemApplyDefendEffect_415C00
func nox_xxx_itemApplyDefendEffect_415C00(it C.int) C.double {
	return C.double(equipmentDefend(objectFromInt(it)))
}

//export nox_xxx_unitGetStrength_4F9FD0
func nox_xxx_unitGetStrength_4F9FD0(u C.int) C.int { return C.int(equipmentStrength(objectFromInt(u))) }
