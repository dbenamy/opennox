//go:build porttest

package legacy

/*
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2489160;
*/
import "C"

import (
	"bytes"
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestAISpellSpec struct {
	Op, Spell, Mode int
	DurationEmpty   bool
	// ID 1 is the word at +1492. These are sparse physical input writes.
	Permissions                                [][2]uint32
	Definitions                                []server.PortTestSpellClassDef
	DurationSpell                              [2]uint32
	DurationOwner                              [2]bool
	SummonAllowed                              bool
	InversionRange                             float64
	Missile, MissileTarget                     bool
	TargetClass, TargetFlags, TargetSubclass   uint32
	TargetCur, TargetMax, SecondCur, SecondMax uint16
	NilTargetHealth, Second                    bool
	SecondPos                                  [2]uint32
	Accuracy                                   uint32
	TargetArg, SelfTarget                      bool
	ArgPos                                     [2]uint32
}
type PortTestAISpellResult struct {
	Return  uint32
	Globals [2]uint32
	Output  [2]uint32
	Intact  bool
}
type portTestAISpellState struct {
	guards                  func() bool
	definitions             func([]server.PortTestSpellClassDef)
	inversion               func(float64)
	duration                []server.DurSpell
	health                  []server.HealthData
	words                   []uint32
	beforeDur, beforeHealth []byte
	beforePermissions       []uint32
	spec                    *PortTestAISpellSpec
}

func portTestAISpellEnvironment(proxy *portTestRoamOwnerServer) func() {
	defs, freeDefs := proxy.core.PortTestAISpellDefs()
	inv, freeInv := proxy.core.PortTestAIInversionRange()
	type guardedDuration struct {
		Left   [2]uint32
		Values [2]server.DurSpell
		Right  [2]uint32
	}
	type guardedHealth struct {
		Left   [2]uint32
		Values [2]server.HealthData
		Right  [2]uint32
	}
	db, freeDur := alloc.New(guardedDuration{})
	hb, freeHealth := alloc.New(guardedHealth{})
	left, right := [2]uint32{0xa5a5a5a5, 0x5a5a5a5a}, [2]uint32{0x12345678, 0xabcdef01}
	db.Left, db.Right, hb.Left, hb.Right = left, right, left, right
	duration, health := db.Values[:], hb.Values[:]
	guards := func() bool { return db.Left == left && db.Right == right && hb.Left == left && hb.Right == right }
	words, freeWords := alloc.Make([]uint32{}, 9)
	oldList := proxy.core.Spells.Dur.List
	oldMissile, oldHeal := *memmap.PtrUint32(0x5D4594, 2489156), C.dword_5d4594_2489160
	oldSummon := Nox_xxx_checkSummonedCreaturesLimit_500D70
	st := &portTestAISpellState{guards: guards, definitions: defs, inversion: inv, duration: duration, health: health, words: words}
	proxy.spells = st
	Nox_xxx_checkSummonedCreaturesLimit_500D70 = func(u *server.Object, ind int) bool {
		proxy.trace = append(proxy.trace, 45, proxy.life.ids[uint32(uintptr(u.CObj()))], uint32(ind))
		return st.spec.SummonAllowed
	}
	return func() {
		Nox_xxx_checkSummonedCreaturesLimit_500D70 = oldSummon
		proxy.core.Spells.Dur.List = oldList
		*memmap.PtrUint32(0x5D4594, 2489156) = oldMissile
		C.dword_5d4594_2489160 = oldHeal
		freeWords()
		freeHealth()
		freeDur()
		freeInv()
		freeDefs()
		proxy.spells = nil
	}
}
func portTestAISpellPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestAISpellSpec) {
	st, ud := proxy.spells, u.UpdateDataMonster()
	st.spec = sp
	defs := make([]server.PortTestSpellClassDef, 136)
	for i := range defs {
		defs[i] = server.PortTestSpellClassDef{Index: uint32(i + 1), Valid: true}
	}
	defs = append(defs, sp.Definitions...)
	st.definitions(defs)
	st.inversion(sp.InversionRange)
	permission := unsafe.Slice((*uint32)(unsafe.Add(u.UpdateData, 1492)), 136)
	clear(permission)
	for _, w := range sp.Permissions {
		if w[0] < 1 || w[0] > 136 {
			panic("spell permission ID")
		}
		permission[w[0]-1] = w[1]
	}
	st.beforePermissions = append(st.beforePermissions[:0], permission...)
	ud.Field330 = math.Float32frombits(sp.Accuracy)
	for i := range st.duration {
		st.duration[i] = server.DurSpell{Spell: sp.DurationSpell[i], Caster16: proxy.combat.target}
		if sp.DurationOwner[i] {
			st.duration[i].Caster16 = u
		}
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(&st.duration[i])))] = uint32(900 + i)
	}
	st.duration[0].Next = &st.duration[1]
	proxy.core.Spells.Dur.List = &st.duration[0]
	if sp.DurationEmpty {
		proxy.core.Spells.Dur.List = nil
	}
	st.beforeDur = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&st.duration[0])), 2*int(unsafe.Sizeof(server.DurSpell{}))))
	*memmap.PtrUint32(0x5D4594, 2489156) = 0x12345678
	C.dword_5d4594_2489160 = 0
	for i := range st.words {
		st.words[i] = 0xa5a5a5a5
	}
	target := proxy.combat.target
	target.ObjClass = object.Class(sp.TargetClass)
	target.ObjFlags = object.Flags(sp.TargetFlags)
	target.ObjSubClass = object.SubClass(sp.TargetSubclass)
	st.health[0] = server.HealthData{Cur: sp.TargetCur, Max: sp.TargetMax}
	st.health[1] = server.HealthData{Cur: sp.SecondCur, Max: sp.SecondMax}
	target.HealthData = &st.health[0]
	if sp.NilTargetHealth {
		target.HealthData = nil
	}
	for i := range st.health {
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(&st.health[i])))] = uint32(950 + i)
	}
	if sp.Missile {
		target.ObjClass = object.ClassMissile
		target.ObjFlags |= object.FlagActive
		if sp.MissileTarget {
			*(*unsafe.Pointer)(unsafe.Add(target.UpdateData, 4)) = u.CObj()
		}
	}
	if sp.SelfTarget {
		target = u
	}
	st.words[1] = uint32(uintptr(target.CObj()))
	st.words[2] = sp.ArgPos[0]
	st.words[3] = sp.ArgPos[1]
	head := ud.AIStackHead()
	head.Args = [4]uintptr{uintptr(sp.Spell), 0, uintptr(sp.ArgPos[0]), uintptr(sp.ArgPos[1])}
	if sp.TargetArg {
		head.Args[2] = uintptr(target.CObj())
	}
	if proxy.combat.target.ObjFlags.Has(object.FlagActive) {
		proxy.combat.target.ObjFlags &^= object.FlagPartitioned
		proxy.combat.target.ObjIndexBase = server.ObjectIndex{}
		proxy.combat.target.ObjIndex = [4]server.ObjectIndex{}
		proxy.combat.target.ObjIndexCur = 0
		proxy.core.Map.AddObjectToIndex(proxy.combat.target)
	}
	if sp.Second {
		w := proxy.combat.weapon
		*w = *proxy.combat.target
		w.HealthData = &st.health[1]
		w.PosVec = types.Pointf{X: math.Float32frombits(sp.SecondPos[0]), Y: math.Float32frombits(sp.SecondPos[1])}
		w.NewPos = w.PosVec
		w.PrevPos = w.PosVec
		w.ObjFlags &^= object.FlagPartitioned
		w.ObjIndexBase = server.ObjectIndex{}
		w.ObjIndex = [4]server.ObjectIndex{}
		w.ObjIndexCur = 0
		proxy.core.Map.AddObjectToIndex(w)
	}
	st.beforeHealth = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&st.health[0])), 2*int(unsafe.Sizeof(server.HealthData{}))))
	for i, b := range proxy.combat.extra {
		proxy.combat.before[i] = bytes.Clone(b)
	}
}
func portTestAISpellCall(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestAISpellSpec) uint32 {
	t := proxy.combat.target
	if sp.SelfTarget {
		t = u
	}
	a := (*[3]uint32)(unsafe.Pointer(&proxy.spells.words[1]))
	out := (*types.Pointf)(unsafe.Pointer(&proxy.spells.words[5]))
	asWord := func(v bool) uint32 {
		if v {
			return 1
		}
		return 0
	}
	switch sp.Op {
	case 0:
		return asWord(monsterCastBusy(u))
	case 1:
		return asWord(monsterCastInversion(u))
	case 2:
		monsterMagicMissile(t, u)
	case 3:
		return asWord(monsterBuffSelf(u))
	case 4:
		return asWord(monsterSpellEnchantActive(u, spell.ID(sp.Spell)))
	case 5:
		return asWord(monsterSummonSpell(spell.ID(sp.Spell)))
	case 6:
		return asWord(monsterSummonActive(u))
	case 7:
		return asWord(monsterCastRelated2(u, t))
	case 8:
		return asWord(monsterCastOffensive(u, t))
	case 9:
		return asWord(monsterCastRelated(u))
	case 10:
		return asWord(monsterHealSomeone(u))
	case 11:
		monsterHealCandidate(t, u)
	case 12:
		monsterCastSpell(sp.Spell, u, a)
	case 13:
		monsterActionCast(u, sp.Mode)
	case 14:
		monsterCastRecoil(u, t, out)
	}

	return 0 // private char/pointer scratch results are unused by every caller.
}
func portTestAISpellTrace(proxy *portTestRoamOwnerServer, u *server.Object, rv uint32, normalize func(uint32) uint32) *PortTestAISpellResult {
	st := proxy.spells
	r := &PortTestAISpellResult{Return: rv, Globals: [2]uint32{*memmap.PtrUint32(0x5D4594, 2489156), normalize(uint32(C.dword_5d4594_2489160))}, Output: [2]uint32{st.words[5], st.words[6]}, Intact: true}
	r.Intact = st.guards() && bytes.Equal(st.beforeDur, unsafe.Slice((*byte)(unsafe.Pointer(&st.duration[0])), len(st.beforeDur))) && bytes.Equal(st.beforeHealth, unsafe.Slice((*byte)(unsafe.Pointer(&st.health[0])), len(st.beforeHealth)))
	for _, i := range []int{0, 4, 7, 8} {
		r.Intact = r.Intact && st.words[i] == 0xa5a5a5a5
	}
	permission := unsafe.Slice((*uint32)(unsafe.Add(u.UpdateData, 1492)), 136)
	for i, v := range permission {
		r.Intact = r.Intact && v == st.beforePermissions[i]
	}
	return r
}
