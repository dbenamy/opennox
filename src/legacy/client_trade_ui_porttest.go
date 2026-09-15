//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"image"
	"unsafe"
)

// PortTestTradeUI invokes the actual Go quantity-dialog and trade-window owners.
func PortTestTradeUI(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	w := func() *gui.Window { return uiInventoryWindowValue(uint32(a[0])) }
	point := func() image.Point { p := (*[2]int32)(unsafe.Pointer(a[0])); return image.Pt(int(p[0]), int(p[1])) }
	switch op {
	case 0:
		uiAmountToggle()
		return 0
	case 1:
		return uint32(uiAmountMouse(w(), int(a[1]), a[2], 0))
	case 2:
		return uint32(uiAmountCancel())
	case 3:
		return uint32(uiAmountInit())
	case 4:
		return uint32(uiAmountDraw(w()))
	case 5:
		return uint32(uiAmountPanel(w(), int(a[1]), a[2], a[3]))
	case 6:
		uiAmountFree()
		return 0
	case 7:
		return uint32(uiAmountShow((*uint16)(unsafe.Pointer(a[0])), int(a[1]), int(a[2]), uint32(a[3]), uint32(a[4]), unsafe.Pointer(a[5]), uint32(a[6]), uint32(a[7]), unsafe.Pointer(a[8]), unsafe.Pointer(a[9])))
	case 8:
		return uint32(uiAmountPosition(int(a[0]), int(a[1])))
	case 9:
		return uiAmountPrice(uint32(a[0]), uint32(a[1]))
	case 10:
		return uint32(uiTradeRequest(14))
	case 11:
		return uint32(uiTradeMouse(w(), int(a[1]), a[2], 0))
	case 12:
		return uint32(uiTradeRemoveRequest(uiInventoryDrawable(uint32(a[0]))))
	case 13:
		return uiInventoryPointer(unsafe.Pointer(uiTradeHit(0, point())))
	case 14:
		return uint32(uiTradePanel(w(), int(a[1]), a[2], a[3]))
	case 15:
		return uint32(uiTradeRequest(17))
	case 16:
		return uint32(uiTradeDraw())
	case 17:
		return uint32(uiTradeHover(uiInventoryPackedPoint(a[2])))
	case 18:
		return uiInventoryPointer(unsafe.Pointer(uiTradeHit(1, point())))
	case 19:
		return uint32(uiTradeDestroy())
	case 20:
		return uiTradeActive()
	case 21:
		return uint32(uiTradeShow())
	case 22:
		return uint32(uiTradeStart(unsafe.Pointer(a[0])))
	case 23:
		return uint32(uiTradeReset())
	case 24:
		return uint32(uiTradeFinish())
	case 25:
		return uiTradeRemoveCode((*uiTradeCell)(unsafe.Pointer(a[0])), uint32(a[1]))
	case 26:
		return uint32(uiTradeFindCode((*uiTradeCell)(unsafe.Pointer(a[0])), uint32(a[1])))
	case 27:
		return uiTradeAdd(unsafe.Pointer(a[0]))
	case 28:
		return uint32(uiTradeCellFits(uint32(a[0]), (*uiTradeCell)(unsafe.Pointer(a[1]))))
	case 29:
		return uiInventoryPointer(unsafe.Pointer(uiTradeFindCell(0, uint32(a[0]))))
	case 30:
		return uiInventoryPointer(unsafe.Pointer(uiTradeFindCell(1, uint32(a[0]))))
	case 31:
		return uint32(uiTradeMoney(unsafe.Pointer(a[0])))
	case 32:
		return uiTradeAcceptance(unsafe.Pointer(a[0]))
	case 33:
		return uint32(uiTradePrepare())
	case 34:
		return uint32(uiTradeInit())
	case 35:
		return uiTradeRemove(unsafe.Pointer(a[0]))
	default:
		panic("trade UI operation")
	}
}
