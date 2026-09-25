//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
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
		if sp.Controls.Reports != nil {
			extra = append(extra, "TeamBase", "SilverKey", "GoldKey")
		}
		if sp.Controls.SpellLifecycle != nil {
			extra = append(extra, spellLifecycleTypeNames...)
			if sp.Controls.SpellLifecycle.Effects != nil {
				extra = append(extra, spellEffectsTypeNames...)
				if sp.Controls.SpellLifecycle.Effects.Sustained != nil {
					extra = append(extra, sustainedTypeNames...)
				}
			}
		}
	}
	restore := p.proxy.core.PortTestAttackTypes(sp.ProjectileSpeed, sp.MissingTypes, extra...)
	restoreCollision := p.projectileCollisionPrepare()
	restoreDamage := p.damagePrepare()
	restoreState := p.objectStatePrepare()
	restoreReward := p.rewardPrepare()
	restoreControls := p.controlsPrepare()
	oldRange, oldHit, oldTarget := dword_5d4594_2488652, dword_5d4594_2488656, dword_5d4594_2488660
	dword_5d4594_2488652 = uint32(sp.NearestRange)
	dword_5d4594_2488656 = 0
	dword_5d4594_2488660 = 0
	return func() {
		dword_5d4594_2488652, dword_5d4594_2488656, dword_5d4594_2488660 = oldRange, oldHit, oldTarget
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
	p.reservedFunctionIDs += 12
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
			m.AttackPreHit52.Fnc = modifierKey(modifierIDRecoilEffect)
		}
		if sp.PreMask&(1<<i) != 0 {
			m.AttackPreHit52.Fnc = C.attackEffectPtr()
		}
	}
	dword_5d4594_2488660 = uint32(uintptr(p.temporaryRef(sp.NearestTarget).CObj()))
	p.projectileCollisionItems()
	p.damageItems()
	p.objectStateItems()
	p.rewardItems()
	p.controlsItems()
}
func attackCallDirect(op int, u, t, it, ammo *server.Object, record unsafe.Pointer, value uint32) uint32 {
	r := (*attackRecord)(record)
	switch op {
	case 0:
		return uint32(attackPreEffects(t, u, it, r))
	case 1:
		return uint32(attackTrace(u, r))
	case 2:
		attackHit(t, r)
		return 0
	case 3:
		attackNearest(t, u)
		return 0
	case 4:
		return uint32(attackItemEffects(it, u, r))
	case 5:
		return uint32(attackPlayer(u))
	case 6:
		return uint32(int32(attackWarcry(t, u)))
	case 7:
		return uint32(attackBow(u, it))
	case 8:
		return uint32(attackShoot(u, ammo, it, value))
	case 9:
		return uint32(attackShotEffects(u, it, t))
	case 10:
		return uint32(attackReload(equipmentObject(unsafe.Pointer(u.CObj())), 128))
	case 11:
		return uint32(attackReload(equipmentObject(unsafe.Pointer(u.CObj())), 2))
	default:
		return 0
	}
}

func (p *portTestShopPools) attackAction(a PortTestShopAction) uint32 {
	tmp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	sp := tmp.World.Objectives.Attack
	var it *server.Object
	if a.Item >= 0 {
		it = p.items[a.Item].u
	}
	p.temporary.result = attackCallDirect(a.Op-900, p.temporaryRef(sp.Actor), p.temporaryRef(tmp.Target), it, p.temporaryRef(sp.Ammo), p.temporary.world.objectives.attack.record, a.Value)
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
	out = append(out, uint32(dword_5d4594_2488652), uint32(dword_5d4594_2488656), p.normalize(uint32(dword_5d4594_2488660)), uint32(C.attackN()))
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
