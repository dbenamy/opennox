//go:build porttest

package legacy

/*
#include "client__gui__window.h"
#include "GAME2_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// Exercises retained C entrypoints and native adapters for retired fixture-only bridges.
func PortTestWindowHelper(op int, w *gui.Window, a, b int, other unsafe.Pointer, flag int) (int, [2]uint32) {
	p := (*nox_window)(w.C())
	addr := C.int(uintptr(w.C()))
	x, y := C.uint(0x13572468), C.uint(0x89abcdef)
	ret := 0
	switch op {
	case 0:
		ret = int(C.nox_gui_getWindowOffs_46AA20(p, &x, &y))
	case 1:
		if flag != 0 {
			ret = int(C.nox_client_wndGetPosition_46AA60(p, &x, &x))
		} else {
			ret = int(C.nox_client_wndGetPosition_46AA60(p, &x, &y))
		}
	case 2:
		ret = int(portUIWindow_nox_window_get_size(p, (*int32)(unsafe.Pointer(&x)), (*int32)(unsafe.Pointer(&y))))
	case 3:
		if C.nox_xxx_wndPointInWnd_46AAB0((*C.uint)(w.C()), C.int(a), C.int(b)) {
			ret = 1
		}
	case 4:
		ret = int(portUIWindow_sub_46AB20((*uint32)(w.C()), int32(a), int32(b)))
	case 5:
		ret = int(portUIWindow_nox_xxx_wnd_46ABB0(p, int32(a)))
	case 6:
		ret = int(portUIWindow_nox_xxx_wnd_46AD60(int32(addr), int32(a)))
	case 7:
		ret = int(portUIWindow_nox_xxx_wndClearFlag_46AD80(int32(addr), int32(a)))
	case 8:
		ret = int(portUIWindow_nox_xxx_wndGetFlags_46ADA0(int32(addr)))
	case 9:
		ret = int(C.nox_window_is_child(p, (*nox_window)(other)))
	case 10:
		ret = int(portUIWindow_nox_xxx_wnd_46B280(int32(addr), int32(uintptr(other))))
	case 11:
		ret = int(portUIWindow_sub_46AE10(int32(addr), int32(a)))
	case 12:
		ret = int(portUIWindow_nox_xxx_wndSetOffsetMB_46AE40(int32(addr), int32(a), int32(b)))
	case 13:
		ret = int(portUIWindow_nox_xxx_wndSetIcon_46AE60(int32(addr), int32(uintptr(other))))
	case 14:
		ret = int(portUIWindow_nox_xxx_wndSetIconLit_46AEA0(int32(addr), int32(uintptr(other))))
	case 15:
		ret = int(portUIWindow_sub_46AEC0(int32(addr), int32(uintptr(other))))
	case 16:
		ret = int(C.sub_46AEE0(addr, C.int(uintptr(other))))
	case 17:
		ret = int(uintptr(unsafe.Pointer(C.sub_46AF00(w.C()))))
	case 18:
		ret = int(uintptr(uiWindowFont(w)))
	case 19:
		ret = w.CopyDrawData((*gui.WindowData)(other))
	case 20:
		ret = int(uintptr(uiWindowChildAt(w, a, b).C()))
	case 21:
		ret = uiWindowHidden(w)
	case 22:
		portUIWindow_sub_46ACE0((*uint32)(w.C()), int32(a), int32(b), int32(flag))
	case 23:
		portUIWindow_sub_46AD20((*uint32)(w.C()), int32(a), int32(b), int32(flag))
	case 24:
		ret = int(portUIWindow_nox_xxx_wndRetNULL_46A8A0())
	case 25:
		ret = 0
	default:
		panic("window helper case")
	}
	return ret, [2]uint32{uint32(x), uint32(y)}
}

// Native adapters preserve the retired export boundaries, including output write order.
func portUIWindow_nox_window_get_size(w *nox_window, x, y *int32) int32 {
	if w == nil {
		*x = 0
		*y = 0
		return -2
	}
	*x = int32(asWindow(w).SizeVal.X)
	*y = int32(asWindow(w).SizeVal.Y)
	return 0
}

func portUIWindow_nox_xxx_wndClearFlag_46AD80(ptr, mask int32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	old := w.Flags
	w.Flags &^= gui.StatusFlags(mask)
	return int32(old)
}

func portUIWindow_nox_xxx_wndGetFlags_46ADA0(ptr int32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	return int32(w.Flags)
}

func portUIWindow_nox_xxx_wndRetNULL_46A8A0() int32 { return 0 }

func portUIWindow_nox_xxx_wndSetIconLit_46AEA0(ptr, img int32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	w.DrawData().HlImageHnd = noxrender.ImageHandle(unsafe.Pointer(uintptr(uint32(img))))
	return 0
}

func portUIWindow_nox_xxx_wndSetIcon_46AE60(ptr, img int32) int32 {
	return int32(uiWindowSetBackgroundImage(uint32(ptr), uint32(img)))
}

func portUIWindow_nox_xxx_wndSetOffsetMB_46AE40(ptr, x, y int32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w != nil {
		w.DrawData().ImgPtVal = image.Pt(int(x), int(y))
	}
	return ptr
}

func portUIWindow_nox_xxx_wnd_46ABB0(w *nox_window, enabled int32) int32 {
	return int32(uiWindowEnable(asWindow(w), int(enabled)))
}

func portUIWindow_nox_xxx_wnd_46AD60(ptr, mask int32) int32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(ptr))))
	if w == nil {
		return -2
	}
	old := w.Flags
	w.Flags |= gui.StatusFlags(mask)
	return int32(old)
}

func portUIWindow_nox_xxx_wnd_46B280(ptr, owner int32) int32 {
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

func portUIWindow_sub_46AB20(w *uint32, width, height int32) int32 {
	return int32(uiWindowResize((*gui.Window)(unsafe.Pointer(w)), int(width), int(height)))
}

func portUIWindow_sub_46ACE0(w *uint32, first, last, hidden int32) {
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

func portUIWindow_sub_46AD20(w *uint32, first, last, enabled int32) {
	win := (*gui.Window)(unsafe.Pointer(w))
	for i := int(first); i <= int(last); i++ {
		uiWindowEnable(win.ChildByID(uint(i)), int(enabled))
		if i == int(last) {
			break
		}
	}
}

func portUIWindow_sub_46AE10(ptr, on int32) int32 {
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

func portUIWindow_sub_46AEC0(ptr, img int32) int32 {
	return int32(uiWindowSetSelectedImage(uint32(ptr), uint32(img)))
}
