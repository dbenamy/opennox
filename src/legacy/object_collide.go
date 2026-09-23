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
	server.RegisterObjectCollide("DefaultCollide", C.nox_objectCollideDefault, 0)
	server.RegisterObjectCollide("MonsterCollide", C.nox_xxx_collideMonsterEventProc_4E83B0, 0)
	server.RegisterObjectCollide("PlayerCollide", C.nox_xxx_collidePlayer_4E8460, 0)
	server.RegisterObjectCollide("ProjectileCollide", C.nox_xxx_collideProjectileGeneric_4E87B0, 8)
	server.RegisterObjectCollide("ProjectileSparkCollide", C.nox_xxx_collideProjectileSpark_4E8880, 8)
	server.RegisterObjectCollide("DoorCollide", C.nox_xxx_collideDoor_4E8AC0, 0)
	server.RegisterObjectCollide("PickupCollide", C.nox_xxx_collidePickup_4E8DF0, 0)
	server.RegisterObjectCollide("ExitCollide", C.nox_xxx_collideExit_4E9090, 88)
	server.RegisterObjectCollide("DamageCollide", C.nox_xxx_collideDamage_4E9430, 8)
	server.RegisterObjectCollide("ManaDrainCollide", C.nox_xxx_collideManadrain_4E9490, 8)
	server.RegisterObjectCollide("BombCollide", C.nox_xxx_collideBomb_4E96F0, 8)
	server.RegisterObjectCollide("SparkExplosionCollide", C.nox_xxx_fireballCollide_4E9AC0, 1)
	server.RegisterObjectCollide("ChestCollide", C.nox_xxx_collideChest_4E9C40, 0)
	server.RegisterObjectCollide("WallReflectCollide", C.nox_xxx_collideSulphurShot2_4E9D80, 8)
	server.RegisterObjectCollide("WallReflectSparkCollide", C.nox_xxx_collideWallReflectSpark_4EA200, 8)
	server.RegisterObjectCollide("PixieCollide", C.nox_xxx_collidePixie_4EA080, 8)
	server.RegisterObjectCollide("OwnCollide", C.sub_4EA2C0, 0)
	server.RegisterObjectCollide("SparkCollide", C.nox_xxx_collideSpark_4EA300, 8)
	server.RegisterObjectCollide("BarrelCollide", C.sub_4EAAA0, 0)
	server.RegisterObjectCollide("AudioEventCollide", C.sub_4EAAD0, 4)
	server.RegisterObjectCollide("TriggerCollide", C.nox_xxx_collideTrigger_54FCD0, 0)
	server.RegisterObjectCollide("TeleportCollide", C.sub_4EACA0, 8)
	server.RegisterObjectCollide("ElevatorCollide", C.nox_objectCollideDefault, 8)
	server.RegisterObjectCollide("AwardSpellCollide", C.nox_xxx_collideSpellPedestal_4EAD20, 4)
	server.RegisterObjectCollide("DieCollide", C.nox_xxx_collideDie_4E99B0, 0)
	server.RegisterObjectCollide("GlyphCollide", C.nox_xxx_collideGlyph_4E9A00, 0)
	server.RegisterObjectCollide("SpellProjectileCollide", C.nox_xxx_spellFlyCollide_4E9500, 0)
	server.RegisterObjectCollide("BoomCollide", C.nox_xxx_collideBoom_4E9770, 0)
	server.RegisterObjectCollide("SignCollide", C.nox_xxx_collideSign_4EAB40, 0)
	server.RegisterObjectCollide("PentagramCollide", C.nox_xxx_collidePentagram_4EAB20, 0)
	server.RegisterObjectCollide("SpiderSpitCollide", C.nox_xxx_collideWebbing_4EA380, 0)
	server.RegisterObjectCollide("DeathBallCollide", C.nox_xxx_collideDeathBall_4E9E90, 0)
	server.RegisterObjectCollide("DeathBallFragmentCollide", C.nox_xxx_collideDeathBallFragment_4E9FE0, 0)
	server.RegisterObjectCollide("TelekinesisCollide", C.nox_objectCollideDefault, 0)
	server.RegisterObjectCollide("FistCollide", C.nox_xxx_collideFist_4EADF0, 0)
	server.RegisterObjectCollide("TeleportWakeCollide", C.nox_xxx_collideTeleportWake_4EAE30, 8)
	server.RegisterObjectCollide("FlagCollide", C.sub_4EA400, 0)
	server.RegisterObjectCollide("ChakramInMotionCollide", C.nox_xxx_collideChakram_4EAF00, 0)
	server.RegisterObjectCollide("ArrowCollide", C.nox_xxx_collideArrow_4EB490, 8)
	server.RegisterObjectCollide("MonsterArrowCollide", C.nox_xxx_collideMonsterArrow_4EB800, 8)
	server.RegisterObjectCollide("BearTrapCollide", C.nox_xxx_collideBearTrap_4EB890, 0)
	server.RegisterObjectCollide("PoisonGasTrapCollide", C.nox_xxx_collidePoisonGasTrap_4EB910, 0)
	server.RegisterObjectCollide("TrapDoorCollide", C.nox_xxx_collideTrapDoor_4EAB60, 28)
	server.RegisterObjectCollide("BallCollide", C.nox_xxx_collideBall_4EBA00, 0)
	server.RegisterObjectCollide("HomeBaseCollide", C.nox_xxx_collideHomeBase_4EBB80, 0)
	server.RegisterObjectCollide("CrownCollide", C.sub_4EBB50, 0)
	server.RegisterObjectCollide("UndeadKillerCollide", C.nox_xxx_collideUndeadKiller_4EBD40, 4)
	server.RegisterObjectCollide("YellowStarShotCollide", C.nox_xxx_collideSulphurShot_4E9E50, 8)
	server.RegisterObjectCollide("MimicCollide", C.nox_xxx_collideMimic_4E83D0, 0)
	server.RegisterObjectCollide("HarpoonCollide", C.nox_xxx_collideHarpoon_4EB6A0, 8)
	server.RegisterObjectCollide("MonsterGeneratorCollide", C.nox_xxx_collideMonsterGen_4EBE10, 0)
	server.RegisterObjectCollide("SoulGateCollide", C.sub_4EBE40, 4)
	server.RegisterObjectCollide("AnkhCollide", C.nox_xxx_collideAnkhQuest_4EBF40, 0)

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
