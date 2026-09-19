package legacy

/*
#include "GAME1_2.h"
#include "client__shell__noxworld.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"unsafe"
)

func browserColumns() []*gui.Window {
	return []*gui.Window{browserWindow(uint32(browserUI.playersColumn)), browserWindow(uint32(browserUI.modeColumn)), browserWindow(uint32(browserUI.mapColumn)), browserWindow(uint32(browserUI.pingColumn)), browserWindow(uint32(browserUI.statusColumn))}
}
func browserColumnsSend(code int, a, b uintptr) int {
	ret := 0
	for _, w := range browserColumns() {
		ret = optionsSend(w, code, a, b)
	}
	return ret
}
func browserColumnsInit() int {
	main := browserWindow(uint32(browserUI.gameList))
	main.SetFunc94(browserListEvent)
	main.SetFunc93(browserListInput)
	for _, w := range browserColumns() {
		w.SetParent(main)
		w.SetFunc93(browserListInput)
	}
	d := uiListData(asWindow(browserUI.detailList))
	(*gui.Window)(d.Field_7).SetID(10035)
	(*gui.Window)(d.Field_8).SetID(10036)
	(*gui.Window)(d.Field_9).SetID(10032)
	d = uiListData(main)
	browserColumnsSend(16408, uintptr(d.Field_7), 0)
	browserColumnsSend(16409, uintptr(d.Field_8), 0)
	return browserColumnsSend(16410, uintptr(d.Field_9), 0)
}
func browserListInput(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	code := ev.EventCode()
	if code == 19 || code == 20 {
		p := uiListData(w).Field_7
		if code == 20 {
			p = uiListData(w).Field_8
		}
		optionsSend(browserWindow(uint32(browserUI.gameList)), 16391, uintptr(p), 0)
		browserColumnsSend(16391, uintptr(p), 0)
		return nil
	}
	return uiListSingleInput(w, ev)
}
func browserListEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	code := ev.EventCode()
	a, b := ev.EventArgsC()
	switch code {
	case 23:
		return gui.RawEventResp(1)
	case 16400:
		id := browserWindow(uint32(a)).ID()
		if id >= 10038 && id <= 10042 {
			browserColumnsSend(16403, b, 0)
			if uint32(b) < uint32(browserUI.resultCount) {
				p := browserListAt(int32(b))
				browserUI.selected = p
				pos := GetClient().GetMousePos()
				point := [2]uint32{uint32(pos.X), uint32(pos.Y)}
				browserInfoPopup(&point, p)
			}
		}
	case 16403, 16412:
		browserColumnsSend(code, a, 0)
	case 16398, 16399:
		browserColumnsSend(code, a, b)
	case 16384:
		main := browserWindow(uint32(browserUI.gameList))
		if browserWindow(uint32(a)) == main.ChildByID(10043) || browserWindow(uint32(a)) == main.ChildByID(10044) {
			browserColumnsSend(code, a, 0)
		}
	case 16393:
		d := uiListData(w)
		uiListEvent(w, ev)
		browserColumnsSend(16412, uintptr(uiListIndex(d)), 0)
	}
	return uiListEvent(w, ev)
}
func browserInfoPosition(x, y int32, out *[2]uint32) {
	px, py := x-65, y-20
	if px+130 > 600 {
		px = 470
	}
	if py+120 > 451 {
		py = 331
	}
	low := uint32(27)
	if browserUI.listMode == 1 {
		low = 55
	}
	if uint32(py) < low {
		py = int32(low)
	}
	if px < 216 {
		px = 216
	}
	out[0], out[1] = uint32(px), uint32(py)
}
func browserInfoPopup(point *[2]uint32, record unsafe.Pointer) {
	browserInfoPosition(int32(point[0]), int32(point[1]), point)
	w := browserWindow(uint32(uintptr(browserUI.detailPanel)))
	w.SetParent(nil)
	w.ShowModal()
	w.StackPush()
	w.SetPos(image.Pt(int(int32(point[0])), int(int32(point[1]))))
	browserDetails(record)
	browserUI.hasSelection = 1
	*memmap.PtrUint16(0x5D4594, 814604) = *(*uint16)(unsafe.Add(record, 109))
	if noxflags.HasEngine(noxflags.EngineNoRendering) {
		uiWindowEnable(browserWindow(memmap.Uint32(0x5D4594, 815008)), 0)
	}
}

func sub_438480() C.int { return C.int(browserColumnsInit()) }

func sub_438EF0(w *C.uint32_t, code C.int, a C.uint, b C.int) C.int {
	return C.int(gui.EventRespInt(browserListInput((*gui.Window)(unsafe.Pointer(w)), &gui.RawEvent{Event: int(code), Arg1: uintptr(a), Arg2: uintptr(uint32(b))})))
}

func sub_439050(w C.int, code C.uint, a *C.int, b C.uint) C.int {
	return C.int(gui.EventRespInt(browserListEvent(browserWindow(uint32(w)), &gui.RawEvent{Event: int(code), Arg1: uintptr(unsafe.Pointer(a)), Arg2: uintptr(b)})))
}

func sub_439450(x, y C.int, out *C.uint32_t) *C.uint32_t {
	browserInfoPosition(int32(x), int32(y), (*[2]uint32)(unsafe.Pointer(out)))
	return out
}

func nox_client_gui_serverInfoBlockCheckExp_439370(point *C.int2, record C.int) {
	browserInfoPopup((*[2]uint32)(unsafe.Pointer(point)), unsafe.Pointer(uintptr(uint32(record))))
}
