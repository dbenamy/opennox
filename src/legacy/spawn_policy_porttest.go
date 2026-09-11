//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME4_1.h"
extern void* nox_alloc_spawn_2386216;
extern void* nox_alloc_monsterList_2386220;
extern uint32_t dword_5d4594_2386212;
extern uint32_t dword_5d4594_2386224;
extern uint32_t dword_5d4594_2386228;
*/
import "C"

import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Add this optional field to PortTestGeneratorSpec during integration. Op 9
// dispatches this fixture; existing generator operations and hashes remain
// unchanged.
//
//	SpawnPolicy *PortTestGeneratorSpawnPolicySpec
//
// The outer fixture supplies real players, player update/info fields, map rays,
// and balance. This fixture owns only C objects, their 2200-byte update areas,
// allocator-list state, and all raw snapshots.
type PortTestGeneratorSpawnPolicyOp uint8

const (
	PortTestGeneratorSpawnRegister PortTestGeneratorSpawnPolicyOp = iota
	PortTestGeneratorSpawnRemove
	PortTestGeneratorSpawnNonZombieCleanup
	PortTestGeneratorSpawnFarCull
	PortTestGeneratorSpawnVisibleCull
	PortTestGeneratorSpawnAdmission
	PortTestGeneratorSpawnTick
	PortTestGeneratorSpawnGlyphCleanup
	PortTestGeneratorSpawnCandidate
)

type PortTestGeneratorSpawnPolicyRecord struct {
	// Setup is intentionally raw enough for the root oracle to define every C
	// gate. Object and update buffers are C-backed and guarded by this fixture.
	Indexed                bool
	Class, Subclass, Flags uint32
	TypeInd                uint16
	Pos                    [2]uint32
	// UpdateWords are copied at offsets 0,4,... before the operation. The
	// caller configures owner +2192, generator +2196, and any player-visible
	// fields needed by the selected operation.
	UpdateWords map[uintptr]uint32
}

type PortTestGeneratorSpawnPolicyAction struct {
	Op     PortTestGeneratorSpawnPolicyOp
	Record int
}

type PortTestGeneratorSpawnPolicySpec struct {
	ReserveMonster uint8
	Players        []PortTestSpawnPlayer
	Records        []PortTestGeneratorSpawnPolicyRecord // bounded to 100
	Registers      []int                                // initial 50E030 child records
	Actions        []PortTestGeneratorSpawnPolicyAction
	// Point is used by each admission action. The integration caller sets core
	// frame/FPS before a Tick action, so 50D890 sees its native modulo gates.
	Point [2]uint32
}

type portTestGeneratorSpawnPolicyNode struct {
	Object, Next, Prev uint32 // normalized: object IDs are 10000+i
}

type portTestGeneratorMonsterListNode struct {
	Object, Mask, Next, Prev uint32
	Distances                [32]uint32
}

type PortTestGeneratorSpawnPolicyResult struct {
	Calls          []uint32
	Glyphs         [][]uint32
	GlyphData      [][]uint32
	GeneratorCount byte
	Return         uint32
	Spawn          []portTestGeneratorSpawnPolicyNode
	MonsterList    []portTestGeneratorMonsterListNode
	Counters       [32]uint32
	Objects        [][]uint32 // exactly the 772-byte C ABI prefix per record
	Updates        [][]uint32
	// Heads are standalone C vardefs, not memmap storage.
	SpawnHead, MonsterListHead, MonsterListCount uint32
	Occupied, GlyphCache                         uint32
	GuardsOK                                     bool
}

const (
	portTestSpawnPolicyMaxRecords         = 100
	portTestSpawnPolicyObjectSnapshotSize = 772  // C object ABI; snapshot only this prefix
	portTestSpawnPolicyUpdateSize         = 2200 // monster update ABI
	portTestSpawnPolicyGuard              = 8
)

type portTestGeneratorSpawnPolicy struct {
	proxy *portTestRoamOwnerServer
	// Reuse the existing generator object fixture's real SpawnClass reset.
	resetSpawn func()

	objects, updates []byte
	objectStride     int
	records          []*server.Object
	objID, updateID  map[uintptr]uint32
	nodeID           map[uintptr]uint32
	oldIDs           map[uint32]uint32
	hadID            map[uint32]bool
	oldCounters      [32]uint32
	oldOccupied      uint32
	oldGlyphCache    uint32
}

// portTestGeneratorSpawnPolicyFixture is meant to be called from
// portTestGeneratorObjectsEnvironment after resetSpawn is installed. The
// allocator fixture owns the class allocations; this helper owns only its raw
// record slabs and map-normalization entries.
func portTestGeneratorSpawnPolicyFixture(proxy *portTestRoamOwnerServer, resetSpawn func()) (f *portTestGeneratorSpawnPolicy, restore func()) {
	f = &portTestGeneratorSpawnPolicy{
		proxy: proxy, resetSpawn: resetSpawn,
		objID: make(map[uintptr]uint32), updateID: make(map[uintptr]uint32), nodeID: make(map[uintptr]uint32),
		oldIDs: make(map[uint32]uint32), hadID: make(map[uint32]bool),
	}
	f.objectStride = int(unsafe.Sizeof(server.Object{})) + 2*portTestSpawnPolicyGuard
	objs, freeObjects := alloc.Make([]byte{}, portTestSpawnPolicyMaxRecords*f.objectStride)
	updates, freeUpdates := alloc.Make([]byte{}, portTestSpawnPolicyMaxRecords*(portTestSpawnPolicyUpdateSize+2*portTestSpawnPolicyGuard))
	f.objects, f.updates = objs, updates
	for i := range f.oldCounters {
		f.oldCounters[i] = *memmap.PtrUint32(0x5D4594, 2386232+4*uintptr(i))
	}
	f.oldOccupied = *memmap.PtrUint32(0x5D4594, 2386208)
	f.oldGlyphCache = *memmap.PtrUint32(0x5D4594, 2386360)
	f.records = make([]*server.Object, portTestSpawnPolicyMaxRecords)
	for i := range f.records {
		op := unsafe.Pointer(&f.objects[i*f.objectStride+portTestSpawnPolicyGuard])
		up := unsafe.Pointer(&f.updates[i*(portTestSpawnPolicyUpdateSize+2*portTestSpawnPolicyGuard)+portTestSpawnPolicyGuard])
		f.records[i] = (*server.Object)(op)
		f.objID[uintptr(op)] = 10000 + uint32(i)
		f.updateID[uintptr(up)] = 10100 + uint32(i)
	}
	restore = func() {
		f.unindex()
		f.cleanupGlyphs()
		// A caller must not restore while a raw child is still linked: 50D7E0
		// frees allocator nodes, then the C slabs can be released safely.
		f.resetSpawn()
		for i := range f.oldCounters {
			*memmap.PtrUint32(0x5D4594, 2386232+4*uintptr(i)) = f.oldCounters[i]
		}
		*memmap.PtrUint32(0x5D4594, 2386208) = f.oldOccupied
		*memmap.PtrUint32(0x5D4594, 2386360) = f.oldGlyphCache
		for key, had := range f.hadID {
			if had {
				proxy.life.ids[key] = f.oldIDs[key]
			} else {
				delete(proxy.life.ids, key)
			}
		}
		freeUpdates()
		freeObjects()
	}
	return f, restore
}

func (f *portTestGeneratorSpawnPolicy) rememberID(ptr unsafe.Pointer, id uint32) {
	key := uint32(uintptr(ptr))
	if _, done := f.hadID[key]; !done {
		f.oldIDs[key], f.hadID[key] = f.proxy.life.ids[key]
	}
	f.proxy.life.ids[key] = id
}

func (f *portTestGeneratorSpawnPolicy) prepare(actor *server.Object, spec PortTestGeneratorSpawnPolicySpec) {
	f.unindex()
	f.cleanupGlyphs()
	if len(spec.Records) > len(f.records) {
		panic("too many generator spawn-policy records")
	}
	f.resetSpawn()
	for i := uint8(0); i < spec.ReserveMonster; i++ {
		if alloc.AsClass(C.nox_alloc_monsterList_2386220).NewObject() == nil {
			panic("invalid monster-list reserve")
		}
	}
	clear(f.nodeID)
	for i := 0; i < 32; i++ {
		*memmap.PtrUint32(0x5D4594, 2386232+4*uintptr(i)) = 0
	}
	*memmap.PtrUint32(0x5D4594, 2386208) = 0
	*memmap.PtrUint32(0x5D4594, 2386360) = 0
	// Allocator handles may be reached through raw state later; normalize them
	// before any output is captured even though they are no longer result fields.
	f.rememberID(unsafe.Pointer(C.nox_alloc_spawn_2386216), 9001)
	f.rememberID(unsafe.Pointer(C.nox_alloc_monsterList_2386220), 9002)
	for i, rec := range spec.Records {
		oi := i * f.objectStride
		ui := i * (portTestSpawnPolicyUpdateSize + 2*portTestSpawnPolicyGuard)
		ob, ub := f.objects[oi:oi+f.objectStride], f.updates[ui:ui+portTestSpawnPolicyUpdateSize+2*portTestSpawnPolicyGuard]
		clear(ob)
		clear(ub)
		for j := 0; j < portTestSpawnPolicyGuard; j++ {
			ob[j], ob[len(ob)-portTestSpawnPolicyGuard+j] = 0xa5, 0x5a
			ub[j], ub[len(ub)-portTestSpawnPolicyGuard+j] = 0xa5, 0x5a
		}
		u := f.records[i]
		// Preserve the actor's Go extension/server handle so retained C callbacks
		// that return through Go object services remain valid. Only the C ABI
		// prefix is later compared.
		*u = *actor
		u.InvFirstItem, u.InvHolder, u.InvNextItem, u.Field125 = nil, nil, nil, nil
		u.ObjOwner, u.Obj130, u.ObjNext, u.ObjPrev = nil, nil, nil, nil
		u.InitData, u.UseData.Ptr, u.HealthData = nil, nil, nil
		u.ObjIndexBase = server.ObjectIndex{}
		u.ObjIndex = [4]server.ObjectIndex{}
		u.ObjIndexCur = 0
		u.TypeInd = rec.TypeInd
		u.ObjClass = object.Class(rec.Class)
		u.ObjSubClass = object.SubClass(rec.Subclass)
		u.ObjFlags = object.Flags(rec.Flags)
		u.PosVec = types.Pointf{X: *(*float32)(unsafe.Pointer(&rec.Pos[0])), Y: *(*float32)(unsafe.Pointer(&rec.Pos[1]))}
		u.UpdateData = unsafe.Pointer(&ub[portTestSpawnPolicyGuard])
		for off, value := range rec.UpdateWords {
			if off+4 > portTestSpawnPolicyUpdateSize || off&3 != 0 {
				panic("bad generator spawn-policy update offset")
			}
			*(*uint32)(unsafe.Add(u.UpdateData, off)) = value
		}
		f.rememberID(u.CObj(), 10000+uint32(i))
		f.rememberID(u.UpdateData, 10100+uint32(i))
		if rec.Indexed {
			u.NewPos = u.PosVec
			u.Shape.Kind = server.ShapeKindCircle
			u.Shape.Circle.R = 5
			u.Shape.Circle.R2 = 25
			f.rememberID(unsafe.Pointer(&u.ObjIndexBase), uint32(40000+i*8))
			for j := range u.ObjIndex {
				f.rememberID(unsafe.Pointer(&u.ObjIndex[j]), uint32(40001+i*8+j))
			}
			f.proxy.core.Map.AddObjectToIndex(u)
		}

	}
}

// dispatch uses the unchanged original C family. Register assumes the caller
// configured generator update +86 and child update +2192/2196 prerequisites;
// player dimensions/status/list wiring belongs to the shared lifecycle fixture.
//
// Do not test 50DFB0 by treating Object Class/Flags bits as ordinary raw
// masks: its C signature is float* and it converts float-view fields at +8
// and +16 numerically before the byte/flag checks. Direct callback cases must
// construct float values with the intended conversion results.
func (f *portTestGeneratorSpawnPolicy) dispatch(gen *server.Object, spec PortTestGeneratorSpawnPolicySpec, action PortTestGeneratorSpawnPolicyAction) uint32 {
	var u *server.Object
	switch action.Op {
	case PortTestGeneratorSpawnRegister, PortTestGeneratorSpawnRemove, PortTestGeneratorSpawnNonZombieCleanup, PortTestGeneratorSpawnGlyphCleanup, PortTestGeneratorSpawnCandidate:
		if action.Record < 0 || action.Record >= len(spec.Records) {
			panic("invalid generator spawn-policy action record")
		}
		u = f.records[action.Record]
	}
	switch action.Op {
	case PortTestGeneratorSpawnRegister:
		return uint32(spawnPolicyRegister(gen, u))
	case PortTestGeneratorSpawnRemove:
		C.sub_50E140(C.int(uintptr(u.CObj())))
	case PortTestGeneratorSpawnNonZombieCleanup:
		spawnPolicyDeathRelease(u)
	case PortTestGeneratorSpawnGlyphCleanup:
		spawnPolicyGlyphRelease(u)
	case PortTestGeneratorSpawnCandidate:
		spawnPolicyCandidate(u, &f.proxy.life.players[0])
	case PortTestGeneratorSpawnFarCull:
		spawnPolicyFarCull()
	case PortTestGeneratorSpawnVisibleCull:
		return spawnPolicyVisibleCull()
	case PortTestGeneratorSpawnAdmission:
		p := types.Pointf{X: *(*float32)(unsafe.Pointer(&spec.Point[0])), Y: *(*float32)(unsafe.Pointer(&spec.Point[1]))}
		return uint32(spawnPolicyAdmission(gen, p))
	case PortTestGeneratorSpawnTick:
		return spawnPolicyTick()
	default:
		panic("generator spawn-policy operation")
	}
	return 0
}

// run permits zero Records for cull/admission/tick operations. It registers its
// initial records, then returns one fully normalized result after each action.
// Primary's root oracle supplies player/configuration setup through the shared
// optional SpawnPolicy.Players fixture field.
func (f *portTestGeneratorSpawnPolicy) run(gen *server.Object, spec PortTestGeneratorSpawnPolicySpec) []PortTestGeneratorSpawnPolicyResult {
	f.prepare(gen, spec)
	for _, i := range spec.Registers {
		_ = f.dispatch(gen, spec, PortTestGeneratorSpawnPolicyAction{Op: PortTestGeneratorSpawnRegister, Record: i})
	}
	out := make([]PortTestGeneratorSpawnPolicyResult, 0, len(spec.Actions))
	for _, action := range spec.Actions {
		start := len(f.proxy.trace)
		ret := f.dispatch(gen, spec, action)
		snap := f.snapshot(spec, ret)
		snap.Calls = append([]uint32(nil), f.proxy.trace[start:]...)
		snap.GeneratorCount = *(*byte)(unsafe.Add(gen.UpdateData, 86))
		out = append(out, snap)
	}
	return out
}

func (f *portTestGeneratorSpawnPolicy) normalize(p uint32) uint32 {
	if p == 0 {
		return 0
	}
	if v, ok := f.objID[uintptr(p)]; ok {
		return v
	}
	if v, ok := f.updateID[uintptr(p)]; ok {
		return v
	}
	if v, ok := f.proxy.life.ids[p]; ok {
		return v
	}
	return p
}

func policyWords(p unsafe.Pointer, n int, norm func(uint32) uint32) []uint32 {
	out := make([]uint32, 0, n/4)
	for b := unsafe.Slice((*byte)(p), n); len(b) != 0; b = b[4:] {
		out = append(out, norm(binary.LittleEndian.Uint32(b)))
	}
	return out
}

func (f *portTestGeneratorSpawnPolicy) registerListNodes() {
	used := make(map[uint32]bool, len(f.nodeID))
	for _, id := range f.nodeID {
		used[id] = true
	}
	nextSpawn, nextMonster := uint32(20000), uint32(21000)
	allocID := func(next *uint32) uint32 {
		for used[*next] {
			*next++
		}
		id := *next
		used[id] = true
		*next++
		return id
	}
	for p, n := uintptr(C.dword_5d4594_2386212), 0; p != 0; n++ {
		if n == 96 {
			panic("SpawnClass traversal exceeds capacity")
		}
		if _, ok := f.nodeID[p]; !ok {
			id := allocID(&nextSpawn)
			f.nodeID[p] = id
			f.rememberID(unsafe.Pointer(p), id)
		}
		p = uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(p), 4)))
	}
	for p, n := uintptr(C.dword_5d4594_2386224), 0; p != 0; n++ {
		if n == 96 {
			panic("MonsterList traversal exceeds capacity")
		}
		if _, ok := f.nodeID[p]; !ok {
			id := allocID(&nextMonster)
			f.nodeID[p] = id
			f.rememberID(unsafe.Pointer(p), id)
		}
		p = uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(p), 140)))
	}
}

func (f *portTestGeneratorSpawnPolicy) snapshot(spec PortTestGeneratorSpawnPolicySpec, ret uint32) (out PortTestGeneratorSpawnPolicyResult) {
	out.Return = ret
	out.GuardsOK = true
	f.registerListNodes()
	out.SpawnHead = f.normalize(uint32(C.dword_5d4594_2386212))
	out.MonsterListHead = f.normalize(uint32(C.dword_5d4594_2386224))
	out.MonsterListCount = uint32(C.dword_5d4594_2386228)
	out.Occupied = *memmap.PtrUint32(0x5D4594, 2386208)
	out.GlyphCache = *memmap.PtrUint32(0x5D4594, 2386360)
	for i := range out.Counters {
		out.Counters[i] = *memmap.PtrUint32(0x5D4594, 2386232+4*uintptr(i))
	}

	// Bounded list traversals make topology corruption a fixture failure rather
	// than a test hang. Spawn: [object,next,prev]; monster: object/mask,
	// distances[32], next@140, prev@144.
	for p, n := uintptr(C.dword_5d4594_2386212), 0; p != 0; n++ {
		if n == 96 {
			panic("SpawnClass traversal exceeds capacity")
		}
		w := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 3)
		out.Spawn = append(out.Spawn, portTestGeneratorSpawnPolicyNode{f.normalize(w[0]), f.normalize(w[1]), f.normalize(w[2])})
		p = uintptr(w[1])
	}
	for p, n := uintptr(C.dword_5d4594_2386224), 0; p != 0; n++ {
		if n == 96 {
			panic("MonsterList traversal exceeds capacity")
		}
		w := unsafe.Slice((*uint32)(unsafe.Pointer(p)), 37)
		node := portTestGeneratorMonsterListNode{Object: f.normalize(w[0]), Mask: w[2], Next: f.normalize(w[35]), Prev: f.normalize(w[36])}
		copy(node.Distances[:], w[3:35])
		out.MonsterList = append(out.MonsterList, node)
		p = uintptr(w[35])
	}
	for i := range spec.Records {
		j := 0
		for it := f.records[i].InvFirstItem; it != nil; it = it.InvNextItem {
			if j > 3 {
				panic("spawn glyph inventory cycle")
			}
			f.rememberID(it.CObj(), uint32(30000+i*4+j))
			f.rememberID(it.InitData, uint32(31000+i*4+j))
			out.Glyphs = append(out.Glyphs, policyWords(it.CObj(), 772, f.normalize))
			out.GlyphData = append(out.GlyphData, policyWords(it.InitData, 64, f.normalize))
			j++
		}
	}

	for i := range spec.Records {
		oi, ui := i*f.objectStride, i*(portTestSpawnPolicyUpdateSize+2*portTestSpawnPolicyGuard)
		ob, ub := f.objects[oi:oi+f.objectStride], f.updates[ui:ui+portTestSpawnPolicyUpdateSize+2*portTestSpawnPolicyGuard]
		for j := 0; j < portTestSpawnPolicyGuard; j++ {
			out.GuardsOK = out.GuardsOK && ob[j] == 0xa5 && ob[len(ob)-portTestSpawnPolicyGuard+j] == 0x5a && ub[j] == 0xa5 && ub[len(ub)-portTestSpawnPolicyGuard+j] == 0x5a
		}
		out.Objects = append(out.Objects, policyWords(unsafe.Pointer(&ob[portTestSpawnPolicyGuard]), portTestSpawnPolicyObjectSnapshotSize, f.normalize))
		out.Updates = append(out.Updates, policyWords(unsafe.Pointer(&ub[portTestSpawnPolicyGuard]), portTestSpawnPolicyUpdateSize, f.normalize))
	}
	return out
}

// The fixture does not inject allocation failure into alloc.NewClass: the
// retained allocator panics on allocation failure. Fixed-pool exhaustion is
// exercised through the real NewObject service.

func portTestSpawnPolicyEnvironment(proxy *portTestRoamOwnerServer) func() {
	f, free := portTestGeneratorSpawnPolicyFixture(proxy, proxy.callbacks.generation.objects.resetSpawn)
	proxy.callbacks.generation.objects.policy = f
	return free
}

func (f *portTestGeneratorSpawnPolicy) cleanupGlyphs() {
	for _, u := range f.records {
		for it := u.InvFirstItem; it != nil; {
			next := it.InvNextItem
			it.InvFirstItem, it.InvNextItem, it.Field125, it.InvHolder = nil, nil, nil, nil
			for _, p := range []unsafe.Pointer{it.InitData, it.UseData.Ptr, it.UpdateData, unsafe.Pointer(it.HealthData), it.Field189} {
				if p != nil {
					alloc.FreePtr(p)
				}
			}
			it.InitData, it.UseData.Ptr, it.UpdateData, it.HealthData, it.Field189 = nil, nil, nil, nil, nil
			f.proxy.core.Objs.FreeObject(it)
			it = next
		}
		u.InvFirstItem = nil
	}
}

func (f *portTestGeneratorSpawnPolicy) unindex() {
	for _, u := range f.records {
		if !u.ObjFlags.Has(object.FlagPartitioned) {
			continue
		}
		for i := uint32(0); i < u.ObjIndexCur; i++ {
			f.proxy.core.Map.Sub5178E0(true, &u.ObjIndex[i])
		}
		f.proxy.core.Map.Sub5178E0(false, &u.ObjIndexBase)
		u.ObjFlags &^= object.FlagPartitioned
	}
}
