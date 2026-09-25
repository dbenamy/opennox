package legacy

/*
#include "defs.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"unsafe"
)

//export sub_4CA720
func sub_4CA720(unused, drawable C.int) C.int {
	return C.int(effectOrbitUpdate((*client.Drawable)(unsafe.Pointer(uintptr(uint32(drawable))))))
}
