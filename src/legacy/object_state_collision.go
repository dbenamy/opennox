package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "server__script__script.h"
uint32_t nox_xxx_wallFlags(int i);
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func stateCloseDoor(u *server.Object, p unsafe.Pointer) {
	if u.ObjClass&0x80 != 0 && *equipmentWord(u.UpdateData, 16) == *equipmentWord(p, 0) && *equipmentWord(u.UpdateData, 20) == *equipmentWord(p, 4) {
		*(*byte)(unsafe.Add(u.UpdateData, 1)) = 0
		if bool(C.nox_common_gameFlags_check_40A5C0(4096)) {
			stateDoorNotify(u)
		}
	}
}
func stateDoorNotify(u *server.Object) int32 {
	*(*byte)(unsafe.Add(u.UpdateData, 48)) = 1
	return int32(C.sub_4D6A20(255, inventoryInt(u)))
}
func stateMonsterCollision(u, t *server.Object) unsafe.Pointer {
	return C.nox_xxx_scriptCallByEventBlock_502490(unsafe.Add(u.UpdateData, 1272), t.CObj(), u.CObj(), 22)
}
func stateMimicCollision(u, t *server.Object) unsafe.Pointer {
	if t != nil && t.ObjFlags&0x8000 == 0 && t.ObjClass&6 != 0 && C.nox_xxx_unitIsEnemyTo_5330C0(asObjectC(u), asObjectC(t)) != 0 && !u.MonsterActionIsScheduled(15) {
		if st := u.MonsterPushAction(43); st != nil {
			*equipmentWord(unsafe.Pointer(st), 4) = GetServer().S().Frame()
		}
		if st := u.MonsterPushAction(15); st != nil {
			*equipmentWord(unsafe.Pointer(st), 4) = *equipmentWord(t.CObj(), 56)
			*equipmentWord(unsafe.Pointer(st), 8) = *equipmentWord(t.CObj(), 60)
			*equipmentWord(unsafe.Pointer(st), 12) = GetServer().S().Frame()
		}
	}
	return stateMonsterCollision(u, t)
}
func stateChargeImpact(u, t *server.Object) bool {
	if t == nil {
		return true
	}
	class := t.ObjClass
	if class&6 == 0 || t.HealthData.Cur == 0 && t.HealthData.Max != 0 {
		if class&0x400000 == 0 && int8(t.ObjFlags) >= 0 && float64(t.Mass) <= float64(u.Mass) {
			return false
		}
	}
	if class&0x80 != 0 {
		return *(*byte)(unsafe.Add(t.UpdateData, 1)) != 0
	}
	return t.ObjFlags&9 == 0 && class&1 == 0
}
func stateChargeMoveBack(u *server.Object) {
	C.nox_xxx_unitMove_4E7010(asObjectC(u), (*C.float2)(unsafe.Pointer(&u.PrevPos)))
}
func stateChargeStun(u *server.Object) {
	duration := floatToInt32(float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerStunDuration"))))
	C.nox_xxx_buffApplyTo_4FF380(asObjectC(u), 5, C.short(duration), 5)
}
func stateCharge(u, t *server.Object) {
	C.nox_xxx_playerSetState_4FA020(asObjectC(u), 13)
	C.nox_xxx_earthquakeSend_4D9110((*C.float)(unsafe.Pointer(&u.PosVec)), 10)
	Sub_4FC300(u, 1)
	if t != nil {
		damage := floatToInt32(float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerDamage"))))
		if t.ObjClass&0x400000 == 0 {
			C.sub_4E86E0(inventoryInt(u), (*C.float)(t.CObj()))
		}
		projectileDamage(t, u.FindOwnerChainPlayer(), u, damage, 2)
		if t.ObjClass&0x20006 != 0 {
			stateChargeMoveBack(u)
			return
		}
		stateChargeStun(u)
	} else {
		wall := *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 296))
		if wall != nil && C.nox_xxx_wallFlags(C.int(*(*byte)(unsafe.Add(wall, 1))))&5 == 0 {
			stateChargeMoveBack(u)
			return
		}
		C.nox_xxx_aud_501960(171, asObjectC(u), 0, 0)
		stateChargeStun(u)
		x, y := projectileGrid(u.NewPos)
		projectileWall(u, x, y, 100, 2)
	}
	pain := floatToInt32(float32(float64(C.nox_xxx_gamedataGetFloat_419D40(internCStr("BerserkerPainRatio"))) * float64(u.HealthData.Cur)))
	if pain < 1 {
		pain = 1
	}
	resourceDamage(u, pain)
	stateChargeMoveBack(u)
}
func statePlayerCollision(u, t *server.Object) {
	if GetServer().S().Abils.IsActive(u, 1) && stateChargeImpact(u, t) {
		stateCharge(u, t)
	}
	if t != nil && t.ObjClass&4 != 0 && t.ObjFlags&0x8000 == 0 && u.Buffs&(1<<16) != 0 && uint32(u.BuffsDur[16]) < 14*GetServer().S().TickRate() {
		C.nox_xxx_buffApplyTo_4FF380(asObjectC(t), 16, C.short(uint16(15*GetServer().S().TickRate())), C.char(u.BuffsPower[16]))
		C.nox_xxx_spellBuffOff_4FF5B0(asObjectC(u), 16)
	}
}
