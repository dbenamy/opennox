package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func sub_460380() unsafe.Pointer { return unsafe.Pointer(uintptr(quickbarClearAbilities())) }

func nox_xxx_cliPrepareGameplay1_460E60() C.int { return C.int(quickbarPrepare()) }

func sub_460EB0(id C.int, value C.char) { quickbarSetFlash(int(id), byte(value)) }

func sub_461090(id, state C.int) *C.char {
	return (*C.char)(unsafe.Pointer(quickbarAbilityState(uint32(id), uint32(state))))
}

func sub_4610D0(id C.uchar) *C.char { return (*C.char)(unsafe.Pointer(quickbarResetAbility(byte(id)))) }

func sub_461120(id, on C.int) *C.char {
	return (*C.char)(unsafe.Pointer(quickbarAbilityFlags(uint32(id), uint32(on))))
}

func nox_xxx_netAbilityRewardCli_4611E0(id, rank C.int, notify *C.char) {
	quickbarAbilityReward(int(id), int(rank), uintptr(unsafe.Pointer(notify)))
}

//export nox_xxx_quickbarButtonBook_45F3F0
func nox_xxx_quickbarButtonBook_45F3F0() C.int { return C.int(quickbarBookTooltip()) }

//export sub_45F480
func sub_45F480(w C.int) C.int { return C.int(quickbarDirectionTooltip(bookWindow(uint32(w)))) }
