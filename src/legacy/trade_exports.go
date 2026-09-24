package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

func nox_xxx_servShopStart_50EF10_trade(left, right int32) *uint32 {
	return (*uint32)(unsafe.Pointer(tradeStart(objectFromWord(uint32(left)), objectFromWord(uint32(right)))))
}

func nox_xxx_tradeP2PAddOffer2_50F820_trade(s, u int32, bits float32) int32 {
	item := (*server.Object)(unsafe.Pointer(uintptr(math.Float32bits(float32(bits)))))
	return int32(tradeAddOffer(shopSessionFromWord(uint32(s)), objectFromWord(uint32(u)), item))
}

func sub_5100C0_trade(u int32, s *uint32, code int32) {
	tradeBuy(objectFromWord(uint32(u)), (*shopSession)(unsafe.Pointer(s)), uint32(code))
}

func sub_510640_trade(u, s, typ int32, count *float32) *float32 {
	return (*float32)(unsafe.Pointer(uintptr(tradeBuyMany(objectFromWord(uint32(u)), shopSessionFromWord(uint32(s)), int32(typ), uint32(uintptr(unsafe.Pointer(count)))))))
}

func sub_5109C0_trade(u *int32, s int32, code *uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(tradeSaleQuote((*server.Object)(unsafe.Pointer(u)), shopSessionFromWord(uint32(s)), uint32(uintptr(unsafe.Pointer(code)))))))
}

func sub_510BE0_trade(u *int32, s int32, code *uint32) *uint32 {
	return (*uint32)(unsafe.Pointer(uintptr(tradeSell((*server.Object)(unsafe.Pointer(u)), shopSessionFromWord(uint32(s)), uint32(uintptr(unsafe.Pointer(code)))))))
}
