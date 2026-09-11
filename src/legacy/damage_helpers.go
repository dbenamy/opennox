package legacy

/*
#include "GAME1.h"
#include "common__random.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5_2.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

func damageTypeByName(name string) int32 {
	for i := 0; i < 18; i++ {
		if name == alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 200728+uintptr(i*4)))) {
			return int32(i)
		}
	}
	return 18
}
func damageReflect(u, t *server.Object) int32 {
	dx, dy := float64(u.PosVec.X)-float64(t.PosVec.X), float64(u.PosVec.Y)-float64(t.PosVec.Y)
	u.Direction1 += 128
	dir := int16(u.Direction1)
	negY := float32(-dy)
	length := float32(math.Sqrt(dy*dy+dx*dx) + 0.1)
	proj := (dy*float64(u.VelVec.Y) + dx*float64(u.VelVec.X)) / float64(length)
	px, py := float32(proj*dx/float64(length)), float32(proj*dy/float64(length))
	tangent := (float64(negY)*float64(u.VelVec.X) + dx*float64(u.VelVec.Y)) / float64(length)
	tx := float32(tangent * float64(negY) / float64(length))
	u.VelVec.X = tx - px
	u.VelVec.Y = float32(dx*tangent/float64(length) - float64(py))
	if dir >= 256 {
		u.Direction1 = server.Dir16(uint16(dir - 256))
	}
	u.NewPos = u.PrevPos
	return int32(uintptr(u.CObj()))
}
func damageDefend(u, source, weapon *server.Object, damage *int32, kind int32) int32 {
	result := int32(uintptr(u.CObj()))
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		result = int32(it.ObjFlags)
		if it.ObjFlags&0x100 == 0 {
			continue
		}
		for _, m := range unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[2:] {
			if m != nil && m.Defend76.Fnc != nil {
				data, free := alloc.New([2]int32{})
				*data = [2]int32{*damage, kind}
				ccall.CallVoidPtr6(m.Defend76.Fnc, unsafe.Pointer(m), it.CObj(), u.CObj(), weapon.CObj(), source.CObj(), unsafe.Pointer(data))
				*damage = data[0]
				free()
			}
		}
		result = 0
	}
	return result
}
func damagePre(u, source, weapon *server.Object, damage *int32) int32 {
	for _, m := range unsafe.Slice((**server.ModifierEff)(weapon.InitData), 4) {
		if m != nil && m.AttackPreDmg64.Fnc != nil {
			ccall.CallVoidPtr5(m.AttackPreDmg64.Fnc, unsafe.Pointer(m), weapon.CObj(), source.CObj(), u.CObj(), unsafe.Pointer(damage))
		}
	}
	return 0
}
func damageMelee(u, it *server.Object) bool {
	if it == nil {
		return u.ObjClass&4 != 0 || u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0
	}
	if it.ObjClass&0x1000 != 0 {
		if it.ObjSubClass&0x47f0000 == 0 || *(*byte)(unsafe.Add(it.UseData.Ptr, 96))&2 != 0 {
			return true
		}
	} else if it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x47f00fe == 0 {
		return true
	}
	return it.ObjClass&2 != 0
}
func damageFriendlyWeapon(it *server.Object) bool {
	return it != nil && it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x4000 != 0
}
func damageFraction(u *server.Object, out *int32, value float32) {
	if u == nil || out == nil {
		return
	}
	var fraction *float32
	if u.ObjClass&4 != 0 {
		fraction = (*float32)(unsafe.Add(u.UpdateData, 84))
	} else if u.ObjClass&2 != 0 {
		fraction = (*float32)(unsafe.Add(u.UpdateData, 4))
	}
	if fraction != nil {
		value = float32(float64(value) + float64(*fraction))
	}
	n := floatToInt32(value)
	*out = n
	if fraction != nil {
		*fraction = float32(float64(value) - float64(n))
	}
}
func damageItemReport(ind int32, it *server.Object, old, now uint16) int32 {
	result := int32(uintptr(unsafe.Pointer(it.HealthData)))
	if it.HealthData == nil {
		return result
	}
	if bool(C.nox_common_gameFlags_check_40A5C0(2048)) {
		return int32(C.nox_xxx_itemReportHealth_4D87A0(C.int(ind), asObjectC(it)))
	}
	before := C.sub_57B190(C.ushort(old), C.ushort(it.HealthData.Max))
	result = int32(C.sub_57B190(C.ushort(now), C.ushort(it.HealthData.Max)))
	if int32(before) != result {
		result = int32(C.nox_xxx_itemReportHealth_4D87A0(C.int(ind), asObjectC(it)))
	}
	return result
}
func damageDurability(it, holder, source, weapon *server.Object, value float32, kind int32, weaponOnly bool) {
	if it == nil || it.HealthData == nil {
		return
	}
	if weaponOnly && it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x7800000 != 0 {
		return
	}
	fraction := (*float32)(it.UpdateData)
	if m := unsafe.Slice((**server.ModifierEff)(it.InitData), 4)[1]; m != nil && m.Defend76.Fnc != nil {
		data, free := alloc.New(float32(0))
		*data = value
		ccall.CallVoidPtr6(m.Defend76.Fnc, unsafe.Pointer(m), it.CObj(), holder.CObj(), weapon.CObj(), source.CObj(), unsafe.Pointer(data))
		value = *data
		free()
	}
	value = float32(float64(value) + float64(*fraction))
	n := floatToInt32(value)
	*fraction = float32(float64(value) - float64(n))
	if n <= 0 {
		return
	}
	old := it.HealthData.Cur
	projectileDamage(it, source, weapon, n, kind)
	if (holder != nil || !weaponOnly) && holder.ObjClass&4 != 0 && old != it.HealthData.Cur {
		damageItemReport(int32(uint8(holder.UpdateDataPlayer().Player.PlayerInd)), it, old, it.HealthData.Cur)
	}
}
func damageInventory(u, source, weapon *server.Object, amount int32, kindBits float32) {
	var coeff float32
	if u.ObjClass&4 != 0 {
		coeff = *(*float32)(unsafe.Add(u.UpdateData, 228))
	} else if u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0 {
		coeff = *(*float32)(unsafe.Add(u.UpdateData, 2072))
	} else {
		return
	}
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x2000000 != 0 && it.ObjFlags&0x100 != 0 {
			value := float32(float64(C.nox_xxx_itemApplyDefendEffect_415C00(inventoryInt(it))) / float64(coeff) * float64(amount))
			damageDurability(it, u, source, weapon, value, int32(math.Float32bits(kindBits)), false)
		}
	}
}
func damageConductivity(u *server.Object) float64 {
	value := float32(1)
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		if it.ObjClass&0x2000000 == 0 || *(*byte)(unsafe.Add(it.CObj(), 24))&0x10 == 0 || it.ObjFlags&0x100 == 0 || C.sub_4133D0(asObjectC(it)) != 0 {
			continue
		}
		var add float64
		if it.ObjSubClass&0x2000000 != 0 {
			add = float64(memmap.Float32(0x587000, 201108))
		} else {
			add = float64(C.sub_415BD0(inventoryInt(it)))
		}
		value = float32(add + float64(value))
	}
	return float64(value)
}
func damageBlockingItem(u, source, weapon *server.Object, mask int32, value float32, kind int32, weaponOnly bool) int32 {
	it := u.InvFirstItem
	for it != nil {
		if it.ObjFlags&0x100 != 0 && uint32(mask)&uint32(it.ObjSubClass) != 0 {
			break
		}
		it = it.InvNextItem
	}
	actual := weapon
	if actual == nil {
		actual = source
	}
	damageDurability(it, u, source, actual, value, kind, weaponOnly)
	if !weaponOnly && it == nil {
		return 0
	}
	if it.ObjFlags&0x20 == 0 {
		return 0
	}
	if u.ObjClass&4 != 0 {
		C.nox_xxx_playerSetState_4FA020(asObjectC(u), 13)
	} else {
		C.nox_xxx_monsterPopAction_50A160(asObjectC(u))
	}
	return 1
}
