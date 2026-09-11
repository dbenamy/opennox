//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME5.h"
#include "GAME4_3.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
static void porttest_ai_original(int action, int mode, void* p) {
 int u = (int)p;
 if (mode == 3) { if (action < 4) nox_ai_action_pop_532100(u); return; }
 if (mode != 0) return;
 switch(action) {
 case 0: sub_545210(u); break;
 case 1: sub_545300(u); break;
 case 2: sub_545340(u); break;
 case 3: sub_5453E0(u); break;
 case 4: nox_xxx_mobActionRandomWalk_545020(u); break;
 case 5: nox_xxx_mobActionConfuse_545140(u); break;
 }
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"math"
	"slices"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestAIActionSpec struct {
	Action, Mode       int
	Seed               int
	Direction, Desired uint16
	X, Y, TX, TY       uint32
	Flags, Status      uint32
	Speed, Run, Melee  uint32
	Missile            byte
	Stack              int8
	NilTarget          bool
}
type PortTestAIActionResult struct {
	Changes                            []uint32
	Direction, Desired                 uint16
	LogicIndex, OtherIndex             int
	Stack                              int8
	StackChanged, GuardsOK, ReadOnlyOK bool
}

var portTestAITypes = [...]ai.ActionType{ai.ACTION_FACE_LOCATION, ai.ACTION_FACE_OBJECT, ai.ACTION_FACE_ANGLE, ai.ACTION_SET_ANGLE, ai.ACTION_RANDOM_WALK, ai.ACTION_CONFUSED}

func portTestAIDir(d int) (float32, float32) {
	if d == 0 {
		return math.Float32frombits(0xbff6495b), 0
	}
	x, y := float32(d%31-15)/8, float32((d*7)%29-14)/8
	if d%11 == 0 {
		x = 0.1
	}
	if d%13 == 0 {
		y = -0.1
	}
	return x, y
}

func PortTestAIActions(specs []PortTestAIActionSpec, registered bool) (out []PortTestAIActionResult, restored bool) {
	core := new(server.Server)
	core.SetFrame(123)
	oldGet, oldFlags := GetServer, noxflags.GetEngine()
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	noxflags.UnsetEngine(noxflags.EngineShowAI)
	defer func() { GetServer = oldGet; noxflags.ResetEngine(); noxflags.SetEngine(oldFlags) }()
	// Guarded C-owned storage for every object/state the real callbacks touch.
	makeBlock := func(n int) ([]byte, func()) { return alloc.Make([]byte{}, n+16) }
	ob, fo := makeBlock(int(unsafe.Sizeof(server.Object{})))
	defer fo()
	ub, fu := makeBlock(int(unsafe.Sizeof(server.MonsterUpdateData{})))
	defer fu()
	db, fd := makeBlock(int(unsafe.Sizeof(server.MonsterDef{})))
	defer fd()
	tb, ft := makeBlock(int(unsafe.Sizeof(server.Object{})))
	defer ft()
	obj := (*server.Object)(unsafe.Pointer(&ob[8]))
	ud := (*server.MonsterUpdateData)(unsafe.Pointer(&ub[8]))
	def := (*server.MonsterDef)(unsafe.Pointer(&db[8]))
	target := (*server.Object)(unsafe.Pointer(&tb[8]))
	detach := server.PortTestAttachAI(core, obj, target)
	defer detach()
	// One grid/table fixture is shared by every action in the batch.
	rows, fr := alloc.Make([]uint32{}, 130)
	defer fr()
	cells, fc := alloc.Make([]uint32{}, 128*(128*11+4))
	defer fc()
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
			cells[off+1], cells[off+6] = 3, 9
			if (x*17+y*13+1)%5 == 0 {
				cells[off+1] = 6
			}
			if (x*7+y*11+2)%7 < 3 {
				cells[off+6] = 6
			}
			cells[off+5], cells[off+10] = 0, 0
		}
	}
	wantRows, wantCells := slices.Clone(rows), slices.Clone(cells)
	oldGrid := C.ptr_5D4594_2650668
	installed := (**C.obj_5D4594_2650668_t)(unsafe.Pointer(&rows[1]))
	C.ptr_5D4594_2650668 = installed
	dirs := unsafe.Slice(memmap.PtrUint8(0x587000, 194128), 2064)
	oldDirs := bytes.Clone(dirs)
	for i := range dirs {
		dirs[i] = 0x59
	}
	for d := 0; d < 256; d++ {
		x, y := portTestAIDir(d)
		binary.LittleEndian.PutUint32(dirs[8+d*8:], math.Float32bits(x))
		binary.LittleEndian.PutUint32(dirs[12+d*8:], math.Float32bits(y))
	}
	wantDirs := bytes.Clone(dirs)
	defer func() {
		C.ptr_5D4594_2650668 = oldGrid
		copy(dirs, oldDirs)
		restored = C.ptr_5D4594_2650668 == oldGrid && bytes.Equal(dirs, oldDirs)
	}()
	for _, sp := range specs {
		// Preserve server registration while resetting all public object bytes.
		clear(ob[:8])
		clear(ob[8 : 8+772])
		clear(ob[len(ob)-8:])
		clear(tb[:8])
		clear(tb[8 : 8+772])
		clear(tb[len(tb)-8:])
		clear(ub)
		clear(db)
		for _, b := range [][]byte{ob, ub, db, tb} {
			for i := 0; i < 8; i++ {
				b[i] = byte(0xa0 + i)
				b[len(b)-8+i] = byte(0xc0 + i)
			}
		}
		obj.ObjClass = object.ClassMonster
		obj.ObjSubClass = object.SubClass(sp.Flags)
		obj.UpdateData = unsafe.Pointer(ud)
		obj.Direction1, obj.Direction2 = server.Dir16(sp.Direction), server.Dir16(sp.Desired)
		obj.PosVec = types.Pointf{X: math.Float32frombits(sp.X), Y: math.Float32frombits(sp.Y)}
		obj.SpeedCur = math.Float32frombits(sp.Speed)
		obj.ForceVec = types.Pointf{X: 123, Y: -456}
		target.PosVec = types.Pointf{X: math.Float32frombits(sp.TX), Y: math.Float32frombits(sp.TY)}
		ud.MonsterDef = def
		ud.StatusFlags = object.MonsterStatus(sp.Status)
		ud.AIStackInd = sp.Stack
		def.RunMultiplier96 = math.Float32frombits(sp.Run)
		def.MeleeAttackRange112 = math.Float32frombits(sp.Melee)
		def.MissileName148[0] = sp.Missile
		// Keep untouched move-audio code on its silent monster branch.
		def.MoveSndFrameA100, def.MoveSndFrameB104 = 250, 251
		ud.Field2, ud.Field67, ud.Field74, ud.Field91 = 91, 92, 93, 94
		ud.Field124, ud.Field137 = 95, 96
		for i := 0; i <= int(sp.Stack); i++ {
			ud.AIStack[i].Action = uint32(ai.ACTION_WAIT)
		}
		head := &ud.AIStack[sp.Stack]
		head.Action = uint32(portTestAITypes[sp.Action])
		head.Args[0], head.Args[1] = uintptr(sp.TX), uintptr(sp.TY)
		if sp.Action == 1 {
			head.Args[0] = uintptr(unsafe.Pointer(target))
			if sp.NilTarget {
				head.Args[0] = 0
			}
		}
		if sp.Action == 2 || sp.Action == 3 {
			head.Args[0] = uintptr(sp.TX)
		}
		core.Rand.Logic, core.Rand.Other = prand.New(sp.Seed), prand.New(sp.Seed+1)
		core.AI.StackChanged = false
		beforeO, beforeU := bytes.Clone(ob), bytes.Clone(ub)
		beforeD, beforeT := bytes.Clone(db), bytes.Clone(tb)
		if registered {
			a := server.GetAIAction(portTestAITypes[sp.Action])
			switch sp.Mode {
			case 0:
				a.Update(obj)
			case 1:
				a.Start(obj)
			case 2:
				a.End(obj)
			case 3:
				a.Cancel(obj)
			}
		} else {
			C.porttest_ai_original(C.int(sp.Action), C.int(sp.Mode), unsafe.Pointer(obj))
		}
		r := PortTestAIActionResult{Direction: uint16(obj.Direction1), Desired: uint16(obj.Direction2), LogicIndex: core.Rand.Logic.Index(), OtherIndex: core.Rand.Other.Index(), Stack: ud.AIStackInd, StackChanged: core.AI.StackChanged, GuardsOK: true, ReadOnlyOK: bytes.Equal(db, beforeD) && bytes.Equal(tb, beforeT)}
		for region, pair := range [][2][]byte{{beforeO, ob}, {beforeU, ub}} {
			for j := 8; j < len(pair[0])-8; j += 4 {
				v := binary.LittleEndian.Uint32(pair[1][j:])
				if v != binary.LittleEndian.Uint32(pair[0][j:]) {
					r.Changes = append(r.Changes, uint32(region*4096+j-8), v)
				}
			}
		}
		for _, b := range [][]byte{ob, ub, db, tb} {
			for i := 0; i < 8; i++ {
				r.GuardsOK = r.GuardsOK && b[i] == byte(0xa0+i) && b[len(b)-8+i] == byte(0xc0+i)
			}
		}
		r.ReadOnlyOK = r.ReadOnlyOK && C.ptr_5D4594_2650668 == installed
		out = append(out, r)
	}
	if !slices.Equal(rows, wantRows) || !slices.Equal(cells, wantCells) || !bytes.Equal(dirs, wantDirs) {
		panic("AI fixture grid/table mutated")
	}
	return out, false
}
