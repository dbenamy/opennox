package legacy

/*
#include "GAME2.h"
#include "client__gui__guibrief.h"
*/
import "C"
import (
	"unsafe"
)

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

//export nox_client_lockScreenBriefing_450160
func nox_client_lockScreenBriefing_450160(chapter, begin C.int, mode C.char) C.int {
	return C.int(briefingShow(int(chapter), int(begin), byte(mode)))
}
