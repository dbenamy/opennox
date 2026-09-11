package legacy

/*
#include "GAME1.h"
#include "GAME3_3.h"
#include "GAME4_3.h"
*/
import "C"
import "unsafe"

//export sub_417F50
func sub_417F50(a C.int) C.int { return C.int(objectiveBallReset(objectFromInt(a))) }

//export nox_xxx_pickupFlagCtf_4EA490
func nox_xxx_pickupFlagCtf_4EA490(a, b C.int) { objectiveCTFPickup(objectFromInt(a), objectFromInt(b)) }

//export sub_4EB9B0
func sub_4EB9B0(a, b C.int) C.int {
	return inventoryInt(objectiveRememberOwner(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collideBall_4EBA00
func nox_xxx_collideBall_4EBA00(a, b C.int) { objectiveBallCollide(objectFromInt(a), objectFromInt(b)) }

//export sub_4EBB50
func sub_4EBB50(a, b C.int) C.int {
	return C.int(objectiveCrownCollide(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collideHomeBase_4EBB80
func nox_xxx_collideHomeBase_4EBB80(a, b C.int) C.short {
	return C.short(objectiveHomeBase(objectFromInt(a), objectFromInt(b)))
}

//export sub_4ECBD0
func sub_4ECBD0(a C.int) C.int { return C.int(objectiveFlagID(objectFromInt(a))) }

//export sub_4ECC00
func sub_4ECC00(a **C.char) C.int { return C.int(objectiveColor(unsafe.Pointer(a))) }

//export nox_xxx_updateObelisk_53C580
func nox_xxx_updateObelisk_53C580(a C.int) C.int { return C.int(objectiveObelisk(objectFromInt(a))) }

//export nox_xxx_updateFlag_53DDF0
func nox_xxx_updateFlag_53DDF0(a C.int) C.int { return C.int(objectiveFlagUpdate(objectFromInt(a))) }

//export nox_xxx_updateGameBall_53DF40
func nox_xxx_updateGameBall_53DF40(a C.int) { objectiveBallUpdate(objectFromInt(a)) }

//export nox_xxx_updateCrown_53E1D0
func nox_xxx_updateCrown_53E1D0(a C.int) { objectiveCrownUpdate(objectFromInt(a)) }

//export sub_4EA400
func sub_4EA400(a, b C.int) { objectiveFlagCollide(objectFromInt(a), objectFromInt(b)) }

//export sub_4EA7A0
func sub_4EA7A0(a C.int) C.int { return C.int(objectivePickupBuffs(objectFromInt(a))) }

//export sub_4EA800
func sub_4EA800(a, b C.int) C.short {
	return C.short(objectiveFlagBallScore(objectFromInt(a), objectFromInt(b)))
}
