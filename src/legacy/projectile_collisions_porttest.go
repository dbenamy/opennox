//go:build porttest

package legacy

/*
#include "GAME3_3.h"
#include "GAME4_3.h"
static void* projectileCollisionFunction(int id){switch(id){







case 7: return (void*)sub_4E9A30;












case 20: return (void*)sub_4EB250;
case 21: return (void*)sub_4EB340;
case 22: return (void*)sub_4EB3E0;




default:return 0;}}
static uint32_t projectileCollisionCall(int id,int u,int t,void* normal){switch(id){







case 7: return sub_4E9A30((nox_object_t*)u,(nox_object_t*)t);












case 20: return sub_4EB250(u);
case 21: sub_4EB340((float*)t,u+56);return 0;
case 22: sub_4EB3E0(u);return 0;




default:return 0;}}
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
)

type PortTestProjectileCollisionSpec struct {
	DefinitionRefs      map[int]int
	Normal              *[2]uint32
	WallContact         bool
	WallXY              [2]int32
	CollideRefs         []map[int]int
	Globals             map[int]uint32
	GlobalRefs          map[int]int
	RegisteredCollision bool
}

var projectileCollisionTypeNames = []string{"ThrowingStone", "ImpShot", "ClosedBearTrap", "ToxicCloud"}

var projectileCollisionRegistryNames = map[int]string{
	0: "ProjectileCollide", 1: "ProjectileSparkCollide", 2: "DamageCollide", 3: "ManaDrainCollide",
	4: "BombCollide", 5: "BoomCollide", 6: "DieCollide", 8: "SparkExplosionCollide",
	9: "WallReflectCollide", 10: "YellowStarShotCollide", 11: "DeathBallFragmentCollide", 12: "PixieCollide",
	13: "WallReflectSparkCollide", 14: "OwnCollide", 15: "SparkCollide", 16: "SpiderSpitCollide",
	17: "FistCollide", 18: "TeleportWakeCollide", 19: "ChakramInMotionCollide", 23: "ArrowCollide",
	24: "MonsterArrowCollide", 25: "BearTrapCollide", 26: "PoisonGasTrapCollide",
}

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
	oldTarget, oldContact := dword_5d4594_1567928, dword_5d4594_2488620
	dword_5d4594_1567928 = 0
	dword_5d4594_2488620 = uint32(bool2int(sp.WallContact))
	*memmap.PtrUint32(0x5d4594, 2488612) = uint32(sp.WallXY[0])
	*memmap.PtrUint32(0x5d4594, 2488616) = uint32(sp.WallXY[1])
	return func() {
		for i, off := range projectileCollisionOffsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
		dword_5d4594_1567928 = oldTarget
		dword_5d4594_2488620 = oldContact
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
		p.identify(portTestProjectileCollisionFunction(i), 87000+uint32(i))
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
	if sp.Collision != nil && sp.Collision.RegisteredCollision {
		name := projectileCollisionRegistryNames[a.Op-1000]
		if name != "" {
			legacyNormal := (*types.Pointf)(unsafe.Pointer(normal))
			PortTestRegisteredCollision(p.temporaryRef(sp.Actor), p.temporaryRef(tmp.Target), legacyNormal, name)
			p.temporary.result = 0 // all owner callbacks are invoked for side effects here
			return p.temporary.result
		}
	}
	p.temporary.result = portTestProjectileCollisionCall(int(a.Op-1000), p.temporaryRef(sp.Actor), p.temporaryRef(tmp.Target), (*types.Pointf)(normal))
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
	return append(out, p.normalize(uint32(dword_5d4594_1567928)), uint32(dword_5d4594_2488620))
}

func portTestProjectileCollisionFunction(id int) unsafe.Pointer {
	switch id {
	case 0:
		return collisionKey(collisionIdentityProjectile)
	case 1:
		return collisionKey(collisionIdentityProjectileSpark)
	case 2:
		return collisionKey(collisionIdentityDamage)
	case 3:
		return collisionKey(collisionIdentityManaDrain)
	case 4:
		return collisionKey(collisionIdentityBomb)
	case 5:
		return collisionKey(collisionIdentityBoom)
	case 6:
		return collisionKey(collisionIdentityDie)
	case 8:
		return collisionKey(collisionIdentitySparkExplosion)
	case 9:
		return collisionKey(collisionIdentityWallReflect)
	case 10:
		return collisionKey(collisionIdentityYellowStarShot)
	case 11:
		return collisionKey(collisionIdentityDeathBallFragment)
	case 12:
		return collisionKey(collisionIdentityPixie)
	case 13:
		return collisionKey(collisionIdentityWallReflectSpark)
	case 14:
		return collisionKey(collisionIdentityOwn)
	case 15:
		return collisionKey(collisionIdentitySpark)
	case 16:
		return collisionKey(collisionIdentitySpiderSpit)
	case 17:
		return collisionKey(collisionIdentityFist)
	case 18:
		return collisionKey(collisionIdentityTeleportWake)
	case 19:
		return collisionKey(collisionIdentityChakramInMotion)
	case 23:
		return collisionKey(collisionIdentityArrow)
	case 24:
		return collisionKey(collisionIdentityMonsterArrow)
	case 25:
		return collisionKey(collisionIdentityBearTrap)
	case 26:
		return collisionKey(collisionIdentityPoisonGasTrap)
	default:
		return C.projectileCollisionFunction(C.int(id))
	}
}

func portTestProjectileCollisionCall(id int, u, target *server.Object, normal *types.Pointf) uint32 {
	switch id {
	case 0:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityProjectile), u, target, normal)
	case 1:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityProjectileSpark), u, target, normal)
	case 2:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityDamage), u, target, normal)
	case 3:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityManaDrain), u, target, normal)
	case 4:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityBomb), u, target, normal)
	case 5:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityBoom), u, target, normal)
	case 6:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityDie), u, target, normal)
	case 8:
		return server.PortTestCollisionResult(collisionKey(collisionIdentitySparkExplosion), u, target, normal)
	case 9:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityWallReflect), u, target, normal)
	case 10:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityYellowStarShot), u, target, normal)
	case 11:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityDeathBallFragment), u, target, normal)
	case 12:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityPixie), u, target, normal)
	case 13:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityWallReflectSpark), u, target, normal)
	case 14:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityOwn), u, target, normal)
	case 15:
		return server.PortTestCollisionResult(collisionKey(collisionIdentitySpark), u, target, normal)
	case 16:
		return server.PortTestCollisionResult(collisionKey(collisionIdentitySpiderSpit), u, target, normal)
	case 17:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityFist), u, target, normal)
	case 18:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityTeleportWake), u, target, normal)
	case 19:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityChakramInMotion), u, target, normal)
	case 23:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityArrow), u, target, normal)
	case 24:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityMonsterArrow), u, target, normal)
	case 25:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityBearTrap), u, target, normal)
	case 26:
		return server.PortTestCollisionResult(collisionKey(collisionIdentityPoisonGasTrap), u, target, normal)
	default:
		return uint32(C.projectileCollisionCall(C.int(id), C.int(inventoryInt(u)), C.int(inventoryInt(target)), unsafe.Pointer(normal)))
	}
}
