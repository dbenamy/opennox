//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include <stdint.h>
extern uint32_t dword_5d4594_1549844,dword_5d4594_1550912,dword_5d4594_1550916;
static uint32_t* growthGlobal(int i) {
 switch(i) {case 0:return &dword_5d4594_1549844;case 1:return &dword_5d4594_1550912;case 2:return &dword_5d4594_1550916;}
 return 0;
}
static uint32_t growthInvoke(int op,uint32_t* v) {
 switch(op) {
 case 0:return (uintptr_t)nox_xxx_mapgen_Doors_4D4790();
 case 1:return (uintptr_t)nox_xxx_mapGenMkSmallRoom_4D4F40((uint32_t*)(uintptr_t)v[0]);
 case 2:sub_4D52F0();return 0;
 case 3:return sub_4D5350((uint32_t*)(uintptr_t)v[0],v[1],v[2],v[3],v[4]);
 case 4:return nox_xxx_mapGenFillRoom_4D53B0(v[0],v[1],v[2],v[3],v[4]);
 case 5:return sub_4D5630(v[0],v[1],v[2],v[3],v[4]);
 case 6:return sub_4D5D20((uint32_t*)(uintptr_t)v[0]);
 }
 return 0;
}
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func PortTestMapGrowth(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	ext := &paintTestExtension{globals: map[string]*uint32{}}
	for i, name := range []string{"growthMergeRate", "growthFrontier", "growthStart"} {
		ext.globals[name] = (*uint32)(unsafe.Pointer(C.growthGlobal(C.int(i))))
	}
	for i := 0; i < 5; i++ {
		ext.globals[fmt.Sprintf("roomGlobal%d", i)] = mapRoomGlobalWord(i)
	}
	ext.globals["waypointKind"] = memmap.PtrUint32(0x973F18, 35972)
	ext.globals["waypointAutoConnect"] = memmap.PtrUint32(0x973F18, 35976)
	var frees []func()
	for _, name := range []string{"growthInitGrid", "gameFlags", "waypointHead", "waypointPending"} {
		p, free := alloc.New(uint32(0))
		ext.globals[name] = p
		frees = append(frees, free)
	}
	defer func() {
		for _, free := range frees {
			free()
		}
	}()
	tables := blobdata.PortTestMapGrowthTables()
	for off, data := range tables {
		for i := 0; i < len(data); i += 4 {
			ext.globals[fmt.Sprintf("growthTable%d", off+uintptr(i))] = memmap.PtrUint32(0x587000, off+uintptr(i))
		}
	}
	cfg := memmap.PtrOff(0x5D4594, 1549796)
	saved := bytes.Clone(unsafe.Slice((*byte)(cfg), 1116))
	defer copy(unsafe.Slice((*byte)(cfg), 1116), saved)
	var active *paintTestFixture
	ext.setup = func(p *server.PortTestPaintOwners) func() {
		p.GrowthTypes()
		oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
		themeTestOnAllocate = func(ptr unsafe.Pointer, size int) {
			// Generator-owned room/exclusion records. Object-pool allocations retain
			// their server owner and are captured by the existing painting fixture.
			if active != nil && (size == 376 || size == 28) {
				r := active.register(ptr, size, "growthAllocation", false)
				active.owned[r] = true
			}
		}
		themeTestOnRelease = func(ptr unsafe.Pointer) {
			if active != nil {
				if r := active.known(ptr); r != nil {
					if r.size == 376 || r.size == 28 {
						trace := active.allocate(r.size+4, "growthReleased")
						words := unsafe.Slice((*uint32)(trace.ptr), r.size/4+1)
						words[0] = r.id
						for i, value := range unsafe.Slice((*uint32)(ptr), r.size/4) {
							words[i+1] = active.norm(value)
						}
					}
					r.alive = false
					delete(active.owned, r)
				}
			}
		}
		return func() { themeObserve(false, 0); themeTestOnAllocate = oldAlloc; themeTestOnRelease = oldFree }
	}
	ext.before = func(f *paintTestFixture, sp PortTestPaintSpec) {
		active = f
		copy(unsafe.Slice((*byte)(cfg), 1116), unsafe.Slice((*byte)(f.slots[1].ptr), 1116))
		f.register(cfg, 1116, "growthConfig", false)
		for off, data := range tables {
			for i := 0; i < len(data); i += 4 {
				name := fmt.Sprintf("growthTable%d", off+uintptr(i))
				if _, ok := sp.Globals[name]; !ok {
					*ext.globals[name] = binary.LittleEndian.Uint32(data[i:])
				}
			}
		}
		if *ext.globals["growthInitGrid"] != 0 {
			if radius := *mapRoomWord(cfg, 68); radius < 0 || radius > 32 {
				panic("growth grid bounds")
			}
			if mapRoomGridInit(cfg) == 0 {
				panic("growth grid allocation")
			}
			n := int(*mapRoomGlobalWord(2))
			r := f.register(mapRoomPointer(*mapRoomGlobalWord(0)), 4*n, "growthGridRows", false)
			f.owned[r] = true
			for i := 0; i < n; i++ {
				r := f.register(unsafe.Pointer(mapRoomGridRow(int32(i))), 20*n, "growthGrid", false)
				f.owned[r] = true
			}
			for r := mapRoomHead(); r != nil; r = r.Next {
				mapRoomOccupy(r)
			}
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(*ext.globals["gameFlags"]))
		themeObserve(true, 0)
		defer themeObserve(false, 0)
		return uint32(C.growthInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
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
	}
	ext.finish = func(f *paintTestFixture) {
		themeObserve(false, 0)
		active = nil
		f.owners.S.Nox_xxx_waypoint_5799C0()
		f.owners.S.Nox_xxx_waypointDeleteAll_579DD0()
	}
	return portTestMapPainting(cases, owner, ext)
}
