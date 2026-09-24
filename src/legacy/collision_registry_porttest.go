//go:build porttest

package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
#include "GAME5.h"

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
		{"DefaultCollide", collisionKey(collisionIdentityDefault), 0},
		{"MonsterCollide", collisionKey(collisionIdentityMonster), 0},
		{"PlayerCollide", collisionKey(collisionIdentityPlayer), 0},
		{"ProjectileCollide", collisionKey(collisionIdentityProjectile), 8},
		{"ProjectileSparkCollide", collisionKey(collisionIdentityProjectileSpark), 8},
		{"DoorCollide", collisionKey(collisionIdentityDoor), 0},
		{"PickupCollide", collisionKey(collisionIdentityPickup), 0},
		{"ExitCollide", collisionKey(collisionIdentityExit), 88},
		{"DamageCollide", collisionKey(collisionIdentityDamage), 8},
		{"ManaDrainCollide", collisionKey(collisionIdentityManaDrain), 8},
		{"BombCollide", collisionKey(collisionIdentityBomb), 8},
		{"SparkExplosionCollide", collisionKey(collisionIdentitySparkExplosion), 1},
		{"ChestCollide", collisionKey(collisionIdentityChest), 0},
		{"WallReflectCollide", collisionKey(collisionIdentityWallReflect), 8},
		{"WallReflectSparkCollide", collisionKey(collisionIdentityWallReflectSpark), 8},
		{"PixieCollide", collisionKey(collisionIdentityPixie), 8},
		{"OwnCollide", collisionKey(collisionIdentityOwn), 0},
		{"SparkCollide", collisionKey(collisionIdentitySpark), 8},
		{"BarrelCollide", collisionKey(collisionIdentityBarrel), 0},
		{"AudioEventCollide", collisionKey(collisionIdentityAudioEvent), 4},
		{"TriggerCollide", collisionKey(collisionIdentityTrigger), 0},
		{"TeleportCollide", collisionKey(collisionIdentityTeleport), 8},
		{"ElevatorCollide", collisionKey(collisionIdentityDefault), 8},
		{"AwardSpellCollide", collisionKey(collisionIdentityAwardSpell), 4},
		{"DieCollide", collisionKey(collisionIdentityDie), 0},
		{"GlyphCollide", collisionKey(collisionIdentityGlyph), 0},
		{"SpellProjectileCollide", collisionKey(collisionIdentitySpellProjectile), 0},
		{"BoomCollide", collisionKey(collisionIdentityBoom), 0},
		{"SignCollide", collisionKey(collisionIdentitySign), 0},
		{"PentagramCollide", collisionKey(collisionIdentityPentagram), 0},
		{"SpiderSpitCollide", collisionKey(collisionIdentitySpiderSpit), 0},
		{"DeathBallCollide", collisionKey(collisionIdentityDeathBall), 0},
		{"DeathBallFragmentCollide", collisionKey(collisionIdentityDeathBallFragment), 0},
		{"TelekinesisCollide", collisionKey(collisionIdentityDefault), 0},
		{"FistCollide", collisionKey(collisionIdentityFist), 0},
		{"TeleportWakeCollide", collisionKey(collisionIdentityTeleportWake), 8},
		{"FlagCollide", collisionKey(collisionIdentityFlag), 0},
		{"ChakramInMotionCollide", collisionKey(collisionIdentityChakramInMotion), 0},
		{"ArrowCollide", collisionKey(collisionIdentityArrow), 8},
		{"MonsterArrowCollide", collisionKey(collisionIdentityMonsterArrow), 8},
		{"BearTrapCollide", collisionKey(collisionIdentityBearTrap), 0},
		{"PoisonGasTrapCollide", collisionKey(collisionIdentityPoisonGasTrap), 0},
		{"TrapDoorCollide", collisionKey(collisionIdentityTrapDoor), 28},
		{"BallCollide", collisionKey(collisionIdentityBall), 0},
		{"HomeBaseCollide", collisionKey(collisionIdentityHomeBase), 0},
		{"CrownCollide", collisionKey(collisionIdentityCrown), 0},
		{"UndeadKillerCollide", collisionKey(collisionIdentityUndeadKiller), 4},
		{"YellowStarShotCollide", collisionKey(collisionIdentityYellowStarShot), 8},
		{"MimicCollide", collisionKey(collisionIdentityMimic), 0},
		{"HarpoonCollide", collisionKey(collisionIdentityHarpoon), 8},
		{"MonsterGeneratorCollide", collisionKey(collisionIdentityMonsterGenerator), 0},
		{"SoulGateCollide", collisionKey(collisionIdentitySoulGate), 4},
		{"AnkhCollide", collisionKey(collisionIdentityAnkh), 0},
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
