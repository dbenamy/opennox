//go:build porttest

package legacy

/*
#include "GAME1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
*/
import "C"

import (
	"math"
	"unsafe"
)

// PortTestGridBounds uses a null grid to assert rejection occurs before access.
// Inputs that expose the original C bug are run only in isolated child processes.
func PortTestGridBounds(xBits, yBits uint32) int {
	old := C.ptr_5D4594_2650668
	C.ptr_5D4594_2650668 = nil
	defer func() { C.ptr_5D4594_2650668 = old }()
	point := [2]C.float{C.float(math.Float32frombits(xBits)), C.float(math.Float32frombits(yBits))}
	return int(C.nox_xxx_tileNFromPoint_411160((*C.float2)(unsafe.Pointer(&point[0]))))
}
