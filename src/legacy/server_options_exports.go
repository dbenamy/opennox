package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_guiServerOptsLoad_457500
func nox_xxx_guiServerOptsLoad_457500() C.int { return C.int(serverOptionsConstruct()) }

//export nox_xxx_guiServerOptionsTryHide_4574D0
func nox_xxx_guiServerOptionsTryHide_4574D0() C.int { return C.int(serverOptionsTryClose()) }

func nox_xxx_guiServerOptionsHide_4597E0(value C.int) *C.int {
	return (*C.int)(unsafe.Pointer(serverOptionsClose(int(value))))
}

//export sub_457460
func sub_457460(data C.int) C.int {
	return C.int(serverOptionsLimits(serverOptionsRecord(unsafe.Pointer(uintptr(uint32(data))))))
}

//export sub_459AA0
func sub_459AA0(data unsafe.Pointer) *C.char {
	return (*C.char)(unsafe.Pointer(serverOptionsRead(serverOptionsRecord(data))))
}

func sub_459C30() C.int { return C.int(serverOptionsRefresh()) }

//export sub_459D50
func sub_459D50(value C.int) C.int { return C.int(serverOptionsDirty(int(value))) }

//export sub_459D80
func sub_459D80(value C.int) C.int { return C.int(serverOptionsVisible(int(value))) }

func sub_459DA0() C.int { return C.int(bool2int(serverOptionsRoot != 0)) }

// Tooltip callbacks are stored by the GUI through its existing C function-pointer ABI.
//
//export nox_xxx_options_457AA0
func nox_xxx_options_457AA0(_ C.int, data *C.uint8_t) C.int {
	return C.int(serverOptionsTooltip(false, byte(*data)))
}

//export nox_xxx_options_457B00
func nox_xxx_options_457B00(_ C.int, data *C.uint8_t) C.int {
	return C.int(serverOptionsTooltip(true, byte(*data)))
}
