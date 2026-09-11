//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2487948;
extern uint32_t dword_587000_261388;
extern uint32_t dword_5d4594_1565628, dword_5d4594_1565632;
static int pt_combat_strikes, pt_combat_actor, pt_combat_result;
static int pt_combat_strike(int obj) {
 pt_combat_strikes++; pt_combat_actor=obj; return pt_combat_result;
}
static void* pt_combat_strike_ptr(void) { return (void*)pt_combat_strike; }
static void pt_combat_reset(int result) { pt_combat_strikes=0; pt_combat_actor=0; pt_combat_result=result; }
static int pt_combat_count(void) { return pt_combat_strikes; }
static int pt_combat_last(void) { return pt_combat_actor; }
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestCombatResult struct {
	Extra                               []uint32
	Intact                              bool
	Actions                             []server.AIStackItem
	Sounds                              []uint32
	Strikes                             int
	Projectile                          []uint32
	Selected, Nearest, Cooldown, Status uint32
	Stamina                             byte
}
type PortTestCombatSpec struct {
	Friendly, Shield, PlayerUpdate                                       bool
	Radius                                                               uint32
	Op, Phase, Strike, Wall                                              int
	Cooldown, ArgFrame, MeleeRange, MissileRange, AttackFrame            uint32
	DelayMin, DelayMax, TargetFlags, TargetClass, Weapon, WeaponSubclass uint32
	Target, Velocity                                                     [2]uint32
	Direction                                                            uint16
	Anim, Progress, Done, Stamina                                        byte
	Playerlike, Sound, Shoot, TargetArg, Scan, Killable                  bool
	ScanCode, EnemyCode, HeardCode, Heard                                uint32
}

var portTestCombatActions = [...]ai.ActionType{ai.ACTION_FIGHT, ai.ACTION_BLOCK_ATTACK, ai.ACTION_BLOCK_FINISH, ai.ACTION_WEAPON_BLOCK, ai.ACTION_MELEE_ATTACK, ai.ACTION_MISSILE_ATTACK}

type portTestCombatScriptBase = NoxScript

type portTestCombatState struct {
	portTestCombatScriptBase
	proxy                 *portTestRoamOwnerServer
	actor, target, weapon *server.Object
	sounds                []uint32
	extra                 [][]byte
	before                [][]byte
}

func (s *portTestRoamOwnerServer) NoxScriptC() NoxScript { return s.combat }
func (s *portTestCombatState) ScriptCallback(b *server.ScriptCallback, caller, trigger *server.Object, event server.ScriptEventType) unsafe.Pointer {
	if s.proxy.state != nil {
		id := func(u *server.Object) uint32 {
			if u == nil {
				return 0
			}
			return s.proxy.life.ids[uint32(uintptr(u.CObj()))]
		}
		s.proxy.trace = append(s.proxy.trace, 10, uint32(event), uint32(uintptr(unsafe.Pointer(b))-uintptr(trigger.UpdateData)), id(caller), id(trigger))
		return nil
	}
	s.proxy.trace = append(s.proxy.trace, 10, uint32(event), uint32(uintptr(unsafe.Pointer(b))-uintptr(s.actor.UpdateData)), uint32(bool2int(caller == s.target)), uint32(bool2int(trigger == s.actor)))
	return nil
}
func (s *portTestRoamOwnerServer) CreateObjectAt(obj, owner server.Obj, p types.Pointf) {
	u := server.ToObject(obj)
	if s.life != nil {
		s.life.created = append(s.life.created, u)
		s.life.ids[uint32(uintptr(u.CObj()))] = uint32(1000 + len(s.life.created))
		if s.callbacks != nil {
			if u.InitData != nil {
				s.life.ids[uint32(uintptr(u.InitData))] = uint32(4000 + len(s.life.created))
			}
			if u.UseData.Ptr != nil {
				s.life.ids[uint32(uintptr(u.UseData.Ptr))] = uint32(5000 + len(s.life.created))
			}
		}
		if s.callbacks != nil && u.UpdateData != nil {
			s.life.ids[uint32(uintptr(u.UpdateData))] = uint32(3000 + len(s.life.created))
		}
		if u.Field189 != nil {
			s.life.ids[uint32(uintptr(u.Field189))] = uint32(2000 + len(s.life.created))
		}
		s.trace = append(s.trace, 31, uint32(u.TypeInd), math.Float32bits(p.X), math.Float32bits(p.Y), uint32(bool2int(owner == nil)))
		u.PosVec = p
		return
	}
	s.trace = append(s.trace, 11, math.Float32bits(p.X), math.Float32bits(p.Y), uint32(bool2int(owner == s.combat.actor)))
	u.PosVec = p
	s.combatProjectile = u
}
func (s *portTestRoamOwnerServer) DelayedDelete(u *server.Object) {
	if s.life != nil {
		s.trace = append(s.trace, 32, s.life.ids[uint32(uintptr(u.CObj()))])
		return
	}
	s.trace = append(s.trace, 12)
	s.combatProjectile = u
}

func portTestCombatEnvironment(proxy *portTestRoamOwnerServer) func() {
	offsets := []uintptr{2487684, 2487944, 2487952, 2487956, 2487988, 1565652, 1565656, 1565636, 1567708, 1565640}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	freeProjectiles := proxy.core.PortTestCombatProjectileType("porttest-combat-projectile", 3)
	oldDX, oldDY := C.dword_5d4594_1565628, C.dword_5d4594_1565632
	oldPtr, oldRadius := C.dword_5d4594_2487948, C.dword_587000_261388
	oldTables := make(map[uintptr][]byte)
	for off, b := range blobdata.PortTestCombatTables() {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(b))
		oldTables[off] = bytes.Clone(dst)
		copy(dst, b)
	}
	sounds, freeSounds := alloc.Make([]uint32{}, 19)
	for i := range sounds {
		sounds[i] = uint32(300 + i)
	}
	weapon, freeWeapon := alloc.New(server.Object{})
	extra := make([][]byte, 3)
	frees := make([]func(), 3)
	for i, n := range []int{int(unsafe.Sizeof(server.MonsterUpdateData{})), 5000, 5000} {
		extra[i], frees[i] = alloc.Make([]byte{}, n+16)
	}
	proxy.combat = &portTestCombatState{extra: extra, proxy: proxy, sounds: sounds, weapon: weapon}
	return func() {
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		C.dword_5d4594_2487948, C.dword_587000_261388 = oldPtr, oldRadius
		C.dword_5d4594_1565628, C.dword_5d4594_1565632 = oldDX, oldDY
		for off, b := range oldTables {
			copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(b)), b)
		}
		freeSounds()
		freeWeapon()
		for _, f := range frees {
			f()
		}
		freeProjectiles()
	}
}
func portTestCombatPrepare(proxy *portTestRoamOwnerServer, u, target *server.Object, health *server.HealthData, sp *PortTestCombatSpec) {
	ud := u.UpdateDataMonster()
	def := ud.MonsterDef
	s := proxy.combat
	s.actor, s.target = u, target
	for _, b := range s.extra {
		clear(b)
		for i := 0; i < 8; i++ {
			b[i] = 0xa5
			b[len(b)-8+i] = 0x5a
		}
	}
	target.UpdateData = unsafe.Pointer(&s.extra[0][8])
	if sp.Friendly {
		(*server.MonsterUpdateData)(target.UpdateData).StatusFlags = object.MonStatusMorphed
	}
	if sp.PlayerUpdate {
		ud.Field516 = uint32(uintptr(unsafe.Pointer(s.weapon)))
		put := func(b []byte, off int, p unsafe.Pointer) {
			binary.LittleEndian.PutUint32(b[8+off:], uint32(uintptr(p)))
		}
		put(s.extra[1], 292, u.UpdateData)
		put(s.extra[1], 276, unsafe.Pointer(&s.extra[2][8]))
		put(s.extra[1], 104, unsafe.Pointer(s.weapon))
		*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 2180)) = unsafe.Pointer(&s.extra[1][8])
	}

	proxy.combatProjectile = nil
	proxy.core.PortTestCombatAudioReset()
	u.ObjFlags = object.FlagActive
	*(*uint32)(unsafe.Add(u.CObj(), 176)) = sp.Radius
	if sp.Playerlike {
		u.ObjSubClass = object.SubClass(0x10)
	}
	u.Direction1 = server.Dir16(sp.Direction)
	ud.Field128 = sp.Cooldown
	ud.Field120_1 = sp.Anim
	ud.Field120_2 = sp.Progress
	ud.Field120_3 = sp.Done
	ud.Field282_0 = sp.Stamina
	ud.Field514 = sp.Weapon
	*s.weapon = server.Object{}
	s.weapon.ObjSubClass = object.SubClass(sp.WeaponSubclass)
	if sp.WeaponSubclass != 0 {
		ud.Field516 = uint32(uintptr(unsafe.Pointer(s.weapon)))
	}
	if sp.Sound {
		ud.SoundSet122 = unsafe.Pointer(&s.sounds[0])
	}
	def.MeleeAttackRange112 = math.Float32frombits(sp.MeleeRange)
	def.MissileAttackRange212 = math.Float32frombits(sp.MissileRange)
	def.MeleeAttackFrame108 = sp.AttackFrame
	def.MissileAttackFrame216 = sp.AttackFrame
	def.MeleeAttackDelayMin128 = sp.DelayMin
	def.MeleeAttackDelayMax132 = sp.DelayMax
	def.MissileAttackDelayMin220 = sp.DelayMin
	def.MissileAttackDelayMax224 = sp.DelayMax
	if sp.Shoot {
		copy(def.MissileName148[:], "porttest-combat-projectile")
	}
	if sp.Strike >= 0 {
		def.MeleeStrikeFunc236 = C.pt_combat_strike_ptr()
	}
	C.pt_combat_reset(C.int(sp.Strike))
	target.PosVec = types.Pointf{X: math.Float32frombits(sp.Target[0]), Y: math.Float32frombits(sp.Target[1])}
	target.NewPos = target.PosVec
	target.PrevPos = target.PosVec
	target.VelVec = types.Pointf{X: math.Float32frombits(sp.Velocity[0]), Y: math.Float32frombits(sp.Velocity[1])}
	target.ObjFlags = object.Flags(sp.TargetFlags)
	target.ObjClass = object.Class(sp.TargetClass)
	target.NetCode = sp.ScanCode
	target.Shape.Kind = server.ShapeKindCenter
	if sp.Shield {
		*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 696)) = C.pt_combat_strike_ptr()
		*(*unsafe.Pointer)(unsafe.Add(target.CObj(), 696)) = C.pt_combat_strike_ptr()
		u.ObjFlags |= 0x80
	}
	if sp.Killable {
		target.HealthData = health
		health.Cur = 10
		health.Max = 10
	}
	ud.Field300 = sp.EnemyCode
	*(*uint32)(unsafe.Add(u.UpdateData, 392)) = sp.HeardCode
	*(*uint32)(unsafe.Add(u.UpdateData, 388)) = sp.Heard
	head := ud.AIStackHead()
	head.Action = uint32(ai.ACTION_FIGHT)
	if sp.Op < len(portTestCombatActions) {
		head.Action = uint32(portTestCombatActions[sp.Op])
	}
	head.Args = [4]uintptr{uintptr(sp.Target[0]), uintptr(sp.Target[1]), uintptr(sp.ArgFrame), 0}
	if sp.Op != 0 {
		head.Args[2] = 0
	}
	if sp.Op == 1 {
		head.Args[0] = uintptr(sp.ArgFrame)
	}
	if sp.TargetArg {
		head.Args[2] = uintptr(unsafe.Pointer(target))
	}
	if sp.Scan {
		proxy.core.Map.AddObjectToIndex(target)
	}
	for _, off := range []uintptr{2487684, 2487944, 2487952, 2487956, 2487988} {
		*memmap.PtrUint32(0x5D4594, off) = 0x12345678
	}
	C.dword_5d4594_2487948 = 0
	if sp.Op == 10 {
		*memmap.PtrUint32(0x5D4594, 2487952) = math.Float32bits(51)
	}
	C.dword_587000_261388 = C.uint32_t(math.Float32bits(50))
	s.before = nil
	for _, b := range s.extra {
		s.before = append(s.before, bytes.Clone(b))
	}
}
func portTestCombatCall(u *server.Object, sp *PortTestCombatSpec) {
	if sp.Op < len(portTestCombatActions) {
		a := server.GetAIAction(portTestCombatActions[sp.Op])
		switch sp.Phase {
		case 0:
			a.Update(u)
		case 1:
			a.Start(u)
		case 2:
			a.End(u)
		case 3:
			a.Cancel(u)
		}
		return
	}
	t := GetServer().(*portTestRoamOwnerServer).combat.target
	switch sp.Op {
	case 6:
		combatChase(u, t)
	case 7:
		combatMeleeChain(u, t)
	case 8:
		combatMissileChain(u, t)
	case 9:
		combatChoose(u, t)
	case 10:
		combatScan(t, u)
	}
}
func portTestCombatTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestCombatResult {
	ud := proxy.combat.actor.UpdateDataMonster()
	r := &PortTestCombatResult{Intact: true, Strikes: int(C.pt_combat_count()), Selected: normalize(uint32(C.dword_5d4594_2487948)), Nearest: *memmap.PtrUint32(0x5D4594, 2487952), Cooldown: ud.Field128, Status: uint32(ud.StatusFlags), Stamina: ud.Field282_0}
	r.Actions = append([]server.AIStackItem(nil), ud.AIStack[:ud.AIStackInd+1]...)
	for i := range r.Actions {
		for j, v := range r.Actions[i].Args {
			r.Actions[i].Args[j] = uintptr(normalize(uint32(v)))
		}
	}
	proxy.trace = append(proxy.trace, 13, uint32(C.pt_combat_count()), uint32(bool2int(uint32(C.pt_combat_last()) == uint32(uintptr(unsafe.Pointer(proxy.combat.actor))))))
	for _, off := range []uintptr{2487684, 2487944, 2487952, 2487956, 2487988} {
		proxy.trace = append(proxy.trace, uint32(off), normalize(*memmap.PtrUint32(0x5D4594, off)))
	}
	proxy.trace = append(proxy.trace, 2487948, normalize(uint32(C.dword_5d4594_2487948)))
	for _, ev := range proxy.core.PortTestCombatAudioSnapshot() {
		r.Sounds = append(r.Sounds, uint32(ev.ID))
		proxy.trace = append(proxy.trace, 14, uint32(ev.ID), uint32(bool2int(ev.Obj == proxy.combat.actor)), math.Float32bits(ev.Pos.X), math.Float32bits(ev.Pos.Y), uint32(ev.Kind), ev.Code, uint32(bool2int(ev.ByPos)))
	}
	if u := proxy.combatProjectile; u != nil {
		proxy.trace = append(proxy.trace, 15)
		b := unsafe.Slice((*byte)(unsafe.Pointer(u)), 772)
		for i := 0; i < len(b); i += 4 {
			v := normalize(binary.LittleEndian.Uint32(b[i:]))
			proxy.trace = append(proxy.trace, v)
			r.Projectile = append(r.Projectile, v)
		}
		proxy.core.Objs.FreeObject(u)
	}

	for region, b := range proxy.combat.extra {
		before := proxy.combat.before[region]
		if !bytes.Equal(b[:8], before[:8]) || !bytes.Equal(b[len(b)-8:], before[len(b)-8:]) {
			r.Intact = false
		}
		for off := 8; off < len(b)-8; off += 4 {
			v := binary.LittleEndian.Uint32(b[off:])
			if v != binary.LittleEndian.Uint32(before[off:]) {
				r.Extra = append(r.Extra, uint32(region*8192+off-8), normalize(v))
			}
		}
	}
	return r
}
