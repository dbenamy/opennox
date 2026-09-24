//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME3_2.h"

*/
import "C"
import (
	"bytes"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestMapCatalogEntry struct {
	Name    string
	Enabled int32
	Flags   uint32
	Tag     uint16
}
type PortTestMapCatalogState struct {
	Count   uint32
	Clock   uint32
	Last    uint32
	Quest   []byte
	Cycle   []byte
	Counts  [6]uint32
	Indices [6]uint32
}
type PortTestMapCatalog struct {
	restore []func()
	core    *server.Server
}

func PortTestMapCatalogOpen(seed int) *PortTestMapCatalog {
	f := &PortTestMapCatalog{}
	save := func(p unsafe.Pointer, n int) {
		raw := unsafe.Slice((*byte)(p), n)
		saved := bytes.Clone(raw)
		f.restore = append(f.restore, func() { copy(raw, saved) })
		clear(raw)
	}
	save(unsafe.Pointer(&legacyGlobals.nox_common_maplist), 12)
	listClear((*legacyListNode)(unsafe.Pointer(&legacyGlobals.nox_common_maplist)))
	save(unsafe.Pointer(&dword_5d4594_1548476), 4)
	save(unsafe.Pointer(&dword_5d4594_1548480), 4)
	save(memmap.PtrOff(0x5D4594, 1524108), 1024)
	save(memmap.PtrOff(0x5D4594, 1525132), 128*32)
	save(memmap.PtrOff(0x5D4594, 1529228), 6*25*128)
	save(memmap.PtrOff(0x5D4594, 1548428), 24)
	save(memmap.PtrOff(0x5D4594, 1548452), 24)
	save(memmap.PtrOff(0x5D4594, 1548484), 4)
	save(memmap.PtrOff(0x587000, 191832), 176)
	copy(unsafe.Slice(memmap.PtrUint8(0x587000, 191832), 176), blobdata.PortTestMapCatalogTable())
	for i, off := range []uintptr{191884, 191900, 191916, 191936, 191956, 191968} {
		*memmap.PtrPtr(0x587000, 191832+8*uintptr(i)) = memmap.PtrOff(0x587000, off)
	}
	old := GetServer
	core := new(server.Server)
	core.Rand.Logic = prand.New(seed)
	core.Rand.Other = prand.New(seed + 1)
	f.core = core
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	f.restore = append(f.restore, func() { GetServer = old })
	return f
}
func (f *PortTestMapCatalog) Close() {
	mapCatalogFree()
	for i := len(f.restore) - 1; i >= 0; i-- {
		f.restore[i]()
	}
}
func (f *PortTestMapCatalog) Add(v PortTestMapCatalogEntry) {
	p := (*Nox_map_list_item)(C.calloc(1, C.size_t(unsafe.Sizeof(Nox_map_list_item{}))))
	if p == nil {
		panic("map fixture allocation")
	}
	Sub_425770(p)
	copy(p.Name[:len(p.Name)-1], v.Name)
	p.Field_6 = v.Enabled
	p.Field_7 = v.Flags
	p.Field_8_2 = v.Tag
	Nox_common_maplist_add_4D0760(p)
}
func (f *PortTestMapCatalog) Entries() (out []PortTestMapCatalogEntry) {
	for p := C.nox_common_maplist_first_4D09B0(); p != nil; p = C.nox_common_maplist_next_4D09C0(p) {
		v := (*Nox_map_list_item)(unsafe.Pointer(p))
		out = append(out, PortTestMapCatalogEntry{GoStringP(unsafe.Pointer(&v.Name[0])), v.Field_6, v.Field_7, v.Field_8_2})
	}
	return
}
func (f *PortTestMapCatalog) BuildQuest() int { return mapQuestBuild() }
func (f *PortTestMapCatalog) ResetQuest()     { mapQuestReset() }
func (f *PortTestMapCatalog) ChooseQuest() (string, int) {
	p := C.nox_xxx_getQuestMapFile_4D0F60()
	return GoString(p), f.core.Rand.Logic.Index()
}
func (f *PortTestMapCatalog) Played(name *string) {
	if name == nil {
		mapQuestPlayed(nil)
	} else {
		mapQuestPlayed((*byte)(unsafe.Pointer(internCStr(*name))))
	}
}
func (f *PortTestMapCatalog) State() (s PortTestMapCatalogState) {
	s.Count = uint32(dword_5d4594_1548476)
	s.Clock = uint32(dword_5d4594_1548480)
	s.Last = *memmap.PtrUint32(0x587000, 191880)
	s.Quest = bytes.Clone(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1525132), 128*32))
	s.Cycle = bytes.Clone(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1529228), 6*25*128))
	copy(s.Counts[:], unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548428), 6))
	copy(s.Indices[:], unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548452), 6))
	return
}
func (f *PortTestMapCatalog) RestoreState(s PortTestMapCatalogState) {
	dword_5d4594_1548476 = uint32(s.Count)
	dword_5d4594_1548480 = uint32(s.Clock)
	*memmap.PtrUint32(0x587000, 191880) = s.Last
	copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1525132), 128*32), s.Quest)
	copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 1529228), 6*25*128), s.Cycle)
	copy(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548428), 6), s.Counts[:])
	copy(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1548452), 6), s.Indices[:])
}
func PortTestMapCycleGroup(name string) int {
	return mapCycleGroup(GoString((*C.char)(internCStr(name))))
}
func PortTestMapCycleMask(group int) uint32      { return mapCycleMask(group) }
func PortTestMapCycleFlagGroup(flags uint32) int { return mapCycleFlagGroup(flags) }
func PortTestMapCycleSetIndex(flags uint32, index int) int {
	return mapCycleSetIndex(flags, uint32(index))
}
func PortTestMapCycleGetIndex(flags uint32) int { return int(mapCycleGetIndex(flags)) }
func PortTestMapCycleStrip(data []byte) []byte {
	if data == nil {
		mapCycleStrip(nil)
		return nil
	}
	p := C.CBytes(data)
	defer C.free(p)
	mapCycleStrip((*byte)(p))
	return bytes.Clone(unsafe.Slice((*byte)(p), len(data)))
}
func PortTestMapCycleNext() string     { return GoStringP(unsafe.Pointer(mapCycleNext())) }
func PortTestMapCycleLoad()            { mapCycleLoad() }
func PortTestMapCycleEnabled() int     { return int(C.sub_4D0D70()) }
func PortTestMapCycleEnable(v int) int { return int(C.sub_4D0D90(C.int(v))) }
func PortTestMapCycleReset()           { mapCycleReset() }

func PortTestMapCyclePath() string { return GoStringP(memmap.PtrOff(0x5D4594, 1524108)) }
func PortTestMapCycleOpenHandles() int {
	files.RLock()
	defer files.RUnlock()
	return len(files.byHandle)
}

func PortTestMapCatalogNilNext() bool { return C.nox_common_maplist_next_4D09C0(nil) == nil }
