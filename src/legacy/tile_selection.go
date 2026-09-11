package legacy

/*
#include "GAME4_1.h"
#include "noxstring.h"
extern nox_tileDef_t nox_tile_defs_arr[176];
extern uint32_t dword_5d4594_3835348;
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func tileDefinitionsAll() *[176]server.TileDef {
	return (*[176]server.TileDef)(unsafe.Pointer(&C.nox_tile_defs_arr[0]))
}

func selectTileName(name *C.char) bool {
	selected := memmap.PtrUint32(0x973F18, 35912)
	found := false
	// Scan physical storage, including entries beyond nox_tile_def_cnt. Keep
	// the shared C comparator until string handling is ported; it uses the C
	// locale's byte folding, rather than Unicode case folding.
	tiles := tileDefinitionsAll()
	for i := range tiles {
		p := &tiles[i]
		if C.nox_strcmpi((*C.char)(unsafe.Pointer(&p.NameBuf[0])), name) == 0 {
			*selected = uint32(i)
			found = true
		}
	}
	if C.nox_strcmpi(name, (*C.char)(unsafe.Pointer(alloc.InternCString("NONE")))) == 0 {
		*selected = 255
		return true
	}
	if !found {
		*selected = 0
	}
	return found
}

func selectTileImage(index int32) bool {
	selected := memmap.PtrUint32(0x973F18, 35912)
	if index < 0 || index >= 176 {
		*selected = 0
		return false
	}
	*selected = uint32(index)
	return true
}

func selectTileVariation(variation int32) bool {
	// Existing callers select/reset a valid tile first. NONE (255) is not a
	// valid tile for this operation in the original C either.
	p := &tileDefinitionsAll()[int32(memmap.Uint32(0x973F18, 35912))]
	if variation <= int32(p.Field52)*int32(p.Field53)-1 {
		C.dword_5d4594_3835348 = C.uint32_t(variation)
		return true
	}
	C.dword_5d4594_3835348 = 0
	return false
}

func setTileFlag(value int32) bool {
	if value != 0 && value != 1 {
		return false
	}
	*memmap.PtrUint32(0x973F18, 35916) = uint32(value)
	return true
}

//export nox_xxx_tileGetDefByName_51D4D0
func nox_xxx_tileGetDefByName_51D4D0(name *C.char) C.int {
	return C.int(bool2int(selectTileName(name)))
}

//export nox_xxx_tileCheckImage_51D540
func nox_xxx_tileCheckImage_51D540(index C.int) C.int {
	return C.int(bool2int(selectTileImage(int32(index))))
}

//export nox_xxx_tileCheckImageVari_51D570
func nox_xxx_tileCheckImageVari_51D570(variation C.int) C.int {
	return C.int(bool2int(selectTileVariation(int32(variation))))
}

//export nox_xxx_tile_51D5C0
func nox_xxx_tile_51D5C0(value C.int) C.int {
	return C.int(bool2int(setTileFlag(int32(value))))
}
