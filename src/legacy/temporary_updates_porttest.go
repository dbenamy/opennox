//go:build porttest

package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
void nox_xxx_updateSpark_53ADC0(int a1);
float* nox_xxx_updateProjTrail_53AEC0(int a1);
void nox_xxx_updateLifetime_53B8F0(int unit);
void nox_xxx_spellFlyUpdate_53B940(int a1);
void nox_xxx_updateAntiSpellProj_53BB00(int a1);
void sub_53BD10(int a1, int a2);
int nox_xxx_updateMagicMissile_53BDA0(int a1);
void nox_xxx_updateBlackPowderBarrel_53C9A0(float* a1);
void nox_xxx_updateOneSecondDie_53CB60(int a1);
void nox_xxx_updateWaterBarrel_53CB90(int a1);
void nox_xxx_waterBarrel_53CC30(float* a1, int a2);
void nox_xxx_updateSelfDestruct_53CC90(int a1);
void nox_xxx_updateBlackPowderBurn_53CCB0(int a1);
void nox_xxx_updateDeathBallFragment_53D220(int a1);
void nox_xxx_updateMoonglow_53D270(int a1);
void nox_xxx_updateTelekinesis_53D330(int a1);
void nox_xxx_updateFist_53D400(int a1);
void nox_xxx_updateFlameCleanse_53D510(int a1);
void nox_xxx_updateMeteorShower_53D5A0(float* a2);
void nox_xxx_meteorExplode_53D6E0(int a6);
void nox_xxx_updateToxicCloud_53D850(int a1);
void sub_53D8C0(int a1, int a2);
void nox_xxx_updateSmallToxicCloud_53D960(int a1);
void nox_xxx_toxicCloudPoison_53D9D0(int a1, int a2);
void nox_xxx_updateArachnaphobia_53DA60(int* a1);
void nox_xxx_updateExpire_53DB00(int a1);
int* nox_xxx_updateBreak_53DB30(uint32_t* a1);
int* nox_xxx_updateOpen_53DBB0(uint32_t* a1);
void nox_xxx_updateBreakAndRemove_53DC30(uint32_t* a1);
void nox_xxx_updateChakramInMotion_53DCC0(int a1);
float* nox_xxx_createSpark_54FD80(float a1, float a2, int a3, int a4, float a5, float a6, float a7, int a8);
static uint32_t tempCalls[8192]; static int tempCount; static int tempReturn;
static void tempReset(int ret) { tempCount=0; tempReturn=ret; }
static void tempDie(int u) { tempCalls[tempCount++]=1; tempCalls[tempCount++]=u; }
static int tempCollide(int u,int a,int b) { tempCalls[tempCount++]=2;tempCalls[tempCount++]=u;tempCalls[tempCount++]=a;tempCalls[tempCount++]=b;return tempReturn; }
static void* tempDiePtr(void){return tempDie;} static void* tempCollidePtr(void){return tempCollide;}
static int tempN(void){return tempCount;} static uint32_t tempValue(int i){return tempCalls[i];}
static void* tempFunction(int id) {switch(id){
case 0:return (void*)nox_xxx_updateSpark_53ADC0;
case 1:return (void*)nox_xxx_updateProjTrail_53AEC0;
case 2:return (void*)nox_xxx_updateLifetime_53B8F0;
case 3:return (void*)nox_xxx_spellFlyUpdate_53B940;
case 4:return (void*)nox_xxx_updateAntiSpellProj_53BB00;
case 5:return (void*)sub_53BD10;
case 6:return (void*)nox_xxx_updateMagicMissile_53BDA0;
case 7:return (void*)nox_xxx_updateBlackPowderBarrel_53C9A0;
case 8:return (void*)nox_xxx_updateOneSecondDie_53CB60;
case 9:return (void*)nox_xxx_updateWaterBarrel_53CB90;
case 10:return (void*)nox_xxx_waterBarrel_53CC30;
case 11:return (void*)nox_xxx_updateSelfDestruct_53CC90;
case 12:return (void*)nox_xxx_updateBlackPowderBurn_53CCB0;
case 13:return (void*)nox_xxx_updateDeathBallFragment_53D220;
case 14:return (void*)nox_xxx_updateMoonglow_53D270;
case 15:return (void*)nox_xxx_updateTelekinesis_53D330;
case 16:return (void*)nox_xxx_updateFist_53D400;
case 17:return (void*)nox_xxx_updateFlameCleanse_53D510;
case 18:return (void*)nox_xxx_updateMeteorShower_53D5A0;
case 19:return (void*)nox_xxx_meteorExplode_53D6E0;
case 20:return (void*)nox_xxx_updateToxicCloud_53D850;
case 21:return (void*)sub_53D8C0;
case 22:return (void*)nox_xxx_updateSmallToxicCloud_53D960;
case 23:return (void*)nox_xxx_toxicCloudPoison_53D9D0;
case 24:return (void*)nox_xxx_updateArachnaphobia_53DA60;
case 25:return (void*)nox_xxx_updateExpire_53DB00;
case 26:return (void*)nox_xxx_updateBreak_53DB30;
case 27:return (void*)nox_xxx_updateOpen_53DBB0;
case 28:return (void*)nox_xxx_updateBreakAndRemove_53DC30;
case 29:return (void*)nox_xxx_updateChakramInMotion_53DCC0;
case 30:return (void*)nox_xxx_createSpark_54FD80;
default:return 0;}}
static uint32_t tempCall(int id, nox_object_t* u, nox_object_t* target, int value,int side){switch(id){
case 0: nox_xxx_updateSpark_53ADC0((int)u);return 0;
case 1: return (uint32_t)nox_xxx_updateProjTrail_53AEC0((int)u);
case 2: nox_xxx_updateLifetime_53B8F0((int)u);return 0;
case 3: nox_xxx_spellFlyUpdate_53B940((int)u);return 0;
case 4: nox_xxx_updateAntiSpellProj_53BB00((int)u);return 0;
case 5: sub_53BD10((int)target,(int)u);return 0;
case 6: return (uint32_t)nox_xxx_updateMagicMissile_53BDA0((int)u);
case 7: nox_xxx_updateBlackPowderBarrel_53C9A0((float*)u);return 0;
case 8: nox_xxx_updateOneSecondDie_53CB60((int)u);return 0;
case 9: nox_xxx_updateWaterBarrel_53CB90((int)u);return 0;
case 10: nox_xxx_waterBarrel_53CC30((float*)target,(int)((float*)u+14));return 0;
case 11: nox_xxx_updateSelfDestruct_53CC90((int)u);return 0;
case 12: nox_xxx_updateBlackPowderBurn_53CCB0((int)u);return 0;
case 13: nox_xxx_updateDeathBallFragment_53D220((int)u);return 0;
case 14: nox_xxx_updateMoonglow_53D270((int)u);return 0;
case 15: nox_xxx_updateTelekinesis_53D330((int)u);return 0;
case 16: nox_xxx_updateFist_53D400((int)u);return 0;
case 17: nox_xxx_updateFlameCleanse_53D510((int)u);return 0;
case 18: nox_xxx_updateMeteorShower_53D5A0((float*)u);return 0;
case 19: nox_xxx_meteorExplode_53D6E0((int)u);return 0;
case 20: nox_xxx_updateToxicCloud_53D850((int)u);return 0;
case 21: sub_53D8C0((int)target,(int)u);return 0;
case 22: nox_xxx_updateSmallToxicCloud_53D960((int)u);return 0;
case 23: nox_xxx_toxicCloudPoison_53D9D0((int)target,(int)u);return 0;
case 24: nox_xxx_updateArachnaphobia_53DA60((int*)u);return 0;
case 25: nox_xxx_updateExpire_53DB00((int)u);return 0;
case 26: return (uint32_t)nox_xxx_updateBreak_53DB30((uint32_t*)u);
case 27: return (uint32_t)nox_xxx_updateOpen_53DBB0((uint32_t*)u);
case 28: nox_xxx_updateBreakAndRemove_53DC30((uint32_t*)u);return 0;
case 29: nox_xxx_updateChakramInMotion_53DCC0((int)u);return 0;
case 30: return (uint32_t)nox_xxx_createSpark_54FD80(*((float*)u+14),*((float*)u+15),value,side,*((float*)u+20),*((float*)u+21),*((float*)u+27),((uint32_t*)u)[127]);
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

const (
	PortTestTemporary53ADC0 = 600
	PortTestTemporary53AEC0 = 601
	PortTestTemporary53B8F0 = 602
	PortTestTemporary53B940 = 603
	PortTestTemporary53BB00 = 604
	PortTestTemporary53BD10 = 605
	PortTestTemporary53BDA0 = 606
	PortTestTemporary53C9A0 = 607
	PortTestTemporary53CB60 = 608
	PortTestTemporary53CB90 = 609
	PortTestTemporary53CC30 = 610
	PortTestTemporary53CC90 = 611
	PortTestTemporary53CCB0 = 612
	PortTestTemporary53D220 = 613
	PortTestTemporary53D270 = 614
	PortTestTemporary53D330 = 615
	PortTestTemporary53D400 = 616
	PortTestTemporary53D510 = 617
	PortTestTemporary53D5A0 = 618
	PortTestTemporary53D6E0 = 619
	PortTestTemporary53D850 = 620
	PortTestTemporary53D8C0 = 621
	PortTestTemporary53D960 = 622
	PortTestTemporary53D9D0 = 623
	PortTestTemporary53DA60 = 624
	PortTestTemporary53DB00 = 625
	PortTestTemporary53DB30 = 626
	PortTestTemporary53DBB0 = 627
	PortTestTemporary53DC30 = 628
	PortTestTemporary53DCC0 = 629
	PortTestTemporary54FD80 = 630
)

// References use 0=nil, 1=resource unit, 2=combat target, 3+=inventory items.
var PortTestTemporaryServer func(*server.Server) func()

type PortTestTemporaryUpdatesSpec struct {
	ItemWords, UpdateWords []map[int]uint32
	ItemRefs, UpdateRefs   []map[int]int
	Indexed                []int
	Updatable              int
	DieCallback            bool
	CollideReturn          uint32
	Target                 int
	MissingTypes           []string
}
type portTestTemporaryUpdates struct {
	result  uint32
	indexed []*server.Object
}

var temporaryCacheOffsets = []uintptr{2488664, 2488668, 2488672, 2488684, 2488688, 2488708, 2488636}

func (p *portTestShopPools) temporaryRef(i int) *server.Object {
	switch i {
	case 0:
		return nil
	case 1:
		return p.resources.unit
	case 2:
		return p.proxy.combat.target
	default:
		return p.items[i-3].u
	}
}
func (p *portTestShopPools) temporaryPrepare() func() {
	p.temporary = nil
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	if sp == nil {
		return func() {}
	}
	if p.effectsUse == nil {
		panic("temporary updates require effects fixture")
	}
	p.temporary = &portTestTemporaryUpdates{}
	restore := p.proxy.core.PortTestTemporaryTypes(sp.MissingTypes)
	restoreServer := PortTestTemporaryServer(p.proxy.core)
	old := make([]uint32, len(temporaryCacheOffsets))
	for i, off := range temporaryCacheOffsets {
		old[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	scorch := unsafe.Slice(memmap.PtrUint32(0x587000, 276828), 6)
	oldScorch := append([]uint32(nil), scorch...)
	for i, name := range []string{"temporaryscorch0", "temporaryscorch1", "temporaryscorch2"} {
		scorch[2*i] = uint32(p.proxy.core.Types.IndByID(name))
	}
	*memmap.PtrUint32(0x5d4594, 2488636) = 1
	*memmap.PtrUint32(0x5d4594, 2488672) = 1287568416 // Actual initial nearest-search bound.
	oldUpdatable := p.proxy.core.Objs.UpdatableList
	oldDeleted := p.proxy.core.Objs.DeletedList
	p.proxy.core.Objs.DeletedList = nil
	C.tempReset(C.int(sp.CollideReturn))
	return func() {
		for _, it := range p.items {
			p.temporaryUnindex(it.u)
		}
		for _, u := range p.temporary.indexed {
			p.temporaryUnindex(u)
		}
		p.proxy.core.Objs.UpdatableList = oldUpdatable
		p.proxy.core.Objs.DeletedList = oldDeleted
		for i, off := range temporaryCacheOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		restoreServer()
		copy(scorch, oldScorch)
		restore()
	}
}
func (p *portTestShopPools) temporaryUnindex(u *server.Object) {
	if u == nil || u.ObjFlags&object.FlagPartitioned == 0 {
		return
	}
	for i := 0; i < int(u.ObjIndexCur); i++ {
		p.proxy.core.Map.Sub5178E0(true, &u.ObjIndex[i])
	}
	p.proxy.core.Map.Sub5178E0(false, &u.ObjIndexBase)
	u.ObjFlags &^= object.FlagPartitioned
	u.ObjIndexCur = 0
}
func (p *portTestShopPools) temporaryItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	if sp == nil {
		return
	}
	for i := 0; i < 31; i++ {
		p.identify(C.tempFunction(C.int(i)), 68000+uint32(i))
	}
	p.identify(C.tempDiePtr(), 68100)
	p.identify(C.tempCollidePtr(), 68101)
	// Preserve spatial-node identity across address-space randomization.
	objects := []*server.Object{p.resources.unit, p.proxy.combat.target, p.proxy.combat.actor}
	for _, it := range p.items {
		objects = append(objects, it.u)
	}
	for i, u := range objects {
		p.identify(unsafe.Pointer(&u.ObjIndexBase), 690000+uint32(i*100+90))
		for j := range u.ObjIndex {
			p.identify(unsafe.Pointer(&u.ObjIndex[j]), 690000+uint32(i*100+j))
		}
	}
	// Remove pre-existing index nodes before changing spatial inputs.
	for _, i := range sp.Indexed {
		p.temporaryUnindex(p.temporaryRef(i))
	}
	for i, it := range p.items {
		u := it.u
		apply := func(ptr unsafe.Pointer, words map[int]uint32, refs map[int]int, size int) {
			for off, v := range words {
				if off < 0 || off+4 > size || off%4 != 0 {
					panic("temporary word offset")
				}
				*(*uint32)(unsafe.Add(ptr, off)) = v
			}
			for off, id := range refs {
				if off < 0 || off+4 > size || off%4 != 0 {
					panic("temporary reference offset")
				}
				*(*unsafe.Pointer)(unsafe.Add(ptr, off)) = p.temporaryRef(id).CObj()
			}
		}
		var words map[int]uint32
		var refs map[int]int
		if i < len(sp.ItemWords) {
			words = sp.ItemWords[i]
		}
		if i < len(sp.ItemRefs) {
			refs = sp.ItemRefs[i]
		}
		apply(u.CObj(), words, refs, 772)
		words = nil
		refs = nil
		if i < len(sp.UpdateWords) {
			words = sp.UpdateWords[i]
		}
		if i < len(sp.UpdateRefs) {
			refs = sp.UpdateRefs[i]
		}
		apply(u.UpdateData, words, refs, 64)
		u.Damage = p.proxy.combat.target.Damage
		u.Collide = C.tempCollidePtr()
		if sp.DieCallback {
			u.Death = C.tempDiePtr()
		}
	}
	if sp.Updatable != 0 {
		u := p.temporaryRef(sp.Updatable)
		u.IsUpdatable = 1
		p.proxy.core.Objs.UpdatableList = u
	}
	for _, i := range sp.Indexed {
		u := p.temporaryRef(i)
		u.NewPos = u.PosVec
		u.Shape.Kind = server.ShapeKindCenter
		p.proxy.core.Map.AddObjectToIndex(u)
		p.temporary.indexed = append(p.temporary.indexed, u)
	}
}
func (p *portTestShopPools) temporaryAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	p.temporary.result = uint32(C.tempCall(C.int(a.Op-600), asObjectC(p.items[a.Item].u), asObjectC(p.temporaryRef(sp.Target)), C.int(a.Value), C.int(a.Side)))
	return p.temporary.result
}
func (p *portTestShopPools) temporarySnapshot() []uint32 {
	if p.temporary == nil {
		return nil
	}
	out := []uint32{p.normalize(p.temporary.result), p.normalize(uint32(uintptr(unsafe.Pointer(p.proxy.core.Objs.UpdatableList)))), uint32(C.tempN())}
	for i := 0; i < int(C.tempN()); i++ {
		out = append(out, p.normalize(uint32(C.tempValue(C.int(i)))))
	}
	for _, off := range temporaryCacheOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	for _, it := range p.items {
		for _, v := range unsafe.Slice((*uint32)(it.u.UpdateData), 16) {
			out = append(out, p.normalize(v))
		}
	}
	return out
}
