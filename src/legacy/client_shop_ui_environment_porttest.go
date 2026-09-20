//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_2.h"
#include "client__gui__guishop.h"
void sub_479680(void);
extern uint32_t dword_5d4594_1098456;
extern uint32_t dword_5d4594_1098576;
extern uint32_t dword_5d4594_1098580;
extern uint32_t dword_5d4594_1098592;
extern uint32_t dword_5d4594_1098596;
extern uint32_t dword_5d4594_1098600;
extern uint32_t dword_5d4594_1098604;
extern uint32_t dword_5d4594_1098616;
extern uint32_t dword_5d4594_1098620;
*/
import "C"
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
		"dword_5d4594_1098456": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098456)),
		"dword_5d4594_1098576": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098576)),
		"dword_5d4594_1098580": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098580)),
		"dword_5d4594_1098592": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098592)),
		"dword_5d4594_1098596": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098596)),
		"dword_5d4594_1098600": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098600)),
		"dword_5d4594_1098604": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098604)),
		"dword_5d4594_1098616": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098616)),
		"dword_5d4594_1098620": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1098620)),
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
		C.sub_478030,
		C.sub_478040,
		nil,
		nil,
		nil,
		nil,
		nil,
		C.sub_478850,
		nil,
		nil,
		nil,
		C.sub_478E50,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		C.sub_479590,
		C.sub_4795A0,
		C.sub_479690,
		nil,
		nil,
		C.sub_479810,
		C.sub_479820,
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
		C.sub_479680,
		nil,
		nil,
	}
}
