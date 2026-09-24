package legacy

/*
#include "GAME3.h"
#include "client__shell__selcolor.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// These callbacks remain addressable by the existing C-layout animation slots.
//
//export nox_game_showSelClass_4A4840
func nox_game_showSelClass_4A4840() C.int { return C.int(characterShowClass()) }

//export nox_game_showSelColor_4A5D00
func nox_game_showSelColor_4A5D00() C.int { return C.int(characterShowColor()) }

//export sub_4A4970
func sub_4A4970() C.int { return C.int(characterClassStart()) }

//export sub_4A49A0
func sub_4A49A0() C.int { return C.int(characterClassDone()) }

//export sub_4A6890
func sub_4A6890() C.int { return C.int(characterColorStart()) }

//export sub_4A6C90
func sub_4A6C90() C.int { return C.int(characterColorDone()) }

func characterClassNextColor() { characterUI.classAnim.Func13Ptr = C.nox_game_showSelColor_4A5D00 }
func characterMenuEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	return gui.RawEventResp(Sub_4A18E0(w, ev.EventCode(), int(a), int(b)))
}
func characterShowClass() int {
	questProgressReset("*:*")
	Sub_4A1BE0(1)
	characterUI.classHost = unsafe.Pointer(Nox_xxx_getHostInfoPtr_431770())
	GetClient().GameAddStateCode(600)
	w := Nox_new_window_from_file("SelClass.wnd", characterClassEvent)
	characterUI.classRoot = w
	if w == nil {
		return 0
	}
	w.SetFunc93(characterMenuEvent)
	a := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -460, 0, 20, 0, -40)
	characterUI.classAnim = a
	if a == nil {
		return 0
	}
	a.StateID = 600
	a.Func12Ptr = C.sub_4A4970
	a.FncDoneOutPtr = C.sub_4A49A0
	for _, id := range []uint{601, 603, 602} {
		w.ChildByID(id).SetDraw(func(w *gui.Window, _ *gui.WindowData) int { return characterClassDraw(w) })
	}
	*memmap.PtrUint32(0x5D4594, 1307728) = uint32(uintptr(w.ChildByID(610).C()))
	*memmap.PtrUint32(0x5D4594, 1307740) = 0
	Sub_4A19F0("OptsBack.wnd:Back")
	quickbarClearSlots()
	return 1
}
func characterShowColor() int {
	GetClient().GameAddStateCode(700)
	characterUI.colorHost = unsafe.Pointer(Nox_xxx_getHostInfoPtr_431770())
	*(*byte)(unsafe.Add(characterUI.colorHost, 67)) = 0
	w := Nox_new_window_from_file("SelColor.wnd", characterColorEvent)
	characterUI.colorRoot = w
	if w == nil {
		return 0
	}
	w.SetFunc93(characterMenuEvent)
	a := Nox_gui_makeAnimation_43C5B0(w, 0, 0, 0, -440, 0, 20, 0, -40)
	characterUI.colorAnim = a
	if a == nil {
		return 0
	}
	a.StateID = 700
	a.Func12Ptr = C.sub_4A6890
	a.FncDoneOutPtr = C.sub_4A6C90
	characterSetup()
	draw := func(w *gui.Window, _ *gui.WindowData) int { return characterPaletteDraw(w) }
	for id := uint(720); id <= 729; id++ {
		w.ChildByID(id).SetDraw(draw)
	}
	for id := uint(761); id <= 792; id++ {
		w.ChildByID(id).SetDraw(draw)
	}
	if characterUI.defaults != 0 {
		optionsSend(characterUI.controls[14], 16414, uintptr(unsafe.Pointer(characterDefaultName())), 0)
	}
	characterUI.palette = w.ChildByID(760)
	characterUI.palette.SetFunc94(characterColorEvent)
	characterUI.palette.SetFunc93(characterPaletteOutside)
	characterUI.palette.SetParent(nil)
	Sub_4BFAD0()
	w.ChildByID(740).SetDraw(func(w *gui.Window, _ *gui.WindowData) int { return characterPreview(w) })
	return 1
}
