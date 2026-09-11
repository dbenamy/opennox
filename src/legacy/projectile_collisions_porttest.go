//go:build porttest

package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_1567928, dword_5d4594_2488620;
static void* projectileCollisionFunction(int id){switch(id){
case 0: return (void*)nox_xxx_collideProjectileGeneric_4E87B0;
case 1: return (void*)nox_xxx_collideProjectileSpark_4E8880;
case 2: return (void*)nox_xxx_collideDamage_4E9430;
case 3: return (void*)nox_xxx_collideManadrain_4E9490;
case 4: return (void*)nox_xxx_collideBomb_4E96F0;
case 5: return (void*)nox_xxx_collideBoom_4E9770;
case 6: return (void*)nox_xxx_collideDie_4E99B0;
case 7: return (void*)sub_4E9A30;
case 8: return (void*)nox_xxx_fireballCollide_4E9AC0;
case 9: return (void*)nox_xxx_collideSulphurShot2_4E9D80;
case 10: return (void*)nox_xxx_collideSulphurShot_4E9E50;
case 11: return (void*)nox_xxx_collideDeathBallFragment_4E9FE0;
case 12: return (void*)nox_xxx_collidePixie_4EA080;
case 13: return (void*)nox_xxx_collideWallReflectSpark_4EA200;
case 14: return (void*)sub_4EA2C0;
case 15: return (void*)nox_xxx_collideSpark_4EA300;
case 16: return (void*)nox_xxx_collideWebbing_4EA380;
case 17: return (void*)nox_xxx_collideFist_4EADF0;
case 18: return (void*)nox_xxx_collideTeleportWake_4EAE30;
case 19: return (void*)nox_xxx_collideChakram_4EAF00;
case 20: return (void*)sub_4EB250;
case 21: return (void*)sub_4EB340;
case 22: return (void*)sub_4EB3E0;
case 23: return (void*)nox_xxx_collideArrow_4EB490;
case 24: return (void*)nox_xxx_collideMonsterArrow_4EB800;
case 25: return (void*)nox_xxx_collideBearTrap_4EB890;
case 26: return (void*)nox_xxx_collidePoisonGasTrap_4EB910;
default:return 0;}}
static uint32_t projectileCollisionCall(int id,int u,int t,void* normal){switch(id){
case 0: nox_xxx_collideProjectileGeneric_4E87B0(u,t);return 0;
case 1: nox_xxx_collideProjectileSpark_4E8880(u,t);return 0;
case 2: nox_xxx_collideDamage_4E9430(u,t);return 0;
case 3: nox_xxx_collideManadrain_4E9490(u,t);return 0;
case 4: nox_xxx_collideBomb_4E96F0(u,t);return 0;
case 5: nox_xxx_collideBoom_4E9770(u,t,(float*)normal);return 0;
case 6: nox_xxx_collideDie_4E99B0(u,t);return 0;
case 7: return sub_4E9A30((nox_object_t*)u,(nox_object_t*)t);
case 8: nox_xxx_fireballCollide_4E9AC0(u,t);return 0;
case 9: nox_xxx_collideSulphurShot2_4E9D80(u,t,(float*)normal);return 0;
case 10: nox_xxx_collideSulphurShot_4E9E50(u,t,(int)normal);return 0;
case 11: nox_xxx_collideDeathBallFragment_4E9FE0(u,t,(float*)normal);return 0;
case 12: nox_xxx_collidePixie_4EA080(u,t,(float*)normal);return 0;
case 13: nox_xxx_collideWallReflectSpark_4EA200(u,t,(float2*)normal);return 0;
case 14: sub_4EA2C0(u,t);return 0;
case 15: nox_xxx_collideSpark_4EA300(u,t,(float*)normal);return 0;
case 16: nox_xxx_collideWebbing_4EA380(u,t);return 0;
case 17: nox_xxx_collideFist_4EADF0(u,t);return 0;
case 18: nox_xxx_collideTeleportWake_4EAE30(u,t);return 0;
case 19: nox_xxx_collideChakram_4EAF00(u,t,(float*)normal);return 0;
case 20: return sub_4EB250(u);
case 21: sub_4EB340((float*)t,u+56);return 0;
case 22: sub_4EB3E0(u);return 0;
case 23: nox_xxx_collideArrow_4EB490(u,t);return 0;
case 24: nox_xxx_collideMonsterArrow_4EB800(u,t);return 0;
case 25: nox_xxx_collideBearTrap_4EB890((int*)u,t);return 0;
case 26: nox_xxx_collidePoisonGasTrap_4EB910((int*)u,t);return 0;
default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

type PortTestProjectileCollisionSpec struct {
	DefinitionRefs map[int]int
	Normal         *[2]uint32
	WallContact    bool
	WallXY         [2]int32
	CollideRefs    []map[int]int
	Globals        map[int]uint32
	GlobalRefs     map[int]int
}

var projectileCollisionTypeNames = []string{"ThrowingStone", "ImpShot", "ClosedBearTrap", "ToxicCloud"}

var projectileCollisionOffsets = []uintptr{1567836, 1567840, 1567924, 1567932, 1567948, 1567952, 1567964, 1567968, 1567972, 1567976, 1567980, 1567984, 1568000, 2488612, 2488616}

func (p *portTestShopPools) projectileCollisionPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Collision
	if sp == nil {
		return func() {}
	}
	if sp.Normal != nil {
		ptr := p.objectiveRegion(8)
		copy(unsafe.Slice((*uint32)(ptr), 2), sp.Normal[:])
		p.temporary.world.objectives.attack.collisionNormal = ptr
	}
	var old []uint32
	for _, off := range projectileCollisionOffsets {
		old = append(old, *memmap.PtrUint32(0x5d4594, off))
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	oldTarget, oldContact := C.dword_5d4594_1567928, C.dword_5d4594_2488620
	C.dword_5d4594_1567928 = 0
	C.dword_5d4594_2488620 = C.uint32_t(bool2int(sp.WallContact))
	*memmap.PtrUint32(0x5d4594, 2488612) = uint32(sp.WallXY[0])
	*memmap.PtrUint32(0x5d4594, 2488616) = uint32(sp.WallXY[1])
	return func() {
		for i, off := range projectileCollisionOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		C.dword_5d4594_1567928 = oldTarget
		C.dword_5d4594_2488620 = oldContact
	}
}

func (p *portTestShopPools) projectileCollisionItems() {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Collision
	if sp == nil {
		return
	}
	for typ, id := range sp.DefinitionRefs {
		d := p.proxy.core.Modif.Nox_xxx_getProjectileClassById413250(typ)
		if d == nil {
			panic("missing projectile modifier definition")
		}
		d.TypeInd = uint32(p.temporaryRef(id).TypeInd)
	}
	for i := 0; i < 27; i++ {
		p.identify(C.projectileCollisionFunction(C.int(i)), 87000+uint32(i))
	}
	for _, id := range []int{1, 2, 100, 101, 102} {
		if u := p.temporaryRef(id); u != nil {
			u.Damage = p.proxy.combat.target.Damage
		}
	}
	for i, refs := range sp.CollideRefs {
		for off, id := range refs {
			if off < 0 || off+4 > 64 || off%4 != 0 {
				panic("projectile collide ref offset")
			}
			*temporaryRefWord(p.items[i].u.CollideData, off) = p.temporaryRef(id)
		}
	}
	valid := func(off int) {
		for _, x := range projectileCollisionOffsets {
			if x == uintptr(off) {
				return
			}
		}
		panic("projectile global offset")
	}
	for off, v := range sp.Globals {
		valid(off)
		*memmap.PtrUint32(0x5d4594, uintptr(off)) = v
	}
	for off, id := range sp.GlobalRefs {
		valid(off)
		*memmap.PtrUint32(0x5d4594, uintptr(off)) = uint32(uintptr(p.temporaryRef(id).CObj()))
	}
}

func (p *portTestShopPools) projectileCollisionAction(a PortTestShopAction) uint32 {
	tmp := p.proxy.callbacks.shop.spec.TemporaryUpdates
	sp := tmp.World.Objectives.Attack
	normal := p.temporary.world.objectives.attack.collisionNormal
	p.temporary.result = uint32(C.projectileCollisionCall(C.int(a.Op-1000), inventoryInt(p.temporaryRef(sp.Actor)), inventoryInt(p.temporaryRef(tmp.Target)), normal))
	return p.temporary.result
}

func (p *portTestShopPools) projectileCollisionSnapshot(out []uint32) []uint32 {
	if p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Collision == nil {
		return out
	}
	for _, off := range projectileCollisionOffsets {
		out = append(out, p.normalize(*memmap.PtrUint32(0x5d4594, off)))
	}
	if n := p.temporary.world.objectives.attack.collisionNormal; n != nil {
		out = append(out, unsafe.Slice((*uint32)(n), 2)...)
	}
	return append(out, p.normalize(uint32(C.dword_5d4594_1567928)), uint32(C.dword_5d4594_2488620))
}
