//go:build porttest

package legacy

/*
#include "common/alloc/classes/alloc_class.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type portTestLifecycleEntry struct {
	Name    string
	Pointer unsafe.Pointer
	Size    uintptr
	Create  bool
}

func portTestLifecycleEntries() []portTestLifecycleEntry {
	return []portTestLifecycleEntry{
		{"MonsterCreate", C.nox_xxx_monsterCreateFn_54C480, 0, true},
		{"ArmorCreate", C.sub_54C950, 0, true},
		{"WeaponCreate", C.nox_xxx_createWeapon_54C710, 0, true},
		{"ObeliskCreate", C.nox_xxx_createFnObelisk_54CA10, 0, true},
		{"AnimCreate", C.nox_xxx_createFnAnim_54CA50, 0, true},
		{"TriggerCreate", C.nox_xxx_createTrigger_54CA60, 0, true},
		{"MonsterGeneratorCreate", C.nox_xxx_createMonsterGen_54CA90, 0, true},
		{"RewardMarkerCreate", C.nox_xxx_createRewardMarker_54CAC0, 0, true},
		{"MonsterInit", C.nox_xxx_unitMonsterInit_4F0040, 0, false},
		{"PlayerInit", C.nox_xxx_unitInitPlayer_4EFE80, 0, false},
		{"SparkInit", C.nox_xxx_unitSparkInit_4F0390, 0, false},
		{"FrogInit", C.nox_xxx_initFrog_4F03B0, 0, false},
		{"ChestInit", C.nox_xxx_initChest_4F0400, 0, false},
		{"BoulderInit", C.nox_xxx_unitBoulderInit_4F0420, 0, false},
		{"BreakInit", C.nox_xxx_breakInit_4F0570, 0, false},
		{"MonsterGeneratorInit", C.nox_xxx_unitInitGenerator_4F0590, 0, false},
		{"ShopkeeperInit", C.nox_xxx_unitMonsterInit_4F0040, unsafe.Sizeof(server.ShopkeeperInitData{}), false},
		{"SkullInit", C.sub_4F0450, 8, false},
		{"DirectionInit", C.sub_4F0490, 8, false},
		{"GoldInit", C.nox_xxx_unitInitGold_4F04B0, unsafe.Sizeof(server.GoldInitData{}), false},
	}
}
func PortTestLifecycleRegistryNames(create bool) []string {
	var out []string
	for _, v := range portTestLifecycleEntries() {
		if v.Create == create {
			PortTestLifecycleRegistryPointer(v.Name, create)
			out = append(out, v.Name)
		}
	}
	return out
}
func PortTestLifecycleRegistryPointer(name string, create bool) unsafe.Pointer {
	for _, v := range portTestLifecycleEntries() {
		if v.Name == name && v.Create == create {
			p, n := server.PortTestLifecycleRegistry(name, create)
			if p != v.Pointer || n != v.Size {
				panic("lifecycle address/size mismatch: " + name)
			}
			return p
		}
	}
	panic("unknown lifecycle callback: " + name)
}

var portTestInitCounts = map[string]int{}

func PortTestRegisteredInit(u *server.Object, name string) {
	u.Init = PortTestLifecycleRegistryPointer(name, false)
	u.CallInit()
	portTestInitCounts[name]++
}
func PortTestInitRegistryTakeCounts() map[string]int {
	out := portTestInitCounts
	portTestInitCounts = map[string]int{}
	return out
}
