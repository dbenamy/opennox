package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// Raw ABI helpers retain access to owned windows whose destruction flag is set.
// The higher-level GUI getters intentionally hide some of that state.
func uiWindowPosition(w *gui.Window) image.Point {
	p := w.Off
	for w = w.Parent(); w != nil; w = w.Parent() {
		p = p.Add(w.Off)
	}
	return p
}
func uiWindowResize(w *gui.Window, width, height int) int {
	if w == nil {
		return -2
	}
	w.SizeVal = image.Pt(width, height)
	w.EndPos = w.Off.Add(w.SizeVal)
	w.Func94(gui.AsWindowEvent(16388, uintptr(width), uintptr(height)))
	return 0
}
func uiWindowEnable(w *gui.Window, enabled int) int {
	if w == nil {
		return -2
	}
	if enabled != 0 {
		w.Flags |= 8
	} else {
		w.Flags &^= 8
	}
	for c := w.Field100Ptr; c != nil; c = c.Prev() {
		uiWindowEnable(c, enabled)
	}
	return 0
}
func uiWindowText(w *gui.Window) unsafe.Pointer {
	if w == nil {
		return nil
	}
	style := w.DrawData().Style
	code := 0
	if style&0x800 != 0 {
		code = 16386
	} else if style&0x80 != 0 {
		code = 16413
	} else {
		return nil
	}
	return unsafe.Pointer(uintptr(gui.EventRespInt(w.Func94(gui.AsWindowEvent(code, 0, 0)))))
}
func uiWindowChildAt(w *gui.Window, x, y int) *gui.Window {
	if w == nil {
		return nil
	}
	for {
		found := false
		for c := w.Field100Ptr; c != nil; c = c.Prev() {
			p := uiWindowPosition(c)
			if x >= p.X && uint32(x) <= uint32(p.X)+uint32(c.SizeVal.X) && y >= p.Y && uint32(y) <= uint32(p.Y)+uint32(c.SizeVal.Y) && c.Flags&16 == 0 {
				w = c
				found = true
				break
			}
		}
		if !found {
			return w
		}
	}
}

//export nox_gui_getWindowOffs_46AA20
func nox_gui_getWindowOffs_46AA20(w *nox_window, x, y *C.uint) C.int {
	if w == nil {
		*x = 0
		*y = 0
		return -2
	}
	*x = C.uint(asWindow(w).Off.X)
	*y = C.uint(asWindow(w).Off.Y)
	return 0
}

//export nox_client_wndGetPosition_46AA60
func nox_client_wndGetPosition_46AA60(w *nox_window, x, y *C.uint) C.int {
	if w == nil {
		return -2
	}
	*x = C.uint(asWindow(w).Off.X)
	*y = C.uint(asWindow(w).Off.Y)
	for i := asWindow(w).Parent(); i != nil; i = i.Parent() {
		*x += C.uint(i.Off.X)
		*y += C.uint(i.Off.Y)
	}
	return 0
}

//export nox_xxx_wndPointInWnd_46AAB0
func nox_xxx_wndPointInWnd_46AAB0(w *C.uint, x, y C.int) C.bool {
	return C.bool(uiWindowPointIn((*gui.Window)(unsafe.Pointer(w)), int32(x), int32(y)))
}

func uiWindowPointIn(w *gui.Window, x, y int32) bool {
	p, size := image.Point{}, image.Point{}
	if w != nil {
		win := w
		p = uiWindowPosition(win)
		size = win.SizeVal
	}
	return int(x) >= p.X && int(x) <= p.X+size.X && int(y) >= p.Y && int(y) <= p.Y+size.Y
}

//export nox_window_is_child
func nox_window_is_child(parent, child *nox_window) C.int {
	if parent == nil || child == nil {
		return 0
	}
	for c := asWindow(child).Parent(); c != nil; c = c.Parent() {
		if c == asWindow(parent) {
			return 1
		}
	}
	return 0
}

func uiWindowSetBackgroundImage(ptr, img uint32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(ptr)))
	if w == nil {
		return -2
	}
	w.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(img)))
	return 0
}

func uiWindowSetSelectedImage(ptr, img uint32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(ptr)))
	if w == nil {
		return -2
	}
	w.DrawData().SelImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(img)))
	return 0
}

//export sub_46AEE0
func sub_46AEE0(ptr, text C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w != nil {
		w.Func94(gui.AsWindowEvent(16385, uintptr(uint32(text)), 0))
	}
	return 0
}

//export sub_46AF00
func sub_46AF00(ptr unsafe.Pointer) *wchar2_t { return (*wchar2_t)(uiWindowText((*gui.Window)(ptr))) }

func uiWindowFont(w *gui.Window) unsafe.Pointer {
	ptr := w.C()
	if ptr == nil {
		return nil
	}
	return (*gui.Window)(ptr).DrawData().FontPtr
}

func uiWindowHidden(w *gui.Window) int {
	if w == nil {
		return 1
	}
	for win := w; win != nil; win = win.Parent() {
		if win.Flags&16 != 0 {
			return 1
		}
	}
	return 0
}
