package legacy

/*
#include "defs.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME2_3.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "server__object__health.h"
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

type lifecycleAIAction struct{ typ ai.ActionType }

func (a lifecycleAIAction) Type() ai.ActionType { return a.typ }
func (a lifecycleAIAction) Start(u *server.Object) {
	switch a.typ {
	case ai.ACTION_DYING:
		lifecycleDyingStart(u)
	case ai.ACTION_DEAD:
		lifecycleDeadStart(u)
	}
}
func (a lifecycleAIAction) Update(u *server.Object) {
	switch a.typ {
	case ai.ACTION_HUNT:
		lifecycleHunt(u)
	case ai.ACTION_PICKUP_OBJECT:
		lifecyclePickup(u)
	case ai.ACTION_DYING, ai.ACTION_GET_UP:
		if u.UpdateDataMonster().Field120_3 != 0 {
			u.MonsterPopAction()
		}
	case ai.ACTION_DEAD:
		lifecycleDeadUpdate(u)
	}
}
func (a lifecycleAIAction) End(u *server.Object) {
	if a.typ == ai.ACTION_DYING {
		lifecycleBurnDeleteCheck(u)
	}
}
func (a lifecycleAIAction) Cancel(*server.Object) {}

func init() {
	for _, typ := range []ai.ActionType{ai.ACTION_HUNT, ai.ACTION_PICKUP_OBJECT, ai.ACTION_DYING, ai.ACTION_DEAD, ai.ACTION_GET_UP} {
		server.RegisterAIAction(lifecycleAIAction{typ})
	}
}

func lifecycleHunt(u *server.Object) {
	u.UpdateDataMonster().Aggression = math.Float32frombits(0x3f547ae1)
	if st := u.MonsterPushAction(ai.ACTION_ROAM); st != nil {
		st.Args[0] = 0
		st.Args[2] = st.Args[2]&^0xff | 0x80
	}
}
func lifecyclePickup(u *server.Object) {
	h := u.UpdateDataMonster().AIStackHead()
	if t := h.ArgObj(0); t != nil {
		dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		if dx*dx+dy*dy < 5625 && GetServer().S().CanInteract(u, t, 0) {
			Nox_xxx_inventoryServPlace_4F36F0(u, t, 1, 1)
			// C reloads the action argument after inventory placement.
			if t = h.ArgObj(0); t != nil && t.SubClass()&0x10 != 0 {
				effectsUse(u, t)
			}
		}
	}
	u.MonsterPopAction()
}
func lifecycleDyingStart(u *server.Object) {
	ud := u.UpdateDataMonster()
	combatSound(u, 15)
	GetServer().NoxScriptC().ScriptCallback(&ud.ScriptDeath, nil, u, server.ScriptEventType(7))
	if p := ud.MonsterDef.DieFunc228; p != nil {
		ccall.CallIntPtr(p, u.CObj())
	}
}
func lifecycleIsZombie(u *server.Object) bool {
	return monsterIsZombie(u)
}
func lifecycleBurnDeleteCheck(u *server.Object) {
	if lifecycleIsZombie(u) && u.UpdateDataMonster().StatusFlags&0x80000 != 0 {
		lifecycleBurnDelete(u)
	}
}
func lifecycleBurnDelete(u *server.Object) {
	if owner := u.ObjOwner; owner != nil && owner.Class().Has(object.ClassPlayer) {
		pl := owner.UpdateDataPlayer().Player
		u.ObjSubClass &^= 0x80
		C.nox_xxx_netFxShield_0_4D9200(C.int(pl.PlayerInd), C.int(uintptr(u.CObj())))
		GetServer().S().Players.Nox_xxx_netUnmarkMinimapObj_417300(ntype.PlayerInd(pl.PlayerInd), u, 1)
	}
	C.nox_xxx_soloMonsterKillReward_4EE500_obj_health(C.int(uintptr(u.CObj())))
	Nox_xxx_sMakeScorch_537AF0(u.PosVec, 1)
	GetServer().DelayedDelete(u)
}
func lifecycleDeadStart(u *server.Object) {
	ud := u.UpdateDataMonster()
	if u.Field131 == 14 && u.SubClass()&0x10000 != 0 {
		lifecycleReleasedSoul(u)
	}
	u.VelVec = types.Pointf{}
	u.ForceVec = types.Pointf{}
	u.Pos24 = types.Pointf{}
	if p := ud.MonsterDef.DeadFunc232; p != nil {
		ccall.CallVoidPtr(p, u.CObj())
	}
	if lifecycleIsZombie(u) {
		lo := C.nox_float2int(C.float(GetServer().S().Balance.FloatInd("ZombieDeadDuration", 0)))
		hi := C.nox_float2int(C.float(GetServer().S().Balance.FloatInd("ZombieDeadDuration", 1)))
		ud.Field123 = uint32(nox_common_randomInt_415FA0(int(lo), int(hi)))
		u.ObjFlags |= 0x10
	} else {
		u.ObjFlags |= 0x18
	}
}
func lifecycleReleasedSoul(u *server.Object) {
	id := memmap.PtrUint32(0x5D4594, 2489456)
	if *id == 0 {
		*id = uint32(GetServer().S().Types.IndByID("ReleasedSoul"))
	}
	if soul := GetServer().S().NewObjectByTypeInd(int(*id)); soul != nil {
		GetServer().CreateObjectAt(soul, nil, u.PosVec)
		soul.Direction1, soul.Direction2 = u.Direction1, u.Direction1
	}
}
func lifecycleDeadUpdate(u *server.Object) {
	ud := u.UpdateDataMonster()
	if lifecycleIsZombie(u) {
		if ud.StatusFlags&0x80000 != 0 {
			C.nox_xxx_netSparkExplosionFx_5231B0((*C.float)(unsafe.Pointer(&u.PosVec)), 100)
			lifecycleBurnDelete(u)
		} else if GetServer().S().Frame()-ud.Field137 > ud.Field123 && ud.StatusFlags&0x100000 == 0 && ud.CurrentEnemy != nil {
			lifecycleRaiseZombie(u)
		}
		return
	}
	GetServer().S().Objs.RemoveFromUpdatable(u)
	u.Update = nil // C clears Object+744; Object+748 UpdateData remains live.
	lifecycleReset(u)
	if ud.MonsterDef.StatusFlags92&1 != 0 {
		GetServer().DelayedDelete(u)
	}
}
func lifecycleReset(u *server.Object) {
	ud := u.UpdateDataMonster()
	ud.Field523_2 = 0
	if n := int32(ud.Field74); n > 0 {
		clear(unsafe.Slice(&ud.Field75, int(n)))
	}
	ud.Field74, ud.Field91 = 0, 0
	if n := int(ud.Field282_1); n > 0 {
		clear(unsafe.Slice(&ud.Field283, n))
	}
	ud.Field282_1 = 0
	ud.CurrentEnemy = nil
}
func lifecycleRaiseZombie(u *server.Object) uint32 {
	if !lifecycleIsZombie(u) {
		return 0
	}
	if got := u.UpdateDataMonster().AIStackHead().Type(); got != ai.ACTION_DEAD {
		return uint32(got)
	}
	u.MonsterPopAction()
	u.MonsterPushAction(ai.DEPENDENCY_UNINTERRUPTABLE)
	u.MonsterPushAction(ai.ACTION_GET_UP)
	GetServer().S().Audio.EventObj(469, u, 0, 0)
	resourceRestoreHP(u)
	u.ObjFlags &= 0xffff7fa7
	return uint32(u.ObjFlags)
}

func lifecycleCandidate(off uintptr) *server.Object {
	return (*server.Object)(unsafe.Pointer(uintptr(*memmap.PtrUint32(0x5D4594, off))))
}
func lifecycleSetCandidate(off uintptr, u *server.Object) {
	*memmap.PtrUint32(0x5D4594, off) = uint32(uintptr(u.CObj()))
}
func lifecycleFoodSearch(u *server.Object, radius float32, usable bool) *server.Object {
	var candidateOff, minOff uintptr
	if usable {
		candidateOff, minOff = 2489440, 2489448
	} else {
		candidateOff, minOff = 2489452, 2489444
	}
	*memmap.PtrUint32(0x5D4594, candidateOff) = 0
	*memmap.PtrFloat32(0x5D4594, minOff) = 10000000
	GetServer().S().Map.EachObjInCircle(u.PosVec, radius, func(t *server.Object) bool {
		if usable {
			if !t.Class().Has(0x1000000) || !Nox_xxx_playerClassCanUseItem_57B3D0(t, player.Class(0)) {
				return true
			}
		} else {
			sub := uint32(t.SubClass())
			if !t.Class().Has(0x10) || sub&4 != 0 || (sub&0x80 != 0 && u.Poison540 == 0) ||
				(sub&8 != 0 && (sub&0x10 == 0 || !noxflags.HasGame(0x2000) || u.SubClass()&0x10 == 0)) {
				return true
			}
		}
		if !GetServer().S().CanInteract(u, t, 0) {
			return true
		}
		dx, dy := float64(t.PosVec.X)-float64(u.PosVec.X), float64(t.PosVec.Y)-float64(u.PosVec.Y)
		if dist := dx*dx + dy*dy; dist < float64(*memmap.PtrFloat32(0x5D4594, minOff)) {
			*memmap.PtrFloat32(0x5D4594, minOff) = float32(dist)
			lifecycleSetCandidate(candidateOff, t)
		}
		return true
	})
	return lifecycleCandidate(candidateOff)
}

//export nox_xxx_mobRaiseZombie_534AB0
func nox_xxx_mobRaiseZombie_534AB0(a1 C.int) C.uint {
	return C.uint(lifecycleRaiseZombie(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1))))))
}

//export nox_xxx_mobSearchEdible_544A00
func nox_xxx_mobSearchEdible_544A00(a1 *nox_object_t, r C.float) C.int {
	if u := lifecycleFoodSearch(asObjectS(a1), float32(r), false); u != nil {
		return C.int(uintptr(u.CObj()))
	}
	return 0
}

//export sub_544AE0
func sub_544AE0(a1 C.int, r C.float) C.int {
	if u := lifecycleFoodSearch(asObjectS((*nox_object_t)(unsafe.Pointer(uintptr(a1)))), float32(r), true); u != nil {
		return C.int(uintptr(u.CObj()))
	}
	return 0
}
