package legacy

/*
#include "common/alloc/classes/alloc_class.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
#include "GAME5.h"

*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

var (
	Sub_4E3B80 func(ind int) bool
)

func init() {
	server.DefaultDamage = C.nox_xxx_damageDefaultProc_4E0B30
	server.DefaultDamageSound = C.nox_xxx_soundDefaultDamageSound_532E20
	server.DefaultXfer = C.nox_xxx_XFerDefault_4F49A0

	server.RegisterObjectCreateGo("MonsterCreate", C.nox_xxx_monsterCreateFn_54C480, func(u *server.Object) { nox_xxx_monsterCreateFn_54C480(asObjectC(u)) })
	server.RegisterObjectCreateGo("ArmorCreate", C.sub_54C950, func(u *server.Object) { sub_54C950(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("WeaponCreate", C.nox_xxx_createWeapon_54C710, func(u *server.Object) { nox_xxx_createWeapon_54C710(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("ObeliskCreate", C.nox_xxx_createFnObelisk_54CA10, func(u *server.Object) { nox_xxx_createFnObelisk_54CA10(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("AnimCreate", C.nox_xxx_createFnAnim_54CA50, func(u *server.Object) { nox_xxx_createFnAnim_54CA50(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("TriggerCreate", C.nox_xxx_createTrigger_54CA60, func(u *server.Object) { nox_xxx_createTrigger_54CA60(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("MonsterGeneratorCreate", C.nox_xxx_createMonsterGen_54CA90, func(u *server.Object) { nox_xxx_createMonsterGen_54CA90(C.int(uintptr(u.CObj()))) })
	server.RegisterObjectCreateGo("RewardMarkerCreate", C.nox_xxx_createRewardMarker_54CAC0, func(u *server.Object) { nox_xxx_createRewardMarker_54CAC0(C.int(uintptr(u.CObj()))) })

	server.RegisterObjectInitGo("MonsterInit", C.nox_xxx_unitMonsterInit_4F0040, func(u *server.Object) { nox_xxx_unitMonsterInit_4F0040(asObjectC(u)) }, 0)
	server.RegisterObjectInitGo("PlayerInit", C.nox_xxx_unitInitPlayer_4EFE80, func(u *server.Object) { nox_xxx_unitInitPlayer_4EFE80(asObjectC(u)) }, 0)
	server.RegisterObjectInitGo("SparkInit", C.nox_xxx_unitSparkInit_4F0390, func(u *server.Object) { nox_xxx_unitSparkInit_4F0390(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectInitGo("FrogInit", C.nox_xxx_initFrog_4F03B0, func(u *server.Object) { nox_xxx_initFrog_4F03B0(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectInitGo("ChestInit", C.nox_xxx_initChest_4F0400, func(u *server.Object) { nox_xxx_initChest_4F0400(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectInitGo("BoulderInit", C.nox_xxx_unitBoulderInit_4F0420, func(u *server.Object) { nox_xxx_unitBoulderInit_4F0420((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectInitGo("BreakInit", C.nox_xxx_breakInit_4F0570, func(u *server.Object) { nox_xxx_breakInit_4F0570(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectInitGo("MonsterGeneratorInit", C.nox_xxx_unitInitGenerator_4F0590, func(u *server.Object) { nox_xxx_unitInitGenerator_4F0590(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectInitGo("ShopkeeperInit", C.nox_xxx_unitMonsterInit_4F0040, func(u *server.Object) { nox_xxx_unitMonsterInit_4F0040(asObjectC(u)) }, unsafe.Sizeof(server.ShopkeeperInitData{}))
	server.RegisterObjectInitGo("SkullInit", C.sub_4F0450, func(u *server.Object) { sub_4F0450(C.int(uintptr(u.CObj()))) }, 8)
	server.RegisterObjectInitGo("DirectionInit", C.sub_4F0490, func(u *server.Object) { sub_4F0490(C.int(uintptr(u.CObj()))) }, 8)
	server.RegisterObjectInitGo("GoldInit", C.nox_xxx_unitInitGold_4F04B0, func(u *server.Object) { nox_xxx_unitInitGold_4F04B0(C.int(uintptr(u.CObj()))) }, unsafe.Sizeof(server.GoldInitData{}))
}

func nox_xxx_newObjectWithTypeInd_4E3450(ind int) *nox_object_t {
	s := GetServer().S()
	return asObjectC(s.NewObjectByTypeInd(ind))
}

func nox_xxx_getUnitName_4E39D0(cobj *nox_object_t) *C.char {
	return (*C.char)(internCStr(GetServer().S().Types.ByInd(int(asObjectS(cobj).TypeInd)).ID()))
}

func sub_4E3B80(ind int) int { return bool2int(Sub_4E3B80(ind)) }

func nox_xxx_getUnitNameByThingType_4E3A80(ind int) *C.char {
	if ind == 0 {
		return nil
	}
	return (*C.char)(internCStr(GetServer().S().Types.ByInd(ind).ID()))
}

func nox_xxx_newObjectByTypeID_4E3810(cstr *C.char) *nox_object_t {
	obj := GetServer().S().NewObjectByTypeID(GoString(cstr))
	if obj == nil {
		return nil
	}
	return asObjectC(obj)
}

func Get_nox_xxx_XFerInvLight_4F5AA0() unsafe.Pointer {
	return unsafe.Pointer(C.nox_xxx_XFerInvLight_4F5AA0)
}
func Get_nox_xxx_unitInitGold_4F04B0() unsafe.Pointer {
	return unsafe.Pointer(C.nox_xxx_unitInitGold_4F04B0)
}
func Nox_call_objectType_new_go(a1 unsafe.Pointer, a2 *server.Object) {
	ccall.CallVoidPtr(a1, a2.CObj())
}
