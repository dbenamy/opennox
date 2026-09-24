//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type portTestUpdateRegistration struct {
	name    string
	address unsafe.Pointer
	size    uintptr
}

func portTestUpdateRegistrations() []portTestUpdateRegistration {
	return []portTestUpdateRegistration{
		{"PlayerUpdate", updateIdentityKey(updateIDPlayer), unsafe.Sizeof(server.PlayerUpdateData{})},
		{"ProjectileUpdate", updateIdentityKey(updateIDProjectile), 0},
		{"SpellProjectileUpdate", updateIdentityKey(updateIDSpellProjectile), unsafe.Sizeof(server.SpellProjectileUpdateData{})},
		{"AntiSpellProjectileUpdate", updateIdentityKey(updateIDAntiSpellProjectile), 28},
		{"DoorUpdate", updateIdentityKey(updateIDDoor), 52},
		{"SparkUpdate", updateIdentityKey(updateIDSpark), 16},
		{"ProjectileTrailUpdate", updateIdentityKey(updateIDProjectileTrail), 0},
		{"PushUpdate", updateIdentityKey(updateIDPush), 12},
		{"TriggerUpdate", updateIdentityKey(updateIDTrigger), 60},
		{"ToggleUpdate", updateIdentityKey(updateIDToggle), 60},
		{"MonsterUpdate", updateIdentityKey(updateIDMonster), unsafe.Sizeof(server.MonsterUpdateData{})},
		{"LoopAndDamageUpdate", updateIdentityKey(updateIDLoopAndDamage), 16},
		{"ElevatorUpdate", updateIdentityKey(updateIDElevator), 20},
		{"ElevatorShaftUpdate", updateIdentityKey(updateIDElevatorShaft), 16},
		{"PhantomPlayerUpdate", updateIdentityKey(updateIDPhantomPlayer), 0},
		{"ObeliskUpdate", updateIdentityKey(updateIDObelisk), unsafe.Sizeof(server.ObeliskUpdateData{})},
		{"LifetimeUpdate", updateIdentityKey(updateIDLifetime), 4},
		{"MagicMissileUpdate", updateIdentityKey(updateIDMagicMissile), 28},
		{"PixieUpdate", updateIdentityKey(updateIDPixie), 28},
		{"SkullUpdate", updateIdentityKey(updateIDSkull), 52},
		{"PentagramUpdate", updateIdentityKey(updateIDPentagram), 24},
		{"InvisiblePentagramUpdate", updateIdentityKey(updateIDInvisiblePentagram), 24},
		{"SwitchUpdate", updateIdentityKey(updateIDSwitch), 0},
		{"BlowUpdate", updateIdentityKey(updateIDBlow), 0},
		{"MoverUpdate", updateIdentityKey(updateIDMover), 36},
		{"BlackPowderBarrelUpdate", updateIdentityKey(updateIDBlackPowderBarrel), 0},
		{"OneSecondDieUpdate", updateIdentityKey(updateIDOneSecondDie), 0},
		{"WaterBarrelUpdate", updateIdentityKey(updateIDWaterBarrel), 0},
		{"SelfDestructUpdate", updateIdentityKey(updateIDSelfDestruct), 0},
		{"BlackPowderBurnUpdate", updateIdentityKey(updateIDBlackPowderBurn), 0},
		{"DeathBallUpdate", updateIdentityKey(updateIDDeathBall), 0},
		{"DeathBallFragmentUpdate", updateIdentityKey(updateIDDeathBallFragment), 0},
		{"MoonglowUpdate", updateIdentityKey(updateIDMoonglow), 0},
		{"SentryGlobeUpdate", updateIdentityKey(updateIDSentryGlobe), 12},
		{"TelekinesisUpdate", updateIdentityKey(updateIDTelekinesis), 0},
		{"FistUpdate", updateIdentityKey(updateIDFist), 4},
		{"MeteorShowerUpdate", updateIdentityKey(updateIDMeteorShower), 4},
		{"MeteorUpdate", updateIdentityKey(updateIDMeteor), 4},
		{"ToxicCloudUpdate", updateIdentityKey(updateIDToxicCloud), 4},
		{"SmallToxicCloudUpdate", updateIdentityKey(updateIDSmallToxicCloud), 4},
		{"ArachnaphobiaUpdate", updateIdentityKey(updateIDArachnaphobia), 0},
		{"ExpireUpdate", updateIdentityKey(updateIDExpire), 0},
		{"BreakUpdate", updateIdentityKey(updateIDBreak), 0},
		{"OpenUpdate", updateIdentityKey(updateIDOpen), 0},
		{"BreakAndRemoveUpdate", updateIdentityKey(updateIDBreakAndRemove), 0},
		{"ChakramInMotionUpdate", updateIdentityKey(updateIDChakramInMotion), 28},
		{"FlagUpdate", updateIdentityKey(updateIDFlag), 12},
		{"TrapDoorUpdate", updateIdentityKey(updateIDTrapDoor), 0},
		{"BallUpdate", updateIdentityKey(updateIDBall), 32},
		{"CrownUpdate", updateIdentityKey(updateIDCrown), 12},
		{"UndeadKillerUpdate", updateIdentityKey(updateIDUndeadKiller), 0},
		{"HarpoonUpdate", updateIdentityKey(updateIDHarpoon), 4},
		{"MonsterGeneratorUpdate", updateIdentityKey(updateIDMonsterGenerator), 164},
	}
}
func portTestUpdateRegistryPointer(name string) (unsafe.Pointer, uint32) {
	for i, entry := range portTestUpdateRegistrations() {
		if entry.name != name {
			continue
		}
		ptr, size := server.PortTestUnitGameplayRegistration(name, false)
		if ptr != entry.address || size != entry.size {
			panic("update registration identity/size mismatch: " + name)
		}
		return ptr, 99400 + uint32(i)
	}
	panic("unmapped update callback: " + name)
}
func PortTestUpdateRegistryNames() []string {
	var names []string
	for _, entry := range portTestUpdateRegistrations() {
		portTestUpdateRegistryPointer(entry.name)
		names = append(names, entry.name)
	}
	return names
}

var portTestUpdateCounts = make(map[string]uint64)

// Invokes the actual configured owner method on both sides of the conversion.
func PortTestRegisteredUpdate(u *server.Object, name string) {
	ptr, _ := portTestUpdateRegistryPointer(name)
	u.Update = ptr
	u.CallUpdate()
	portTestUpdateCounts[name]++
}
func PortTestUpdateRegistryTakeCounts() map[string]uint64 {
	out := portTestUpdateCounts
	portTestUpdateCounts = make(map[string]uint64)
	return out
}

var portTestUpdateActions = map[int]string{
	600: "SparkUpdate",
	601: "ProjectileTrailUpdate",
	602: "LifetimeUpdate",
	603: "SpellProjectileUpdate",
	604: "AntiSpellProjectileUpdate",
	606: "MagicMissileUpdate",
	607: "BlackPowderBarrelUpdate",
	608: "OneSecondDieUpdate",
	609: "WaterBarrelUpdate",
	611: "SelfDestructUpdate",
	612: "BlackPowderBurnUpdate",
	613: "DeathBallFragmentUpdate",
	614: "MoonglowUpdate",
	615: "TelekinesisUpdate",
	616: "FistUpdate",
	618: "MeteorShowerUpdate",
	619: "MeteorUpdate",
	620: "ToxicCloudUpdate",
	622: "SmallToxicCloudUpdate",
	624: "ArachnaphobiaUpdate",
	625: "ExpireUpdate",
	626: "BreakUpdate",
	627: "OpenUpdate",
	628: "BreakAndRemoveUpdate",
	629: "ChakramInMotionUpdate",
	700: "DoorUpdate",
	701: "PushUpdate",
	702: "ToggleUpdate",
	703: "TriggerUpdate",
	704: "LoopAndDamageUpdate",
	705: "SwitchUpdate",
	706: "ElevatorShaftUpdate",
	709: "ElevatorUpdate",
	711: "PhantomPlayerUpdate",
	712: "PentagramUpdate",
	714: "InvisiblePentagramUpdate",
	716: "BlowUpdate",
	718: "TrapDoorUpdate",
	808: "ObeliskUpdate",
	809: "FlagUpdate",
	810: "BallUpdate",
	811: "CrownUpdate",
}

func (p *portTestShopPools) registeredUpdateAction(a PortTestShopAction) bool {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	if !sp.RegisteredUpdates {
		return false
	}
	name, ok := portTestUpdateActions[a.Op]
	if !ok {
		return false
	}
	ptr, id := portTestUpdateRegistryPointer(name)
	p.identify(ptr, id)
	PortTestRegisteredUpdate(p.items[a.Item].u, name)
	p.temporary.result = 0 // Object.CallUpdate intentionally discards callback returns.
	return true
}
