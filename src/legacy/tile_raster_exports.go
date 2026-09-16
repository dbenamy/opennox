package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// Both remaining C callbacks use the actual three-argument void dispatch ABI.
//
//export nox_xxx_tileDraw_4815E0
func nox_xxx_tileDraw_4815E0(pos *C.int2, img, tile C.uint32_t) {
	tileRasterTexture(*(*image.Point)(unsafe.Pointer(pos)), noxrender.ImageHandle(unsafe.Pointer(uintptr(img))))
}

//export sub_481770
func sub_481770(pos *C.int2, img, tile C.uint32_t) {
	tileRasterFill(*(*image.Point)(unsafe.Pointer(pos)), uint16(tile))
}
