//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__gui__guicon.h"
#include <stdint.h>
static int port_test_interaction_console(wchar2_t* format, wchar2_t* wide, char* narrow, int number) {
	return nox_gui_console_Printf_450C00(NOX_CONSOLE_RED, format, wide, narrow, number);
}
*/
import "C"
import "unsafe"

func PortTestClientInteractionCall(op string, args ...uintptr) uint64 {
	var a [7]uintptr
	copy(a[:], args)
	id := -1
	switch op {
	case "nox_xxx_clientTalk_42E7B0":
		id = 0
	case "nox_xxx_clientCollideOrUse_42E810":
		id = 1
	case "nox_xxx_clientTrade_42E850":
		id = 2
	case "sub_430AA0":
		id = 3
	case "nox_client_mousePriKey_430AF0":
		id = 4
	case "nox_xxx_cursor_430B00":
		id = 5
	case "nox_client_setMousePos_430B10":
		id = 6
	case "nox_xxx_initTime_435570":
		id = 7
	case "nox_client_drawable_testBuff_4356C0":
		id = 8
	case "sub_435700":
		id = 9
	case "nox_xxx_cliToggleObsWindow_4357A0":
		id = 10
	case "sub_435F60":
		id = 11
	case "sub_436550":
		id = 12
	case "sub_437100":
		id = 13
	case "nox_xxx_playerAnimCheck_4372B0":
		id = 14
	case "nox_xxx_clientIsObserver_4372E0":
		id = 15
	case "sub_445450":
		id = 16
	case "nox_xxx_drawMessageLines_445530":
		id = 17
	case "nox_xxx_guiChatMode_4456E0":
		id = 18
	case "nox_xxx_guiChatShowHide_445730":
		id = 19
	case "sub_445770":
		id = 20
	case "sub_469FA0":
		id = 21
	case "nox_client_chatStart_46A430":
		id = 22
	case "sub_46A4A0":
		id = 23
	case "nox_xxx_cmdSayDo_46A4B0":
		id = 24
	case "sub_46A5D0":
		id = 25
	case "sub_46A6A0":
		id = 26
	case "sub_46A730":
		id = 27
	case "sub_46A7E0":
		id = 28
	case "sub_46A820":
		id = 29
	case "sub_46A860":
		id = 30
	case "nox_xxx_packetGetMarshall_476F40":
		id = 31
	case "nox_xxx_clientEnumHover_476FA0":
		id = 32
	case "nox_xxx_clientOnCursorHover_477050":
		id = 33
	case "nox_xxx_guiCursor_477600":
		id = 34
	case "sub_479950":
		id = 35
	case "sub_4799A0":
		id = 36
	case "nox_xxx_guiDialog_479B00":
		id = 37
	case "sub_479BE0":
		id = 38
	case "sub_479C40":
		id = 39
	case "sub_479CB0":
		id = 40
	case "sub_479D00":
		id = 41
	case "sub_479D10":
		id = 42
	case "sub_47A260":
		id = 43
	case "nox_xxx_showObserverWindow_48CA70":
		id = 44
	case "sub_48D4B0":
		id = 45
	case "sub_49B3E0":
		id = 46
	case "sub_49B420":
		id = 47
	case "sub_49B490":
		id = 48
	case "sub_49B6B0":
		id = 49
	case "nox_xxx_consoleEsc_49B7A0":
		id = 50
	case "nox_xxx_wnd_49C760":
		id = 51
	case "sub_49C7A0":
		id = 52
	case "sub_49C810":
		id = 53
	case "nox_xxx_clientReportSecondaryWeapon_4BF010":
		id = 55
	case "sub_4BF7E0":
		id = 56
	case "sub_4BF9F0":
		id = 57
	case "sub_4BFAD0":
		id = 58
	case "sub_4BFB70":
		id = 59
	case "sub_4BFBB0":
		id = 60
	case "sub_4BFBF0":
		id = 61
	case "sub_4BFC70":
		id = 62
	case "sub_4BFC90":
		id = 63
	case "sub_4BFCD0":
		id = 64
	case "sub_4BFD10":
		id = 65
	case "sub_4BFD30":
		id = 66
	case "sub_4C3390":
		id = 67
	case "sub_4C3410":
		id = 68
	case "sub_4C3460":
		id = 69
	case "sub_4C34A0":
		id = 70
	case "nox_xxx_cliShowHelpGui_49C560":
		id = 71
	case "nox_xxx_clientPickup_46C140":
		id = 72
	case "nox_xxx_printCentered_445490":
		id = 73
	case "sub_49B4B0":
		id = 74
	case "sub_49B6E0":
		id = 75
	case "nox_xxx_guiChatIconLoad_445650":
		id = 76
	case "sub_48C9F0":
		id = 77
	case "sub_48C980":
		id = 78
	}
	if id < 0 {
		panic(op)
	}
	switch id {
	case 0:
		nox_xxx_clientTalk_42E7B0((*nox_drawable)(unsafe.Pointer(a[0])))
		return 0
	case 1:
		nox_xxx_clientCollideOrUse_42E810((*nox_drawable)(unsafe.Pointer(a[0])))
		return 0
	case 2:
		nox_xxx_clientTrade_42E850((*nox_drawable)(unsafe.Pointer(a[0])))
		return 0
	case 3:
		return uint64(sub_430AA0(C.int(a[0])))
	case 4:
		return uint64(nox_client_mousePriKey_430AF0())
	case 5:
		return uint64(nox_xxx_cursor_430B00())
	case 6:
		nox_client_setMousePos_430B10(C.int(a[0]), C.int(a[1]))
		return 0
	case 7:
		return uint64(nox_xxx_initTime_435570())
	case 8:
		return uint64(bool2int(bool(nox_client_drawable_testBuff_4356C0((*nox_drawable)(unsafe.Pointer(a[0])), C.char(a[1])))))
	case 9:
		return uint64(uintptr(unsafe.Pointer(sub_435700((*C.ushort)(unsafe.Pointer(a[0])), C.int(a[1])))))
	case 10:
		return uint64(nox_xxx_cliToggleObsWindow_4357A0())
	case 11:
		return uint64(sub_435F60())
	case 12:
		return uint64(sub_436550())
	case 13:
		sub_437100()
		return 0
	case 14:
		return uint64(nox_xxx_playerAnimCheck_4372B0())
	case 15:
		return uint64(nox_xxx_clientIsObserver_4372E0())
	case 16:
		return uint64(uintptr(unsafe.Pointer(sub_445450())))
	case 17:
		return uint64(nox_xxx_drawMessageLines_445530())
	case 18:
		return uint64(nox_xxx_guiChatMode_4456E0((*C.int)(unsafe.Pointer(a[0]))))
	case 19:
		return uint64(nox_xxx_guiChatShowHide_445730(C.int(a[0])))
	case 20:
		return uint64(sub_445770())
	case 21:
		return uint64(sub_469FA0())
	case 22:
		nox_client_chatStart_46A430(C.int(a[0]))
		return 0
	case 23:
		return uint64(sub_46A4A0())
	case 24:
		return uint64(nox_xxx_cmdSayDo_46A4B0((*C.ushort)(unsafe.Pointer(a[0])), C.int(a[1])))
	case 25:
		return uint64(sub_46A5D0((*C.uint32_t)(unsafe.Pointer(a[0])), C.int(a[1])))
	case 26:
		return uint64(sub_46A6A0())
	case 27:
		return uint64(uintptr(unsafe.Pointer(sub_46A730())))
	case 28:
		return uint64(sub_46A7E0((*C.uint32_t)(unsafe.Pointer(a[0])), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case 29:
		return uint64(sub_46A820(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case 30:
		return uint64(sub_46A860())
	case 31:
		return uint64(nox_xxx_packetGetMarshall_476F40())
	case 32:
		nox_xxx_clientEnumHover_476FA0()
		return 0
	case 33:
		nox_xxx_clientOnCursorHover_477050(C.int(a[0]), C.int(a[1]))
		return 0
	case 34:
		return uint64(nox_xxx_guiCursor_477600())
	case 35:
		return uint64(sub_479950())
	case 36:
		return uint64(sub_4799A0())
	case 37:
		return uint64(nox_xxx_guiDialog_479B00(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 38:
		return uint64(sub_479BE0((*C.uint32_t)(unsafe.Pointer(a[0])), C.int(a[1]), C.uint(a[2]), C.int(a[3])))
	case 39:
		return uint64(sub_479C40((*C.uint32_t)(unsafe.Pointer(a[0])), C.int(a[1])))
	case 40:
		return uint64(sub_479CB0(C.int(a[0]), C.int(a[1])))
	case 41:
		return uint64(sub_479D00())
	case 42:
		return uint64(sub_479D10())
	case 43:
		return uint64(sub_47A260())
	case 44:
		return uint64(nox_xxx_showObserverWindow_48CA70(C.int(a[0])))
	case 45:
		return uint64(sub_48D4B0(C.int(a[0])))
	case 46:
		return uint64(sub_49B3E0())
	case 47:
		return uint64(sub_49B420(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 48:
		return uint64(sub_49B490())
	case 49:
		return uint64(sub_49B6B0())
	case 50:
		nox_xxx_consoleEsc_49B7A0()
		return 0
	case 51:
		return uint64(nox_xxx_wnd_49C760(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 52:
		return uint64(sub_49C7A0())
	case 53:
		return uint64(sub_49C810())
	case 55:
		return uint64(nox_xxx_clientReportSecondaryWeapon_4BF010(C.int(a[0])))
	case 56:
		return uint64(sub_4BF7E0((*C.uint32_t)(unsafe.Pointer(a[0]))))
	case 57:
		return uint64(sub_4BF9F0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]), C.int(a[6])))
	case 58:
		return uint64(sub_4BFAD0())
	case 59:
		sub_4BFB70(C.int(a[0]))
		return 0
	case 60:
		sub_4BFBB0(C.uint(a[0]))
		return 0
	case 61:
		return uint64(sub_4BFBF0())
	case 62:
		return uint64(sub_4BFC70())
	case 63:
		return uint64(sub_4BFC90())
	case 64:
		return uint64(sub_4BFCD0(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 65:
		sub_4BFD10()
		return 0
	case 66:
		return uint64(sub_4BFD30())
	case 67:
		return uint64(sub_4C3390())
	case 68:
		return uint64(sub_4C3410((*C.int)(unsafe.Pointer(a[0]))))
	case 69:
		return uint64(sub_4C3460(C.int(a[0])))
	case 70:
		return uint64(sub_4C34A0())
	case 71:
		return uint64(uintptr(unsafe.Pointer(nox_xxx_cliShowHelpGui_49C560())))
	case 72:
		nox_xxx_clientPickup_46C140((*nox_drawable)(unsafe.Pointer(a[0])))
		return 0
	case 73:
		nox_xxx_printCentered_445490((*C.ushort)(unsafe.Pointer(a[0])))
		return 0
	case 74:
		return uint64(sub_49B4B0((*C.ushort)(unsafe.Pointer(a[0]))))
	case 75:
		return uint64(sub_49B6E0())
	case 76:
		return uint64(nox_xxx_guiChatIconLoad_445650())
	case 77:
		return uint64(sub_48C9F0((*C.int)(unsafe.Pointer(a[0]))))
	case 78:
		return uint64(sub_48C980())
	}
	return 0
}

func PortTestClientInteractionConsole(format, wide *uint16, narrow *byte, number int32) int {
	return int(C.port_test_interaction_console((*C.wchar2_t)(unsafe.Pointer(format)), (*C.wchar2_t)(unsafe.Pointer(wide)), (*C.char)(unsafe.Pointer(narrow)), C.int(number)))
}
