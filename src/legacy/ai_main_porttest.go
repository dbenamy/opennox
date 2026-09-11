//go:build porttest

package legacy

/*
#include "GAME5.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2489460;
extern uint32_t dword_5d4594_2491592;
*/
import "C"

import (
	"bytes"
	"math"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// The raw writes describe C layout inputs, not a Go decision oracle. Offsets
// address monster update data; action words are ordered bottom to top.
type PortTestMainSpec struct {
	TargetPlayer                            bool
	Started                                 bool
	Nearest                                 uint32
	Op                                      int // idle, dangerous, unwind, main owner, dodge, shield owner, shield candidate
	Words                                   [][2]uint32
	Actions                                 []uint32
	Class, Flags2, DefFlags                 uint32
	Tile                                    uint32
	Host, CursorObject                      bool
	Cursor                                  uint32
	Mouse                                   [2]int32
	PlayerBusy                              [2]uint32
	TargetType, TargetSubclass              uint32
	Missile, SecondMissile                  bool
	TargetNewPos, SecondPos, SecondVelocity [2]uint32
	CacheCloud, CacheSmall, Danger          uint32
}
type PortTestMainResult struct {
	Return  uint32
	Globals []uint32
	Intact  bool
}
type portTestMainState struct {
	host       func(*server.Object)
	grid       func(uint32)
	gridIntact func() bool
}

func portTestMainEnvironment(proxy *portTestRoamOwnerServer) func() {
	_, freeTypes := proxy.core.PortTestMainTypes()
	host, freeHost := proxy.core.PortTestMainHost()
	grid, gridIntact, freeGrid := portTestMainGrid()
	offsets := []uintptr{2489468, 2489472, 1096672, 2491596, 2491600, 2491604}
	old := make([]uint32, len(offsets))
	for i, o := range offsets {
		old[i] = *memmap.PtrUint32(0x5D4594, o)
	}
	danger, cursorOwner := C.dword_5d4594_2489460, C.dword_5d4594_2491592
	oldCast := Nox_xxx_castSpellByUser_4FDD20
	Nox_xxx_castSpellByUser_4FDD20 = func(spell int, u *server.Object, args unsafe.Pointer) int {
		words := unsafe.Slice((*uint32)(args), 3)
		proxy.trace = append(proxy.trace, 40, uint32(spell), proxy.life.ids[uint32(uintptr(u.CObj()))])
		for _, v := range words {
			if id, ok := proxy.life.ids[v]; ok {
				v = id
			}
			proxy.trace = append(proxy.trace, v)
		}
		return 1
	}
	proxy.main = &portTestMainState{host: host, grid: grid, gridIntact: gridIntact}
	return func() {
		Nox_xxx_castSpellByUser_4FDD20 = oldCast
		for i, o := range offsets {
			*memmap.PtrUint32(0x5D4594, o) = old[i]
		}
		C.dword_5d4594_2489460, C.dword_5d4594_2491592 = danger, cursorOwner
		freeGrid()
		freeHost()
		freeTypes()
		proxy.main = nil
	}
}
func portTestMainPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestMainSpec) {
	ud := u.UpdateDataMonster()
	for _, w := range sp.Words {
		if w[0]%4 != 0 || w[0] >= uint32(unsafe.Sizeof(*ud)) {
			panic("invalid main fixture word")
		}
		*(*uint32)(unsafe.Add(u.UpdateData, w[0])) = w[1]
	}
	if len(sp.Actions) != 0 {
		ud.AIStackInd = int8(len(sp.Actions) - 1)
		for i, a := range sp.Actions {
			ud.AIStack[i].Action = a
		}
	}
	if sp.Started {
		ud.AIStackHead().Field5 = 1
	}
	if sp.Op == 6 {
		*memmap.PtrUint32(0x5D4594, 2487956) = 0
		*memmap.PtrUint32(0x5D4594, 2487988) = sp.Nearest
	}
	u.ObjClass = object.Class(sp.Class)
	u.Field5 = sp.Flags2
	ud.MonsterDef.StatusFlags92 = object.MonsterStatus(sp.DefFlags)
	proxy.main.grid(sp.Tile)
	for i, v := range []uint32{sp.CacheCloud, sp.CacheSmall, sp.Cursor, 0, 0, 0} {
		*memmap.PtrUint32(0x5D4594, []uintptr{2489468, 2489472, 1096672, 2491596, 2491600, 2491604}[i]) = v
	}
	C.dword_5d4594_2489460 = C.uint32_t(sp.Danger)
	C.dword_5d4594_2491592 = 0
	p := &proxy.life.players[0]
	pi := p.UpdateDataPlayer().Player
	*(*int32)(unsafe.Add(unsafe.Pointer(pi), 2284)) = sp.Mouse[0]
	*(*int32)(unsafe.Add(unsafe.Pointer(pi), 2288)) = sp.Mouse[1]
	*(*uint32)(unsafe.Add(p.UpdateData, 280)) = sp.PlayerBusy[0]
	*(*uint32)(unsafe.Add(p.UpdateData, 284)) = sp.PlayerBusy[1]
	p.ObjFlags = 0
	proxy.main.host(nil)
	if sp.Host {
		proxy.main.host(p)
	}
	if sp.CursorObject {
		u.ObjFlags |= object.FlagActive
		u.NewPos = u.PosVec
		u.Shape.Kind = server.ShapeKindCircle
		u.Shape.Circle.R = 5
		u.Shape.Circle.R2 = 25
		proxy.core.Map.AddObjectToIndex(u)
	}
	t, w := proxy.combat.target, proxy.combat.weapon
	t.TypeInd = uint16(sp.TargetType)
	if sp.TargetPlayer {
		t.UpdateData = p.UpdateData
		proxy.life.ids[uint32(uintptr(p.UpdateData))] = 850
	}
	t.ObjSubClass = object.SubClass(sp.TargetSubclass)
	t.PrevPos = types.Pointf{X: math.Float32frombits(sp.TargetNewPos[0]), Y: math.Float32frombits(sp.TargetNewPos[1])}
	if sp.Missile {
		t.ObjFlags |= object.FlagActive
		t.ObjClass = object.ClassMissile
		proxy.core.Map.AddObjectToIndex(t)
	}
	if sp.SecondMissile {
		*w = *t
		w.PosVec = types.Pointf{X: math.Float32frombits(sp.SecondPos[0]), Y: math.Float32frombits(sp.SecondPos[1])}
		w.PrevPos = w.PosVec
		w.NewPos = w.PosVec
		w.ObjFlags &^= object.FlagPartitioned
		w.VelVec = types.Pointf{X: math.Float32frombits(sp.SecondVelocity[0]), Y: math.Float32frombits(sp.SecondVelocity[1])}
		w.ObjIndex = [4]server.ObjectIndex{}
		w.ObjIndexBase = server.ObjectIndex{}
		w.ObjIndexCur = 0 // second spatial-index node is independent
		proxy.core.Map.AddObjectToIndex(w)
	}
	for i, v := range []*server.Object{u, t, w} {
		if v.Use.Ptr != nil {
			proxy.life.ids[uint32(uintptr(v.Use.Ptr))] = 881
		}
		cb := *(*uint32)(unsafe.Add(v.CObj(), 696))
		if cb != 0 {
			proxy.life.ids[cb] = 880
		}
		proxy.life.ids[uint32(uintptr(unsafe.Pointer(&v.ObjIndexBase)))] = uint32(810 + i*8)
		for j := range v.ObjIndex {
			proxy.life.ids[uint32(uintptr(unsafe.Pointer(&v.ObjIndex[j])))] = uint32(811 + i*8 + j)
		}
	}
	proxy.life.configurePlayers(1)
	for i, b := range proxy.combat.extra {
		proxy.combat.before[i] = bytes.Clone(b)
	}
}
func portTestMainCall(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestMainSpec) uint32 {
	switch sp.Op {
	case 0:
		C.nox_xxx_mobAction_5469B0(asObjectC(u))
	case 1:
		return uint32(int32(C.nox_xxx_unitIsDangerous_547120(asObjectC(proxy.combat.target), asObjectC(u))))
	case 2:
		C.nox_xxx_monsterPopAttackActions_5471B0(C.int(uintptr(u.CObj())))
	case 3:
		C.nox_xxx_monsterMainAIFn_547210(asObjectC(u))
	case 4:
		return uint32(C.nox_xxx_monsterCheckDodgeables_547C50(C.int(uintptr(u.CObj()))))
	case 5:
		return uint32(C.nox_xxx_monsterTestBlockShield_533E70(asObjectC(u)))
	case 6:
		C.sub_533EB0(C.int(uintptr(proxy.combat.target.CObj())), C.int(uintptr(u.CObj())))
	}
	return 0
}
func portTestMainTrace(proxy *portTestRoamOwnerServer, rv uint32, normalize func(uint32) uint32) *PortTestMainResult {
	r := &PortTestMainResult{Return: normalize(rv), Intact: proxy.main.gridIntact()}
	for _, o := range []uintptr{2489468, 2489472, 1096672, 2491596, 2491600, 2491604} {
		r.Globals = append(r.Globals, normalize(*memmap.PtrUint32(0x5D4594, o)))
	}
	r.Globals = append(r.Globals, uint32(C.dword_5d4594_2489460), normalize(uint32(C.dword_5d4594_2491592)))
	r.Globals = append(r.Globals, normalize(*memmap.PtrUint32(0x5D4594, 2487956)), *memmap.PtrUint32(0x5D4594, 2487988))
	return r
}
