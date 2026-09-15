//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guiinv.h"
int sub_467CA0(void);
*/
import "C"

import "unsafe"

// PortTestInventoryWindow invokes original inventory lifecycle and input owners.
func PortTestInventoryWindow(op int, a, b, c, d uintptr) uint32 {
	switch op {
	case 0:
		return uint32(C.sub_462740())
	case 1:
		return uint32(C.sub_464BD0(C.int(a), C.int(b), C.uint(c)))
	case 2:
		return uint32(C.sub_466160())
	case 3:
		return uint32(C.sub_4661D0())
	case 4:
		return uint32(C.nox_xxx_inventroryOnHovewerSub_4667E0(C.int(a), C.int(b), C.uint(c)))
	case 5:
		return uint32(C.nox_xxx_inventoryDrawAllMB_463430(C.int(a)))
	case 6:
		return uint32(C.sub_464770(C.int(a), C.int(b), C.uint(c)))
	case 7:
		return uint32(C.nox_xxx_XorEaxEaxSub_464BA0())
	case 8:
		return uint32(C.nox_xxx_inventoryWndProc_464BB0(C.int(a), C.int(b)))
	case 9:
		return uint32(C.nox_xxx_clientTradeMB_4657E0((*C.uint32_t)(unsafe.Pointer(a))))
	case 10:
		C.sub_4658A0(C.int(a), (*C.int2)(unsafe.Pointer(b)))
		return 0
	case 11:
		return uint32(C.sub_465990((*C.uint32_t)(unsafe.Pointer(a))))
	case 12:
		return uint32(C.sub_465CA0())
	case 13:
		C.sub_465CD0((*C.uint32_t)(unsafe.Pointer(a)), C.int(b), C.int(c), C.int(d))
		return 0
	case 14:
		return uint32(C.sub_465DE0(C.int(a)))
	case 15:
		return uint32(C.nox_xxx_wndCreateInventoryMB_465E00())
	case 16:
		return uint32(C.nox_xxx_movEax1Sub_4661C0())
	case 17:
		return uint32(C.sub_466220(C.int(a), C.int(b), (*C.int)(unsafe.Pointer(c)), C.int(d)))
	case 18:
		return uint32(C.sub_466550(C.int(a), C.uint(b)))
	case 19:
		return uint32(C.sub_466620(C.int(a), C.int(b), C.uint(c)))
	case 20:
		return uint32(C.sub_466950(C.int(a)))
	case 21:
		return uint32(C.sub_466BA0((*C.uint32_t)(unsafe.Pointer(a)), C.int(b), C.uint(c), C.int(d)))
	case 22:
		return uint32(C.sub_466BF0(C.int(a), C.int(b), C.uint(c), C.int(d)))
	case 23:
		return uint32(C.sub_466C40(C.int(a)))
	case 24:
		return uint32(C.sub_466ED0(C.int(a)))
	case 25:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_inventoryLoadImages_467050())))
	case 26:
		return uint32(C.sub_467650())
	case 27:
		return uint32(C.sub_467980())
	case 28:
		return uint32(C.sub_467BB0())
	case 29:
		return uint32(C.sub_467C10())
	case 30:
		return uint32(C.nox_client_toggleInventory_467C60())
	case 31:
		return uint32(C.sub_467C80())
	case 32:
		return uint32(C.sub_467CA0())
	case 33:
		return uint32(C.sub_467CD0())
	default:
		panic("inventory window operation")
	}
}
