//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include <stdint.h>
extern uint32_t dword_5d4594_1599596,dword_5d4594_2487524;
extern uint32_t dword_5d4594_1549844,dword_5d4594_1550912,dword_5d4594_1550916;
static uint32_t* orchestrationGlobal(int i) {
 switch(i) {case 0:return &dword_5d4594_1549844;case 1:return &dword_5d4594_1550912;case 2:return &dword_5d4594_1550916;}
 return 0;
}
static uint32_t orchestrationInvoke(int op,uint32_t* args) {
 switch(op) {
 case 0:return sub_4D42E0((char*)(uintptr_t)args[0]);
 case 1:return (uintptr_t)nox_xxx_getRandMapName_4D4310();
 case 2:return nox_xxx_mapGenStep_4D44E0();
 case 3:return nox_xxx_mapGenStart_4D4320();
 case 4:return nox_xxx_mapGenStartAlt_4D5F30();
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

type PortTestMapOrchestrationServices struct {
	Progress func(string)
	Before   func(*server.Server)
	After    func(*server.Server, int, uint32)
	Save     func(string, int) int
}

func PortTestMapOrchestration(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func()), services ...*PortTestMapOrchestrationServices) []PortTestPaintResult {
	var svc *PortTestMapOrchestrationServices
	if len(services) != 0 {
		svc = services[0]
	}
	if svc != nil {
		oldSave := orchestrationTestSave
		orchestrationTestSave = svc.Save
		defer func() { orchestrationTestSave = oldSave }()
		metadata := memmap.PtrOff(0x973F18, 2408)
		savedMetadata := bytes.Clone(unsafe.Slice((*byte)(metadata), 1464))
		defer copy(unsafe.Slice((*byte)(metadata), 1464), savedMetadata)
		savedObjects := *memmap.PtrUint32(0x5D4594, 1550924)
		defer func() { *memmap.PtrUint32(0x5D4594, 1550924) = savedObjects }()
	}

	ext := &paintTestExtension{globals: map[string]*uint32{}}
	ext.globals["areaCount"] = (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599596))
	ext.globals["themeLine"] = populationBlob(2487520)
	ext.globals["themeTemplate"] = (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487524))
	for i := 0; i < 3; i++ {
		ext.globals[fmt.Sprintf("ambient%d", i)] = memmap.PtrUint32(0x587000, 142296+uintptr(4*i))
	}
	oldDebug := Sub_57C490_2
	Sub_57C490_2 = func(string) {}
	if svc != nil && svc.Progress != nil {
		Sub_57C490_2 = svc.Progress
	}
	defer func() { Sub_57C490_2 = oldDebug }()
	name := memmap.PtrOff(0x587000, 197860)
	savedName := bytes.Clone(unsafe.Slice((*byte)(name), 64))
	defer copy(unsafe.Slice((*byte)(name), 64), savedName)

	for i, name := range []string{"growthMergeRate", "growthFrontier", "growthStart"} {
		ext.globals[name] = (*uint32)(unsafe.Pointer(C.orchestrationGlobal(C.int(i))))
	}
	for i := 0; i < 14; i++ {
		ext.globals[fmt.Sprintf("populationGlobal%d", i)] = populationGlobal(i)
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
	for off, data := range blobdata.PortTestMapPopulationTables() {
		tables[off] = data
	}
	for off, data := range tables {
		for i := 0; i < len(data); i += 4 {
			ext.globals[fmt.Sprintf("growthTable%d", off+uintptr(i))] = memmap.PtrUint32(0x587000, off+uintptr(i))
		}
	}
	cfg := memmap.PtrOff(0x5D4594, 1549796)
	saved := bytes.Clone(unsafe.Slice((*byte)(cfg), 1116))
	defer copy(unsafe.Slice((*byte)(cfg), 1116), saved)
	token := memmap.PtrOff(0x5D4594, 2487264)
	table := memmap.PtrOff(0x587000, 253144)
	oldToken := bytes.Clone(unsafe.Slice((*byte)(token), 256))
	oldTable := bytes.Clone(unsafe.Slice((*byte)(table), 368))
	defer func() {
		copy(unsafe.Slice((*byte)(token), 256), oldToken)
		copy(unsafe.Slice((*byte)(table), 368), oldTable)
	}()
	copy(unsafe.Slice((*byte)(table), 368), blobdata.PortTestMapThemeTableData())
	*memmap.PtrPtr(0x587000, 253144) = memmap.PtrOff(0x587000, 253296)
	*memmap.PtrPtr(0x587000, 253148) = memmap.PtrOff(0x587000, 253304)
	*memmap.PtrPtr(0x587000, 253152) = memmap.PtrOff(0x587000, 253312)
	*memmap.PtrPtr(0x587000, 253156) = memmap.PtrOff(0x587000, 253320)
	*memmap.PtrPtr(0x587000, 253160) = memmap.PtrOff(0x587000, 253328)
	*memmap.PtrPtr(0x587000, 253164) = memmap.PtrOff(0x587000, 253336)
	*memmap.PtrPtr(0x587000, 253172) = memmap.PtrOff(0x587000, 253340)
	*memmap.PtrPtr(0x587000, 253176) = memmap.PtrOff(0x587000, 253348)
	*memmap.PtrPtr(0x587000, 253180) = memmap.PtrOff(0x587000, 253356)
	*memmap.PtrPtr(0x587000, 253184) = memmap.PtrOff(0x587000, 253364)
	*memmap.PtrPtr(0x587000, 253188) = memmap.PtrOff(0x587000, 253372)
	*memmap.PtrPtr(0x587000, 253192) = memmap.PtrOff(0x587000, 253380)
	*memmap.PtrPtr(0x587000, 253200) = memmap.PtrOff(0x587000, 253388)
	*memmap.PtrPtr(0x587000, 253204) = memmap.PtrOff(0x587000, 253396)
	*memmap.PtrPtr(0x587000, 253208) = memmap.PtrOff(0x587000, 253412)
	*memmap.PtrPtr(0x587000, 253216) = memmap.PtrOff(0x587000, 253432)
	*memmap.PtrPtr(0x587000, 253220) = memmap.PtrOff(0x587000, 253440)
	*memmap.PtrPtr(0x587000, 253224) = memmap.PtrOff(0x587000, 253448)
	*memmap.PtrPtr(0x587000, 253228) = memmap.PtrOff(0x587000, 253464)
	*memmap.PtrPtr(0x587000, 253232) = memmap.PtrOff(0x587000, 253472)
	*memmap.PtrPtr(0x587000, 253236) = memmap.PtrOff(0x587000, 253480)
	*memmap.PtrPtr(0x587000, 253244) = memmap.PtrOff(0x587000, 253488)
	*memmap.PtrPtr(0x587000, 253248) = memmap.PtrOff(0x587000, 253496)

	var active *paintTestFixture
	var gridRowIDs map[uint32]uint32
	var traceBlocks [][]uint32
	ext.setup = func(p *server.PortTestPaintOwners) func() {
		p.OrchestrationTypes()
		oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
		themeTestOnAllocate = func(ptr unsafe.Pointer, size int) {
			// Generator-owned room/exclusion records. Object-pool allocations retain
			// their server owner and are captured by the existing painting fixture.
			if active != nil && orchestrationAllocationSize(size) {
				r := active.register(ptr, size, "orchestrationAllocation", false)
				active.owned[r] = true
			}
		}
		themeTestOnRelease = func(ptr unsafe.Pointer) {
			if active != nil {
				if r := active.known(ptr); r != nil {
					n := int(*mapRoomGlobalWord(2))
					table := mapRoomPointer(*mapRoomGlobalWord(0))
					isGridIndex := table == ptr && n > 0 && n <= 65 && r.size == 4*n
					if !isGridIndex && n > 0 && n <= 65 && r.size == 20*n && table != nil {
						for _, row := range unsafe.Slice((*uint32)(table), n) {
							if row == mapRoomRaw(ptr) {
								gridRowIDs[row] = r.id
								break
							}
						}
					}
					if orchestrationAllocationSize(r.size) {
						// Keep diagnostic storage outside the engine's C allocator so
						// recording releases cannot reuse a disposed room or grid row.
						buffer := make([]uint32, r.size/4+5)
						for i := len(buffer) - 4; i < len(buffer); i++ {
							buffer[i] = 0x5a5a5a5a
						}
						traceBlocks = append(traceBlocks, buffer)
						active.register(unsafe.Pointer(&buffer[0]), r.size+4, "orchestrationReleased", true).captureOnly = true
						words := buffer[:len(buffer)-4]
						words[0] = r.id
						for i, value := range unsafe.Slice((*uint32)(ptr), r.size/4) {
							// Rows are freed before their index. Keep their identities from
							// before release; trace allocations can reuse those addresses.
							if isGridIndex && value != 0 {
								id, ok := gridRowIDs[value]
								if !ok {
									panic("untracked released generator grid row")
								}
								words[i+1] = id
							} else {
								words[i+1] = active.norm(value)
							}
						}
					}
					if isGridIndex {
						clear(gridRowIDs)
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
		gridRowIDs = make(map[uint32]uint32)
		traceBlocks = nil
		clear(unsafe.Slice((*byte)(name), 64))
		clear(unsafe.Slice((*byte)(token), 256))
		f.register(token, 256, "themeToken", false)
		f.register(table, 368, "themeTables", false)
		f.register(name, 64, "mapName", false)
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
	previousBefore := ext.before
	ext.before = func(f *paintTestFixture, sp PortTestPaintSpec) {
		previousBefore(f, sp)
		if svc != nil {
			clear(unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464))
			*memmap.PtrUint32(0x5D4594, 1550924) = 0
			if svc.Before != nil {
				svc.Before(f.owners.S)
			}
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(*ext.globals["gameFlags"]))
		themeObserve(true, 1700000000)
		defer themeObserve(false, 0)
		return uint32(C.orchestrationInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
	}
	ext.after = func(f *paintTestFixture, op int, ret uint32) {
		if svc != nil && svc.After != nil {
			svc.After(f.owners.S, op, ret)
		}
		*ext.globals["gameFlags"] = uint32(noxflags.GetGame())
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

// Register generator-owned allocations, including the currently allocated grid.
// Server objects and their update records retain their normal owner.
func orchestrationAllocationSize(size int) bool {
	if n := int(*mapRoomGlobalWord(2)); n > 0 && n <= 65 && (size == 4*n || size == 20*n) {
		return true
	}
	switch size {
	case 0, 376, 28, 8192, 224, 128, 24, 100, 2060, 12, 160, 156:
		return true
	}
	return false
}
