//go:build porttest

package legacy

/*
#include "GAME4_1.h"
extern uint32_t dword_5d4594_2490504;
*/
import "C"

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestPathSpec struct {
	Graph                                                    bool
	Op, Wall                                                 int
	Path                                                     [32][2]uint32
	Count, Index, Status, Frame, ObjFlags, Speed, Multiplier uint32
	Target, Cache, WaypointCache                             [2]uint32
	WaypointCount, WaypointIndex                             uint32
	Waypoints                                                [16]byte
	Start, End                                               byte
}

func portTestPathPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestPathSpec, raw func(byte) uint32) {
	ud := u.UpdateDataMonster()
	u.ObjFlags = object.Flags(sp.ObjFlags)
	u.SpeedCur = math.Float32frombits(sp.Speed)
	ud.MonsterDef.RunMultiplier96 = math.Float32frombits(sp.Multiplier)
	ud.Field2, ud.Field67, ud.Field71, ud.Field70 = sp.Count, sp.Index, sp.Status, sp.Frame
	for i, p := range sp.Path {
		ud.Path[i] = types.Pointf{X: math.Float32frombits(p[0]), Y: math.Float32frombits(p[1])}
	}
	ud.Field68 = types.Pointf{X: math.Float32frombits(sp.Cache[0]), Y: math.Float32frombits(sp.Cache[1])}
	ud.Field92, ud.Field93 = sp.WaypointCache[0], sp.WaypointCache[1]
	ud.Field74, ud.Field91 = sp.WaypointCount, sp.WaypointIndex
	words := unsafe.Slice(&ud.Field75, 16)
	for i, id := range sp.Waypoints {
		words[i] = raw(id)
	}
	head := ud.AIStackHead()
	head.Action = uint32(ai.ACTION_MOVE_TO)
	head.Args[0], head.Args[1], head.Args[2] = uintptr(sp.Target[0]), uintptr(sp.Target[1]), 0
	*memmap.PtrUint32(0x5D4594, 2386204) = 0x12345678
	*memmap.PtrUint32(0x5D4594, 2490500) = 0x24681357
	proxy.trace = append(proxy.trace, 5, uint32(bool2int(proxy.core.MapTraceRayAt(types.Pointf{X: 100, Y: 100}, types.Pointf{X: 200, Y: 100}, nil, nil, 132))))
	proxy.pathEndpoints = [2]*server.Waypoint{roamWaypoint(raw(sp.Start)), roamWaypoint(raw(sp.End))}
	proxy.endpointCall = 0
	C.dword_5d4594_2490504 = 1
	scratch := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2489476), 256)
	for i := range scratch {
		scratch[i] = 0xab000000 + uint32(i)
	}
	if sp.Graph {
		for id := byte(1); id < 34; id++ {
			wp := roamWaypoint(raw(id))
			wp.Flags = 1
			wp.Flags2 = 128
			wp.PointsCnt = 0
			if id < 33 {
				wp.PointsCnt = 1
				wp.Points[0].Waypoint = roamWaypoint(raw(id + 1))
			}
		}
	}

}
func portTestPathCall(u *server.Object, sp *PortTestPathSpec) int {
	switch sp.Op {
	case 0:
		return int(C.nox_xxx_creatureActuallyMove_50D3B0((*C.float)(u.CObj())))
	case 1:
		return int(C.nox_xxx_creatureSetMovePath_50D5A0(C.int(uintptr(u.CObj()))))
	case 2:
		return int(C.sub_50D2E0(C.int(uintptr(u.CObj()))))
	case 3:
		return int(C.sub_50D2A0(C.int(uintptr(u.CObj())), C.int(uintptr(unsafe.Pointer(&u.UpdateDataMonster().AIStackHead().Args[0])))))
	}
	panic("invalid path operation")
}

func portTestPathEnvironment() func() {
	epsilon := memmap.PtrFloat64(0x581450, 10288)
	oldEpsilon := *epsilon
	// Runtime blob_581450.dat stores float32(0.01) promoted to double here.
	*epsilon = float64(float32(.01))
	oldEpoch := C.dword_5d4594_2490504
	scratch := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2489476), 256)
	old := append([]uint32(nil), scratch...)
	return func() {
		unchanged := *epsilon == float64(float32(.01))
		*epsilon = oldEpsilon
		C.dword_5d4594_2490504 = oldEpoch
		copy(scratch, old)
		if !unchanged {
			panic("path execution changed read-only epsilon")
		}
	}
}
func portTestPathGraphState(norm func(uint32) uint32) []uint32 {
	out := []uint32{uint32(C.dword_5d4594_2490504)}
	for i, v := range unsafe.Slice(memmap.PtrUint32(0x5D4594, 2489476), 256) {
		if v != 0xab000000+uint32(i) {
			out = append(out, uint32(i), norm(v))
		}
	}
	return out
}
