//go:build porttest

package legacy

/*
#include "defs.h"

#include "GAME3_1.h"
#include "client__gui__guiinput.h"
#include "client__shell__inputcfg__inputcfg.h"








*/
import "C"
import "unsafe"
import "github.com/opennox/opennox/v1/client/gui"

// PortTestBindingWords owns live globals independently of the backing memory blob.
func PortTestBindingWords() (map[int]*uint32, func()) {
	w := map[int]*uint32{
		1321224: (*uint32)(unsafe.Pointer(&dword_5d4594_1321224)),
		1321228: (*uint32)(unsafe.Pointer(&dword_5d4594_1321228)),
		1321232: (*uint32)(unsafe.Pointer(&dword_5d4594_1321232)),
		1321236: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1321236)),
		1321240: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1321240)),
		1321244: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1321244)),
		1321248: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1321248)),
		1321252: (*uint32)(unsafe.Pointer(&dword_5d4594_1321252)),
		1522604: (*uint32)(unsafe.Pointer(&dword_5d4594_1522604)),
		1522612: (*uint32)(unsafe.Pointer(&dword_5d4594_1522612)),
		1522616: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1522616)),
		1522620: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1522620)),
		1522624: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1522624)),
		1522628: (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1522628)),
		1522632: (*uint32)(unsafe.Pointer(&dword_5d4594_1522632)),
	}
	old := make(map[int]uint32)
	for n, p := range w {
		old[n] = *p
		*p = 0
	}
	return w, func() {
		for n, p := range w {
			*p = old[n]
		}
	}
}

// PortTestBindingAssign invokes the production Go keyboard/mouse assignment owners.
func PortTestBindingAssign(menu, mouse bool, key uint32) int {
	return bindingEditor(menu).assign(key, mouse)
}

func PortTestBindingModal(menu bool, win uint32, event int, key uint32, state int) int {
	return gui.EventRespInt(bindingEditor(menu).modalEvent((*gui.Window)(unsafe.Pointer(uintptr(win))), &gui.RawEvent{Event: event, Arg1: uintptr(key), Arg2: uintptr(uint32(state))}))
}

func PortTestBindingApply(menu bool) { bindingEditor(menu).apply() }

func PortTestBindingDimensions() [2]*int32 {
	return [2]*int32{(*int32)(unsafe.Pointer(&nox_win_width)), (*int32)(unsafe.Pointer(&nox_win_height))}
}
func PortTestBindingConstruct(menu bool) int { return bindingEditor(menu).construct() }

func PortTestBindingDestroy() int { return bindingDestroy() }

func PortTestBindingInvoke(op string, a [4]uint32) uint32 {
	w := (*gui.Window)(unsafe.Pointer(uintptr(a[0])))
	ev := &gui.RawEvent{Event: int(a[1]), Arg1: uintptr(a[2]), Arg2: uintptr(a[3])}
	switch op {
	case "sub_4C3500":
		return bindingYesNo()
	case "sub_4C35B0":
		return uint32(bindingClose(int(a[0])))
	case "sub_4C4260":
		bindingShow()
		return 0
	case "sub_4C4280":
		return uint32(bindingVisible())
	case "sub_4CBB70":
		return uint32(bindingMenuBack())
	case "sub_4CBBB0":
		return uint32(bindingMenuDone())
	case "sub_4C3A60", "sub_4CC140":
		return uint32(gui.EventRespInt(bindingFilter(w, ev)))
	case "sub_4C3CD0":
		return uint32(gui.EventRespInt(bindingInGame.route(w, ev)))
	case "sub_4CBF60":
		return uint32(gui.EventRespInt(bindingMenu.route(w, ev)))
	default:
		panic(op)
	}
}

func PortTestBindingAnimationWord() *uint32 {
	return (*uint32)(unsafe.Pointer(&legacyGlobals.nox_wnd_xxx_1522608))
}
