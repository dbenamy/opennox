//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME4_3.h"
#include "GAME3_2.h"
extern uint32_t dword_5d4594_2488652,dword_5d4594_2488656,dword_5d4594_2488660;
static uint32_t attackTrace[8192];static int attackCount,attackSet;static uint32_t attackOutput;static uint32_t attackFrontMask;
static void attackReset(int set,uint32_t output){attackFrontMask=255;attackCount=0;attackSet=set;attackOutput=output;}
static int attackEffect(int m,int it,int target,int actor,int record){
 attackTrace[attackCount++]=m;attackTrace[attackCount++]=it;attackTrace[attackCount++]=target;attackTrace[attackCount++]=actor;
 // Bytes 5..7 and 33..35 are unused padding in the original stack record.
 memcpy(&attackTrace[attackCount],(void*)record,36);
 attackTrace[attackCount+1]&=255;attackTrace[attackCount+8]&=attackFrontMask;attackCount+=9;
 if(attackSet) *(uint32_t*)record=attackOutput;
 return 0;
}
// Collision owners leave the entire Front word unused and uninitialized.
static void attackCollisionRecord(void){attackFrontMask=0;}
static void* attackEffectPtr(void){return attackEffect;}
static int attackN(void){return attackCount;}
static uint32_t attackValue(int i){return attackTrace[i];}
static void* attackFunction(int id){switch(id){
case 0:return (void*)nox_xxx_playerPreAttackEffects_538290;
case 1:return (void*)nox_xxx_playerTraceAttack_538330;
case 2:return (void*)sub_538510;
case 3:return (void*)sub_5386A0;
case 4:return (void*)nox_xxx_itemApplyAttackEffect_538840;
case 5:return (void*)nox_xxx_playerAttack_538960;
case 6:return (void*)nox_xxx_warcryStunMonsters_539B90;
case 7:return (void*)nox_xxx_shootBowCrossbow1_539BD0;
case 8:return (void*)nox_xxx_shootBowCrossbow2_539D80;
case 9:return (void*)nox_xxx_shootApplyEffects_539F40;
case 10:return (void*)sub_539FB0;
case 11:return (void*)nox_xxx_playerTryReloadQuiver_539FF0;
default:return 0;}}
static uint32_t attackCall(int id,nox_object_t* u,nox_object_t* t,nox_object_t* it,nox_object_t* ammo,void* record,uint32_t value){switch(id){
case 0:return nox_xxx_playerPreAttackEffects_538290((int)t,(int)u,(int)it,(int)record);
case 1:return nox_xxx_playerTraceAttack_538330((int)u,(int)record);
case 2:sub_538510((int)t,(int)record);return 0;
case 3:sub_5386A0((int)t,(int)u);return 0;
case 4:return nox_xxx_itemApplyAttackEffect_538840((int)it,(int)u,(int)record);
case 5:return nox_xxx_playerAttack_538960(u);
case 6:return (uint32_t)nox_xxx_warcryStunMonsters_539B90((int)t,(int)u);
case 7:return nox_xxx_shootBowCrossbow1_539BD0((int)u,(int)it);
case 8:return (uint32_t)nox_xxx_shootBowCrossbow2_539D80((int)u,(int)ammo,(int)it,(char*)value);
case 9:return nox_xxx_shootApplyEffects_539F40((int)u,(int)it,(int)t);
case 10:return sub_539FB0((uint32_t*)u);
case 11:return nox_xxx_playerTryReloadQuiver_539FF0((uint32_t*)u);
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

const (
	PortTestAttack538290 = 900
	PortTestAttack538330 = 901
	PortTestAttack538510 = 902
	PortTestAttack5386A0 = 903
	PortTestAttack538840 = 904
	PortTestAttack538960 = 905
	PortTestAttack539B90 = 906
	PortTestAttack539BD0 = 907
	PortTestAttack539D80 = 908
	PortTestAttack539F40 = 909
	PortTestAttack539FB0 = 910
	PortTestAttack539FF0 = 911
)

type PortTestAttackSpec struct {
	Controls                             *PortTestPlayerControlsSpec
	Reward                               *PortTestRewardSpec
	State                                *PortTestObjectStateSpec
	Damage                               *PortTestDamageSpec
	Collision                            *PortTestProjectileCollisionSpec
	Actor, Ammo                          int
	RecordWords, ActorWords, UpdateWords map[int]uint32
	RecordRefs, ActorRefs, UpdateRefs    map[int]int
	AttackMask, PreMask, RecoilMask      uint32
	SetOutput                            bool
	Output, NearestRange                 uint32
	NearestTarget                        int
	ProjectileSpeed                      float32
	MissingTypes                         []string
}
type portTestAttack struct {
	controls                *portTestPlayerControls
	reward                  *portTestReward
	record, collisionNormal unsafe.Pointer
	damage                  *portTestDamage
	state                   *portTestObjectState
}

func (p *portTestShopPools) attackPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	if sp == nil {
		return func() {}
	}
	recordSize := 64
	if sp.Controls != nil {
		recordSize = 80
	}
	p.temporary.world.objectives.attack = &portTestAttack{record: p.objectiveRegion(recordSize)}
	var extra []string
	if sp.Collision != nil {
		extra = projectileCollisionTypeNames
	}
	if sp.Damage != nil {
		extra = append(extra, damageTypeNames...)
	}
	if sp.State != nil {
		extra = append(extra, objectStateTypeNames()...)
	}
	if sp.Controls != nil {
		extra = append(extra, controlsTypeNames...)
	}
	restore := p.proxy.core.PortTestAttackTypes(sp.ProjectileSpeed, sp.MissingTypes, extra...)
	restoreCollision := p.projectileCollisionPrepare()
	restoreDamage := p.damagePrepare()
	restoreState := p.objectStatePrepare()
	restoreReward := p.rewardPrepare()
	restoreControls := p.controlsPrepare()
	oldRange, oldHit, oldTarget := C.dword_5d4594_2488652, C.dword_5d4594_2488656, C.dword_5d4594_2488660
	C.dword_5d4594_2488652 = C.uint32_t(sp.NearestRange)
	C.dword_5d4594_2488656 = 0
	C.dword_5d4594_2488660 = 0
	return func() {
		C.dword_5d4594_2488652, C.dword_5d4594_2488656, C.dword_5d4594_2488660 = oldRange, oldHit, oldTarget
		restoreControls()
		restoreReward()
		restoreState()
		restoreDamage()
		restoreCollision()
		restore()
	}
}
func (p *portTestShopPools) attackItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack
	if sp == nil {
		return
	}
	r := p.temporary.world.objectives.attack
	C.attackReset(C.int(bool2int(sp.SetOutput)), C.uint32_t(sp.Output))
	if sp.Collision != nil {
		C.attackCollisionRecord()
	}
	p.identify(C.attackEffectPtr(), 86000)
	for i := 0; i < 12; i++ {
		p.identify(C.attackFunction(C.int(i)), 86001+uint32(i))
	}
	apply := func(ptr unsafe.Pointer, size int, words map[int]uint32, refs map[int]int) {
		for off, v := range words {
			if off < 0 || off+4 > size || off%4 != 0 {
				panic("attack word offset")
			}
			*(*uint32)(unsafe.Add(ptr, off)) = v
		}
		for off, id := range refs {
			if off < 0 || off+4 > size || off%4 != 0 {
				panic("attack ref offset")
			}
			*(*unsafe.Pointer)(unsafe.Add(ptr, off)) = p.temporaryRef(id).CObj()
		}
	}
	recordSize := 64
	if sp.Controls != nil {
		recordSize = 80
	}
	apply(r.record, recordSize, sp.RecordWords, sp.RecordRefs)
	u := p.temporaryRef(sp.Actor)
	if u != nil {
		apply(u.CObj(), 772, sp.ActorWords, sp.ActorRefs)
		size := 64
		if u.ObjClass&4 != 0 {
			size = int(unsafe.Sizeof(server.PlayerUpdateData{}))
		} else if u.ObjClass&2 != 0 {
			size = 2200
		}
		apply(u.UpdateData, size, sp.UpdateWords, sp.UpdateRefs)
	}
	for i, it := range p.items {
		*equipmentWord(it.u.InitData, 4) = p.temporary.world.objectives.initWords[i]
		it.u.Damage = p.proxy.combat.target.Damage
	}
	for i := 0; i < 4; i++ {
		m := (*server.ModifierEff)(unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*144))
		if sp.AttackMask&(1<<i) != 0 {
			m.Attack40.Fnc = C.attackEffectPtr()
		}
		if sp.RecoilMask&(1<<i) != 0 {
			m.AttackPreHit52.Fnc = C.nox_xxx_recoilEffect_4E0640
		}
		if sp.PreMask&(1<<i) != 0 {
			m.AttackPreHit52.Fnc = C.attackEffectPtr()
		}
	}
	C.dword_5d4594_2488660 = C.uint32_t(uintptr(p.temporaryRef(sp.NearestTarget).CObj()))
	p.projectileCollisionItems()
	p.damageItems()
	p.objectStateItems()
	p.rewardItems()
	p.controlsItems()
}
func (p *portTestShopPools) attackAction(a PortTestShopAction) uint32 {
	tmp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	sp := tmp.World.Objectives.Attack
	var it *server.Object
	if a.Item >= 0 {
		it = p.items[a.Item].u
	}
	p.temporary.result = uint32(C.attackCall(C.int(a.Op-900), asObjectC(p.temporaryRef(sp.Actor)), asObjectC(p.temporaryRef(tmp.Target)), asObjectC(it), asObjectC(p.temporaryRef(sp.Ammo)), p.temporary.world.objectives.attack.record, C.uint32_t(a.Value)))
	if a.Op == PortTestAttack539B90 {
		// This retained short is an input object's truncated address, not a scalar.
		// Normalize known low-word addresses just as the fixture normalizes full pointers.
		for _, ref := range []int{tmp.Target, sp.Actor} {
			obj := p.temporaryRef(ref)
			if obj != nil && p.temporary.result == uint32(int32(int16(uintptr(obj.CObj())))) {
				p.temporary.result = p.normalize(uint32(uintptr(obj.CObj())))
				break
			}
		}
	}
	return p.temporary.result
}
func (p *portTestShopPools) attackSnapshot(out []uint32) []uint32 {
	if p.temporary.world.objectives.attack == nil {
		return out
	}
	out = append(out, uint32(C.dword_5d4594_2488652), uint32(C.dword_5d4594_2488656), p.normalize(uint32(C.dword_5d4594_2488660)), uint32(C.attackN()))
	for i := 0; i < int(C.attackN()); i++ {
		out = append(out, p.normalize(uint32(C.attackValue(C.int(i)))))
	}
	for _, u := range p.proxy.life.created {
		if u.InitData != nil {
			size := 0
			if typ := p.proxy.core.Types.ByInd(int(u.TypeInd)); typ != nil {
				size = int(typ.InitDataSize)
			}
			if size == 36 {
				for _, v := range unsafe.Slice((*byte)(u.InitData), size)[20:] {
					if v != 0x5a {
						panic("attack projectile init guard")
					}
				}
			}
			for _, v := range unsafe.Slice((*uint32)(u.InitData), size/4) {
				out = append(out, p.normalize(v))
			}
		}
	}
	return p.controlsSnapshot(p.rewardSnapshot(p.objectStateSnapshot(p.damageSnapshot(p.projectileCollisionSnapshot(out)))))
}
