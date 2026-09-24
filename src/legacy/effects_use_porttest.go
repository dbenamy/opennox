//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
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
	ExpectedClock                                *[2]uint32 // Optional assertion after all shared setup has run.
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
		*memmap.PtrPtr(0x587000, 200160+uintptr(20*i)) = portTestEffectsFunction(id - 499)
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
	p.reservedFunctionIDs += 1
	if clock := sp.ExpectedClock; clock != nil && (p.proxy.core.Frame() != clock[0] || uint32(p.proxy.core.TickRate()) != clock[1]) {
		panic("effects fixture clock was overwritten")
	}
	if p.equipment == nil {
		panic("effects fixture requires equipment")
	}

	for i := 1; i <= 43; i++ {
		if fn := portTestEffectsFunction(i); fn != nil {
			p.identify(fn, 66100+uint32(i))
		}
	}
	for i, s := range sp.Modifiers {
		m := (*server.ModifierEff)(unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*144))
		for off, v := range s.Words {
			if off < 8 || off >= 136 || off%4 != 0 {
				panic("effects descriptor word offset")
			}
			*(*uint32)(unsafe.Add(unsafe.Pointer(m), off)) = v
		}
		m.Engage112 = portTestEffectsFunction(s.Engage)
		m.Disengage116 = portTestEffectsFunction(s.Disengage)
		m.Attack40.Fnc = portTestEffectsFunction(s.Attack)
		m.Defend76.Fnc = portTestEffectsFunction(s.Defend)
		m.DefendCollide88.Fnc = portTestEffectsFunction(s.Collide)
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
	if a.Op == PortTestEffects53F8E0 {
		p.effectsUse.result = uint64(uint32(effectsUse(u, it)))
	} else if a.Op == PortTestEffects53F290 {
		p.effectsUse.result = uint64(uint32(effectsLesserFireball(u, it)))
	} else if a.Op == PortTestEffects53F4F0 {
		p.effectsUse.result = uint64(uint32(effectsWandCast(u, it)))
	} else if a.Op == PortTestEffects53F670 {
		p.effectsUse.result = uint64(uint32(effectsFireWand(u, it)))
	} else {
		p.effectsUse.result = portTestEffectsCall(a.Op-500, u, it, p.effectsUseTarget(sp.Target), (*server.ModifierEff)(mod), p.effectsUse.scalar, int(int32(a.Value)), a.Side, p.inventory.pos)
	}
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
		if p.spellEffectsActive() && u.HealthData != nil {
			p.identify(unsafe.Pointer(u.HealthData), 96000+uint32(i))
		}
		appendWords(u.CObj(), 193)
		if p.spellEffectsActive() && u.HealthData != nil {
			appendWords(unsafe.Pointer(u.HealthData), 5)
		}
		if u.UpdateData != nil {
			size := 80
			if p.spellEffectsActive() {
				if typ := p.proxy.core.Types.ByInd(int(u.TypeInd)); typ != nil && typ.UpdateDataSize >= 16 {
					size = int(typ.UpdateDataSize)
				}
			}
			for _, v := range unsafe.Slice((*byte)(u.UpdateData), size)[size-16:] {
				if v != 0xa5 {
					panic("effects Spark update guard")
				}
			}
			appendWords(u.UpdateData, size/4)
		}
		if u.CollideData != nil {
			size := 20
			if typ := p.proxy.core.Types.ByInd(int(u.TypeInd)); typ != nil && typ.CollideDataSize >= 16 {
				size = int(typ.CollideDataSize)
			}
			for _, v := range unsafe.Slice((*byte)(u.CollideData), size)[size-16:] {
				if v != 0x5a {
					panic("effects Spark collide guard")
				}
			}
			appendWords(u.CollideData, size/4)
		}
	}
	if p.proxy.callbacks.shop.spec.EffectsUse.ExpectedClock != nil {
		out = append(out, p.proxy.core.Frame(), uint32(p.proxy.core.TickRate()))
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

func portTestEffectsFunction(id int) unsafe.Pointer {
	switch id {
	case 37:
		return itemIdentityKey(itemIDWandUse)
	case 39:
		return itemIdentityKey(itemIDWandCastUse)
	case 40:
		return itemIdentityKey(itemIDFireWandUse)
	}
	return modifierFunctionKey(id)
}

func modifierFunctionKey(id int) unsafe.Pointer {
	switch id {
	case 1:
		return modifierKey(modifierIDBrillianceEngage)
	case 2:
		return modifierKey(modifierIDBrillianceDisengage)
	case 3:
		return modifierTestKey(2)
	case 4:
		return modifierKey(modifierIDSpeedEngage)
	case 5:
		return modifierKey(modifierIDSpeedDisengage)
	case 6:
		return modifierKey(modifierIDFireProtectEngage)
	case 7:
		return modifierKey(modifierIDFireProtectDisengage)
	case 8:
		return modifierKey(modifierIDLightningProtectEngage)
	case 9:
		return modifierKey(modifierIDLightningProtectDisengage)
	case 10:
		return modifierKey(modifierIDPoisonProtectEngage)
	case 11:
		return modifierKey(modifierIDPoisonProtectDisengage)
	case 12:
		return modifierTestKey(1)
	case 13:
		return modifierTestKey(0)
	case 14:
		return modifierTestKey(3)
	case 15:
		return modifierKey(modifierIDRegenerationEngage)
	case 16:
		return modifierKey(modifierIDRegenerationDisengage)
	case 17:
		return modifierKey(modifierIDRegenerationUpdate)
	case 18:
		return modifierKey(modifierIDContinualReplenishmentUpdate)
	case 19:
		return modifierKey(modifierIDArmorMultiplierEffect)
	case 20:
		return modifierKey(modifierIDDurabilityMultiplierEffect)
	case 21:
		return modifierKey(modifierIDInversionEffect)
	case 22:
		return modifierTestKey(7)
	case 23:
		return modifierKey(modifierIDGripEffect)
	case 24:
		return modifierKey(modifierIDDamageMultiplierEffect)
	case 25:
		return modifierKey(modifierIDStunEffect)
	case 26:
		return modifierKey(modifierIDRecoilEffect)
	case 27:
		return modifierKey(modifierIDConfuseEffect)
	case 28:
		return modifierKey(modifierIDLightningEffect)
	case 29:
		return modifierKey(modifierIDDrainManaEffect)
	case 30:
		return modifierKey(modifierIDVampirismEffect)
	case 31:
		return modifierKey(modifierIDPoisonEffect)
	case 32:
		return modifierKey(modifierIDSympathyEffect)
	case 33:
		return modifierTestKey(5)
	case 34:
		return modifierKey(modifierIDProjectileSpeedEffect)
	case 35:
		return modifierTestKey(6)
	case 36:
		return modifierTestKey(4)
	case 38:
		return modifierTestKey(8)
	case 42:
		return modifierKey(modifierIDReadinessEffect)
	case 43:
		return modifierKey(modifierIDReplenishmentEffect)
	default:
		return nil
	}
}

func portTestEffectsCall(op int, u, it, target *server.Object, m *server.ModifierEff, scalar *uint32, value, side int, pos *types.Pointf) uint64 {
	var data unsafe.Pointer
	if scalar != nil {
		data = unsafe.Pointer(scalar)
	}
	switch op {
	case 0:
		effectsEngageFlag(u, 8, 75)
	case 1:
		effectsDisengageFlag(u, 8, 76)
	case 2:
		return uint64(uint32(effectsInventory(u, byte(value))))
	case 3:
		effectsSpeed(m, u, true)
	case 4:
		effectsSpeed(m, u, false)
	case 5:
		effectsEngageFlag(u, 1, 102)
	case 6:
		if u != nil && it != nil {
			effectsDisengageFlag(u, 1, 103)
		}
	case 7:
		effectsEngageFlag(u, 4, 106)
	case 8:
		effectsDisengageFlag(u, 4, 107)
	case 9:
		effectsEngageFlag(u, 2, 110)
	case 10:
		effectsDisengageFlag(u, 2, 111)
	case 11:
		return math.Float64bits(effectsProtection(u, modifierKey(modifierIDFireProtectEngage), 17, "FireSpellProtection", .5, .60000002))
	case 12:
		return math.Float64bits(effectsProtection(u, modifierKey(modifierIDLightningProtectEngage), 20, "ElectricitySpellProtection", .5, .60000002))
	case 13:
		return math.Float64bits(effectsProtection(u, modifierKey(modifierIDPoisonProtectEngage), 18, "PoisonSpellProtection", .69999999, .89999998))
	case 14:
		effectsEngageFlag(u, 32, 123)
	case 15:
		if u != nil && u.ObjClass&4 != 0 {
			effectsDisengageFlag(u, 32, 124)
		}
	case 16:
		effectsRegeneration(m, it)
	case 17:
		effectsReplenish(m, it)
	case 18, 19:
		v := (*float32)(data)
		if op == 18 {
			*v = float32(float64(m.Defend76.Valf) * float64(*v))
		} else {
			*v = float32((1 - float64(m.Defend76.Valf) + 1) * float64(*v))
		}
		return uint64(uint32(uintptr(data)))
	case 20:
		return uint64(uint32(effectsGrip(m, (*int32)(data), true)))
	case 21:
		return uint64(uint32(effectsGripSearch(u, u, it, target)))
	case 22:
		return uint64(uint32(effectsGrip(m, (*int32)(data), false)))
	case 23:
		v := (*float32)(data)
		*v = float32(float64(m.Attack40.Valf) * float64(*v))
		return uint64(uint32(uintptr(data)))
	case 24:
		effectsStatus(m, u, target, true)
	case 25:
		effectsRecoil(m, it, target)
	case 26:
		effectsStatus(m, u, target, false)
	case 27:
		effectsLightning(m, it, u, target)
	case 28:
		modifierDrainNative(m, u, target, data)
	case 29:
		modifierVampirismNative(m, u, target, data)
	case 30:
		effectsPoison(m, u, target)
	case 31:
		modifierSympathyNative(m, u, target, data)
	case 32:
		return uint64(uint32(effectsReadiness(it)))
	case 33:
		target.SpeedCur = float32(float64(m.Attack40.Valf) * float64(target.SpeedCur))
		return uint64(uint32(uintptr(unsafe.Pointer(target))))
	case 34:
		return uint64(uint32(effectsRecharge(it, int32(value))))
	case 35:
		return uint64(uint32(effectsRechargeRate(it)))
	case 37:
		return uint64(effectsWandShot(u, value, *pos, uint32(side)))
	}
	return 0
}
