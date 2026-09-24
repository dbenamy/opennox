package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func nox_xxx_guiServerOptionsHide_4597E0(value C.int) *C.int {
	return (*C.int)(unsafe.Pointer(serverOptionsClose(int(value))))
}

func sub_459C30() C.int { return C.int(serverOptionsRefresh()) }

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
