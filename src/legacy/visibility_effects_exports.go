package legacy

/*
#include "defs.h"
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
)

//export nox_xxx_netSendPointFx_522FF0
func nox_xxx_netSendPointFx_522FF0(code C.char, pos *C.float2) C.int {
	return C.int(visibilityFXPoint(byte(code), *(*types.Pointf)(unsafe.Pointer(pos))))
}

//export nox_xxx_sendArrowTrapFX_5238A0
func nox_xxx_sendArrowTrapFX_5238A0(pos *C.float, extra C.char) {
	visibilityFXArrowTrap(*(*types.Pointf)(unsafe.Pointer(pos)), byte(extra))
}

//export nox_xxx_netUpdateObjectSpecial_527E50
func nox_xxx_netUpdateObjectSpecial_527E50(a, b *C.nox_object_t) C.int {
	return C.int(visibilitySpecialUpdate((*server.Object)(unsafe.Pointer(a)), (*server.Object)(unsafe.Pointer(b))))
}

//export nox_xxx_frameCounterSetCopyToNextFrame_5281D0
func nox_xxx_frameCounterSetCopyToNextFrame_5281D0() C.int { return C.int(visibilityFrameCopy(true)) }

//export nox_xxx_netObjectOutOfSight_528A60
func nox_xxx_netObjectOutOfSight_528A60(to C.int, u *C.uint32_t) C.int {
	return C.int(visibilityOutOfSight(int(to), (*server.Object)(unsafe.Pointer(u))))
}

//export nox_xxx_netObjectInShadows_528A90
func nox_xxx_netObjectInShadows_528A90(to C.int, u *C.uint32_t) C.int {
	return C.int(visibilityInShadows(int(to), (*server.Object)(unsafe.Pointer(u))))
}
