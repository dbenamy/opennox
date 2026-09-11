//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME3_3.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2386576, dword_5d4594_3484, dword_5d4594_2523804;
extern uint32_t dword_5d4594_2488720, dword_5d4594_2488724, dword_5d4594_2488728;
static uint32_t invTrace[4097], invUseDelete, invDropResult;
static uint32_t* invTracePtr(void) {return invTrace;}
static void invReset(int useDelete,int dropResult) {memset(invTrace,0,sizeof(invTrace));invUseDelete=useDelete;invDropResult=dropResult;}
static int invUse(nox_object_t* u,nox_object_t* it) {
 uint32_t i=1+6*invTrace[0]++;
 if(i+5<4097) {invTrace[i]=1;invTrace[i+1]=(uint32_t)u;invTrace[i+2]=(uint32_t)it;}
 if(invUseDelete) it->obj_flags|=0x20;
 return 1;
}
static int invDrop(nox_object_t* u,nox_object_t* it,float2* pos) {
 uint32_t i=1+6*invTrace[0]++;
 if(i+5<4097) {invTrace[i]=2;invTrace[i+1]=(uint32_t)u;invTrace[i+2]=(uint32_t)it;memcpy(&invTrace[i+3],pos,8);}
 return invDropResult;
}
static void* invUsePtr(void){return invUse;}
static void* invDropPtr(void){return invDrop;}
static uint64_t invCall(int op,nox_object_t* u,nox_object_t* it,int value,int arg,float radius,float2* pos) {
 int up=(int)(uintptr_t)u, ip=(int)(uintptr_t)it;
 switch(op) {
 case 0: sub_4ED0C0(u, it);return 0;
 case 1: return (uint32_t)(uintptr_t)nox_xxx_dropDefault_4ED290(u,it,pos);
 case 2: return (uint32_t)(uintptr_t)nox_GlyphDrop_4ED500(up,ip,pos);
 case 3: return (uint32_t)(uintptr_t)nox_xxx_dropTrap_4ED580(up,ip,pos);
 case 4: return (uint32_t)(uintptr_t)nox_xxx_dropCrown_4ED5E0(up,ip,(int*)pos);
 case 5: return (uint32_t)(uintptr_t)nox_xxx_dropTreasure_4ED710(up,ip,(int*)pos);
 case 6: return (uint32_t)(uintptr_t)nox_xxx_drop_4ED790(u,it,pos);
 case 7: return (uint32_t)(uintptr_t)nox_xxx_drop_4ED810(up,ip,(float*)pos);
 case 8: return (uint32_t)(uintptr_t)nox_xxx_invForceDropItem_4ED930(up,(uint32_t*)it);
 case 9: return (uint32_t)(uintptr_t)sub_4ED970(radius,(float2*)((char*)u+56),pos);
 case 10: return (uint32_t)(uintptr_t)nox_xxx_dropAllItems_4EDA40((uint32_t*)u);
 case 12: return (uint32_t)(uintptr_t)sub_4EDDE0(up,(uint32_t*)it,(int*)pos);
 case 13: return (uint32_t)(uintptr_t)nox_xxx_dropFood_4EDE50(up,ip,(int*)pos);
 case 14: nox_xxx_chest_4EDF00(up,ip);return 0;
 case 16: return (uint32_t)(uintptr_t)nox_xxx_dropAnkhTradable_4EE370(up,ip,(int*)pos);
 case 17: nox_xxx_inventoryPutImpl_4F3070(u,it,value);return 0;
 case 18: return (uint32_t)(uintptr_t)nox_xxx_pickupFood_4F3350(up,ip,value);
 case 19: return (uint32_t)(uintptr_t)sub_4F3400(up,ip,value);
 case 20: return (uint32_t)(uintptr_t)nox_xxx_pickupUse_4F34D0(up,ip,value);
 case 21: return (uint32_t)(uintptr_t)nox_xxx_pickupTrap_4F3510(up,ip,value);
 case 22: return (uint32_t)(uintptr_t)nox_xxx_pickupTreasure_4F3580(up,ip,value);
 case 23: return (uint32_t)(uintptr_t)nox_xxx_pickupAmmo_4F3B00(up,it,value,arg);
 case 24: return (uint32_t)(uintptr_t)nox_xxx_pickupSpellbook_4F3C60(up,ip,value);
 case 25: return (uint32_t)(uintptr_t)nox_xxx_pickupAbilitybook_4F3CE0(up,ip,value);
 case 26: return (uint32_t)(uintptr_t)sub_4F3DD0(up,ip);
 case 27: return (uint32_t)(uintptr_t)sub_53A720(up,it,value,arg);
 case 28: return (uint32_t)(uintptr_t)nox_xxx_sendMsgOblivionPickup_53A9C0(up,it,value,arg);
 case 29: return (uint32_t)(uintptr_t)nox_xxx_dropWeapon_53AB10(up,(uint32_t*)it,(int*)pos);
 case 30: return (uint32_t)(uintptr_t)nox_xxx_pickupArmor_53E7F0(up,ip,value,arg);
 case 31: return (uint32_t)(uintptr_t)nox_xxx_dropArmor_53EB70(up,(uint32_t*)it,(int*)pos);
 case 32: return (uint32_t)(uintptr_t)nox_xxx_ItemIsDroppable_53EBF0(ip);
 }
 return 0;
}
*/
import "C"
import (
	"bytes"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestInventorySpec struct {
	NilDrop     bool
	TeamMembers bool
	Teams       [4]byte // Three players, then the first item.
	CrownTimes  [3]uint32

	ServerFlags, Gameplay, Treasure, TreasureMax uint32
	DropTable                                    [][3]uint32
	FoodDrop, FoodPickup                         [][3]uint32

	Blocked               bool
	WeaponBits, ArmorBits map[uint16]uint32
	Owned                 []int
	Materials             []uint16
	ItemTypes             []string
	Shape                 [5]uint32
	WallMode              int

	Linked                                                            []int
	Weights                                                           []byte
	Carry                                                             uint16
	Position, Target                                                  types.Pointf
	Radius                                                            float32
	PickupInsert, UseDelete, DropResult, DefaultDrop, NilItem, NilPos bool
}
type portTestInventory struct {
	playerHandle uint32
	blocks       [][]byte
	frees        []func()

	calls  []uint32
	pos    *types.Pointf
	result uint64
}

const PortTestInventory4ED0C0 = 300
const PortTestInventory4ED290 = 301
const PortTestInventory4ED500 = 302
const PortTestInventory4ED580 = 303
const PortTestInventory4ED5E0 = 304
const PortTestInventory4ED710 = 305
const PortTestInventory4ED790 = 306
const PortTestInventory4ED810 = 307
const PortTestInventory4ED930 = 308
const PortTestInventory4ED970 = 309
const PortTestInventory4EDA40 = 310
const PortTestInventory4EDCD0 = 311
const PortTestInventory4EDDE0 = 312
const PortTestInventory4EDE50 = 313
const PortTestInventory4EDF00 = 314
const PortTestInventory4EE2A0 = 315
const PortTestInventory4EE370 = 316
const PortTestInventory4F3070 = 317
const PortTestInventory4F3350 = 318
const PortTestInventory4F3400 = 319
const PortTestInventory4F34D0 = 320
const PortTestInventory4F3510 = 321
const PortTestInventory4F3580 = 322
const PortTestInventory4F3B00 = 323
const PortTestInventory4F3C60 = 324
const PortTestInventory4F3CE0 = 325
const PortTestInventory4F3DD0 = 326
const PortTestInventory53A720 = 327
const PortTestInventory53A9C0 = 328
const PortTestInventory53AB10 = 329
const PortTestInventory53E7F0 = 330
const PortTestInventory53EB70 = 331
const PortTestInventory53EBF0 = 332

func (p *portTestShopPools) inventoryPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.Inventory
	if sp == nil {
		return func() {}
	}
	if p.resources == nil {
		panic("inventory fixture requires resources")
	}
	pos, freePos := alloc.New(types.Pointf{})
	*pos = sp.Target
	p.inventory = &portTestInventory{pos: pos, frees: []func(){freePos}}
	// Borrowed player slabs need the same server association as pool objects
	// when retained ownership code evaluates hostility and minimap visibility.
	var units []*server.Object
	for i := range p.proxy.life.players {
		units = append(units, &p.proxy.life.players[i])
	}
	restorePlayers := server.PortTestAttachAI(p.proxy.core, units...)
	p.inventory.playerHandle = *(*uint32)(unsafe.Add(units[0].CObj(), 772))
	for i, u := range units {
		u.TeamVal.ID = server.TeamID(sp.Teams[i])

		if sp.TeamMembers {
			p.identify(unsafe.Pointer(&u.TeamVal), 620000+uint32(i))
		}
		*(*uint32)(unsafe.Add(u.UpdateData, 264)) = sp.CrownTimes[i]
	}
	restoreServer := p.proxy.core.PortTestInventoryEnvironment(sp.Blocked, sp.WeaponBits, sp.ArmorBits)
	if sp.TeamMembers {
		for _, u := range units {
			if team := p.proxy.core.Teams.ByID(u.TeamVal.ID); team != nil {
				head := (*uint32)(unsafe.Add(unsafe.Pointer(team), 44))
				u.TeamVal.Field0 = *head
				*head = uint32(uintptr(unsafe.Pointer(&u.TeamVal)))
			}
		}
	}

	walls, unchanged, freeWalls := p.proxy.core.PortTestPathWalls()
	walls(sp.WallMode)
	oldDecay := C.dword_5d4594_2386576
	C.dword_5d4594_2386576 = 0
	C.invReset(C.int(bool2int(sp.UseDelete)), C.int(bool2int(sp.DropResult)))
	p.identify(C.nox_xxx_dropDefault_4ED290, 55000)
	p.identify(C.invUsePtr(), 55001)
	p.identify(C.invDropPtr(), 55002)
	oldState := Nox_xxx_playerSetState_4FA020
	Nox_xxx_playerSetState_4FA020 = func(u *server.Object, state server.PlayerState) bool {
		p.inventory.calls = append(p.inventory.calls, 4, p.normalize(uint32(uintptr(u.CObj()))), uint32(state))
		return true
	}
	oldServerFlags, oldPause := C.dword_5d4594_3484, C.dword_5d4594_2523804
	C.dword_5d4594_3484, C.dword_5d4594_2523804 = C.uint32_t(sp.ServerFlags), 1
	oldGameplay := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
	noxflags.SetGamePlay(noxflags.GameplayFlag(sp.Gameplay))
	oldTreasure := *memmap.PtrUint32(0x5D4594, 1548528)
	*memmap.PtrUint32(0x5D4594, 1548528) = sp.TreasureMax
	oldHeads := [32]*server.MinimapItem{}
	otherPlayers := map[int][]byte{}
	for i := range oldHeads {
		pl := p.proxy.core.Players.ByIndRaw(ntype.PlayerInd(i))
		oldHeads[i] = pl.Field4580
		if i != 1 && i != 7 && i != 31 {
			otherPlayers[i] = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl))))
			p.identify(unsafe.Pointer(pl), 610000+uint32(i))
		}
		pl.Field4580 = nil
	}
	for i := range p.proxy.life.players {
		p.proxy.life.players[i].UpdateDataPlayer().Player.Field2152 = sp.Treasure
	}
	var foodRestores []func()
	for i, entries := range [][][3]uint32{sp.FoodDrop, sp.FoodPickup} {
		off := uintptr(205704)
		if i == 1 {
			off = 215640
		}
		words := unsafe.Slice(memmap.PtrUint32(0x587000, off), 10)
		old := append([]uint32(nil), words...)
		clear(words)
		if len(entries) > 4 {
			panic("food table capacity")
		}
		for j, v := range entries {
			words[2*j] = v[0]
			words[2*j+1] = uint32(uint16(v[1])) | uint32(uint16(v[2]))<<16
		}
		foodRestores = append(foodRestores, func() { copy(words, old) })
	}
	oldPickup := Nox_xxx_pickupDefault_4F31E0
	Nox_xxx_pickupDefault_4F31E0 = func(u, it *server.Object, a3, a4 int) bool {
		p.inventory.calls = append(p.inventory.calls, 3, p.normalize(uint32(uintptr(u.CObj()))), p.normalize(uint32(uintptr(it.CObj()))), uint32(a3), uint32(a4))
		result := p.proxy.callbacks.shop.spec.Resources.PickupResult
		if result && sp.PickupInsert {
			C.nox_xxx_inventoryPutImpl_4F3070(asObjectC(u), asObjectC(it), C.int(a3))
		}
		return result
	}
	if u := p.resources.unit; u != nil {
		u.PosVec = sp.Position
		for i, v := range sp.Shape {
			*(*uint32)(unsafe.Add(u.CObj(), 172+4*i)) = v
		}
		*(*uint16)(unsafe.Add(u.CObj(), 490)) = sp.Carry
	}
	offsets := []uintptr{1568244, 1568248, 1568252, 1568256, 2488712, 2488716}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		v := memmap.PtrUint32(0x5D4594, off)
		old[i] = *v
		*v = 0
	}
	oldShield1, oldShield2, oldDrop := C.dword_5d4594_2488720, C.dword_5d4594_2488724, C.dword_5d4594_2488728
	C.dword_5d4594_2488720, C.dword_5d4594_2488724, C.dword_5d4594_2488728 = 0, 0, 1
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 279432), 48)
	oldTable := append([]uint32(nil), table...)
	clear(table)
	if len(sp.DropTable) > 15 {
		panic("drop table capacity")
	}
	for i, v := range sp.DropTable {
		copy(table[3*i:3*i+3], v[:])
	}
	return func() {
		restorePlayers()
		Nox_xxx_pickupDefault_4F31E0 = oldPickup
		Nox_xxx_playerSetState_4FA020 = oldState
		C.dword_5d4594_3484, C.dword_5d4594_2523804 = oldServerFlags, oldPause
		noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
		noxflags.SetGamePlay(oldGameplay)
		*memmap.PtrUint32(0x5D4594, 1548528) = oldTreasure
		for _, restore := range foodRestores {
			restore()
		}
		for i, head := range oldHeads {
			pl := p.proxy.core.Players.ByIndRaw(ntype.PlayerInd(i))
			seen := map[*server.MinimapItem]bool{}
			for m := pl.Field4580; m != nil && !seen[m]; {
				seen[m] = true
				next := m.Field8
				alloc.Free(m)
				m = next
			}
			pl.Field4580 = head
		}
		for i, b := range otherPlayers {
			pl := p.proxy.core.Players.ByIndRaw(ntype.PlayerInd(i))
			copy(unsafe.Slice((*byte)(unsafe.Pointer(pl)), len(b)), b)
		}
		if !unchanged() {
			panic("inventory wall input modified")
		}
		freeWalls()
		restoreServer()
		C.dword_5d4594_2386576 = oldDecay
		for _, o := range p.items {
			if o.alive {
				o.u.UpdateData, o.u.InitData, o.u.UseData.Ptr, o.u.HealthData = nil, nil, nil, nil
			}
			o.init, o.health = nil, nil // guarded regions belong to this fixture
		}
		for _, free := range p.inventory.frees {
			free()
		}
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		C.dword_5d4594_2488720, C.dword_5d4594_2488724, C.dword_5d4594_2488728 = oldShield1, oldShield2, oldDrop
		copy(table, oldTable)
		p.inventory = nil
	}
}
func (p *portTestShopPools) inventoryItems() {
	sp := p.proxy.callbacks.shop.spec.Inventory
	if sp == nil {
		return
	}
	for i, o := range p.items {
		u := o.u
		region := func(size int, id uint32) unsafe.Pointer {
			b, free := alloc.Make([]byte{}, size+16)
			for j := 0; j < 8; j++ {
				b[j] = 0xa5
				b[len(b)-8+j] = 0x5a
			}
			p.inventory.blocks = append(p.inventory.blocks, b)
			p.inventory.frees = append(p.inventory.frees, free)
			ptr := unsafe.Pointer(&b[8])
			p.identify(ptr, id)
			return ptr
		}
		init := region(64, 58000+uint32(i))
		copy(unsafe.Slice((*byte)(init), 20), unsafe.Slice((*byte)(u.InitData), 20))
		alloc.FreePtr(u.InitData)
		u.InitData, o.init, o.initSize = init, init, 64
		for j, enabled := range p.proxy.callbacks.shop.spec.Items[i].Mods {
			if enabled {
				*(*uintptr)(unsafe.Add(init, 4*j)) = uintptr(unsafe.Add(p.proxy.callbacks.shop.ptr(8), j*int(unsafe.Sizeof(server.ModifierEff{}))))
			}
		}
		use := region(128, 56000+uint32(i))
		if u.UseData.Ptr != nil {
			copy(unsafe.Slice((*byte)(use), 128), unsafe.Slice((*byte)(u.UseData.Ptr), 128))
			alloc.FreePtr(u.UseData.Ptr)
		}
		u.UseData.Ptr, o.useSize = use, 128
		if u.HealthData != nil {
			hp := region(int(unsafe.Sizeof(server.HealthData{})), 59000+uint32(i))
			*(*server.HealthData)(hp) = *u.HealthData
			alloc.FreePtr(unsafe.Pointer(u.HealthData))
			u.HealthData = (*server.HealthData)(hp)
			o.health = hp
		}
		if i < len(sp.ItemTypes) && sp.ItemTypes[i] != "" {
			u.TypeInd = uint16(p.proxy.core.Types.IndByID(sp.ItemTypes[i]))
		}
		size := 64
		if u.ObjClass&2 != 0 {
			size = 2200
		}
		b, free := alloc.Make([]byte{}, size+16)
		for j := 0; j < 8; j++ {
			b[j] = 0xa5
			b[len(b)-8+j] = 0x5a
		}
		p.inventory.blocks = append(p.inventory.blocks, b)
		p.inventory.frees = append(p.inventory.frees, free)
		u.UpdateData = unsafe.Pointer(&b[8])
		p.identify(u.UpdateData, 57000+uint32(i))
		if u.ObjClass&2 != 0 {
			copy(b[8:len(b)-8], unsafe.Slice((*byte)(p.proxy.combat.actor.UpdateData), 2200))
			p.identify(*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 484)), 54004)
			p.identify(*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 488)), 54005)
		}
		if i < len(sp.Materials) {
			u.Material = sp.Materials[i]
		}
		u.Use.Ptr = C.invUsePtr()
		u.Drop.Ptr = C.invDropPtr()
		if sp.DefaultDrop {
			u.Drop.Ptr = C.nox_xxx_dropDefault_4ED290
		}
		if sp.NilDrop {
			u.Drop.Ptr = nil
		}

		if i < len(sp.Weights) {
			*(*byte)(unsafe.Add(u.CObj(), 488)) = sp.Weights[i]
		}
		u.PosVec = sp.Target
		if i == 0 {
			u.TeamVal.ID = server.TeamID(sp.Teams[3])
		}
	}
	u := p.resources.unit
	// These links are fixture inputs; ownership notifications belong to the
	// operations under test, not construction of the initial inventory.
	if u != nil {
		u.Field129 = nil
	}
	for _, index := range sp.Owned {
		it := p.items[index].u
		it.ObjOwner = u
		it.Field128 = u.Field129
		u.Field129 = it
	}
	for i, index := range sp.Linked {
		it := p.items[index].u
		it.InvHolder = u
		it.Field125 = nil
		it.InvNextItem = nil
		if i == 0 {
			u.InvFirstItem = it
		} else {
			prev := p.items[sp.Linked[i-1]].u
			prev.InvNextItem = it
			it.Field125 = prev
		}
	}
}
func (p *portTestShopPools) inventoryAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec.Inventory
	u := p.resources.unit
	var it *server.Object
	if !sp.NilItem && a.Item >= 0 {
		it = p.items[a.Item].u
	}
	pos := (*C.float2)(unsafe.Pointer(p.inventory.pos))
	if sp.NilPos {
		pos = nil
	}
	var out uint64
	switch a.Op {
	case PortTestInventory4EDCD0:
		out = uint64(bool2int(inventoryDropEligible(u, it)))
	case PortTestInventory4EE2A0:
		out = math.Float64bits(inventoryShapeRadius(u))
	default:
		out = uint64(C.invCall(C.int(a.Op-300), asObjectC(u), asObjectC(it), C.int(a.Value), C.int(a.Side), C.float(sp.Radius), pos))
	}
	p.inventory.result = out
	if a.Op == PortTestInventory4ED970 {
		if uint32(out) != uint32(uintptr(unsafe.Pointer(pos))) {
			panic("inventory placement return pointer")
		}
		p.inventory.result = 55003
		return 55003
	}
	return uint32(out)
}
func (p *portTestShopPools) inventorySnapshot() []uint32 {
	if p.inventory == nil {
		return nil
	}
	r := p.inventory
	var minimap []*server.MinimapItem
	for i := 0; i < 32; i++ {
		head := p.proxy.core.Players.ByIndRaw(ntype.PlayerInd(i)).Field4580
		seen := map[*server.MinimapItem]bool{}
		for m := head; m != nil && !seen[m]; m = m.Field8 {
			seen[m] = true
			p.identify(unsafe.Pointer(m), uint32(600000+i*100+len(seen)))
			minimap = append(minimap, m)
		}
	}
	out := []uint32{p.normalize(uint32(r.result)), uint32(r.result >> 32), math.Float32bits(r.pos.X), math.Float32bits(r.pos.Y)}
	b := unsafe.Slice((*uint32)(unsafe.Pointer(C.invTracePtr())), 4097)
	if b[0] > 680 {
		panic("inventory callback trace overflow")
	}
	for _, v := range b[:1+6*b[0]] {
		out = append(out, p.normalize(v))
	}
	out = append(out, r.calls...)
	for _, off := range []uintptr{1568244, 1568248, 1568252, 1568256, 2488712, 2488716} {
		out = append(out, *memmap.PtrUint32(0x5D4594, off))
	}
	out = append(out, uint32(C.dword_5d4594_2488720), uint32(C.dword_5d4594_2488724), uint32(C.dword_5d4594_2488728))
	out = append(out, p.normalize(uint32(C.dword_5d4594_2386576)))
	for _, b := range p.inventory.blocks {
		for i := 0; i < 8; i++ {
			if b[i] != 0xa5 || b[len(b)-8+i] != 0x5a {
				panic("inventory item UD guard")
			}
		}
		for i := 8; i < len(b)-8; i += 4 {
			out = append(out, p.normalize(*(*uint32)(unsafe.Pointer(&b[i]))))
		}
	}
	out = append(out, uint32(noxflags.GetGame()), uint32(noxflags.GetEngine()), uint32(noxflags.GetGamePlay()), uint32(C.dword_5d4594_3484))
	for _, m := range minimap {
		out = append(out, p.normalize(uint32(uintptr(unsafe.Pointer(m)))), m.Field0, p.normalize(uint32(uintptr(m.Field4.CObj()))), p.normalize(uint32(uintptr(unsafe.Pointer(m.Field8)))), p.normalize(uint32(uintptr(unsafe.Pointer(m.Field12)))))
	}
	for i := 0; i < 32; i++ {
		if i == 1 || i == 7 || i == 31 {
			continue
		}
		pl := p.proxy.core.Players.ByIndRaw(ntype.PlayerInd(i))
		if pl.Active == 0 && pl.Field4580 == nil {
			continue
		}
		out = append(out, 610000+uint32(i))
		for _, v := range unsafe.Slice((*uint32)(unsafe.Pointer(pl)), int(unsafe.Sizeof(*pl))/4) {
			out = append(out, p.normalize(v))
		}
	}
	return out
}
