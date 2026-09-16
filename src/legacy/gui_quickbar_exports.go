package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_book_45DBE0
func nox_xxx_book_45DBE0(kind unsafe.Pointer, id, slot C.int) unsafe.Pointer {
	return unsafe.Pointer(quickbarBookSlot(uintptr(kind), uint32(id), int(slot)))
}

//export nox_xxx_clientUpdateButtonRow_45E110
func nox_xxx_clientUpdateButtonRow_45E110(row C.int) C.int { return C.int(quickbarSelectRow(int(row))) }

//export nox_xxx_buttonsGetSelectedRow_45E180
func nox_xxx_buttonsGetSelectedRow_45E180() C.int { return C.int(quickbarMain().Selected) }

//export sub_4602F0
func sub_4602F0() unsafe.Pointer { return unsafe.Pointer(uintptr(quickbarClearSlots())) }

//export sub_460380
func sub_460380() unsafe.Pointer { return unsafe.Pointer(uintptr(quickbarClearAbilities())) }

//export nox_client_trapSetSelect_4604B0
func nox_client_trapSetSelect_4604B0(row C.int) C.int { return C.int(quickbarTrapSelect(int(row))) }

//export sub_4604E0
func sub_4604E0() C.int { return C.int(*quickbarByte(1048140)) }

//export sub_460660
func sub_460660() C.int { return C.int(quickbarCancelCapture()) }

//export nox_xxx_quickBarClose_4606B0
func nox_xxx_quickBarClose_4606B0() C.int { return C.int(quickbarCloseExpanded()) }

//export sub_460940
func sub_460940(_ unsafe.Pointer) C.int { return C.int(quickbarSave()) }

//export nox_xxx_cliPrepareGameplay1_460E60
func nox_xxx_cliPrepareGameplay1_460E60() C.int { return C.int(quickbarPrepare()) }

//export sub_460EA0
func sub_460EA0(show C.int) C.int { return C.int(quickbarVisible(show != 0)) }

//export sub_460EB0
func sub_460EB0(id C.int, value C.char) { quickbarSetFlash(int(id), byte(value)) }

//export sub_461090
func sub_461090(id, state C.int) *C.char {
	return (*C.char)(unsafe.Pointer(quickbarAbilityState(uint32(id), uint32(state))))
}

//export sub_4610D0
func sub_4610D0(id C.uchar) *C.char { return (*C.char)(unsafe.Pointer(quickbarResetAbility(byte(id)))) }

//export sub_461120
func sub_461120(id, on C.int) *C.char {
	return (*C.char)(unsafe.Pointer(quickbarAbilityFlags(uint32(id), uint32(on))))
}

//export nox_xxx_netAbilityRewardCli_4611E0
func nox_xxx_netAbilityRewardCli_4611E0(id, rank C.int, notify *C.char) {
	quickbarAbilityReward(int(id), int(rank), uintptr(unsafe.Pointer(notify)))
}

//export nox_xxx_quickbarButtonBook_45F3F0
func nox_xxx_quickbarButtonBook_45F3F0() C.int { return C.int(quickbarBookTooltip()) }

//export sub_45F480
func sub_45F480(w C.int) C.int { return C.int(quickbarDirectionTooltip(bookWindow(uint32(w)))) }
