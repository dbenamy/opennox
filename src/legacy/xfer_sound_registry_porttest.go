//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

var portTestXferSoundKey byte
var portTestXferSoundWords [3]uint32
var portTestXferSoundResult int32

func portTestXferSoundRecord(a, b unsafe.Pointer) int32 {
	portTestXferSoundWords[0]++
	portTestXferSoundWords[1], portTestXferSoundWords[2] = uint32(uintptr(a)), uint32(uintptr(b))
	return portTestXferSoundResult
}
func init() {
	key := unsafe.Pointer(&portTestXferSoundKey)
	server.PortTestRegisterXferCallback(key, func(u *server.Object, arg unsafe.Pointer) int {
		return int(portTestXferSoundRecord(unsafe.Pointer(u), arg))
	})
	server.PortTestRegisterDamageSoundCallback(key, func(u, v *server.Object) { portTestXferSoundRecord(unsafe.Pointer(u), unsafe.Pointer(v)) })
	server.PortTestRegisterUseCallback(key, func(u, v *server.Object) int32 { return portTestXferSoundRecord(unsafe.Pointer(u), unsafe.Pointer(v)) })
	client.RegisterDrawableUpdateCallbackGo(key, func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
		return portTestXferSoundRecord(unsafe.Pointer(vp), unsafe.Pointer(dr))
	})
}

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
	words := &portTestXferSoundWords
	result := &portTestXferSoundResult
	oldWords, oldResult := *words, *result
	*words = [3]uint32{}
	*result = 0
	return unsafe.Pointer(&portTestXferSoundKey), words, result, func() { *words = oldWords; *result = oldResult }
}

// Exercise the unchanged owner that chooses and invokes the damage sound slot.
func PortTestDamageSoundOwner(u, source, weapon *server.Object) int32 {
	return damageDefault(u, source, weapon, 12, 0)
}
