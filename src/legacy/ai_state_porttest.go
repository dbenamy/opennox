//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_3.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// PortTestMonsterStateSpec supplies raw input bits for one original-C helper.
// Op selects the following C entry point:
// 0 533790; 1..16 534220,280,2C0,300,320,340,390,3C0,400,440,470,710,750,
// 780,7A0,7C0; 17 7F0; 18 810; 19 840; 20 950; 21 A10; 22 A40.
// Type selects ordinary(0), Mimic(1), Plant(2), Zombie(3), or VileZombie(4).
// Direction is intentionally a signed raw 16-bit table index. All Float fields
// are IEEE float32 bits, except Frame/FPS and status fields which are raw uint32.
type PortTestMonsterStateSpec struct {
	Order                                                               int
	Source                                                              int // 0 nil, 1 ordinary owner, 2 player
	Broadcast, Own, Second, SecondEnabled, AnimData, NilHealth, NilUnit bool
	FrameA, FrameB, NetCode, FallbackAnim                               uint32
	Op, Type                                                            int
	Stack                                                               int8
	Action                                                              uint32
	Status, Subclass, ObjFlags, PlayerWeapon, PlayerShield              uint32
	Speed, Aggression, MeleeRange                                       uint32
	HealthCur, HealthMax                                                uint16
	Poison                                                              byte
	Frame, FPS, Deadline                                                uint32
	Direction                                                           int16
	Pos, Arg                                                            [2]uint32
	// For mimic: HeadArg is action argument position; MimicAge is update Field137.
	HeadArg, MimicAge uint32
	// Def/animation state consumed by capability and mimic helpers.
	MissileName byte
	// CacheSeed preserves a caller-selected existing cache value to test lazy
	// lookup versus reuse. Zero causes the corresponding C lazy lookup.
	MimicCache, PlantCache, ZombieCache, VileZombieCache uint32
}

type PortTestMonsterStateResult struct {
	Return    uint64
	Caches    [4]uint32
	Animation []uint32
	Objects   [][]uint32
	Command   uint32
	Intact    bool
}

type portTestMonsterState struct {
	types        server.PortTestMonsterStateTypeIDs
	anim, before []byte
	source       *server.Object
	spec         *PortTestMonsterStateSpec
}

func portTestMonsterStateEnvironment(proxy *portTestRoamOwnerServer) func() {
	ids, freeTypes := proxy.core.PortTestMonsterStateTypes()
	offsets := [...]uintptr{2488524, 2488528, 2488532, 2488536, 2487964, 2487968, 2487972, 2487976, 1570280}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	anim, freeAnim := alloc.Make([]byte{}, 272)
	weaponTable := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 215824)), 4*server.PlayerWeaponCnt)
	oldWeaponTable := bytes.Clone(weaponTable)
	for i := 0; i < server.PlayerWeaponCnt; i++ {
		binary.LittleEndian.PutUint32(weaponTable[4*i:], uint32(2+i))
	}
	oldObserve := Nox_xxx_playerObserveMonster_4DDE80
	Nox_xxx_playerObserveMonster_4DDE80 = func(owner, unit *server.Object) {
		proxy.trace = append(proxy.trace, 34, proxy.life.ids[uint32(uintptr(owner.CObj()))], proxy.life.ids[uint32(uintptr(unit.CObj()))])
	}
	proxy.state = &portTestMonsterState{types: ids, anim: anim}
	return func() {
		copy(weaponTable, oldWeaponTable)
		Nox_xxx_playerObserveMonster_4DDE80 = oldObserve
		freeAnim()
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		proxy.state = nil
		freeTypes()
	}
}

func portTestMonsterStatePrepare(proxy *portTestRoamOwnerServer, u *server.Object, h *server.HealthData, sp *PortTestMonsterStateSpec) {
	st, ud := proxy.state, u.UpdateDataMonster()
	st.spec = sp
	clear(st.anim)
	for i := 0; i < 8; i++ {
		st.anim[i] = 0xa5
		st.anim[len(st.anim)-8+i] = 0x5a
	}
	for i := 0; i < 64; i++ {
		binary.LittleEndian.PutUint32(st.anim[8+4*i:], uint32(i+100))
	}
	st.before = bytes.Clone(st.anim)
	for i := 0; i < 16; i++ {
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(&st.anim[8+16*i])))] = uint32(700 + i)
	}
	for _, off := range []uintptr{2487964, 2487968, 2487972, 2487976} {
		*memmap.PtrUint32(0x5D4594, off) = 0xabcdef01
	}
	*memmap.PtrUint32(0x5D4594, 1570280) = 0

	// Type cache inputs are raw on purpose; tests can seed a mismatching cache.
	for i, v := range [...]uint32{sp.MimicCache, sp.PlantCache, sp.ZombieCache, sp.VileZombieCache} {
		*memmap.PtrUint32(0x5D4594, [...]uintptr{2488524, 2488528, 2488532, 2488536}[i]) = v
	}
	u.TypeInd = 1
	switch sp.Type {
	case 1:
		u.TypeInd = uint16(st.types.Mimic)
	case 2:
		u.TypeInd = uint16(st.types.Plant)
	case 3:
		u.TypeInd = uint16(st.types.Zombie)
	case 4:
		u.TypeInd = uint16(st.types.VileZombie)
	}
	u.ObjClass = object.ClassMonster
	u.ObjSubClass = object.SubClass(sp.Subclass)
	u.ObjFlags = object.Flags(sp.ObjFlags)
	u.Poison540 = sp.Poison
	u.PosVec.X, u.PosVec.Y = math.Float32frombits(sp.Pos[0]), math.Float32frombits(sp.Pos[1])
	u.SpeedBase = math.Float32frombits(sp.Speed)
	u.Direction1 = server.Dir16(sp.Direction)
	u.Direction2 = server.Dir16(sp.Direction)
	u.HealthData = h
	if sp.NilHealth {
		u.HealthData = nil
	}
	h.Cur, h.Max = sp.HealthCur, sp.HealthMax
	ud.AIStackInd = sp.Stack
	if sp.Stack < 0 {
		ud.AIStackInd = -1
	} else {
		ud.AIStack[sp.Stack].Action = sp.Action
	}
	head := ud.AIStackHead()
	if head != nil {
		head.Args = [4]uintptr{uintptr(sp.Arg[0]), uintptr(sp.Arg[1]), 0, 0}
	}
	ud.StatusFlags = object.MonsterStatus(sp.Status)
	ud.Aggression = math.Float32frombits(sp.Aggression)
	ud.Field137 = sp.MimicAge
	ud.Field514 = sp.PlayerWeapon
	ud.Field515 = sp.PlayerShield
	d := ud.MonsterDef
	d.MeleeAttackRange112 = math.Float32frombits(sp.MeleeRange)
	d.MissileName148[0] = sp.MissileName
	d.MoveSndFrameA100 = sp.FrameA
	d.MoveSndFrameB104 = sp.FrameB
	u.NetCode = sp.NetCode
	ud.Field517 = sp.FallbackAnim
	if sp.AnimData {
		ud.Field119 = (*[16]server.MonsterAnim)(unsafe.Pointer(&st.anim[8]))
	}

	proxy.core.SetFrame(sp.Frame)
	proxy.core.SetTickRate(sp.FPS)
	// AI stack entries use raw action argument words. This exact write is needed
	// for morph's current-action position test; it is not an action construction oracle.
	if sp.Stack >= 0 {
		ud.AIStack[sp.Stack].Args[0] = uintptr(sp.Arg[0])
	}
	// 534810 reads Field127 (+508) as its deadline.
	ud.Field127 = sp.Deadline
	t, w := proxy.combat.target, proxy.combat.weapon
	st.source = nil
	if sp.Source == 1 {
		st.source = w
	} else if sp.Source == 2 {
		st.source = &proxy.life.players[0]
	}
	for i := range proxy.life.players {
		proxy.life.players[i].Field129 = nil
	}
	u.ObjOwner = nil
	if sp.Own {
		u.ObjOwner = st.source
	}
	if sp.Op == 23 {
		u.Field128 = nil
		if sp.Second {
			t.ObjClass = object.ClassMonster
			t.TypeInd = 1
			t.ObjFlags = 0
			t.ObjOwner = u.ObjOwner
			tu := t.UpdateDataMonster()
			tu.MonsterDef = d
			tu.SoundSet122 = ud.SoundSet122
			tu.AIStackInd = 1
			tu.AIStack[0].Action = 1
			tu.AIStack[1].Action = 1
			tu.StatusFlags = 0
			if sp.SecondEnabled {
				tu.StatusFlags = 0x80
			}
			t.SpeedBase = u.SpeedBase
			u.Field128 = t
		}
		if st.source != nil {
			st.source.Field129 = u
		}
	}
	proxy.life.configurePlayers(1)
	for i, b := range proxy.combat.extra {
		proxy.combat.before[i] = bytes.Clone(b)
	}

}

func portTestMonsterStateCall(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestMonsterStateSpec) uint64 {
	p := C.int(uintptr(u.CObj()))
	switch sp.Op {
	case 0:
		return uint64(C.nox_xxx_mobActionToAnimation_533790(p))
	case 1:
		return uint64(C.nox_xxx_monsterCanMelee_534220(p))
	case 2:
		return uint64(C.nox_xxx_monsterCanShoot_534280(p))
	case 3:
		return uint64(C.nox_xxx_monsterHasShield_5342C0(p))
	case 4:
		return uint64(C.nox_xxx_monsterCanCast_534300((*C.nox_object_t)(u.CObj())))
	case 5:
		return uint64(C.nox_xxx_monsterIsMoveing_534320(p))
	case 6:
		return uint64(C.sub_534340(p))
	case 7:
		return uint64(C.nox_xxx_monsterCanAttackAtWill_534390((*C.nox_object_t)(u.CObj())))
	case 8:
		return uint64(C.sub_5343C0(p))
	case 9:
		return uint64(C.sub_534400(p))
	case 10:
		return uint64(C.sub_534440(p))
	case 11:
		return math.Float64bits(float64(C.sub_534470(p)))
	case 12:
		return uint64(C.sub_534710(p))
	case 13:
		return uint64(uint32(C.sub_534750(p)))
	case 14:
		return uint64(uint32(C.sub_534780(p)))
	case 15:
		return uint64(C.sub_5347A0((*C.nox_object_t)(u.CObj())))
	case 16:
		return uint64(C.sub_5347C0(p))
	case 17:
		return uint64(C.nox_xxx_isNotPoisoned_5347F0(p))
	case 18:
		return uint64(C.nox_xxx_mobGetMoveAttemptTime_534810((*C.nox_object_t)(u.CObj())))
	case 19:
		return uint64(C.nox_xxx_unitIsMimic_534840(p))
	case 20:
		C.nox_xxx_monsterMimicCheckMorph_534950((*C.nox_object_t)(u.CObj()))
	case 21:
		return uint64(C.nox_xxx_unitIsPlant_534A10(p))
	case 22:
		return uint64(C.nox_xxx_unitIsZombie_534A40(p))
	case 23:
		target := u
		if sp.Broadcast || sp.NilUnit {
			target = nil
		}
		C.nox_xxx_orderUnit_533900(asObjectC(proxy.state.source), asObjectC(target), C.int(sp.Order))
	case 24:
		return uint64(uintptr(unsafe.Pointer(C.nox_xxx_unitNPCActionToAnim_533D00(p))))
	case 25:
		C.nox_xxx_monsterMoveAudio_534030(p)
	case 26, 27:
		point, free := alloc.New(types.Pointf{})
		*point = types.Pointf{X: math.Float32frombits(sp.Arg[0]), Y: math.Float32frombits(sp.Arg[1])}
		defer free()
		if sp.Op == 26 {
			return uint64(C.sub_534120(p, (*C.float2)(unsafe.Pointer(point))))
		}
		if sp.NilUnit {
			p = 0
		}
		C.nox_xxx_mobCalcDir_533CC0(p, (*C.float)(unsafe.Pointer(point)))
	default:
		panic("invalid monster-state operation")
	}
	return 0
}

func portTestMonsterStateTrace(proxy *portTestRoamOwnerServer, u *server.Object, ret uint64, normalize func(uint32) uint32) *PortTestMonsterStateResult {
	r := &PortTestMonsterStateResult{Return: ret, Intact: bytes.Equal(proxy.state.anim, proxy.state.before)}
	if proxy.state.spec.Op == 24 {
		if uint32(ret) == uint32(uintptr(memmap.PtrOff(0x5D4594, 2487964))) {
			r.Return = 800
		} else {
			r.Return = uint64(normalize(uint32(ret)))
		}
	}
	for _, off := range []uintptr{2487964, 2487968, 2487972, 2487976, 1570280} {
		r.Animation = append(r.Animation, *memmap.PtrUint32(0x5D4594, off))
	}
	for _, obj := range []*server.Object{proxy.combat.target, proxy.combat.weapon} {
		words := make([]uint32, 193)
		for i := range words {
			words[i] = normalize(*(*uint32)(unsafe.Add(obj.CObj(), 4*i)))
		}
		r.Objects = append(r.Objects, words)
	}
	pl := proxy.life.players[0].UpdateDataPlayer().Player
	r.Command = pl.SummonOrderAll
	pl.SummonOrderAll = 0

	for i, off := range [...]uintptr{2488524, 2488528, 2488532, 2488536} {
		r.Caches[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	return r
}
