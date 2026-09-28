//go:build porttest

package legacy

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

var portTestLifecycleKeys [3]byte
var portTestLifecycleCalls [24]uint32
var portTestLifecycleCount int

func portTestLifecycleKey(i int) unsafe.Pointer { return unsafe.Pointer(&portTestLifecycleKeys[i]) }
func portTestLifecycleRecord(words ...uint32) {
	for _, w := range words {
		portTestLifecycleCalls[portTestLifecycleCount] = w
		portTestLifecycleCount++
	}
}
func portTestLifecycleDie(u *server.Object) int32 {
	portTestLifecycleRecord(1, *(*uint32)(unsafe.Add(unsafe.Pointer(u), 80)))
	return 1
}
func init() {
	monsterCallbackHandlers[portTestLifecycleKey(0)] = portTestLifecycleDie
	monsterCallbackHandlers[portTestLifecycleKey(1)] = func(u *server.Object) int32 {
		portTestLifecycleRecord(2)
		portTestLifecycleRecord(unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(u), 80)), 6)...)
		return 0
	}
	server.RegisterObjectUpdateCallbackGo(portTestLifecycleKey(0), func(u *server.Object) { portTestLifecycleDie(u) })
	server.PortTestRegisterUseCallback(portTestLifecycleKey(2), func(u, t *server.Object) int32 {
		portTestLifecycleRecord(3, uint32(uintptr(unsafe.Pointer(u))), uint32(uintptr(unsafe.Pointer(t))))
		return 1
	})
}

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
	oldDecay, oldHead, oldTail, oldSize, oldMask := motionDecayHead, dword_5d4594_1565512, dword_5d4594_1565516, dword_5d4594_1565520, dword_5d4594_2649712
	*memmap.PtrUint32(0x5D4594, 1565508) = 0
	dword_5d4594_1565512 = 0
	dword_5d4594_1565516 = 0
	dword_5d4594_2649712 = 2
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
		motionDecayHead, dword_5d4594_1565512, dword_5d4594_1565516, dword_5d4594_1565520, dword_5d4594_2649712 = oldDecay, oldHead, oldTail, oldSize, oldMask
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
	dword_5d4594_1565512 = 0
	dword_5d4594_1565516 = 0
	clear(unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565524)), 64))
	motionDecayHead = 0
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
		d.DieFunc228 = portTestLifecycleKey(0)
		d.DeadFunc232 = portTestLifecycleKey(1)
	}
	portTestLifecycleCount = 0
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
		t.Use.Ptr = portTestLifecycleKey(2)
		proxy.combat.weapon.Use.Ptr = portTestLifecycleKey(2)
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
		u.Update = portTestLifecycleKey(0)
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
	switch sp.Op {
	case 5:
		return int(lifecycleRaiseZombie(u))
	case 6:
		lifecycleReset(u)
	case 7:
		lifecycleReleasedSoul(u)
	case 8:
		lifecycleBurnDelete(u)
	case 9:
		return portTestLifecycleFoodSearch(u, math.Float32frombits(sp.Range), false)
	case 10:
		return portTestLifecycleFoodSearch(u, math.Float32frombits(sp.Range), true)
	}
	return 0
}
func portTestLifecycleTrace(proxy *portTestRoamOwnerServer, h *server.HealthData, normalize func(uint32) uint32) *PortTestLifecycleResult {
	r := &PortTestLifecycleResult{Updatable: normalize(uint32(uintptr(unsafe.Pointer(proxy.core.Objs.UpdatableList)))), Decay: normalize(uint32(motionDecayHead))}
	for i := 0; i < portTestLifecycleCount; i++ {
		r.Calls = append(r.Calls, normalize(portTestLifecycleCalls[i]))
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
		if proxy.callbacks != nil {
			if p.InitData != nil {
				alloc.FreePtr(p.InitData)
				p.InitData = nil
			}
			if p.UseData.Ptr != nil {
				alloc.FreePtr(p.UseData.Ptr)
				p.UseData.Ptr = nil
			}
		}
		if proxy.callbacks != nil && p.UpdateData != nil {
			alloc.FreePtr(p.UpdateData)
			p.UpdateData = nil
		}
		if proxy.callbacks != nil && (proxy.callbacks.generation != nil || proxy.callbacks.shop != nil && proxy.callbacks.shop.pools.spellEffectsActive()) && p.HealthData != nil {
			alloc.FreePtr(unsafe.Pointer(p.HealthData))
			p.HealthData = nil
		}

		if proxy.callbacks != nil && proxy.callbacks.shop != nil && proxy.callbacks.shop.spec.EffectsUse != nil && p.CollideData != nil {
			alloc.FreePtr(p.CollideData)
			p.CollideData = nil
		}
		if proxy.callbacks != nil && proxy.callbacks.shop != nil {
			sp := proxy.callbacks.shop.spec.TemporaryUpdates
			if sp != nil && sp.World != nil && sp.World.Objectives != nil && sp.World.Objectives.Attack != nil {
				// The full created prefix above retains the transferred inventory head.
				// Shop teardown already freed its items; created objects are freed separately.
				p.InvFirstItem = nil
			}
		}
		proxy.core.Objs.FreeObject(p)
	}
	if proxy.callbacks != nil && proxy.callbacks.shop != nil && proxy.callbacks.shop.spec.EffectsUse != nil {
		if proxy.core.Objs.Alive != proxy.callbacks.shop.pools.initialAlive {
			panic("effects fixture projectile leak")
		}
	}
	for p := uint32(dword_5d4594_1565512); p != 0; {
		b := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(p))), 416)
		packet := bytes.Clone(b[251 : 251+int(b[401])])
		if proxy.callbacks != nil && proxy.callbacks.shop != nil && proxy.callbacks.shop.spec != nil && proxy.callbacks.shop.spec.Engine != nil {
			portTestTradePacketDefined(packet)
		}
		r.Packets = append(r.Packets, packet)
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

// Preserve the original signed 32-bit pointer result conversion.
func portTestLifecycleFoodSearch(actor *server.Object, radius float32, second bool) int {
	return int(int32(uintptr(unsafe.Pointer(lifecycleFoodSearch(actor, radius, second)))))
}
