package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_traceRay_5374B0
func nox_xxx_traceRay_5374B0(ray *C.float4) C.int {
	return C.int(bool2int(spatialRay((*[4]float32)(unsafe.Pointer(ray)))))
}
