package legacy

import (
	"unsafe"
)

func nox_xxx_clientQuestWinScreen_450770(a int32) int32 {
	return int32(briefingWinReport(unsafe.Pointer(uintptr(uint32(a)))))
}

func nox_client_showQuestBriefing2_450980(a, b int32) int32 {
	return int32(briefingSelection(unsafe.Pointer(uintptr(uint32(a))), int(b), false))
}

func nox_client_showQuestBriefing_450A30(a, b int32) int32 {
	return int32(briefingSelection(unsafe.Pointer(uintptr(uint32(a))), int(b), true))
}

func nox_client_lockScreenBriefing_450160(chapter, begin int32, mode int8) int32 {
	return int32(briefingShow(int(chapter), int(begin), byte(mode)))
}
