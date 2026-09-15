package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export sub_430B50
func sub_430B50(x1, y1, x2, y2 C.int) C.int {
	return C.int(uiRenderBounds(int(x1), int(y1), int(x2), int(y2)))
}

//export sub_49CD30
func sub_49CD30(x, y, w, h, color, width C.int) {
	uiRenderBorder(int(x), int(y), int(w), int(h), int(color), int(width))
}

//export sub_49D1C0
func sub_49D1C0(p unsafe.Pointer, color, size C.int) C.int {
	uiRenderFill(p, uint32(color), int32(size))
	return 0
}

//export nox_client_copyRect_49F6F0
func nox_client_copyRect_49F6F0(x, y, w, h C.int) C.int {
	return C.int(uiRenderCopyRect(int(x), int(y), int(w), int(h)))
}
