package legacy

/*
#include <stdlib.h>
#include <stdint.h>
#include <time.h>
extern uint32_t dword_5d4594_2487524;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"os"
	"unsafe"
)

// Theme records remain owned by the shared C allocator while the rest of the
// map generator consumes them. Numeric conversion and clock retain libc semantics.
func mapThemeAlloc(n, size uint32) uint32 { return mapRoomRaw(C.calloc(C.size_t(n), C.size_t(size))) }
func mapThemeFree(p uint32)               { C.free(mapRoomPointer(p)) }
func mapThemeTemplate() *uint32           { return (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487524)) }
func mapThemeInt(s string) int32 {
	b := append([]byte(s), 0)
	return int32(C.atoi((*C.char)(unsafe.Pointer(&b[0]))))
}
func mapThemeFloat(s string) float64 {
	b := append([]byte(s), 0)
	return float64(C.atof((*C.char)(unsafe.Pointer(&b[0]))))
}
func mapThemeTable(off uintptr, s string) int {
	for i := 0; ; i++ {
		p := *memmap.PtrUint32(0x587000, off+uintptr(4*i))
		if p == 0 {
			return -1
		}
		if mapThemeLower(populationString(p)) == mapThemeLower(s) {
			return i
		}
	}
}
func mapThemeReadByte(f uint32) (byte, bool) {
	file := fileByHandle((*FILE)(mapRoomPointer(f)))
	var b [1]byte
	_, err := file.Bin.Read(b[:])
	if err != nil {
		file.Err = err
	}
	return b[0], file.Err == nil
}
func mapThemeFile(cfg, name uint32) uint32 {
	clear(unsafe.Slice(mapThemeByte(cfg, 0), 1116))
	defaults := map[int]uint32{4: 5, 8: 2, 12: 3, 16: 1, 20: 3, 24: 40, 28: 20, 32: 10, 36: 5, 40: 20, 48: 100, 52: 100, 56: 1, 64: 1157234688, 72: 100, 76: uint32(C.time(nil))}
	for off, v := range defaults {
		*populationWord(cfg, off) = v
	}
	*populationBlob(2487520) = 0
	file, err := binfile.BinfileOpen("mapgen/"+populationString(name)+".thm", binfile.ReadOnly)
	if err != nil {
		if !os.IsNotExist(err) {
			binfile.Log.Println(err)
		}
		return 0
	}
	f := mapRoomRaw(unsafe.Pointer(NewFileHandle(file.File)))
	if err = file.SetKey(1); err != nil {
		binfile.Log.Println(err)
		return 0
	}
	result := uint32(1)
	for mapThemeNext(f) {
		switch mapThemeLower(mapThemeText()) {
		case "algorithm_data":
			result = mapThemeAlgorithm(cfg, f)
		case "exit":
			result = mapThemeExit(cfg, f)
		case "prefab":
			result = mapThemePrefabs(cfg, f)
		case "spell_set":
			result = mapThemeSpells(cfg, f)
		case "weapon_set":
			result = mapThemeEquipment(cfg, f, false)
		case "armor_set":
			result = mapThemeEquipment(cfg, f, true)
		case "decor":
			result = mapThemeDecor(cfg, f)
		case "ambient_light":
			for _, off := range []int{536, 540, 544} {
				if !mapThemeNext(f) {
					result = 0
					break
				}
				*populationWord(cfg, off) = uint32(mapThemeInt(mapThemeText()))
			}
		default:
			result = 0
		}
		if result == 0 {
			break
		}
	}
	file.Close()
	if result != 1 {
		return result
	}
	if mapThemeValidate(cfg+88) == 0 || mapThemeValidate(cfg+120) == 0 {
		return 0
	}
	if *populationWord(cfg, 184) != 0 && mapThemeValidate(cfg+184) == 0 {
		return 0
	}
	return 1
}

//export nox_xxx_mapGenReadTheme_51E260
func nox_xxx_mapGenReadTheme_51E260(cfg *C.int, name C.int) C.int {
	return C.int(mapThemeFile(mapRoomRaw(unsafe.Pointer(cfg)), uint32(name)))
}

//export sub_520D50
func sub_520D50(cfg *C.uint32_t) *C.uint32_t {
	mapThemeCleanup(mapRoomRaw(unsafe.Pointer(cfg)))
	return nil
}
