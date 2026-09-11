package legacy

/*
#include "GAME4_3.h"
#include "noxstring.h"
extern uint32_t dword_5d4594_251572;
extern uint32_t dword_5d4594_2489436;
extern uint32_t dword_5d4594_3835356;
extern uint32_t dword_5d4594_3835360;
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func findBorderName(name *C.char) int32 {
	count := int32(C.dword_5d4594_251572)
	if name == nil || count <= 0 {
		return -1
	}
	want := C.GoString(name)
	// The loader caps active count at the physical 64 rows. Unlike tile
	// lookup, this scan is exact, case-sensitive and stops at the first match.
	for i := int32(0); i < count; i++ {
		p := (*C.char)(memmap.PtrOff(0x85B3FC, 28644+60*uintptr(i)))
		if C.GoString(p) == want {
			return i
		}
	}
	return -1
}

func selectBorderName(name *C.char) bool {
	C.dword_5d4594_2489436 = 0
	// Only the sentinel is case-insensitive; preserve the shared C locale
	// comparator rather than applying Unicode folding to border names.
	if C.nox_strcmpi((*C.char)(unsafe.Pointer(alloc.InternCString("NONE"))), name) == 0 {
		C.dword_5d4594_3835356 = 255
		return true
	}
	return selectBorderPrimary(findBorderName(name))
}

func selectBorderPrimary(index int32) bool {
	if index < 0 || index >= int32(C.dword_5d4594_251572) {
		return false
	}
	C.dword_5d4594_3835356 = C.uint32_t(index)
	C.dword_5d4594_2489436 = 1
	return true
}

func selectBorderVariation(variation int32) bool {
	if C.dword_5d4594_2489436 == 0 {
		return true
	}
	count := int32(C.dword_5d4594_251572)
	selected := uint32(C.dword_5d4594_3835356)
	// Approved fix: validate the selected border's limit, not the unrelated
	// row indexed by variation. Reject invalid selected rows before access.
	if count <= 0 || selected >= uint32(count) || selected >= 64 || variation < 0 ||
		variation >= int32(memmap.Uint16(0x85B3FC, 28688+60*uintptr(selected))) {
		return false
	}
	C.dword_5d4594_3835360 = C.uint32_t(variation)
	return true
}

//export sub_544020
func sub_544020(name *C.char) C.int {
	return C.int(bool2int(selectBorderName(name)))
}

//export nox_xxx_tileCheckByte3_544070
func nox_xxx_tileCheckByte3_544070(index C.int) C.int {
	return C.int(bool2int(selectBorderPrimary(int32(index))))
}

//export nox_xxx_tileCheckByte4_5440A0
func nox_xxx_tileCheckByte4_5440A0(variation C.int) C.int {
	return C.int(bool2int(selectBorderVariation(int32(variation))))
}
