package legacy

/*
#include "GAME4_2.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_1568868;
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func controlFindStart(out *types.Pointf, u *server.Object) {
	s := GetServer().S()
	if C.dword_5d4594_1568868 == 0 {
		C.dword_5d4594_1568868 = C.uint32_t(s.Types.IndByID("PlayerStart"))
	}
	if u == nil {
		return
	}
	team := int32(0)
	if C.nox_xxx_servObjectHasTeam_419130(C.int(uintptr(unsafe.Add(u.CObj(), 48)))) != 0 {
		team = int32(*controlByte(u.CObj(), 52))
		C.nox_xxx_getTeamByID_418AB0(C.int(team))
	}
	var last, best *server.Object
	var choices []*server.Object
	for it := s.Objs.List; it != nil; it = it.ObjNext {
		if uint32(it.TypeInd) == uint32(C.dword_5d4594_1568868) {
			last = it
			if controlStartEligible(it, team) {
				choices = append(choices, it)
			}
		}
	}
	if len(choices) == 0 {
		if last == nil {
			*out = types.Pointf{X: 2000, Y: 2000}
		} else {
			*out = last.PosVec
		}
		return
	}
	// Pick the start farthest from its nearest enemy. With no enemies, choose
	// randomly among eligible starts, preserving traversal order and RNG use.
	noEnemies := true
	maxDist := float32(0)
	for _, it := range choices {
		minDist := float32(10000000)
		for p := s.Players.FirstUnit(); p != nil; p = s.Players.NextUnit(p) {
			if p != u && GetServer().S().IsEnemyTo(u, p) {
				dx := float64(it.PosVec.X) - float64(p.PosVec.X)
				dy := float64(it.PosVec.Y) - float64(p.PosVec.Y)
				dist := dy*dy + dx*dx
				if dist < float64(minDist) {
					minDist = float32(dist)
				}
				noEnemies = false
			}
		}
		if minDist > maxDist {
			maxDist = minDist
			best = it
		}
	}
	if noEnemies || best == nil {
		best = choices[s.Rand.Logic.IntClamp(0, len(choices)-1)]
	}
	*out = best.PosVec
}
func controlDropBall(u *server.Object) int32 {
	typ := stateType(1568872, "GameBall")
	for it := u.Field129; it != nil; it = it.Field128 {
		if uint32(it.TypeInd) != typ {
			continue
		}
		it.ObjFlags &^= 0x40
		C.nox_xxx_objectApplyForce_52DF80((*C.float)(unsafe.Pointer(&u.PosVec)), (*C.nox_object_t)(it.CObj()), 100)
		*controlPtr(it.CObj(), 520) = nil
		GetServer().S().ObjSetOwner(nil, it)
		C.nox_xxx_aud_501960(926, (*C.nox_object_t)(u.CObj()), 0, 0)
		C.sub_4E8290(1, 0)
		return 1
	}
	return 0
}
func controlNearStart(u *server.Object, out *types.Pointf) int32 {
	*out = u.PosVec
	var result C.int
	for i := 0; i < 32; i++ {
		C.sub_4ED970(60, (*C.float2)(unsafe.Pointer(&u.PosVec)), (*C.float2)(unsafe.Pointer(out)))
		result = C.nox_xxx_mapTileAllowTeleport_411A90((*C.float2)(unsafe.Pointer(out)))
		if result == 0 {
			break
		}
	}
	return int32(result)
}
func controlWalkWaypoint(u *server.Object) int32 {
	ptr := controlPtr(u.UpdateData, 168+4*int(*controlByte(u.UpdateData, 181)))
	it := (*server.Object)(*ptr)
	if it == nil {
		return 0
	}
	dx := float32(float64(it.PosVec.X) - float64(u.PosVec.X))
	dy := float64(it.PosVec.Y) - float64(u.PosVec.Y)
	pos := types.Pointf{X: dx, Y: float32(dy)}
	if dy*float64(pos.Y)+float64(dx)*float64(dx) < 100 {
		GetServer().DelayedDelete(it)
		*ptr = nil
		return 0
	}
	u.Direction2 = server.Dir16(C.nox_xxx_math_509ED0((*C.float2)(unsafe.Pointer(&pos))))
	if u.Buffs&(1<<3) != 0 {
		u.Direction2 = server.Dir16(controlConfusedDirection(u))
	}
	off := uintptr(int(int16(u.Direction2)) * 8)
	for i := 0; i < 2; i++ {
		f := (*float32)(unsafe.Add(u.CObj(), 88+4*i))
		*f = float32(float64(*memmap.PtrFloat32(0x587000, 194136+off+uintptr(i*4)))*float64(*(*float32)(unsafe.Add(u.CObj(), 544))) + float64(*f))
	}
	return 1
}
func controlInputAttack(u *server.Object) {
	if !controlAimsAtEnemy(u) {
		return
	}
	d := u.UpdateData
	weapon := *equipmentWord(controlPlayer(u), 4)
	setAttack := func() { *equipmentWord(u.CObj(), 136) = GetServer().S().Frame(); *controlByte(d, 236) = 0 }
	if weapon == 0 {
		if *controlByte(d, 88) != 1 {
			C.nox_xxx_playerSetState_4FA020((*C.nox_object_t)(u.CObj()), 1)
		}
		return
	}
	if weapon&0x47f0000 != 0 && controlActionState(u) != 29 {
		it := controlObject(d, 104)
		data := it.UseData.Ptr
		if *controlByte(data, 108) != 0 || *controlByte(data, 109) == 0 {
			setAttack()
			C.nox_xxx_playerSetState_4FA020((*C.nox_object_t)(u.CObj()), 1)
			C.nox_xxx_useByNetCode_53F8E0(inventoryInt(u), inventoryInt(it))
		} else if controlSubStamina(u, 45) != 0 {
			*equipmentWord(data, 96) |= 2
			setAttack()
			C.nox_xxx_playerSetState_4FA020((*C.nox_object_t)(u.CObj()), 1)
		}
	} else if *controlByte(d, 88) != 1 {
		cost := controlWeaponStamina(weapon)
		if controlSubStamina(u, cost) != 0 {
			setAttack()
			if C.nox_xxx_playerSetState_4FA020((*C.nox_object_t)(u.CObj()), 1) == 0 {
				controlAdjustStamina(u, -int8(cost))
			}
		}
	}
	C.nox_xxx_spellBuffOff_4FF5B0((*C.nox_object_t)(u.CObj()), 0)
	C.nox_xxx_spellBuffOff_4FF5B0((*C.nox_object_t)(u.CObj()), 23)
	C.nox_xxx_spellCancelDurSpell_4FEB10(67, (*C.nox_object_t)(u.CObj()))
}
func controlFollowEnemy(u *server.Object) int32 {
	if u == nil {
		return 0
	}
	it := controlObject(u.CObj(), 520)
	if it == nil {
		return 0
	}
	it = it.FindOwnerChainPlayer()
	if it.ObjFlags&0x20 != 0 || it.ObjClass&2 != 0 || it.ObjClass&4 != 0 && *controlByte(controlPlayer(it), 3680)&1 != 0 {
		return 0
	}
	C.nox_xxx_playerCameraFollow_4E6060((*C.nox_object_t)(u.CObj()), (*C.nox_object_t)(it.CObj()))
	return 1
}
