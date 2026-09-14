//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
static uint32_t gameplayTextInvoke(int op, uint32_t* a) {
 switch(op) {
 case 0: return nox_xxx_netSendLineMessage_4D9EB0(a[0], (wchar2_t*)a[1], a[2], a[3]);
 case 1: return nox_xxx_printToAll_4D9FD0(a[0], (wchar2_t*)a[1], a[2], a[3]);
 default: return 0xDEADBEEF;
 }
}
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func gameplayTextInvoke(op int, args [5]uint32) uint32 {
	ptr := func(i int) unsafe.Pointer { return unsafe.Pointer(uintptr(args[i])) }
	obj := func(i int) *server.Object { return (*server.Object)(ptr(i)) }
	switch op {
	case 0, 1:
		return uint32(C.gameplayTextInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
	case 2:
		return uint32(gameplayTextInformation(int(args[0]), int(args[1]), ptr(2)))
	case 3:
		return uint32(gameplayTextInformationAll(int(args[0]), ptr(1)))
	case 4:
		gameplayTextPrivate(obj(0), (*byte)(ptr(1)), byte(args[2]))
		return 0
	case 5:
		return uint32(gameplayTextPrivateAll((*byte)(ptr(0))))
	case 6:
		return uint32(uintptr(GetServer().S().Players.FirstUnit().CObj()))
	case 7:
		return uint32(uintptr(GetServer().S().Players.NextUnit(obj(0)).CObj()))
	case 8:
		return uint32(bool2int(gameplayTextByteEncoding(gameplayTextUnits((*uint16)(ptr(0))))))
	case 9:
		return uint32(gameplayTextChat(obj(0), gameplayTextUnits((*uint16)(ptr(1))), uint16(args[2])))
	default:
		panic("text operation")
	}
}
