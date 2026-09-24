//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "MixPatch.h"
static uint32_t eqTrace[4097],eqReturn[4],eqOutput[4];
static uint32_t* eqTracePtr(void){return eqTrace;}
static void eqReset(void){memset(eqTrace,0,sizeof(eqTrace));}
static void eqConfigure(int i,uint32_t ret,uint32_t out){eqReturn[i]=ret;eqOutput[i]=out;}
static uint32_t eqRecord(int kind,void* mod,void* u,void* it,uint32_t before,uint32_t after){
 uint32_t n=1+6*eqTrace[0]++;
 if(n+5<4097){eqTrace[n]=kind;eqTrace[n+1]=(uint32_t)mod;eqTrace[n+2]=(uint32_t)u;eqTrace[n+3]=(uint32_t)it;eqTrace[n+4]=before;eqTrace[n+5]=after;}
 return after;
}
static int eqEngage(void* mod,void* u,void* it){return eqRecord(1,mod,u,it,0,eqReturn[((uint32_t*)mod)[1]-40]);}
static int eqDisengage(void* mod,void* u,void* it){return eqRecord(2,mod,u,it,0,eqReturn[((uint32_t*)mod)[1]-40]);}
static void eqDefend(void* mod,void* u,int a,void* it,int b,float* value){
 uint32_t before;memcpy(&before,value,4);uint32_t out=eqOutput[((uint32_t*)mod)[1]-40];
 eqRecord(3,mod,u,it,before,out);eqRecord(4,mod,(void*)(uintptr_t)a,(void*)(uintptr_t)b,0,0);memcpy(value,&out,4);
}
static void* eqEngagePtr(void){return eqEngage;}
static void* eqDisengagePtr(void){return eqDisengage;}
static void* eqDefendPtr(void){return eqDefend;}

*/
import "C"

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestEquipmentDef struct {
	Words    map[int]uint32
	Type     uint16
	Armor    bool
	Strength uint16
	Coeff    uint32
}
type PortTestEquipmentEffect struct {
	Engage, Disengage, Defend bool
	Return, Output            uint32
}
type PortTestEquipmentSpec struct {
	ColdTable                                  bool
	GameEx                                     uint32
	Cheat, NilUnit, HolderOnly                 bool
	Strength                                   uint32
	NPCStrength, PlayerState                   byte
	ActiveWeapon, SecondaryWeapon, SavedShield int // One-based item indices; zero is nil.
	ArmorFlags, WeaponFlags                    uint32
	ActiveAbilities                            uint32
	Threshold                                  uint64
	Definitions                                []PortTestEquipmentDef
	Effects                                    [4]PortTestEquipmentEffect
	ItemTeamWord                               []uint32
	TableNames                                 []string
	TableFlags                                 []uint32
}
type portTestEquipment struct {
	blocks [][]byte
	frees  []func()
	result uint64
}

func (p *portTestShopPools) equipmentPrepare() func() {
	sp := p.proxy.callbacks.shop.spec.Equipment
	if sp == nil {
		return func() {}
	}
	if p.inventory == nil {
		panic("equipment requires inventory fixture")
	}
	p.equipment = &portTestEquipment{}
	C.eqReset()
	p.identify(C.eqEngagePtr(), 65000)
	p.identify(C.eqDisengagePtr(), 65001)
	p.identify(C.eqDefendPtr(), 65002)
	oldGameEx, oldCheat, oldThreshold := gameex_flags, nox_cheat_allowall, qword_581450_9512
	gameex_flags = uint32(sp.GameEx)
	nox_cheat_allowall = int32(bool2int(sp.Cheat))
	qword_581450_9512 = uint64(sp.Threshold)
	core := p.proxy.core
	oldWeapon, oldArmor := core.Modif.Dword_5d4594_251600, core.Modif.Dword_5d4594_251608
	core.Modif.Dword_5d4594_251600, core.Modif.Dword_5d4594_251608 = nil, nil
	for i, d := range sp.Definitions {
		ptr := p.equipmentRegion(int(unsafe.Sizeof(server.Modifier{})), 630000+uint32(i))
		m := (*server.Modifier)(ptr)
		m.TypeInd = uint32(d.Type)
		m.ReqStrength60 = d.Strength
		m.DamageCoeffOrArmor64 = math.Float32frombits(d.Coeff)
		for off, v := range d.Words {
			if off < 36 || off > 76 || off%4 != 0 {
				panic("equipment definition word offset")
			}
			*(*uint32)(unsafe.Add(ptr, off)) = v
		}
		head := &core.Modif.Dword_5d4594_251600
		if d.Armor {
			head = &core.Modif.Dword_5d4594_251608
		}
		m.Next80 = *head
		if *head != nil {
			(*head).Prev84 = m
		}
		*head = m
	}
	u := p.resources.unit
	if u != nil {
		if u.ObjClass&4 != 0 {
			pl := u.UpdateDataPlayer().Player
			*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 2239)) = sp.Strength
			*(*uint32)(unsafe.Pointer(pl)) = sp.ArmorFlags
			*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4)) = sp.WeaponFlags
			*(*byte)(unsafe.Add(u.UpdateData, 88)) = sp.PlayerState
		} else if u.ObjClass&2 != 0 {
			*(*byte)(unsafe.Add(u.UpdateData, 1324)) = sp.NPCStrength
		}
		p.identify(unsafe.Add(u.CObj(), 688), 65003) // npc sync helper's interior return pointer
	}
	oldAbilities := core.Abils.ByUnit
	core.Abils.Reset()
	if u != nil {
		ad := core.Abils.GetFor(u)
		for i := 1; i <= 2; i++ {
			if sp.ActiveAbilities&(1<<i) != 0 {
				ad.ExecList = &server.ExecAbilityClass{Abil: server.Ability(i), Active: 1, Next: ad.ExecList}
			}
		}
	}
	if len(sp.TableNames) > 15 {
		panic("equipment drop table capacity")
	}
	if sp.ColdTable {
		dword_5d4594_2488728 = 0
	}
	if sp.TableNames != nil {
		table := unsafe.Slice(memmap.PtrUint32(0x587000, 279432), 48)
		clear(table)
		for i, name := range sp.TableNames {
			ptr := p.equipmentRegion((len(name)+1+3)&^3, 630500+uint32(i))
			copy(unsafe.Slice((*byte)(ptr), len(name)), name)
			table[3*i] = uint32(uintptr(ptr))
			table[3*i+1] = 0xdeadbeef
			if i < len(sp.TableFlags) {
				table[3*i+2] = sp.TableFlags[i]
			}
		}
	}
	return func() {
		gameex_flags, nox_cheat_allowall, qword_581450_9512 = oldGameEx, oldCheat, oldThreshold
		core.Modif.Dword_5d4594_251600, core.Modif.Dword_5d4594_251608 = oldWeapon, oldArmor
		core.Abils.ByUnit = oldAbilities
		for _, free := range p.equipment.frees {
			free()
		}
		p.equipment = nil
	}
}
func (p *portTestShopPools) equipmentRegion(size int, id uint32) unsafe.Pointer {
	b, free := alloc.Make([]byte{}, size+16)
	for i := 0; i < 8; i++ {
		b[i] = 0xa5
		b[len(b)-8+i] = 0x5a
	}
	p.equipment.blocks = append(p.equipment.blocks, b)
	p.equipment.frees = append(p.equipment.frees, free)
	ptr := unsafe.Pointer(&b[8])
	p.identify(ptr, id)
	return ptr
}
func (p *portTestShopPools) equipmentItems() {
	sp := p.proxy.callbacks.shop.spec.Equipment
	if sp == nil {
		return
	}
	for i, e := range sp.Effects {
		m := (*server.ModifierEff)(unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*int(unsafe.Sizeof(server.ModifierEff{}))))
		if e.Engage {
			m.Engage112 = C.eqEngagePtr()
		}
		if e.Disengage {
			m.Disengage116 = C.eqDisengagePtr()
		}
		if e.Defend {
			m.Defend76.Fnc = C.eqDefendPtr()
		}
		C.eqConfigure(C.int(i), C.uint32_t(e.Return), C.uint32_t(e.Output))
	}
	if sp.HolderOnly {
		for _, item := range p.items {
			item.u.InvHolder = p.resources.unit
		}
	}
	for i, v := range sp.ItemTeamWord {
		p.items[i].u.TeamVal.Field0 = v
	}
	u := p.resources.unit
	if u != nil && u.ObjClass&4 != 0 {
		ptr := func(index int) uint32 {
			if index == 0 {
				return 0
			}
			return uint32(uintptr(p.items[index-1].u.CObj()))
		}
		*(*uint32)(unsafe.Add(u.UpdateData, 104)) = ptr(sp.ActiveWeapon)
		*(*uint32)(unsafe.Add(u.UpdateData, 108)) = ptr(sp.SecondaryWeapon)
		*(*uint32)(unsafe.Add(unsafe.Pointer(u.UpdateDataPlayer().Player), 2500)) = ptr(sp.SavedShield)
	}
}

// equipmentFixtureCall mirrors eqCall while routing fixture operations to the Go owners.
func equipmentFixtureCall(op int, u, it *server.Object, value, side int32) uint64 {
	switch op {
	case 0:
		return uint64(uint32(equipmentNPCDequipWeapon(u, it)))
	case 1:
		equipmentDequipAmmo(u, int(value), int(side))
	case 2:
		return uint64(uint32(equipmentDequipWeapon(u, it, int(value), int(side))))
	case 3:
		return uint64(uint32(equipmentNPCEquipWeapon(u, it)))
	case 4:
		equipmentRemoveShields(u)
	case 5:
		return uint64(uint32(equipmentEquipWeapon(u, it, int(value), int(side))))
	case 6:
		return uint64(uint32(equipmentEquipBow(u)))
	case 7:
		equipmentPickupSound(u, it)
	case 8:
		equipmentDropSound(it)
	case 10:
		return uint64(uint32(equipmentArmorMask(it)))
	case 11:
		return uint64(uint32(equipmentRecalculate(u)))
	case 12:
		return uint64(uint32(equipmentNPCDequipArmor(u, it)))
	case 13:
		return uint64(uint32(equipmentDequipArmor(u, it, int(value), int(side))))
	case 14:
		return uint64(uint32(equipmentNPCEquipArmor(u, it)))
	case 15:
		equipmentRemoveWeapons(u)
	case 16:
		return uint64(uint32(equipmentEquipArmor(u, it, int(value), int(side))))
	case 17:
		return uint64(uint32(uintptr(unsafe.Pointer(equipmentSameArmor(u, it)))))
	case 18:
		equipmentArmorDropSound(it)
	case 19:
		equipmentInitDropTable()
	case 20:
		return uint64(uint32(equipmentDropPolicy(it, int(value))))
	case 21:
		return uint64(uint32(uintptr(equipmentNPCSync(u, it, int(value)))))
	case 22:
		return uint64(uint32(equipmentCount(u, int(value))))
	case 23:
		return uint64(uint32(equipmentDuplicate(u, it)))
	case 24:
		return uint64(uint32(equipmentTryEquip(u, it)))
	case 25:
		return uint64(uint32(equipmentTryDequip(u, it)))
	case 26:
		return uint64(uint32(equipmentEffects(it, u, true)))
	case 27:
		return uint64(uint32(equipmentEffects(it, u, false)))
	case 28:
		return uint64(bool2int(equipmentCheckStrength(u, it)))
	case 29:
		equipmentSaveShield(u)
	case 30:
		return uint64(uint32(uintptr(unsafe.Pointer(equipmentFindShield(u)))))
	case 31:
		return math.Float64bits(equipmentDefend(it))
	case 32:
		return uint64(uint32(equipmentStrength(u)))
	}
	return 0
}

func (p *portTestShopPools) equipmentAction(a PortTestShopAction) uint32 {
	sp := p.proxy.callbacks.shop.spec
	u := p.resources.unit
	if sp.Equipment.NilUnit {
		u = nil
	}
	var it *server.Object
	if !sp.Inventory.NilItem && a.Item >= 0 {
		it = p.items[a.Item].u
	}
	if a.Op == 409 {
		equipmentSecondary(u, it)
		p.equipment.result = 0
	} else {
		p.equipment.result = equipmentFixtureCall(a.Op-400, u, it, int32(a.Value), int32(a.Side))
	}
	return uint32(p.equipment.result)
}
func (p *portTestShopPools) equipmentSnapshot() []uint32 {
	if p.equipment == nil {
		return nil
	}
	s := p.equipment
	out := []uint32{p.normalize(uint32(s.result)), uint32(s.result >> 32), uint32(gameex_flags), uint32(nox_cheat_allowall), uint32(qword_581450_9512), uint32(qword_581450_9512 >> 32)}
	trace := unsafe.Slice((*uint32)(unsafe.Pointer(C.eqTracePtr())), 4097)
	if trace[0] > 680 {
		panic("equipment callback trace overflow")
	}
	for _, v := range trace[:1+6*trace[0]] {
		out = append(out, p.normalize(v))
	}
	for _, b := range s.blocks {
		if !bytes.Equal(b[:8], bytes.Repeat([]byte{0xa5}, 8)) || !bytes.Equal(b[len(b)-8:], bytes.Repeat([]byte{0x5a}, 8)) {
			panic("equipment definition guard")
		}
		for i := 8; i < len(b)-8; i += 4 {
			out = append(out, p.normalize(*(*uint32)(unsafe.Pointer(&b[i]))))
		}
	}
	for _, v := range unsafe.Slice(memmap.PtrUint32(0x587000, 279432), 48) {
		out = append(out, p.normalize(v))
	}
	if u := p.resources.unit; u != nil {
		ad := p.proxy.core.Abils.GetFor(u)
		for _, v := range ad.Cooldowns {
			out = append(out, uint32(v))
		}
		for it := ad.ExecList; it != nil; it = it.Next {
			out = append(out, uint32(it.Abil), it.Frame, it.Active)
		}
	}
	return out
}

const PortTestEquipment53A030 = 400

const PortTestEquipment53A0F0 = 401

const PortTestEquipment53A140 = 402

const PortTestEquipment53A2C0 = 403

const PortTestEquipment53A3D0 = 404

const PortTestEquipment53A420 = 405

const PortTestEquipment53A680 = 406

const PortTestEquipment53A6C0 = 407

const PortTestEquipment53AAB0 = 408

const PortTestEquipment53AB90 = 409

const PortTestEquipment53E2D0 = 410

const PortTestEquipment53E300 = 411

const PortTestEquipment53E3A0 = 412

const PortTestEquipment53E430 = 413

const PortTestEquipment53E520 = 414

const PortTestEquipment53E600 = 415

const PortTestEquipment53E650 = 416

const PortTestEquipment53E7B0 = 417

const PortTestEquipment53EAE0 = 418

const PortTestEquipment53EC40 = 419

const PortTestEquipment53EC80 = 420

const PortTestEquipment4E4B20 = 421

const PortTestEquipment4E7D30 = 422

const PortTestEquipment4E7EC0 = 423

const PortTestEquipment4F2F70 = 424

const PortTestEquipment4F2FB0 = 425

const PortTestEquipment4F2FF0 = 426

const PortTestEquipment4F3030 = 427

const PortTestEquipment4F3180 = 428

const PortTestEquipment980523 = 429

const PortTestEquipment9805EB = 430

const PortTestEquipment415C00 = 431

const PortTestEquipment4F9FD0 = 432
