//go:build porttest

package legacy

/*
#include <stdlib.h>
#include <stdint.h>
extern unsigned int dword_5d4594_2487884,dword_5d4594_2487932,nox_xxx_energyBoltTarget_5d4594_2487880;
extern unsigned int nox_xxx_lightningOwner_5d4594_2487900,nox_xxx_lightningTargetArrayIndex_5d4594_2487904,nox_xxx_lightningTarget_5d4594_2487908,nox_xxx_lightningClosestTargetDistance_5d4594_2487912;
static uint32_t* sustainedGlobals[]={&dword_5d4594_2487884,&dword_5d4594_2487932,&nox_xxx_lightningOwner_5d4594_2487900,&nox_xxx_lightningTargetArrayIndex_5d4594_2487904,&nox_xxx_lightningTarget_5d4594_2487908,&nox_xxx_lightningClosestTargetDistance_5d4594_2487912,&nox_xxx_energyBoltTarget_5d4594_2487880};
static uint32_t sustainedGlobalGet(int i){return *sustainedGlobals[i];}
static void sustainedGlobalSet(int i,uint32_t v){*sustainedGlobals[i]=v;}
#include <stdint.h>
#include <string.h>
#include "GAME4_2.h"
#include "GAME4_3.h"
int nox_xxx_spellDrainMana_52E210(float a1);
int nox_xxx_spellEnergyBoltStop_52E820(int a1);
int nox_xxx_spellEnergyBoltTick_52E850(float a1);
int nox_xxx_firewalkTick_52ED40(float* a1);
int sub_52EF30(int a1);
int sub_52EFD0(int a1);
int sub_52F1D0(int a1);
int sub_52F220(int* a1);
int sub_52F2E0(float a1);
int sub_52F460(float a1);
int nox_xxx_castShield1_52F5A0(uint32_t* a1);
int sub_52F650(int a1);
int sub_52F670(int a1);
int nox_xxx_onStartLightning_52F820(int a1);
int nox_xxx_onFrameLightning_52F8A0(float a1);
char sub_530100(uint32_t* a1);
int nox_xxx_spellTagCreature_530160(uint32_t* a1);
unsigned int sub_530250(int a1);
int sub_530270(int a1);
int nox_xxx_spellBlink2_530310(uint32_t* a1);
int nox_xxx_spellBlink1_530380(int* a1);
int sub_5305D0(uint32_t* a1);
int sub_530650(int* a1);
int nox_xxx_castTele_530820(int a1);
int sub_530880(int* a1);
int nox_xxx_castTTT_530B70(int* a1);
int sub_530CA0(int a1);
int sub_530D30(int* a1);
int nox_xxx_manaBomb_530F90(uint32_t* a1);
int nox_xxx_manaBombBoom_5310C0(int* a1);
int sub_531290(int a1);
int nox_xxx_spellTurnUndeadCreate_531310(uint32_t* a1);
int nox_xxx_spellTurnUndeadUpdate_531410();
int nox_xxx_spellTurnUndeadDelete_531420(int a1);
int sub_531490(uint32_t* a1);
int sub_5314F0(int a1);
int sub_531560(int a1);
int nox_xxx_plasmaSmth_531580(int a1);
int nox_xxx_plasmaShot_531600(int a1);
int sub_5319E0(int a1);
int nox_xxx_spellCreateMoonglow_531A00(uint32_t* a1);
int sub_531AF0(int a1);
static uint32_t sustainedInvoke(int op,int u,int a,void* record,void* output,int q){float fp;memcpy(&fp,&record,4);switch(op){
case 0: return (uint32_t)(nox_xxx_spellDrainMana_52E210(fp));
case 5: return (uint32_t)(nox_xxx_spellEnergyBoltStop_52E820((int)record));
case 6: return (uint32_t)(nox_xxx_spellEnergyBoltTick_52E850(fp));
case 8: return (uint32_t)(nox_xxx_firewalkTick_52ED40((float*)record));
case 9: return (uint32_t)(sub_52EF30((int)record));
case 10: return (uint32_t)(sub_52EFD0((int)record));
case 11: return (uint32_t)(sub_52F1D0((int)record));
case 12: return (uint32_t)(sub_52F220((int*)record));
case 13: return (uint32_t)(sub_52F2E0(fp));
case 14: return (uint32_t)(sub_52F460(fp));
case 15: return (uint32_t)(nox_xxx_castShield1_52F5A0((uint32_t*)record));
case 16: return (uint32_t)(sub_52F650((int)record));
case 17: return (uint32_t)(sub_52F670((int)record));
case 20: return (uint32_t)(nox_xxx_onStartLightning_52F820((int)record));
case 21: return (uint32_t)(nox_xxx_onFrameLightning_52F8A0(fp));
case 24: return (uint32_t)(sub_530100((uint32_t*)record));
case 25: return (uint32_t)(nox_xxx_spellTagCreature_530160((uint32_t*)record));
case 26: return (uint32_t)(sub_530250((int)record));
case 27: return (uint32_t)(sub_530270((int)record));
case 28: return (uint32_t)(nox_xxx_spellBlink2_530310((uint32_t*)record));
case 29: return (uint32_t)(nox_xxx_spellBlink1_530380((int*)record));
case 31: return (uint32_t)(sub_5305D0((uint32_t*)record));
case 32: return (uint32_t)(sub_530650((int*)record));
case 33: return (uint32_t)(nox_xxx_castTele_530820((int)record));
case 34: return (uint32_t)(sub_530880((int*)record));
case 35: return (uint32_t)(nox_xxx_castTTT_530B70((int*)record));
case 36: return (uint32_t)(sub_530CA0((int)record));
case 37: return (uint32_t)(sub_530D30((int*)record));
case 38: return (uint32_t)(nox_xxx_manaBomb_530F90((uint32_t*)record));
case 39: return (uint32_t)(nox_xxx_manaBombBoom_5310C0((int*)record));
case 40: return (uint32_t)(sub_531290((int)record));
case 41: return (uint32_t)(nox_xxx_spellTurnUndeadCreate_531310((uint32_t*)record));
case 42: return (uint32_t)(nox_xxx_spellTurnUndeadUpdate_531410());
case 43: return (uint32_t)(nox_xxx_spellTurnUndeadDelete_531420((int)record));
case 44: return (uint32_t)(sub_531490((uint32_t*)record));
case 45: return (uint32_t)(sub_5314F0((int)record));
case 46: return (uint32_t)(sub_531560((int)record));
case 47: return (uint32_t)(nox_xxx_plasmaSmth_531580((int)record));
case 48: return (uint32_t)(nox_xxx_plasmaShot_531600((int)record));
case 50: return (uint32_t)(sub_5319E0((int)record));
case 51: return (uint32_t)(nox_xxx_spellCreateMoonglow_531A00((uint32_t*)record));
case 52: return (uint32_t)(sub_531AF0((int)record));
default:abort();}}
static void* sustainedFunction(int op){switch(op){
case 0:return nox_xxx_spellDrainMana_52E210;
case 5:return nox_xxx_spellEnergyBoltStop_52E820;
case 6:return nox_xxx_spellEnergyBoltTick_52E850;
case 8:return nox_xxx_firewalkTick_52ED40;
case 9:return sub_52EF30;
case 10:return sub_52EFD0;
case 11:return sub_52F1D0;
case 12:return sub_52F220;
case 13:return sub_52F2E0;
case 14:return sub_52F460;
case 15:return nox_xxx_castShield1_52F5A0;
case 16:return sub_52F650;
case 17:return sub_52F670;
case 20:return nox_xxx_onStartLightning_52F820;
case 21:return nox_xxx_onFrameLightning_52F8A0;
case 24:return sub_530100;
case 25:return nox_xxx_spellTagCreature_530160;
case 26:return sub_530250;
case 27:return sub_530270;
case 28:return nox_xxx_spellBlink2_530310;
case 29:return nox_xxx_spellBlink1_530380;
case 31:return sub_5305D0;
case 32:return sub_530650;
case 33:return nox_xxx_castTele_530820;
case 34:return sub_530880;
case 35:return nox_xxx_castTTT_530B70;
case 36:return sub_530CA0;
case 37:return sub_530D30;
case 38:return nox_xxx_manaBomb_530F90;
case 39:return nox_xxx_manaBombBoom_5310C0;
case 40:return sub_531290;
case 41:return nox_xxx_spellTurnUndeadCreate_531310;
case 42:return nox_xxx_spellTurnUndeadUpdate_531410;
case 43:return nox_xxx_spellTurnUndeadDelete_531420;
case 44:return sub_531490;
case 45:return sub_5314F0;
case 46:return sub_531560;
case 47:return nox_xxx_plasmaSmth_531580;
case 48:return nox_xxx_plasmaShot_531600;
case 50:return sub_5319E0;
case 51:return nox_xxx_spellCreateMoonglow_531A00;
case 52:return sub_531AF0;
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestSustainedSpellsSpec struct {
	RecordDamageRefs []int
	Caches           map[uintptr]uint32
	CacheRefs        map[uintptr]int
	HealTable        [5]float32
	DurationEdges    map[int]map[int]int
	Globals          [7]uint32
	GlobalRefs       map[int]int
}

var sustainedTypeNames = []string{"SmallFlame", "Flame", "ForceOfNatureCharge", "ManaBombCharge", "TeleportWake", "UndeadKiller", "HecubahWithOrb", "Moonglow"}
var sustainedCacheOffsets = []uintptr{2487820, 2487824, 2487828, 2487832, 2487836, 2487840, 2487844, 2487848, 2487852, 2487856, 2487860, 2487868, 2487872, 2487876, 2487888, 2487892, 2487896, 2487916, 2487920, 2487924, 2487928, 2487936, 2487940}

func (p *portTestShopPools) sustainedSpec() *PortTestSustainedSpellsSpec {
	return p.spellLifeSpec().Effects.Sustained
}
func (p *portTestShopPools) sustainedPrepare() func() {
	sp := p.sustainedSpec()
	if sp == nil {
		return func() {}
	}
	// The older AI fixture aliases target player update storage to player zero.
	// Duration effects need distinct caster/donor state, including player metadata.
	target := p.proxy.combat.target
	originalUpdate := target.UpdateData
	if target.ObjClass&4 != 0 && originalUpdate != nil {
		ud := (*server.PlayerUpdateData)(p.objectiveRegion(int(unsafe.Sizeof(server.PlayerUpdateData{}))))
		*ud = *target.UpdateDataPlayer()
		pl := (*server.Player)(p.objectiveRegion(int(unsafe.Sizeof(server.Player{}))))
		*pl = *ud.Player
		pl.PlayerUnit = target
		pl.PlayerInd = 1
		ud.Player = pl
		target.UpdateData = unsafe.Pointer(ud)
	}
	old := make([]uint32, len(sustainedCacheOffsets))
	for i, off := range sustainedCacheOffsets {
		old[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = sp.Caches[off]
	}
	var heal [5]float32
	for i := range heal {
		off := uintptr(260364 + 4*i)
		heal[i] = *memmap.PtrFloat32(0x587000, off)
		*memmap.PtrFloat32(0x587000, off) = sp.HealTable[i]
	}
	var oldGlobals [7]uint32
	for i := range oldGlobals {
		oldGlobals[i] = uint32(C.sustainedGlobalGet(C.int(i)))
		C.sustainedGlobalSet(C.int(i), C.uint32_t(sp.Globals[i]))
	}
	oldPlasma := *memmap.PtrUint32(0x587000, 260404)
	*memmap.PtrUint32(0x587000, 260404) = 1209810944
	return func() {
		target.UpdateData = originalUpdate
		for i, off := range sustainedCacheOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		for i := range heal {
			*memmap.PtrFloat32(0x587000, uintptr(260364+4*i)) = heal[i]
		}
		for i, v := range oldGlobals {
			C.sustainedGlobalSet(C.int(i), C.uint32_t(v))
		}
		*memmap.PtrUint32(0x587000, 260404) = oldPlasma
	}
}
func (p *portTestShopPools) sustainedItems() {
	sp := p.sustainedSpec()
	if sp == nil {
		return
	}
	for _, ref := range sp.RecordDamageRefs {
		p.temporaryRef(ref).Damage = p.proxy.combat.target.Damage
	}
	for i := 0; i < 53; i++ {
		if fn := C.sustainedFunction(C.int(i)); fn != nil {
			p.identify(fn, 97000+uint32(i))
		}
	}
	for i, ref := range sp.GlobalRefs {
		C.sustainedGlobalSet(C.int(i), C.uint32_t(uintptr(p.temporaryRef(ref).CObj())))
	}
	for off, ref := range sp.CacheRefs {
		*memmap.PtrPtr(0x5d4594, off) = p.temporaryRef(ref).CObj()
	}
	for i, edges := range sp.DurationEdges {
		for off, j := range edges {
			*controlPtr(p.spellLifeState().durations[i].C(), off) = p.spellLifeState().durations[j].C()
		}
	}
}
func (p *portTestShopPools) sustainedAction(a PortTestShopAction) uint32 {
	sp := p.spellLifeSpec()
	st := p.spellLifeState()
	ctrl := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	record, output := st.record, st.effects.output
	if sp.NullRecord {
		record = nil
	}
	if sp.Effects.NullOutput {
		output = nil
	}
	u, t := p.temporaryRef(ctrl.Actor), p.temporaryRef(sp.Effects.Args[0])
	q := int32(sp.Effects.Ints[0])
	var out uint32
	switch a.Op - 1700 {
	case 1:
		out = sustainedTransferMana(u, t, q)
	case 2:
		out = controlRaw(sustainedFindMana(*(*types.Pointf)(record), u))
	case 3:
		sustainedManaCandidate(u, t)
	case 4:
		out = sustainedHasMana(u)
	case 7:
		sustainedEnergyCandidate(u, t)
	case 18:
		sustainedShieldAbsorb(u, q)
	case 19:
		sustainedShieldDamage(u, (*int32)(output), q, t)
	case 22:
		sustainedLightningCandidate(u, t)
	case 23:
		sustainedLightningTrapHit(u, t)
	case 30:
		out = sustainedTeleportWake(u, (*types.Pointf)(record), (*types.Pointf)(output))
	case 49:
		sustainedPlasmaCandidate(u, t)
	default:
		out = uint32(C.sustainedInvoke(C.int(a.Op-1700), C.int(uintptr(u.CObj())), C.int(uintptr(t.CObj())), record, output, C.int(q)))
	}
	p.sustainedChildren(true)
	ctrlState := p.temporary.world.objectives.attack.controls
	ctrlState.result = uint64(out)
	return p.normalize(out)
}
func (p *portTestShopPools) sustainedSnapshot(out []uint32) []uint32 {
	if p.sustainedSpec() == nil {
		return out
	}
	for _, off := range sustainedCacheOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	for i := 0; i < 7; i++ {
		out = append(out, p.normalize(uint32(C.sustainedGlobalGet(C.int(i)))))
	}
	out = append(out, *memmap.PtrUint32(0x587000, 260404))
	children := p.sustainedChildren(false)
	out = append(out, uint32(len(children)))
	for _, ptr := range children {
		for _, v := range unsafe.Slice((*uint32)(ptr), 30) {
			out = append(out, p.normalize(v))
		}
	}
	return out
}

func (p *portTestShopPools) sustainedChildren(identify bool) []unsafe.Pointer {
	var out []unsafe.Pointer
	seen := map[unsafe.Pointer]bool{}
	for _, off := range []int{104, 108} {
		for ptr := *controlPtr(p.spellLifeState().record, off); ptr != nil; ptr = *controlPtr(ptr, 116) {
			if seen[ptr] {
				break
			}
			seen[ptr] = true
			out = append(out, ptr)
		}
	}
	if identify {
		for _, ptr := range out {
			// Allocator addresses can reuse an identity from an earlier case. The real
			// duration allocation ID remains stable across ticks and independent of ASLR.
			p.identify(ptr, 98000+uint32((*server.DurSpell)(ptr).ID))
		}
	}
	return out
}

type PortTestSustainedSpellsResult struct {
	Record, Output        []uint32
	Children, CreatedInit [][]uint32
	Globals               []uint32
}

func (p *portTestShopPools) sustainedDetails() *PortTestSustainedSpellsResult {
	words := func(ptr unsafe.Pointer, n int) []uint32 {
		out := make([]uint32, n)
		for i, v := range unsafe.Slice((*uint32)(ptr), n) {
			out[i] = p.normalize(v)
		}
		return out
	}
	st := p.spellLifeState()
	out := &PortTestSustainedSpellsResult{Record: words(st.record, 40), Output: words(st.effects.output, 8)}
	for _, ptr := range p.sustainedChildren(false) {
		out.Children = append(out.Children, words(ptr, 30))
	}
	for _, u := range p.proxy.life.created {
		if u.InitData != nil {
			for _, v := range unsafe.Slice((*byte)(u.InitData), 36)[20:] {
				if v != 0x5a {
					panic("sustained init guard")
				}
			}
			out.CreatedInit = append(out.CreatedInit, words(u.InitData, 9))
		}
	}
	for i := 0; i < 7; i++ {
		out.Globals = append(out.Globals, p.normalize(uint32(C.sustainedGlobalGet(C.int(i)))))
	}
	return out
}
