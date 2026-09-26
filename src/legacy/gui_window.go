package legacy

/*
#include <stdlib.h>
#include "client__gui__window.h"

*/
import "C"
import (
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

func Nox_xxx_wnd_46ABB0(p *gui.Window, v int) int {
	return uiWindowEnable(p, v)
}
