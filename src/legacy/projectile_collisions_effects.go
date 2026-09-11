package legacy

/*
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func projectileFriendly(u, t *server.Object) bool {
	if !noxflags.HasGame(4096) {
		return false
	}
	owner := u.FindOwnerChainPlayer()
	return owner != nil && t != nil && owner.ObjClass&4 != 0 && t.ObjClass&4 != 0 && !GetServer().S().IsEnemyTo(owner, t)
}
func projectileSplash(u, exclude *server.Object, radius, inner float32, damage, kind int32) {
	C.nox_xxx_mapDamageUnitsAround_4E25B0((*C.float)(unsafe.Pointer(&u.PosVec)), C.float(radius), C.float(inner), C.int(damage), C.int(kind), asObjectC(u), asObjectC(exclude))
}
func projectilePush(u *server.Object, radius, inner, force float32) {
	C.nox_xxx_mapPushUnitsAround_52E040(unsafe.Pointer(&u.PosVec), C.float(radius), C.float(inner), C.float(force), asObjectC(u), 0, 0)
}
func projectileBoom(u, t *server.Object, n *types.Pointf) {
	init := memmap.PtrUint32(0x5d4594, 1567964)
	if *init == 0 {
		*memmap.PtrUint32(0x5d4594, 1567968) = uint32(floatToInt32(float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr("MagicMissileDamage")))))
		*memmap.PtrUint32(0x5d4594, 1567972) = uint32(floatToInt32(float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr("MagicMissileSplashDamage")))))
		for _, v := range []struct {
			off  uintptr
			name string
		}{{1567976, "MagicMissileRange"}, {1567980, "MagicMissilePushRange"}, {1567984, "MagicMissileForce"}} {
			*memmap.PtrFloat32(0x5d4594, v.off) = float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr(v.name)))
		}
		*init = 1
	}
	if projectileFriendly(u, t) {
		return
	}
	if u != nil {
		projectileFX(134, u)
	}
	if t != nil {
		if t.ObjClass&4 != 0 {
			if C.nox_xxx_checkInversionEffect_4FA4F0(inventoryInt(t), inventoryInt(u)) != 0 {
				Nox_xxx_changeOwner_52BE40(u, t)
				return
			}
			if t.HasEnchant(27) && projectileFront(t, u) {
				Nox_xxx_changeOwner_52BE40(u, t)
				inventorySound(122, t, 0, 0)
				return
			}
		}
		projectileDamage(t, u.FindOwnerChainPlayer(), u, int32(memmap.Uint32(0x5d4594, 1567968)), 7)
		C.nox_xxx_sMakeScorch_537AF0((*C.float)(unsafe.Pointer(&t.PosVec)), 0)
	} else if n != nil {
		collisionReflect(n, &u.VelVec)
		u.Direction2 = server.Dir16(C.nox_xxx_math_509ED0((*C.float2)(unsafe.Pointer(&u.VelVec))))
		u.VelVec.X = float32(float64(u.VelVec.X) * 0.5)
		u.VelVec.Y = float32(float64(u.VelVec.Y) * 0.5)
		projectileContact(u, int32(memmap.Uint32(0x5d4594, 1567968)), 7)
		return
	}
	projectileSplash(u, nil, memmap.Float32(0x5d4594, 1567976), 5, int32(memmap.Uint32(0x5d4594, 1567972)), 7)
	radius := memmap.Float32(0x5d4594, 1567980)
	projectilePush(u, radius, radius, memmap.Float32(0x5d4594, 1567984))
	inventorySound(84, u, 0, 0)
	GetServer().DelayedDelete(u)
}
func projectileFireball(u, t *server.Object) {
	data := u.CollideData
	active := true
	if t != nil && t.HasEnchant(27) && projectileFront(t, u) {
		damageReflect(u, t)
		core := GetServer().S()
		core.ObjClearOwner(u)
		core.ObjSetOwner(t, u)
		active = false
		inventorySound(122, t, 0, 0)
	}
	if projectileFriendly(u, t) || !active {
		return
	}
	damage := int32(*(*byte)(data))
	projectilePush(u, float32(float64(damage)*0.66666669), 0, 50)
	if t != nil {
		projectileDamage(t, u.FindOwnerChainPlayer(), u, int32(*(*byte)(data))>>1, 1)
	}
	damage = int32(*(*byte)(data))
	inner := float32(15)
	if noxflags.HasGame(2048) {
		inner = 0
	}
	projectileSplash(u, t, float32(float64(damage)*0.33333334), inner, damage>>1, 1)
	C.nox_xxx_netSparkExplosionFx_5231B0((*C.float)(unsafe.Pointer(&u.PosVec)), C.char(*(*byte)(data)))
	inventorySound(42, u, 0, 0)
	C.nox_xxx_sMakeScorch_537AF0((*C.float)(unsafe.Pointer(&u.PosVec)), 2)
	GetServer().DelayedDelete(u)
}
func projectilePixie(u, t *server.Object, n *types.Pointf) {
	if t != nil {
		if !GetServer().S().IsEnemyTo(u, t) || t.ObjClass&0x20006 == 0 || t.ObjFlags&0x8020 != 0 {
			return
		}
		owner := u.ObjOwner
		if owner != nil && owner.ObjClass&4 != 0 && owner.ObjFlags&2 != 0 {
			return
		}
		if t.ObjClass&4 != 0 {
			if C.nox_xxx_checkInversionEffect_4FA4F0(inventoryInt(t), inventoryInt(u)) != 0 {
				Nox_xxx_changeOwner_52BE40(u, t)
				return
			}
			if t.HasEnchant(27) && projectileFront(t, u) {
				Nox_xxx_changeOwner_52BE40(u, t)
				inventorySound(122, t, 0, 0)
				return
			}
		}
		projectileDamage(t, u.FindOwnerChainPlayer(), u, int32(*equipmentWord(u.CollideData, 0)), 11)
		inventorySound(96, u, 0, 0)
		GetServer().DelayedDelete(u)
	} else if n != nil {
		collisionReflect(n, &u.VelVec)
		u.Direction2 = server.Dir16(C.nox_xxx_math_509ED0((*C.float2)(unsafe.Pointer(&u.VelVec))))
		u.NewPos.X = float32(float64(u.VelVec.X) + float64(u.NewPos.X))
		y := float64(u.VelVec.Y) + float64(u.NewPos.Y)
		u.NewPos.Y = float32(y)
		gy := floatToInt32(float32(y * 0.043478262))
		gx := floatToInt32(float32(float64(u.NewPos.X) * 0.043478262))
		projectileWall(u, gx, gy, int32(*equipmentWord(u.CollideData, 0)), 11)
	} else {
		GetServer().DelayedDelete(u)
	}
}
