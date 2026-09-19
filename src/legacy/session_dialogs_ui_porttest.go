//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1_3.h"
#include "GAME2_2.h"
#include "GAME3.h"
#include "client__gui__guiquit.h"
int sub_489FF0(int,int,const void*);
extern uint32_t dword_5d4594_1193380,dword_5d4594_1193384;
extern uint32_t dword_5d4594_1309748,dword_5d4594_1309756;
extern uint32_t dword_5d4594_826028,dword_5d4594_826032;
extern void* dword_5d4594_826036;
extern uint32_t dword_5d4594_1305680,dword_5d4594_1305684;
extern nox_window* nox_wnd_quitMenu_825760;
static uintptr_t port_test_motd_free(int slot) { return (uintptr_t)sub_446490(slot); }
static uintptr_t port_test_motd_add(uint8_t* text) { return (uintptr_t)nox_xxx_motdAddSomeTextMB_446730(text); }
*/
import "C"
import "unsafe"

func PortTestSessionDialogWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"otherDialogA":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305680)),
		"otherDialogB":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305684)),
		"filter":         (*uint32)(unsafe.Pointer(&C.dword_5d4594_1193380)),
		"filterControls": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1193384)),
		"disconnect":     (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309748)),
		"disconnectIcon": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309756)),
		"motd":           (*uint32)(unsafe.Pointer(&C.dword_5d4594_826028)),
		"motdList":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_826032)),
		"motdFile":       (*uint32)(unsafe.Pointer(&C.dword_5d4594_826036)),
		"quit":           (*uint32)(unsafe.Pointer(&C.nox_wnd_quitMenu_825760)),
	}
	old := make(map[string]uint32)
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

func PortTestSessionDialogCall(op string, w unsafe.Pointer, code int, a, b uint32) uintptr {
	switch op {
	case "filterOpen":
		return uintptr(unsafe.Pointer(C.sub_489B80(C.int(uintptr(w)))))
	case "filterSave":
		return uintptr(uint32(C.sub_489870()))
	case "filterControls":
		C.sub_489DC0()
		return 0
	case "filterClose":
		return uintptr(uint32(C.sub_489FB0()))
	case "filterConfig":
		return uintptr(uint32(C.sub_489FF0(C.int(a), C.int(b), w)))
	case "filterEvent":
		return uintptr(uint32(C.nox_xxx_windowMplayFilterProc_489E70(C.int(uintptr(w)), C.int(code), (*C.int)(unsafe.Pointer(uintptr(a))), C.int(b))))
	case "motdOpen":
		return uintptr(uint32(C.nox_xxx_guiMotdLoad_4465C0()))
	case "motdClose":
		return uintptr(uint32(C.sub_446780()))
	case "motdShown":
		return uintptr(uint32(C.sub_446950()))
	case "motdShow":
		C.nox_xxx_motd_4467F0()
		return 0
	case "motdEvent":
		return uintptr(uint32(C.sub_4466C0(C.int(uintptr(w)), C.int(code), C.int(a), C.int(b))))
	case "motdAdd":
		return uintptr(C.port_test_motd_add((*C.uint8_t)(w)))
	case "motdRead":
		return uintptr(uint32(C.nox_motd_4463E0(C.int(a))))
	case "motdFree":
		return uintptr(C.port_test_motd_free(C.int(a)))
	case "disconnectOpen":
		return uintptr(uint32(C.sub_4AB260()))
	case "disconnectClose":
		return uintptr(uint32(C.sub_4AB470()))
	case "disconnectIcon":
		return uintptr(uint32(C.sub_4AB4A0(C.int(a))))
	case "disconnectShow":
		return uintptr(uint32(C.sub_4AB4D0(C.int(a))))
	case "disconnectInput":
		return uintptr(uint32(C.sub_4AB340(C.int(uintptr(w)), C.int(code), C.int(a), C.int(b))))
	case "disconnectEvent":
		return uintptr(uint32(C.sub_4AB390(C.int(uintptr(w)), C.int(code), (*C.int)(unsafe.Pointer(uintptr(a))), C.int(b))))
	case "disconnectDraw":
		return uintptr(uint32(C.sub_4AB420((*C.int)(w))))
	case "quitShown":
		return uintptr(C.nox_gui_xxx_check_446360())
	case "quitColors":
		return uintptr(unsafe.Pointer(C.sub_445FF0()))
	case "quitCapture":
		return uintptr(uint32(C.nox_xxx_quitDialogNo_445B30()))
	case "quitHide":
		C.sub_445C20()
		return 0
	case "quitToggle":
		C.sub_445C40()
		return 0
	case "quitEvent":
		return uintptr(uint32(C.nox_xxx_menuGameOnButton_445840((*C.uint32_t)(w), C.int(code), (*C.int)(unsafe.Pointer(uintptr(a))), C.int(b))))
	}
	panic(op)
}
