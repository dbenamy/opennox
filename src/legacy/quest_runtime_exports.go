package legacy

/*
#include <stdint.h>
#include "defs.h"
*/
import "C"
import "unsafe"

//export sub_4D6000
func sub_4D6000(a *nox_object_t) C.int { return C.int(questRuntimeReset(asObjectS(a))) }

//export sub_4D60B0
func sub_4D60B0() C.int { return C.int(questRuntimeResetAll()) }

//export sub_4D60E0
func sub_4D60E0(a C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(questRuntimeStageComplete(objectFromInt(a)))))
}

//export sub_4D6130
func sub_4D6130(a C.int) C.int { return C.int(questRuntimeIncrement(objectFromInt(a), 4660, 2)) }

//export sub_4D61F0
func sub_4D61F0(a C.int) C.int { return C.int(questRuntimeIncrement(objectFromInt(a), 4672, 16)) }

//export sub_4D6540
func sub_4D6540(a C.int) C.uint { return C.uint(questRuntimePlayerScore(int(a))) }

//export nox_xxx_isQuest_4D6F50
func nox_xxx_isQuest_4D6F50() C.int { return C.int(questRuntimeWord(1556160)) }

//export nox_xxx_setQuest_4D6F60
func nox_xxx_setQuest_4D6F60(a C.int) C.int { return C.int(questRuntimeSetWord(1556160, uint32(a))) }

//export sub_4D6F70
func sub_4D6F70() C.int { return C.int(questRuntimeWord(1556164)) }

//export sub_4D6F80
func sub_4D6F80(a C.int) C.int { return C.int(questRuntimeSetWord(1556164, uint32(a))) }

//export sub_4D6FA0
func sub_4D6FA0() C.int { return C.int(questRuntimeWord(1556104)) }

//export sub_4D72D0
func sub_4D72D0(a C.int) C.int { return C.int(questRuntimePreviousStage(uint32(a))) }

//export sub_4D7450
func sub_4D7450(a C.int, b C.short) C.int {
	return C.int(questRuntimeHighestMessage(int(a), uint16(b)))
}

//export sub_4D75E0
func sub_4D75E0() C.int { return C.int(questRuntimeWord(1556120)) }

//export sub_4D76E0
func sub_4D76E0(a C.int) C.int { return C.int(questRuntimeSetWord(1556124, uint32(a))) }

//export sub_4D7A60
func sub_4D7A60(a C.int) C.int { return C.int(questRuntimeDepartureStamp(int(a))) }

//export sub_4E3CA0
func sub_4E3CA0() C.double { return C.double(questRuntimeFloat(202024)) }

//export nox_game_getQuestStage_4E3CC0
func nox_game_getQuestStage_4E3CC0() C.int { return C.int(questRuntimeStage()) }

//export nox_game_setQuestStage_4E3CD0
func nox_game_setQuestStage_4E3CD0(a C.int) { questRuntimeSetStage(uint32(a)) }

//export nox_xxx_player_4E3CE0
func nox_xxx_player_4E3CE0() C.int { return C.int(questRuntimeCount()) }

//export sub_4E3D50
func sub_4E3D50() C.int { return C.int(questRuntimeDifficulty()) }

//export sub_4E3DD0
func sub_4E3DD0() { questRuntimeScaleHealth() }

//export sub_4E40F0
func sub_4E40F0() C.double { return C.double(questRuntimeFloat(202036)) }
