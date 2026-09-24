package legacy

/*
#include "defs.h"
*/
import "C"

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
