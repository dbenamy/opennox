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

func nox_xxx_servShopStart_50EF10_trade(left, right int32) *uint32 {
	return (*uint32)(unsafe.Pointer(tradeStart(objectFromInt(C.int(left)), objectFromInt(C.int(right)))))
}

func nox_xxx_tradeP2PAddOffer2_50F820_trade(s, u int32, bits float32) int32 {
	item := (*server.Object)(unsafe.Pointer(uintptr(math.Float32bits(float32(bits)))))
	return int32(tradeAddOffer(shopSessionFromInt(C.int(s)), objectFromInt(C.int(u)), item))
}

func sub_5100C0_trade(u int32, s *uint32, code int32) {
	tradeBuy(objectFromInt(C.int(u)), (*shopSession)(unsafe.Pointer(s)), uint32(code))
}

func sub_510640_trade(u, s, typ int32, count *float32) *float32 {
	return (*float32)(unsafe.Pointer(uintptr(tradeBuyMany(objectFromInt(C.int(u)), shopSessionFromInt(C.int(s)), int32(typ), uint32(uintptr(unsafe.Pointer(count)))))))
}

func sub_5109C0_trade(u *int32, s int32, code *uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(tradeSaleQuote((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(C.int(s)), uint32(uintptr(unsafe.Pointer(code)))))))
}

func sub_510BE0_trade(u *int32, s int32, code *uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(tradeSell((*server.Object)(unsafe.Pointer(u)), shopSessionFromInt(C.int(s)), uint32(uintptr(unsafe.Pointer(code)))))))
}
