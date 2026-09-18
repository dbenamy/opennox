package legacy

/*
#include <stdint.h>
#include "defs.h"
#include "GAME1.h"
#include "GAME3_2.h"
#include "GAME4_3.h"
#include "server__script__file.h"
#include "noxstring.h"
extern uint32_t nox_tile_def_cnt;
extern uint32_t dword_5d4594_3835396, dword_5d4594_3835312;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func prefabGlobal(i int) *uint32 {
	switch i {
	case prefabSelected:
		return (*uint32)(unsafe.Pointer(&C.dword_5d4594_3835396))
	case prefabInstance:
		return (*uint32)(unsafe.Pointer(&C.dword_5d4594_3835312))
	default:
		return &prefabState[i]
	}
}
func prefabCompareNames(a, b uint32) int {
	return int(C.nox_strcmpi((*C.char)(mapRoomPointer(a)), (*C.char)(mapRoomPointer(b))))
}
func prefabTileDefinitionCount() int    { return int(C.nox_tile_def_cnt) }
func prefabSetBounds(bounds *[8]uint32) { C.sub_4D3C80((*C.uint32_t)(unsafe.Pointer(bounds))) }
func prefabInteresting(bounds *[4]int32, index uint32) uint32 {
	return uint32(C.nox_xxx_interesting_xfer_4D0010((*C.uint32_t)(unsafe.Pointer(bounds)), C.int(index)))
}
func prefabAdjustScript(instance uint32, dx, dy int32) {
	C.sub_542BF0(C.int(instance), C.int(dx), C.int(dy))
	offset := [2]int32{dx, dy}
	C.sub_543110((*C.char)(memmap.PtrOff(0x973F18, 30760)), (*C.int)(unsafe.Pointer(&offset)))
}
func prefabCombineScripts(a, b, c string) {
	x, freeX := alloc.CString(a)
	defer freeX()
	y, freeY := alloc.CString(b)
	defer freeY()
	z, freeZ := alloc.CString(c)
	defer freeZ()
	Nox_script_readWriteZzz_541670(x, y, z)
}
func prefabTransferWallData(src, dst *server.Wall) {
	if dst.Flags4&4 != 0 {
		C.sub_4107A0(dst.Data)
		dst.Data = nil
		dst.Flags4 &^= 4
	}
	if src.Flags4&12 != 0 {
		dst.Field10 = src.Field10
	}
	if src.Flags4&4 != 0 && src.Data != nil {
		dst.Flags4 |= 4
		dst.Data = src.Data
		src.Data = nil
		data := (*[8]uint32)(dst.Data)
		data[1], data[2], data[3] = uint32(dst.X5), uint32(dst.Y6), mapRoomRaw(dst.C())
		C.nox_xxx_wallSecretBlock_410760((*C.uint32_t)(dst.Data))
	}
	if src.Flags4&8 != 0 && dst.Flags4&8 == 0 {
		dst.Flags4 |= 8
		GetServer().S().Walls.AddBreakable(dst)
	}
}
