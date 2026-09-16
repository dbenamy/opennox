package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

//export nox_xxx_updDrawColorlight_4CE390
func nox_xxx_updDrawColorlight_4CE390(vp *C.uint32_t, dr C.int) C.int {
	return C.int(colorLightUpdate((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(uintptr(uint32(dr))))))
}
