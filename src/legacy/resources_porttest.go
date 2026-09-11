//go:build porttest

package legacy

/*
#include "GAME3_3.h"
#include "GAME4.h"
int nox_xxx_pickupGold_4F3A60_obj_pickup(int unit, int item, int flags);
static uint32_t resourceDieCount, resourceDieUnit;
static void resourceDie(void* unit) {resourceDieCount++;resourceDieUnit=(uint32_t)unit;}
static void* resourceDiePtr(void) {return resourceDie;}
static void resourceDieReset(void) {resourceDieCount=0;resourceDieUnit=0;}
static uint32_t resourceDieCalls(void) {return resourceDieCount;}
static uint32_t resourceDieWho(void) {return resourceDieUnit;}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestResourceSpec struct {
	ExtraClass, SyncSeed                                 uint32
	OtherHolder, OtherOwner                              bool
	PickupResult                                         bool
	Buffs                                                uint32
	Harpoon                                              bool
	PlayerClass                                          byte
	Subject                                              int // 0 nil, 1 player, 2 simple object, 3 monster
	HP, OldHP, MaxHP                                     uint16
	Mana, OldMana, MaxMana                               uint16
	NoHealth, God, Owner, Holder, Protected, DieCallback bool
	Flags, Subclass, GoldItem, PoisonTime, Status        uint32
	Poison                                               byte
	PoisonTimer                                          uint16
}

const (
	PortTestResourceSetHP = 200 + iota
	PortTestResourceAdjustHP
	PortTestResourceInformOwner
	PortTestResourceDamage
	PortTestResourceRestoreHP
	PortTestResourceHPHistory
	PortTestResourceGetHP
	PortTestResourceGetMaxHP
	PortTestResourceSetMaxHP
	PortTestResourcePoison
	PortTestResourcePoisonReduce
	PortTestResourcePoisonRemove
	PortTestResourcePoisonSet
	PortTestResourceManaAdd
	PortTestResourceManaSub
	PortTestResourceGetMana
	PortTestResourceGetMaxMana
	PortTestResourceSetMaxMana
	PortTestResourceManaRefresh
	PortTestResourceGoldAdd
	PortTestResourceGoldSub
	PortTestResourceGoldSet
	PortTestResourceGetGold
	PortTestResourceObjectGold
	PortTestResourceGoldPickup
)

type portTestResources struct {
	pickupCalls  []uint32
	harpoonCalls []uint32
	unit         *server.Object
	blocks       [][]byte
	protection   func() []uint32
}

func (p *portTestShopPools) resourcePrepare() func() {
	sp := p.proxy.callbacks.shop.spec.Resources
	if sp == nil {
		return func() {}
	}
	r := &portTestResources{}
	p.resources = r
	oldPickup := Nox_xxx_pickupDefault_4F31E0
	Nox_xxx_pickupDefault_4F31E0 = func(u, item *server.Object, a3, a4 int) bool {
		r.pickupCalls = append(r.pickupCalls, p.normalize(uint32(uintptr(u.CObj()))), p.normalize(uint32(uintptr(item.CObj()))), uint32(a3), uint32(a4))
		return sp.PickupResult
	}
	oldHarpoon := Nox_xxx_harpoonBreakForPlr_537520
	Nox_xxx_harpoonBreakForPlr_537520 = func(u *server.Object) {
		r.harpoonCalls = append(r.harpoonCalls, p.normalize(uint32(uintptr(u.CObj()))))
	}
	oldEngine := noxflags.GetEngine()
	noxflags.ResetEngine()
	if sp.God {
		noxflags.SetEngine(noxflags.EngineGodMode)
	}
	C.resourceDieReset()
	var frees []func()
	region := func(size int, id uint32) unsafe.Pointer {
		b, free := alloc.Make([]byte{}, size+16)
		frees = append(frees, free)
		for i := 0; i < 8; i++ {
			b[i] = 0xa5
			b[len(b)-8+i] = 0x5a
		}
		r.blocks = append(r.blocks, b)
		ptr := unsafe.Pointer(&b[8])
		p.identify(ptr, id)
		return ptr
	}
	hp := (*server.HealthData)(region(int(unsafe.Sizeof(server.HealthData{})), 54001))
	hp.Cur, hp.Field2, hp.Max, hp.Field16 = sp.HP, sp.OldHP, sp.MaxHP, sp.PoisonTime
	u := p.proxy.callbacks.shop.npc()
	if sp.Subject == 1 {
		u = &p.proxy.life.players[0]
	}
	oldUnit := *u
	var savedUD []byte
	if sp.Subject == 1 {
		savedUD = bytes.Clone(unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.PlayerUpdateData{}))))
	} else {
		ud := region(2200, 54002)
		copy(unsafe.Slice((*byte)(ud), 2200), unsafe.Slice((*byte)(p.proxy.combat.actor.UpdateData), 2200))
		u.UpdateData = ud
		// These shared records are real retained dependencies, not fixture copies.
		// Their addresses vary with ASLR; preserve identity in full UD snapshots.
		p.identify(*(*unsafe.Pointer)(unsafe.Add(ud, 484)), 54004)
		p.identify(*(*unsafe.Pointer)(unsafe.Add(ud, 488)), 54005)
		u.ObjClass = object.ClassSimple
		if sp.Subject == 3 {
			u.ObjClass = object.ClassMonster
		}
		u.TypeInd = 23
		u.IDPtr = nil
		u.InitData = nil
	}
	u.ObjClass |= object.Class(sp.ExtraClass)
	if sp.SyncSeed != 0 {
		for i := 0; i < 32; i++ {
			*(*uint32)(unsafe.Add(u.CObj(), 560+4*i)) = sp.SyncSeed + uint32(i*12345)
		}
	}
	u.ObjFlags = object.Flags(sp.Flags)
	u.ObjSubClass = object.SubClass(sp.Subclass)
	u.HealthData = hp
	if sp.NoHealth {
		u.HealthData = nil
	}
	u.Buffs = sp.Buffs
	u.BuffsDur[18], u.BuffsPower[18] = 1, 1
	u.Poison540 = sp.Poison
	*(*uint16)(unsafe.Add(u.CObj(), 542)) = sp.PoisonTimer
	u.InvFirstItem, u.InvNextItem, u.Field125, u.InvHolder = nil, nil, nil, nil
	// Inputs do not award experience through a retained damage-source chain.
	*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 520)) = nil
	*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 508)) = nil
	*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 724)) = nil
	if sp.DieCallback {
		ptr := C.resourceDiePtr()
		*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 724)) = ptr
		p.identify(ptr, 54003)
	}
	if sp.Owner {
		owner := &p.proxy.life.players[1]
		*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 508)) = owner.CObj()
		*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 512)) = nil
		*(*unsafe.Pointer)(unsafe.Add(owner.CObj(), 516)) = u.CObj()
	}
	if sp.OtherOwner {
		*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 508)) = p.proxy.callbacks.shop.item().CObj()
	}
	if sp.OtherHolder {
		u.InvHolder = p.proxy.callbacks.shop.item()
	}
	if sp.Holder {
		u.InvHolder = &p.proxy.life.players[1]
	}
	if sp.Subject == 1 {
		*(*byte)(unsafe.Add(unsafe.Pointer(u.UpdateDataPlayer().Player), 2251)) = sp.PlayerClass
		if sp.Harpoon {
			*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 132)) = p.proxy.life.players[1].CObj()
		}
		w := unsafe.Slice((*uint16)(u.UpdateData), 5)
		w[2], w[3], w[4] = sp.Mana, sp.OldMana, sp.MaxMana
		*(*uint32)(unsafe.Add(unsafe.Pointer(u.UpdateDataPlayer().Player), 3680)) = sp.Status
	}
	r.unit = u
	p.identify(u.CObj(), 54000)
	if sp.Subject == 0 {
		r.unit = nil
	}
	restoreStrings := p.proxy.core.PortTestTradeStrings()
	snapshotProtection, freeProtection := portTestResourceProtection(p, sp)
	r.protection = snapshotProtection
	return func() {
		freeProtection()
		Nox_xxx_harpoonBreakForPlr_537520 = oldHarpoon
		Nox_xxx_pickupDefault_4F31E0 = oldPickup
		restoreStrings()
		*u = oldUnit
		if savedUD != nil {
			copy(unsafe.Slice((*byte)(u.UpdateData), len(savedUD)), savedUD)
		}
		for _, free := range frees {
			free()
		}
		p.resources = nil
		noxflags.ResetEngine()
		noxflags.SetEngine(oldEngine)
	}
}
func (p *portTestShopPools) resourceAction(a PortTestShopAction) uint32 {
	u := p.resources.unit
	var ptr unsafe.Pointer
	if u != nil {
		ptr = u.CObj()
	}
	unit := C.int(uintptr(ptr))
	value := C.int(a.Value)
	switch a.Op {
	case PortTestResourceSetHP:
		return uint32(C.nox_xxx_unitSetHP_4E4560(asObjectC(u), C.ushort(a.Value)))
	case PortTestResourceAdjustHP:
		C.nox_xxx_unitAdjustHP_4EE460(asObjectC(u), value)
	case PortTestResourceInformOwner:
		resourceInformOwner(u)
	case PortTestResourceDamage:
		C.nox_xxx_unitDamageClear_4EE5E0(asObjectC(u), value)
	case PortTestResourceRestoreHP:
		C.nox_xxx_unitHPsetOnMax_4EE6F0(unit)
	case PortTestResourceHPHistory:
		C.nox_xxx_playerHP_4EE730(unit)
	case PortTestResourceGetHP:
		return uint32(C.nox_xxx_unitGetHP_4EE780(asObjectC(u)))
	case PortTestResourceGetMaxHP:
		return uint32(C.nox_xxx_unitGetMaxHP_4EE7A0(unit))
	case PortTestResourceSetMaxHP:
		return uint32(C.nox_xxx_unitSetMaxHP_4EE7C0(unit, C.short(a.Value)))
	case PortTestResourcePoison:
		return uint32(C.nox_xxx_activatePoison_4EE7E0(unit, value, C.int(a.Item)))
	case PortTestResourcePoisonReduce:
		C.nox_xxx_updatePoison_4EE8F0(asObjectC(u), value)
	case PortTestResourcePoisonRemove:
		C.nox_xxx_removePoison_4EE9D0(asObjectC(u))
	case PortTestResourcePoisonSet:
		C.nox_xxx_setSomePoisonData_4EEA90(unit, value)
	case PortTestResourceManaAdd:
		out := uint32(C.nox_xxx_playerManaAdd_4EEB80(asObjectC(u), C.short(a.Value)))
		// The unsupported nonplayer path returns the low sixteen pointer bits.
		if u != nil && u.ObjClass&4 == 0 {
			if out != uint32(uint16(uintptr(ptr))) {
				panic("mana-add pointer return")
			}
			return 54000
		}
		return out
	case PortTestResourceManaSub:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_playerManaSub_4EEBF0(unit, value))))
	case PortTestResourceGetMana:
		return uint32(C.nox_xxx_unitGetOldMana_4EEC80(unit))
	case PortTestResourceGetMaxMana:
		return uint32(C.nox_xxx_playerGetMaxMana_4EECB0(unit))
	case PortTestResourceSetMaxMana:
		return uint32(C.nox_xxx_playerSetMaxMana_4EECD0(unit, C.short(a.Value)))
	case PortTestResourceManaRefresh:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_playerManaRefresh_4EECF0(unit))))
	case PortTestResourceGoldAdd:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_playerAddGold_4FA590(unit, value))))
	case PortTestResourceGoldSub:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_playerSubGold_4FA5D0(unit, C.uint(a.Value)))))
	case PortTestResourceGoldSet:
		resourceSetGold(u, int32(value))
	case PortTestResourceGetGold:
		return uint32(C.nox_xxx_playerGetGold_4FA6B0(unit))
	case PortTestResourceObjectGold:
		return uint32(C.nox_object_getGold_4FA6D0(asObjectC(u)))
	case PortTestResourceGoldPickup:
		item := p.items[a.Item].u
		*(*uint32)(item.InitData) = p.proxy.callbacks.shop.spec.Resources.GoldItem
		return uint32(C.nox_xxx_pickupGold_4F3A60_obj_pickup(unit, C.int(uintptr(item.CObj())), value))
	default:
		panic("resource fixture operation")
	}
	return 0
}
func (p *portTestShopPools) resourceSnapshot() (data [][]uint32, messages [][]byte) {
	r := p.resources
	if r == nil {
		return
	}
	words := func(ptr unsafe.Pointer, size int) []uint32 {
		out := make([]uint32, size/4)
		b := unsafe.Slice((*byte)(ptr), size)
		for i := range out {
			out[i] = p.normalize(binary.LittleEndian.Uint32(b[i*4:]))
		}
		return out
	}
	data = append(data, append([]uint32{uint32(C.resourceDieCalls()), p.normalize(uint32(C.resourceDieWho()))}, r.harpoonCalls...))
	if r.unit != nil {
		data = append(data, words(r.unit.CObj(), int(unsafe.Sizeof(server.Object{}))-8))
	}
	for _, b := range r.blocks {
		for i := 0; i < 8; i++ {
			if b[i] != 0xa5 || b[len(b)-8+i] != 0x5a {
				panic("resource fixture guard")
			}
		}
		data = append(data, words(unsafe.Pointer(&b[8]), len(b)-16))
	}
	data = append(data, r.protection())
	if len(r.pickupCalls) != 0 {
		data = append(data, append([]uint32(nil), r.pickupCalls...))
	}
	for _, ind := range []ntype.PlayerInd{1, 7, 31} {
		messages = append(messages, p.proxy.core.NetList.CopyPacketsA(ind, 1))
	}
	return
}
