//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_1.h"
#include "client__gui__guitrade.h"
extern uint32_t dword_5d4594_1320932;
extern uint32_t dword_5d4594_1320936;
extern uint32_t dword_5d4594_1320940;
extern uint32_t dword_5d4594_1320944;
extern uint32_t dword_5d4594_1320948;
extern uint32_t dword_5d4594_1320968;
extern uint32_t dword_5d4594_1320972;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

type PortTestTradeUICell = uiTradeCell

func PortTestTradeUICells() [2][]PortTestTradeUICell {
	if unsafe.Sizeof(PortTestTradeUICell{}) != 140 || unsafe.Offsetof(PortTestTradeUICell{}.Value) != 136 {
		panic("trade cell layout")
	}
	return [2][]PortTestTradeUICell{
		unsafe.Slice((*PortTestTradeUICell)(memmap.PtrOff(0x5D4594, 1319284)), 4),
		unsafe.Slice((*PortTestTradeUICell)(memmap.PtrOff(0x5D4594, 1320308)), 4),
	}
}
func PortTestTradeUIWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1320932": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320932)),
		"dword_5d4594_1320936": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320936)),
		"dword_5d4594_1320940": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320940)),
		"dword_5d4594_1320944": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320944)),
		"dword_5d4594_1320948": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320948)),
		"dword_5d4594_1320968": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320968)),
		"dword_5d4594_1320972": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1320972)),
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
func PortTestTradeUICallbacks() []unsafe.Pointer {
	return []unsafe.Pointer{
		nil,
		nil,
		C.sub_4BFE40,
		nil,
		nil,
		nil,
		nil,
		C.nox_gui_itemAmountDialog_4C0430,
		nil,
		C.sub_4C05F0,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		C.sub_4C1120,
		nil,
		nil,
		nil,
		nil,
		C.nox_xxx_netP2PStartTrade_4C1320,
		nil,
		C.sub_4C1590,
		nil,
		nil,
		C.nox_xxx_tradeClientAddItem_4C1790,
		nil,
		nil,
		nil,
		C.sub_4C1B50,
		C.sub_4C1BC0,
		C.nox_xxx_prepareP2PTrade_4C1BF0,
		nil,
		C.sub_4C15D0,
	}
}
