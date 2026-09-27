package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export sub_4C3260
func sub_4C3260() C.int { return C.int(bool2int(summonFirst() != nil)) }

//export sub_4C2CE0
func sub_4C2CE0() C.int { return C.int(summonCommandTooltip()) }

//export sub_4C2C20
func sub_4C2C20(w *C.uint32_t, _ C.int, arg C.uint) C.int {
	Nox_xxx_cursorSetTooltip_4776B0(GoWStringP(summonSlotTooltip(AsWindowP(unsafe.Pointer(w)), bookPoint(uint32(arg)))))
	return 1
}
