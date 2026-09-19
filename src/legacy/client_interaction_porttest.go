//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__gui__guicon.h"
#include <stdint.h>
void nox_xxx_clientTalk_42E7B0(nox_drawable* a1p);
void nox_xxx_clientCollideOrUse_42E810(nox_drawable* a1p);
void nox_xxx_clientTrade_42E850(nox_drawable* a1p);
int sub_430AA0(int a1);
int nox_client_mousePriKey_430AF0();
int nox_xxx_cursor_430B00();
void nox_client_setMousePos_430B10(int x, int y);
long long nox_xxx_initTime_435570();
bool nox_client_drawable_testBuff_4356C0(nox_drawable* dr, char a2);
wchar2_t* sub_435700(wchar2_t* a1, int a2);
int nox_xxx_cliToggleObsWindow_4357A0();
int sub_435F60();
int sub_436550();
void sub_437100();
int nox_xxx_playerAnimCheck_4372B0();
int nox_xxx_clientIsObserver_4372E0();
wchar2_t* sub_445450();
int nox_xxx_drawMessageLines_445530();
int nox_xxx_guiChatMode_4456E0(int* a1);
int nox_xxx_guiChatShowHide_445730(int a1);
int sub_445770();
int sub_469FA0();
void nox_client_chatStart_46A430(int a1);
int sub_46A4A0();
size_t nox_xxx_cmdSayDo_46A4B0(wchar2_t* a1, int a2);
int sub_46A5D0(uint32_t* a1, int a2);
int sub_46A6A0();
uint32_t* sub_46A730();
int sub_46A7E0(uint32_t* a1, int a2, int a3, int a4);
int sub_46A820(int a1, int a2, int a3, int a4);
int sub_46A860();
unsigned int nox_xxx_packetGetMarshall_476F40();
void nox_xxx_clientEnumHover_476FA0();
void nox_xxx_clientOnCursorHover_477050(int arg0, int a2);
int nox_xxx_guiCursor_477600();
int sub_479950();
int sub_4799A0();
int nox_xxx_guiDialog_479B00(int a1, int a2, int* a3, int a4);
int sub_479BE0(uint32_t* a1, int a2, unsigned int a3, int a4);
int sub_479C40(uint32_t* a1, int a2);
int sub_479CB0(int a1, int a2);
int sub_479D00();
int sub_479D10();
int sub_47A260();
int nox_xxx_showObserverWindow_48CA70(int a1);
int sub_48D4B0(int a1);
int sub_49B3E0();
int sub_49B420(int a1, int a2, int* a3, int a4);
int sub_49B490();
int sub_49B6B0();
void nox_xxx_consoleEsc_49B7A0();
int nox_xxx_wnd_49C760(int a1, int a2, int* a3, int a4);
int sub_49C7A0();
int sub_49C810();
int sub_49CB40();
int nox_xxx_clientReportSecondaryWeapon_4BF010(int a1);
short sub_4BF7E0(uint32_t* a1);
short sub_4BF9F0(int a1, int a2, int a3, int a4, int a5, int a6, int a7);
int sub_4BFAD0();
void sub_4BFB70(int a1);
void sub_4BFBB0(uint32_t* a1);
int sub_4BFBF0();
int sub_4BFC70();
int sub_4BFC90();
int sub_4BFCD0(int a1, int a2, int* a3, int a4);
void sub_4BFD10();
int sub_4BFD30();
int sub_4C3390();
int sub_4C3410(int* a1);
int sub_4C3460(int a1);
int sub_4C34A0();
uint32_t* nox_xxx_cliShowHelpGui_49C560();
void nox_xxx_clientPickup_46C140(nox_drawable* a1p);
void nox_xxx_printCentered_445490(wchar2_t* a1);
int sub_49B4B0(unsigned short* a1);
int sub_49B6E0();
int nox_xxx_guiChatIconLoad_445650();
int sub_48C9F0(int* a1);
int sub_48C980();
static int port_test_interaction_console(wchar2_t* format, wchar2_t* wide, char* narrow, int number) {
	return nox_gui_console_Printf_450C00(NOX_CONSOLE_RED, format, wide, narrow, number);
}
static uint64_t port_test_interaction(int op,uintptr_t a0,uintptr_t a1,uintptr_t a2,uintptr_t a3,uintptr_t a4,uintptr_t a5,uintptr_t a6) {
switch(op) {
case 0: nox_xxx_clientTalk_42E7B0((nox_drawable*)a0); return 0;
case 1: nox_xxx_clientCollideOrUse_42E810((nox_drawable*)a0); return 0;
case 2: nox_xxx_clientTrade_42E850((nox_drawable*)a0); return 0;
case 3: return (uint64_t)sub_430AA0((int)a0);
case 4: return (uint64_t)nox_client_mousePriKey_430AF0();
case 5: return (uint64_t)nox_xxx_cursor_430B00();
case 6: nox_client_setMousePos_430B10((int)a0,(int)a1); return 0;
case 7: return (uint64_t)nox_xxx_initTime_435570();
case 8: return (uint64_t)nox_client_drawable_testBuff_4356C0((nox_drawable*)a0,(char)a1);
case 9: return (uint64_t)(uintptr_t)sub_435700((wchar2_t*)a0,(int)a1);
case 10: return (uint64_t)nox_xxx_cliToggleObsWindow_4357A0();
case 11: return (uint64_t)sub_435F60();
case 12: return (uint64_t)sub_436550();
case 13: sub_437100(); return 0;
case 14: return (uint64_t)nox_xxx_playerAnimCheck_4372B0();
case 15: return (uint64_t)nox_xxx_clientIsObserver_4372E0();
case 16: return (uint64_t)(uintptr_t)sub_445450();
case 17: return (uint64_t)nox_xxx_drawMessageLines_445530();
case 18: return (uint64_t)nox_xxx_guiChatMode_4456E0((int*)a0);
case 19: return (uint64_t)nox_xxx_guiChatShowHide_445730((int)a0);
case 20: return (uint64_t)sub_445770();
case 21: return (uint64_t)sub_469FA0();
case 22: nox_client_chatStart_46A430((int)a0); return 0;
case 23: return (uint64_t)sub_46A4A0();
case 24: return (uint64_t)nox_xxx_cmdSayDo_46A4B0((wchar2_t*)a0,(int)a1);
case 25: return (uint64_t)sub_46A5D0((uint32_t*)a0,(int)a1);
case 26: return (uint64_t)sub_46A6A0();
case 27: return (uint64_t)(uintptr_t)sub_46A730();
case 28: return (uint64_t)sub_46A7E0((uint32_t*)a0,(int)a1,(int)a2,(int)a3);
case 29: return (uint64_t)sub_46A820((int)a0,(int)a1,(int)a2,(int)a3);
case 30: return (uint64_t)sub_46A860();
case 31: return (uint64_t)nox_xxx_packetGetMarshall_476F40();
case 32: nox_xxx_clientEnumHover_476FA0(); return 0;
case 33: nox_xxx_clientOnCursorHover_477050((int)a0,(int)a1); return 0;
case 34: return (uint64_t)nox_xxx_guiCursor_477600();
case 35: return (uint64_t)sub_479950();
case 36: return (uint64_t)sub_4799A0();
case 37: return (uint64_t)nox_xxx_guiDialog_479B00((int)a0,(int)a1,(int*)a2,(int)a3);
case 38: return (uint64_t)sub_479BE0((uint32_t*)a0,(int)a1,(unsigned int)a2,(int)a3);
case 39: return (uint64_t)sub_479C40((uint32_t*)a0,(int)a1);
case 40: return (uint64_t)sub_479CB0((int)a0,(int)a1);
case 41: return (uint64_t)sub_479D00();
case 42: return (uint64_t)sub_479D10();
case 43: return (uint64_t)sub_47A260();
case 44: return (uint64_t)nox_xxx_showObserverWindow_48CA70((int)a0);
case 45: return (uint64_t)sub_48D4B0((int)a0);
case 46: return (uint64_t)sub_49B3E0();
case 47: return (uint64_t)sub_49B420((int)a0,(int)a1,(int*)a2,(int)a3);
case 48: return (uint64_t)sub_49B490();
case 49: return (uint64_t)sub_49B6B0();
case 50: nox_xxx_consoleEsc_49B7A0(); return 0;
case 51: return (uint64_t)nox_xxx_wnd_49C760((int)a0,(int)a1,(int*)a2,(int)a3);
case 52: return (uint64_t)sub_49C7A0();
case 53: return (uint64_t)sub_49C810();
case 54: return (uint64_t)sub_49CB40();
case 55: return (uint64_t)nox_xxx_clientReportSecondaryWeapon_4BF010((int)a0);
case 56: return (uint64_t)sub_4BF7E0((uint32_t*)a0);
case 57: return (uint64_t)sub_4BF9F0((int)a0,(int)a1,(int)a2,(int)a3,(int)a4,(int)a5,(int)a6);
case 58: return (uint64_t)sub_4BFAD0();
case 59: sub_4BFB70((int)a0); return 0;
case 60: sub_4BFBB0((uint32_t*)a0); return 0;
case 61: return (uint64_t)sub_4BFBF0();
case 62: return (uint64_t)sub_4BFC70();
case 63: return (uint64_t)sub_4BFC90();
case 64: return (uint64_t)sub_4BFCD0((int)a0,(int)a1,(int*)a2,(int)a3);
case 65: sub_4BFD10(); return 0;
case 66: return (uint64_t)sub_4BFD30();
case 67: return (uint64_t)sub_4C3390();
case 68: return (uint64_t)sub_4C3410((int*)a0);
case 69: return (uint64_t)sub_4C3460((int)a0);
case 70: return (uint64_t)sub_4C34A0();
case 71: return (uint64_t)(uintptr_t)nox_xxx_cliShowHelpGui_49C560();
case 72: nox_xxx_clientPickup_46C140((nox_drawable*)a0); return 0;
case 73: nox_xxx_printCentered_445490((wchar2_t*)a0); return 0;
case 74: return (uint64_t)sub_49B4B0((unsigned short*)a0);
case 75: return (uint64_t)sub_49B6E0();
case 76: return (uint64_t)nox_xxx_guiChatIconLoad_445650();
case 77: return (uint64_t)sub_48C9F0((int*)a0);
case 78: return (uint64_t)sub_48C980();
}
return 0;
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
	case "sub_49CB40":
		id = 54
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
	return uint64(C.port_test_interaction(C.int(id), C.uintptr_t(a[0]), C.uintptr_t(a[1]), C.uintptr_t(a[2]), C.uintptr_t(a[3]), C.uintptr_t(a[4]), C.uintptr_t(a[5]), C.uintptr_t(a[6])))
}

func PortTestClientInteractionConsole(format, wide *uint16, narrow *byte, number int32) int {
	return int(C.port_test_interaction_console((*C.wchar2_t)(unsafe.Pointer(format)), (*C.wchar2_t)(unsafe.Pointer(wide)), (*C.char)(unsafe.Pointer(narrow)), C.int(number)))
}
