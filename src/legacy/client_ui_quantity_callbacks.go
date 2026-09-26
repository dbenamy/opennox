package legacy

import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

const (
	uiAmountDialogExport = iota
	uiAmountDrop
	uiAmountShopBuy
	uiAmountShopSell
	uiAmountShopSellCancel
	uiAmountShopRepair
	uiAmountShopRepairCancel
	uiAmountShopActive
	uiAmountShopTooltip
	uiAmountTradeTooltip
)

var uiAmountIdentitySlots [10]byte

func uiAmountNativeKey(id int) unsafe.Pointer { return unsafe.Pointer(&uiAmountIdentitySlots[id]) }

func init() {
	gui.RegisterTooltipCallbackGo(uiAmountNativeKey(uiAmountShopTooltip), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiShopHover(uiInventoryPackedPoint(arg))
	})
	gui.RegisterTooltipCallbackGo(uiAmountNativeKey(uiAmountTradeTooltip), func(_ *gui.Window, _ *gui.WindowData, arg uintptr) {
		uiTradeHover(uiInventoryPackedPoint(arg))
	})
}

// uiAmountCall retains the arbitrary foreign callback path. Native keys share the
// original five-word call frame; handlers intentionally ignore unused words.
func uiAmountCall(fn unsafe.Pointer, a1, a2, a3, a4, a5 uintptr) {
	switch fn {
	case uiAmountNativeKey(uiAmountDrop):
		p := (*[2]int32)(unsafe.Pointer(a1))
		uiInventoryDropQuantity(image.Pt(int(p[0]), int(p[1])), uint32(a2), uint32(a3), int(int32(a4)))
	case uiAmountNativeKey(uiAmountShopBuy):
		uiShopBuyAccept(uint16(a2), uint32(a3), uint32(a4))
	case uiAmountNativeKey(uiAmountShopSell):
		uiShopSellAccept(uint16(a2), uint16(a3), uint32(a4))
	case uiAmountNativeKey(uiAmountShopSellCancel):
		*uiShopWord(1098616) = 0
	case uiAmountNativeKey(uiAmountShopRepair):
		uiShopRepairAccept(uint16(a2))
	case uiAmountNativeKey(uiAmountShopRepairCancel):
		*uiShopWord(1098620) = 0
	default:
		ccall.CallVoidUPtr5(fn, a1, a2, a3, a4, a5)
	}
}
