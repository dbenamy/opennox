//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guiinv.h"
*/
import "C"

import (
	"image"
	"unsafe"
)

// PortTestInventoryWindow invokes inventory lifecycle owners and retained C interfaces.
func PortTestInventoryWindow(op int, a, b, c, d uintptr) uint32 {
	switch op {
	case 0:
		return uint32(C.sub_462740())
	case 1:
		return uint32(uiInventoryMainEvents(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 2:
		return uint32(C.sub_466160())
	case 3:
		return uint32(C.sub_4661D0())
	case 4:
		return uint32(C.nox_xxx_inventroryOnHovewerSub_4667E0(C.int(a), C.int(b), C.uint(c)))
	case 5:
		return uint32(uiInventoryDrawWindow(uiInventoryWindowValue(uint32(a))))
	case 6:
		return uint32(uiInventoryAlternateEvents(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 7:
		return 0
	case 8:
		return uint32(uiInventoryWindowAdmission(uiInventoryWindowValue(uint32(a)), int(b), 0, 0))
	case 9:
		return uint32(uiInventoryTradeClick(portInventoryWindowPoint(a)))
	case 10:
		uiInventoryStartDrag(uiInventoryWindowValue(uint32(a)), portInventoryWindowPoint(b))
		return 0
	case 11:
		return uint32(uiInventoryEquipmentAt(portInventoryWindowPoint(a)))
	case 12:
		return uint32(uiInventoryOpenIdentify())
	case 13:
		C.sub_465CD0((*C.uint32_t)(unsafe.Pointer(a)), C.int(b), C.int(c), C.int(d))
		return 0
	case 14:
		return uint32(sub_465DE0(C.int(a)))
	case 15:
		return uint32(uiInventoryCreateWindow())
	case 16:
		return 1
	case 17:
		return uint32(uiInventoryPanelEvents(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 18:
		return uint32(uiInventoryToggleButton(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 19:
		return uint32(C.sub_466620(C.int(a), C.int(b), C.uint(c)))
	case 20:
		return uint32(uiInventoryNewScrollControls(uiInventoryWindowValue(uint32(a))))
	case 21:
		return uint32(uiInventoryThumbEvents(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 22:
		return uint32(uiInventoryTrackEvents(uiInventoryWindowValue(uint32(a)), int(b), c, d))
	case 23:
		return uint32(uiInventoryNewModeControls(uiInventoryWindowValue(uint32(a))))
	case 24:
		return uint32(uiInventoryNewIdentifyWindow(uiInventoryWindowValue(uint32(a))))
	case 25:
		return uint32(uiInventoryLoadImages())
	case 26:
		return uint32(C.sub_467650())
	case 27:
		return uint32(uiInventoryResetWindow())
	case 28:
		return uint32(C.sub_467BB0())
	case 29:
		return uint32(C.sub_467C10())
	case 30:
		return uint32(uiInventoryToggleWindow())
	case 31:
		return uint32(C.sub_467C80())
	case 32:
		return uint32(uiInventoryResetClosedScroll())
	case 33:
		return uint32(uiInventoryCancelDrag())
	default:
		panic("inventory window operation")
	}
}

func portInventoryWindowPoint(p uintptr) image.Point {
	v := (*[2]int32)(unsafe.Pointer(p))
	return image.Pt(int(v[0]), int(v[1]))
}
