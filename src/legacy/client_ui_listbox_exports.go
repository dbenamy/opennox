package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

//export nox_xxx_wndListboxProcWithoutData10_4A28E0
func nox_xxx_wndListboxProcWithoutData10_4A28E0(w *C.uint, code C.int, a C.uint, b C.int) C.int {
	return C.int(gui.EventRespInt(uiListSingleInput((*gui.Window)(unsafe.Pointer(w)), &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))})))
}

//export nox_xxx_wndListboxProcPre_4A30D0
func nox_xxx_wndListboxProcPre_4A30D0(w *nox_window, code, a C.uint, b C.int) C.int {
	return C.int(gui.EventRespInt(uiListEvent(asWindow(w), &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))})))
}

//export sub_4A4800
func sub_4A4800(d C.int) C.int {
	return C.int(uiListIndex((*gui.ScrollListBoxData)(unsafe.Pointer(uintptr(uint32(d))))))
}

//export nox_gui_newScrollListBox_4A4310
func nox_gui_newScrollListBox_4A4310(parent *nox_window, flags, x, y, width, height, draw C.int, opts *C.nox_scrollListBox_data) *nox_window {
	return (*nox_window)(uiListNew(asWindow(parent), gui.StatusFlags(flags), int(x), int(y), int(width), int(height), (*gui.WindowData)(unsafe.Pointer(uintptr(uint32(draw)))), (*gui.ScrollListBoxData)(unsafe.Pointer(opts))).C())
}
