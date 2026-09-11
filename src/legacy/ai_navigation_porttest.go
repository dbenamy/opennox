//go:build porttest

package legacy

/*
#include "GAME5.h"
*/
import "C"

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestNavigationSpec struct {
	Op, Phase                                               int
	Speed, Multiplier, TX, TY, Follow, Resume, FleeRange    uint32
	Cur, Max                                                uint16
	Previous                                                uint32
	PathCount, PathIndex, PathStatus, PathFrame, RetryFrame uint32
	OneShot, GameFlags                                      uint32
	Generator                                               int
	NoOwner, Food, TargetArg                                bool
}

var portTestNavigationActions = [...]ai.ActionType{ai.ACTION_MOVE_TO, ai.ACTION_FAR_MOVE_TO, ai.ACTION_DODGE, ai.ACTION_FLEE, ai.ACTION_MOVE_TO_HOME, ai.ACTION_RETREAT, ai.ACTION_RETREAT_TO_MASTER}

func portTestNavigationPrepare(proxy *portTestRoamOwnerServer, u, target *server.Object, health *server.HealthData, sp *PortTestNavigationSpec) {
	ud := u.UpdateDataMonster()
	u.HealthData = health
	health.Cur, health.Max = sp.Cur, sp.Max
	u.SpeedCur = math.Float32frombits(sp.Speed)
	u.Frame134 = 0
	u.ObjOwner = target
	if sp.NoOwner {
		u.ObjOwner = nil
	}
	target.PosVec = types.Pointf{X: math.Float32frombits(sp.TX), Y: math.Float32frombits(sp.TY)}
	ud.MonsterDef.RunMultiplier96 = math.Float32frombits(sp.Multiplier)
	ud.Field329 = math.Float32frombits(sp.Follow)
	ud.ResumeLevel = math.Float32frombits(sp.Resume)
	ud.FleeRange = math.Float32frombits(sp.FleeRange)
	ud.Field2, ud.Field67, ud.Field71, ud.Field70, ud.Field135 = sp.PathCount, sp.PathIndex, sp.PathStatus, sp.PathFrame, sp.RetryFrame
	for i := range ud.Path {
		ud.Path[i] = u.PosVec
	}
	head := ud.AIStackHead()
	head.Action = uint32(ai.ACTION_RETREAT)
	if sp.Op < len(portTestNavigationActions) {
		head.Action = uint32(portTestNavigationActions[sp.Op])
	}
	head.Args[0], head.Args[1], head.Args[2] = uintptr(sp.TX), uintptr(sp.TY), 0
	if sp.TargetArg {
		head.Args[2] = uintptr(unsafe.Pointer(target))
	}
	if ud.AIStackInd > 0 {
		ud.AIStack[ud.AIStackInd-1].Action = sp.Previous
	}
	*memmap.PtrUint32(0x5D4594, 2490500) = sp.OneShot
	*memmap.PtrUint32(0x5D4594, 2489452) = 0
	*memmap.PtrUint32(0x5D4594, 2489444) = 0
	noxflags.ResetGame()
	noxflags.SetGame(noxflags.GameFlag(sp.GameFlags))
	proxy.generator = sp.Generator
	if sp.Food {
		target.ObjClass = object.ClassFood
		target.ObjFlags = object.FlagActive
		target.NewPos = target.PosVec
		proxy.core.Map.AddObjectToIndex(target)
	}
}
func portTestNavigationCall(u *server.Object, sp *PortTestNavigationSpec) uint32 {
	if sp.Op < len(portTestNavigationActions) {
		a := server.GetAIAction(portTestNavigationActions[sp.Op])
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
	case 7:
		return uint32(C.nox_xxx_monsterCanResumeAttack_545520(C.int(uintptr(u.CObj()))))
	case 8:
		return uint32(C.sub_545580(C.int(uintptr(u.CObj()))))
	case 9:
		return uint32(C.nox_xxx_monsterCanCast2_5455B0(C.int(uintptr(u.CObj()))))
	case 10:
		C.nox_xxx_mobRetreatCheckEdibles_5455E0(C.int(uintptr(u.CObj())))
	default:
		panic("invalid navigation operation")
	}
	return 0
}
func (s *portTestRoamOwnerServer) Nox_xxx_generateRetreatPath_50CA00(path []types.Pointf, u *server.Object, p *types.Pointf) int {
	s.trace = append(s.trace, 4, uint32(len(path)), math.Float32bits(p.X), math.Float32bits(p.Y))
	for i := 0; i < min(len(path), max(0, s.generator)); i++ {
		path[i] = u.PosVec
	}
	return s.generator
}
