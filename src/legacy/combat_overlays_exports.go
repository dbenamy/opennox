package legacy

/*
#include "GAME2_3.h"
*/
import "C"

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

//export nox_xxx_unitSpriteCheckAlly_4951F0
func nox_xxx_unitSpriteCheckAlly_4951F0(code C.int) C.int {
	return C.int(bool2int(combatAllyLookup(uint32(code)) != nil))
}
