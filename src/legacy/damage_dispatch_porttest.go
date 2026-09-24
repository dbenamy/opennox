//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME3_2.h"
#include "GAME3_3.h"
static uint32_t damageEvents[2048];static int damageN,damageSet;static uint32_t damageOutput;
static void damageReset(int set,uint32_t out){damageN=0;damageSet=set;damageOutput=out;}
static void damageDefend(int m,int it,int u,int a,int w,int* data){
 damageEvents[damageN++]=1;damageEvents[damageN++]=m;damageEvents[damageN++]=it;damageEvents[damageN++]=u;damageEvents[damageN++]=a;damageEvents[damageN++]=w;damageEvents[damageN++]=data[0];damageEvents[damageN++]=data[1];
 if(damageSet)data[0]=damageOutput;
}
static void damageDefendScalar(int m,int it,int u,int a,int w,int* data){
 damageEvents[damageN++]=3;damageEvents[damageN++]=m;damageEvents[damageN++]=it;damageEvents[damageN++]=u;damageEvents[damageN++]=a;damageEvents[damageN++]=w;damageEvents[damageN++]=data[0];
 if(damageSet)data[0]=damageOutput;
}
static void* damageDefendScalarPtr(void){return damageDefendScalar;}
static void damagePre(int m,int u,int a,int it,int* data){
 damageEvents[damageN++]=2;damageEvents[damageN++]=m;damageEvents[damageN++]=u;damageEvents[damageN++]=a;damageEvents[damageN++]=it;damageEvents[damageN++]=*data;
 if(damageSet)*data=damageOutput;
}
static void* damageDefendPtr(void){return damageDefend;}
static void* damagePrePtr(void){return damagePre;}
static int damageCount(void){return damageN;}
static uint32_t damageEvent(int i){return damageEvents[i];}
*/
import "C"
import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestDamageSpec struct {
	Registry                  string
	RegistryValue             bool
	Source, Weapon, Other     int
	Amount, Kind              int32
	PlayerIndex               int
	FloatBits                 uint32
	Name                      string
	DefendMask, PreDamageMask uint32
	SetOutput                 bool
	Output                    uint32
}
type portTestDamage struct{ result uint64 }

var damageTypeNames = []string{"SmallFist", "MediumFist", "LargeFist", "Meteor", "ToxicCloud", "SmallToxicCloud"}
var damageNameOffsets = []uintptr{200800, 200816, 200832, 200848, 200864, 200880, 200896, 200920, 200940, 200952, 200968, 200980, 200996, 201008, 201028, 201044, 201064, 201080}
var damageOffsets = []uintptr{1563316, 1563324, 1563328, 1563332, 1563336, 1563340}

func PortTestDamageNames() (out []string) {
	data := blobdata.PortTestDamageTable()
	for _, off := range damageNameOffsets {
		raw := data[off-200728:]
		out = append(out, string(raw[:bytes.IndexByte(raw, 0)]))
	}
	return
}

func (p *portTestShopPools) damagePrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Damage
	if sp == nil {
		return func() {}
	}
	p.temporary.world.objectives.attack.damage = &portTestDamage{}
	raw := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 200728)), 384)
	oldRaw := bytes.Clone(raw)
	copy(raw, blobdata.PortTestDamageTable())
	var oldNames []unsafe.Pointer
	for i, off := range damageNameOffsets {
		ptr := memmap.PtrPtr(0x587000, 200728+uintptr(i*4))
		oldNames = append(oldNames, *ptr)
		*ptr = memmap.PtrOff(0x587000, off)
	}
	var old []uint32
	for _, off := range damageOffsets {
		old = append(old, *memmap.PtrUint32(0x5d4594, off))
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	cached := dword_5d4594_1563320
	dword_5d4594_1563320 = 0
	return func() {
		for i, old := range oldNames {
			*memmap.PtrPtr(0x587000, 200728+uintptr(i*4)) = old
		}
		for i, off := range damageOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		dword_5d4594_1563320 = cached
		copy(raw, oldRaw)
	}
}

func (p *portTestShopPools) damageItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Damage
	if sp == nil {
		return
	}
	C.damageReset(C.int(bool2int(sp.SetOutput)), C.uint32_t(sp.Output))
	// Retired non-registry adapters occupied fifteen entries in the capture map.
	// Keep later dynamically allocated IDs stable without identifying nil.
	p.reservedFunctionIDs += 15
	for i := 0; i < 26; i++ {
		if key := portTestDamageFunction(i); key != nil {
			p.identify(key, 88000+uint32(i))
		}
	}
	p.identify(C.damageDefendPtr(), 88050)
	p.identify(C.damagePrePtr(), 88051)
	p.identify(C.damageDefendScalarPtr(), 88052)
	for _, id := range []int{1, 2, 100, 101, 102} {
		if u := p.temporaryRef(id); u != nil {
			u.Damage = p.proxy.combat.target.Damage
		}
	}
	if sp.Registry != "" {
		actor := p.temporaryRef(p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Actor)
		if actor == nil {
			panic("damage registry fixture requires an actor")
		}
		actor.Damage = server.PortTestDamageRegistry(sp.Registry)
	}
	for i := 0; i < 4; i++ {
		m := (*server.ModifierEff)(unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*144))
		if sp.DefendMask&(1<<i) != 0 {
			m.Defend76.Fnc = C.damageDefendPtr()
			if i == 1 {
				m.Defend76.Fnc = C.damageDefendScalarPtr()
			}
		}
		if sp.PreDamageMask&(1<<i) != 0 {
			m.AttackPreDmg64.Fnc = C.damagePrePtr()
		}
	}
}

func (p *portTestShopPools) damageAction(a PortTestShopAction) uint32 {
	attack := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	sp := attack.Damage
	state := p.temporary.world.objectives.attack
	if sp.Registry != "" {
		actor := p.temporaryRef(attack.Actor)
		if actor.Damage != portTestDamageFunction(a.Op-1100) {
			panic("damage registry name/address mismatch")
		}
		if sp.RegistryValue {
			state.damage.result = uint64(projectileDamage(actor, p.temporaryRef(sp.Source), p.temporaryRef(sp.Weapon), sp.Amount, sp.Kind))
		} else {
			state.damage.result = uint64(bool2int(actor.CallDamage(p.temporaryRef(sp.Source), p.temporaryRef(sp.Weapon), int(sp.Amount), object.DamageType(sp.Kind))))
		}
	} else {
		state.damage.result = portTestDamageInvoke(a.Op-1100, p.temporaryRef(attack.Actor), p.temporaryRef(sp.Source), p.temporaryRef(sp.Weapon), p.temporaryRef(sp.Other), sp.Amount, sp.Kind, sp.FloatBits, state.record, (*byte)(unsafe.Pointer(internCStr(sp.Name))), int32(sp.PlayerIndex))
	}
	p.temporary.result = uint32(state.damage.result)
	return p.temporary.result
}

func (p *portTestShopPools) damageSnapshot(out []uint32) []uint32 {
	d := p.temporary.world.objectives.attack.damage
	if d == nil {
		return out
	}
	out = append(out, *(*uint32)(p.temporary.world.objectives.attack.record))
	out = append(out, p.normalize(uint32(d.result)), uint32(d.result>>32), uint32(C.damageCount()))
	for i := 0; i < int(C.damageCount()); i++ {
		out = append(out, p.normalize(uint32(C.damageEvent(C.int(i)))))
	}
	for _, off := range damageOffsets {
		out = append(out, *memmap.PtrUint32(0x5d4594, off))
	}
	return append(out, uint32(dword_5d4594_1563320))
}

func portTestDamageFunction(id int) unsafe.Pointer {
	switch id {
	case 2:
		return damageIdentityKey(damageIDDefault)
	case 8:
		return damageIdentityKey(damageIDBall)
	case 9:
		return damageIdentityKey(damageIDWeapon)
	case 10:
		return damageIdentityKey(damageIDArmor)
	case 14:
		return damageIdentityKey(damageIDPlayer)
	case 20:
		return damageIdentityKey(damageIDSkeleton)
	case 21:
		return damageIdentityKey(damageIDStone)
	case 22:
		return damageIdentityKey(damageIDMechGolem)
	case 23:
		return damageIdentityKey(damageIDFlammable)
	case 24:
		return damageIdentityKey(damageIDBlackPowder)
	case 25:
		return damageIdentityKey(damageIDMonsterGenerator)
	default:
		return nil
	}
}
func portTestDamageInvoke(op int, u, t, it, other *server.Object, amount, kind int32, floatBits uint32, record unsafe.Pointer, name *byte, playerIndex int32) uint64 {
	f := math.Float32frombits(floatBits)
	switch op {
	case 0:
		return uint64(uint32(damageTypeByName(alloc.GoString(name))))
	case 1:
		return uint64(uint32(damageReflect(u, t)))
	case 2:
		return uint64(uint32(damageDefault(u, t, it, amount, kind)))
	case 3:
		damageBall(u, t, amount)
	case 4:
		return uint64(uint32(damageDefend(u, t, it, (*int32)(record), kind)))
	case 5:
		return uint64(uint32(damagePre(u, t, it, (*int32)(record))))
	case 6:
		return uint64(uint32(bool2int(damageMelee(u, it))))
	case 7:
		return uint64(uint32(bool2int(damageFriendlyWeapon(it))))
	case 8:
		return 0
	case 9:
		return uint64(uint32(damageWeapon(u, t, it, amount, kind)))
	case 10:
		return uint64(uint32(damageArmor(u, t, it, amount, kind)))
	case 11:
		damageDurability(u, t, it, other, f, kind, true)
	case 12:
		return uint64(uint32(damageItemReport(playerIndex, t, uint16(amount), uint16(kind))))
	case 13:
		damageDurability(u, t, it, other, f, kind, false)
	case 14:
		return uint64(uint32(damagePlayer(u, t, it, amount, kind)))
	case 15:
		damageFraction(u, (*int32)(record), f)
	case 16:
		damageInventory(u, t, it, amount, f)
	case 17:
		return math.Float64bits(damageConductivity(u))
	case 18:
		return uint64(uint32(damageBlockingItem(u, t, it, amount, f, kind, true)))
	case 19:
		return uint64(uint32(damageBlockingItem(u, t, it, amount, f, kind, false)))
	case 20:
		return uint64(uint32(damageSkeleton(u, t, it, amount, kind)))
	case 21:
		return uint64(uint32(damageDefault(u, t, it, amount, kind)))
	case 22:
		return uint64(uint32(damageMechGolem(u, t, it, amount, kind)))
	case 23:
		return uint64(uint32(damageFlammable(u, t, it, amount, kind)))
	case 24:
		return uint64(uint32(damageBlackPowder(u, t, it, amount, kind)))
	case 25:
		return uint64(uint32(damageGenerator(u, t, it, amount, kind)))
	}
	return 0
}
