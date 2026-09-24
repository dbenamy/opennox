//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
int nox_xxx_mapReadWriteObjData_4F4530(nox_object_t* a1p, int a2);
int nox_xxx_xfer_4F3E30(unsigned short a1, nox_object_t* a2, int a3);
#include "GAME3_2.h"
static int xferSoundObserveReturn;
static unsigned int xferSoundObserveWords[3];
static int xferSoundObserve(void* a,void* b) {
 xferSoundObserveWords[0]++;xferSoundObserveWords[1]=(uintptr_t)a;xferSoundObserveWords[2]=(uintptr_t)b;
 return xferSoundObserveReturn;
}
static void* xferSoundObservePtr(void) {return xferSoundObserve;}
static unsigned int* xferSoundObserveData(void) {return xferSoundObserveWords;}
static int* xferSoundObserveResult(void) {return &xferSoundObserveReturn;}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestXferRegistryNames() []string {
	var out []string
	for _, e := range portTestXferEntries() {
		PortTestXferRegistryPointer(e.Name)
		out = append(out, e.Name)
	}
	return out
}

type portTestXferEntry struct {
	Name    string
	Pointer unsafe.Pointer
}

func portTestXferEntries() []portTestXferEntry {
	return []portTestXferEntry{
		{"DefaultXfer", C.nox_xxx_XFerDefault_4F49A0},
		{"SpellPagePedestalXfer", C.nox_xxx_XFerSpellPagePedistal_4F4A20},
		{"SpellRewardXfer", C.nox_xxx_XFerSpellReward_4F5F30},
		{"AbilityRewardXfer", C.nox_xxx_XFerAbilityReward_4F6240},
		{"FieldGuideXfer", C.nox_xxx_XFerFieldGuide_4F6390},
		{"ReadableXfer", C.nox_xxx_XFerReadable_4F4AB0},
		{"ExitXfer", C.nox_xxx_XFerExit_4F4B90},
		{"DoorXfer", C.nox_xxx_XFerDoor_4F4CB0},
		{"TriggerXfer", C.nox_xxx_unitTriggerXfer_4F4E50},
		{"MonsterXfer", C.nox_xxx_XFerMonster_528DB0},
		{"HoleXfer", C.nox_xxx_XFerHole_4F51D0},
		{"TransporterXfer", C.nox_xxx_XFerTransporter_4F5300},
		{"ElevatorXfer", C.nox_xxx_XFerElevator_4F53D0},
		{"ElevatorShaftXfer", C.nox_xxx_XFerElevatorShaft_4F54A0},
		{"MoverXfer", C.nox_xxx_XFerMover_4F5730},
		{"GlyphXfer", C.nox_xxx_XFerGlyph_4F5890},
		{"InvisibleLightXfer", C.nox_xxx_XFerInvLight_4F5AA0},
		{"SentryXfer", C.nox_xxx_XFerSentry_4F5E50},
		{"WeaponXfer", C.nox_xxx_XFerWeapon_4F64A0},
		{"ArmorXfer", C.nox_xxx_XFerArmor_4F6860},
		{"TeamXfer", C.nox_xxx_XFerTeam_4F6D20},
		{"GoldXfer", C.nox_xxx_XFerGold_4F6EC0},
		{"AmmoXfer", C.nox_xxx_XFerAmmo_4F6B20},
		{"NPCXfer", C.nox_xxx_XFerNPC_52ADE0},
		{"ObeliskXfer", C.nox_xxx_XFerObelisk_4F6F60},
		{"ToxicCloudXfer", C.nox_xxx_XFerToxicCloud_4F70A0},
		{"MonsterGeneratorXfer", C.nox_xxx_XFerMonsterGen_4F7130},
		{"RewardMarkerXfer", C.nox_xxx_XFerRewardMarker_4F74D0},
	}
}
func PortTestXferRegistryPointer(name string) unsafe.Pointer {
	for _, e := range portTestXferEntries() {
		if e.Name == name {
			p := server.PortTestXferSoundRegistry(name, false)
			if p != e.Pointer {
				panic("xfer registry identity: " + name)
			}
			return p
		}
	}
	panic("unknown xfer registry name: " + name)
}
func PortTestDamageSoundRegistryPointer(name string) unsafe.Pointer {
	var want unsafe.Pointer
	switch name {
	case "DefaultDamageSound":
		want = damageIdentityKey(damageIDDefaultSound)
	case "PlayerDamageSound":
		want = damageIdentityKey(damageIDPlayerSound)
	default:
		panic("unknown damage sound registry name")
	}
	p := server.PortTestXferSoundRegistry(name, true)
	if p != want {
		panic("damage sound registry identity")
	}
	return p
}
func PortTestXferSoundRawObserver() (unsafe.Pointer, *[3]uint32, *int32, func()) {
	words := (*[3]uint32)(unsafe.Pointer(C.xferSoundObserveData()))
	result := (*int32)(unsafe.Pointer(C.xferSoundObserveResult()))
	oldWords, oldResult := *words, *result
	*words = [3]uint32{}
	*result = 0
	return C.xferSoundObservePtr(), words, result, func() { *words = oldWords; *result = oldResult }
}

// Exercise the unchanged owner that chooses and invokes the damage sound slot.
func PortTestDamageSoundOwner(u, source, weapon *server.Object) int32 {
	return damageDefault(u, source, weapon, 12, 0)
}
