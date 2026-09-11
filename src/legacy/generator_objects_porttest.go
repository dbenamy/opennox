//go:build porttest

package legacy

/*
#include "GAME5.h"
extern uint32_t dword_5d4594_1556136;
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

type PortTestGeneratorObjectsResult struct {
	Inventory                                                                                   *PortTestGeneratorInventoryResult `json:",omitempty"`
	Source, Destination, SourceData, DestinationData, Player, CreatedData, CreatedHealth, Spawn []uint32
	DeathFrame, QuestState                                                                      uint32
	Intact                                                                                      bool
}
type portTestGeneratorObjects struct {
	inventory                            *portTestGeneratorInventoryFixture
	beforePlayerObject, beforePlayerData []byte
	source, destination                  *server.Object
	raw                                  [][]byte
	resetSpawn                           func()
	snapshotSpawn                        func() portTestGeneratorSpawnAllocatorSnapshot
	beforePlayer                         []byte
}

func portTestGeneratorObjectsEnvironment(proxy *portTestRoamOwnerServer) func() {
	st := &portTestGeneratorObjects{}
	inv, freeInv := portTestGeneratorInventoryFixtureNew(proxy)
	st.inventory = inv
	var frees []func()
	for _, n := range []int{int(unsafe.Sizeof(server.Object{})), int(unsafe.Sizeof(server.Object{})), 2200, 2200} {
		b, free := alloc.Make([]byte{}, n+16)
		st.raw = append(st.raw, b)
		frees = append(frees, free)
	}
	st.source = (*server.Object)(unsafe.Pointer(&st.raw[0][8]))
	st.destination = (*server.Object)(unsafe.Pointer(&st.raw[1][8]))
	reset, snapshot, freeSpawn := portTestGeneratorSpawnAllocator()
	st.resetSpawn = reset
	st.snapshotSpawn = snapshot
	oldDeath := C.dword_5d4594_1556136
	oldQuest := *memmap.PtrUint32(0x5D4594, 1556120)
	oldScale := *memmap.PtrUint32(0x587000, 202036)
	oldCount := *memmap.PtrUint32(0x5D4594, 2386208)
	proxy.callbacks.generation.objects = st
	return func() {
		freeInv()
		freeSpawn()
		for i := len(frees) - 1; i >= 0; i-- {
			frees[i]()
		}
		C.dword_5d4594_1556136 = oldDeath
		*memmap.PtrUint32(0x5D4594, 1556120) = oldQuest
		*memmap.PtrUint32(0x587000, 202036) = oldScale
		*memmap.PtrUint32(0x5D4594, 2386208) = oldCount
	}
}
func portTestGeneratorObjectsPrepare(proxy *portTestRoamOwnerServer, u *server.Object, sp *PortTestGeneratorSpec) {
	st := proxy.callbacks.generation.objects
	st.resetSpawn()
	for _, b := range st.raw {
		clear(b)
		for i := 0; i < 8; i++ {
			b[i] = 0xa5
			b[len(b)-8+i] = 0x5a
		}
	}
	for i, t := range []*server.Object{st.source, st.destination} {
		*t = *u
		t.ObjClass = object.ClassMonster
		t.ObjSubClass = 0
		t.ObjFlags = 4
		t.TypeInd = 25
		t.InvFirstItem = nil
		t.InvHolder = nil
		t.ObjOwner = nil
		t.Obj130 = nil
		t.ObjNext = nil
		t.ObjPrev = nil
		t.InitData = nil
		t.UseData.Ptr = nil
		t.HealthData = nil
		t.UpdateData = unsafe.Pointer(&st.raw[2+i][8])
		t.Direction1 = 31
		t.Direction2 = 99
		proxy.life.ids[uint32(uintptr(t.CObj()))] = uint32(8001 + i)
		proxy.life.ids[uint32(uintptr(t.UpdateData))] = uint32(8101 + i)
	}
	if sp.Op == 6 || sp.Op == 8 {
		for _, t := range []*server.Object{st.source, st.destination} {
			*(*uint32)(unsafe.Add(t.UpdateData, 1504)) = sp.SourceSpellWord
			if sp.UseDef {
				def := u.UpdateDataMonster().MonsterDef
				def.HealthQuest72 = sp.DefHealth
				proxy.life.ids[uint32(uintptr(unsafe.Pointer(def)))] = 8401
				*(*unsafe.Pointer)(unsafe.Add(t.UpdateData, 484)) = unsafe.Pointer(def)
			}
		}
	}
	if sp.Beholder {
		st.source.TypeInd = 24
	}
	if sp.Op == 7 {
		for i := 8; i < len(st.raw[2])-8; i++ {
			st.raw[2][i] = byte(i*17 + int(sp.SourceFill))
		}
		for i := 8; i < len(st.raw[3])-8; i++ {
			st.raw[3][i] = 0xcc
		}
	}
	player := proxy.life.players[0].UpdateDataPlayer().Player
	po := &proxy.life.players[0]
	st.beforePlayerObject = bytes.Clone(unsafe.Slice((*byte)(po.CObj()), int(unsafe.Sizeof(*po))))
	st.beforePlayerData = bytes.Clone(unsafe.Slice((*byte)(po.UpdateData), int(unsafe.Sizeof(*po.UpdateDataPlayer()))))
	st.beforePlayer = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(player)), int(unsafe.Sizeof(*player))))
	proxy.life.ids[uint32(uintptr(unsafe.Pointer(player)))] = 991
	if sp.Op >= 4 {
		po.PosVec = types.Pointf{X: math.Float32frombits(sp.PlayerPos[0]), Y: math.Float32frombits(sp.PlayerPos[1])}
		po.ObjFlags = object.Flags(sp.PlayerFlags)
		*(*uint16)(unsafe.Add(unsafe.Pointer(player), 10)) = sp.View[0]
		*(*uint16)(unsafe.Add(unsafe.Pointer(player), 12)) = sp.View[1]
		*(*uint32)(unsafe.Add(unsafe.Pointer(player), 3680)) = sp.PlayerStatus
		*(*uint32)(unsafe.Add(unsafe.Pointer(player), 4792)) = uint32(bool2int(sp.Joined))
		clear(unsafe.Slice((*byte)(u.UpdateData), 164))
		*(*uint32)(unsafe.Add(u.UpdateData, 92)) = sp.Flags
		if sp.Killer {
			u.Obj130 = &proxy.life.players[0]
		} else {
			u.Obj130 = nil
		}
	}
	if sp.Op == 8 {
		if sp.Level > 2 || sp.Sources < 1 || sp.Sources > 4 {
			panic("invalid packed generator sources")
		}
		u.ObjFlags = object.Flags(sp.GenObjFlags)
		u.Field5 = sp.GenXStatus
		for i := byte(0); i < sp.Sources; i++ {
			t := st.source
			if i&1 != 0 {
				t = st.destination
				t.TypeInd = 24
			}
			*(**server.Object)(unsafe.Add(u.UpdateData, 16*sp.Level+4*uint32(i))) = t
		}
		*(*byte)(unsafe.Add(u.UpdateData, 80+sp.Level)) = sp.RateClass
		*(*byte)(unsafe.Add(u.UpdateData, 86)) = sp.Current
		*(*byte)(unsafe.Add(u.UpdateData, 87)) = sp.Limit
		*(*uint32)(unsafe.Add(u.UpdateData, 88)) = sp.LastSpawn
	}
	if sp.Inventory != nil {
		if sp.Op != 7 {
			panic("inventory initially requires direct copy")
		}
		st.inventory.prepare(st.source, st.destination, *sp.Inventory)
	}
	C.dword_5d4594_1556136 = 0x12345678
	*memmap.PtrUint32(0x5D4594, 1556120) = sp.QuestState
	*memmap.PtrUint32(0x587000, 202036) = sp.HealthScale
	*memmap.PtrUint32(0x5D4594, 2386208) = 0
}
func portTestGeneratorObjectsCall(proxy *portTestRoamOwnerServer, u *server.Object) uint32 {
	st := proxy.callbacks.generation
	sp := st.spec
	o := st.objects
	p := unsafe.Pointer(&st.point[2])
	switch sp.Op {
	case 4:
		C.nox_xxx_dieMonsterGen_54E630(combatPtr(u))
	case 5:
		return uint32(C.nox_xxx_mobGeneratorPick_54EBA0((*C.uint32_t)(u.CObj()), (*C.float2)(p), combatPtr(o.source)))
	case 6:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_mobGeneratorSpawn_54F070(combatPtr(u), C.int(uintptr(p)), combatPtr(o.source)))))
	case 8:
		rv := int8(C.nox_xxx_updateMonsterGenerator_54E930((*C.uint32_t)(u.CObj())))
		// The update ABI returns only the low byte of a newly allocated pointer.
		// Verify that transport before normalizing the allocator-dependent byte.
		if len(proxy.life.created) == 1 {
			child := proxy.life.created[0]
			if rv != int8(uintptr(child.CObj())) {
				panic("generator update lost spawn pointer return byte")
			}
			return 0xffff0001
		}
		return uint32(int32(rv))
	case 7:
		C.nox_xxx_unitCreatureCopyUC_54F2B0(combatPtr(o.source), combatPtr(o.destination))
	default:
		panic("generator object operation")
	}
	return 0
}
func portTestGeneratorObjectsTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestGeneratorObjectsResult {
	st := proxy.callbacks.generation.objects
	var inventory *PortTestGeneratorInventoryResult
	if proxy.callbacks.generation.spec.Inventory != nil {
		v := st.inventory.trace(normalize)
		inventory = &v
	}
	words := func(p unsafe.Pointer, n int) []uint32 {
		var v []uint32
		b := unsafe.Slice((*byte)(p), n)
		for i := 0; i < n; i += 4 {
			v = append(v, normalize(binary.LittleEndian.Uint32(b[i:])))
		}
		return v
	}
	r := &PortTestGeneratorObjectsResult{Inventory: inventory, Intact: inventory == nil || inventory.Intact, Source: words(st.source.CObj(), 772), Destination: words(st.destination.CObj(), 772), SourceData: words(st.source.UpdateData, 2200), DestinationData: words(st.destination.UpdateData, 2200), DeathFrame: uint32(C.dword_5d4594_1556136), QuestState: *memmap.PtrUint32(0x5D4594, 1556120)}
	for _, b := range st.raw {
		for i := 0; i < 8; i++ {
			r.Intact = r.Intact && b[i] == 0xa5 && b[len(b)-8+i] == 0x5a
		}
	}
	player := proxy.life.players[0].UpdateDataPlayer().Player
	r.Player = words(unsafe.Pointer(player), int(unsafe.Sizeof(*player)))
	copy(unsafe.Slice((*byte)(unsafe.Pointer(player)), len(st.beforePlayer)), st.beforePlayer)
	po := &proxy.life.players[0]
	copy(unsafe.Slice((*byte)(po.UpdateData), len(st.beforePlayerData)), st.beforePlayerData)
	copy(unsafe.Slice((*byte)(po.CObj()), len(st.beforePlayerObject)), st.beforePlayerObject)
	for _, n := range st.snapshotSpawn().Spawn {
		r.Spawn = append(r.Spawn, normalize(uint32(n.Object)), uint32(n.Prev), uint32(n.Next))
	}
	for i, t := range proxy.life.created {
		if t.HealthData != nil {
			proxy.life.ids[uint32(uintptr(unsafe.Pointer(t.HealthData)))] = uint32(8201 + i)
			r.CreatedHealth = append(r.CreatedHealth, words(unsafe.Pointer(t.HealthData), int(unsafe.Sizeof(*t.HealthData)))...)
		}
		if t.UpdateData != nil && t.Class().Has(object.ClassMonster) {
			ref := *(*unsafe.Pointer)(unsafe.Add(t.UpdateData, 2196))
			if ref != nil {
				proxy.life.ids[uint32(uintptr(ref))] = uint32(8301 + i)
			}
			r.CreatedData = append(r.CreatedData, words(t.UpdateData, 2200)...)
		}
	}
	return r
}
