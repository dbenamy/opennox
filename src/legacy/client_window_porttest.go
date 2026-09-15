//go:build porttest

package legacy

/*
#include "client__gui__window.h"
#include "GAME2_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

// Calls the production ABI; contains no reference implementation.
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
		ret = int(C.nox_window_get_size(p, (*C.int)(unsafe.Pointer(&x)), (*C.int)(unsafe.Pointer(&y))))
	case 3:
		if C.nox_xxx_wndPointInWnd_46AAB0((*C.uint)(w.C()), C.int(a), C.int(b)) {
			ret = 1
		}
	case 4:
		ret = int(C.sub_46AB20((*C.uint)(w.C()), C.int(a), C.int(b)))
	case 5:
		ret = int(C.nox_xxx_wnd_46ABB0(p, C.int(a)))
	case 6:
		ret = int(C.nox_xxx_wnd_46AD60(addr, C.int(a)))
	case 7:
		ret = int(C.nox_xxx_wndClearFlag_46AD80(addr, C.int(a)))
	case 8:
		ret = int(C.nox_xxx_wndGetFlags_46ADA0(addr))
	case 9:
		ret = int(C.nox_window_is_child(p, (*nox_window)(other)))
	case 10:
		ret = int(C.nox_xxx_wnd_46B280(addr, C.int(uintptr(other))))
	case 11:
		ret = int(C.sub_46AE10(addr, C.int(a)))
	case 12:
		ret = int(C.nox_xxx_wndSetOffsetMB_46AE40(addr, C.int(a), C.int(b)))
	case 13:
		ret = int(C.nox_xxx_wndSetIcon_46AE60(addr, C.int(uintptr(other))))
	case 14:
		ret = int(C.nox_xxx_wndSetIconLit_46AEA0(addr, C.int(uintptr(other))))
	case 15:
		ret = int(C.sub_46AEC0(addr, C.int(uintptr(other))))
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
		C.sub_46ACE0((*C.uint)(w.C()), C.int(a), C.int(b), C.int(flag))
	case 23:
		C.sub_46AD20((*C.uint)(w.C()), C.int(a), C.int(b), C.int(flag))
	case 24:
		ret = int(C.nox_xxx_wndRetNULL_46A8A0())
	case 25:
		ret = 0
	default:
		panic("window helper case")
	}
	return ret, [2]uint32{uint32(x), uint32(y)}
}
