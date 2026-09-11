//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_3.h"
#include "GAME3_3.h"
extern uint32_t dword_5d4594_2386576;
extern uint32_t dword_5d4594_1565512, dword_5d4594_1565516, dword_5d4594_1565520;
extern uint32_t dword_5d4594_2649712;
static uint32_t pt_life_calls[24];
static int pt_life_count;
static int pt_life_die(uint32_t* u) {pt_life_calls[pt_life_count++]=1;pt_life_calls[pt_life_count++]=u[20];return 1;}
static void pt_life_dead(uint32_t* u) {pt_life_calls[pt_life_count++]=2;for(int i=20;i<26;i++)pt_life_calls[pt_life_count++]=u[i];}
static int pt_life_use(uint32_t* u,uint32_t* t) {pt_life_calls[pt_life_count++]=3;pt_life_calls[pt_life_count++]=(uint32_t)u;pt_life_calls[pt_life_count++]=(uint32_t)t;return 1;}
static void* pt_life_die_ptr(void) {return (void*)pt_life_die;}
static void* pt_life_dead_ptr(void) {return (void*)pt_life_dead;}
static void* pt_life_use_ptr(void) {return (void*)pt_life_use;}
static void pt_life_reset(void) {pt_life_count=0;}
static int pt_life_n(void) {return pt_life_count;}
static uint32_t pt_life_value(int i) {return pt_life_calls[i];}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/player"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

type PortTestLifecycleSpec struct {
	HeadAction                                                          uint32
	Second                                                              bool
	SecondPos                                                           [2]uint32
	Op, Phase, Kind                                                     int
	Callback, Soul, DeleteDef, PlayerOwner, Updatable, Place, Swap, Use bool
	HasTarget, Eligible                                                 bool
	ItemSubclass                                                        uint32
	Poison                                                              byte
	ObjectFlags, Subclass, GameFlags, Frame137, Duration, HistoryCount  uint32
	SeenCount                                                           byte
	Motion                                                              [6]uint32
	DelayMin, DelayMax                                                  float64
	Cur, Max                                                            uint16
	Range                                                               uint32
}
type PortTestLifecycleResult struct {
	Calls            []uint32
	Created          [][]uint32
	Packets          [][]byte
	Health           []uint32
	Updatable, Decay uint32
}
type portTestLifecycleState struct {
	ids              map[uint32]uint32
	types            server.PortTestLifecycleTypeIDs
	configurePlayers func(int)
	playersUnchanged func() bool
	players          []server.Object
	created          []*server.Object
	durationRestore  func()
	spec             *PortTestLifecycleSpec
}

var portTestLifecycleActions = [...]ai.ActionType{ai.ACTION_HUNT, ai.ACTION_PICKUP_OBJECT, ai.ACTION_DYING, ai.ACTION_DEAD, ai.ACTION_GET_UP}

func portTestLifecycleEnvironment(proxy *portTestRoamOwnerServer) func() {
	freeHandles := handles.PortTestInit()
	ids, freeTypes := proxy.core.PortTestLifecycleTypes()
	players, configure, unchanged, freePlayers := proxy.core.PortTestEscortPlayers()
	proxy.life = &portTestLifecycleState{types: ids, players: players, configurePlayers: configure, playersUnchanged: unchanged}
	oldNet := proxy.core.NetList
	proxy.core.NetList = netlist.New()
	proxy.core.NetList.Init()
	offsets := []uintptr{2488532, 2488536, 2489456, 2488636, 2489440, 2489448, 1565508}
	old := make([]uint32, len(offsets))
	for i, off := range offsets {
		old[i] = *memmap.PtrUint32(0x5D4594, off)
	}
	seq := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565524)), 64)
	oldSeq := bytes.Clone(seq)
	scorch := []uintptr{276828, 276836, 276844}
	oldScorch := make([]uint32, 3)
	for i, off := range scorch {
		oldScorch[i] = *memmap.PtrUint32(0x587000, off)
	}
	oldDecay, oldHead, oldTail, oldSize, oldMask := C.dword_5d4594_2386576, C.dword_5d4594_1565512, C.dword_5d4594_1565516, C.dword_5d4594_1565520, C.dword_5d4594_2649712
	*memmap.PtrUint32(0x5D4594, 1565508) = 0
	C.dword_5d4594_1565512 = 0
	C.dword_5d4594_1565516 = 0
	C.dword_5d4594_2649712 = 2
	oldEligible := Nox_xxx_playerClassCanUseItem_57B3D0
	Nox_xxx_playerClassCanUseItem_57B3D0 = func(t *server.Object, cl player.Class) bool {
		proxy.trace = append(proxy.trace, 33, uint32(cl))
		return proxy.life.spec.Eligible
	}
	oldPlace := Nox_xxx_inventoryServPlace_4F36F0
	Nox_xxx_inventoryServPlace_4F36F0 = func(u, t *server.Object, a, b int) bool {
		proxy.trace = append(proxy.trace, 30, uint32(bool2int(u == proxy.combat.actor)), uint32(bool2int(t == proxy.combat.target)), uint32(a), uint32(b))
		if proxy.life.spec.Swap {
			u.UpdateDataMonster().AIStackHead().Args[0] = uintptr(unsafe.Pointer(proxy.combat.weapon))
		}
		return proxy.life.spec.Place
	}
	return func() {
		if proxy.life.durationRestore != nil {
			proxy.life.durationRestore()
		}
		if p := *memmap.PtrPtr(0x5D4594, 1565508); p != nil {
			alloc.AsClass(p).Free()
		}
		Nox_xxx_playerClassCanUseItem_57B3D0 = oldEligible
		Nox_xxx_inventoryServPlace_4F36F0 = oldPlace
		for i, off := range offsets {
			*memmap.PtrUint32(0x5D4594, off) = old[i]
		}
		for i, off := range scorch {
			*memmap.PtrUint32(0x587000, off) = oldScorch[i]
		}
		copy(seq, oldSeq)
		C.dword_5d4594_2386576, C.dword_5d4594_1565512, C.dword_5d4594_1565516, C.dword_5d4594_1565520, C.dword_5d4594_2649712 = oldDecay, oldHead, oldTail, oldSize, oldMask
		proxy.core.NetList.Free()
		proxy.core.NetList = oldNet
		freePlayers()
		freeTypes()
		freeHandles()
	}
}
func portTestLifecyclePrepare(proxy *portTestRoamOwnerServer, u, t *server.Object, h *server.HealthData, sp *PortTestLifecycleSpec) {
	st := proxy.life
	st.spec = sp
	st.created = nil
	if st.durationRestore != nil {
		st.durationRestore()
	}
	st.durationRestore = proxy.core.PortTestZombieDeadDuration(sp.DelayMin, sp.DelayMax)
	st.configurePlayers(1)
	proxy.core.NetList.ResetAll()
	if p := *memmap.PtrPtr(0x5D4594, 1565508); p != nil {
		alloc.AsClass(p).FreeAllObjects()
	}
	C.dword_5d4594_1565512 = 0
	C.dword_5d4594_1565516 = 0
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565524)), 64))
	C.dword_5d4594_2386576 = 0
	*memmap.PtrUint32(0x5D4594, 2488532) = 0
	*memmap.PtrUint32(0x5D4594, 2488536) = 0
	*memmap.PtrUint32(0x5D4594, 2489456) = 0
	*memmap.PtrUint32(0x5D4594, 2488636) = 1
	for i, id := range []int{st.types.ScorchSmall, st.types.ScorchMedium, st.types.ScorchLarge} {
		*memmap.PtrUint32(0x587000, []uintptr{276828, 276836, 276844}[i]) = uint32(id)
	}
	*memmap.PtrUint32(0x5D4594, 2489452) = 0
	*memmap.PtrUint32(0x5D4594, 2489444) = 0
	*memmap.PtrUint32(0x5D4594, 2489440) = 0
	*memmap.PtrUint32(0x5D4594, 2489448) = 0
	ud := u.UpdateDataMonster()
	d := ud.MonsterDef
	u.TypeInd = 1
	if sp.Kind == 1 {
		u.TypeInd = uint16(st.types.Zombie)
	} else if sp.Kind == 2 {
		u.TypeInd = uint16(st.types.VileZombie)
	}
	u.ObjFlags = object.Flags(sp.ObjectFlags)
	u.ObjSubClass = object.SubClass(sp.Subclass)
	u.Poison540 = sp.Poison
	t.ObjSubClass = object.SubClass(sp.ItemSubclass)
	u.ObjOwner = nil
	if sp.PlayerOwner {
		u.ObjOwner = &st.players[0]
	}
	u.HealthData = h
	h.Cur = sp.Cur
	h.Max = sp.Max
	if sp.Soul {
		u.Field131 = 14
	}
	if sp.DeleteDef {
		d.StatusFlags92 |= 1
	}
	if sp.Callback {
		d.DieFunc228 = C.pt_life_die_ptr()
		d.DeadFunc232 = C.pt_life_dead_ptr()
	}
	C.pt_life_reset()
	for i, v := range sp.Motion {
		*(*uint32)(unsafe.Add(u.CObj(), 80+4*i)) = v
	}
	ud.Field137 = sp.Frame137
	ud.Field123 = sp.Duration
	ud.Field74 = sp.HistoryCount
	ud.Field282_1 = sp.SeenCount
	ud.Field523_2 = 0xaa
	for i := 0; i < 16; i++ {
		*(*uint32)(unsafe.Add(u.UpdateData, 300+4*i)) = 0xabba0000 + uint32(i)
	}
	for i := 0; i < 16; i++ {
		*(*uint32)(unsafe.Add(u.UpdateData, 1132+4*i)) = 0xbcde0000 + uint32(i)
	}
	head := ud.AIStackHead()
	head.Action = uint32(ai.ACTION_DEAD)
	head.Args = [4]uintptr{}
	if sp.Op < len(portTestLifecycleActions) {
		head.Action = uint32(portTestLifecycleActions[sp.Op])
	}
	if sp.HeadAction != 0 {
		head.Action = sp.HeadAction
	}
	if sp.HasTarget {
		head.Args[0] = uintptr(t.CObj())
	}
	if sp.Use {
		t.Use.Ptr = C.pt_life_use_ptr()
		proxy.combat.weapon.Use.Ptr = C.pt_life_use_ptr()
		proxy.combat.weapon.ObjSubClass = 0x10
	}
	if sp.Second {
		w := proxy.combat.weapon
		w.PosVec = types.Pointf{X: math.Float32frombits(sp.SecondPos[0]), Y: math.Float32frombits(sp.SecondPos[1])}
		w.NewPos = w.PosVec
		w.PrevPos = w.PosVec
		w.Shape.Kind = server.ShapeKindCenter
		w.ObjFlags = 4
		w.ObjClass = t.ObjClass
		w.ObjSubClass = t.ObjSubClass
		proxy.core.Map.AddObjectToIndex(w)
	}
	proxy.core.Objs.UpdatableList = nil
	if sp.Updatable {
		u.IsUpdatable = 1
		proxy.core.Objs.UpdatableList = u
		u.Update = C.pt_life_die_ptr()
	}
	noxflags.ResetGame()
	noxflags.SetGame(noxflags.GameFlag(sp.GameFlags))
}
func portTestLifecycleCall(u *server.Object, sp *PortTestLifecycleSpec) int {
	if sp.Op < len(portTestLifecycleActions) {
		a := server.GetAIAction(portTestLifecycleActions[sp.Op])
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
		return 0
	}
	p := C.int(uintptr(u.CObj()))
	switch sp.Op {
	case 5:
		return int(C.nox_xxx_mobRaiseZombie_534AB0(p))
	case 6:
		lifecycleReset(u)
	case 7:
		lifecycleReleasedSoul(u)
	case 8:
		lifecycleBurnDelete(u)
	case 9:
		return int(C.nox_xxx_mobSearchEdible_544A00(asObjectC(u), C.float(math.Float32frombits(sp.Range))))
	case 10:
		return int(C.sub_544AE0(p, C.float(math.Float32frombits(sp.Range))))
	}
	return 0
}
func portTestLifecycleTrace(proxy *portTestRoamOwnerServer, h *server.HealthData, normalize func(uint32) uint32) *PortTestLifecycleResult {
	r := &PortTestLifecycleResult{Updatable: normalize(uint32(uintptr(unsafe.Pointer(proxy.core.Objs.UpdatableList)))), Decay: normalize(uint32(C.dword_5d4594_2386576))}
	for i := 0; i < int(C.pt_life_n()); i++ {
		r.Calls = append(r.Calls, normalize(uint32(C.pt_life_value(C.int(i)))))
	}
	for _, p := range proxy.life.created {
		var words []uint32
		b := unsafe.Slice((*byte)(p.CObj()), 772)
		for i := 0; i < len(b); i += 4 {
			words = append(words, normalize(binary.LittleEndian.Uint32(b[i:])))
		}
		r.Created = append(r.Created, words)
		if p.Field189 != nil {
			alloc.FreePtr(p.Field189)
			p.Field189 = nil
		}
		proxy.core.Objs.FreeObject(p)
	}
	for p := uint32(C.dword_5d4594_1565512); p != 0; {
		b := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(p))), 416)
		r.Packets = append(r.Packets, bytes.Clone(b[251:251+int(b[401])]))
		p = binary.LittleEndian.Uint32(b[408:])
	}
	r.Packets = append(r.Packets, proxy.core.NetList.CopyPacketsA(1, netlist.Kind1))
	b := unsafe.Slice((*byte)(h.C()), int(unsafe.Sizeof(*h)))
	for i := 0; i < len(b); i += 4 {
		r.Health = append(r.Health, binary.LittleEndian.Uint32(b[i:]))
	}
	for _, off := range []uintptr{2488532, 2488536, 2489456, 2489440, 2489448, 2489452, 2489444} {
		proxy.trace = append(proxy.trace, uint32(off), normalize(*memmap.PtrUint32(0x5D4594, off)))
	}
	return r
}
