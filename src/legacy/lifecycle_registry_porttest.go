//go:build porttest

package legacy

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
		{"MonsterCreate", lifecycleCreateKey(createIDMonster), 0, true},
		{"ArmorCreate", lifecycleCreateKey(createIDArmor), 0, true},
		{"WeaponCreate", lifecycleCreateKey(createIDWeapon), 0, true},
		{"ObeliskCreate", lifecycleCreateKey(createIDObelisk), 0, true},
		{"AnimCreate", lifecycleCreateKey(createIDAnim), 0, true},
		{"TriggerCreate", lifecycleCreateKey(createIDTrigger), 0, true},
		{"MonsterGeneratorCreate", lifecycleCreateKey(createIDMonsterGenerator), 0, true},
		{"RewardMarkerCreate", lifecycleCreateKey(createIDRewardMarker), 0, true},
		{"MonsterInit", lifecycleInitKey(initIDMonster), 0, false},
		{"PlayerInit", lifecycleInitKey(initIDPlayer), 0, false},
		{"SparkInit", lifecycleInitKey(initIDSpark), 0, false},
		{"FrogInit", lifecycleInitKey(initIDFrog), 0, false},
		{"ChestInit", lifecycleInitKey(initIDChest), 0, false},
		{"BoulderInit", lifecycleInitKey(initIDBoulder), 0, false},
		{"BreakInit", lifecycleInitKey(initIDBreak), 0, false},
		{"MonsterGeneratorInit", lifecycleInitKey(initIDMonsterGenerator), 0, false},
		{"ShopkeeperInit", lifecycleInitKey(initIDMonster), unsafe.Sizeof(server.ShopkeeperInitData{}), false},
		{"SkullInit", lifecycleInitKey(initIDSkull), 8, false},
		{"DirectionInit", lifecycleInitKey(initIDDirection), 8, false},
		{"GoldInit", lifecycleInitKey(initIDGold), unsafe.Sizeof(server.GoldInitData{}), false},
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
