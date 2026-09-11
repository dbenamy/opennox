//go:build porttest

package legacy

/*
#include "GAME5.h"
extern uint32_t dword_5d4594_2491580, dword_5d4594_2491588;
static uint32_t pt_callback_damage[1024];
static int pt_callback_damage_n;
static int pt_callback_mutate;
static uint32_t pt_callback_force_bits;
static int pt_callback_hit(uint32_t* t, uint32_t* a, uint32_t* w, int damage, int kind) {
 int i=pt_callback_damage_n;
 if(i+5<=1024) {pt_callback_damage[i]=(uint32_t)t;pt_callback_damage[i+1]=(uint32_t)a;pt_callback_damage[i+2]=(uint32_t)w;pt_callback_damage[i+3]=damage;pt_callback_damage[i+4]=kind;pt_callback_damage_n+=5;}
 if(pt_callback_mutate) {*(uint32_t*)(*(uint32_t*)(a[187]+484)+120)=pt_callback_force_bits;}
 return 1;
}
static void* pt_callback_hit_ptr(void) {return (void*)pt_callback_hit;}
static void pt_callback_reset(int mutate,uint32_t force) {pt_callback_damage_n=0;pt_callback_mutate=mutate;pt_callback_force_bits=force;}
static int pt_callback_count(void) {return pt_callback_damage_n;}
static uint32_t pt_callback_word(int i) {return pt_callback_damage[i];}
static void pt_callback_area(int target,int actor) {union {int i;float f;} u;u.i=actor;nox_xxx_monsterAttackAreaDamage_549860(target,u.f);}
static int pt_callback_golem(int actor) {union {int i;float f;} u;u.i=actor;return nox_xxx_sendEquakeAfterGolem_549800(u.f);}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

var PortTestCallbackServer func(*server.Server) (Server, func())

type PortTestAICallbackSpec struct {
	MutateOnDamage                                                         bool
	ForceAfterDamage                                                       uint32
	Op                                                                     int
	Range, Force, Damage, DamageType, PoisonChance, PoisonPower, PoisonMax uint32
	TargetBuffs, TargetPoison, TargetRadius, ActorRadius                   uint32
	Nearest, AllTargets                                                    uint32
	DebrisIndex, BoneIndex                                                 uint32
	Enabled                                                                map[string]bool
	CloudLifetime                                                          float64
	SelfTarget                                                             bool
	LootName                                                               string
}
type PortTestAICallbackResult struct {
	PlayerStatus                         uint32
	Modifiers                            []uint32
	Definition                           []uint32
	Return                               uint32
	Globals, Damage, Health, CreatedData []uint32
	Intact                               bool
}
type portTestAICallbackState struct {
	playerStatus uint32
	modifiers    server.PortTestAICallbackModifiers
	configure    func(map[string]bool)
	lifetime     func(float64)
	before       []byte
	spec         *PortTestAICallbackSpec
}

var portTestCallbackOffsets = []uintptr{2491556, 2491560, 2491564, 2491568, 2491572, 2491576, 2491584}

func portTestAICallbackEnvironment(proxy *portTestRoamOwnerServer) func() {
	adapter, freeAdapter := PortTestCallbackServer(proxy.core)
	oldAdapter := proxy.portTestRandomServer.Server
	proxy.portTestRandomServer.Server = adapter
	freeTables := portTestCallbackTablesEnvironment()
	_, configure, freeTypes := proxy.core.PortTestAICallbackTypes()
	mods, freeMods := proxy.core.PortTestAICallbackModifiers()
	lifetime, freeLifetime := proxy.core.PortTestSmallToxicCloudLifetime()
	proxy.callbacks = &portTestAICallbackState{configure: configure, lifetime: lifetime, modifiers: mods}
	old := make([]uint32, len(portTestCallbackOffsets))
	for i, off := range portTestCallbackOffsets {
		old[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	rot, bones := C.dword_5d4594_2491580, C.dword_5d4594_2491588
	regions := []struct {
		off  uintptr
		size int
	}{{287328, 24}, {287976, 36}, {288240, 8}, {288868, 8}}
	before := make([][]byte, len(regions))
	for i, r := range regions {
		before[i] = bytes.Clone(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.size))
	}
	a, freeA := alloc.CString("PortTestDebrisA")
	b, freeB := alloc.CString("PortTestDebrisB")
	for _, off := range []uintptr{287976, 288240, 288868} {
		*memmap.PtrPtr(0x587000, off) = unsafe.Pointer(a)
		*memmap.PtrPtr(0x587000, off+4) = unsafe.Pointer(b)
	}
	*memmap.PtrPtr(0x587000, 287984) = nil
	*memmap.PtrUint32(0x587000, 287344) = 2
	*memmap.PtrUint32(0x587000, 287348) = 2
	*memmap.PtrUint32(0x587000, 287332) = 0x00000302
	return func() {
		for i, r := range regions {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, r.off)), r.size), before[i])
		}
		for i, off := range portTestCallbackOffsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		C.dword_5d4594_2491580, C.dword_5d4594_2491588 = rot, bones
		freeB()
		freeA()
		freeLifetime()
		freeMods()
		freeTypes()
		freeTables()
		proxy.portTestRandomServer.Server = oldAdapter
		freeAdapter()
	}
}
func portTestAICallbackPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestAICallbackSpec) {
	st := proxy.callbacks
	st.playerStatus = *(*uint32)(unsafe.Add(unsafe.Pointer(proxy.life.players[0].UpdateDataPlayer().Player), 3680))
	st.spec = sp
	st.configure(sp.Enabled)
	st.lifetime(sp.CloudLifetime)
	C.pt_callback_reset(C.int(bool2int(sp.MutateOnDamage)), C.uint32_t(sp.ForceAfterDamage))
	for i, p := range []*server.ModifierEff{st.modifiers.WeaponPower1, st.modifiers.Material1, st.modifiers.Material2} {
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(p)))] = uint32(970 + i)
	}
	for _, off := range portTestCallbackOffsets {
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	*memmap.PtrUint32(0x5D4594, 2491568) = sp.AllTargets
	*memmap.PtrUint32(0x5D4594, 2491572) = sp.Nearest
	C.dword_5d4594_2491580 = C.uint32_t(sp.DebrisIndex)
	C.dword_5d4594_2491588 = C.uint32_t(sp.BoneIndex)
	*memmap.PtrFloat32(0x587000, 287328) = 10
	d := u.UpdateDataMonster().MonsterDef
	for i, p := range []unsafe.Pointer{d.MeleeStrikeFunc236, d.DieFunc228, d.DeadFunc232} {
		if p != nil {
			proxy.life.ids[uint32(uintptr(p))] = uint32(965 + i)
		}
	}
	for _, v := range [][2]uint32{{112, sp.Range}, {116, sp.Damage}, {120, sp.Force}, {124, sp.DamageType}, {136, sp.PoisonChance}, {140, sp.PoisonPower}, {144, sp.PoisonMax}} {
		*(*uint32)(unsafe.Add(unsafe.Pointer(d), uintptr(v[0]))) = v[1]
	}
	u.Shape.Circle.R = math.Float32frombits(sp.ActorRadius)
	u.Mass = 1
	for _, t := range []*server.Object{proxy.combat.target, proxy.combat.weapon} {
		t.Damage = C.pt_callback_hit_ptr()
		t.Buffs = sp.TargetBuffs
		t.Poison540 = byte(sp.TargetPoison)
		t.Mass = 1
		t.Shape.Circle.R = math.Float32frombits(sp.TargetRadius)
	}
	proxy.life.ids[uint32(uintptr(C.pt_callback_hit_ptr()))] = 960
	// Snapshot allowed target/extra writes only after all callback inputs are installed.
	for i, b := range proxy.combat.extra {
		proxy.combat.before[i] = bytes.Clone(b)
	}
}
func portTestAICallbackCall(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestAICallbackSpec) uint32 {
	t := proxy.combat.target
	if sp.SelfTarget {
		t = u
	}
	p, q := combatPtr(u), combatPtr(t)
	if sp.Op < 11 {
		return uint32(ccall.CallIntPtr(*memmap.PtrPtr(0x587000, 287100+uintptr(sp.Op)*8), u.CObj()))
	}
	if sp.Op < 20 {
		index := []int{0, 1, 2, 3, 4, 6, 7, 8, 9}[sp.Op-11]
		return uint32(ccall.CallIntPtr(*memmap.PtrPtr(0x587000, 287196+uintptr(index)*8), u.CObj()))
	}
	if sp.Op < 25 {
		return uint32(ccall.CallIntPtr(*memmap.PtrPtr(0x587000, 287284+uintptr(sp.Op-20)*8), u.CObj()))
	}
	switch sp.Op {
	case 25:
		C.sub_549270(q, p)
	case 26:
		return uint32(C.nox_xxx_monsterPickMeleeTarget_549440(p, C.int(sp.AllTargets)))
	case 27:
		C.sub_5494C0((*C.float)(t.CObj()), p)
	case 28:
		return uint32(C.sub_549690(p, q))
	case 29:
		return uint32(C.pt_callback_golem(p))
	case 30:
		C.pt_callback_area(q, p)
	case 31:
		name, free := alloc.CString(sp.LootName)
		defer free()
		C.sub_54A390(p, (*C.char)(unsafe.Pointer(name)), nil, nil, nil, nil, 5)
	case 32:
		C.sub_54A4C0(p)
	}
	return 0
}
func portTestAICallbackTrace(proxy *portTestRoamOwnerServer, rv uint32, normalize func(uint32) uint32) *PortTestAICallbackResult {
	r := &PortTestAICallbackResult{Return: normalize(rv), Intact: true}
	status := (*uint32)(unsafe.Add(unsafe.Pointer(proxy.life.players[0].UpdateDataPlayer().Player), 3680))
	r.PlayerStatus = *status
	*status = proxy.callbacks.playerStatus
	d := proxy.combat.actor.UpdateDataMonster().MonsterDef
	for off := uintptr(0); off < unsafe.Sizeof(*d); off += 4 {
		r.Definition = append(r.Definition, normalize(*(*uint32)(unsafe.Add(unsafe.Pointer(d), off))))
	}
	for _, off := range portTestCallbackOffsets {
		r.Globals = append(r.Globals, normalize(*memmap.PtrUint32(0x5D4594, off)))
	}
	r.Globals = append(r.Globals, uint32(C.dword_5d4594_2491580), uint32(C.dword_5d4594_2491588))
	for i := 0; i < int(C.pt_callback_count()); i++ {
		r.Damage = append(r.Damage, normalize(uint32(C.pt_callback_word(C.int(i)))))
	}
	for _, h := range proxy.spells.health {
		b := unsafe.Slice((*byte)(unsafe.Pointer(&h)), int(unsafe.Sizeof(h)))
		for i := 0; i < len(b); i += 4 {
			r.Health = append(r.Health, binary.LittleEndian.Uint32(b[i:]))
		}
	}
	// Poison is allowed to update health's timestamp; its full result is captured above.
	proxy.spells.beforeHealth = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&proxy.spells.health[0])), 2*int(unsafe.Sizeof(server.HealthData{}))))
	for _, u := range proxy.life.created {
		if u.InitData != nil {
			r.Modifiers = append(r.Modifiers, normalize(uint32(uintptr(u.CObj()))))
			for i := uintptr(0); i < 16; i += 4 {
				r.Modifiers = append(r.Modifiers, normalize(*(*uint32)(unsafe.Add(u.InitData, i))))
			}
			ammo := uint32(0)
			if u.UseData.Ptr != nil {
				ammo = uint32(*(*uint16)(u.UseData.Ptr))
			}
			r.Modifiers = append(r.Modifiers, ammo)
		}
		if u.UpdateData != nil {
			r.CreatedData = append(r.CreatedData, normalize(uint32(uintptr(u.CObj()))), *(*uint32)(u.UpdateData))
		}
	}
	return r
}

func (s *portTestRoamOwnerServer) ApplyForce(u *server.Object, p types.Pointf, force float64) {
	s.trace = append(s.trace, 47, s.life.ids[uint32(uintptr(u.CObj()))], math.Float32bits(p.X), math.Float32bits(p.Y), math.Float32bits(float32(force)), uint32(u.Poison540))
	s.portTestRandomServer.Server.ApplyForce(u, p, force)
}
