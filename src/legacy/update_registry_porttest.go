//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "server__magic__plyrspel.h"


void nox_xxx_updateProjectile_53AC10(nox_object_t* a1);
void nox_xxx_updateDeathBall_53D080(nox_object_t* a1);
void nox_xxx___mkgmtime_538280(nox_object_t* a1);
*/
import "C"

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
		{"PlayerUpdate", C.nox_xxx_updatePlayer_4F8100, unsafe.Sizeof(server.PlayerUpdateData{})},
		{"ProjectileUpdate", C.nox_xxx_updateProjectile_53AC10, 0},
		{"SpellProjectileUpdate", C.nox_xxx_spellFlyUpdate_53B940, unsafe.Sizeof(server.SpellProjectileUpdateData{})},
		{"AntiSpellProjectileUpdate", C.nox_xxx_updateAntiSpellProj_53BB00, 28},
		{"DoorUpdate", C.nox_xxx_updateDoor_53AC50, 52},
		{"SparkUpdate", C.nox_xxx_updateSpark_53ADC0, 16},
		{"ProjectileTrailUpdate", C.nox_xxx_updateProjTrail_53AEC0, 0},
		{"PushUpdate", C.nox_xxx_updatePush_53B030, 12},
		{"TriggerUpdate", C.nox_xxx_updateTrigger_53B1B0, 60},
		{"ToggleUpdate", C.nox_xxx_updateToggle_53B060, 60},
		{"MonsterUpdate", C.nox_xxx_unitUpdateMonster_50A5C0, unsafe.Sizeof(server.MonsterUpdateData{})},
		{"LoopAndDamageUpdate", C.sub_53B300, 16},
		{"ElevatorUpdate", C.nox_xxx_updateElevator_53B5D0, 20},
		{"ElevatorShaftUpdate", C.nox_xxx_updateElevatorShaft_53B380, 16},
		{"PhantomPlayerUpdate", C.nox_xxx_updatePhantomPlayer_53B860, 0},
		{"ObeliskUpdate", C.nox_xxx_updateObelisk_53C580, unsafe.Sizeof(server.ObeliskUpdateData{})},
		{"LifetimeUpdate", C.nox_xxx_updateLifetime_53B8F0, 4},
		{"MagicMissileUpdate", C.nox_xxx_updateMagicMissile_53BDA0, 28},
		{"PixieUpdate", C.nox_xxx_updatePixie_53CD20, 28},
		{"SkullUpdate", C.nox_xxx_updateShootingTrap_54F9A0, 52},
		{"PentagramUpdate", C.nox_xxx_updateTeleportPentagram_53BEF0, 24},
		{"InvisiblePentagramUpdate", C.nox_xxx_updateInvisiblePentagram_53C0C0, 24},
		{"SwitchUpdate", C.nox_xxx_updateSwitch_53B320, 0},
		{"BlowUpdate", C.nox_xxx_updateBlow_53C160, 0},
		{"MoverUpdate", C.nox_xxx_unitUpdateMover_54F740, 36},
		{"BlackPowderBarrelUpdate", C.nox_xxx_updateBlackPowderBarrel_53C9A0, 0},
		{"OneSecondDieUpdate", C.nox_xxx_updateOneSecondDie_53CB60, 0},
		{"WaterBarrelUpdate", C.nox_xxx_updateWaterBarrel_53CB90, 0},
		{"SelfDestructUpdate", C.nox_xxx_updateSelfDestruct_53CC90, 0},
		{"BlackPowderBurnUpdate", C.nox_xxx_updateBlackPowderBurn_53CCB0, 0},
		{"DeathBallUpdate", C.nox_xxx_updateDeathBall_53D080, 0},
		{"DeathBallFragmentUpdate", C.nox_xxx_updateDeathBallFragment_53D220, 0},
		{"MoonglowUpdate", C.nox_xxx_updateMoonglow_53D270, 0},
		{"SentryGlobeUpdate", C.nox_xxx_updateSentryGlobe_510E60, 12},
		{"TelekinesisUpdate", C.nox_xxx_updateTelekinesis_53D330, 0},
		{"FistUpdate", C.nox_xxx_updateFist_53D400, 4},
		{"MeteorShowerUpdate", C.nox_xxx_updateMeteorShower_53D5A0, 4},
		{"MeteorUpdate", C.nox_xxx_meteorExplode_53D6E0, 4},
		{"ToxicCloudUpdate", C.nox_xxx_updateToxicCloud_53D850, 4},
		{"SmallToxicCloudUpdate", C.nox_xxx_updateSmallToxicCloud_53D960, 4},
		{"ArachnaphobiaUpdate", C.nox_xxx_updateArachnaphobia_53DA60, 0},
		{"ExpireUpdate", C.nox_xxx_updateExpire_53DB00, 0},
		{"BreakUpdate", C.nox_xxx_updateBreak_53DB30, 0},
		{"OpenUpdate", C.nox_xxx_updateOpen_53DBB0, 0},
		{"BreakAndRemoveUpdate", C.nox_xxx_updateBreakAndRemove_53DC30, 0},
		{"ChakramInMotionUpdate", C.nox_xxx_updateChakramInMotion_53DCC0, 28},
		{"FlagUpdate", C.nox_xxx_updateFlag_53DDF0, 12},
		{"TrapDoorUpdate", C.nox_xxx_updateTrapDoor_53DE80, 0},
		{"BallUpdate", C.nox_xxx_updateGameBall_53DF40, 32},
		{"CrownUpdate", C.nox_xxx_updateCrown_53E1D0, 12},
		{"UndeadKillerUpdate", C.nox_xxx_updateUndeadKiller_53E190, 0},
		{"HarpoonUpdate", C.nox_xxx_updateHarpoon_54F380, 4},
		{"MonsterGeneratorUpdate", C.nox_xxx_updateMonsterGenerator_54E930, 164},
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
