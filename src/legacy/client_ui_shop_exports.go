package legacy

/*
#include "defs.h"
*/
import "C"

//export sub_478030
func sub_478030() C.int { return C.int(uiShopActive()) }

//export sub_478850
func sub_478850(point C.int, code C.short, typ, count C.int) {
	uiShopBuyAccept(uint16(code), uint32(typ), uint32(count))
}

//export sub_478E50
func sub_478E50(w, event C.int, packed C.uint) C.int {
	return C.int(uiShopHover(uiInventoryPackedPoint(uintptr(packed))))
}

//export sub_479690
func sub_479690(point C.int, code, typ C.short, count C.int) C.int {
	return C.int(uiShopSellAccept(uint16(code), uint16(typ), uint32(count)))
}

//export sub_479680
func sub_479680() { *uiShopWord(1098616) = 0 }

//export sub_479810
func sub_479810() { *uiShopWord(1098620) = 0 }

//export sub_479820
func sub_479820(point C.int, code C.short) C.int { return C.int(uiShopRepairAccept(uint16(code))) }
