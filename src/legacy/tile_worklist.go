package legacy

/*
#include "GAME4_2.h"
extern uint32_t dword_5d4594_2487248;
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

type tileFillEntry struct {
	X, Y, Flags uint32
}

func tileFillStack() *[500]tileFillEntry {
	return (*[500]tileFillEntry)(memmap.PtrOff(0x973F18, 16200))
}

func pushTileFill(x, y, flags, key int32) {
	if x <= 0 || x >= 127 || y <= 0 || y >= 127 || flags&3 == 0 {
		return
	}
	rows := (*[128]*[128]C.obj_5D4594_2650668_t)(unsafe.Pointer(C.ptr_5D4594_2650668))
	cell := &rows[x][y]
	if !(flags&2 != 0 && int32(cell.field_6) == key || flags&1 != 0 && int32(cell.field_1) == key) {
		return
	}
	if flags&1 != 0 && y == 1 || flags&2 != 0 && x == 1 {
		return
	}
	count := uint32(C.dword_5d4594_2487248)
	stack := tileFillStack()
	// C scans using a signed count but tests capacity using the unsigned word.
	for i := int32(0); i < int32(count); i++ {
		p := &stack[i]
		if p.X == uint32(x) && p.Y == uint32(y) && p.Flags == uint32(flags) {
			return
		}
	}
	if count >= 500 {
		*memmap.PtrUint32(0x973F18, 22200) = 1
		return
	}
	C.dword_5d4594_2487248 = C.uint32_t(count + 1)
	stack[count] = tileFillEntry{X: uint32(x), Y: uint32(y), Flags: uint32(flags)}
}

func popTileFill(x, y, flags *uint32) bool {
	count := uint32(C.dword_5d4594_2487248)
	if int32(count) <= 0 {
		return false
	}
	count--
	C.dword_5d4594_2487248 = C.uint32_t(count)
	stack := tileFillStack()
	*x = stack[count].X
	// Output pointers can alias the count or stack. Reload after each write,
	// as C does, instead of snapshotting a triplet before writing outputs.
	*y = stack[uint32(C.dword_5d4594_2487248)].Y
	*flags = stack[uint32(C.dword_5d4594_2487248)].Flags
	return true
}

//export sub_51DD50
func sub_51DD50(x, y, flags, key C.int) {
	pushTileFill(int32(x), int32(y), int32(flags), int32(key))
}

//export sub_51DE30
func sub_51DE30(x, y, flags *C.uint32_t) C.int {
	return C.int(bool2int(popTileFill((*uint32)(unsafe.Pointer(x)), (*uint32)(unsafe.Pointer(y)), (*uint32)(unsafe.Pointer(flags)))))
}
