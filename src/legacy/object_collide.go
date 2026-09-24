package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
#include "GAME5.h"

int nox_objectCollideDefault(int a1, int a2, float* a3);
void nox_xxx_collideDeathBall_4E9E90(nox_object_t* a1, nox_object_t* a2, float* a3);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_collideDeathBall_4E9E90 func(a1, a2 *server.Object, pos *types.Pointf)
	Nox_xxx_castCounterSpell_52BBB0 func(cspl spell.ID, a2, a3, a4 *server.Object, sa *server.SpellAcceptArg, lvl int) int
	Nox_xxx_changeOwner_52BE40      func(a1, a2 *server.Object)
)

func init() {
	server.RegisterObjectCollideGo("DefaultCollide", C.nox_objectCollideDefault, func(u *server.Object, a2, a3 uintptr) {
		nox_objectCollideDefault(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("MonsterCollide", C.nox_xxx_collideMonsterEventProc_4E83B0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideMonsterEventProc_4E83B0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("PlayerCollide", C.nox_xxx_collidePlayer_4E8460, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collidePlayer_4E8460(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("ProjectileCollide", C.nox_xxx_collideProjectileGeneric_4E87B0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideProjectileGeneric_4E87B0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("ProjectileSparkCollide", C.nox_xxx_collideProjectileSpark_4E8880, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideProjectileSpark_4E8880(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("DoorCollide", C.nox_xxx_collideDoor_4E8AC0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideDoor_4E8AC0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("PickupCollide", C.nox_xxx_collidePickup_4E8DF0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collidePickup_4E8DF0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("ExitCollide", C.nox_xxx_collideExit_4E9090, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideExit_4E9090(C.int(uintptr(u.CObj())), C.int(a2), C.int(a3))
	}, 88)
	server.RegisterObjectCollideGo("DamageCollide", C.nox_xxx_collideDamage_4E9430, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideDamage_4E9430(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("ManaDrainCollide", C.nox_xxx_collideManadrain_4E9490, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideManadrain_4E9490(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("BombCollide", C.nox_xxx_collideBomb_4E96F0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideBomb_4E96F0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("SparkExplosionCollide", C.nox_xxx_fireballCollide_4E9AC0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_fireballCollide_4E9AC0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 1)
	server.RegisterObjectCollideGo("ChestCollide", C.nox_xxx_collideChest_4E9C40, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideChest_4E9C40((*C.uint32_t)(u.CObj()), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("WallReflectCollide", C.nox_xxx_collideSulphurShot2_4E9D80, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideSulphurShot2_4E9D80(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 8)
	server.RegisterObjectCollideGo("WallReflectSparkCollide", C.nox_xxx_collideWallReflectSpark_4EA200, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideWallReflectSpark_4EA200(C.int(uintptr(u.CObj())), C.int(a2), (*C.float2)(unsafe.Pointer(a3)))
	}, 8)
	server.RegisterObjectCollideGo("PixieCollide", C.nox_xxx_collidePixie_4EA080, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collidePixie_4EA080(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 8)
	server.RegisterObjectCollideGo("OwnCollide", C.sub_4EA2C0, func(u *server.Object, a2, a3 uintptr) { sub_4EA2C0(C.int(uintptr(u.CObj())), C.int(a2)) }, 0)
	server.RegisterObjectCollideGo("SparkCollide", C.nox_xxx_collideSpark_4EA300, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideSpark_4EA300(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 8)
	server.RegisterObjectCollideGo("BarrelCollide", C.sub_4EAAA0, func(u *server.Object, a2, a3 uintptr) { sub_4EAAA0(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectCollideGo("AudioEventCollide", C.sub_4EAAD0, func(u *server.Object, a2, a3 uintptr) { sub_4EAAD0(C.int(uintptr(u.CObj())), C.int(a2)) }, 4)
	server.RegisterObjectCollideGo("TriggerCollide", C.nox_xxx_collideTrigger_54FCD0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideTrigger_54FCD0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("TeleportCollide", C.sub_4EACA0, func(u *server.Object, a2, a3 uintptr) { sub_4EACA0(C.int(uintptr(u.CObj())), C.int(a2)) }, 8)
	server.RegisterObjectCollideGo("ElevatorCollide", C.nox_objectCollideDefault, func(u *server.Object, a2, a3 uintptr) {
		nox_objectCollideDefault(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 8)
	server.RegisterObjectCollideGo("AwardSpellCollide", C.nox_xxx_collideSpellPedestal_4EAD20, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideSpellPedestal_4EAD20(C.int(uintptr(u.CObj())), C.int(a2))
	}, 4)
	server.RegisterObjectCollideGo("DieCollide", C.nox_xxx_collideDie_4E99B0, func(u *server.Object, a2, a3 uintptr) { nox_xxx_collideDie_4E99B0(C.int(uintptr(u.CObj())), C.int(a2)) }, 0)
	server.RegisterObjectCollideGo("GlyphCollide", C.nox_xxx_collideGlyph_4E9A00, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideGlyph_4E9A00((*nox_object_t)(u.CObj()), (*nox_object_t)(unsafe.Pointer(a2)))
	}, 0)
	server.RegisterObjectCollideGo("SpellProjectileCollide", C.nox_xxx_spellFlyCollide_4E9500, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_spellFlyCollide_4E9500(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("BoomCollide", C.nox_xxx_collideBoom_4E9770, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideBoom_4E9770(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("SignCollide", C.nox_xxx_collideSign_4EAB40, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideSign_4EAB40(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("PentagramCollide", C.nox_xxx_collidePentagram_4EAB20, func(u *server.Object, a2, a3 uintptr) { nox_xxx_collidePentagram_4EAB20(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectCollideGo("SpiderSpitCollide", C.nox_xxx_collideWebbing_4EA380, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideWebbing_4EA380(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("DeathBallCollide", C.nox_xxx_collideDeathBall_4E9E90, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideDeathBall_4E9E90((*nox_object_t)(u.CObj()), (*nox_object_t)(unsafe.Pointer(a2)), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("DeathBallFragmentCollide", C.nox_xxx_collideDeathBallFragment_4E9FE0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideDeathBallFragment_4E9FE0(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("TelekinesisCollide", C.nox_objectCollideDefault, func(u *server.Object, a2, a3 uintptr) {
		nox_objectCollideDefault(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("FistCollide", C.nox_xxx_collideFist_4EADF0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideFist_4EADF0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("TeleportWakeCollide", C.nox_xxx_collideTeleportWake_4EAE30, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideTeleportWake_4EAE30(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("FlagCollide", C.sub_4EA400, func(u *server.Object, a2, a3 uintptr) { sub_4EA400(C.int(uintptr(u.CObj())), C.int(a2)) }, 0)
	server.RegisterObjectCollideGo("ChakramInMotionCollide", C.nox_xxx_collideChakram_4EAF00, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideChakram_4EAF00(C.int(uintptr(u.CObj())), C.int(a2), (*C.float)(unsafe.Pointer(a3)))
	}, 0)
	server.RegisterObjectCollideGo("ArrowCollide", C.nox_xxx_collideArrow_4EB490, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideArrow_4EB490(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("MonsterArrowCollide", C.nox_xxx_collideMonsterArrow_4EB800, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideMonsterArrow_4EB800(C.int(uintptr(u.CObj())), C.int(a2))
	}, 8)
	server.RegisterObjectCollideGo("BearTrapCollide", C.nox_xxx_collideBearTrap_4EB890, func(u *server.Object, a2, a3 uintptr) { nox_xxx_collideBearTrap_4EB890((*C.int)(u.CObj()), C.int(a2)) }, 0)
	server.RegisterObjectCollideGo("PoisonGasTrapCollide", C.nox_xxx_collidePoisonGasTrap_4EB910, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collidePoisonGasTrap_4EB910((*C.int)(u.CObj()), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("TrapDoorCollide", C.nox_xxx_collideTrapDoor_4EAB60, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideTrapDoor_4EAB60(C.int(uintptr(u.CObj())), C.int(a2))
	}, 28)
	server.RegisterObjectCollideGo("BallCollide", C.nox_xxx_collideBall_4EBA00, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideBall_4EBA00(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("HomeBaseCollide", C.nox_xxx_collideHomeBase_4EBB80, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideHomeBase_4EBB80(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("CrownCollide", C.sub_4EBB50, func(u *server.Object, a2, a3 uintptr) { sub_4EBB50(C.int(uintptr(u.CObj())), C.int(a2)) }, 0)
	server.RegisterObjectCollideGo("UndeadKillerCollide", C.nox_xxx_collideUndeadKiller_4EBD40, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideUndeadKiller_4EBD40(C.int(uintptr(u.CObj())), C.int(a2), C.int(a3))
	}, 4)
	server.RegisterObjectCollideGo("YellowStarShotCollide", C.nox_xxx_collideSulphurShot_4E9E50, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideSulphurShot_4E9E50(C.int(uintptr(u.CObj())), C.int(a2), C.int(a3))
	}, 8)
	server.RegisterObjectCollideGo("MimicCollide", C.nox_xxx_collideMimic_4E83D0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideMimic_4E83D0(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("HarpoonCollide", C.nox_xxx_collideHarpoon_4EB6A0, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideHarpoon_4EB6A0((*nox_object_t)(u.CObj()), (*nox_object_t)(unsafe.Pointer(a2)))
	}, 8)
	server.RegisterObjectCollideGo("MonsterGeneratorCollide", C.nox_xxx_collideMonsterGen_4EBE10, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideMonsterGen_4EBE10(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)
	server.RegisterObjectCollideGo("SoulGateCollide", C.sub_4EBE40, func(u *server.Object, a2, a3 uintptr) { sub_4EBE40(C.int(uintptr(u.CObj())), C.int(a2)) }, 4)
	server.RegisterObjectCollideGo("AnkhCollide", C.nox_xxx_collideAnkhQuest_4EBF40, func(u *server.Object, a2, a3 uintptr) {
		nox_xxx_collideAnkhQuest_4EBF40(C.int(uintptr(u.CObj())), C.int(a2))
	}, 0)

	server.RegisterObjectCollideParse("ProjectileCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("ProjectileSparkCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("DamageCollide", resourceObjectParser("collide", "damage"))
	server.RegisterObjectCollideParse("ManaDrainCollide", resourceObjectParser("collide", "mana"))
	server.RegisterObjectCollideParse("SparkExplosionCollide", resourceObjectParser("collide", "spark"))
	server.RegisterObjectCollideParse("WallReflectCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("WallReflectSparkCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("PixieCollide", resourceObjectParser("collide", "projectile"))
	server.RegisterObjectCollideParse("AudioEventCollide", resourceObjectParser("collide", "audio"))
	server.RegisterObjectCollideParse("MonsterArrowCollide", resourceObjectParser("collide", "arrow"))
	server.RegisterObjectCollideParse("YellowStarShotCollide", resourceObjectParser("collide", "projectile"))
}

//export nox_xxx_collideDeathBall_4E9E90
func nox_xxx_collideDeathBall_4E9E90(a1, a2 *nox_object_t, pos *C.float) {
	Nox_xxx_collideDeathBall_4E9E90(asObjectS(a1), asObjectS(a2), (*types.Pointf)(unsafe.Pointer(pos)))
}

func nox_xxx_castCounterSpell_52BBB0(a1 int32, a2, a3, a4 *nox_object_t) {
	Nox_xxx_castCounterSpell_52BBB0(spell.ID(a1), asObjectS(a2), asObjectS(a3), asObjectS(a4), nil, 0)
}

//export nox_xxx_changeOwner_52BE40
func nox_xxx_changeOwner_52BE40(a1, a2 *nox_object_t) {
	Nox_xxx_changeOwner_52BE40(asObjectS(a1), asObjectS(a2))
}
