package legacy

/*
#include <stdint.h>
*/
import "C"

import "unsafe"

func nox_xxx_packetDynamicUnitCode_578B40(value C.int) C.int {
	return C.int(networkDynamicUnitCode(uint32(value)))
}

func networkDynamicUnitCode(code uint32) uint32 {
	if code&0x8000 == 0 {
		return code
	}
	obj := GetServer().S().Objs.GetObjectByInd(int(code &^ 0x8000))
	if obj == nil {
		return 0
	}
	return obj.NetCode
}

//export nox_xxx_netGetUnitByExtent_4ED020
func nox_xxx_netGetUnitByExtent_4ED020(value C.int) C.int {
	obj := GetServer().S().Objs.GetObjectByInd(int(uint32(value)))
	// Preserve the raw 32-bit address ABI used by remaining decompiled C callers.
	return C.int(uintptr(unsafe.Pointer(obj)))
}
