//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"unsafe"
)

// PortTestShopUI dispatches to production Go, including private helpers.
func PortTestShopUI(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	ptr := func(i int) unsafe.Pointer { return unsafe.Pointer(a[i]) }
	point := func(i int) image.Point { p := (*[2]int32)(ptr(i)); return image.Pt(int(p[0]), int(p[1])) }
	switch op {
	case 0:
		return uiShopActive()
	case 1:
		return uint32(uiShopCancelRequest())
	case 2:
		return uiInventoryPointer(uiShopDrawable(uint32(a[0])).C())
	case 3:
		return uiInventoryPointer(unsafe.Pointer(uiShopFind(uint32(a[0]))))
	case 4:
		return uint32(uiShopInit())
	case 5:
		return uint32(uiShopPanel(uiInventoryWindowValue(uint32(a[0])), int(a[1]), a[2], a[3]))
	case 6:
		return uint32(uiShopMouse(uiInventoryWindowValue(uint32(a[0])), int(a[1]), a[2], 0))
	case 7:
		uiShopBuyAccept(uint16(a[1]), uint32(a[2]), uint32(a[3]))
	case 8:
		return uint32(uiShopDraw())
	case 9:
		return uint32(uiShopDrawText(point(0), 1))
	case 10:
		return uint32(uiShopDrawStock())
	case 11:
		return uint32(uiShopHover(uiInventoryPackedPoint(a[2])))
	case 12:
		return uiShopClear()
	case 13:
		return uint32(uiShopDestroy())
	case 14:
		return uiShopPicture(uint32(a[0]))
	case 15:
		uiShopClose()
	case 16:
		return uiShopAdd(uint32(a[0]), uint32(a[1]), uint32(a[2]), uint16(a[3]), ptr(4))
	case 17:
		return uiInventoryPointer(unsafe.Pointer(uiShopFindType(uint32(a[0]))))
	case 18:
		return uiInventoryPointer(unsafe.Pointer(uiShopEmpty()))
	case 19:
		return uiShopRemove(uint32(a[0]))
	case 20:
		return uiTradeRemoveCode((*uiShopCell)(ptr(0)), uint32(a[1]))
	case 21:
		return uiShopMode()
	case 22:
		uiShopSetMode(uint32(a[0]))
	case 23:
		return uint32(uiShopSellAccept(uint16(a[1]), uint16(a[2]), uint32(a[3])))
	case 24:
		return uint32(uiInventoryTrade(24, uint16(a[0])))
	case 25:
		return uint32(uiShopMultiple(25, uint16(a[0]), byte(a[1])))
	case 26:
		*uiShopWord(1098620) = 0
	case 27:
		return uint32(uiShopRepairAccept(uint16(a[1])))
	case 28:
		return uint32(uiInventoryTrade(26, uint16(a[0])))
	case 29:
		return uint32(bool2int(uiShopMode() == 2))
	case 30:
		return uint32(bool2int(uiShopInside(point(0))))
	case 31:
		return uiInventoryPointer(uiShopHit(point(0)).C())
	case 32:
		uiShopBuyShow(point(0))
	case 33:
		uiShopBuySingle(uint32(a[0]), uint16(a[1]))
	case 34:
		uiShopBuyMultiple(uint32(a[0]), uint32(a[1]))
	case 35:
		return uint32(uiShopDrawText(point(0), 3))
	case 36:
		return uint32(uiShopDrawText(point(0), 4))
	case 37:
		return uint32(uiShopStart((*uint16)(ptr(0)), alloc.GoString((*byte)(ptr(1))), uint32(a[2])))
	case 38:
		uiShopGoldWarning(uint32(a[0]))
	case 39:
		*uiShopWord(1098616) = 0
	case 40:
		return uint32(uiShopSellShow(uint32(a[0]), uint32(a[1])))
	case 41:
		uiShopRepairShow(uint32(a[0]), uint32(a[1]))
	}
	return 0
}
