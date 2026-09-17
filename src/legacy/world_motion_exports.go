package legacy

/*
#include "GAME4_1.h"
#include "GAME5.h"
*/
import "C"

//export nox_xxx_updateSentryGlobe_510E60
func nox_xxx_updateSentryGlobe_510E60(a C.int) C.int {
	return C.int(motionSentryUpdate(objectFromInt(a)))
}

//export nox_xxx_unitSetDecayTime_511660
func nox_xxx_unitSetDecayTime_511660(a *C.nox_object_t, delay C.int) C.int {
	return C.int(motionDecaySet(asObjectS(a), int32(delay)))
}

//export sub_51B860
func sub_51B860(a C.int) C.char { return C.char(collisionActivate(objectFromInt(a))) }

//export nox_xxx_unitUpdateMover_54F740
func nox_xxx_unitUpdateMover_54F740(a C.int) { motionMover(objectFromInt(a)) }

//export nox_xxx_updateShootingTrap_54F9A0
func nox_xxx_updateShootingTrap_54F9A0(a C.int) C.int {
	return C.int(motionTrapUpdate(objectFromInt(a)))
}

//export nox_xxx_collideTrigger_54FCD0
func nox_xxx_collideTrigger_54FCD0(a, b C.int) { motionTrigger(objectFromInt(a), objectFromInt(b)) }
