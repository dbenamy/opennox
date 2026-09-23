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
		{"PlayerDie", C.nox_xxx_diePlayer_54D2B0, 0},
		{"PotionDie", C.nox_xxx_diePotion_54CBB0, 0},
		{"ImpEggDie", C.nox_xxx_dieImpEgg_54CAE0, 0},
		{"GlyphDie", C.nox_xxx_dieGlyph_54DF30, 0},
		{"BarrelDie", C.nox_xxx_dieBarrel_54DFA0, 0},
		{"CreateObjectDie", C.nox_xxx_dieCreateObject_54E010, 132},
		{"SpawnObjectDie", C.nox_xxx_dieSpawnObject_54E070, 132},
		{"PolypDie", C.nox_xxx_diePolyp_54CB10, 0},
		{"MarkerDie", C.nox_xxx_dieMarker_54E460, 0},
		{"WeaponDie", C.nox_xxx_dieWeapon_54E370_obj_die, 0},
		{"ArmorDie", C.nox_xxx_dieArmor_54E170_obj_die, 0},
		{"BoulderDie", C.nox_xxx_dieBoulder_54E4B0, 0},
		{"GameBallDie", C.nox_xxx_dieGameBall_54E620, 0},
		{"MonsterGeneratorDie", C.nox_xxx_dieMonsterGen_54E630, 0},
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
