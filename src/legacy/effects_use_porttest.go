//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME3_2.h"
#include "GAME4_3.h"
void nullsub_22(void);
void nullsub_36(void);
void sub_4DFB80(int a1, int a2);
void nox_xxx_effectSpeedDisengage_4DFCA0(int a1, int a2);
void nox_xxx_modifFireProtection_4DFD40(int a1, int a2, int a3);
void sub_4DFDB0(int a1, int a2);
void sub_4DFE10(int a1, int a2);
void nox_xxx_effectRegeneration_4E01D0(int a1, int a2);
void nox_xxx_attribContinualReplen_4E02C0(int a1, uint32_t* a2);
int nox_xxx_unusedCheckGripEffect_4E03F0(int a1, int a2, int a3, int a4);
void nox_xxx_stunEffect_4E04D0(int a1, int a2, int a3, int a4);
void nox_xxx_confuseEffect_4E0670(int a1, int a2, int a3, int a4);
void nox_xxx_drainMEffect_4E0740(int a1, int a2, int a3, int a4, int* a5);
void nox_xxx_vampirismEffect_4E07C0(int a1, int a2, int a3, int a4, int* a5);
void nox_xxx_poisonEffect_4E0850(int a1, int a2, int a3, int a4);
void nox_xxx_sympathyEffect_4E08E0(int a1, int a2, int a3, int a4, int* a5);
static void* effectsFunction(int id){switch(id){
case 1:return (void*)sub_4DFB50;
case 2:return (void*)sub_4DFB80;
case 3:return (void*)nox_xxx_enchantItemTestInventory_4DFBB0;
case 4:return (void*)nox_xxx_effectSpeedEngage_4DFC30;
case 5:return (void*)nox_xxx_effectSpeedDisengage_4DFCA0;
case 6:return (void*)sub_4DFD10;
case 7:return (void*)nox_xxx_modifFireProtection_4DFD40;
case 8:return (void*)nox_xxx_buff_4DFD80;
case 9:return (void*)sub_4DFDB0;
case 10:return (void*)nox_xxx_checkPoisonProtectEnch_4DFDE0;
case 11:return (void*)sub_4DFE10;
case 12:return (void*)nox_xxx_checkFireProtect_4DFE40;
case 13:return (void*)nox_xxx_checkElectrProtect_4DFF40;
case 14:return (void*)nox_xxx_getPoisonDmg_4E0040;
case 15:return (void*)sub_4E0140;
case 16:return (void*)sub_4E0170;
case 17:return (void*)nox_xxx_effectRegeneration_4E01D0;
case 18:return (void*)nox_xxx_attribContinualReplen_4E02C0;
case 19:return (void*)sub_4E0370;
case 20:return (void*)sub_4E0380;
case 21:return (void*)nox_xxx_inversionEffect_4E03D0;
case 22:return (void*)nox_xxx_unusedCheckGripEffect_4E03F0;
case 23:return (void*)nox_xxx_gripEffect_4E0480;
case 24:return (void*)nox_xxx_effectDamageMultiplier_4E04C0;
case 25:return (void*)nox_xxx_stunEffect_4E04D0;
case 26:return (void*)nox_xxx_recoilEffect_4E0640;
case 27:return (void*)nox_xxx_confuseEffect_4E0670;
case 28:return (void*)nox_xxx_lightngEffect_4E06F0;
case 29:return (void*)nox_xxx_drainMEffect_4E0740;
case 30:return (void*)nox_xxx_vampirismEffect_4E07C0;
case 31:return (void*)nox_xxx_poisonEffect_4E0850;
case 32:return (void*)nox_xxx_sympathyEffect_4E08E0;
case 33:return (void*)nox_xxx_itemCheckReadinessEffect_4E0960;
case 34:return (void*)nox_xxx_effectProjectileSpeed_4E09B0;
case 35:return (void*)nox_xxx_rechargeItem_53C520;
case 36:return (void*)nox_xxx_getRechargeRate_53C940;
case 37:return (void*)nox_xxx_useLesserFireballStaff_53F290;
case 38:return (void*)nox_xxx_wandShot_53F480;
case 39:return (void*)nox_xxx_useWandCastSpell_53F4F0;
case 40:return (void*)nox_xxx_useFireWand_53F670;
case 41:return (void*)nox_xxx_useByNetCode_53F8E0;
case 42:return (void*)nullsub_22;case 43:return (void*)nullsub_36;default:return 0;}}
static uint64_t effectsCall(int op,nox_object_t* u,nox_object_t* it,nox_object_t* target,void* mod,uint32_t* scalar,int value,int side,float2* pos){
 int up=(int)u,ip=(int)it,tp=(int)target,mp=(int)mod;
 switch(op){
case 0:sub_4DFB50(mp,up);return 0;
case 1:sub_4DFB80(mp,up);return 0;
case 2:return (uint32_t)nox_xxx_enchantItemTestInventory_4DFBB0(up,(char)value);
case 3:nox_xxx_effectSpeedEngage_4DFC30(mp,up);return 0;
case 4:nox_xxx_effectSpeedDisengage_4DFCA0(mp,up);return 0;
case 5:sub_4DFD10(mp,up);return 0;
case 6:nox_xxx_modifFireProtection_4DFD40(mp,up,ip);return 0;
case 7:nox_xxx_buff_4DFD80(mp,up);return 0;
case 8:sub_4DFDB0(mp,up);return 0;
case 9:nox_xxx_checkPoisonProtectEnch_4DFDE0(mp,up);return 0;
case 10:sub_4DFE10(mp,up);return 0;
case 11:{double d=nox_xxx_checkFireProtect_4DFE40((uint32_t*)u);uint64_t b;memcpy(&b,&d,8);return b;}
case 12:{double d=nox_xxx_checkElectrProtect_4DFF40((uint32_t*)u);uint64_t b;memcpy(&b,&d,8);return b;}
case 13:{double d=nox_xxx_getPoisonDmg_4E0040((uint32_t*)u);uint64_t b;memcpy(&b,&d,8);return b;}
case 14:sub_4E0140(mp,up);return 0;
case 15:sub_4E0170(mp,up);return 0;
case 16:nox_xxx_effectRegeneration_4E01D0(mp,ip);return 0;
case 17:nox_xxx_attribContinualReplen_4E02C0(mp,(uint32_t*)it);return 0;
case 18:return (uint32_t)sub_4E0370(mp,ip,up,tp,up,(float*)scalar);
case 19:return (uint32_t)sub_4E0380(mp,ip,up,tp,up,(float*)scalar);
case 20:return (uint32_t)nox_xxx_inversionEffect_4E03D0(mp,ip,up,tp,up,(int*)scalar);
case 21:return (uint32_t)nox_xxx_unusedCheckGripEffect_4E03F0(up,up,ip,tp);
case 22:return (uint32_t)nox_xxx_gripEffect_4E0480(mp,ip,up,tp,up,(int*)scalar);
case 23:return (uint32_t)nox_xxx_effectDamageMultiplier_4E04C0(mp,ip,up,tp,(float*)scalar);
case 24:nox_xxx_stunEffect_4E04D0(mp,ip,up,tp);return 0;
case 25:nox_xxx_recoilEffect_4E0640(mp,ip,up,tp);return 0;
case 26:nox_xxx_confuseEffect_4E0670(mp,ip,up,tp);return 0;
case 27:nox_xxx_lightngEffect_4E06F0(mp,ip,up,tp);return 0;
case 28:nox_xxx_drainMEffect_4E0740(mp,ip,up,tp,(int*)scalar);return 0;
case 29:nox_xxx_vampirismEffect_4E07C0(mp,ip,up,tp,(int*)scalar);return 0;
case 30:nox_xxx_poisonEffect_4E0850(mp,ip,up,tp);return 0;
case 31:nox_xxx_sympathyEffect_4E08E0(mp,ip,up,tp,(int*)scalar);return 0;
case 32:return (uint32_t)nox_xxx_itemCheckReadinessEffect_4E0960(ip);
case 33:return (uint32_t)nox_xxx_effectProjectileSpeed_4E09B0(mp,ip,up,tp,tp);
case 34:return (uint32_t)nox_xxx_rechargeItem_53C520(ip,value);
case 35:return (uint32_t)nox_xxx_getRechargeRate_53C940((uint32_t*)it);
case 36:return (uint32_t)nox_xxx_useLesserFireballStaff_53F290(up,(uint32_t*)it);
case 37:return (uint32_t)nox_xxx_wandShot_53F480(up,value,(int*)pos,(uint32_t*)side);
case 38:return (uint32_t)nox_xxx_useWandCastSpell_53F4F0(up,(uint32_t*)it);
case 39:return (uint32_t)nox_xxx_useFireWand_53F670(up,ip);
case 40:return (uint32_t)nox_xxx_useByNetCode_53F8E0(up,ip);
}
return 0;
}
*/
import "C"

import (
	"bytes"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestEffectsModifier struct {
	Engage, Disengage, Attack, Defend, Collide int            // One-based function identities; zero is nil.
	Words                                      map[int]uint32 // Byte offsets into the actual 144-byte descriptor.
}
type PortTestEffectsUseSpec struct {
	Balance                                      map[string][]float64
	Projectiles, DisableProjectiles, AcceptSpell bool
	ProjectileSpeed                              uint32
	BuffPower                                    byte
	PlayerWords                                  map[int]uint32
	CastTarget                                   bool
	UnitUpdateWords, TargetUpdateWords           map[int]uint32

	Modifiers                 [4]PortTestEffectsModifier
	Modifier, Target          int
	NilUnit, UnitItem, NilUse bool
	Scalar                    uint32
	UnitWords, TargetWords    map[int]uint32
	ItemWords                 []map[int]uint32
}
type portTestEffectsUse struct {
	projectiles [2]uint16
	accepts     []uint32

	result uint64
	scalar *uint32
}

func (p *portTestShopPools) effectsUsePrepare() func() {
	p.effectsUse = nil
	sp := p.proxy.callbacks.shop.spec.EffectsUse
	if sp == nil {
		return func() {}
	}
	ids, restore := p.proxy.core.PortTestEffectsUseEnvironment(sp.Balance, sp.Projectiles, sp.DisableProjectiles, math.Float32frombits(sp.ProjectileSpeed))
	p.effectsUse = &portTestEffectsUse{projectiles: ids, scalar: (*uint32)(p.equipmentRegion(16, 66001))}
	*p.effectsUse.scalar = sp.Scalar
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 200160)), 120)
	oldTable := bytes.Clone(table)
	copy(table, blobdata.PortTestEffectsInventoryTable())
	for i, id := range []int{PortTestEffects4DFB50, PortTestEffects4DFC30, PortTestEffects4DFD10, PortTestEffects4DFD80, PortTestEffects4DFDE0, PortTestEffects4E0140} {
		*memmap.PtrPtr(0x587000, 200160+uintptr(20*i)) = C.effectsFunction(C.int(id - 499))
	}
	offsets := []uintptr{2488732, 1569740, 1569744}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = *memmap.PtrUint32(0x5d4594, off)
		*memmap.PtrUint32(0x5d4594, off) = 0
	}
	return func() {
		copy(table, oldTable)
		restore()
		for i, off := range offsets {
			*memmap.PtrUint32(0x5d4594, off) = old[i]
		}
	}

}
func (s *portTestRoamOwnerServer) Nox_xxx_spellAccept4FD400(id spell.ID, a, b, c *server.Object, arg *server.SpellAcceptArg, level int) bool {
	if s.callbacks != nil && s.callbacks.shop != nil && s.callbacks.shop.spec.EffectsUse != nil {
		p := s.callbacks.shop.pools
		norm := func(u *server.Object) uint32 { return p.normalize(uint32(uintptr(u.CObj()))) }
		p.effectsUse.accepts = append(p.effectsUse.accepts, uint32(id), norm(a), norm(b), norm(c), norm(arg.Obj), math.Float32bits(arg.Pos.X), math.Float32bits(arg.Pos.Y), uint32(level))
		return s.callbacks.shop.spec.EffectsUse.AcceptSpell
	}
	return s.portTestRandomServer.Server.Nox_xxx_spellAccept4FD400(id, a, b, c, arg, level)
}
func (p *portTestShopPools) effectsUseItems() {
	sp := p.proxy.callbacks.shop.spec.EffectsUse
	if sp == nil {
		return
	}
	if p.equipment == nil {
		panic("effects fixture requires equipment")
	}

	for i := 1; i <= 43; i++ {
		p.identify(C.effectsFunction(C.int(i)), 66100+uint32(i))
	}
	for i, s := range sp.Modifiers {
		m := (*server.ModifierEff)(unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*144))
		for off, v := range s.Words {
			if off < 8 || off >= 136 || off%4 != 0 {
				panic("effects descriptor word offset")
			}
			*(*uint32)(unsafe.Add(unsafe.Pointer(m), off)) = v
		}
		m.Engage112 = C.effectsFunction(C.int(s.Engage))
		m.Disengage116 = C.effectsFunction(C.int(s.Disengage))
		m.Attack40.Fnc = C.effectsFunction(C.int(s.Attack))
		m.Defend76.Fnc = C.effectsFunction(C.int(s.Defend))
		m.DefendCollide88.Fnc = C.effectsFunction(C.int(s.Collide))
	}
	write := func(u *server.Object, words map[int]uint32) {
		for off, v := range words {
			if off < 0 || off >= 772 || off%4 != 0 {
				panic("effects object word offset")
			}
			*(*uint32)(unsafe.Add(u.CObj(), off)) = v
		}
	}
	p.resources.unit.Damage = p.proxy.combat.target.Damage
	for i := range p.resources.unit.BuffsPower {
		p.resources.unit.BuffsPower[i] = sp.BuffPower
	}
	writeUD := func(u *server.Object, words map[int]uint32) {
		for off, v := range words {
			if off < 0 || off >= 2200 || off%4 != 0 {
				panic("effects update word offset")
			}
			*(*uint32)(unsafe.Add(u.UpdateData, off)) = v
		}
	}
	writeUD(p.resources.unit, sp.UnitUpdateWords)
	writeUD(p.effectsUseTarget(sp.Target), sp.TargetUpdateWords)
	if sp.Projectiles {
		copy(unsafe.Slice((*byte)(unsafe.Add(p.items[0].u.UseData.Ptr, 4)), 80), "EffectsBolt\x00")
	}
	if p.resources.unit.ObjClass&4 != 0 {
		pl := unsafe.Pointer(p.resources.unit.UpdateDataPlayer().Player)
		for off, v := range sp.PlayerWords {
			if off < 0 || off >= 5000 || off%4 != 0 {
				panic("effects player word offset")
			}
			*(*uint32)(unsafe.Add(pl, off)) = v
		}
		if sp.CastTarget {
			*(*uintptr)(unsafe.Add(p.resources.unit.UpdateData, 288)) = uintptr(p.effectsUseTarget(sp.Target).CObj())
		}
	}
	write(p.resources.unit, sp.UnitWords)
	write(p.effectsUseTarget(sp.Target), sp.TargetWords)
	for i, w := range sp.ItemWords {
		write(p.items[i].u, w)
	}
	if sp.NilUse {
		for _, it := range p.items {
			it.u.Use.Ptr = nil
		}
	}
}
func (p *portTestShopPools) effectsUseTarget(index int) *server.Object {
	switch index {
	case 0:
		return p.proxy.combat.target
	case 1:
		return p.resources.unit
	case 2:
		return p.items[0].u
	case 3:
		return p.items[1].u
	case 4:
		return nil
	default:
		panic("effects target index")
	}
}
func (p *portTestShopPools) effectsUseAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec.EffectsUse
	u := p.resources.unit
	var it *server.Object
	if !p.proxy.callbacks.shop.spec.Inventory.NilItem && a.Item >= 0 {
		it = p.items[a.Item].u
	}
	if sp.UnitItem {
		u = it
	}
	if sp.NilUnit {
		u = nil
	}
	mod := unsafe.Add(p.proxy.callbacks.shop.ptr(8), sp.Modifier*144)
	if a.Op == PortTestEffects53F480 && a.Value == 0xfffffffe {
		a.Value = uint32(p.effectsUse.projectiles[0])
	}
	p.effectsUse.result = uint64(C.effectsCall(C.int(a.Op-500), asObjectC(u), asObjectC(it), asObjectC(p.effectsUseTarget(sp.Target)), mod, (*C.uint32_t)(unsafe.Pointer(p.effectsUse.scalar)), C.int(a.Value), C.int(a.Side), (*C.float2)(unsafe.Pointer(p.inventory.pos))))
	return uint32(p.effectsUse.result)
}
func (p *portTestShopPools) effectsUseSnapshot() []uint32 {
	if p.effectsUse == nil {
		return nil
	}
	out := []uint32{p.normalize(uint32(p.effectsUse.result)), uint32(p.effectsUse.result >> 32), *p.effectsUse.scalar, *memmap.PtrUint32(0x5d4594, 2488732)}
	out = append(out, *memmap.PtrUint32(0x5d4594, 1569740), *memmap.PtrUint32(0x5d4594, 1569744))
	for _, v := range unsafe.Slice(memmap.PtrUint32(0x587000, 200160), 30) {
		out = append(out, p.normalize(v))
	}
	out = append(out, uint32(len(p.effectsUse.accepts)))
	out = append(out, p.effectsUse.accepts...)
	appendWords := func(ptr unsafe.Pointer, n int) {
		for _, v := range unsafe.Slice((*uint32)(ptr), n) {
			out = append(out, p.normalize(v))
		}
	}
	target := p.effectsUseTarget(p.proxy.callbacks.shop.spec.EffectsUse.Target)
	if target != nil {
		out = append(out, 1)
		appendWords(target.CObj(), 193)
		if target.HealthData != nil {
			appendWords(unsafe.Pointer(target.HealthData), 5)
		}
	} else {
		out = append(out, 0)
	}
	for i, u := range p.proxy.life.created {
		if u.CollideData != nil {
			p.identify(u.CollideData, 67000+uint32(i))
		}
		appendWords(u.CObj(), 193)
		if u.UpdateData != nil {
			for _, v := range unsafe.Slice((*byte)(u.UpdateData), 80)[64:] {
				if v != 0xa5 {
					panic("effects Spark update guard")
				}
			}
			appendWords(u.UpdateData, 20)
		}
		if u.CollideData != nil {
			for _, v := range unsafe.Slice((*byte)(u.CollideData), 20)[4:] {
				if v != 0x5a {
					panic("effects Spark collide guard")
				}
			}
			appendWords(u.CollideData, 5)
		}
	}
	return out
}

const PortTestEffects4DFB50 = 500
const PortTestEffects4DFB80 = 501
const PortTestEffects4DFBB0 = 502
const PortTestEffects4DFC30 = 503
const PortTestEffects4DFCA0 = 504
const PortTestEffects4DFD10 = 505
const PortTestEffects4DFD40 = 506
const PortTestEffects4DFD80 = 507
const PortTestEffects4DFDB0 = 508
const PortTestEffects4DFDE0 = 509
const PortTestEffects4DFE10 = 510
const PortTestEffects4DFE40 = 511
const PortTestEffects4DFF40 = 512
const PortTestEffects4E0040 = 513
const PortTestEffects4E0140 = 514
const PortTestEffects4E0170 = 515
const PortTestEffects4E01D0 = 516
const PortTestEffects4E02C0 = 517
const PortTestEffects4E0370 = 518
const PortTestEffects4E0380 = 519
const PortTestEffects4E03D0 = 520
const PortTestEffects4E03F0 = 521
const PortTestEffects4E0480 = 522
const PortTestEffects4E04C0 = 523
const PortTestEffects4E04D0 = 524
const PortTestEffects4E0640 = 525
const PortTestEffects4E0670 = 526
const PortTestEffects4E06F0 = 527
const PortTestEffects4E0740 = 528
const PortTestEffects4E07C0 = 529
const PortTestEffects4E0850 = 530
const PortTestEffects4E08E0 = 531
const PortTestEffects4E0960 = 532
const PortTestEffects4E09B0 = 533
const PortTestEffects53C520 = 534
const PortTestEffects53C940 = 535
const PortTestEffects53F290 = 536
const PortTestEffects53F480 = 537
const PortTestEffects53F4F0 = 538
const PortTestEffects53F670 = 539
const PortTestEffects53F8E0 = 540
