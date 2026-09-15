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
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

type PortTestTradeUICell struct {
	Drawable *client.Drawable
	Count    uint32
	Codes    [32]uint32
	Value    uint32
}

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
		C.sub_4BFD40,
		C.sub_4BFDD0,
		C.sub_4BFE40,
		C.nox_gui_itemAmount_init_4BFEF0,
		C.sub_4C0030,
		C.sub_4C01C0,
		C.nox_gui_itemAmount_free_4C03E0,
		C.nox_gui_itemAmountDialog_4C0430,
		C.sub_4C0560,
		C.sub_4C05F0,
		C.nox_xxx_func_4C0610,
		C.sub_4C0630,
		C.nox_xxx_clientTrade_0_4C08E0,
		C.sub_4C0910,
		C.sub_4C0C90,
		C.nox_xxx_clientTrade_4C0CE0,
		C.sub_4C0D00,
		C.sub_4C1120,
		C.sub_4C11E0,
		C.nox_xxx_closeP2PTradeWnd_4C12A0,
		C.sub_4C12C0,
		C.nox_xxx_showP2PTradeWnd_4C12D0,
		C.nox_xxx_netP2PStartTrade_4C1320,
		C.sub_4C1410,
		C.sub_4C1590,
		C.sub_4C1710,
		C.sub_4C1760,
		C.nox_xxx_tradeClientAddItem_4C1790,
		C.sub_4C18E0,
		C.sub_4C1910,
		C.sub_4C19C0,
		C.sub_4C1B50,
		C.sub_4C1BC0,
		C.nox_xxx_prepareP2PTrade_4C1BF0,
		C.sub_4C09D0,
		C.sub_4C15D0,
	}
}
