//go:build porttest

package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
#include "GAME5.h"

int nox_objectCollideDefault(int a1, int a2, float* a3);
void nox_xxx_collideDeathBall_4E9E90(nox_object_t* a1, nox_object_t* a2, float* a3);
static unsigned int collisionRegistryRaw[4];
static void collisionRegistryObserve(void* u, void* target, void* normal) {
 collisionRegistryRaw[0]++;
 collisionRegistryRaw[1]=(uintptr_t)u;
 collisionRegistryRaw[2]=(uintptr_t)target;
 collisionRegistryRaw[3]=(uintptr_t)normal;
}
static void* collisionRegistryObservePtr(void) { return collisionRegistryObserve; }
static unsigned int* collisionRegistryObserveData(void) { return collisionRegistryRaw; }
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"runtime"
	"unsafe"
)

type portTestCollisionEntry struct {
	Name    string
	Pointer unsafe.Pointer
	Size    uintptr
}

func portTestCollisionEntries() []portTestCollisionEntry {
	return []portTestCollisionEntry{
		{"DefaultCollide", C.nox_objectCollideDefault, 0},
		{"MonsterCollide", C.nox_xxx_collideMonsterEventProc_4E83B0, 0},
		{"PlayerCollide", C.nox_xxx_collidePlayer_4E8460, 0},
		{"ProjectileCollide", C.nox_xxx_collideProjectileGeneric_4E87B0, 8},
		{"ProjectileSparkCollide", C.nox_xxx_collideProjectileSpark_4E8880, 8},
		{"DoorCollide", C.nox_xxx_collideDoor_4E8AC0, 0},
		{"PickupCollide", C.nox_xxx_collidePickup_4E8DF0, 0},
		{"ExitCollide", C.nox_xxx_collideExit_4E9090, 88},
		{"DamageCollide", C.nox_xxx_collideDamage_4E9430, 8},
		{"ManaDrainCollide", C.nox_xxx_collideManadrain_4E9490, 8},
		{"BombCollide", C.nox_xxx_collideBomb_4E96F0, 8},
		{"SparkExplosionCollide", C.nox_xxx_fireballCollide_4E9AC0, 1},
		{"ChestCollide", C.nox_xxx_collideChest_4E9C40, 0},
		{"WallReflectCollide", C.nox_xxx_collideSulphurShot2_4E9D80, 8},
		{"WallReflectSparkCollide", C.nox_xxx_collideWallReflectSpark_4EA200, 8},
		{"PixieCollide", C.nox_xxx_collidePixie_4EA080, 8},
		{"OwnCollide", C.sub_4EA2C0, 0},
		{"SparkCollide", C.nox_xxx_collideSpark_4EA300, 8},
		{"BarrelCollide", C.sub_4EAAA0, 0},
		{"AudioEventCollide", C.sub_4EAAD0, 4},
		{"TriggerCollide", C.nox_xxx_collideTrigger_54FCD0, 0},
		{"TeleportCollide", C.sub_4EACA0, 8},
		{"ElevatorCollide", C.nox_objectCollideDefault, 8},
		{"AwardSpellCollide", C.nox_xxx_collideSpellPedestal_4EAD20, 4},
		{"DieCollide", C.nox_xxx_collideDie_4E99B0, 0},
		{"GlyphCollide", C.nox_xxx_collideGlyph_4E9A00, 0},
		{"SpellProjectileCollide", C.nox_xxx_spellFlyCollide_4E9500, 0},
		{"BoomCollide", C.nox_xxx_collideBoom_4E9770, 0},
		{"SignCollide", C.nox_xxx_collideSign_4EAB40, 0},
		{"PentagramCollide", C.nox_xxx_collidePentagram_4EAB20, 0},
		{"SpiderSpitCollide", C.nox_xxx_collideWebbing_4EA380, 0},
		{"DeathBallCollide", C.nox_xxx_collideDeathBall_4E9E90, 0},
		{"DeathBallFragmentCollide", C.nox_xxx_collideDeathBallFragment_4E9FE0, 0},
		{"TelekinesisCollide", C.nox_objectCollideDefault, 0},
		{"FistCollide", C.nox_xxx_collideFist_4EADF0, 0},
		{"TeleportWakeCollide", C.nox_xxx_collideTeleportWake_4EAE30, 8},
		{"FlagCollide", C.sub_4EA400, 0},
		{"ChakramInMotionCollide", C.nox_xxx_collideChakram_4EAF00, 0},
		{"ArrowCollide", C.nox_xxx_collideArrow_4EB490, 8},
		{"MonsterArrowCollide", C.nox_xxx_collideMonsterArrow_4EB800, 8},
		{"BearTrapCollide", C.nox_xxx_collideBearTrap_4EB890, 0},
		{"PoisonGasTrapCollide", C.nox_xxx_collidePoisonGasTrap_4EB910, 0},
		{"TrapDoorCollide", C.nox_xxx_collideTrapDoor_4EAB60, 28},
		{"BallCollide", C.nox_xxx_collideBall_4EBA00, 0},
		{"HomeBaseCollide", C.nox_xxx_collideHomeBase_4EBB80, 0},
		{"CrownCollide", C.sub_4EBB50, 0},
		{"UndeadKillerCollide", C.nox_xxx_collideUndeadKiller_4EBD40, 4},
		{"YellowStarShotCollide", C.nox_xxx_collideSulphurShot_4E9E50, 8},
		{"MimicCollide", C.nox_xxx_collideMimic_4E83D0, 0},
		{"HarpoonCollide", C.nox_xxx_collideHarpoon_4EB6A0, 8},
		{"MonsterGeneratorCollide", C.nox_xxx_collideMonsterGen_4EBE10, 0},
		{"SoulGateCollide", C.sub_4EBE40, 4},
		{"AnkhCollide", C.nox_xxx_collideAnkhQuest_4EBF40, 0},
	}
}
func PortTestCollisionRegistryNames() []string {
	var out []string
	for _, e := range portTestCollisionEntries() {
		PortTestCollisionRegistryPointer(e.Name)
		out = append(out, e.Name)
	}
	return out
}
func PortTestCollisionRegistryPointer(name string) unsafe.Pointer {
	for _, e := range portTestCollisionEntries() {
		if e.Name == name {
			p, size := server.PortTestWorldCollisionRegistry(name)
			if p != e.Pointer || size != e.Size {
				panic("collision registry identity/size: " + name)
			}
			return p
		}
	}
	panic("unknown collision registry name: " + name)
}

var portTestCollisionCounts = map[string]int{}

func PortTestRegisteredCollision(u, target *server.Object, normal *types.Pointf, name string) {
	old := u.Collide
	u.Collide = PortTestCollisionRegistryPointer(name)
	defer func() { u.Collide = old }()
	server.PortTestCollisionWith(u, target, normal)
	runtime.KeepAlive(target)
	runtime.KeepAlive(normal)
	portTestCollisionCounts[name]++
}
func PortTestCollisionRegistryTakeCounts() map[string]int {
	out := portTestCollisionCounts
	portTestCollisionCounts = map[string]int{}
	return out
}

func PortTestCollisionRawObserver() (unsafe.Pointer, *[4]uint32, func()) {
	p := (*[4]uint32)(unsafe.Pointer(C.collisionRegistryObserveData()))
	old := *p
	*p = [4]uint32{}
	return C.collisionRegistryObservePtr(), p, func() { *p = old }
}
