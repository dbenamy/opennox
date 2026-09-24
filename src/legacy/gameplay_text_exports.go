package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
typedef const char gameplay_text_const_char;
typedef const nox_object_t gameplay_text_const_object;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func nox_xxx_netInformTextMsg_4DA0F0(to, kind C.int, data *C.int) C.int {
	return C.int(gameplayTextInformation(int(to), int(kind), unsafe.Pointer(data)))
}

func nox_xxx_netInformTextMsg2_4DA180(kind C.int, data *C.uint8_t) C.int {
	return C.int(gameplayTextInformationAll(int(kind), unsafe.Pointer(data)))
}

func nox_xxx_netPriMsgToPlayer_4DA2C0(u *C.nox_object_t, text *C.char, flag C.char) {
	gameplayTextPrivate((*server.Object)(unsafe.Pointer(u)), (*byte)(unsafe.Pointer(text)), byte(flag))
}

//export nox_xxx_netPrintLineToAll_4DA390
func nox_xxx_netPrintLineToAll_4DA390(text *C.gameplay_text_const_char) C.int {
	return C.int(gameplayTextPrivateAll((*byte)(unsafe.Pointer(text))))
}
