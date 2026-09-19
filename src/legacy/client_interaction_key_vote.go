package legacy

/*
#include "defs.h"
extern uint32_t dword_5d4594_1319056, dword_5d4594_1319060;
extern uint32_t dword_5d4594_1321216;
extern int nox_win_width,nox_win_height;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"unsafe"
)

func interactionKeyRoot() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1319060)))
}
func interactionKeyState(v uint32) {
	if C.dword_5d4594_1319056 == 0 && v == 1 {
		Nox_xxx_clientPlaySoundSpecial_452D80(1022, 100)
	}
	C.dword_5d4594_1319056 = C.uint32_t(v)
}
func interactionKeyShow() int {
	if w := interactionKeyRoot(); w != nil {
		w.Show()
		uiWindowEnable(w, 1)
		w.SetPos(image.Pt(int(C.nox_win_width)/2-w.SizeVal.X/2, int(C.nox_win_height)/2-w.SizeVal.Y/2))
		GetClient().Cli().GUI.Focus(nil)
	}
	return 0
}
func interactionKeyHide() int {
	if w := interactionKeyRoot(); w != nil {
		w.Hide()
		GetClient().Cli().GUI.Focus(nil)
	}
	return 0
}
func interactionKeyUpdate(value uintptr) {
	state := uint32(C.dword_5d4594_1319056)
	if state != 0 {
		if state == 1 && value == 0 {
			interactionKeyHide()
			interactionKeyState(0)
		}
	} else if value == 1 {
		interactionKeyShow()
		interactionKeyState(1)
	}
}
func interactionKeyEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16391 {
		a, _ := ev.EventArgsC()
		id := -2
		if child := (*gui.Window)(unsafe.Pointer(a)); child != nil {
			id = int(child.ID())
		}
		Nox_xxx_clientPlaySoundSpecial_452D80(766, 100)
		if id == 10803 {
			interactionKeyHide()
		}
	}
	return nil
}
func interactionKeyOpen() int {
	w := Nox_new_window_from_file("SKey.wnd", interactionKeyEvent)
	C.dword_5d4594_1319060 = C.uint32_t(uintptr(w.C()))
	if w == nil {
		return 0
	}
	interactionKeyState(0)
	interactionKeyHide()
	return 1
}
func interactionKeyDestroy() {
	nox_xxx_windowDestroyMB_46C4E0((*nox_window)(interactionKeyRoot().C()))
	C.dword_5d4594_1319060 = 0
	interactionKeyState(0)
}
func interactionVoteOpen() int {
	img := Nox_xxx_gLoadImg("VoteInProgress")
	*memmap.PtrUint32(0x5D4594, 1321220) = uint32(uintptr(img.C()))
	w := GetClient().Cli().GUI.NewWindowRaw(nil, 136, int(C.nox_win_width)-50, int(C.nox_win_height)/2-100, 50, 50, nil)
	C.dword_5d4594_1321216 = C.uint32_t(uintptr(w.C()))
	w.DrawData().BgImageHnd = img.C()
	w.SetAllFuncs(nil, interactionIconDraw, nil)
	w.Hide()
	return 1
}
func interactionVoteHide(v int) int {
	return nox_window_set_hidden((*nox_window)(unsafe.Pointer(uintptr(C.dword_5d4594_1321216))), v)
}
func interactionVoteDestroy() int {
	result := nox_xxx_windowDestroyMB_46C4E0((*nox_window)(unsafe.Pointer(uintptr(C.dword_5d4594_1321216))))
	C.dword_5d4594_1321216 = 0
	return result
}

//export sub_4BFB70
func sub_4BFB70(v C.int) { interactionKeyState(uint32(v)) }

//export sub_4BFBB0
func sub_4BFBB0(v *C.uint32_t) { interactionKeyUpdate(uintptr(unsafe.Pointer(v))) }

//export sub_4BFBF0
func sub_4BFBF0() C.int { return C.int(interactionKeyShow()) }

//export sub_4BFC70
func sub_4BFC70() C.int { return C.int(interactionKeyHide()) }

//export sub_4BFC90
func sub_4BFC90() C.int { return C.int(interactionKeyOpen()) }

//export sub_4BFCD0
func sub_4BFCD0(w, code C.int, a *C.int, b C.int) C.int {
	return C.int(gui.EventRespInt(interactionKeyEvent((*gui.Window)(unsafe.Pointer(uintptr(uint32(w)))), gui.AsWindowEvent(int(code), uintptr(unsafe.Pointer(a)), uintptr(uint32(b))))))
}

//export sub_4BFD10
func sub_4BFD10() { interactionKeyDestroy() }

//export sub_4BFD30
func sub_4BFD30() C.int { return C.int(C.dword_5d4594_1319056) }

//export sub_4C3390
func sub_4C3390() C.int { return C.int(interactionVoteOpen()) }

//export sub_4C3410
func sub_4C3410(w *C.int) C.int {
	return C.int(interactionIconDraw((*gui.Window)(unsafe.Pointer(w)), nil))
}

//export sub_4C3460
func sub_4C3460(v C.int) C.int { return C.int(interactionVoteHide(int(v))) }

//export sub_4C34A0
func sub_4C34A0() C.int { return C.int(interactionVoteDestroy()) }
