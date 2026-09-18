package legacy

/*
#include "GAME4_1.h"
*/
import "C"
import "unsafe"

//export sub_515C80
func sub_515C80(u C.int, value *C.uint8_t) C.int {
	return C.int(monsterControlByte(objectFromInt(u), (*byte)(unsafe.Pointer(value))))
}

//export sub_516FC0
func sub_516FC0() { monsterPendingResolve() }
