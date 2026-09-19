package legacy

/*
#include "defs.h"
extern uint32_t dword_5d4594_1064856, dword_5d4594_1064860;
extern uint32_t dword_5d4594_1064864, dword_5d4594_1064868;
extern int nox_win_width, nox_win_height;
void nox_xxx_consoleEsc_49B7A0();
*/
import "C"

import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func interactionChatRoot() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1064856)))
}
func interactionChatEdit() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1064860)))
}
func interactionChatData() unsafe.Pointer { return unsafe.Pointer(uintptr(C.dword_5d4594_1064864)) }
func interactionChatStart(team uint32) {
	if noxflags.HasGame(2048) || C.dword_5d4594_1064868 != 0 {
		return
	}
	data := interactionChatData()
	*(*uint16)(data) = 0
	*(*uint16)(unsafe.Add(data, 1052)) = 0
	interactionChatRoot().ShowModal()
	interactionChatRoot().StackPush()
	GetClient().Cli().GUI.Focus(interactionChatEdit())
	C.dword_5d4594_1064868 = 1
	*memmap.PtrUint32(0x5D4594, 1064872) = team
}
func interactionSay(text *uint16, team int) uint32 {
	code := ClientPlayerNetCode()
	dr := GetClient().Cli().Objs.ByNetCodeDynamic(code)
	units := gameplayTextUnits(text)
	original := len(units)
	for len(units) != 0 && units[0] == ' ' {
		units = units[1:]
	}
	if len(units) == 0 {
		return uint32(original)
	}
	x, y := uint16(0xffff), uint16(0xffff)
	if dr != nil {
		x, y = uint16(dr.PosVec.X), uint16(dr.PosVec.Y)
	}
	var flags byte
	if team != 0 {
		flags = 1
	}
	// The shared formatter already preserves raw UTF-16 and byte-count wrapping.
	msg := gameplayTextMessage(units, flags, uint16(code), x, y, 0)
	return uint32(bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg)))
}
func interactionChatDraw(w *gui.Window, d *gui.WindowData) int {
	root, edit := interactionChatRoot(), interactionChatEdit()
	root.ShowModal()
	GetClient().Cli().GUI.Focus(edit)
	r := GetClient().R2()
	face := r.GetFonts().AsFont(nil)
	data := interactionChatData()
	width := int32(r.GetStringSizeWrapped(face, alloc.GoString16((*uint16)(data)), 0).X) + int32(r.GetStringSizeWrapped(face, alloc.GoString16((*uint16)(unsafe.Add(data, 512))), 0).X) + 10
	if width < 100 {
		width = 100
	} else if width > 320 {
		width = 320
	}
	root.SetPos(image.Pt(int((int32(C.nox_win_width)-width)/2), root.Off.Y))
	uiWindowResize(w, int(width), 20)
	return uiEntryDraw(w, d, false)
}
func interactionChatClose() int {
	root := interactionChatRoot()
	if uiWindowHidden(root) != 0 {
		return 0
	}
	edit := interactionChatEdit()
	g := GetClient().Cli().GUI
	if g.Focused() == edit {
		g.Focus(nil)
	}
	root.StackPop()
	root.Hide()
	root.Flags &^= 8
	edit.Flags &^= 8
	g.ValYYY = 1
	C.dword_5d4594_1064868 = 0
	return 1
}
func interactionChatOpen() *gui.Window {
	x, y := int32(C.nox_win_width)/2, 2*int32(C.nox_win_height)/3
	*memmap.PtrUint32(0x5D4594, 1064876) = uint32(x)
	*memmap.PtrUint32(0x5D4594, 1064880) = uint32(y)
	root := Nox_new_window_from_file("GuiChat.wnd", interactionChatEvent)
	C.dword_5d4594_1064856 = C.uint32_t(uintptr(root.C()))
	if root == nil {
		return nil
	}
	root.SetPos(image.Pt(int(x), int(y)))
	edit := root.ChildByID(9201)
	C.dword_5d4594_1064860 = C.uint32_t(uintptr(edit.C()))
	if edit == nil {
		return nil
	}
	edit.SetDraw(interactionChatDraw)
	edit.SetFunc93(interactionChatKey)
	C.dword_5d4594_1064864 = C.uint32_t(uintptr(edit.WidgetData))
	return root
}
func interactionChatKey(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	a, b := ev.EventArgsC()
	if ev.EventCode() != 21 || a != 1 {
		return uiEntryInput(w, ev)
	}
	if b == 2 {
		C.nox_xxx_consoleEsc_49B7A0()
	}
	return gui.RawEventResp(1)
}
func interactionChatEvent(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
	if ev.EventCode() == 16415 {
		data := interactionChatData()
		if *(*uint16)(unsafe.Add(data, 1052)) != 0 {
			interactionSay((*uint16)(data), int(memmap.Uint32(0x5D4594, 1064872)))
		}
		interactionChatClose()
	}
	return nil
}
func interactionChatDestroy() int {
	result := 0
	if root := interactionChatRoot(); root != nil {
		root.Destroy()
		C.dword_5d4594_1064856 = 0
	}
	C.dword_5d4594_1064860 = 0
	C.dword_5d4594_1064864 = 0
	C.dword_5d4594_1064868 = 0
	*memmap.PtrUint32(0x5D4594, 1064872) = 0
	return result
}

//export sub_469FA0
func sub_469FA0() C.int { return C.int(memmap.Uint32(0x5D4594, 1064848)) }

//export nox_client_chatStart_46A430
func nox_client_chatStart_46A430(v C.int) { interactionChatStart(uint32(v)) }

//export sub_46A4A0
func sub_46A4A0() C.int { return C.int(C.dword_5d4594_1064868) }

//export nox_xxx_cmdSayDo_46A4B0
func nox_xxx_cmdSayDo_46A4B0(text *C.ushort, v C.int) C.size_t {
	return C.size_t(interactionSay((*uint16)(unsafe.Pointer(text)), int(v)))
}

//export sub_46A5D0
func sub_46A5D0(w *C.uint32_t, d C.int) C.int {
	return C.int(interactionChatDraw((*gui.Window)(unsafe.Pointer(w)), (*gui.WindowData)(unsafe.Pointer(uintptr(uint32(d))))))
}

//export sub_46A6A0
func sub_46A6A0() C.int { return C.int(interactionChatClose()) }

//export sub_46A730
func sub_46A730() *C.uint32_t { return (*C.uint32_t)(interactionChatOpen().C()) }

//export sub_46A7E0
func sub_46A7E0(w *C.uint32_t, code, a, b C.int) C.int {
	return C.int(gui.EventRespInt(interactionChatKey((*gui.Window)(unsafe.Pointer(w)), &gui.RawEvent{Event: int(code), Arg1: uintptr(uint32(a)), Arg2: uintptr(uint32(b))})))
}

//export sub_46A820
func sub_46A820(w, code, a, b C.int) C.int {
	return C.int(gui.EventRespInt(interactionChatEvent((*gui.Window)(unsafe.Pointer(uintptr(uint32(w)))), &gui.RawEvent{Event: int(code), Arg1: uintptr(uint32(a)), Arg2: uintptr(uint32(b))})))
}

//export sub_46A860
func sub_46A860() C.int { return C.int(interactionChatDestroy()) }
