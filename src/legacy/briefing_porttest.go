//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "client__gui__guibrief.h"
*/
import "C"
import "unsafe"

// PortTestBriefing calls the real presentation owners and their remaining callers.
func PortTestBriefing(op int, a, b, c uintptr) uint32 {
	switch op {
	case 0:
		return uint32(uintptr(unsafe.Pointer(C.sub_44E410())))
	case 1:
		return uint32(C.sub_44E8E0(C.int(a), C.int(b)))
	case 2:
		return uint32(C.sub_44F0F0(C.int(a), C.int(b)))
	case 3:
		return uint32(C.sub_44F300(C.int(a), C.int(b)))
	case 4:
		return uint32(C.nox_xxx_clientQuestWinScreen_450770(C.int(a)))
	case 5:
		return uint32(C.nox_client_showQuestBriefing2_450980(C.int(a), C.int(b)))
	case 6:
		return uint32(C.nox_client_showQuestBriefing_450A30(C.int(a), C.int(b)))
	case 7:
		return uint32(uintptr(unsafe.Pointer(C.sub_44E110())))
	case 8:
		return uint32(C.sub_450960(unsafe.Pointer(a), unsafe.Pointer(b)))
	case 9:
		return uint32(uintptr(unsafe.Pointer(C.sub_450AD0((*C.char)(unsafe.Pointer(a))))))
	case 10:
		return uint32(C.sub_450AF0(C.int(a)))
	case 11:
		C.nox_gui_setQuestStage_450B00(C.int(a))
	case 12:
		return uint32(C.nox_gui_getQuestStage_450B10())
	case 13:
		return uint32(uintptr(unsafe.Pointer(C.sub_44E560())))
	case 14:
		return uint32(C.sub_4505E0())
	case 15:
		return uint32(C.nox_client_lockScreenBriefing_450160(C.int(a), C.int(b), C.char(c)))
	}
	return 0
}
func PortTestBriefingCallbacks() []unsafe.Pointer {
	return []unsafe.Pointer{C.sub_44E8E0, C.sub_44F0F0, C.sub_44F300, C.nox_xxx_wndProc_44E6E0, C.nox_client_wndQuestBriefProc_44E630, C.sub_44E6F0}
}
