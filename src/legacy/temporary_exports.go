package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_updateSpark_53ADC0
func nox_xxx_updateSpark_53ADC0(a1 C.int) { temporarySparkUpdate(objectFromInt(a1)) }

//export nox_xxx_updateProjTrail_53AEC0
func nox_xxx_updateProjTrail_53AEC0(a1 C.int) *C.float {
	return (*C.float)(temporaryTrail(objectFromInt(a1)).CObj())
}

//export nox_xxx_updateLifetime_53B8F0
func nox_xxx_updateLifetime_53B8F0(unit C.int) { temporaryLifetime(objectFromInt(unit)) }

//export nox_xxx_spellFlyUpdate_53B940
func nox_xxx_spellFlyUpdate_53B940(a1 C.int) { temporarySpellFly(objectFromInt(a1)) }

//export nox_xxx_updateAntiSpellProj_53BB00
func nox_xxx_updateAntiSpellProj_53BB00(a1 C.int) { temporaryAntiSpell(objectFromInt(a1)) }

//export sub_53BD10
func sub_53BD10(a1 C.int, a2 C.int) { temporaryAntiCandidate(objectFromInt(a1), objectFromInt(a2)) }

//export nox_xxx_updateMagicMissile_53BDA0
func nox_xxx_updateMagicMissile_53BDA0(a1 C.int) C.int {
	return C.int(temporaryMagicMissile(objectFromInt(a1)))
}

//export nox_xxx_updateBlackPowderBarrel_53C9A0
func nox_xxx_updateBlackPowderBarrel_53C9A0(a1 *C.float) {
	temporaryPowderBarrel((*server.Object)(unsafe.Pointer(a1)))
}

//export nox_xxx_updateOneSecondDie_53CB60
func nox_xxx_updateOneSecondDie_53CB60(a1 C.int) { temporaryOneSecond(objectFromInt(a1)) }

//export nox_xxx_updateWaterBarrel_53CB90
func nox_xxx_updateWaterBarrel_53CB90(a1 C.int) { temporaryWaterBarrel(objectFromInt(a1)) }

//export nox_xxx_waterBarrel_53CC30
func nox_xxx_waterBarrel_53CC30(a1 *C.float, a2 C.int) {
	temporaryWaterCandidate((*server.Object)(unsafe.Pointer(a1)), *(*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))))
}

//export nox_xxx_updateSelfDestruct_53CC90
func nox_xxx_updateSelfDestruct_53CC90(a1 C.int) { temporarySelfDestruct(objectFromInt(a1)) }

//export nox_xxx_updateBlackPowderBurn_53CCB0
func nox_xxx_updateBlackPowderBurn_53CCB0(a1 C.int) { temporaryPowderBurn(objectFromInt(a1)) }

//export nox_xxx_updateDeathBallFragment_53D220
func nox_xxx_updateDeathBallFragment_53D220(a1 C.int) { temporaryDeathFragment(objectFromInt(a1)) }

//export nox_xxx_updateMoonglow_53D270
func nox_xxx_updateMoonglow_53D270(a1 C.int) { temporaryMoonglow(objectFromInt(a1)) }

//export nox_xxx_updateTelekinesis_53D330
func nox_xxx_updateTelekinesis_53D330(a1 C.int) { temporaryTelekinesis(objectFromInt(a1)) }

//export nox_xxx_updateFist_53D400
func nox_xxx_updateFist_53D400(a1 C.int) { temporaryFist(objectFromInt(a1)) }

//export nox_xxx_updateFlameCleanse_53D510
func nox_xxx_updateFlameCleanse_53D510(a1 C.int) { temporaryFlameCleanse(objectFromInt(a1)) }

//export nox_xxx_updateMeteorShower_53D5A0
func nox_xxx_updateMeteorShower_53D5A0(a2 *C.float) {
	temporaryMeteorShower((*server.Object)(unsafe.Pointer(a2)))
}

//export nox_xxx_meteorExplode_53D6E0
func nox_xxx_meteorExplode_53D6E0(a6 C.int) { temporaryMeteorExplode(objectFromInt(a6)) }

//export nox_xxx_updateToxicCloud_53D850
func nox_xxx_updateToxicCloud_53D850(a1 C.int) { temporaryCloud(objectFromInt(a1), false) }

//export sub_53D8C0
func sub_53D8C0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(a1), objectFromInt(a2), false)
}

//export nox_xxx_updateSmallToxicCloud_53D960
func nox_xxx_updateSmallToxicCloud_53D960(a1 C.int) { temporaryCloud(objectFromInt(a1), true) }

//export nox_xxx_toxicCloudPoison_53D9D0
func nox_xxx_toxicCloudPoison_53D9D0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(a1), objectFromInt(a2), true)
}

//export nox_xxx_updateArachnaphobia_53DA60
func nox_xxx_updateArachnaphobia_53DA60(a1 *C.int) {
	temporaryArachnaphobia((*server.Object)(unsafe.Pointer(a1)))
}

//export nox_xxx_updateExpire_53DB00
func nox_xxx_updateExpire_53DB00(a1 C.int) { temporaryExpire(objectFromInt(a1)) }

//export nox_xxx_updateBreak_53DB30
func nox_xxx_updateBreak_53DB30(a1 *C.uint32_t) *C.int {
	return (*C.int)(unsafe.Pointer(uintptr(temporaryBreak((*server.Object)(unsafe.Pointer(a1)), false))))
}

//export nox_xxx_updateOpen_53DBB0
func nox_xxx_updateOpen_53DBB0(a1 *C.uint32_t) *C.int {
	return (*C.int)(unsafe.Pointer(uintptr(temporaryBreak((*server.Object)(unsafe.Pointer(a1)), true))))
}

//export nox_xxx_updateBreakAndRemove_53DC30
func nox_xxx_updateBreakAndRemove_53DC30(a1 *C.uint32_t) {
	temporaryBreakRemove((*server.Object)(unsafe.Pointer(a1)))
}

//export nox_xxx_updateChakramInMotion_53DCC0
func nox_xxx_updateChakramInMotion_53DCC0(a1 C.int) { temporaryChakram(objectFromInt(a1)) }

//export nox_xxx_createSpark_54FD80
func nox_xxx_createSpark_54FD80(a1 C.float, a2 C.float, a3 C.int, a4 C.int, a5 C.float, a6 C.float, a7 C.float, a8 C.int) *C.float {
	return (*C.float)(temporarySpark(types.Ptf(float32(a1), float32(a2)), types.Ptf(float32(a5), float32(a6)), int32(a3), int32(a4), float32(a7), objectFromInt(a8)).CObj())
}
