package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

//export nox_xxx_wndEditProc_487D70
func nox_xxx_wndEditProc_487D70(w *nox_window, code, a, b C.int) C.int {
	return C.int(gui.EventRespInt(uiEntryInput(asWindow(w), &gui.RawEvent{Event: int(code), Arg1: uintptr(uint32(a)), Arg2: uintptr(uint32(b))})))
}

//export nox_xxx_wndEditDrawNoImage_488160
func nox_xxx_wndEditDrawNoImage_488160(w, d C.int) C.int {
	return C.int(uiEntryDraw((*gui.Window)(unsafe.Pointer(uintptr(uint32(w)))), (*gui.WindowData)(unsafe.Pointer(uintptr(uint32(d)))), false))
}
