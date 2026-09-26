//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// Shop stock shares the original 140-byte drawable/count/codes/price layout.
type PortTestShopUICell = uiTradeCell

func PortTestShopUICells() []PortTestShopUICell {
	return unsafe.Slice((*PortTestShopUICell)(memmap.PtrOff(0x5D4594, 1098636)), 60)
}
func PortTestShopUIWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1098456": (*uint32)(unsafe.Pointer(&dword_5d4594_1098456)),
		"dword_5d4594_1098576": (*uint32)(unsafe.Pointer(&dword_5d4594_1098576)),
		"dword_5d4594_1098580": (*uint32)(unsafe.Pointer(&dword_5d4594_1098580)),
		"dword_5d4594_1098592": (*uint32)(unsafe.Pointer(&dword_5d4594_1098592)),
		"dword_5d4594_1098596": (*uint32)(unsafe.Pointer(&dword_5d4594_1098596)),
		"dword_5d4594_1098600": (*uint32)(unsafe.Pointer(&dword_5d4594_1098600)),
		"dword_5d4594_1098604": (*uint32)(unsafe.Pointer(&dword_5d4594_1098604)),
		"dword_5d4594_1098616": (*uint32)(unsafe.Pointer(&dword_5d4594_1098616)),
		"dword_5d4594_1098620": (*uint32)(unsafe.Pointer(&dword_5d4594_1098620)),
	}
	old := make(map[string]uint32, len(words))
	for n, p := range words {
		old[n] = *p
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}
func PortTestShopUICallbacks() []unsafe.Pointer {
	return []unsafe.Pointer{
		uiAmountNativeKey(uiAmountShopActive),
		clientUICallbackKey(clientUICallbackID_sub_478040),
		nil,
		nil,
		nil,
		nil,
		nil,
		uiAmountNativeKey(uiAmountShopBuy),
		nil,
		nil,
		nil,
		uiAmountNativeKey(uiAmountShopTooltip),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		clientUICallbackKey(clientUICallbackID_sub_479590),
		clientUICallbackKey(clientUICallbackID_sub_4795A0),
		uiAmountNativeKey(uiAmountShopSell),
		nil,
		nil,
		uiAmountNativeKey(uiAmountShopRepairCancel),
		uiAmountNativeKey(uiAmountShopRepair),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		uiAmountNativeKey(uiAmountShopSellCancel),
		nil,
		nil,
	}
}
