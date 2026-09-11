//go:build porttest

package legacy

/*
#include "GAME5.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestRoamSpec struct {
	Op, Seed                   int
	Index, Insert, Count, Mask byte
	Stack                      int8
	History                    [16]byte
	Neighbors                  [32]byte
	Flags                      [34]byte
	Enabled                    [34]bool
}
type PortTestRoamResult struct {
	History            [16]byte
	Index, Arg, Field2 uint32
	Return             int
	Stack              int8
	Logic, Other       int
	Changed            bool
	Changes            []uint32
	Intact             bool
}

// PortTestRoam runs a shared guarded fixture through the live C entries and Go owners.
// Pointer-bearing output is normalized to stable waypoint IDs before hashing.
func PortTestRoam(specs []PortTestRoamSpec) []PortTestRoamResult {
	core := new(server.Server)
	core.SetFrame(123)
	restoreTypes := core.PortTestObjectInitSize(1, 0)
	defer restoreTypes()
	oldGet, oldFlags := GetServer, noxflags.GetEngine()
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	noxflags.UnsetEngine(noxflags.EngineShowAI)
	defer func() { GetServer = oldGet; noxflags.ResetEngine(); noxflags.SetEngine(oldFlags) }()
	ob, fo := alloc.Make([]byte{}, int(unsafe.Sizeof(server.Object{}))+16)
	defer fo()
	ub, fu := alloc.Make([]byte{}, int(unsafe.Sizeof(server.MonsterUpdateData{}))+16)
	defer fu()
	stride := int(unsafe.Sizeof(server.Waypoint{})) + 16
	wb, fw := alloc.Make([]byte{}, 34*stride)
	defer fw()
	obj := (*server.Object)(unsafe.Pointer(&ob[8]))
	ud := (*server.MonsterUpdateData)(unsafe.Pointer(&ub[8]))
	detach := server.PortTestAttachAI(core, obj)
	defer detach()
	raw := func(id byte) uint32 {
		if id == 0 {
			return 0
		}
		return uint32(uintptr(unsafe.Pointer(&wb[int(id)*stride+8])))
	}
	ids := map[uint32]uint32{0: 0}
	for id := byte(1); id < 34; id++ {
		ids[raw(id)] = uint32(id)
	}
	normalize := func(v uint32) uint32 {
		if id, ok := ids[v]; ok {
			return id
		}
		return v
	}
	put := func(b []byte, off int, v uint32) { binary.LittleEndian.PutUint32(b[off:], v) }
	get := func(b []byte, off int) uint32 { return binary.LittleEndian.Uint32(b[off:]) }
	guard := func(b []byte) {
		for i := 0; i < 8; i++ {
			b[i] = 0xa5
			b[len(b)-8+i] = 0x5a
		}
	}
	intact := func(b []byte) bool {
		for i := 0; i < 8; i++ {
			if b[i] != 0xa5 || b[len(b)-8+i] != 0x5a {
				return false
			}
		}
		return true
	}
	out := make([]PortTestRoamResult, 0, len(specs))
	for _, sp := range specs {
		clear(ob[8:780])
		clear(ub)
		clear(wb)
		guard(ob)
		guard(ub)
		for i := 0; i < 34; i++ {
			guard(wb[i*stride : (i+1)*stride])
		}
		obj.TypeInd = 1
		obj.ObjClass = object.ClassMonster
		obj.UpdateData = unsafe.Pointer(ud)
		ud.AIStackInd = sp.Stack
		ud.Field2 = 99
		ud.Field91 = uint32(sp.Index)
		for i := 0; i <= int(sp.Stack); i++ {
			ud.AIStack[i].Action = uint32(ai.ACTION_WAIT)
		}
		head := &ud.AIStack[sp.Stack]
		head.Action = uint32(ai.ACTION_ROAM)
		head.Args[0] = uintptr(raw(33))
		head.Args[2] = uintptr(sp.Mask)
		for i, id := range sp.History {
			put(ub, 8+300+4*i, raw(id))
		}
		for id := byte(1); id < 34; id++ {
			b := wb[int(id)*stride+8 : (int(id)+1)*stride-8]
			b[477] = sp.Flags[id]
			if sp.Enabled[id] {
				put(b, 480, 1)
			}
		}
		root := wb[8 : stride-8]
		root[476] = sp.Count
		for i, id := range sp.Neighbors {
			put(root, 92+i*8, raw(id))
		}
		core.Rand.Logic, core.Rand.Other = prand.New(sp.Seed), prand.New(sp.Seed+1)
		core.AI.StackChanged = false
		beforeO, beforeU, beforeW := bytes.Clone(ob), bytes.Clone(ub), bytes.Clone(wb)
		op, up, wp := C.int(uintptr(unsafe.Pointer(obj))), C.int(uintptr(unsafe.Pointer(ud))), C.int(uintptr(unsafe.Pointer(&root[0])))
		ret := 0
		switch sp.Op {
		case 0:
			server.GetAIAction(ai.ACTION_ROAM).Start(obj)
		case 1:
			server.GetAIAction(ai.ACTION_ROAM).Cancel(obj)
		case 2:
			C.sub_545B00(up, C.int(raw(sp.Insert)))
		case 3:
			ret = int(normalize(uint32(C.sub_545B60(up, C.uchar(sp.Mask)))))
		case 4:
			ret = int(normalize(roamWaypointWord(roamSuccessor(ud, (*server.Waypoint)(unsafe.Pointer(&root[0])), sp.Mask))))
		case 5:
			ret = int(C.nox_xxx_monsterRoamDeadEnd_545BB0(op, wp))
		default:
			panic("invalid roam operation")
		}
		r := PortTestRoamResult{Index: ud.Field91, Arg: normalize(uint32(head.Args[0])), Field2: ud.Field2, Return: ret, Stack: ud.AIStackInd, Logic: core.Rand.Logic.Index(), Other: core.Rand.Other.Index(), Changed: core.AI.StackChanged, Intact: intact(ob) && intact(ub) && bytes.Equal(wb, beforeW)}
		for i := range r.History {
			r.History[i] = byte(normalize(get(ub, 8+300+4*i)))
		}
		for region, pair := range [][2][]byte{{beforeO, ob}, {beforeU, ub}} {
			for j := 8; j < len(pair[0])-8; j += 4 {
				a, b := get(pair[0], j), get(pair[1], j)
				if a != b {
					r.Changes = append(r.Changes, uint32(region*4096+j-8), normalize(b))
				}
			}
		}
		out = append(out, r)
	}
	return out
}
