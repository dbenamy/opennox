//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME3_2.h"
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1563320;
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
static void* damageFunction(int id){switch(id){
case 0:return (void*)nox_xxx_parseDamageTypeByName_4E0A00;
case 1:return (void*)nox_xxx_projectileReflect_4E0A70;
case 2:return (void*)nox_xxx_damageDefaultProc_4E0B30;
case 3:return (void*)nox_xxx_gameballOnPlayerDamage_4E1230;
case 4:return (void*)nox_xxx_itemApplyDefendEffect2_4E1320;
case 5:return (void*)nox_xxx_itemApplyPreDamageEffect_4E13B0;
case 6:return (void*)sub_4E1400;
case 7:return (void*)sub_4E1470;
case 8:return (void*)sub_4E14A0;
case 9:return (void*)sub_4E14B0;
case 10:return (void*)nox_xxx_damageArmor_4E1500;
case 11:return (void*)nox_xxx_playerDamageWeapon_4E1560;
case 12:return (void*)nox_xxx_itemDestroyed_4E1650;
case 13:return (void*)nox_xxx_equipDamage_4E16D0;
case 14:return (void*)nox_server_handler_PlayerDamage_4E17B0;
case 15:return (void*)nox_xxx_playerDecrementHPMana_4E20F0;
case 16:return (void*)nox_xxx_playerDamageItems_4E2180;
case 17:return (void*)sub_4E2220;
case 18:return (void*)sub_4E22A0;
case 19:return (void*)sub_4E2330;
case 20:return (void*)sub_4E23C0;
case 21:return (void*)sub_4E24B0;
case 22:return (void*)sub_4E24E0;
case 23:return (void*)nox_xxx_damageFlammable_4E2520;
case 24:return (void*)nox_xxx_damageBlackPowder_4E2560;
case 25:return (void*)nox_xxx_damageMonsterGen_4E27D0;
default:return 0;}}
static uint64_t damageCall(int id,int u,int t,int it,int other,int damage,int kind,uint32_t floatBits,void* record,char* name,int playerIndex){
 float f;memcpy(&f,&floatBits,4);switch(id){
case 0:{return (uint32_t)nox_xxx_parseDamageTypeByName_4E0A00(name);}
case 1:{return (uint32_t)nox_xxx_projectileReflect_4E0A70(u,t);}
case 2:{return (uint32_t)nox_xxx_damageDefaultProc_4E0B30(u,t,it,damage,kind);}
case 3:{nox_xxx_gameballOnPlayerDamage_4E1230(u,t,damage);return 0;}
case 4:{return (uint32_t)nox_xxx_itemApplyDefendEffect2_4E1320(u,t,it,(int*)record,kind);}
case 5:{return (uint32_t)nox_xxx_itemApplyPreDamageEffect_4E13B0(u,t,it,(int)record);}
case 6:{return (uint32_t)sub_4E1400(u,(uint32_t*)it);}
case 7:{return (uint32_t)sub_4E1470(it);}
case 8:{return (uint32_t)sub_4E14A0();}
case 9:{return (uint32_t)sub_4E14B0(u,t,it,damage,kind);}
case 10:{return (uint32_t)nox_xxx_damageArmor_4E1500(u,t,it,damage,kind);}
case 11:{nox_xxx_playerDamageWeapon_4E1560(u,t,it,other,f,kind);return 0;}
case 12:{return (uint32_t)nox_xxx_itemDestroyed_4E1650(playerIndex,(uint32_t*)t,(unsigned short)damage,(unsigned short)kind);}
case 13:{nox_xxx_equipDamage_4E16D0(u,t,it,other,f,kind);return 0;}
case 14:{return (uint32_t)nox_server_handler_PlayerDamage_4E17B0(u,t,it,damage,kind);}
case 15:{nox_xxx_playerDecrementHPMana_4E20F0(u,(int)record,f);return 0;}
case 16:{nox_xxx_playerDamageItems_4E2180(u,t,it,damage,f);return 0;}
case 17:{double d=sub_4E2220(u);uint64_t bits;memcpy(&bits,&d,8);return bits;}
case 18:{return (uint32_t)sub_4E22A0(u,t,it,damage,f,kind);}
case 19:{return (uint32_t)sub_4E2330(u,t,it,damage,f,kind);}
case 20:{return (uint32_t)sub_4E23C0(u,t,it,damage,kind);}
case 21:{return (uint32_t)sub_4E24B0(u,t,it,damage,kind);}
case 22:{return (uint32_t)sub_4E24E0(u,t,it,damage,kind);}
case 23:{return (uint32_t)nox_xxx_damageFlammable_4E2520(u,t,it,damage,kind);}
case 24:{return (uint32_t)nox_xxx_damageBlackPowder_4E2560(u,t,it,damage,kind);}
case 25:{return (uint32_t)nox_xxx_damageMonsterGen_4E27D0(u,t,it,damage,kind);}
default:return 0;}}
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestDamageSpec struct {
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
	cached := C.dword_5d4594_1563320
	C.dword_5d4594_1563320 = 0
	return func() {
		for i, old := range oldNames {
			*memmap.PtrPtr(0x587000, 200728+uintptr(i*4)) = old
		}
		for i, off := range damageOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		C.dword_5d4594_1563320 = cached
		copy(raw, oldRaw)
	}
}

func (p *portTestShopPools) damageItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Damage
	if sp == nil {
		return
	}
	C.damageReset(C.int(bool2int(sp.SetOutput)), C.uint32_t(sp.Output))
	for i := 0; i < 26; i++ {
		p.identify(C.damageFunction(C.int(i)), 88000+uint32(i))
	}
	p.identify(C.damageDefendPtr(), 88050)
	p.identify(C.damagePrePtr(), 88051)
	p.identify(C.damageDefendScalarPtr(), 88052)
	for _, id := range []int{1, 2, 100, 101, 102} {
		if u := p.temporaryRef(id); u != nil {
			u.Damage = p.proxy.combat.target.Damage
		}
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
	state.damage.result = uint64(C.damageCall(C.int(a.Op-1100), inventoryInt(p.temporaryRef(attack.Actor)), inventoryInt(p.temporaryRef(sp.Source)), inventoryInt(p.temporaryRef(sp.Weapon)), inventoryInt(p.temporaryRef(sp.Other)), C.int(sp.Amount), C.int(sp.Kind), C.uint32_t(sp.FloatBits), state.record, internCStr(sp.Name), C.int(sp.PlayerIndex)))
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
	return append(out, uint32(C.dword_5d4594_1563320))
}
