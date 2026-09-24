package legacy

/*
#include "client__gui__window.h"
#include "GAME1_2.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "client__shell__inputcfg__inputcfg.h"


*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/strman"

	"github.com/opennox/opennox/v1/client/gui"
)

var (
	Nox_video_setMenuOptions func(win *gui.Window)
	Nox_gui_menu_proc_ext    func(id int) int
	Sub_4A19F0               func(name strman.ID)
	Sub_4AAA10               func() int
	Sub_4C3A90               func(a1, a2 int, a3 unsafe.Pointer, a4 int) int
	Sub_4CBE70               func(a1, a2 int, a3 unsafe.Pointer, a4 int) int
	Sub_4A1A40               func(a1 int)
)

//export sub_4AAA10
func sub_4AAA10() int { return Sub_4AAA10() }

func Sub_4CBD30() {
	bindingMenu.apply()
}

func Sub_430AA0(v int) {
	sub_430AA0(C.int(v))
}

func Sub_4C35B0(v int) {
	bindingClose(v)
}

func Get_nox_wnd_xxx_1309740() *gui.Anim {
	return asGUIAnim(legacyGlobals.nox_wnd_xxx_1309740)
}

func Get_dword_5d4594_1309720() *gui.Window {
	return AsWindowP(unsafe.Pointer(uintptr(legacyGlobals.dword_5d4594_1309720)))
}
