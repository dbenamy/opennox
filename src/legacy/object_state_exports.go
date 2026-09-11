package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

//export nox_xxx_unitNeedSync_4E44F0
func nox_xxx_unitNeedSync_4E44F0(a *C.nox_object_t) { asObjectS(a).NeedSync() }

//export sub_4E4500
func sub_4E4500(a *C.nox_object_t, b, c, d C.int) *C.int {
	u := asObjectS(a)
	u.Sub_4E4500(uint32(b), uint32(c), d != 0)
	return (*C.int)(stateSyncEnd(u))
}

//export nox_xxx_unitSetOnOff_4E4670
func nox_xxx_unitSetOnOff_4E4670(a, b C.int) *C.int {
	return (*C.int)(stateOnOff(objectFromInt(a), b != 0))
}

//export nox_xxx_unitRaise_4E46F0
func nox_xxx_unitRaise_4E46F0(a *C.nox_object_t, b C.float) { stateRaise(asObjectS(a), float32(b)) }

//export nox_xxx_servMarkObjAnimFrame_4E4880
func nox_xxx_servMarkObjAnimFrame_4E4880(a, b C.int) *C.int {
	return (*C.int)(stateAnimation(objectFromInt(a), uint32(b)))
}

//export nox_xxx_setUnitBuffFlags_4E48F0
func nox_xxx_setUnitBuffFlags_4E48F0(a, b C.int) *C.int {
	return (*C.int)(stateBuffs(objectFromInt(a), uint32(b)))
}

//export nox_xxx_modifSetItemAttrs_4E4990
func nox_xxx_modifSetItemAttrs_4E4990(a *C.nox_object_t, b *C.int) *C.int {
	return (*C.int)(stateAttributes(asObjectS(a), unsafe.Pointer(b)))
}

//export nox_xxx_objectGetMass_4E4A70
func nox_xxx_objectGetMass_4E4A70(a C.int) C.double { return C.double(objectFromInt(a).Mass) }

//export nox_xxx_playerRemoveSpawnedStuff_4E5AD0
func nox_xxx_playerRemoveSpawnedStuff_4E5AD0(a *C.nox_object_t) { stateRemoveSpawned(asObjectS(a)) }

//export nox_xxx_isUnit_4E5B50
func nox_xxx_isUnit_4E5B50(a *C.nox_object_t) C.int {
	return C.int(bool2int(stateIsUnit(asObjectS(a))))
}

//export sub_4E5B80
func sub_4E5B80(a *C.nox_object_t) C.int { return C.int(bool2int(stateIsPixie(asObjectS(a)))) }

//export sub_4E5BF0
func sub_4E5BF0(a C.int) { stateCleanup(int32(a)) }

//export sub_4E6BD0
func sub_4E6BD0(a C.int) C.int {
	u := objectFromInt(a)
	return C.int(bool2int(u.HealthData != nil && GetServer().S().Frame()-u.Frame134 <= 1))
}

//export nox_xxx_calcDistance_4E6C00
func nox_xxx_calcDistance_4E6C00(a, b *C.nox_object_t) C.double {
	return C.double(stateDistance(asObjectS(a), asObjectS(b)))
}

//export sub_4E6CE0
func sub_4E6CE0(a, b *C.float2) C.int {
	return C.int(stateDirection((*types.Pointf)(unsafe.Pointer(a)), (*types.Pointf)(unsafe.Pointer(b))))
}

//export nox_server_testTwoPointsAndDirection_4E6E50
func nox_server_testTwoPointsAndDirection_4E6E50(a *C.float2, b C.int, c *C.float2) C.int {
	return C.int(stateFront((*types.Pointf)(unsafe.Pointer(a)), int32(b), (*types.Pointf)(unsafe.Pointer(c))))
}

//export nox_xxx_teleportToMB_4E7190
func nox_xxx_teleportToMB_4E7190(a *C.uint8_t, b *C.float) {
	stateTeleport(asObjectS((*C.nox_object_t)(unsafe.Pointer(a))), (*types.Pointf)(unsafe.Pointer(b)))
}

//export nox_xxx_objectUnkUpdateCoords_4E7290
func nox_xxx_objectUnkUpdateCoords_4E7290(a *C.nox_object_t) C.int {
	asObjectS(a).Nox_xxx_objectUnkUpdateCoords_4E7290()
	return C.int(uintptr(unsafe.Pointer(a)))
}

//export nox_xxx_spawnSomeBarrel_4E7470
func nox_xxx_spawnSomeBarrel_4E7470(a, b C.int) {
	stateLoot(objectFromInt(a), (*types.Pointf)(unsafe.Pointer(uintptr(b))))
}

//export sub_4E7540
func sub_4E7540(a, b *C.nox_object_t) { stateRememberAttacker(asObjectS(a), asObjectS(b)) }

//export nox_xxx_objectSetOn_4E75B0
func nox_xxx_objectSetOn_4E75B0(a *C.nox_object_t) C.char { return C.char(stateOn(asObjectS(a))) }

//export nox_xxx_objectSetOff_4E7600
func nox_xxx_objectSetOff_4E7600(a *C.nox_object_t) C.int { return C.int(stateOff(asObjectS(a))) }

//export sub_4E7700
func sub_4E7700(a C.int) C.int { return C.int(stateChecksum(objectFromInt(a))) }

//export nox_xxx_inventoryGetFirst_4E7980
func nox_xxx_inventoryGetFirst_4E7980(a C.int) C.int {
	return inventoryInt(objectFromInt(a).InvFirstItem)
}

//export nox_xxx_inventoryGetNext_4E7990
func nox_xxx_inventoryGetNext_4E7990(a C.int) C.int {
	if a == 0 {
		return 0
	}
	return inventoryInt(objectFromInt(a).InvNextItem)
}

//export sub_4E79B0
func sub_4E79B0(a C.int) C.int { *memmap.PtrUint32(0x5d4594, 1567712) = uint32(a); return a }

//export nox_xxx_unitFreeze_4E79C0
func nox_xxx_unitFreeze_4E79C0(a *C.nox_object_t, b C.int) C.char {
	return C.char(stateFreeze(asObjectS(a), int32(b)))
}

//export nox_xxx_unitUnFreeze_4E7A60
func nox_xxx_unitUnFreeze_4E7A60(a *C.nox_object_t, b C.int) C.char {
	return C.char(stateUnfreeze(asObjectS(a), int32(b)))
}

//export nox_xxx_unitBecomePet_4E7B00
func nox_xxx_unitBecomePet_4E7B00(a, b C.int) { statePet(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_monsterRemoveMonitors_4E7B60
func nox_xxx_monsterRemoveMonitors_4E7B60(a, b *C.nox_object_t) {
	stateRemoveMonitors(asObjectS(a), asObjectS(b))
}

//export sub_4E7BC0
func sub_4E7BC0(a C.int) C.int {
	if a == 0 {
		return 0
	}
	return C.int(uint32(objectFromInt(a).ObjClass) >> 2 & 1)
}

//export nox_xxx_unitIsCrown_4E7BE0
func nox_xxx_unitIsCrown_4E7BE0(a C.int) C.int {
	return C.int(bool2int(stateOwns(objectFromInt(a), 1567716, "Crown")))
}

//export nox_xxx_unitIsGameball_4E7C30
func nox_xxx_unitIsGameball_4E7C30(a C.int) C.int {
	return C.int(bool2int(stateOwns(objectFromInt(a), 1567720, "GameBall")))
}

//export nox_xxx_unitCountSlaves_4E7CF0
func nox_xxx_unitCountSlaves_4E7CF0(a, b, c C.int) C.int {
	return C.int(stateCount(objectFromInt(a), uint32(b), uint32(c)))
}

//export sub_4E7DE0
func sub_4E7DE0(a C.int, b *C.nox_object_t) C.int {
	return C.int(bool2int(stateEqual(objectFromInt(a), asObjectS(b))))
}

//export nox_xxx_unitPostCreateNotify_4E7F10
func nox_xxx_unitPostCreateNotify_4E7F10(a *C.nox_object_t) *C.char {
	statePostCreate(asObjectS(a))
	return nil
}

//export sub_4E8110
func sub_4E8110(a C.int) *C.char { statePlayerVisibility(int32(a)); return nil }

//export sub_4E81D0
func sub_4E81D0(a *C.nox_object_t) C.int { return C.int(stateResetPixie(asObjectS(a))) }

//export nox_xxx_fnFindCloseDoors_4E8340
func nox_xxx_fnFindCloseDoors_4E8340(a *C.float, b C.int) {
	stateCloseDoor(asObjectS((*C.nox_object_t)(unsafe.Pointer(a))), unsafe.Pointer(uintptr(b)))
}

//export sub_4E8390
func sub_4E8390(a C.int) C.int { return C.int(stateDoorNotify(objectFromInt(a))) }

//export nox_xxx_collideMonsterEventProc_4E83B0
func nox_xxx_collideMonsterEventProc_4E83B0(a, b C.int) *C.uchar {
	return (*C.uchar)(stateMonsterCollision(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collideMimic_4E83D0
func nox_xxx_collideMimic_4E83D0(a, b C.int) *C.uchar {
	return (*C.uchar)(stateMimicCollision(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collidePlayer_4E8460
func nox_xxx_collidePlayer_4E8460(a, b C.int) {
	statePlayerCollision(objectFromInt(a), objectFromInt(b))
}

//export nox_objectCollideDefault
func nox_objectCollideDefault(a, b C.int, c *C.float) C.int { return 0 }
