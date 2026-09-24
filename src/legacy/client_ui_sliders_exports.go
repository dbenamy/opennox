package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func nox_xxx_wndScrollBoxDraw_4B4BA0(win, code int32, a uint32, b int32) int32 {
	// Keep the C caller's raw key-state and packed-coordinate words intact.
	ev := &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))}
	return int32(gui.EventRespInt(uiSliderInput((*gui.Window)(unsafe.Pointer(uintptr(uint32(win)))), ev, false)))
}
