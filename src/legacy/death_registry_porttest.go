//go:build porttest

package legacy

/*
#include "server__object__die__die.h"
#include "GAME4_3.h"
#include "GAME5.h"
static void* ptDeathForwardObject;
static int ptDeathForwardCount;
static void ptDeathForward(void* obj) { ptDeathForwardObject=obj; ptDeathForwardCount++; }
static void* ptDeathForwardPtr(void) { return (void*)ptDeathForward; }
static void ptDeathForwardReset(void) { ptDeathForwardObject=0; ptDeathForwardCount=0; }
static void* ptDeathForwardGetObject(void) { return ptDeathForwardObject; }
static int ptDeathForwardGetCount(void) { return ptDeathForwardCount; }
*/
import "C"

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type portTestDeathRegistration struct {
	name    string
	address unsafe.Pointer
	size    uintptr
}

func portTestDeathRegistrations() []portTestDeathRegistration {
	return []portTestDeathRegistration{
		{"PlayerDie", deathKey(deathIdentityPlayer), 0},
		{"PotionDie", deathKey(deathIdentityPotion), 0},
		{"ImpEggDie", deathKey(deathIdentityImpEgg), 0},
		{"GlyphDie", deathKey(deathIdentityGlyph), 0},
		{"BarrelDie", deathKey(deathIdentityBarrel), 0},
		{"CreateObjectDie", deathKey(deathIdentityCreateObject), 132},
		{"SpawnObjectDie", deathKey(deathIdentitySpawnObject), 132},
		{"PolypDie", deathKey(deathIdentityPolyp), 0},
		{"MarkerDie", deathKey(deathIdentityMarker), 0},
		{"WeaponDie", deathKey(deathIdentityWeapon), 0},
		{"ArmorDie", deathKey(deathIdentityArmor), 0},
		{"BoulderDie", deathKey(deathIdentityBoulder), 0},
		{"GameBallDie", deathKey(deathIdentityGameBall), 0},
		{"MonsterGeneratorDie", deathKey(deathIdentityMonsterGenerator), 0},
	}
}

func portTestDeathRegistryPointer(name string) (unsafe.Pointer, uint32) {
	for i, entry := range portTestDeathRegistrations() {
		if entry.name != name {
			continue
		}
		ptr, size := server.PortTestDeathRegistry(name)
		if ptr != entry.address || size != entry.size {
			panic("death registration identity/size mismatch: " + name)
		}
		return ptr, 99300 + uint32(i)
	}
	panic("unmapped death callback: " + name)
}

func PortTestDeathRegistryNames() []string {
	var names []string
	for _, entry := range portTestDeathRegistrations() {
		portTestDeathRegistryPointer(entry.name)
		names = append(names, entry.name)
	}
	return names
}

// This existing production owner is unchanged as a test entrypoint across the
// dispatch refactor. No copied raw invocation can bypass the new registry.
func PortTestDeathProjectile(u, target *server.Object) { projectileDie(u, target) }

func PortTestDeathRegisteredProjectile(u *server.Object, name string) {
	ptr, _ := portTestDeathRegistryPointer(name)
	u.Death = ptr
	target, free := alloc.New(server.Object{})
	target.ObjClass = object.ClassPlayer
	defer free()
	if itemOwnerSameTeam(u, target) {
		panic("death fixture target unexpectedly allied")
	}
	projectileDie(u, target)
}

func PortTestDeathForwardReset() unsafe.Pointer {
	C.ptDeathForwardReset()
	return C.ptDeathForwardPtr()
}
func PortTestDeathForwardSnapshot() (unsafe.Pointer, int) {
	return C.ptDeathForwardGetObject(), int(C.ptDeathForwardGetCount())
}
