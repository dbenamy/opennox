//go:build porttest

package legacy

/*
#include "defs.h"
#include <stdio.h>
#include <stdint.h>
extern uint32_t dword_5d4594_2487524;
int nox_xxx_mapGenReadTheme_51E260(int* a1, int a2);
int nox_xxx_mapGenReadLine_51E540(FILE* a1, uint8_t* a2);
int sub_51E570(FILE* a1, uint8_t* a2);
int sub_51E630(FILE* a1);
int sub_51E670(FILE* a1);
int sub_51E720(FILE* a1);
int sub_51E780(FILE* a1);
int sub_51E800(int a1, uint32_t* a2);
int sub_51EAF0(int a1, uint32_t* a2);
int nox_xxx_genReadAlgData_51EBB0(int a1, FILE* a2);
int nox_xxx_genReadSpellSet_51EFB0(int a1, FILE* a2);
int nox_xxx_genReadWeaponSet_51F030(int a1, FILE* a2);
void nox_xxx_mapGenFreeStr_51F1F0(void* lpMem);
int sub_51F230(int a1, FILE* a2);
int nox_xxx_genReadArmorSet_51F640(int a1, FILE* a2);
int nox_xxx_genReadExit_51F800(int a1, FILE* a2);
int nox_xxx_genReadDecor_51F9F0(uint32_t* a1, FILE* a2);
char* nox_xxx_genDecorReadWallFloor_51FE00(int a1, FILE* a2);
int sub_51FEC0(int a1, int a2, FILE* a3);
int nox_xxx_genDecorReadDecorSet_51FFA0(int a1, FILE* a2);
uint32_t* nox_xxx_gen_520380(FILE* a1);
uint32_t* nox_xxx_gen_5205B0(FILE* a1);
int nox_xxx_genDecorReadCopy_520660(uint32_t* a1, const char* a2, FILE* a3);
int nox_xxx_genDecorReadOccurConstraint_520810(int a1, FILE* a2);
int nox_xxx_genDecorReadOccurLimit_5208D0(int a1, FILE* a2);
int nox_xxx_genDecorReadFrequency_520910(int a1, FILE* a2);
int nox_xxx_genDecorReadRoomSizeCon_5209F0(int a1, FILE* a2);
int nox_xxx_genDecorReadDoor_520A90(int a1, FILE* a2);
int nox_xxx_genDecorReadDoubleDoor_520AB0(int a1, FILE* a2);
int nox_xxx_mapgenCheckSettings_520AD0(int* a1);
int nox_xxx_genReadPrefab_520BF0(int a1, FILE* a2);
char* sub_520CE0(int a1, FILE* a2);
uint32_t* sub_520D50(uint32_t* a1);
static uint32_t themeInvokeC(int op, uint32_t *v) {
 switch (op) {
case 0: return (uint32_t)(uintptr_t)(nox_xxx_mapGenReadTheme_51E260((int*)(uintptr_t)v[0], (int)(uintptr_t)v[1]));
case 1: return (uint32_t)(uintptr_t)(nox_xxx_mapGenReadLine_51E540((FILE*)(uintptr_t)v[0], (uint8_t*)(uintptr_t)v[1]));
case 2: return (uint32_t)(uintptr_t)(sub_51E570((FILE*)(uintptr_t)v[0], (uint8_t*)(uintptr_t)v[1]));
case 3: return (uint32_t)(uintptr_t)(sub_51E630((FILE*)(uintptr_t)v[0]));
case 4: return (uint32_t)(uintptr_t)(sub_51E670((FILE*)(uintptr_t)v[0]));
case 5: return (uint32_t)(uintptr_t)(sub_51E720((FILE*)(uintptr_t)v[0]));
case 6: return (uint32_t)(uintptr_t)(sub_51E780((FILE*)(uintptr_t)v[0]));
case 7: return (uint32_t)(uintptr_t)(sub_51E800((int)(uintptr_t)v[0], (uint32_t*)(uintptr_t)v[1]));
case 8: return (uint32_t)(uintptr_t)(sub_51EAF0((int)(uintptr_t)v[0], (uint32_t*)(uintptr_t)v[1]));
case 9: return (uint32_t)(uintptr_t)(nox_xxx_genReadAlgData_51EBB0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 10: return (uint32_t)(uintptr_t)(nox_xxx_genReadSpellSet_51EFB0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 11: return (uint32_t)(uintptr_t)(nox_xxx_genReadWeaponSet_51F030((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 12: nox_xxx_mapGenFreeStr_51F1F0((void*)(uintptr_t)v[0]); return 0;
case 13: return (uint32_t)(uintptr_t)(sub_51F230((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 14: return (uint32_t)(uintptr_t)(nox_xxx_genReadArmorSet_51F640((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 15: return (uint32_t)(uintptr_t)(nox_xxx_genReadExit_51F800((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 16: return (uint32_t)(uintptr_t)(nox_xxx_genReadDecor_51F9F0((uint32_t*)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 17: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadWallFloor_51FE00((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 18: return (uint32_t)(uintptr_t)(sub_51FEC0((int)(uintptr_t)v[0], (int)(uintptr_t)v[1], (FILE*)(uintptr_t)v[2]));
case 19: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadDecorSet_51FFA0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 20: return (uint32_t)(uintptr_t)(nox_xxx_gen_520380((FILE*)(uintptr_t)v[0]));
case 21: return (uint32_t)(uintptr_t)(nox_xxx_gen_5205B0((FILE*)(uintptr_t)v[0]));
case 22: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadCopy_520660((uint32_t*)(uintptr_t)v[0], (const char*)(uintptr_t)v[1], (FILE*)(uintptr_t)v[2]));
case 23: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadOccurConstraint_520810((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 24: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadOccurLimit_5208D0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 25: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadFrequency_520910((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 26: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadRoomSizeCon_5209F0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 27: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadDoor_520A90((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 28: return (uint32_t)(uintptr_t)(nox_xxx_genDecorReadDoubleDoor_520AB0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 29: return (uint32_t)(uintptr_t)(nox_xxx_mapgenCheckSettings_520AD0((int*)(uintptr_t)v[0]));
case 30: return (uint32_t)(uintptr_t)(nox_xxx_genReadPrefab_520BF0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 31: return (uint32_t)(uintptr_t)(sub_520CE0((int)(uintptr_t)v[0], (FILE*)(uintptr_t)v[1]));
case 32: return (uint32_t)(uintptr_t)(sub_520D50((uint32_t*)(uintptr_t)v[0]));
 }
 return 0;
}
*/
import "C"
import (
	"bytes"
	"io"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Reserved action values resolve to real file and token-buffer pointers.
const PortTestThemeFile = 0x60000001
const PortTestThemeToken = 0x60000002

func PortTestMapTheme(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	ext := &paintTestExtension{globals: map[string]*uint32{
		"themeLine":     memmap.PtrUint32(0x5D4594, 2487520),
		"themeTemplate": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487524)),
	}}
	for _, name := range []string{"themeInputPath", "themeFile", "themeOffset", "themeError", "themeClock", "themePlayers", "themeOpenFiles", "themeClosedFiles"} {
		p, free := alloc.New(uint32(0))
		defer free()
		ext.globals[name] = p
	}
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
	var input *binfile.Binfile
	var oldHandles map[unsafe.Pointer]bool
	var seenHandles map[unsafe.Pointer]bool
	var playersUnchanged func() bool
	var restorePlayers func()
	ext.setup = func(p *server.PortTestPaintOwners) func() {
		oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
		themeTestOnAllocate = func(ptr unsafe.Pointer, size int) {
			if active != nil {
				r := active.register(ptr, size, "themeAllocation", false)
				active.owned[r] = true
			}
		}
		themeTestOnRelease = func(ptr unsafe.Pointer) {
			if active != nil {
				if r := active.known(ptr); r != nil {
					r.alive = false
					delete(active.owned, r)
				}
			}
		}
		return func() { themeObserve(false, 0); themeTestOnAllocate = oldAlloc; themeTestOnRelease = oldFree }
	}
	ext.before = func(f *paintTestFixture, sp PortTestPaintSpec) {
		active = f
		seenHandles = map[unsafe.Pointer]bool{}
		if ptr := *ext.globals["themePlayers"]; ptr != 0 {
			count := *populationWord(ptr, 0)
			if count > 32 {
				panic("theme fixture player count")
			}
			classes := make([]byte, count)
			for i := range classes {
				classes[i] = byte(*populationWord(ptr, 4+4*i))
			}
			playersUnchanged, restorePlayers = f.owners.S.PortTestThemePlayers(classes)
		}
		clear(unsafe.Slice((*byte)(token), 256))
		f.register(token, 256, "themeToken", false)
		f.register(table, 368, "themeTables", false)
		oldHandles = map[unsafe.Pointer]bool{}
		files.RLock()
		for h := range files.byHandle {
			oldHandles[h] = true
		}
		files.RUnlock()
		input = nil
		if path := *ext.globals["themeInputPath"]; path != 0 {
			var err error
			input, err = binfile.BinfileOpen(populationString(path), binfile.ReadOnly)
			if err != nil {
				panic(err)
			}
			handle := NewFileHandle(input.File)
			*ext.globals["themeFile"] = mapRoomRaw(unsafe.Pointer(handle))
			f.register(unsafe.Pointer(handle), 0, "themeFile", false)
			seenHandles[unsafe.Pointer(handle)] = true
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		for i, v := range args {
			switch v {
			case PortTestThemeFile:
				args[i] = *ext.globals["themeFile"]
			case PortTestThemeToken:
				args[i] = mapRoomRaw(token)
			}
		}
		themeObserve(true, *ext.globals["themeClock"])
		defer themeObserve(false, 0)
		return uint32(C.themeInvokeC(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
	}
	ext.after = func(f *paintTestFixture, op int, ret uint32) {
		if input != nil {
			offset, err := input.File.File.Seek(0, io.SeekCurrent)
			if err != nil {
				panic(err)
			}
			*ext.globals["themeOffset"] = uint32(offset)
			*ext.globals["themeError"] = uint32(bool2int(input.File.Err != nil))
		}
		if playersUnchanged != nil {
			f.intact = playersUnchanged() && f.intact
		}
		*ext.globals["themeOpenFiles"], *ext.globals["themeClosedFiles"] = 0, 0
		files.RLock()
		defer files.RUnlock()
		for h, file := range files.byHandle {
			if !oldHandles[h] {
				if file.File.Fd() == ^uintptr(0) {
					*ext.globals["themeClosedFiles"]++
				} else {
					*ext.globals["themeOpenFiles"]++
				}
			}
			if !oldHandles[h] && !seenHandles[h] {
				seenHandles[h] = true
				r := f.register(h, 0, "themeFile", false)
				if file.File.Fd() == ^uintptr(0) {
					r.alive = false
				}
			}
		}
	}
	ext.finish = func(f *paintTestFixture) {
		themeObserve(false, 0)
		active = nil
		if restorePlayers != nil {
			restorePlayers()
			restorePlayers = nil
			playersUnchanged = nil
		}
		files.Lock()
		defer files.Unlock()
		for h, file := range files.byHandle {
			if !oldHandles[h] {
				_ = file.Close()
				delete(files.byHandle, h)
			}
		}
	}
	return portTestMapPainting(cases, owner, ext)
}
