package legacy

/*
#include <stdint.h>
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

//export nox_xxx_netGetUnitCodeCli_578B00
func nox_xxx_netGetUnitCodeCli_578B00(a1 C.int) C.uint {
	if a1 == 0 {
		return 0
	}
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a1))))
	code := dr.NetCode32
	if code >= 0x8000 {
		return 0
	}
	if uint32(dr.ObjClass)&0x20400000 != 0 {
		code |= 0x8000
	}
	return C.uint(code)
}

//export nox_xxx_netClearHighBit_578B30
func nox_xxx_netClearHighBit_578B30(a1 C.short) C.int {
	return C.int(uint16(a1) & 0x7fff)
}

//export nox_xxx_netTestHighBit_578B70
func nox_xxx_netTestHighBit_578B70(a1 C.uint) C.uint {
	return (a1 >> 15) & 1
}
