package legacy

/*
#include <stdlib.h>
#include "client__gui__window.h"

*/
import "C"
import (
	"image"
	"log/slog"
	"runtime/debug"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

type nox_window = C.nox_window

func AsWindowP(win unsafe.Pointer) *gui.Window {
	w := (*gui.Window)(win)
	if false && cgoSafe && w.ID() == DeadWord {
		slog.Error("memory corruption detected")
		debug.PrintStack()
		C.abort()
	}
	return w
}

func asWindow(win *nox_window) *gui.Window {
	return AsWindowP(unsafe.Pointer(win))
}

func get_dword_5d4594_3799468() int {
	return GetClient().Cli().GUI.ValXXX
}

//export nox_xxx_wndGetID_46B0A0
func nox_xxx_wndGetID_46B0A0(win *nox_window) int {
	if win == nil {
		return -2
	}
	return int(asWindow(win).ID())
}

//export nox_gui_winSetFunc96_46B070
func nox_gui_winSetFunc96_46B070(win *nox_window, fnc unsafe.Pointer) {
	asWindow(win).SetTooltipFunc(fnc)
}

//export nox_window_call_field_94_fnc
func nox_window_call_field_94_fnc(p *nox_window, a2, a3, a4 int, file *C.char, line int) int {
	if p == nil {
		return 0
	}
	if guiDebug {
		guiLog.Printf("nox_window_call_field_94(%p, %x, %x, %x): %s:%d", p, a2, a3, a4, GoString(file), line)
	}
	r := asWindow(p).Func94(gui.AsWindowEvent(a2, uintptr(a3), uintptr(a4)))
	if r == nil {
		return 0
	}
	return int(r.EventRespC())
}

//export nox_window_call_field_93
func nox_window_call_field_93(p *nox_window, a2, a3, a4 int) int {
	if p == nil {
		return 0
	}
	r := asWindow(p).Func93(gui.AsWindowEvent(a2, uintptr(a3), uintptr(a4)))
	if r == nil {
		return 0
	}
	return int(r.EventRespC())
}

//export nox_xxx_wndGetChildByID_46B0C0
func nox_xxx_wndGetChildByID_46B0C0(root *nox_window, id int) *nox_window {
	return (*nox_window)(asWindow(root).ChildByID(uint(id)).C())
}

func nox_xxx_windowDestroyMB_46C4E0(a1 *nox_window) int {
	win := asWindow(a1)
	if win == nil {
		return -2
	}
	win.Destroy()
	return 0
}

func nox_window_set_hidden(p *nox_window, hidden int) int {
	if p == nil {
		return -2
	}
	win := asWindow(p)
	if hidden != 0 {
		win.Hide()
	} else {
		win.Show()
	}
	return 0
}

//export nox_xxx_wndShowModalMB_46A8C0
func nox_xxx_wndShowModalMB_46A8C0(p *nox_window) int {
	return asWindow(p).ShowModal()
}

//export nox_window_setPos_46A9B0
func nox_window_setPos_46A9B0(p *nox_window, x, y int) int {
	win := asWindow(p)
	if win == nil {
		return -2
	}
	win.SetPos(image.Point{X: x, Y: y})
	return 0
}

//export wndIsShown_nox_xxx_wndIsShown_46ACC0
func wndIsShown_nox_xxx_wndIsShown_46ACC0(p *nox_window) int {
	if p == nil {
		return 1
	}
	win := asWindow(p)
	is := win.GetFlags().IsHidden()
	return bool2int(is)
}

//export nox_xxx_wnd_46C6E0
func nox_xxx_wnd_46C6E0(p *nox_window) int {
	return asWindow(p).StackPop()
}

func Nox_xxx_wnd_46ABB0(p *gui.Window, v int) int {
	return uiWindowEnable(p, v)
}
