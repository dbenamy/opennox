package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export sub_4C1CA0
func sub_4C1CA0(command C.int) C.int { return C.int(summonSetCommand(uint32(command))) }

func nox_xxx_cliSummonCreat_4C2E50(code, typ, quiet C.int) C.char {
	return C.char(summonAdd(uint32(code), uint32(typ), quiet != 0))
}

func nox_xxx_cliSummonOnDieOrBanish_4C3140(code C.int, quiet unsafe.Pointer) {
	summonRemove(uint32(code), quiet != nil)
}

//export sub_4C3260
func sub_4C3260() C.int { return C.int(bool2int(summonFirst() != nil)) }

//export sub_4C2CE0
func sub_4C2CE0() C.int { return C.int(summonCommandTooltip()) }

//export sub_4C2C20
func sub_4C2C20(w *C.uint32_t, _ C.int, arg C.uint) C.int {
	Nox_xxx_cursorSetTooltip_4776B0(GoWStringP(summonSlotTooltip(AsWindowP(unsafe.Pointer(w)), bookPoint(uint32(arg)))))
	return 1
}
