package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_unitSparkInit_4F0390
func nox_xxx_unitSparkInit_4F0390(a C.int) *C.uint32_t {
	return (*C.uint32_t)(rewardInitSpark(objectFromInt(a)))
}

//export nox_xxx_initFrog_4F03B0
func nox_xxx_initFrog_4F03B0(a C.int) C.int { return C.int(rewardInitFrog(objectFromInt(a))) }

//export nox_xxx_initChest_4F0400
func nox_xxx_initChest_4F0400(a C.int) *C.int {
	u := objectFromInt(a)
	rewardInitBreakable(u)
	return (*C.int)(u.CObj())
}

//export nox_xxx_unitBoulderInit_4F0420
func nox_xxx_unitBoulderInit_4F0420(a *C.uint32_t) *C.uint32_t {
	u := equipmentObject(unsafe.Pointer(a))
	u.Pos39 = u.PosVec
	return a
}

//export sub_4F0450
func sub_4F0450(a C.int) C.int { return C.int(rewardInitDirection(objectFromInt(a), true)) }

//export sub_4F0490
func sub_4F0490(a C.int) C.int { return C.int(rewardInitDirection(objectFromInt(a), false)) }

//export nox_xxx_unitInitGold_4F04B0
func nox_xxx_unitInitGold_4F04B0(a C.int) C.int { return C.int(rewardInitGold(objectFromInt(a))) }

//export nox_xxx_breakInit_4F0570
func nox_xxx_breakInit_4F0570(a C.int) *C.int {
	u := objectFromInt(a)
	rewardInitBreakable(u)
	return (*C.int)(u.CObj())
}

//export nox_xxx_unitInitGenerator_4F0590
func nox_xxx_unitInitGenerator_4F0590(a C.int) C.int {
	return C.int(rewardInitGenerator(objectFromInt(a)))
}
