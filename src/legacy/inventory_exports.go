package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"unsafe"
)

func inventoryPos(p unsafe.Pointer) *types.Pointf { return (*types.Pointf)(p) }

//export sub_4ED0C0
func sub_4ED0C0(u, it *nox_object_t) { inventoryRemove(asObjectS(u), asObjectS(it)) }

//export nox_xxx_inventoryPutImpl_4F3070
func nox_xxx_inventoryPutImpl_4F3070(u, it *nox_object_t, arg C.int) {
	inventoryInsert(asObjectS(u), asObjectS(it), int(arg))
}

//export nox_xxx_dropDefault_4ED290
func nox_xxx_dropDefault_4ED290(u, it *nox_object_t, p *C.float2) C.int {
	return C.int(inventoryDefaultDrop(asObjectS(u), asObjectS(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_GlyphDrop_4ED500
func nox_GlyphDrop_4ED500(u, it C.int, p *C.float2) C.int {
	return C.int(inventoryGlyphDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_dropTrap_4ED580
func nox_xxx_dropTrap_4ED580(u, it C.int, p *C.float2) C.int {
	return C.int(inventoryTrapDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_dropCrown_4ED5E0
func nox_xxx_dropCrown_4ED5E0(u, it C.int, p *C.int) C.int {
	return C.int(inventoryCrownDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_dropTreasure_4ED710
func nox_xxx_dropTreasure_4ED710(u, it C.int, p *C.int) C.int {
	return C.int(inventoryTreasureDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_drop_4ED790
func nox_xxx_drop_4ED790(u, it *nox_object_t, p *C.float2) C.int {
	return C.int(inventoryDrop(asObjectS(u), asObjectS(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_drop_4ED810
func nox_xxx_drop_4ED810(u, it C.int, p *C.float) C.int {
	return C.int(inventoryTargetDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_invForceDropItem_4ED930
func nox_xxx_invForceDropItem_4ED930(u C.int, it *C.uint32_t) C.int {
	return C.int(inventoryForceDrop(objectFromInt(u), asObjectS((*nox_object_t)(unsafe.Pointer(it)))))
}

//export sub_4ED970
func sub_4ED970(radius C.float, origin, p *C.float2) *C.float2 {
	inventoryRandomPlacement(float32(radius), inventoryPos(unsafe.Pointer(origin)), inventoryPos(unsafe.Pointer(p)))
	return p
}

//export nox_xxx_dropAllItems_4EDA40
func nox_xxx_dropAllItems_4EDA40(u *C.uint32_t) *C.uint32_t {
	return resourceReturnPointer(inventoryDropAll(asObjectS((*nox_object_t)(unsafe.Pointer(u)))))
}

//export sub_4EDDE0
func sub_4EDDE0(u C.int, it *C.uint32_t, p *C.int) C.int {
	return C.int(inventoryPotionDrop(objectFromInt(u), asObjectS((*nox_object_t)(unsafe.Pointer(it))), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_dropFood_4EDE50
func nox_xxx_dropFood_4EDE50(u, it C.int, p *C.int) C.int {
	return C.int(inventoryFoodDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}

//export nox_xxx_chest_4EDF00
func nox_xxx_chest_4EDF00(u, it C.int) { inventoryChest(objectFromInt(u), objectFromInt(it)) }

//export nox_xxx_pickupFood_4F3350
func nox_xxx_pickupFood_4F3350(u, it, arg C.int) C.int {
	return C.int(inventoryFoodPickup(objectFromInt(u), objectFromInt(it), int(arg)))
}

//export sub_4F3400
func sub_4F3400(u, it, arg C.int) C.int {
	return C.int(inventoryCrownPickup(objectFromInt(u), objectFromInt(it), int(arg)))
}

//export nox_xxx_pickupUse_4F34D0
func nox_xxx_pickupUse_4F34D0(u, it, arg C.int) C.int {
	return C.int(inventoryUsePickup(objectFromInt(u), objectFromInt(it), int(arg)))
}

//export nox_xxx_pickupTrap_4F3510
func nox_xxx_pickupTrap_4F3510(u, it, arg C.int) C.int {
	return C.int(inventoryTrapPickup(objectFromInt(u), objectFromInt(it), int(arg)))
}

//export nox_xxx_pickupTreasure_4F3580
func nox_xxx_pickupTreasure_4F3580(u, it, arg C.int) C.int {
	return C.int(inventoryTreasurePickup(objectFromInt(u), objectFromInt(it), int(arg)))
}

//export nox_xxx_pickupAmmo_4F3B00
func nox_xxx_pickupAmmo_4F3B00(u C.int, it *nox_object_t, arg, equip C.int) C.int {
	return C.int(inventoryAmmoPickup(objectFromInt(u), asObjectS(it), int(arg), int(equip)))
}

//export nox_xxx_pickupSpellbook_4F3C60
func nox_xxx_pickupSpellbook_4F3C60(u, it, arg C.int) C.int {
	return C.int(inventoryBookPickup(objectFromInt(u), objectFromInt(it), int(arg), false))
}

//export nox_xxx_pickupAbilitybook_4F3CE0
func nox_xxx_pickupAbilitybook_4F3CE0(u, it, arg C.int) C.int {
	return C.int(inventoryBookPickup(objectFromInt(u), objectFromInt(it), int(arg), true))
}

//export sub_4F3DD0
func sub_4F3DD0(u, it C.int) C.int {
	return C.int(inventoryAnkhPickup(objectFromInt(u), objectFromInt(it)))
}

//export sub_53A720
func sub_53A720(u C.int, it *nox_object_t, arg, equip C.int) C.int {
	return C.int(inventoryWeaponPickup(objectFromInt(u), asObjectS(it), int(arg), int(equip)))
}

//export nox_xxx_sendMsgOblivionPickup_53A9C0
func nox_xxx_sendMsgOblivionPickup_53A9C0(u C.int, it *nox_object_t, arg, equip C.int) C.int {
	return C.int(inventoryOblivionPickup(objectFromInt(u), asObjectS(it), int(arg), int(equip)))
}

//export nox_xxx_dropWeapon_53AB10
func nox_xxx_dropWeapon_53AB10(u C.int, it *C.uint32_t, p *C.int) C.int {
	return C.int(inventoryEquipmentDrop(objectFromInt(u), asObjectS((*nox_object_t)(unsafe.Pointer(it))), inventoryPos(unsafe.Pointer(p)), false))
}

//export nox_xxx_pickupArmor_53E7F0
func nox_xxx_pickupArmor_53E7F0(u, it, arg, equip C.int) C.int {
	return C.int(inventoryArmorPickup(objectFromInt(u), objectFromInt(it), int(arg), int(equip)))
}

//export nox_xxx_dropArmor_53EB70
func nox_xxx_dropArmor_53EB70(u C.int, it *C.uint32_t, p *C.int) C.int {
	return C.int(inventoryEquipmentDrop(objectFromInt(u), asObjectS((*nox_object_t)(unsafe.Pointer(it))), inventoryPos(unsafe.Pointer(p)), true))
}

//export nox_xxx_ItemIsDroppable_53EBF0
func nox_xxx_ItemIsDroppable_53EBF0(it C.int) C.int {
	return C.int(bool2int(inventoryDroppable(objectFromInt(it))))
}

//export nox_xxx_dropAnkhTradable_4EE370
func nox_xxx_dropAnkhTradable_4EE370(u, it C.int, p *C.int) C.int {
	return C.int(inventoryDefaultDrop(objectFromInt(u), objectFromInt(it), inventoryPos(unsafe.Pointer(p))))
}
