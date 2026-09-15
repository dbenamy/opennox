//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
*/
import "C"

import (
	"math"
	"unsafe"
)

// PortTestUIInventoryCall exercises the production client inventory ABI.
func PortTestUIInventoryCall(op int, a, b, c uint32) uint32 {
	switch op {
	case 0:
		return uint32(C.sub_4615C0())
	case 1:
		return uint32(C.sub_461600(C.int(a)))
	case 2:
		return uint32(C.sub_461930())
	case 3:
		return uint32(uintptr(unsafe.Pointer(C.sub_461EF0(C.int(a)))))
	case 4:
		return uint32(C.sub_4673F0(C.int(a), C.int(b)))
	case 5:
		return uint32(C.sub_467410(C.int(a)))
	case 6:
		return uint32(C.sub_467420(C.char(a)))
	case 7:
		return uint32(C.sub_467430())
	case 8:
		return uint32(C.sub_467440(C.int(a)))
	case 9:
		return uint32(C.sub_467450(C.int(a)))
	case 10:
		return uint32(C.sub_467470(C.int(a), C.float(math.Float32frombits(b))))
	case 11:
		return uint32(C.sub_467490(C.int(a)))
	case 12:
		return uint32(C.sub_4674A0())
	case 13:
		C.nox_window_set_visible_unk5(C.int(a))
		return 0
	case 14:
		uiInventoryUsePotion(a)
		return 0
	case 15:
		return uint32(uintptr(unsafe.Pointer(uiInventoryFindType(a))))
	case 16:
		return uint32(C.sub_467590())
	case 17:
		return uint32(uiInventoryMode())
	case 18:
		return uint32(C.sub_4675E0(C.int(a), C.short(b), C.short(c)))
	case 19:
		C.sub_467680()
		return 0
	case 20:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_wndGetHandle_4676A0())))
	case 21:
		return uint32(C.sub_4676D0(C.int(a)))
	case 22:
		return uint32(C.sub_467700(C.int(a)))
	case 23:
		return uint32(C.sub_467740(C.int(a)))
	case 24:
		return uint32(C.sub_467810(C.int(a), C.int(b)))
	case 25:
		return uint32(uiInventoryTypeCount(a))
	case 26:
		return uint32(uintptr(unsafe.Pointer(C.sub_467870(C.int(a), C.int(b)))))
	case 27:
		return uint32(C.sub_4678B0())
	case 28:
		return uint32(C.sub_4678C0())
	case 29:
		return uint32(uintptr(unsafe.Pointer(uiInventorySelectedWeapon())))
	case 30:
		return uint32(uintptr(unsafe.Pointer(C.sub_467930(C.int(a), C.int(b), C.int(c)))))
	default:
		panic("client inventory operation")
	}
}
