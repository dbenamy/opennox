package legacy

/*
#include <stdint.h>
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

func shopSessionFromInt(v C.int) *shopSession {
	return (*shopSession)(unsafe.Pointer(uintptr(uint32(v))))
}

//export nox_xxx_shopGetItemCost_50E3D0
func nox_xxx_shopGetItemCost_50E3D0(mode, session C.int, bits C.float) C.int {
	u := (*server.Object)(unsafe.Pointer(uintptr(math.Float32bits(float32(bits)))))
	return C.int(shopPrice(int(mode), shopSessionFromInt(session), u))
}

//export sub_50F3A0
func sub_50F3A0(s *C.uint32_t) { shopCancelTrade((*shopSession)(unsafe.Pointer(s))) }

//export nox_xxx_shopExit_50F4C0
func nox_xxx_shopExit_50F4C0(s *C.uint32_t) { shopExit((*shopSession)(unsafe.Pointer(s))) }

//export nox_xxx_tradeAccept_50F5A0
func nox_xxx_tradeAccept_50F5A0(s, u C.int) { shopAccept(shopSessionFromInt(s), objectFromInt(u)) }

//export nox_xxx_tradeP2PAddOfferMB_50FE20
func nox_xxx_tradeP2PAddOfferMB_50FE20(s C.int, code C.int) C.int {
	return C.int(shopWithdraw(shopSessionFromInt(s), uint32(code)))
}

//export sub_5108D0
func sub_5108D0(u, s, code C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(shopRepairQuote(objectFromInt(u), shopSessionFromInt(s), uint32(code)))))
}

//export sub_510AE0
func sub_510AE0(u *C.int, s C.int, code *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(shopRepair((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(s), uint32(uintptr(unsafe.Pointer(code)))))))
}

//export sub_510D10
func sub_510D10(u *C.int, s, typ C.int, count C.uint) {
	shopSell((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(s), int32(typ), uint32(count))
}

//export nox_xxx_shopCancelSession_510DC0
func nox_xxx_shopCancelSession_510DC0(s unsafe.Pointer) { shopCancel((*shopSession)(s)) }

//export sub_510DE0
func sub_510DE0(u, code C.int) C.int { return C.int(shopLookup(objectFromInt(u), uint32(code))) }

//export sub_510E20
func sub_510E20(ind C.int) { shopPlayerCleanup(int(ind)) }
