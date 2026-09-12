package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2487712;
*/
import "C"
import (
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func spellEffectObject(p unsafe.Pointer, off int) *server.Object {
	return (*server.Object)(*controlPtr(p, off))
}
func spellEffectPos(p unsafe.Pointer, off int) types.Pointf {
	return *(*types.Pointf)(unsafe.Add(p, off))
}
func spellEffectScalar(name string) float64 {
	return float64(C.nox_xxx_gamedataGetFloat_419D40(internCStr(name)))
}
func spellEffectTable(name string, i int32) float64 {
	return float64(C.nox_xxx_gamedataGetFloatTable_419D70(internCStr(name), C.int(i)))
}
func spellEffectAudio(id, phase int32, u *server.Object) {
	s := GetServer().S()
	s.Audio.EventObj(s.Spells.DefByInd(spell.ID(id)).GetAudio(int(phase)), u, 0, 0)
}
func spellEffectPosAudio(id, phase int32, pos types.Pointf) {
	s := GetServer().S()
	s.Audio.EventPos(s.Spells.DefByInd(spell.ID(id)).GetAudio(int(phase)), pos, 0, 0)
}
func spellEffectInform(u *server.Object) {
	if u.ObjClass&4 != 0 {
		v := C.int(2)
		C.nox_xxx_netInformTextMsg_4DA0F0(C.int(*controlByte(controlPlayer(u), 2064)), 0, &v)
	}
}
func spellEffectAlert(source, target *server.Object) {
	C.sub_4E7540(asObjectC(source), asObjectC(target))
}
func spellEffectGlyphType() uint32 {
	if C.dword_5d4594_2487712 == 0 {
		C.dword_5d4594_2487712 = C.uint32_t(GetServer().S().Types.IndByID("Glyph"))
	}
	return uint32(C.dword_5d4594_2487712)
}
func spellEffectCreate(u, owner *server.Object, pos types.Pointf) {
	C.nox_xxx_createAt_4DAA50(asObjectC(u), asObjectC(owner), C.float(pos.X), C.float(pos.Y))
}
func spellEffectNew(id uint32) *server.Object {
	return asObjectS(C.nox_xxx_newObjectWithTypeInd_4E3450(C.int(id)))
}
func spellEffectTrace(from, to types.Pointf, flags int) bool {
	return GetServer().S().MapTraceRayAt(from, to, nil, nil, server.MapTraceFlags(flags))
}
func spellEffectPlacement(from, to types.Pointf) bool {
	ray := [4]float32{from.X, from.Y, to.X, to.Y}
	return byte(C.nox_xxx_traceRay_5374B0((*C.float4)(unsafe.Pointer(&ray)))) != 0
}
func spellEffectSummonCost(id int32, u *server.Object) int32 {
	if u != nil && u.ObjClass&4 != 0 {
		return int32(*memmap.PtrUint32(0x587000, 217668+uintptr(id)*4))
	}
	return 0
}
func spellEffectRestoreHealth(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	resourceRestoreHP(u)
	GetServer().S().Audio.EventObj(754, u, 0, 0)
	return 1
}
func spellEffectRestoreMana(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 0 {
		resourceAddMana(u, int16(*controlHalf(u.UpdateData, 8)))
		GetServer().S().Audio.EventObj(755, u, 0, 0)
	}
	return 1
}
func spellEffectConfuse(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	spellLifeApplyBuff(u, 3, int16(floatToInt32(float32(spellEffectScalar("ConfuseEnchantDuration")))), int8(level))
	spellEffectAlert(b, u)
	return 1
}
func spellEffectStun(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	dur := int16(floatToInt32(float32(spellEffectScalar("StunEnchantDuration"))))
	buff := int32(5)
	if u.ObjClass&4 != 0 {
		if *controlByte(controlPlayer(u), 2251) == 0 {
			buff = 4
		}
	} else if u.ObjClass&2 != 0 && u.Mass > 15 {
		buff = 4
	}
	spellLifeApplyBuff(u, buff, dur, int8(level))
	spellEffectAlert(b, u)
	return 1
}
func spellEffectShock(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	typ := spellEffectGlyphType()
	if c != nil && uint32(c.TypeInd) == typ {
		projectileDamage(u, b, b, floatToInt32(float32(spellEffectTable("ShockTrapDamage", level-1))), 9)
	} else {
		spellLifeApplyBuff(u, 22, int16(floatToInt32(float32(spellEffectScalar("ShockEnchantDuration")))), int8(level))
	}
	return 1
}
func spellEffectPoison(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	resourcePoison(u, level, level)
	spellEffectAlert(b, u)
	return 1
}
func spellEffectCurePoison(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	if u.Poison540 != 0 {
		if int32(u.Poison540) > level {
			resourceReducePoison(u, level)
			resourcePriority(u, "ExecSpel.c:PoisonCure")
		} else {
			resourceRemovePoison(u)
			resourcePriority(u, "ExecSpel.c:PoisonClean")
		}
		spellEffectAudio(id, 1, u)
	} else if u != a {
		spellEffectAudio(id, 1, u)
	} else {
		spellLifeRefundMana(u, int16(GetServer().S().Spells.ManaCost(spell.ID(id), 1)))
	}
	return 1
}
func spellEffectTelekinesis(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 0 {
		if v := GetServer().S().NewObjectByTypeID("TelekinesisHand"); v != nil {
			spellEffectCreate(v, u, u.PosVec)
			spellLifeApplyBuff(u, 24, int16(uint16(GetServer().S().TickRate())*20), int8(level))
			GetServer().S().Spells.Dur.CancelFor(24, u)
			GetServer().S().Spells.Dur.CancelFor(43, u)
			spellEffectAudio(id, 1, u)
		}
	}
	return 1
}
func spellEffectLesserHeal(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	if resourceGetHP(u) == resourceGetMaxHP(u) && a == u {
		spellLifeRefundMana(b, int16(GetServer().S().Spells.ManaCost(spell.ID(id), 1)))
		return 1
	}
	v := float32(spellEffectScalar("LesserHealAmount"))
	if b != nil && b.ObjClass&4 != 0 {
		var max float64
		switch *controlByte(controlPlayer(b), 2251) {
		case 0:
			max = float64(GetServer().S().Players.Mult.Warrior.Health)
		case 1:
			max = float64(GetServer().S().Players.Mult.Wizard.Health)
		case 2:
			max = float64(GetServer().S().Players.Mult.Conjurer.Health)
		default:
			resourceAdjustHP(u, floatToInt32(v))
			spellEffectAudio(id, 1, u)
			return 1
		}
		v = float32(max * float64(v))
	}
	resourceAdjustHP(u, floatToInt32(v))
	spellEffectAudio(id, 1, u)
	return 1
}
func spellEffectFumble(id int32, a, b, c *server.Object, record unsafe.Pointer, level int32) int32 {
	u := spellEffectObject(record, 0)
	if u == nil {
		return 0
	}
	if u.ObjClass&4 != 0 || u.ObjClass&2 != 0 && u.ObjSubClass&0x10 != 0 {
		for it := u.InvFirstItem; it != nil; {
			next := it.InvNextItem
			if it.ObjFlags&0x100 != 0 && (it.ObjClass&0x1001000 != 0 || it.ObjClass&0x2000000 != 0 && it.ObjSubClass&2 != 0) {
				inventoryForceDrop(u, it)
			}
			it = next
		}
		typ := stateType(2487728, "GameBall")
		for it := u.Field129; it != nil; it = it.Field128 {
			if uint32(it.TypeInd) == typ {
				C.nox_xxx_objectApplyForce_52DF80((*C.float)(unsafe.Pointer(&u.PosVec)), asObjectC(it), 100)
				it.SetOwner(nil)
				GetServer().S().Audio.EventObj(926, u, 0, 0)
				break
			}
		}
	} else if u.ObjClass&2 == 0 || u.ObjSubClass&0x2000 == 0 {
		inventoryDropAll(u)
		C.nox_xxx_objectApplyForce_52DF80((*C.float)(unsafe.Pointer(&c.PosVec)), asObjectC(u), 50)
	}
	spellEffectAudio(id, 1, u)
	return 1
}
