package legacy

/*
#include "GAME4.h"
*/
import "C"
import "unsafe"

//export nox_xxx_summonStart_500DA0
func nox_xxx_summonStart_500DA0(a C.int) C.int {
	return C.int(spellEffectSummonStart(unsafe.Pointer(uintptr(a))))
}

//export nox_xxx_summonFinish_5010D0
func nox_xxx_summonFinish_5010D0(a C.int) C.int {
	return C.int(spellEffectSummonFinish(unsafe.Pointer(uintptr(a))))
}

//export nox_xxx_summonCancel_5011C0
func nox_xxx_summonCancel_5011C0(a C.int) { spellEffectSummonCancel(unsafe.Pointer(uintptr(a))) }

//export nox_xxx_charmCreature1_5011F0
func nox_xxx_charmCreature1_5011F0(a *C.int) C.int {
	return C.int(spellEffectCharmStart(unsafe.Pointer(a)))
}

//export nox_xxx_charmCreatureFinish_5013E0
func nox_xxx_charmCreatureFinish_5013E0(a *C.int) C.int {
	return C.int(spellEffectCharmFinish(unsafe.Pointer(a)))
}

//export nox_xxx_charmCreature2_501690
func nox_xxx_charmCreature2_501690(a C.int) C.int {
	return C.int(spellEffectCharmCancel(unsafe.Pointer(uintptr(a))))
}
