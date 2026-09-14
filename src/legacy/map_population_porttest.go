//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include "GAME3_3.h"
int sub_4D42C0();
int nox_xxx_mapGenSpellIdByName_51E1D0(const char* a1);
int nox_xxx_mapgen_521C60(int a1, int a2);
int sub_521CB0(int a1, int a2, int a3, int a4);
void nox_xxx_mapgen_521FE0(int a1, int a2, uint32_t* a3);
uint32_t* nox_xxx_mapGenMakeSpellbook_5220E0(int a1, const char* a2);
uint32_t* nox_xxx_mapGenMakeEnchantedItem_5221A0(char* a1, char* a2, int a3);
uint32_t* sub_522300(int a1, uint32_t* a2);
void nox_xxx_mapgen_522340(int a1, int a2);
void nox_xxx_mapgen_522370(int a1, int a2, int* a3);
float* nox_xxx_mapgen_5224B0(int a1, int a2, int a3);
int nox_xxx_mapGenMakeMonsterInRoom_522810(float* a1, char* a2);
void nox_xxx_mapGenFinishPopulate_5228B0_mapgen_populate(int a1);
int nox_xxx_mapGenMakeExit_522A40(int a1);
float* nox_xxx_mapgen_522AD0(float* a1, int a2);
float* sub_522C80(float* a1);
float* sub_522D30(int a1);
unsigned char* nox_xxx_mapGenTryNextRoom_522F40(uint32_t* a1);
void nox_xxx_mapGenSetFlags_5235F0(char a1);
double sub_5259D0();
int sub_5259E0();
void sub_5259F0(int a1, int a2, float a3);
float* sub_525AF0(int a1);
void sub_525BF0(int a1);
float* sub_525C90();
int nox_xxx_mapGen_InPrefab1_525D20(int a1);
void sub_525DF0(float a1);
int nox_xxx_mapGenPrefabMkRoom_526100(int a1, int a2);
long long sub_526260(int a1, float* a2);
int sub_5262F0(int a1, int a2);
int sub_526550(int a1, int a2);
int nox_xxx_mapGen_InPrefab2_5266F0(int a1);
int nox_xxx_mapGenPlacePrefabs_526830(int a1);
int sub_5268F0(const char* a1);
char* sub_526950();
void sub_526A90();
char* sub_526AA0(int a1);
int sub_527D50(int a1, char* a2);
extern uint32_t dword_5d4594_1550916;
extern uint32_t dword_5d4594_2487564;
extern uint32_t dword_5d4594_2487568;
extern uint32_t dword_5d4594_2487576;
extern uint32_t dword_5d4594_2487580;
extern uint32_t dword_5d4594_2487584;
extern uint32_t dword_5d4594_2487620;
extern uint32_t dword_5d4594_2487624;
extern uint32_t dword_5d4594_2487628;
extern uint32_t dword_5d4594_2487632;
extern uint32_t dword_5d4594_2487652;
extern uint32_t dword_5d4594_2487656;
extern uint32_t dword_5d4594_2487672;
extern uint32_t dword_5d4594_2487676;
extern uint32_t dword_5d4594_2491616;
extern void* dword_5d4594_1599588;
extern void* dword_5d4594_1599592;
extern uint32_t dword_5d4594_1599576;
extern uint32_t dword_5d4594_1599596;
extern uint32_t dword_5d4594_1599480;
extern uint32_t dword_5d4594_1599476;
extern uint32_t dword_5d4594_1599540;
extern uint32_t dword_5d4594_3835396;
static float populationFloat(uint32_t x){float f; memcpy(&f,&x,4);return f;}
static uint32_t populationInvoke(int op,const uint32_t* v,uint32_t* high){
*high=0; switch(op){
case 0: return (uint32_t)sub_4D42C0();
case 1: return (uint32_t)nox_xxx_mapGenSpellIdByName_51E1D0((const char*)(uintptr_t)v[0]);
case 2: return (uint32_t)nox_xxx_mapgen_521C60((int)v[0],(int)v[1]);
case 3: return (uint32_t)sub_521CB0((int)v[0],(int)v[1],(int)v[2],(int)v[3]);
case 4: nox_xxx_mapgen_521FE0((int)v[0],(int)v[1],(uint32_t*)(uintptr_t)v[2]);return 0;
case 5: return (uint32_t)(uintptr_t)nox_xxx_mapGenMakeSpellbook_5220E0((int)v[0],(const char*)(uintptr_t)v[1]);
case 6: return (uint32_t)(uintptr_t)nox_xxx_mapGenMakeEnchantedItem_5221A0((char*)(uintptr_t)v[0],(char*)(uintptr_t)v[1],(int)v[2]);
case 7: return (uint32_t)(uintptr_t)sub_522300((int)v[0],(uint32_t*)(uintptr_t)v[1]);
case 8: nox_xxx_mapgen_522340((int)v[0],(int)v[1]);return 0;
case 9: nox_xxx_mapgen_522370((int)v[0],(int)v[1],(int*)(uintptr_t)v[2]);return 0;
case 10: return (uint32_t)(uintptr_t)nox_xxx_mapgen_5224B0((int)v[0],(int)v[1],(int)v[2]);
case 11: return (uint32_t)nox_xxx_mapGenMakeMonsterInRoom_522810((float*)(uintptr_t)v[0],(char*)(uintptr_t)v[1]);
case 12: nox_xxx_mapGenFinishPopulate_5228B0_mapgen_populate((int)v[0]);return 0;
case 13: return (uint32_t)nox_xxx_mapGenMakeExit_522A40((int)v[0]);
case 14: return (uint32_t)(uintptr_t)nox_xxx_mapgen_522AD0((float*)(uintptr_t)v[0],(int)v[1]);
case 15: return (uint32_t)(uintptr_t)sub_522C80((float*)(uintptr_t)v[0]);
case 16: return (uint32_t)(uintptr_t)sub_522D30((int)v[0]);
case 17: return (uint32_t)(uintptr_t)nox_xxx_mapGenTryNextRoom_522F40((uint32_t*)(uintptr_t)v[0]);
case 18: nox_xxx_mapGenSetFlags_5235F0((char)v[0]);return 0;
case 19: { double d=sub_5259D0();uint64_t b;memcpy(&b,&d,8);*high=b>>32;return (uint32_t)b;}
case 20: return (uint32_t)sub_5259E0();
case 21: sub_5259F0((int)v[0],(int)v[1],populationFloat(v[2]));return 0;
case 22: return (uint32_t)(uintptr_t)sub_525AF0((int)v[0]);
case 23: sub_525BF0((int)v[0]);return 0;
case 24: return (uint32_t)(uintptr_t)sub_525C90();
case 25: return (uint32_t)nox_xxx_mapGen_InPrefab1_525D20((int)v[0]);
case 26: sub_525DF0(populationFloat(v[0]));return 0;
case 27: return (uint32_t)nox_xxx_mapGenPrefabMkRoom_526100((int)v[0],(int)v[1]);
case 28: { uint64_t b=(uint64_t)sub_526260((int)v[0],(float*)(uintptr_t)v[1]);*high=b>>32;return (uint32_t)b;}
case 29: return (uint32_t)sub_5262F0((int)v[0],(int)v[1]);
case 30: return (uint32_t)sub_526550((int)v[0],(int)v[1]);
case 31: return (uint32_t)nox_xxx_mapGen_InPrefab2_5266F0((int)v[0]);
case 32: return (uint32_t)nox_xxx_mapGenPlacePrefabs_526830((int)v[0]);
case 33: return (uint32_t)sub_5268F0((const char*)(uintptr_t)v[0]);
case 34: return (uint32_t)(uintptr_t)sub_526950();
case 35: sub_526A90();return 0;
case 36: return (uint32_t)(uintptr_t)sub_526AA0((int)v[0]);
case 37: return (uint32_t)sub_527D50((int)v[0],(char*)(uintptr_t)v[1]);
default:abort();}}
extern uint32_t dword_5d4594_2487244;
extern uint32_t dword_5d4594_1599532;
extern uint32_t dword_5d4594_1599556;
extern uint32_t dword_5d4594_1599548;
extern uint32_t dword_5d4594_1599644;
extern uint32_t dword_5d4594_3835312;
static uint32_t* populationGlobal(int i){switch(i){
case 0:return &dword_5d4594_1550916;
case 1:return &dword_5d4594_2487564;
case 2:return &dword_5d4594_2487568;
case 3:return &dword_5d4594_2487576;
case 4:return &dword_5d4594_2487580;
case 5:return &dword_5d4594_2487584;
case 6:return &dword_5d4594_2487620;
case 7:return &dword_5d4594_2487624;
case 8:return &dword_5d4594_2487628;
case 9:return &dword_5d4594_2487632;
case 10:return &dword_5d4594_2487652;
case 11:return &dword_5d4594_2487656;
case 12:return &dword_5d4594_2487672;
case 13:return &dword_5d4594_2487676;
case 14:return &dword_5d4594_1599576;
case 15:return &dword_5d4594_1599596;
case 16:return &dword_5d4594_1599480;
case 17:return &dword_5d4594_1599476;
case 18:return &dword_5d4594_1599540;
case 19:return &dword_5d4594_3835396;
case 20:return &dword_5d4594_2487244;
case 21:return &dword_5d4594_1599532;
case 22:return &dword_5d4594_1599556;
case 23:return &dword_5d4594_1599548;
case 24:return &dword_5d4594_1599644;
case 25:return &dword_5d4594_3835312;
case 26:return (uint32_t*)&dword_5d4594_1599588;
case 27:return (uint32_t*)&dword_5d4594_1599592;
case 28:return &dword_5d4594_2491616;
default:abort();}}
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Dispatches the production implementation and snapshots shared engine owners.
func PortTestMapPopulation(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	ext := &paintTestExtension{globals: map[string]*uint32{}}
	for i, n := range []string{"dword_5d4594_1550916", "dword_5d4594_2487564", "dword_5d4594_2487568", "dword_5d4594_2487576", "dword_5d4594_2487580", "dword_5d4594_2487584", "dword_5d4594_2487620", "dword_5d4594_2487624", "dword_5d4594_2487628", "dword_5d4594_2487632", "dword_5d4594_2487652", "dword_5d4594_2487656", "dword_5d4594_2487672", "dword_5d4594_2487676", "dword_5d4594_1599576", "dword_5d4594_1599596", "dword_5d4594_1599480", "dword_5d4594_1599476", "dword_5d4594_1599540", "dword_5d4594_3835396", "dword_5d4594_2487244", "dword_5d4594_1599532", "dword_5d4594_1599556", "dword_5d4594_1599548", "dword_5d4594_1599644", "dword_5d4594_3835312", "dword_5d4594_1599588", "dword_5d4594_1599592", "dword_5d4594_2491616"} {
		ext.globals[n] = (*uint32)(unsafe.Pointer(C.populationGlobal(C.int(i))))
	}
	for i := 0; i < 5; i++ {
		ext.globals[fmt.Sprintf("roomGlobal%d", i)] = mapRoomGlobalWord(i)
	}
	for off := uint32(2487572); off <= 2487676; off += 4 {
		ext.globals[fmt.Sprintf("blob%d", off)] = memmap.PtrUint32(0x5D4594, uintptr(off))
	}
	ext.globals["waypointKind"] = memmap.PtrUint32(0x973F18, 35972)
	ext.globals["waypointAutoConnect"] = memmap.PtrUint32(0x973F18, 35976)
	ext.globals["progressThreshold"] = memmap.PtrUint32(0x587000, 254948)
	ext.globals["guiSequence"] = memmap.PtrUint32(0x5D4594, 1309668)
	for _, n := range []string{"returnHigh", "gameFlags", "ticks", "occupancyEnabled", "decodedCache", "prefabLoadCount", "prefabLoadIndex", "waypointHead", "waypointPending", "modifier", "modifierName"} {
		p, free := alloc.New(uint32(0))
		defer free()
		ext.globals[n] = p
	}
	for off := uintptr(1599484); off < 1599532; off += 4 {
		ext.globals[fmt.Sprintf("cacheBlob%d", off)] = memmap.PtrUint32(0x5D4594, off)
	}
	for off := uintptr(2491608); off <= 2491620; off += 4 {
		ext.globals[fmt.Sprintf("connection%d", off)] = memmap.PtrUint32(0x5D4594, off)
	}
	for _, off := range []uintptr{1599536, 1599544, 1599552, 1599560, 1599568, 1599572} {
		ext.globals[fmt.Sprintf("cacheBlob%d", off)] = memmap.PtrUint32(0x5D4594, off)
	}
	for _, off := range []uintptr{35920, 35924, 35928, 35932, 35936, 35940, 35944, 35960, 35964, 35968, 35980} {
		ext.globals[fmt.Sprintf("selection%d", off)] = memmap.PtrUint32(0x973F18, off)
	}
	for _, off := range []uintptr{739980, 739984, 1564960} {
		ext.globals[fmt.Sprintf("service%d", off)] = memmap.PtrUint32(0x5D4594, off)
	}
	ext.globals["transparentFloor"] = memmap.PtrUint32(0x587000, 229704)
	tables := blobdata.PortTestMapPopulationTables()
	for off, data := range tables {
		for i := 0; i < len(data); i += 4 {
			at := off + uintptr(i)
			ext.globals[fmt.Sprintf("populationTable%d", at)] = memmap.PtrUint32(0x587000, at)
		}
	}
	ext.reset = func(sp PortTestPaintSpec) {
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(sp.Globals["gameFlags"].Value))
	}
	ext.constants = map[uint32]uint32{mapRoomRaw(unsafe.Pointer(C.nox_xxx_XFerExit_4F4B90)): 0x70000003}
	var active *paintTestFixture
	var disposed map[*server.Object]bool
	var restoreModifier func()
	ext.setup = func(p *server.PortTestPaintOwners) func() {
		p.PopulationTypes(unsafe.Pointer(C.nox_xxx_XFerExit_4F4B90))
		oldRelease := mapRoomTestRelease
		mapRoomTestRelease = func(ptr unsafe.Pointer) {
			if active != nil {
				if r := active.known(ptr); r != nil {
					r.alive = false
					delete(active.owned, r)
				}
			}
		}
		oldLoad := populationTestLoad
		populationTestLoad = func(index int32) (int32, bool) {
			if *ext.globals["decodedCache"] == 0 {
				return 0, false
			}
			*ext.globals["prefabLoadCount"]++
			*ext.globals["prefabLoadIndex"] = uint32(index)
			if index < 0 || index >= int32(*ext.globals["dword_5d4594_1599596"]) {
				return 0, true
			}
			*ext.globals["dword_5d4594_3835396"] = uint32(index)
			*ext.globals["dword_5d4594_1599480"] = uint32(index)
			*ext.globals["dword_5d4594_1599476"] = 0
			return 0, true
		}
		ticks := PlatformTicks
		PlatformTicks = func() uint64 { return uint64(*ext.globals["ticks"]) }
		return func() { PlatformTicks = ticks; populationTestLoad = oldLoad; mapRoomTestRelease = oldRelease }
	}
	ext.finish = func(f *paintTestFixture) {
		if restoreModifier != nil {
			restoreModifier()
			restoreModifier = nil
		}
		f.owners.S.Nox_xxx_waypoint_5799C0()
		f.owners.S.Nox_xxx_waypointDeleteAll_579DD0()
		active = nil
	}
	ext.before = func(f *paintTestFixture, sp PortTestPaintSpec) {
		active = f
		disposed = map[*server.Object]bool{}
		if *ext.globals["modifier"] != 0 {
			restoreModifier = f.owners.S.PortTestRewardModifier((*server.ModifierEff)(mapRoomPointer(*ext.globals["modifier"])), (*byte)(mapRoomPointer(*ext.globals["modifierName"])))
		}
		for off, data := range tables {
			for i := 0; i < len(data); i += 4 {
				name := fmt.Sprintf("populationTable%d", off+uintptr(i))
				if _, ok := sp.Globals[name]; !ok {
					*ext.globals[name] = binary.LittleEndian.Uint32(data[i:])
				}
			}
		}
		if *ext.globals["occupancyEnabled"] != 0 {
			cfg := f.slots[1].ptr
			if *mapRoomWord(cfg, 68) > 32 {
				panic("population fixture grid capacity")
			}
			if mapRoomGridInit(cfg) == 0 {
				panic("population fixture grid initialization")
			}
			n := int(*mapRoomGlobalWord(2))
			r := f.register(mapRoomPointer(*mapRoomGlobalWord(0)), 4*n, "occupancyRows", false)
			f.owned[r] = true
			for x := 0; x < n; x++ {
				r = f.register(unsafe.Pointer(mapRoomGridRow(int32(x))), 20*n, "occupancy", false)
				f.owned[r] = true
			}
			for room := mapRoomHead(); room != nil; room = room.Next {
				mapRoomOccupy(room)
			}
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(*ext.globals["gameFlags"]))
		var released *mapRoomTestRegion
		if op == 35 || op == 12 {
			released = f.known(mapRoomPointer(*ext.globals["dword_5d4594_2487672"]))
		}
		cached := map[*mapRoomTestRegion]*mapRoomTestRegion{}
		if op == 32 && *ext.globals["decodedCache"] != 0 {
			for node := mapRoomPointer(*ext.globals["dword_5d4594_1599540"]); node != nil; node = *mapRoomRef(node, 4) {
				if rec := f.known(node); rec != nil {
					cached[rec] = f.known(*mapRoomRef(node, 0))
				}
			}
		}
		ret := uint32(C.populationInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0])), (*C.uint32_t)(unsafe.Pointer(ext.globals["returnHigh"]))))
		if len(cached) != 0 {
			remaining := map[unsafe.Pointer]bool{}
			for node := mapRoomPointer(*ext.globals["dword_5d4594_1599540"]); node != nil; node = *mapRoomRef(node, 4) {
				remaining[node] = true
			}
			for node, object := range cached {
				if !remaining[node.ptr] {
					node.alive = false
					delete(f.owned, node)
					if object != nil {
						object.alive = false
						disposed[(*server.Object)(object.ptr)] = true
					}
				}
			}
		}
		if released != nil {
			released.alive = false
			delete(f.owned, released)
		}
		return ret
	}
	ext.after = func(f *paintTestFixture, op int, ret uint32) {
		for _, head := range []*server.Waypoint{f.owners.S.WPs.List, f.owners.S.WPs.Pending} {
			for wp := head; wp != nil; wp = wp.WpNext {
				if f.known(wp.C()) == nil {
					f.register(wp.C(), int(unsafe.Sizeof(server.Waypoint{})), "waypoint", false)
				}
			}
		}
		*ext.globals["waypointHead"] = mapRoomRaw(unsafe.Pointer(f.owners.S.WPs.List))
		*ext.globals["waypointPending"] = mapRoomRaw(unsafe.Pointer(f.owners.S.WPs.Pending))
		if op == 34 {
			ptr := mapRoomPointer(*ext.globals["dword_5d4594_2487672"])
			if ptr != nil {
				n := int(*ext.globals["dword_5d4594_2487676"]) * 64
				r := f.register(ptr, n, "prefabMetadata", false)
				f.owned[r] = true
			}
		}
		visitRoom := func(ptr unsafe.Pointer) {
			if ptr != nil && f.known(ptr) == nil {
				r := f.register(ptr, 376, "input", false)
				f.owned[r] = true
			}
		}
		for room := mapRoomHead(); room != nil; room = room.Next {
			visitRoom(unsafe.Pointer(room))
		}
		for _, r := range append([]*mapRoomTestRegion(nil), f.regions...) {
			if r.alive && r.kind == "input" && r.size == 160 {
				visitRoom(*mapRoomRef(r.ptr, 148))
			}
		}
		switch op {
		case 5, 6, 7, 11, 14:
			if ret != 0 {
				f.trackObject((*server.Object)(mapRoomPointer(ret)))
			}
		}
		seen := map[*server.Object]bool{}
		var visit func(*server.Object)
		visit = func(u *server.Object) {
			if u == nil || seen[u] || disposed[u] {
				return
			}
			seen[u] = true
			f.trackObject(u)
			for c := u.InvFirstItem; c != nil; c = c.InvNextItem {
				if seen[c] {
					break
				}
				visit(c)
			}
		}
		for _, u := range f.owners.Objects() {
			visit(u)
		}
	}
	return portTestMapPainting(cases, owner, ext)
}
