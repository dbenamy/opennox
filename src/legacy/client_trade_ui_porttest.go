//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_1.h"
#include "client__gui__guitrade.h"
*/
import "C"
import "unsafe"

// PortTestTradeUI invokes the original quantity-dialog and trade-window owners.
func PortTestTradeUI(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	switch op {
	case 0:
		C.sub_4BFD40()
		return 0
	case 1:
		return uint32(C.sub_4BFDD0((*C.uint32_t)(unsafe.Pointer(a[0])), C.int(a[1]), C.uint(a[2])))
	case 2:
		return uint32(C.sub_4BFE40())
	case 3:
		return uint32(C.nox_gui_itemAmount_init_4BFEF0())
	case 4:
		return uint32(C.sub_4C0030(C.int(a[0])))
	case 5:
		return uint32(C.sub_4C01C0(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 6:
		C.nox_gui_itemAmount_free_4C03E0()
		return 0
	case 7:
		return uint32(C.nox_gui_itemAmountDialog_4C0430((*C.wchar2_t)(unsafe.Pointer(a[0])), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), unsafe.Pointer(a[5]), C.int(a[6]), C.int(a[7]), unsafe.Pointer(a[8]), unsafe.Pointer(a[9])))
	case 8:
		return uint32(C.sub_4C0560(C.int(a[0]), C.int(a[1])))
	case 9:
		return uint32(C.sub_4C05F0(C.int(a[0]), C.int(a[1])))
	case 10:
		return uint32(C.nox_xxx_func_4C0610())
	case 11:
		return uint32(C.sub_4C0630(C.int(a[0]), C.uint(a[1]), C.uint(a[2])))
	case 12:
		return uint32(C.nox_xxx_clientTrade_0_4C08E0(C.int(a[0])))
	// These probes require the documented full-pointer C prerequisite first.
	case 13:
		return uint32(uintptr(unsafe.Pointer(C.sub_4C0910((*C.int2)(unsafe.Pointer(a[0]))))))
	case 14:
		return uint32(C.sub_4C0C90(C.int(a[0]), C.int(a[1]), (*C.int)(unsafe.Pointer(a[2])), C.int(a[3])))
	case 15:
		return uint32(C.nox_xxx_clientTrade_4C0CE0())
	case 16:
		return uint32(C.sub_4C0D00())
	case 17:
		return uint32(C.sub_4C1120(C.int(a[0]), C.int(a[1]), C.uint(a[2])))
	case 18:
		return uint32(uintptr(unsafe.Pointer(C.sub_4C11E0((*C.uint32_t)(unsafe.Pointer(a[0]))))))
	case 19:
		return uint32(C.nox_xxx_closeP2PTradeWnd_4C12A0())
	case 20:
		return uint32(C.sub_4C12C0())
	case 21:
		return uint32(C.nox_xxx_showP2PTradeWnd_4C12D0())
	case 22:
		return uint32(C.nox_xxx_netP2PStartTrade_4C1320(C.int(a[0])))
	case 23:
		return uint32(C.sub_4C1410())
	case 24:
		return uint32(C.sub_4C1590())
	case 25:
		return uint32(C.sub_4C1710(C.int(a[0]), C.int(a[1])))
	case 26:
		return uint32(C.sub_4C1760(C.int(a[0]), C.int(a[1])))
	case 27:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_tradeClientAddItem_4C1790(C.int(a[0])))))
	case 28:
		return uint32(C.sub_4C18E0(C.int(a[0]), (*C.uint32_t)(unsafe.Pointer(a[1]))))
	case 29:
		return uint32(uintptr(unsafe.Pointer(C.sub_4C1910(C.int(a[0])))))
	case 30:
		return uint32(uintptr(unsafe.Pointer(C.sub_4C19C0(C.int(a[0])))))
	case 31:
		return uint32(C.sub_4C1B50(C.int(a[0])))
	case 32:
		return uint32(C.sub_4C1BC0(C.int(a[0])))
	case 33:
		return uint32(C.nox_xxx_prepareP2PTrade_4C1BF0())
	case 34:
		return uint32(C.sub_4C09D0())
	case 35:
		return uint32(C.sub_4C15D0(C.int(a[0])))
	default:
		panic("trade UI operation")
	}
}
