//go:build porttest

package legacy

/*
#include "GAME1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
*/
import "C"

import (
	"math"
	"slices"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestGuardEscortSpec struct {
	Players                                                             int
	ScriptName                                                          string
	Pending, UseTerrain, Water                                          bool
	Op                                                                  int
	Kind                                                                uint16
	Direction, Desired                                                  uint16
	NetCode, Sight, Follow, TX, TY                                      uint32
	Heard, HeardFrame                                                   uint32
	Precheck, Seen, SeenPlayer, Reaction, Previous, Unresolved, NoOwner bool
	Name                                                                string
}

var portTestGuardGrids [2]unsafe.Pointer

func portTestGuardEscortEnvironment() func() {
	mimic, plant := memmap.PtrUint32(0x5D4594, 2488524), memmap.PtrUint32(0x5D4594, 2488528)
	oldMimic, oldPlant := *mimic, *plant
	oldGrid, oldGrids := C.ptr_5D4594_2650668, portTestGuardGrids
	*mimic, *plant = 2, 3
	var blocks, wants [][]uint32
	var frees []func()
	for n, tile := range []uint32{3, 6} {
		rows, fr := alloc.Make([]uint32{}, 130)
		cells, fc := alloc.Make([]uint32{}, 128*(128*11+4))
		frees = append(frees, fr, fc)
		for i := range rows {
			rows[i] = 0x31313131
		}
		for i := range cells {
			cells[i] = 0x35353535
		}
		for x := 0; x < 128; x++ {
			rows[x+1] = uint32(uintptr(unsafe.Pointer(&cells[x*(128*11+4)+2])))
			for y := 0; y < 128; y++ {
				off := x*(128*11+4) + 2 + y*11
				cells[off+1] = tile
				cells[off+6] = tile
				cells[off+5] = 0
				cells[off+10] = 0
			}
		}
		portTestGuardGrids[n] = unsafe.Pointer(&rows[1])
		blocks = append(blocks, rows, cells)
		wants = append(wants, slices.Clone(rows), slices.Clone(cells))
	}
	return func() {
		ok := *mimic == 2 && *plant == 3
		for i, b := range blocks {
			ok = ok && slices.Equal(b, wants[i])
		}
		C.ptr_5D4594_2650668 = oldGrid
		portTestGuardGrids = oldGrids
		*mimic, *plant = oldMimic, oldPlant
		for _, f := range frees {
			f()
		}
		if !ok {
			panic("guard/escort fixture changed read-only grid/cache state")
		}
	}
}
func portTestGuardEscortPrepare(proxy *portTestRoamOwnerServer, u, target *server.Object, sp *PortTestGuardEscortSpec) {
	ud := u.UpdateDataMonster()
	u.TypeInd = sp.Kind
	u.ObjSubClass = object.SubClass(0x400)
	if sp.UseTerrain {
		u.ObjSubClass = 0
	}
	grid := 0
	if sp.Water {
		grid = 1
	}
	C.ptr_5D4594_2650668 = (**C.obj_5D4594_2650668_t)(portTestGuardGrids[grid])
	u.NetCode = sp.NetCode
	u.Direction1, u.Direction2 = server.Dir16(sp.Direction), server.Dir16(sp.Desired)
	if !sp.Reaction {
		u.Frame134 = 0
	}
	u.Pos132 = types.Pointf{X: 234, Y: 345}
	u.PrevPos = u.PosVec
	if sp.Previous {
		u.PrevPos.X += 40
	}
	ud.Field137 = 0
	target.PosVec = types.Pointf{X: math.Float32frombits(sp.TX), Y: math.Float32frombits(sp.TY)}
	target.ObjClass = object.ClassMonster
	if sp.SeenPlayer {
		target.ObjClass = object.ClassPlayer
	}
	ud.SightRange = math.Float32frombits(sp.Sight)
	ud.Field329 = math.Float32frombits(sp.Follow)
	ud.Field97 = sp.Heard
	ud.Field101 = sp.HeardFrame
	ud.Field99X = 200
	ud.Field99Y = 100
	if sp.Seen {
		ud.Field282_1 = 1
		ud.Field283 = uint32(uintptr(unsafe.Pointer(target)))
	}
	u.ObjOwner = target
	if sp.NoOwner {
		u.ObjOwner = nil
	}
	name := unsafe.Slice((*byte)(unsafe.Pointer(&ud.Field341)), 76)
	clear(name)
	copy(name, sp.Name)
	head := ud.AIStackHead()
	head.Action = uint32(ai.ACTION_GUARD)
	head.Args[0], head.Args[1], head.Args[2] = uintptr(math.Float32bits(100)), uintptr(math.Float32bits(100)), uintptr(sp.Desired)
	if sp.Op == 1 || sp.Op == 2 || sp.Op == 3 {
		head.Action = uint32(ai.ACTION_ESCORT)
		head.Args[2] = uintptr(unsafe.Pointer(target))
		if sp.Unresolved {
			head.Args[2] = 0
		}
	}
	proxy.precheck = sp.Precheck
}
func portTestGuardEscortCall(u *server.Object, sp *PortTestGuardEscortSpec) uint32 {
	switch sp.Op {
	case 0:
		server.GetAIAction(ai.ACTION_GUARD).Update(u)
	case 1:
		server.GetAIAction(ai.ACTION_ESCORT).Update(u)
	case 2:
		server.GetAIAction(ai.ACTION_ESCORT).End(u)
	case 3:
		server.GetAIAction(ai.ACTION_ESCORT).Cancel(u)
	case 4:
		return uint32(heardSoundAction(u))
	case 5:
		return uint32(investigateHeardSound(u))
	case 6:
		return uint32(uintptr(unsafe.Pointer(escortResolve(u))))
	case 7:
		return uint32(bool2int(u.MonsterLookAtDamager()))
	default:
		panic("invalid guard/escort operation")
	}
	return 0
}
