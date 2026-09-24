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

//export nox_window_get_size
func nox_window_get_size(w *nox_window, x, y *C.int) C.int {
	if w == nil {
		*x = 0
		*y = 0
		return -2
	}
	*x = C.int(asWindow(w).SizeVal.X)
	*y = C.int(asWindow(w).SizeVal.Y)
	return 0
}

//export nox_xxx_wndPointInWnd_46AAB0
func nox_xxx_wndPointInWnd_46AAB0(w *C.uint, x, y C.int) C.bool {
	p, size := image.Point{}, image.Point{}
	if w != nil {
		win := (*gui.Window)(unsafe.Pointer(w))
		p = uiWindowPosition(win)
		size = win.SizeVal
	}
	return C.bool(int(x) >= p.X && int(x) <= p.X+size.X && int(y) >= p.Y && int(y) <= p.Y+size.Y)
}

//export sub_46AB20
func sub_46AB20(w *C.uint, width, height C.int) C.int {
	return C.int(uiWindowResize((*gui.Window)(unsafe.Pointer(w)), int(width), int(height)))
}

//export nox_xxx_wnd_46ABB0
func nox_xxx_wnd_46ABB0(w *nox_window, enabled C.int) C.int {
	return C.int(uiWindowEnable(asWindow(w), int(enabled)))
}

//export nox_xxx_wnd_46AD60
func nox_xxx_wnd_46AD60(ptr, mask C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	old := w.Flags
	w.Flags |= gui.StatusFlags(mask)
	return C.int(old)
}

//export nox_xxx_wndClearFlag_46AD80
func nox_xxx_wndClearFlag_46AD80(ptr, mask C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	old := w.Flags
	w.Flags &^= gui.StatusFlags(mask)
	return C.int(old)
}

//export nox_xxx_wndGetFlags_46ADA0
func nox_xxx_wndGetFlags_46ADA0(ptr C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	return C.int(w.Flags)
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

//export nox_xxx_wnd_46B280
func nox_xxx_wnd_46B280(ptr, owner C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	d := (*gui.Window)(unsafe.Pointer(uintptr(uint32(owner))))
	if d == nil {
		d = w
	}
	w.DrawData().Window = d
	return 0
}

//export sub_46ACE0
func sub_46ACE0(w *C.uint, first, last, hidden C.int) {
	win := (*gui.Window)(unsafe.Pointer(w))
	for i := int(first); i <= int(last); i++ {
		c := win.ChildByID(uint(i))
		if c != nil {
			c.SetHidden(hidden != 0)
		}
		if i == int(last) {
			break
		}
	}
}

//export sub_46AD20
func sub_46AD20(w *C.uint, first, last, enabled C.int) {
	win := (*gui.Window)(unsafe.Pointer(w))
	for i := int(first); i <= int(last); i++ {
		uiWindowEnable(win.ChildByID(uint(i)), int(enabled))
		if i == int(last) {
			break
		}
	}
}

//export nox_xxx_wndRetNULL_46A8A0
func nox_xxx_wndRetNULL_46A8A0() C.int { return 0 }

//export sub_46AE10
func sub_46AE10(ptr, on C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w != nil {
		if on != 0 {
			w.DrawData().Field0 |= 2
		} else {
			w.DrawData().Field0 &^= 2
		}
	}
	return ptr
}

//export nox_xxx_wndSetOffsetMB_46AE40
func nox_xxx_wndSetOffsetMB_46AE40(ptr, x, y C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w != nil {
		w.DrawData().ImgPtVal = image.Pt(int(x), int(y))
	}
	return ptr
}

func uiWindowSetBackgroundImage(ptr, img uint32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(ptr)))
	if w == nil {
		return -2
	}
	w.DrawData().BgImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(img)))
	return 0
}

//export nox_xxx_wndSetIcon_46AE60
func nox_xxx_wndSetIcon_46AE60(ptr, img C.int) C.int {
	return C.int(uiWindowSetBackgroundImage(uint32(ptr), uint32(img)))
}

//export nox_xxx_wndSetIconLit_46AEA0
func nox_xxx_wndSetIconLit_46AEA0(ptr, img C.int) C.int {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	w.DrawData().HlImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(uint32(img))))
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

//export sub_46AEC0
func sub_46AEC0(ptr, img C.int) C.int {
	return C.int(uiWindowSetSelectedImage(uint32(ptr), uint32(img)))
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
