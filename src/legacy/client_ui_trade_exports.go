package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func sub_4BFE40() C.int { return C.int(uiAmountCancel()) }

//export nox_gui_itemAmountDialog_4C0430
func nox_gui_itemAmountDialog_4C0430(title *C.wchar2_t, x, y, code, typ C.int, mods unsafe.Pointer, maximum, extra C.int, accept, cancel unsafe.Pointer) C.int {
	return C.int(uiAmountShow((*uint16)(unsafe.Pointer(title)), int(x), int(y), uint32(code), uint32(typ), mods, uint32(maximum), uint32(extra), accept, cancel))
}

//export sub_4C05F0
func sub_4C05F0(enabled, unit C.int) C.int {
	return C.int(uiAmountPrice(uint32(enabled), uint32(unit)))
}

//export sub_4C1120
func sub_4C1120(w, event C.int, packed C.uint) C.int {
	return C.int(uiTradeHover(uiInventoryPackedPoint(uintptr(packed))))
}

func nox_xxx_netP2PStartTrade_4C1320(data C.int) C.int {
	return C.int(uiTradeStart(unsafe.Pointer(uintptr(uint32(data)))))
}

func sub_4C1590() C.int { return C.int(uiTradeFinish()) }

func nox_xxx_tradeClientAddItem_4C1790(data C.int) C.uint32_t {
	return C.uint32_t(uiTradeAdd(unsafe.Pointer(uintptr(uint32(data)))))
}

func sub_4C1B50(data C.int) C.int { return C.int(uiTradeMoney(unsafe.Pointer(uintptr(uint32(data))))) }

func sub_4C1BC0(data C.int) C.int {
	return C.int(uiTradeAcceptance(unsafe.Pointer(uintptr(uint32(data)))))
}

func nox_xxx_prepareP2PTrade_4C1BF0() C.int { return C.int(uiTradePrepare()) }

func sub_4C15D0(data C.int) C.int { return C.int(uiTradeRemove(unsafe.Pointer(uintptr(uint32(data))))) }
