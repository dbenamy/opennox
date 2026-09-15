//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "client__gui__guibrief.h"
*/
import "C"
import "unsafe"
import "github.com/opennox/opennox/v1/client/gui"

// PortTestBriefing calls the real presentation owners and their remaining callers.
func PortTestBriefing(op int, a, b, c uintptr) uint32 {
	switch op {
	case 0:
		return uint32(uintptr(unsafe.Pointer(briefingLoadChapters())))
	case 1:
		return uint32(briefingDrawStats((*gui.WindowData)(unsafe.Pointer(b))))
	case 2:
		return uint32(briefingDrawTitle((*gui.WindowData)(unsafe.Pointer(b))))
	case 3:
		return uint32(briefingDrawInstructions((*gui.WindowData)(unsafe.Pointer(b))))
	case 4:
		return uint32(briefingWinReport(unsafe.Pointer(a)))
	case 5:
		return uint32(briefingSelection(unsafe.Pointer(a), int(b), false))
	case 6:
		return uint32(briefingSelection(unsafe.Pointer(a), int(b), true))
	case 7:
		return uint32(uintptr(unsafe.Pointer(briefingInitSprites())))
	case 8:
		return uint32(briefingCompare((*briefingScore)(unsafe.Pointer(a)).Total, (*briefingScore)(unsafe.Pointer(b)).Total))
	case 9:
		return briefingSetImage(uint32(a))
	case 10:
		return briefingSetCaption(uint32(a))
	case 11:
		briefingSetStage(uint32(a))
	case 12:
		return briefingStage()
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
