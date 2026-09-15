package legacy

/*
#include "GAME2.h"
#include "client__gui__guibrief.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

//export sub_44E410
func sub_44E410() *C.wchar2_t { return (*C.wchar2_t)(unsafe.Pointer(briefingLoadChapters())) }

//export sub_44E8E0
func sub_44E8E0(a, b C.int) C.int {
	return C.int(briefingDrawStats((*gui.WindowData)(unsafe.Pointer(uintptr(uint32(b))))))
}

//export sub_44F0F0
func sub_44F0F0(a, b C.int) C.int {
	return C.int(briefingDrawTitle((*gui.WindowData)(unsafe.Pointer(uintptr(uint32(b))))))
}

//export sub_44F300
func sub_44F300(a, b C.int) C.int {
	return C.int(briefingDrawInstructions((*gui.WindowData)(unsafe.Pointer(uintptr(uint32(b))))))
}

//export nox_xxx_clientQuestWinScreen_450770
func nox_xxx_clientQuestWinScreen_450770(a C.int) C.int {
	return C.int(briefingWinReport(unsafe.Pointer(uintptr(uint32(a)))))
}

//export nox_client_showQuestBriefing2_450980
func nox_client_showQuestBriefing2_450980(a, b C.int) C.int {
	return C.int(briefingSelection(unsafe.Pointer(uintptr(uint32(a))), int(b), false))
}

//export nox_client_showQuestBriefing_450A30
func nox_client_showQuestBriefing_450A30(a, b C.int) C.int {
	return C.int(briefingSelection(unsafe.Pointer(uintptr(uint32(a))), int(b), true))
}

//export nox_gui_getQuestStage_450B10
func nox_gui_getQuestStage_450B10() C.int { return C.int(briefingStage()) }
