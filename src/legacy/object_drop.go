package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
int nox_objectDropAudEvent_4EE2F0(nox_object_t* a1, nox_object_t* a2, float2* a3);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_objectDropAudEvent_4EE2F0 server.DropFunc
)

func init() {
	registerInventoryDrop("DefaultDrop", C.nox_xxx_dropDefault_4ED290, inventoryDefaultDrop)
	registerInventoryDrop("ArmorDrop", C.nox_xxx_dropArmor_53EB70, func(u, it *server.Object, p *types.Pointf) int { return inventoryEquipmentDrop(u, it, p, true) })
	registerInventoryDrop("WeaponDrop", C.nox_xxx_dropWeapon_53AB10, func(u, it *server.Object, p *types.Pointf) int { return inventoryEquipmentDrop(u, it, p, false) })
	registerInventoryDrop("TreasureDrop", C.nox_xxx_dropTreasure_4ED710, inventoryTreasureDrop)
	registerInventoryDrop("GlyphDrop", C.nox_GlyphDrop_4ED500, inventoryGlyphDrop)
	registerInventoryDrop("PotionDrop", C.sub_4EDDE0, inventoryPotionDrop)
	registerInventoryDrop("TrapDrop", C.nox_xxx_dropTrap_4ED580, inventoryTrapDrop)
	registerInventoryDrop("FoodDrop", C.nox_xxx_dropFood_4EDE50, inventoryFoodDrop)
	registerInventoryDrop("CrownDrop", C.nox_xxx_dropCrown_4ED5E0, inventoryCrownDrop)
	server.RegisterObjectDrop("AudEventDrop", C.nox_objectDropAudEvent_4EE2F0, func(obj, obj2 *server.Object, pos types.Pointf) bool {
		return Nox_objectDropAudEvent_4EE2F0(obj, obj2, pos)
	})
	registerInventoryDrop("AnkhTradableDrop", C.nox_xxx_dropAnkhTradable_4EE370, inventoryDefaultDrop)
}

//export nox_objectDropAudEvent_4EE2F0
func nox_objectDropAudEvent_4EE2F0(cobj1 *nox_object_t, cobj2 *nox_object_t, a3 *C.float2) int {
	return bool2int(Nox_objectDropAudEvent_4EE2F0(asObjectS(cobj1), asObjectS(cobj2), *(*types.Pointf)(unsafe.Pointer(a3))))
}

func Nox_xxx_dropDefault_4ED290(obj1 *server.Object, obj2 *server.Object, a3 *types.Pointf) int {
	return inventoryDefaultDrop(obj1, obj2, a3)
}

var inventoryNativeDrops = make(map[unsafe.Pointer]func(*server.Object, *server.Object, *types.Pointf) int)

func registerInventoryDrop(name string, ptr unsafe.Pointer, fn func(*server.Object, *server.Object, *types.Pointf) int) {
	inventoryNativeDrops[ptr] = fn
	server.RegisterObjectDrop(name, ptr, func(u, it *server.Object, pos types.Pointf) bool { return fn(u, it, &pos) != 0 })
}
