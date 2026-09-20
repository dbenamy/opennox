//go:build porttest

package legacy

/*
#include "GAME1.h"
*/
import "C"

import "unsafe"

func PortTestWorldWallAttach(drawable bool, value unsafe.Pointer, x, y int) unsafe.Pointer {
	if drawable {
		return unsafe.Pointer(C.sub_410390(C.int(uintptr(value)), C.int(x), C.int(y)))
	}
	return unsafe.Pointer(C.nox_xxx_doorAttachWall_410360(C.int(uintptr(value)), C.int(x), C.int(y)))
}
