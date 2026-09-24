//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include "GAME4_2.h"

static unsigned short mapRoomCW(void){unsigned short cw;__asm__ volatile("fnstcw %0":"=m"(cw));return cw;}
static void mapRoomSetCW(unsigned short cw){__asm__ volatile("fldcw %0"::"m"(cw));}
*/
import "C"

import (
	"fmt"
	"github.com/opennox/libs/platform"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"math"
	"runtime"
	"sort"
	"unsafe"
)

type PortTestMapRoomArg struct {
	Slot, Offset int
	Value        uint32
}
type PortTestMapRoomRecord struct {
	Size  int
	Words map[int]uint32
	Refs  map[int]PortTestMapRoomArg
}
type PortTestMapRoomWrite struct {
	Slot  int
	Words map[int]uint32
	Refs  map[int]PortTestMapRoomArg
}
type PortTestMapRoomAction struct {
	Op     int
	Args   [4]PortTestMapRoomArg
	Assign int
	Writes []PortTestMapRoomWrite
}
type PortTestMapRoomSpec struct {
	Seed           uint32
	Radius         int
	NoGrid, Buffer bool
	Records        []PortTestMapRoomRecord
	Actions        []PortTestMapRoomAction
}
type PortTestMapRoomRegion struct {
	ID    uint32
	Kind  string
	Alive bool
	Words []uint32
}
type PortTestMapRoomStep struct {
	Return  [2]uint32
	Globals [5]uint32
	Regions []PortTestMapRoomRegion
	Slots   map[int]uint32
}
type PortTestMapRoomResult struct {
	Constants         [6]uint32
	Intact, ControlOK bool
	Steps             []PortTestMapRoomStep
	RandomTail        int
}
type mapRoomTestRegion struct {
	ptr            unsafe.Pointer
	size           int
	id             uint32
	kind           string
	alive, guarded bool
	captureOnly    bool // Diagnostic storage is never an engine pointer target.
}
type mapRoomTestFixture struct {
	regions      []*mapRoomTestRegion
	slots        map[int]*mapRoomTestRegion
	grid, buffer bool
	intact       bool
}

func portTestMapRoomGlobalOwner(i int) *uint32 {
	switch i {
	case 0:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487532))
	case 1:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487536))
	case 2:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487540))
	case 3:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487556))
	case 4:
		return (*uint32)(unsafe.Pointer(&dword_5d4594_2487560))
	default:
		panic("map room global index")
	}
}
func portTestMapRoomGlobalGet(i int) uint32    { return *portTestMapRoomGlobalOwner(i) }
func portTestMapRoomGlobalSet(i int, v uint32) { *portTestMapRoomGlobalOwner(i) = v }

func (f *mapRoomTestFixture) register(p unsafe.Pointer, n int, kind string, guarded bool) *mapRoomTestRegion {
	if p == nil {
		panic("map room fixture allocation failed")
	}
	r := &mapRoomTestRegion{ptr: p, size: n, id: 10000 + uint32(len(f.regions))*10000, kind: kind, alive: true, guarded: guarded}
	f.regions = append(f.regions, r)
	return r
}
func (f *mapRoomTestFixture) resolve(a PortTestMapRoomArg) uint32 {
	if a.Slot == 0 {
		return a.Value
	}
	r := f.slots[a.Slot]
	if r == nil || !r.alive {
		panic(fmt.Sprintf("map room fixture: unavailable slot %d", a.Slot))
	}
	if a.Offset < 0 || a.Offset >= r.size {
		panic("map room fixture offset")
	}
	return uint32(uintptr(unsafe.Add(r.ptr, a.Offset)))
}
func (f *mapRoomTestFixture) normalize(v uint32) uint32 {
	if v == 0 {
		return 0
	}
	for i := len(f.regions) - 1; i >= 0; i-- {
		r := f.regions[i]
		if r.captureOnly {
			continue
		}
		start := uint32(uintptr(r.ptr))
		if v == start || (v >= start && v-start < uint32(r.size)) {
			return r.id + v - start
		}
	}
	return v
}
func (f *mapRoomTestFixture) find(p unsafe.Pointer) *mapRoomTestRegion {
	for i := len(f.regions) - 1; i >= 0; i-- {
		r := f.regions[i]
		if r.ptr == p && r.alive {
			return r
		}
	}
	panic("map room fixture: untracked allocation")
}
func (f *mapRoomTestFixture) write(slot int, words map[int]uint32, refs map[int]PortTestMapRoomArg) {
	r := f.slots[slot]
	if r == nil || !r.alive {
		panic("map room fixture write slot")
	}
	for off, v := range words {
		if off < 0 || off+4 > r.size {
			panic("map room fixture word offset")
		}
		*(*uint32)(unsafe.Add(r.ptr, off)) = v
	}
	for off, a := range refs {
		if off < 0 || off+4 > r.size {
			panic("map room fixture reference offset")
		}
		*(*uint32)(unsafe.Add(r.ptr, off)) = f.resolve(a)
	}
}
func (f *mapRoomTestFixture) guards() {
	for _, r := range f.regions {
		if r.alive && r.guarded {
			for _, v := range unsafe.Slice((*byte)(unsafe.Add(r.ptr, r.size)), 16) {
				if v != 0x5a {
					f.intact = false
				}
			}
		}
	}
}
func (f *mapRoomTestFixture) registerGrid() {
	n := int(portTestMapRoomGlobalGet(2))
	if n < 0 || n > 65 {
		panic("map room fixture grid size")
	}
	p := unsafe.Pointer(uintptr(portTestMapRoomGlobalGet(0)))
	if p == nil {
		return
	}
	f.register(p, n*4, "grid-index", false)
	for _, v := range unsafe.Slice((*uint32)(p), n) {
		f.register(unsafe.Pointer(uintptr(v)), n*20, "grid-row", false)
	}
	f.grid = true
}
func (f *mapRoomTestFixture) freeGridMark() {
	for _, r := range f.regions {
		if r.alive && (r.kind == "grid-index" || r.kind == "grid-row") {
			r.alive = false
		}
	}
	f.grid = false
}
func (f *mapRoomTestFixture) roomFreeMark(p unsafe.Pointer) {
	for node := *(*unsafe.Pointer)(unsafe.Add(p, 368)); node != nil; {
		next := *(*unsafe.Pointer)(unsafe.Add(node, 24))
		f.find(node).alive = false
		node = next
	}
	f.find(p).alive = false
}
func (f *mapRoomTestFixture) before(op int, args [4]uint32) {
	f.guards()
	p := unsafe.Pointer(uintptr(args[0]))
	switch op {
	case 2:
		if f.grid {
			panic("map room fixture: duplicate grid allocation")
		}
	case 3:
		f.freeGridMark()
	case 9:
		if f.buffer {
			panic("map room fixture: duplicate scratch allocation")
		}
	case 10:
		if ptr := unsafe.Pointer(uintptr(portTestMapRoomGlobalGet(3))); ptr != nil {
			f.find(ptr).alive = false
		}
		f.buffer = false
	case 23:
		f.roomFreeMark(p)
	case 24:
		for node := unsafe.Pointer(uintptr(portTestMapRoomGlobalGet(4))); node != nil; {
			next := *(*unsafe.Pointer)(unsafe.Add(node, 56))
			f.roomFreeMark(node)
			node = next
		}
	case 32:
		for node := *(*unsafe.Pointer)(unsafe.Add(p, 368)); node != nil; node = *(*unsafe.Pointer)(unsafe.Add(node, 24)) {
			if *(*uint32)(node) == 1 {
				f.find(node).alive = false
			}
		}
	}
}
func (f *mapRoomTestFixture) after(op int, result uint64, assign int) {
	p := unsafe.Pointer(uintptr(uint32(result)))
	var r *mapRoomTestRegion
	switch op {
	case 2:
		f.registerGrid()
	case 9:
		if result != 0 {
			r = f.register(unsafe.Pointer(uintptr(portTestMapRoomGlobalGet(3))), 8192, "scratch", false)
			f.buffer = true
		}
	case 21, 22, 47, 48:
		if p != nil {
			r = f.register(p, 376, "room", false)
		}
	case 31:
		if p != nil {
			r = f.register(p, 28, "exclusion", false)
		}
	}
	if assign != 0 {
		if r == nil && p != nil {
			r = f.find(p)
		}
		f.slots[assign] = r
	}
	f.guards()
}
func (f *mapRoomTestFixture) snapshot(ret uint64, op int) PortTestMapRoomStep {
	out := PortTestMapRoomStep{Return: [2]uint32{uint32(ret), uint32(ret >> 32)}, Slots: map[int]uint32{}}
	// Floating and int64 results are raw bits; remaining ABIs return 32-bit values.
	if op != 0 && op != 27 && op != 28 && op != 29 && op != 30 && op != 59 {
		out.Return[0] = f.normalize(out.Return[0])
	}
	for i := range out.Globals {
		out.Globals[i] = f.normalize(uint32(portTestMapRoomGlobalGet(i)))
	}
	for _, r := range f.regions {
		item := PortTestMapRoomRegion{ID: r.id, Kind: r.kind, Alive: r.alive}
		if r.alive {
			for _, v := range unsafe.Slice((*uint32)(r.ptr), r.size/4) {
				item.Words = append(item.Words, f.normalize(v))
			}
		}
		out.Regions = append(out.Regions, item)
	}
	ids := make([]int, 0, len(f.slots))
	for i := range f.slots {
		ids = append(ids, i)
	}
	sort.Ints(ids)
	for _, i := range ids {
		if r := f.slots[i]; r != nil {
			out.Slots[i] = r.id
		} else {
			out.Slots[i] = 0
		}
	}
	return out
}
func PortTestMapRooms(cases []PortTestMapRoomSpec) []PortTestMapRoomResult {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	priorCW := C.mapRoomCW()
	defer C.mapRoomSetCW(priorCW)
	cw := (priorCW &^ 0x0f00) | 0x0200
	C.mapRoomSetCW(cw)
	// Actual startup blob data: epsilon 0.1 and opposite directions 1,0,3,2.
	ptrs := mapRoomConstantPointers()
	values := [6]uint32{2576980378, 1069128089, 1, 0, 3, 2}
	var prior [6]uint32
	for i, ptr := range ptrs {
		prior[i] = *ptr
		*ptr = values[i]
	}
	defer func() {
		for i, ptr := range ptrs {
			*ptr = prior[i]
		}
	}()
	priorPlatform := platform.Get()
	platform.Set(platform.New())
	defer platform.Set(priorPlatform)
	var saved [5]uint32
	for i := range saved {
		saved[i] = uint32(portTestMapRoomGlobalGet(i))
	}
	defer func() {
		for i, v := range saved {
			portTestMapRoomGlobalSet(i, v)
		}
	}()
	out := make([]PortTestMapRoomResult, 0, len(cases))
	for _, sp := range cases {
		out = append(out, portTestMapRoomCase(sp, cw))
	}
	return out
}
func portTestMapRoomCase(sp PortTestMapRoomSpec, cw C.ushort) (out PortTestMapRoomResult) {
	f := &mapRoomTestFixture{slots: map[int]*mapRoomTestRegion{}, intact: true}
	defer func() {
		f.guards()
		for _, r := range f.regions {
			if r.alive {
				C.free(r.ptr)
				r.alive = false
			}
		}
	}()
	for i := 0; i < 5; i++ {
		portTestMapRoomGlobalSet(i, 0)
	}
	platform.RandSeed(int64(sp.Seed))
	for i, spec := range sp.Records {
		if spec.Size <= 0 || spec.Size%4 != 0 {
			panic("map room fixture record size")
		}
		n := (spec.Size + 16 + 255) &^ 255
		p := C.aligned_alloc(256, C.size_t(n))
		if p == nil {
			panic("map room input allocation")
		}
		clear(unsafe.Slice((*byte)(p), n))
		for j := 0; j < 16; j++ {
			*(*byte)(unsafe.Add(p, spec.Size+j)) = 0x5a
		}
		f.slots[i+1] = f.register(p, spec.Size, "input", true)
	}
	for i, spec := range sp.Records {
		f.write(i+1, spec.Words, spec.Refs)
	}
	if !sp.NoGrid {
		if sp.Radius < 0 || sp.Radius > 32 {
			panic("map room fixture radius")
		}
		var config [18]uint32
		config[17] = uint32(sp.Radius)
		if mapRoomGridInit(unsafe.Pointer(&config[0])) == 0 {
			panic("map room grid allocation")
		}
		f.registerGrid()
	}
	if sp.Buffer {
		if mapRoomScratchAlloc() == 0 {
			panic("map room scratch allocation")
		}
		f.register(unsafe.Pointer(uintptr(portTestMapRoomGlobalGet(3))), 8192, "scratch", false)
		f.buffer = true
	}
	for _, action := range sp.Actions {
		for _, w := range action.Writes {
			f.write(w.Slot, w.Words, w.Refs)
		}
		var args [4]uint32
		for i, a := range action.Args {
			args[i] = f.resolve(a)
		}
		var ret uint64
		if action.Op >= 0 {
			f.before(action.Op, args)
			ret = mapRoomPortInvoke(action.Op, args)
			f.after(action.Op, ret, action.Assign)
		}
		out.Steps = append(out.Steps, f.snapshot(ret, action.Op))
	}
	f.guards()
	out.Intact = f.intact
	out.ControlOK = C.mapRoomCW() == cw
	out.RandomTail = platform.RandInt()
	for i, ptr := range mapRoomConstantPointers() {
		out.Constants[i] = *ptr
	}
	return out
}

func mapRoomConstantPointers() [6]*uint32 {
	return [6]*uint32{
		memmap.PtrUint32(0x581450, 10432), memmap.PtrUint32(0x581450, 10436),
		memmap.PtrUint32(0x587000, 254952), memmap.PtrUint32(0x587000, 254956),
		memmap.PtrUint32(0x587000, 254960), memmap.PtrUint32(0x587000, 254964),
	}
}

func mapRoomPortInvoke(op int, args [4]uint32) uint64 {
	switch op {
	case 1:
		return uint64(uint32(mapRoomRaw(unsafe.Pointer(mapRoomCell((*[2]int32)(mapRoomPointer(uint32(args[0]))))))))
	case 4:
		return uint64(uint32(mapRoomOccupy((*mapRoom)(mapRoomPointer(uint32(args[0]))))))
	case 5:
		return uint64(uint32(mapRoomVacate((*mapRoom)(mapRoomPointer(uint32(args[0]))))))
	case 17:
		return uint64(uint32(mapRoomRaw(unsafe.Pointer(mapRoomUpdateBounds((*mapRoom)(mapRoomPointer(uint32(args[0]))))))))
	case 36:
		return uint64(uint32(mapRoomPointExcluded((*mapRoom)(mapRoomPointer(uint32(args[0]))), (*types.Pointf)(mapRoomPointer(uint32(args[1]))))))
	case 50:
		return uint64(uint32(mapRoomRaw(unsafe.Pointer(mapRoomSelectDecoration((*mapRoom)(mapRoomPointer(uint32(args[0]))), mapRoomPointer(uint32(args[1])))))))
	case 51:
		return uint64(uint32(mapRoomDecorMatchesFlags(mapRoomPointer(uint32(args[0])), (*mapRoom)(mapRoomPointer(uint32(args[1]))))))
	case 52:
		return uint64(uint32(mapRoomDecorFitsRoom(mapRoomPointer(uint32(args[0])), (*mapRoom)(mapRoomPointer(uint32(args[1]))))))
	case 53:
		return uint64(uint32(mapRoomConsumeDecoration(mapRoomPointer(uint32(args[0])), mapRoomPointer(uint32(args[1])))))
	case 55:
		return uint64(uint32(mapRoomNear(math.Float32frombits(args[0]), math.Float32frombits(args[1]))))
	case 58:
		return uint64(uint32(mapRoomRandomCentered(int32(args[0]), int32(args[1]))))
	case 0:
		return uint64(mapRoomRound((*types.Pointf)(mapRoomPointer(args[0])), (*[2]int32)(mapRoomPointer(args[1]))))
	case 2:
		return uint64(uint32(mapRoomGridInit(mapRoomPointer(args[0]))))
	case 3:
		mapRoomGridFree()
		return 0
	case 6:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomOverlap((*mapRoom)(mapRoomPointer(args[0]))))))
	case 7:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomAt((*[2]int32)(mapRoomPointer(args[0]))))))
	case 8:
		return uint64(uint32(mapRoomResolveOverlap((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1])))))
	case 9:
		return uint64(mapRoomScratchAlloc())
	case 10:
		mapRoomScratchFree()
		return 0
	case 11:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomHead())))
	case 12:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomNext((*mapRoom)(mapRoomPointer(args[0]))))))
	case 13:
		return uint64(uint32(mapRoomAdd((*mapRoom)(mapRoomPointer(args[0])))))
	case 14:
		return uint64(uint32(mapRoomRemove((*mapRoom)(mapRoomPointer(args[0])))))
	case 15:
		return uint64(uint32(mapRoomWithinBounds(mapRoomPointer(args[0]), (*mapRoom)(mapRoomPointer(args[1])))))
	case 16:
		return uint64(uint32(mapRoomCanPlace(mapRoomPointer(args[0]), (*mapRoom)(mapRoomPointer(args[1])))))
	case 18:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomSetPos((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1]))))))
	case 19:
		return uint64(uint32(mapRoomHasRoomNeighbor((*mapRoom)(mapRoomPointer(args[0])), int32(args[1]))))
	case 20:
		return uint64(uint32(mapRoomConnect((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1])), int32(args[2]))))
	case 21:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomNew(int32(args[0]), int32(args[1])))))
	case 22:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomPrepare(mapRoomPointer(args[0])))))
	case 23:
		mapRoomFree((*mapRoom)(mapRoomPointer(args[0])))
		return 0
	case 24:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomFreeAll())))
	case 25:
		return uint64(uint32(mapRoomConnectBoth((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1])), int32(args[2]))))
	case 26:
		return uint64(uint32(mapRoomIsEntranceSide((*mapRoom)(mapRoomPointer(args[0])), int32(args[1]))))
	case 27:
		return math.Float64bits(mapRoomRandomX((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1]))))
	case 28:
		return math.Float64bits(mapRoomRandomY((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1]))))
	case 29:
		return math.Float64bits(mapRoomRandomReverseX((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1]))))
	case 30:
		return math.Float64bits(mapRoomRandomReverseY((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1]))))
	case 31:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomAddExclusion((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1])), math.Float32frombits(args[2]), math.Float32frombits(args[3])))))
	case 32:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomRemoveGeneratedExclusions((*mapRoom)(mapRoomPointer(args[0]))))))
	case 33:
		return uint64(uint32(mapRoomContainsRect((*mapRoom)(mapRoomPointer(args[0])), (*mapRoomExclusion)(mapRoomPointer(args[1])))))
	case 34:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomFindExclusionOverlap((*mapRoom)(mapRoomPointer(args[0])), (*mapRoomExclusion)(mapRoomPointer(args[1]))))))
	case 35:
		return uint64(uint32(mapRoomRandomPoint((*mapRoom)(mapRoomPointer(args[0])), math.Float32frombits(args[1]), (*types.Pointf)(mapRoomPointer(args[2])))))
	case 37:
		return uint64(uint32(int32(mapRoomRememberPoint((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1]))))))
	case 38:
		return uint64(mapRoomIsHall((*mapRoom)(mapRoomPointer(args[0]))))
	case 39:
		return uint64(uint32(mapRoomHallDirection(int32(args[0]))))
	case 40:
		return uint64(uint32(mapRoomOpposite(int32(args[0]))))
	case 41:
		return uint64(uint32(mapRoomHallSide(int32(args[0]))))
	case 42:
		return uint64(uint32(mapRoomHallOppositeSide(int32(args[0]))))
	case 43:
		return uint64(uint32(mapRoomTrimHall((*mapRoom)(mapRoomPointer(args[0])), (*mapRoom)(mapRoomPointer(args[1])))))
	case 44:
		return uint64(uint32(mapRoomHallExit((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1])))))
	case 45:
		return uint64(uint32(mapRoomHallEntry((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1])))))
	case 46:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomHallCenter((*mapRoom)(mapRoomPointer(args[0])), (*types.Pointf)(mapRoomPointer(args[1]))))))
	case 47:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomNewHall(int32(args[0]), int32(args[1]), int32(args[2])))))
	case 48:
		return uint64(mapRoomRaw(unsafe.Pointer(mapRoomPrepareHall(mapRoomPointer(args[0]), int32(args[1]), int32(args[2])))))
	case 49:
		return uint64(uint32(mapRoomRaw(mapRoomAssignDecoration(mapRoomPointer(args[0]), (*mapRoom)(mapRoomPointer(args[1]))))))
	case 54:
		return uint64(mapRoomAssignRequiredDecorations(mapRoomPointer(args[0])))
	case 56:
		mapRoomSeed(args[0])
		return 0
	case 57:
		return uint64(uint32(mapRoomRandomInt(int32(args[0]), int32(args[1]))))
	case 59:
		return math.Float64bits(mapRoomRandomFloat(math.Float32frombits(args[0]), math.Float32frombits(args[1])))
	default:
		C.abort()
		return 0
	}
}
