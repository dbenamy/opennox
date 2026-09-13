//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2487248;
extern uint32_t dword_5d4594_3835348;
extern uint32_t dword_5d4594_3835352;
extern uint32_t dword_5d4594_3835356;
extern uint32_t dword_5d4594_3835360;
extern uint32_t dword_5d4594_3835364;
extern uint32_t dword_5d4594_3835368;
extern uint32_t dword_5d4594_3835372;
extern uint32_t dword_5d4594_3835388;
extern uint32_t dword_5d4594_3835392;
extern uint32_t dword_5d4594_588084;
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_2489436;
extern void* dword_5d4594_251560;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern nox_tileDef_t nox_tile_defs_arr[176];
extern uint32_t nox_tile_def_cnt;
int* nox_xxx_tileListAddNewSubtile_422160(int a1, int a2, int a3, int a4);
int nox_xxx_tileFreeTileOne_4221E0(void* a1);
int nox_xxx_tileFreeTile_422200(int a1);
unsigned char nox_xxx_wall_42A6C0(unsigned char a1, unsigned char a2);
int nox_xxx_mapGenFixCoords_4D3D90(float2* a1, float2* a2);
int sub_4D3FF0(int a1);
float* sub_51D5E0(float* a1);
int sub_51D8F0(float2* a1);
int sub_51D9C0(int a1, int a2, int a3, int a4, int a5);
int sub_51DA70(int a1, int a2, int a3, int a4, int a5);
int sub_5244D0(int a1);
int sub_524500(float2* a1, int a2);
int sub_524550(int* a1, int a2);
float* sub_5245A0(int a1, float* a2, int a3, int a4);
float* sub_524610(int a1, float* a2, int a3);
int nox_xxx_gen_524680(int a1, int a2, int a3);
int sub_524950(int a1, int a2, float* a3, int* a4);
int sub_5249C0(int a1, int a2, float* a3, int* a4);
void sub_524B50(int a1, int a2, float* a3, int* a4);
void nox_xxx_gen_524E00(int a1, int a2);
int sub_524FB0(int a1, int a2, int a3);
float2* sub_525330(float2* a1, int a2);
float2* sub_525370(float2* a1, int a2);
int sub_5253B0(float* a1);
void nox_xxx_mapgen_525510(int a1, int a2);
void nox_xxx_mapgen_525570(int a1, int a2, int a3, int a4);
int nox_xxx_mapgen_525690(int a1, float2* a2, int a3);
int nox_xxx_mapgen_525740(int a1, float2* a2, int a3);
int nox_xxx_mapgen_525830(int a1, float2* a2, int a3);
int nox_xxx_mapgen_5258E0(int a1, float2* a2, int a3);
int sub_526C40(int a1);
int sub_526C80(int a1);
int sub_526D50(int a1);
int sub_526DD0(float2* a1, int* a2);
int sub_526E60(float* a1);
int sub_527030(float2* a1);
int sub_527380(float* a1);
int sub_527450(uint32_t* a1);
int nox_xxx_mapGenGetObjID_527940(char* a1);
float* nox_xxx_mapGenPlaceObj_5279B0(float2* a1);
float* nox_xxx_mapGenMoveObject_527A10(float* a1, float2* a2);
int nox_xxx_mapGenOrientObj_527C60(int a1, int a2);
int nox_xxx_mapGenFinishSpellbook_527DB0(int a1, char a2);
int sub_543680(float* a1);
int sub_5437E0(int* a1, int a2, int a3);
void sub_543BC0(int a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_tile_543C50(uint32_t* a1, int a2, int a3, int a4, int a5, int a6);
int nox_xxx_tileSubtile_544310(float2* a1);
static uint32_t paintInvoke(int op,const uint32_t* v){switch(op){
case 0: return (uint32_t)(uintptr_t)nox_xxx_tileListAddNewSubtile_422160((int)v[0],(int)v[1],(int)v[2],(int)v[3]);
case 1: return (uint32_t)nox_xxx_tileFreeTileOne_4221E0((void*)(uintptr_t)v[0]);
case 2: return (uint32_t)nox_xxx_tileFreeTile_422200((int)v[0]);
case 3: return (uint32_t)nox_xxx_wall_42A6C0((unsigned char)v[0],(unsigned char)v[1]);
case 4: return (uint32_t)nox_xxx_mapGenFixCoords_4D3D90((float2*)(uintptr_t)v[0],(float2*)(uintptr_t)v[1]);
case 5: return (uint32_t)sub_4D3FF0((int)v[0]);
case 6: return (uint32_t)(uintptr_t)sub_51D5E0((float*)(uintptr_t)v[0]);
case 7: return (uint32_t)sub_51D8F0((float2*)(uintptr_t)v[0]);
case 8: return (uint32_t)sub_51D9C0((int)v[0],(int)v[1],(int)v[2],(int)v[3],(int)v[4]);
case 9: return (uint32_t)sub_51DA70((int)v[0],(int)v[1],(int)v[2],(int)v[3],(int)v[4]);
case 10: return (uint32_t)sub_5244D0((int)v[0]);
case 11: return (uint32_t)sub_524500((float2*)(uintptr_t)v[0],(int)v[1]);
case 12: return (uint32_t)sub_524550((int*)(uintptr_t)v[0],(int)v[1]);
case 13: return (uint32_t)(uintptr_t)sub_5245A0((int)v[0],(float*)(uintptr_t)v[1],(int)v[2],(int)v[3]);
case 14: return (uint32_t)(uintptr_t)sub_524610((int)v[0],(float*)(uintptr_t)v[1],(int)v[2]);
case 15: return (uint32_t)nox_xxx_gen_524680((int)v[0],(int)v[1],(int)v[2]);
case 16: return (uint32_t)sub_524950((int)v[0],(int)v[1],(float*)(uintptr_t)v[2],(int*)(uintptr_t)v[3]);
case 17: return (uint32_t)sub_5249C0((int)v[0],(int)v[1],(float*)(uintptr_t)v[2],(int*)(uintptr_t)v[3]);
case 18: sub_524B50((int)v[0],(int)v[1],(float*)(uintptr_t)v[2],(int*)(uintptr_t)v[3]);return 0;
case 19: nox_xxx_gen_524E00((int)v[0],(int)v[1]);return 0;
case 20: return (uint32_t)sub_524FB0((int)v[0],(int)v[1],(int)v[2]);
case 21: return (uint32_t)(uintptr_t)sub_525330((float2*)(uintptr_t)v[0],(int)v[1]);
case 22: return (uint32_t)(uintptr_t)sub_525370((float2*)(uintptr_t)v[0],(int)v[1]);
case 23: return (uint32_t)sub_5253B0((float*)(uintptr_t)v[0]);
case 24: nox_xxx_mapgen_525510((int)v[0],(int)v[1]);return 0;
case 25: nox_xxx_mapgen_525570((int)v[0],(int)v[1],(int)v[2],(int)v[3]);return 0;
case 26: return (uint32_t)nox_xxx_mapgen_525690((int)v[0],(float2*)(uintptr_t)v[1],(int)v[2]);
case 27: return (uint32_t)nox_xxx_mapgen_525740((int)v[0],(float2*)(uintptr_t)v[1],(int)v[2]);
case 28: return (uint32_t)nox_xxx_mapgen_525830((int)v[0],(float2*)(uintptr_t)v[1],(int)v[2]);
case 29: return (uint32_t)nox_xxx_mapgen_5258E0((int)v[0],(float2*)(uintptr_t)v[1],(int)v[2]);
case 30: return (uint32_t)sub_526C40((int)v[0]);
case 31: return (uint32_t)sub_526C80((int)v[0]);
case 32: return (uint32_t)sub_526D50((int)v[0]);
case 33: return (uint32_t)sub_526DD0((float2*)(uintptr_t)v[0],(int*)(uintptr_t)v[1]);
case 34: return (uint32_t)sub_526E60((float*)(uintptr_t)v[0]);
case 35: return (uint32_t)sub_527030((float2*)(uintptr_t)v[0]);
case 36: return (uint32_t)sub_527380((float*)(uintptr_t)v[0]);
case 37: return (uint32_t)sub_527450((uint32_t*)(uintptr_t)v[0]);
case 38: return (uint32_t)nox_xxx_mapGenGetObjID_527940((char*)(uintptr_t)v[0]);
case 39: return (uint32_t)(uintptr_t)nox_xxx_mapGenPlaceObj_5279B0((float2*)(uintptr_t)v[0]);
case 40: return (uint32_t)(uintptr_t)nox_xxx_mapGenMoveObject_527A10((float*)(uintptr_t)v[0],(float2*)(uintptr_t)v[1]);
case 41: return (uint32_t)nox_xxx_mapGenOrientObj_527C60((int)v[0],(int)v[1]);
case 42: return (uint32_t)nox_xxx_mapGenFinishSpellbook_527DB0((int)v[0],(char)v[1]);
case 43: return (uint32_t)sub_543680((float*)(uintptr_t)v[0]);
case 44: return (uint32_t)sub_5437E0((int*)(uintptr_t)v[0],(int)v[1],(int)v[2]);
case 45: sub_543BC0((int)v[0],(int)v[1],(int)v[2],(int)v[3],(int)v[4],(int)v[5]);return 0;
case 46: return (uint32_t)nox_xxx_tile_543C50((uint32_t*)(uintptr_t)v[0],(int)v[1],(int)v[2],(int)v[3],(int)v[4],(int)v[5]);
case 47: return (uint32_t)nox_xxx_tileSubtile_544310((float2*)(uintptr_t)v[0]);
default:abort();}}
static uint32_t* paintGlobal(int i){switch(i){
case 0:return &dword_5d4594_2487248;
case 1:return &dword_5d4594_3835348;
case 2:return &dword_5d4594_3835352;
case 3:return &dword_5d4594_3835356;
case 4:return &dword_5d4594_3835360;
case 5:return &dword_5d4594_3835364;
case 6:return &dword_5d4594_3835368;
case 7:return &dword_5d4594_3835372;
case 8:return &dword_5d4594_3835388;
case 9:return &dword_5d4594_3835392;
case 10:return &dword_5d4594_588084;
case 11:return &dword_5d4594_251572;
case 12:return &dword_5d4594_2489436;
case 13:return (uint32_t*)&ptr_5D4594_2650668;
case 14:return (uint32_t*)&dword_5d4594_251560;
case 15:return &nox_tile_def_cnt;
default:abort();}}
static void* paintXfer(int i){return i ? (void*)nox_xxx_XFerSpellReward_4F5F30:(void*)nox_xxx_XFerDoor_4F4CB0;}
static unsigned short paintCW(){unsigned short cw;__asm__ __volatile__("fnstcw %0":"=m"(cw));return cw;}
static void paintSetCW(unsigned short cw){__asm__ __volatile__("fldcw %0"::"m"(cw));}
*/
import "C"
import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"runtime"
	"sort"
	"unsafe"

	"github.com/opennox/libs/platform"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
)

type PortTestPaintAction struct {
	Globals map[string]PortTestMapRoomArg
	Op      int
	Args    [6]PortTestMapRoomArg
	Assign  int
	Writes  []PortTestMapRoomWrite
}
type PortTestPaintWall struct {
	X, Y  int
	Words map[int]uint32
	Data  PortTestMapRoomArg
}
type PortTestPaintObject struct {
	Slot, Type       int
	Words, DataWords map[int]uint32
	Refs             map[int]PortTestMapRoomArg
}
type PortTestPaintCell struct {
	X, Y  int
	Words map[int]uint32
	Refs  map[int]PortTestMapRoomArg
}
type PortTestPaintSpec struct {
	ExplicitWallLimits bool
	Seed               int
	Records            []PortTestMapRoomRecord
	Objects            []PortTestPaintObject
	Walls              []PortTestPaintWall
	Cells              []PortTestPaintCell
	Globals            map[string]PortTestMapRoomArg
	Actions            []PortTestPaintAction
	Cycle, Variations  byte
}
type PortTestPaintCellState struct {
	X, Y  int
	Words [11]uint32
}
type PortTestPaintStep struct {
	Rows    [128]uint32
	Return  uint32
	Globals map[string]uint32
	Records []PortTestMapRoomRegion
	Slots   map[int]uint32
	Cells   []PortTestPaintCellState
	Walls   server.PortTestPaintWallState
	Objects [7]uint32
	Tables  map[string][32]byte
}
type PortTestPaintResult struct {
	Intact, ControlOK     bool
	Steps                 []PortTestPaintStep
	RandomTail, LogicTail int
}
type paintTestFixture struct {
	table unsafe.Pointer
	*mapRoomTestFixture
	owners        *server.PortTestPaintOwners
	rows          []unsafe.Pointer
	owned         map[*mapRoomTestRegion]bool
	objectRecords map[*mapRoomTestRegion]*server.Object
	globs         map[string]*uint32
	secret        map[*mapRoomTestRegion]bool
}

func paintGlobals() map[string]*uint32 {
	names := []string{"dword_5d4594_2487248", "dword_5d4594_3835348", "dword_5d4594_3835352", "dword_5d4594_3835356", "dword_5d4594_3835360", "dword_5d4594_3835364", "dword_5d4594_3835368", "dword_5d4594_3835372", "dword_5d4594_3835388", "dword_5d4594_3835392", "dword_5d4594_588084", "dword_5d4594_251572", "dword_5d4594_2489436", "grid", "secretWalls", "tileCount"}
	out := map[string]*uint32{}
	for i, n := range names {
		out[n] = (*uint32)(unsafe.Pointer(C.paintGlobal(C.int(i))))
	}
	for n, off := range map[string]uintptr{"tile": 35912, "tileFlag": 35916, "wall": 35948, "wallDir": 35952, "wallVariation": 35956, "worklistError": 22200} {
		out[n] = memmap.PtrUint32(0x973F18, off)
	}
	return out
}
func (f *paintTestFixture) allocate(size int, kind string) *mapRoomTestRegion {
	n := (size + 16 + 255) &^ 255
	p := C.aligned_alloc(256, C.size_t(n))
	if p == nil {
		panic("painting input allocation")
	}
	clear(unsafe.Slice((*byte)(p), n))
	for i := 0; i < 16; i++ {
		*(*byte)(unsafe.Add(p, size+i)) = 0x5a
	}
	r := f.register(p, size, kind, true)
	f.owned[r] = true
	return r
}
func (f *paintTestFixture) known(p unsafe.Pointer) *mapRoomTestRegion {
	for _, r := range f.regions {
		if r.alive && r.ptr == p {
			return r
		}
	}
	return nil
}
func (f *paintTestFixture) containing(p unsafe.Pointer) *mapRoomTestRegion {
	v := uintptr(p)
	for _, r := range f.regions {
		if r.alive && v >= uintptr(r.ptr) && v-uintptr(r.ptr) < uintptr(r.size) {
			return r
		}
	}
	return nil
}
func (f *paintTestFixture) norm(v uint32) uint32 {
	if v < 4096 {
		return v
	}
	for i := 0; i < 2; i++ {
		if v == mapRoomRaw(C.paintXfer(C.int(i))) {
			return 0x30000000 + uint32(i)
		}
	}
	if n := f.normalize(v); n != v {
		return n
	}
	return f.owners.NormalizeWall(v)
}
func (f *paintTestFixture) cell(x, y int) *[11]uint32 {
	return (*[11]uint32)(unsafe.Add(f.rows[x], y*44))
}
func (f *paintTestFixture) initGrid() {
	table := f.allocate(128*4, "floorRows")
	f.table = table.ptr
	f.rows = make([]unsafe.Pointer, 128)
	for x := range f.rows {
		r := f.allocate(128*44, "floor")
		f.rows[x] = r.ptr
		*mapRoomRef(table.ptr, x*4) = r.ptr
		for y := 0; y < 128; y++ {
			cell := f.cell(x, y)
			cell[1] = 255
			cell[6] = 255
		}
	}
	*f.globs["grid"] = mapRoomRaw(table.ptr)
}
func (f *paintTestFixture) trackObject(u *server.Object) *mapRoomTestRegion {
	if u == nil {
		return nil
	}
	if old := f.known(u.CObj()); old != nil {
		return old
	}
	f.owners.TrackObjects(u)
	r := f.register(u.CObj(), int(unsafe.Sizeof(server.Object{})), "object", false)
	f.objectRecords[r] = u
	if u.UseData.Ptr != nil {
		f.register(u.UseData.Ptr, int(f.owners.S.Types.ByInd(int(u.TypeInd)).UseDataSize), "objectUseData", false)
	}
	if u.UpdateData != nil {
		size := int(f.owners.S.Types.ByInd(int(u.TypeInd)).UpdateDataSize)
		f.register(u.UpdateData, size, "objectData", false)
	}
	return r
}
func (f *paintTestFixture) discover(ret uint32, op int) {
	// Subtiles remain reachable through the free list or floor/input chains. Every
	// new calloc block has ten adjacent 20-byte nodes, all of which stay reachable.
	var order []uintptr
	seen := map[uintptr]bool{}
	visit := func(head uint32) {}
	visit = func(head uint32) {
		for p := mapRoomPointer(head); p != nil; p = *mapRoomRef(p, 16) {
			addr := uintptr(p)
			if seen[addr] {
				return
			}
			seen[addr] = true
			if f.containing(p) == nil {
				order = append(order, addr)
			}
		}
	}
	visit(*f.globs["dword_5d4594_588084"])
	for x := 0; x < 128; x++ {
		for y := 0; y < 128; y++ {
			c := f.cell(x, y)
			visit(c[5])
			visit(c[10])
		}
	}
	if op == 0 {
		visit(ret)
	}
	// Direct subtile-chain calls use small input records as chain heads.
	for _, r := range f.regions {
		if r.alive && r.kind == "input" && r.size == 20 {
			visit(*(*uint32)(unsafe.Add(r.ptr, 16)))
		}
	}
	sorted := append([]uintptr(nil), order...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	bases := map[uintptr]uintptr{}
	for i := 0; i < len(sorted); {
		j := i + 1
		for j < len(sorted) && sorted[j] == sorted[j-1]+20 {
			j++
		}
		if (j-i)%10 != 0 {
			panic(fmt.Sprintf("incomplete subtile allocation: %d nodes", j-i))
		}
		for k := i; k < j; k++ {
			bases[sorted[k]] = sorted[i+((k-i)/10)*10]
		}
		i = j
	}
	registered := map[uintptr]bool{}
	for _, a := range order {
		base := bases[a]
		if !registered[base] {
			r := f.register(unsafe.Pointer(base), 200, "subtilePool", false)
			f.owned[r] = true
			registered[base] = true
		}
	}
	for _, r := range append([]*mapRoomTestRegion(nil), f.regions...) {
		if r.alive && r.kind == "input" && r.size == 376 {
			for p := *mapRoomRef(r.ptr, 368); p != nil; p = *mapRoomRef(p, 24) {
				if f.known(p) == nil {
					child := f.register(p, 28, "exclusion", false)
					f.owned[child] = true
				}
			}
		}
	}
	for u := f.owners.S.Objs.List; u != nil; u = u.ObjNext {
		f.trackObject(u)
	}
	if (op == 39 || op == 40) && ret != 0 {
		f.trackObject((*server.Object)(mapRoomPointer(ret)))
	}
	remaining := map[unsafe.Pointer]bool{}
	for p := mapRoomPointer(*f.globs["secretWalls"]); p != nil; p = *mapRoomRef(p, 0) {
		remaining[p] = true
	}
	for r := range f.secret {
		if !remaining[r.ptr] {
			r.alive = false
			delete(f.owned, r)
		}
	}
}
func (f *paintTestFixture) snapshot(ret uint32) (out PortTestPaintStep) {
	for i, v := range *(*[128]uint32)(f.table) {
		out.Rows[i] = f.norm(v)
		if mapRoomPointer(v) != f.rows[i] {
			f.intact = false
		}
	}
	out.Return = f.norm(ret)
	out.Globals = map[string]uint32{}
	for n, p := range f.globs {
		out.Globals[n] = f.norm(*p)
	}
	for _, r := range f.regions {
		if r.kind == "floor" || r.kind == "floorRows" {
			continue
		}
		v := PortTestMapRoomRegion{ID: r.id, Kind: r.kind, Alive: r.alive}
		if r.alive {
			for _, word := range unsafe.Slice((*uint32)(r.ptr), r.size/4) {
				v.Words = append(v.Words, f.norm(word))
			}
			if u := f.objectRecords[r]; u != nil {
				f.intact = f.owners.NormalizeObject(u, v.Words) && f.intact
			}
		}
		out.Records = append(out.Records, v)
	}
	out.Slots = map[int]uint32{}
	for i, r := range f.slots {
		if r != nil {
			out.Slots[i] = f.norm(mapRoomRaw(r.ptr))
		} else {
			out.Slots[i] = 0
		}
	}
	for x := 0; x < 128; x++ {
		for y := 0; y < 128; y++ {
			c := f.cell(x, y)
			empty := true
			for i, v := range c {
				want := uint32(0)
				if i == 1 || i == 6 {
					want = 255
				}
				if v != want {
					empty = false
					break
				}
			}
			if empty {
				continue
			}
			v := PortTestPaintCellState{X: x, Y: y}
			for i, word := range c {
				v.Words[i] = f.norm(word)
			}
			out.Cells = append(out.Cells, v)
		}
	}
	out.Walls = f.owners.WallSnapshot(f.norm)
	f.intact = f.intact && out.Walls.Intact
	out.Objects = f.owners.ObjectState()
	for i, v := range out.Objects {
		if i >= 5 {
			out.Objects[i] = f.norm(v)
		}
	}
	out.Tables = map[string][32]byte{"tiles": sha256.Sum256(unsafe.Slice((*byte)(unsafe.Pointer(&C.nox_tile_defs_arr[0])), 176*int(unsafe.Sizeof(server.TileDef{})))), "borders": sha256.Sum256(portTestEdgeTable()), "edgeMap": sha256.Sum256(unsafe.Slice(memmap.PtrUint8(0x587000, 282736), 576))}
	out.Tables["wallDirections"] = sha256.Sum256(unsafe.Slice(memmap.PtrUint8(0x587000, 71276), 200))
	out.Tables["monsterAngles"] = sha256.Sum256(unsafe.Slice(memmap.PtrUint8(0x587000, 230052), 40))
	out.Tables["wallMasks"] = sha256.Sum256(unsafe.Slice(memmap.PtrUint8(0x587000, 255052), 64))
	return out
}
func PortTestMapPainting(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	cw := C.paintCW()
	defer C.paintSetCW(cw)
	wantCW := (cw &^ 0x0f00) | 0x0200
	C.paintSetCW(wantCW)
	oldPlatform := platform.Get()
	platform.Set(platform.New())
	defer platform.Set(oldPlatform)
	oldFlags := noxflags.GetGame()
	noxflags.ResetGame()
	defer func() { noxflags.ResetGame(); noxflags.SetGame(oldFlags) }()
	globs := paintGlobals()
	saved := map[string]uint32{}
	for n, p := range globs {
		saved[n] = *p
	}
	defer func() {
		for n, p := range globs {
			*p = saved[n]
		}
	}()
	var restore []func()
	saveBytes := func(p unsafe.Pointer, n int) {
		old := bytes.Clone(unsafe.Slice((*byte)(p), n))
		restore = append(restore, func() { copy(unsafe.Slice((*byte)(p), n), old) })
	}
	saveBytes(unsafe.Pointer(&C.nox_tile_defs_arr[0]), 176*int(unsafe.Sizeof(server.TileDef{})))
	saveBytes(unsafe.Pointer(&portTestEdgeTable()[0]), len(portTestEdgeTable()))
	saveBytes(memmap.PtrOff(0x973F18, 16200), 6000)
	for off, data := range blobdata.PortTestMapPaintingTables() {
		ptr := memmap.PtrOff(0x587000, off)
		saveBytes(ptr, len(data))
		copy(unsafe.Slice((*byte)(ptr), len(data)), data)
	}
	defer func() {
		for i := len(restore) - 1; i >= 0; i-- {
			restore[i]()
		}
	}()
	core := new(server.Server)
	owners := core.PortTestMapPaintingOwners([2]unsafe.Pointer{C.paintXfer(0), C.paintXfer(1)})
	defer owners.Close()
	realOwner, freeOwner := owner(core)
	defer freeOwner()
	oldGet := GetServer
	GetServer = func() Server { return realOwner }
	defer func() { GetServer = oldGet }()
	var out []PortTestPaintResult
	for _, sp := range cases {
		cycle, variations := sp.Cycle, sp.Variations
		if cycle == 0 && !sp.ExplicitWallLimits {
			cycle = 3
		}
		if variations == 0 && !sp.ExplicitWallLimits {
			variations = 4
		}
		owners.Reset(sp.Seed, cycle, variations)
		out = append(out, paintTestCase(sp, owners, globs, wantCW))
	}
	return out
}
func paintTestCase(sp PortTestPaintSpec, owners *server.PortTestPaintOwners, globs map[string]*uint32, cw C.ushort) (out PortTestPaintResult) {
	f := &paintTestFixture{mapRoomTestFixture: &mapRoomTestFixture{slots: map[int]*mapRoomTestRegion{}, intact: true}, owners: owners, owned: map[*mapRoomTestRegion]bool{}, objectRecords: map[*mapRoomTestRegion]*server.Object{}, globs: globs, secret: map[*mapRoomTestRegion]bool{}}
	defer func() {
		for r := range f.owned {
			if r.alive {
				C.free(r.ptr)
				r.alive = false
			}
		}
	}()
	for _, p := range globs {
		*p = 0
	}
	*globs["tileCount"] = 176
	*globs["dword_5d4594_251572"] = 2
	*globs["dword_5d4594_3835356"] = 255
	*globs["dword_5d4594_3835392"] = 500
	tiles := unsafe.Slice((*server.TileDef)(unsafe.Pointer(&C.nox_tile_defs_arr[0])), 176)
	clear(tiles)
	for i, name := range []string{"PaintTile", "PaintTileTwo"} {
		copy(tiles[i].NameBuf[:], name)
		raw := unsafe.Slice((*byte)(unsafe.Pointer(&tiles[i])), int(unsafe.Sizeof(server.TileDef{})))
		raw[52] = 3
		raw[53] = 3
	}
	edges := portTestEdgeTable()
	clear(edges)
	for i, name := range []string{"PaintBorder", "PaintBorderTwo"} {
		row := edges[i*60 : (i+1)*60]
		copy(row, name)
		row[44] = 8
		row[52] = 3
		row[53] = 3
	}
	clear(unsafe.Slice(memmap.PtrUint32(0x973F18, 16200), 1500))
	f.register(memmap.PtrOff(0x973F18, 16200), 6000, "worklist", false)
	platform.RandSeed(int64(sp.Seed))
	f.initGrid()
	for i, spec := range sp.Records {
		f.slots[i+1] = f.allocate(spec.Size, "input")
	}
	for _, spec := range sp.Objects {
		u := owners.S.Objs.NewObject(owners.S.Types.ByInd(spec.Type))
		f.slots[spec.Slot] = f.trackObject(u)
		if u.UpdateData != nil {
			f.slots[spec.Slot+1000] = f.known(u.UpdateData)
		}
	}
	for i, spec := range sp.Records {
		f.write(i+1, spec.Words, spec.Refs)
	}
	for _, spec := range sp.Objects {
		f.write(spec.Slot, spec.Words, spec.Refs)
		if len(spec.DataWords) != 0 {
			f.write(spec.Slot+1000, spec.DataWords, nil)
		}
	}
	for _, spec := range sp.Cells {
		c := f.cell(spec.X, spec.Y)
		for off, v := range spec.Words {
			c[off/4] = v
		}
		for off, a := range spec.Refs {
			c[off/4] = f.resolve(a)
		}
	}
	for n, a := range sp.Globals {
		ptr := globs[n]
		if ptr == nil {
			panic("unknown painting global " + n)
		}
		*ptr = f.resolve(a)
	}
	for _, spec := range sp.Walls {
		w := owners.CreateWall(spec.X, spec.Y)
		if w == nil {
			panic("initial paint wall")
		}
		for off, v := range spec.Words {
			*(*uint32)(unsafe.Add(w.C(), off)) = v
		}
		if a := spec.Data; a.Slot != 0 || a.Value != 0 {
			w.Data = mapRoomPointer(f.resolve(a))
		}
	}
	for p := mapRoomPointer(*globs["secretWalls"]); p != nil; p = *mapRoomRef(p, 0) {
		r := f.known(p)
		if r == nil {
			panic("untracked secret wall")
		}
		f.secret[r] = true
	}
	for _, a := range sp.Actions {
		for n, arg := range a.Globals {
			p := globs[n]
			if p == nil {
				panic("unknown painting action global " + n)
			}
			*p = f.resolve(arg)
		}
		for _, w := range a.Writes {
			f.write(w.Slot, w.Words, w.Refs)
		}
		var args [6]uint32
		for i, v := range a.Args {
			args[i] = f.resolve(v)
		}
		f.guards()
		var ret uint32
		if a.Op >= 0 {
			ret = uint32(C.paintInvoke(C.int(a.Op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
			f.discover(ret, a.Op)
		}
		if a.Assign != 0 {
			if ret == 0 {
				f.slots[a.Assign] = nil
			} else {
				p := mapRoomPointer(ret)
				r := f.known(p)
				if r == nil {
					base := f.containing(p)
					if base == nil {
						panic("unknown painting result allocation")
					}
					r = &mapRoomTestRegion{ptr: p, size: base.size - int(uintptr(p)-uintptr(base.ptr)), id: f.norm(ret), kind: "alias", alive: true}
				}
				f.slots[a.Assign] = r
			}
		}
		f.guards()
		out.Steps = append(out.Steps, f.snapshot(ret))
	}
	out.Intact = f.intact
	out.ControlOK = C.paintCW() == cw
	out.RandomTail = platform.RandInt()
	out.LogicTail = owners.S.Rand.Logic.IntClamp(0, 10000)
	return out
}
