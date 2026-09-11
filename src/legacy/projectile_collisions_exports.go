package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_collideProjectileGeneric_4E87B0
func nox_xxx_collideProjectileGeneric_4E87B0(a, b C.int) {
	projectileGeneric(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideProjectileSpark_4E8880
func nox_xxx_collideProjectileSpark_4E8880(a, b C.int) {
	projectileGenericSpark(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideDamage_4E9430
func nox_xxx_collideDamage_4E9430(a, b C.int) {
	projectileDamageField(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideManadrain_4E9490
func nox_xxx_collideManadrain_4E9490(a, b C.int) {
	projectileManaDrain(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideBomb_4E96F0
func nox_xxx_collideBomb_4E96F0(a, b C.int) { projectileBomb(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideBoom_4E9770
func nox_xxx_collideBoom_4E9770(a, b C.int, n *C.float) {
	projectileBoom(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collideDie_4E99B0
func nox_xxx_collideDie_4E99B0(a, b C.int) { projectileDie(objectFromInt(a), objectFromInt(b)) }

//export sub_4E9A30
func sub_4E9A30(a, b *nox_object_t) C.int {
	return C.int(bool2int(projectileTrapEligible(asObjectS(a), asObjectS(b))))
}

//export nox_xxx_fireballCollide_4E9AC0
func nox_xxx_fireballCollide_4E9AC0(a, b C.int) {
	projectileFireball(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideSulphurShot2_4E9D80
func nox_xxx_collideSulphurShot2_4E9D80(a, b C.int, n *C.float) {
	projectileSulphur(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collideSulphurShot_4E9E50
func nox_xxx_collideSulphurShot_4E9E50(a, b, n C.int) {
	projectileSulphurTrail(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(uintptr(n))))
}

//export nox_xxx_collideDeathBallFragment_4E9FE0
func nox_xxx_collideDeathBallFragment_4E9FE0(a, b C.int, n *C.float) {
	projectileDeathFragment(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collidePixie_4EA080
func nox_xxx_collidePixie_4EA080(a, b C.int, n *C.float) {
	projectilePixie(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collideWallReflectSpark_4EA200
func nox_xxx_collideWallReflectSpark_4EA200(a, b C.int, n *C.float2) {
	projectileWallSpark(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export sub_4EA2C0
func sub_4EA2C0(a, b C.int) { projectileSparkOwner(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideSpark_4EA300
func nox_xxx_collideSpark_4EA300(a, b C.int, n *C.float) {
	projectileSpark(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export nox_xxx_collideWebbing_4EA380
func nox_xxx_collideWebbing_4EA380(a, b C.int) { projectileWeb(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideFist_4EADF0
func nox_xxx_collideFist_4EADF0(a, b C.int) { projectileFist(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideTeleportWake_4EAE30
func nox_xxx_collideTeleportWake_4EAE30(a, b C.int) {
	projectileTeleportWake(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideChakram_4EAF00
func nox_xxx_collideChakram_4EAF00(a, b C.int, n *C.float) {
	projectileChakram(objectFromInt(a), objectFromInt(b), (*types.Pointf)(unsafe.Pointer(n)))
}

//export sub_4EB250
func sub_4EB250(a C.int) C.int { return inventoryInt(projectileChakramSelect(objectFromInt(a))) }

//export sub_4EB340
func sub_4EB340(a *C.float, b C.int) {
	projectileChakramCandidate((*server.Object)(unsafe.Pointer(a)), (*types.Pointf)(unsafe.Pointer(uintptr(b))))
}

//export sub_4EB3E0
func sub_4EB3E0(a C.int) { projectileChakramFallback(objectFromInt(a)) }

//export nox_xxx_collideArrow_4EB490
func nox_xxx_collideArrow_4EB490(a, b C.int) { projectileArrow(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_collideMonsterArrow_4EB800
func nox_xxx_collideMonsterArrow_4EB800(a, b C.int) {
	projectileMonsterArrow(objectFromInt(a), objectFromInt(b))
}

//export nox_xxx_collideBearTrap_4EB890
func nox_xxx_collideBearTrap_4EB890(a *C.int, b C.int) {
	projectileTrap((*server.Object)(unsafe.Pointer(a)), objectFromInt(b), false)
}

//export nox_xxx_collidePoisonGasTrap_4EB910
func nox_xxx_collidePoisonGasTrap_4EB910(a *C.int, b C.int) {
	projectileTrap((*server.Object)(unsafe.Pointer(a)), objectFromInt(b), true)
}
