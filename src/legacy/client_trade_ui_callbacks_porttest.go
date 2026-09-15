//go:build porttest

package legacy

/*
#include "defs.h"
void portTestTradeAccept(int2*, uint32_t, uint32_t, uint32_t, uint32_t);
void portTestTradeCancel(int2*, uint32_t, uint32_t, uint32_t, uint32_t);
*/
import "C"
import "unsafe"

var portTestTradeObserve func([7]uint32)

func PortTestTradeUIObserve(fn func([7]uint32)) (unsafe.Pointer, unsafe.Pointer, func()) {
	old := portTestTradeObserve
	portTestTradeObserve = fn
	return C.portTestTradeAccept, C.portTestTradeCancel, func() { portTestTradeObserve = old }
}

//export portTestTradeAccept
func portTestTradeAccept(pos *C.int2, code, typ, count, extra C.uint32_t) {
	portTestTradeObserve([7]uint32{0, uint32(pos.field_0), uint32(pos.field_4), uint32(code), uint32(typ), uint32(count), uint32(extra)})
}

//export portTestTradeCancel
func portTestTradeCancel(pos *C.int2, code, typ, count, extra C.uint32_t) {
	portTestTradeObserve([7]uint32{1, uint32(pos.field_0), uint32(pos.field_4), uint32(code), uint32(typ), uint32(count), uint32(extra)})
}
