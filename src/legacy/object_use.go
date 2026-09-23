package legacy

/*
#include "GAME4_3.h"

*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_useConsume_53EE10      server.UseFunc
	Nox_xxx_useMushroom_53ECE0     server.UseFunc
	Nox_xxx_useCiderConfuse_53EF00 server.UseFunc
	Nox_xxx_useEnchant_53ED60      server.UseFunc
	Nox_xxx_useCast_53ED90         server.UseFunc
	Nox_xxx_usePotion_53EF70       server.UseFunc
)

func init() {
	server.RegisterObjectUse("ConsumeUse", C.nox_xxx_useConsume_53EE10, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_useConsume_53EE10(obj, obj2)
	}, unsafe.Sizeof(server.ConsumeUseData{}))
	server.RegisterObjectUse("ConsumeConfuseUse", C.nox_xxx_useCiderConfuse_53EF00, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_useCiderConfuse_53EF00(obj, obj2)
	}, unsafe.Sizeof(server.ConsumeUseData{}))
	server.RegisterObjectUse("CastUse", C.nox_xxx_useCast_53ED90, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_useCast_53ED90(obj, obj2)
	}, unsafe.Sizeof(server.CastUseData{}))
	server.RegisterObjectUse("EnchantUse", C.nox_xxx_useEnchant_53ED60, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_useEnchant_53ED60(obj, obj2)
	}, unsafe.Sizeof(server.EnchantUseData{}))
	server.RegisterObjectUse("MushroomUse", C.nox_xxx_useMushroom_53ECE0, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_useMushroom_53ECE0(obj, obj2)
	}, 0)
	server.RegisterObjectUse("PotionUse", C.nox_xxx_usePotion_53EF70, func(obj, obj2 *server.Object) bool {
		return Nox_xxx_usePotion_53EF70(obj, obj2)
	}, unsafe.Sizeof(server.PotionUseData{}))

	server.RegisterObjectUse("FireWandUse", C.nox_xxx_useFireWand_53F670, func(u, it *server.Object) bool { return effectsFireWand(u, it) != 0 }, 0)
	server.RegisterObjectUse("ReadUse", C.nox_xxx_useRead_53F7C0, func(u, it *server.Object) bool { return unitRead(u, it, false) }, 260)
	server.RegisterObjectUse("WarpReadUse", C.sub_53F830, func(u, it *server.Object) bool { return unitRead(u, it, true) }, 260)
	server.RegisterObjectUse("WandUse", C.nox_xxx_useLesserFireballStaff_53F290, func(u, it *server.Object) bool { return effectsLesserFireball(u, it) != 0 }, 116)
	server.RegisterObjectUse("WandCastUse", C.nox_xxx_useWandCastSpell_53F4F0, func(u, it *server.Object) bool { return effectsWandCast(u, it) != 0 }, 116)
	server.RegisterObjectUse("SpellRewardUse", C.nox_xxx_useSpellReward_53F9E0, func(u, it *server.Object) bool { return bookUseSpell(u, it) != 0 }, unsafe.Sizeof(server.SpellRewardUseData{}))
	server.RegisterObjectUse("AbilityRewardUse", C.nox_xxx_useAbilityReward_53FAE0, func(u, it *server.Object) bool { return bookUseAbility(u, it) != 0 }, unsafe.Sizeof(server.AbilityRewardUseData{}))
	server.RegisterObjectUse("FieldGuideUse", C.sub_53F930, func(u, it *server.Object) bool { return bookUseGuide(u, it) != 0 }, unsafe.Sizeof(server.FieldGuideUseData{}))

	server.RegisterObjectUseParse("WandUse", resourceObjectParser("use", "wand"))
	server.RegisterObjectUseParse("WandCastUse", resourceObjectParser("use", "wandcast"))
}

//export nox_xxx_useMushroom_53ECE0
func nox_xxx_useMushroom_53ECE0(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_useMushroom_53ECE0(asObjectS(cobj1), asObjectS(cobj2)))
}

//export nox_xxx_useCiderConfuse_53EF00
func nox_xxx_useCiderConfuse_53EF00(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_useCiderConfuse_53EF00(asObjectS(cobj1), asObjectS(cobj2)))
}

//export nox_xxx_useEnchant_53ED60
func nox_xxx_useEnchant_53ED60(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_useEnchant_53ED60(asObjectS(cobj1), asObjectS(cobj2)))
}

//export nox_xxx_useCast_53ED90
func nox_xxx_useCast_53ED90(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_useCast_53ED90(asObjectS(cobj1), asObjectS(cobj2)))
}

//export nox_xxx_useConsume_53EE10
func nox_xxx_useConsume_53EE10(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_useConsume_53EE10(asObjectS(cobj1), asObjectS(cobj2)))
}

//export nox_xxx_usePotion_53EF70
func nox_xxx_usePotion_53EF70(cobj1 *nox_object_t, cobj2 *nox_object_t) int {
	return bool2int(Nox_xxx_usePotion_53EF70(asObjectS(cobj1), asObjectS(cobj2)))
}

func Get_nox_xxx_usePotion_53EF70() unsafe.Pointer {
	return C.nox_xxx_usePotion_53EF70
}

func Get_nox_xxx_useSpellReward_53F9E0() unsafe.Pointer {
	return C.nox_xxx_useSpellReward_53F9E0
}

func Get_nox_xxx_useAbilityReward_53FAE0() unsafe.Pointer {
	return C.nox_xxx_useAbilityReward_53FAE0
}

func Get_nox_xxx_useEnchant_53ED60() unsafe.Pointer {
	return C.nox_xxx_useEnchant_53ED60
}

func Get_nox_xxx_useCast_53ED90() unsafe.Pointer {
	return C.nox_xxx_useCast_53ED90
}

func Get_sub_53F930() unsafe.Pointer {
	return C.sub_53F930
}
