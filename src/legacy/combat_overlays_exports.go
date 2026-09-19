package legacy

/*
#include "GAME2_3.h"
*/
import "C"
import "unsafe"

//export sub_495060
func sub_495060(code C.int, current, maximum C.short) C.int {
	return C.int(bool2int(combatAllyAdd(uint32(code), uint16(current), uint16(maximum))))
}

//export sub_4950C0
func sub_4950C0(code C.int) C.int { return C.int(bool2int(combatAllyRemove(uint32(code)))) }

//export sub_4950F0
func sub_4950F0(code C.int, flag C.char) C.int {
	return C.int(bool2int(combatAllyFlag(uint32(code), byte(flag))))
}

//export sub_495120
func sub_495120(code C.int, current, maximum C.short) C.int {
	return C.int(bool2int(combatAllyPair(uint32(code), uint16(current), uint16(maximum))))
}

//export sub_495150
func sub_495150(code C.int, current C.short) C.int {
	return C.int(bool2int(combatAllyFirst(uint32(code), uint16(current))))
}

//export nox_xxx_unitSpriteCheckAlly_4951F0
func nox_xxx_unitSpriteCheckAlly_4951F0(code C.int) C.int {
	return C.int(bool2int(combatAllyLookup(uint32(code)) != nil))
}

//export sub_495210
func sub_495210(p C.int) C.int { combatFeedAdd(unsafe.Pointer(uintptr(p))); return 1 }

//export sub_4959B0
func sub_4959B0() { combatFriendClear() }

//export nox_xxx_cliAddObjFriend_4959F0
func nox_xxx_cliAddObjFriend_4959F0(code C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(combatFriendAdd(uint32(code))))
}

//export sub_495A20
func sub_495A20(code C.int) { combatFriendRemove(uint32(code)) }

//export nox_xxx_cliAddHealthChange_49A650
func nox_xxx_cliAddHealthChange_49A650(code C.int, amount C.short) *C.uint16_t {
	return (*C.uint16_t)(unsafe.Pointer(combatHealthAdd(uint32(code), int16(amount))))
}
