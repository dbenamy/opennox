package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func sub_4C5020(packet C.int) C.int {
	return C.int(objectRenderBeamAppend(unsafe.Pointer(uintptr(packet))))
}
