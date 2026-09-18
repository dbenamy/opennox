package legacy

/*
#include "defs.h"
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/libs/types"
)

//export nox_xxx_netSendPointFx_522FF0
func nox_xxx_netSendPointFx_522FF0(code C.char, pos *C.float2) C.int {
	return C.int(visibilityFXPoint(byte(code), *(*types.Pointf)(unsafe.Pointer(pos))))
}

//export nox_xxx_sendArrowTrapFX_5238A0
func nox_xxx_sendArrowTrapFX_5238A0(pos *C.float, extra C.char) {
	visibilityFXArrowTrap(*(*types.Pointf)(unsafe.Pointer(pos)), byte(extra))
}
