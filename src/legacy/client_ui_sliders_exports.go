package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func nox_xxx_wndScrollBoxDraw_4B4BA0(win, code C.int, a C.uint, b C.int) C.int {
	// Keep the C caller's raw key-state and packed-coordinate words intact.
	ev := &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))}
	return C.int(gui.EventRespInt(uiSliderInput((*gui.Window)(unsafe.Pointer(uintptr(uint32(win)))), ev, false)))
}
