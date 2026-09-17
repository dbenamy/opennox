package legacy

/*
#include <stdint.h>
#include <stdbool.h>
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_collideDoor_4E8AC0
func nox_xxx_collideDoor_4E8AC0(a, b C.int) { worldCollideDoor(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collidePickup_4E8DF0
func nox_xxx_collidePickup_4E8DF0(a, b C.int) C.int {
	return C.int(worldCollidePickup(objectFromInt(a), objectFromInt(b)))
}

//export sub_4E8E50
func sub_4E8E50() *C.uchar { return (*C.uchar)(unsafe.Pointer(worldQuestPending())) }

//export sub_4E8E60
func sub_4E8E60() C.int { return C.int(worldQuestCountdown()) }

//export nox_server_questMaybeWarp_4E8F60
func nox_server_questMaybeWarp_4E8F60() C.bool { return C.bool(worldQuestMaybeWarp()) }

//export sub_4E9010
func sub_4E9010() C.int { return C.int(bool2int(worldQuestExitReady())) }

//export nox_xxx_collideExit_4E9090
func nox_xxx_collideExit_4E9090(a, b, n C.int) { worldCollideExit(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_spellFlyCollide_4E9500
func nox_xxx_spellFlyCollide_4E9500(a, b C.int, n *C.float) {
	worldCollideSpellProjectile(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collideChest_4E9C40
func nox_xxx_collideChest_4E9C40(a *C.uint32_t, b C.int) {
	worldCollideChest((*server.Object)(unsafe.Pointer(a)), objectFromInt(b))
}

//export sub_4EAAA0
func sub_4EAAA0(a C.int) { worldCollideBarrel(objectFromInt(a)) }

//export sub_4EAAD0
func sub_4EAAD0(a, b C.int) { worldCollideAudio(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collidePentagram_4EAB20
func nox_xxx_collidePentagram_4EAB20(a C.int) C.int {
	return C.int(worldCollidePentagram(objectFromInt(a)))
}

//export nox_xxx_collideSign_4EAB40
func nox_xxx_collideSign_4EAB40(a, b C.int) { worldCollideSign(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideTrapDoor_4EAB60
func nox_xxx_collideTrapDoor_4EAB60(a, b C.int) { worldCollideTrap(objectFromInt(a), objectFromInt(b)) }

//export sub_4EACA0
func sub_4EACA0(a, b C.int) { worldCollideTeleport(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideSpellPedestal_4EAD20
func nox_xxx_collideSpellPedestal_4EAD20(a, b C.int) C.int {
	return C.int(worldCollideSpellAward(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collideUndeadKiller_4EBD40
func nox_xxx_collideUndeadKiller_4EBD40(a, b, n C.int) {
	worldCollideUndead(objectFromInt(a), objectFromInt(b), n != 0)
}

//export nox_xxx_collideMonsterGen_4EBE10
func nox_xxx_collideMonsterGen_4EBE10(a, b C.int) {
	worldCollideGenerator(objectFromInt(a), objectFromInt(b))
}

//export sub_4EBE40
func sub_4EBE40(a, b C.int) { worldCollideSoulGate(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideAnkhQuest_4EBF40
func nox_xxx_collideAnkhQuest_4EBF40(a, b C.int) {
	worldCollideAnkh(objectFromInt(a), objectFromInt(b))
}
