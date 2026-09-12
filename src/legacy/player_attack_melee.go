package legacy

/*
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2488652,dword_5d4594_2488656,dword_5d4594_2488660;
*/
import "C"
import (
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type attackRecord struct {
	Damage    float32
	Type      byte
	_         [3]byte
	Radius    float32
	Owner     *server.Object
	Pos       types.Pointf
	HitStatic uint32
	Weapon    *server.Object
	Front     byte
	_         [3]byte
}

func attackPreEffects(t, u, it *server.Object, r *attackRecord) int {
	if it == nil {
		return 0
	}
	if !noxflags.HasGamePlay(1) && u != nil && u.ObjClass&6 != 0 && !GetServer().S().IsEnemyTo(u, t) {
		return 0
	}
	if t.HasEnchant(23) || t.HasEnchant(27) {
		return 1
	}
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
		if m != nil && m.AttackPreHit52.Fnc != nil {
			ccall.CallVoidPtr5(m.AttackPreHit52.Fnc, unsafe.Pointer(m), it.CObj(), u.CObj(), t.CObj(), unsafe.Pointer(r))
		}
	}
	return 0
}
func attackItemEffects(it, u *server.Object, r *attackRecord) int {
	for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4) {
		if m != nil && m.Attack40.Fnc != nil {
			ccall.CallVoidPtr5(m.Attack40.Fnc, unsafe.Pointer(m), it.CObj(), u.CObj(), nil, unsafe.Pointer(r))
		}
	}
	return 0
}
func attackHit(t *server.Object, r *attackRecord) {
	if t == nil || r == nil || r.Owner == nil || t == r.Owner {
		return
	}
	u := r.Owner
	if byte(C.nox_server_testTwoPointsAndDirection_4E6E50((*C.float2)(unsafe.Pointer(&u.PosVec)), C.int(int16(u.Direction1)), (*C.float2)(unsafe.Pointer(&t.PosVec))))&r.Front == 0 {
		return
	}
	if t.ObjFlags&0x8040 != 0 || (r.HitStatic == 0 && t.ObjFlags&8 != 0) {
		return
	}
	if t.Material != 0x4000 {
		C.dword_5d4594_2488656 = 1
	}
	if attackRay(u.PosVec, t.PosVec, 5) == 0 {
		return
	}
	attackPreEffects(t, u, r.Weapon, r)
	ccall.CallIntUPtr5(t.Damage, uintptr(t.CObj()), uintptr(r.Owner.CObj()), uintptr(r.Weapon.CObj()), uintptr(uint32(effectsTruncWord(float64(r.Damage)+0.5))), uintptr(r.Type))
	if noxflags.HasGame(2048) && r.Owner.ObjClass&4 != 0 && t.ObjClass&2 == 0 && t.HealthData != nil && t.HealthData.Max != 0 && t.ObjFlags&0x8020 == 0 {
		C.nox_xxx_netSendPointFx_522FF0(-117, (*C.float2)(unsafe.Pointer(&t.PosVec)))
	}
	if r.Weapon == nil {
		u = r.Owner
		var anim byte
		if u.ObjClass&4 != 0 {
			anim = *(*byte)(unsafe.Add(*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 276)), 8))
		} else {
			anim = *(*byte)(unsafe.Add(u.UpdateData, 2068))
		}
		if anim == 25 {
			GetServer().ApplyForce(t, u.PosVec, 20)
		}
	}
}
func attackNearest(t, u *server.Object) {
	if u == nil || t == nil || u == t {
		return
	}
	f := t.ObjFlags
	if f&0x8049 != 0 || (t.ObjClass&6 == 0 && f&0x10 != 0 && f&0x80 == 0) || (t.ObjClass&6 != 0 && !GetServer().S().IsEnemyTo(u, t)) || !GetServer().S().CanInteract(u, t, 1) {
		return
	}
	limit := math.Float32frombits(uint32(C.dword_5d4594_2488652))
	if !(limit > 0) {
		return
	}
	dx := float32(float64(t.PosVec.X) - float64(u.PosVec.X))
	dy64 := float64(t.PosVec.Y) - float64(u.PosVec.Y)
	dy := float32(dy64)
	distance := math.Sqrt(dy64*float64(dy) + float64(dx)*float64(dx))
	if distance == 0 {
		return
	}
	dir := attackDirection(u)
	if !(float64(dy)/distance*float64(dir.Y)+float64(dx)/distance*float64(dir.X) > 0.5) {
		return
	}
	if t.Shape.Kind == 2 {
		distance -= float64(t.Shape.Circle.R)
	} else if t.Shape.Kind == 3 {
		delta := types.Pointf{X: dx, Y: dy}
		n := float64(C.sub_54A990((*C.float2)(unsafe.Pointer(&u.PosVec)), C.float(limit), inventoryInt(t), (*C.float2)(unsafe.Pointer(&delta))))
		if n < 0 {
			return
		}
		distance = float64(limit) - n
	}
	if distance < 0 {
		distance = 0
	}
	old := objectFromInt(C.int(C.dword_5d4594_2488660))
	if (distance < float64(limit) || old != nil && old.ObjClass&2 == 0 && t.ObjClass&2 != 0) && (old == nil || old.ObjClass&2 == 0) {
		C.dword_5d4594_2488652 = C.uint32_t(math.Float32bits(float32(distance)))
		C.dword_5d4594_2488660 = C.uint32_t(uintptr(t.CObj()))
	}
}
func attackTrace(u *server.Object, r *attackRecord) int {
	if u == nil || r == nil {
		return 0
	}
	C.dword_5d4594_2488656 = 0
	C.dword_5d4594_2488660 = 0
	// This dependency includes shape extents; EachObjInCircle uses centers only.
	extra := float64(0)
	if r.Weapon != nil && r.Weapon.ObjSubClass&0x4000 != 0 {
		C.sub_518040(C.int(uintptr(unsafe.Pointer(&r.Pos))), C.float(r.Radius), C.int(uintptr(C.sub_538510)), C.int(uintptr(unsafe.Pointer(r))))
		extra = 25
	} else {
		C.dword_5d4594_2488652 = C.uint32_t(math.Float32bits(r.Radius))
		C.sub_518040(C.int(uintptr(unsafe.Pointer(&u.PosVec))), C.float(r.Radius), C.int(uintptr(C.sub_5386A0)), inventoryInt(u))
		if C.dword_5d4594_2488660 != 0 {
			attackHit(objectFromInt(C.int(C.dword_5d4594_2488660)), r)
		}
	}
	bounds := [4]int32{
		floatToInt32(float32(float64(r.Pos.X)-float64(r.Radius)-extra)) / 23,
		floatToInt32(float32(float64(r.Pos.Y)-float64(r.Radius)-extra)) / 23,
		floatToInt32(float32(float64(r.Radius)+float64(r.Pos.X)+extra)) / 23,
		floatToInt32(float32(float64(r.Pos.Y)+float64(r.Radius)+extra)) / 23,
	}
	it := r.Weapon
	if it == nil {
		it = u
	}
	C.nox_xxx_mapDamageToWalls_534FC0((*C.int4)(unsafe.Pointer(&bounds)), unsafe.Pointer(&r.Pos), C.float(float32(extra+float64(r.Radius))), C.int(effectsTruncWord(float64(r.Damage)+0.5)), C.int(r.Type), it.CObj())
	if r.Weapon != nil && C.dword_5d4594_2488656 != 0 {
		damage := float32(float64(C.nox_xxx_gamedataGetFloat_419D40(internCStr("ItemDamagePercentage"))) * float64(r.Damage))
		target := C.int(C.dword_5d4594_2488660)
		damageDurability(r.Weapon, r.Owner, objectFromInt(target), objectFromInt(target), damage, int32(r.Type), true)
	}
	return int(C.dword_5d4594_2488656)
}
func attackWarcry(t, u *server.Object) int16 {
	if u == nil {
		return 0
	}
	if t == nil {
		return 0
	}
	if t.ObjClass&2 != 0 && t.ObjSubClass&0x20000 != 0 && t.ObjFlags&0x8020 == 0 {
		spellLifeApplyBuff(t, 5, 90, 3)
	}
	return int16(uintptr(t.CObj()))
}
