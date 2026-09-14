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

//export nox_gameplayTextLine
func nox_gameplayTextLine(u C.int, text *C.ushort) C.int {
	return C.int(gameplayTextLine((*server.Object)(unsafe.Pointer(uintptr(uint32(u)))), gameplayTextUnits((*uint16)(unsafe.Pointer(text)))))
}

//export nox_gameplayTextAll
func nox_gameplayTextAll(flags C.char, text *C.ushort) C.int {
	return C.int(gameplayTextAll(byte(flags), gameplayTextUnits((*uint16)(unsafe.Pointer(text)))))
}

//export nox_xxx_netInformTextMsg_4DA0F0
func nox_xxx_netInformTextMsg_4DA0F0(to, kind C.int, data *C.int) C.int {
	return C.int(gameplayTextInformation(int(to), int(kind), unsafe.Pointer(data)))
}

//export nox_xxx_netInformTextMsg2_4DA180
func nox_xxx_netInformTextMsg2_4DA180(kind C.int, data *C.uint8_t) C.int {
	return C.int(gameplayTextInformationAll(int(kind), unsafe.Pointer(data)))
}

//export nox_xxx_netPriMsgToPlayer_4DA2C0
func nox_xxx_netPriMsgToPlayer_4DA2C0(u *C.nox_object_t, text *C.gameplay_text_const_char, flag C.char) {
	gameplayTextPrivate((*server.Object)(unsafe.Pointer(u)), (*byte)(unsafe.Pointer(text)), byte(flag))
}

//export nox_xxx_netPrintLineToAll_4DA390
func nox_xxx_netPrintLineToAll_4DA390(text *C.gameplay_text_const_char) C.int {
	return C.int(gameplayTextPrivateAll((*byte)(unsafe.Pointer(text))))
}

//export nox_xxx_getFirstPlayerUnit_4DA7C0
func nox_xxx_getFirstPlayerUnit_4DA7C0() *C.nox_object_t {
	return (*C.nox_object_t)(GetServer().S().Players.FirstUnit().CObj())
}

//export nox_xxx_getNextPlayerUnit_4DA7F0
func nox_xxx_getNextPlayerUnit_4DA7F0(u *C.gameplay_text_const_object) *C.nox_object_t {
	return (*C.nox_object_t)(GetServer().S().Players.NextUnit((*server.Object)(unsafe.Pointer(u))).CObj())
}

//export nox_xxx_cliCanTalkMB_4100F0
func nox_xxx_cliCanTalkMB_4100F0(text *C.short) C.int {
	return C.int(bool2int(gameplayTextByteEncoding(gameplayTextUnits((*uint16)(unsafe.Pointer(text))))))
}
