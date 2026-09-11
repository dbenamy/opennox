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

//export nox_xxx_servShopStart_50EF10_trade
func nox_xxx_servShopStart_50EF10_trade(left, right C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(tradeStart(objectFromInt(left), objectFromInt(right))))
}

//export nox_xxx_tradeP2PAddOffer2_50F820_trade
func nox_xxx_tradeP2PAddOffer2_50F820_trade(s, u C.int, bits C.float) C.int {
	item := (*server.Object)(unsafe.Pointer(uintptr(math.Float32bits(float32(bits)))))
	return C.int(tradeAddOffer(shopSessionFromInt(s), objectFromInt(u), item))
}

//export sub_5100C0_trade
func sub_5100C0_trade(u C.int, s *C.uint32_t, code C.int) {
	tradeBuy(objectFromInt(u), (*shopSession)(unsafe.Pointer(s)), uint32(code))
}

//export sub_510640_trade
func sub_510640_trade(u, s, typ C.int, count *C.float) *C.float {
	return (*C.float)(unsafe.Pointer(uintptr(tradeBuyMany(objectFromInt(u), shopSessionFromInt(s), int32(typ), uint32(uintptr(unsafe.Pointer(count)))))))
}

//export sub_5109C0_trade
func sub_5109C0_trade(u *C.int, s C.int, code *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(tradeSaleQuote((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(s), uint32(uintptr(unsafe.Pointer(code)))))))
}

//export sub_510BE0_trade
func sub_510BE0_trade(u *C.int, s C.int, code *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(tradeSell((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(s), uint32(uintptr(unsafe.Pointer(code)))))))
}
