//go:build porttest

package legacy

/*
#include <stdint.h>
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
		{"DefaultXfer", xferIdentityKey(xferIDDefault)},
		{"SpellPagePedestalXfer", xferIdentityKey(xferIDSpellPagePedestal)},
		{"SpellRewardXfer", xferIdentityKey(xferIDSpellReward)},
		{"AbilityRewardXfer", xferIdentityKey(xferIDAbilityReward)},
		{"FieldGuideXfer", xferIdentityKey(xferIDFieldGuide)},
		{"ReadableXfer", xferIdentityKey(xferIDReadable)},
		{"ExitXfer", xferIdentityKey(xferIDExit)},
		{"DoorXfer", xferIdentityKey(xferIDDoor)},
		{"TriggerXfer", xferIdentityKey(xferIDTrigger)},
		{"MonsterXfer", xferIdentityKey(xferIDMonster)},
		{"HoleXfer", xferIdentityKey(xferIDHole)},
		{"TransporterXfer", xferIdentityKey(xferIDTransporter)},
		{"ElevatorXfer", xferIdentityKey(xferIDElevator)},
		{"ElevatorShaftXfer", xferIdentityKey(xferIDElevatorShaft)},
		{"MoverXfer", xferIdentityKey(xferIDMover)},
		{"GlyphXfer", xferIdentityKey(xferIDGlyph)},
		{"InvisibleLightXfer", xferIdentityKey(xferIDInvisibleLight)},
		{"SentryXfer", xferIdentityKey(xferIDSentry)},
		{"WeaponXfer", xferIdentityKey(xferIDWeapon)},
		{"ArmorXfer", xferIdentityKey(xferIDArmor)},
		{"TeamXfer", xferIdentityKey(xferIDTeam)},
		{"GoldXfer", xferIdentityKey(xferIDGold)},
		{"AmmoXfer", xferIdentityKey(xferIDAmmo)},
		{"NPCXfer", xferIdentityKey(xferIDNPC)},
		{"ObeliskXfer", xferIdentityKey(xferIDObelisk)},
		{"ToxicCloudXfer", xferIdentityKey(xferIDToxicCloud)},
		{"MonsterGeneratorXfer", xferIdentityKey(xferIDMonsterGenerator)},
		{"RewardMarkerXfer", xferIdentityKey(xferIDRewardMarker)},
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
