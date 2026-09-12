//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include "GAME4_2.h"
static float mapRoomFloat(uint32_t v){float f;memcpy(&f,&v,4);return f;}
static uint64_t mapRoomDouble(double v){uint64_t n;memcpy(&n,&v,8);return n;}
long long nox_xxx_mapGenRoundFloatToPtr_520DF0(float2* a1, uint32_t* a2);
int sub_520E60(int2* a1);
int sub_520EA0(int a1);
void sub_520F80();
int sub_521100(int a1);
int sub_521180(int a1);
int sub_521200(int a1);
int sub_521290(int2* a1);
int sub_5212B0(int a1, uint32_t* a2);
int nox_xxx_mapgenAllocBuffer_5213E0();
void nox_xxx_mapgenFreeBuffer_521400();
void* nox_xxx_mapGenGetTopRoom_521710();
int sub_521720(int a1);
int nox_xxx_mapGenAddNewRoom_521730(uint32_t* a1);
int sub_521760(int a1);
int sub_5217A0(int a1, int a2);
int sub_521820(int a1, int a2);
int nox_xxx_mapGenUpdateRoomRect_521850(int a1);
int nox_xxx_mapGenSetRoomPos_521880(uint32_t* a1, float2* a2);
int sub_5218B0(int a1, int a2);
int sub_521900(int a1, int a2, int a3);
float* nox_xxx_mapGenMakeRoomStruct_521940(int a1, int a2);
float* nox_xxx_mapGenPrepareRoom_521990(int a1);
void sub_521A10(void* lpMem);
uint32_t* nox_xxx_mapGenFreeTopRoom_521A40();
int sub_521A70(int a1, int a2, int a3);
int sub_521AA0(uint32_t* a1, int a2);
double sub_521B00(int a1, int a2);
double sub_521B30(int a1, int a2);
double sub_521B60(int a1, int a2);
double sub_521B90(int a1, int a2);
float* sub_521BC0(int a1, float2* a2, float a3, float a4);
uint32_t* sub_521C10(int a1);
int sub_521EB0(float* a1, float* a2);
int sub_521F10(int a1, float* a2);
int sub_5226D0(int a1, float a2, int a3);
int sub_5227B0(int a1, float* a2);
char sub_522CA0(int a1, float* a2);
int nox_xxx_mapGenCheckRoomType_5238F0(int* a1);
int sub_523920(int a1);
int sub_523960(int a1);
int sub_523970(int a1);
int sub_5239B0(int a1);
int sub_523A10(int a1, float* a2);
int sub_523C30(int a1, int a2);
int sub_523CB0(int a1, int a2);
float* sub_523D30(float* a1, float* a2);
float* sub_523E30(int a1, int a2, int a3);
float* nox_xxx_mapGenMakeHall_523EC0(int a1, int a2, int a3);
int sub_524070(int a1, int a2);
int sub_524090(int a1, int* a2);
int nox_xxx_mapGenDecorChkConstaint_5241C0(int a1, int a2);
int nox_xxx_mapGenChkDecorFillsRoom_5241F0(int a1, int a2);
char nox_xxx_mapGenDecorChkLimit_524220(int* a1, int a2);
int nox_xxx_mapGenMakeRooms_524310(int a1);
int sub_524660(float a1, float a2);
void nox_xxx_mapGenSetRngSeed_526AB0(unsigned int a1);
signed int nox_xxx_mapGenRandFunc_526AC0(int a1, signed int a2);
int nox_xxx_mapGenRandFunc2_526B00(int a1, signed int a2);
double sub_526BC0(float a1, float a2);
static uint64_t mapRoomInvoke(int op,const uint32_t* v){switch(op){
case 0: return (uint64_t)(nox_xxx_mapGenRoundFloatToPtr_520DF0((float2*)(uintptr_t)v[0],(uint32_t*)(uintptr_t)v[1]));
case 1: return (uint32_t)(sub_520E60((int2*)(uintptr_t)v[0]));
case 2: return (uint32_t)(sub_520EA0((int)v[0]));
case 3: sub_520F80(); return 0;
case 4: return (uint32_t)(sub_521100((int)v[0]));
case 5: return (uint32_t)(sub_521180((int)v[0]));
case 6: return (uint32_t)(sub_521200((int)v[0]));
case 7: return (uint32_t)(sub_521290((int2*)(uintptr_t)v[0]));
case 8: return (uint32_t)(sub_5212B0((int)v[0],(uint32_t*)(uintptr_t)v[1]));
case 9: return (uint32_t)(nox_xxx_mapgenAllocBuffer_5213E0());
case 10: nox_xxx_mapgenFreeBuffer_521400(); return 0;
case 11: return (uint32_t)(uintptr_t)(nox_xxx_mapGenGetTopRoom_521710());
case 12: return (uint32_t)(sub_521720((int)v[0]));
case 13: return (uint32_t)(nox_xxx_mapGenAddNewRoom_521730((uint32_t*)(uintptr_t)v[0]));
case 14: return (uint32_t)(sub_521760((int)v[0]));
case 15: return (uint32_t)(sub_5217A0((int)v[0],(int)v[1]));
case 16: return (uint32_t)(sub_521820((int)v[0],(int)v[1]));
case 17: return (uint32_t)(nox_xxx_mapGenUpdateRoomRect_521850((int)v[0]));
case 18: return (uint32_t)(nox_xxx_mapGenSetRoomPos_521880((uint32_t*)(uintptr_t)v[0],(float2*)(uintptr_t)v[1]));
case 19: return (uint32_t)(sub_5218B0((int)v[0],(int)v[1]));
case 20: return (uint32_t)(sub_521900((int)v[0],(int)v[1],(int)v[2]));
case 21: return (uint32_t)(uintptr_t)(nox_xxx_mapGenMakeRoomStruct_521940((int)v[0],(int)v[1]));
case 22: return (uint32_t)(uintptr_t)(nox_xxx_mapGenPrepareRoom_521990((int)v[0]));
case 23: sub_521A10((void*)(uintptr_t)v[0]); return 0;
case 24: return (uint32_t)(uintptr_t)(nox_xxx_mapGenFreeTopRoom_521A40());
case 25: return (uint32_t)(sub_521A70((int)v[0],(int)v[1],(int)v[2]));
case 26: return (uint32_t)(sub_521AA0((uint32_t*)(uintptr_t)v[0],(int)v[1]));
case 27: return mapRoomDouble(sub_521B00((int)v[0],(int)v[1]));
case 28: return mapRoomDouble(sub_521B30((int)v[0],(int)v[1]));
case 29: return mapRoomDouble(sub_521B60((int)v[0],(int)v[1]));
case 30: return mapRoomDouble(sub_521B90((int)v[0],(int)v[1]));
case 31: return (uint32_t)(uintptr_t)(sub_521BC0((int)v[0],(float2*)(uintptr_t)v[1],mapRoomFloat(v[2]),mapRoomFloat(v[3])));
case 32: return (uint32_t)(uintptr_t)(sub_521C10((int)v[0]));
case 33: return (uint32_t)(sub_521EB0((float*)(uintptr_t)v[0],(float*)(uintptr_t)v[1]));
case 34: return (uint32_t)(sub_521F10((int)v[0],(float*)(uintptr_t)v[1]));
case 35: return (uint32_t)(sub_5226D0((int)v[0],mapRoomFloat(v[1]),(int)v[2]));
case 36: return (uint32_t)(sub_5227B0((int)v[0],(float*)(uintptr_t)v[1]));
case 37: return (uint32_t)(sub_522CA0((int)v[0],(float*)(uintptr_t)v[1]));
case 38: return (uint32_t)(nox_xxx_mapGenCheckRoomType_5238F0((int*)(uintptr_t)v[0]));
case 39: return (uint32_t)(sub_523920((int)v[0]));
case 40: return (uint32_t)(sub_523960((int)v[0]));
case 41: return (uint32_t)(sub_523970((int)v[0]));
case 42: return (uint32_t)(sub_5239B0((int)v[0]));
case 43: return (uint32_t)(sub_523A10((int)v[0],(float*)(uintptr_t)v[1]));
case 44: return (uint32_t)(sub_523C30((int)v[0],(int)v[1]));
case 45: return (uint32_t)(sub_523CB0((int)v[0],(int)v[1]));
case 46: return (uint32_t)(uintptr_t)(sub_523D30((float*)(uintptr_t)v[0],(float*)(uintptr_t)v[1]));
case 47: return (uint32_t)(uintptr_t)(sub_523E30((int)v[0],(int)v[1],(int)v[2]));
case 48: return (uint32_t)(uintptr_t)(nox_xxx_mapGenMakeHall_523EC0((int)v[0],(int)v[1],(int)v[2]));
case 49: return (uint32_t)(sub_524070((int)v[0],(int)v[1]));
case 50: return (uint32_t)(sub_524090((int)v[0],(int*)(uintptr_t)v[1]));
case 51: return (uint32_t)(nox_xxx_mapGenDecorChkConstaint_5241C0((int)v[0],(int)v[1]));
case 52: return (uint32_t)(nox_xxx_mapGenChkDecorFillsRoom_5241F0((int)v[0],(int)v[1]));
case 53: return (uint32_t)(nox_xxx_mapGenDecorChkLimit_524220((int*)(uintptr_t)v[0],(int)v[1]));
case 54: return (uint32_t)(nox_xxx_mapGenMakeRooms_524310((int)v[0]));
case 55: return (uint32_t)(sub_524660(mapRoomFloat(v[0]),mapRoomFloat(v[1])));
case 56: nox_xxx_mapGenSetRngSeed_526AB0((unsigned int)v[0]); return 0;
case 57: return (uint32_t)(nox_xxx_mapGenRandFunc_526AC0((int)v[0],(signed int)v[1]));
case 58: return (uint32_t)(nox_xxx_mapGenRandFunc2_526B00((int)v[0],(signed int)v[1]));
case 59: return mapRoomDouble(sub_526BC0(mapRoomFloat(v[0]),mapRoomFloat(v[1])));
default:abort();}}

extern uint32_t dword_5d4594_2487532,dword_5d4594_2487536,dword_5d4594_2487540,dword_5d4594_2487556,dword_5d4594_2487560;
static uint32_t* mapRoomGlobals[]={&dword_5d4594_2487532,&dword_5d4594_2487536,&dword_5d4594_2487540,&dword_5d4594_2487556,&dword_5d4594_2487560};
static uint32_t mapRoomGlobalGet(int i){return *mapRoomGlobals[i];}
static void mapRoomGlobalSet(int i,uint32_t v){*mapRoomGlobals[i]=v;}
static unsigned short mapRoomCW(void){unsigned short cw;__asm__ volatile("fnstcw %0":"=m"(cw));return cw;}
static void mapRoomSetCW(unsigned short cw){__asm__ volatile("fldcw %0"::"m"(cw));}
*/
import "C"

import (
	"fmt"
	"github.com/opennox/libs/platform"
	"github.com/opennox/opennox/v1/common/memmap"
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
}
type mapRoomTestFixture struct {
	regions      []*mapRoomTestRegion
	slots        map[int]*mapRoomTestRegion
	grid, buffer bool
	intact       bool
}

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
		start := uint32(uintptr(r.ptr))
		if v >= start && v-start < uint32(r.size) {
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
	n := int(C.mapRoomGlobalGet(2))
	if n < 0 || n > 65 {
		panic("map room fixture grid size")
	}
	p := unsafe.Pointer(uintptr(C.mapRoomGlobalGet(0)))
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
		if ptr := unsafe.Pointer(uintptr(C.mapRoomGlobalGet(3))); ptr != nil {
			f.find(ptr).alive = false
		}
		f.buffer = false
	case 23:
		f.roomFreeMark(p)
	case 24:
		for node := unsafe.Pointer(uintptr(C.mapRoomGlobalGet(4))); node != nil; {
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
			r = f.register(unsafe.Pointer(uintptr(C.mapRoomGlobalGet(3))), 8192, "scratch", false)
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
		out.Globals[i] = f.normalize(uint32(C.mapRoomGlobalGet(C.int(i))))
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
		saved[i] = uint32(C.mapRoomGlobalGet(C.int(i)))
	}
	defer func() {
		for i, v := range saved {
			C.mapRoomGlobalSet(C.int(i), C.uint32_t(v))
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
		C.mapRoomGlobalSet(C.int(i), 0)
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
		if C.sub_520EA0(C.int(uintptr(unsafe.Pointer(&config[0])))) == 0 {
			panic("map room grid allocation")
		}
		f.registerGrid()
	}
	if sp.Buffer {
		if C.nox_xxx_mapgenAllocBuffer_5213E0() == 0 {
			panic("map room scratch allocation")
		}
		f.register(unsafe.Pointer(uintptr(C.mapRoomGlobalGet(3))), 8192, "scratch", false)
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
			ret = uint64(C.mapRoomInvoke(C.int(action.Op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
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
