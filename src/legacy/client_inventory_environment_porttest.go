//go:build porttest

package legacy

/*
#include "defs.h"
extern uint32_t dword_5d4594_1062552;
extern uint32_t dword_5d4594_1049864;
extern uint32_t dword_5d4594_1062480;
extern uint32_t dword_5d4594_1062488;
extern nox_window* nox_win_unk5;
extern nox_window* dword_5d4594_1062452;
*/
import "C"

import "unsafe"

// PortTestUIInventoryWords borrows actual inventory state and restores it on cleanup.
func PortTestUIInventoryWords() ([]*uint32, func()) {
	words := []*uint32{
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062552)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1049864)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062480)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062488)),
		(*uint32)(unsafe.Pointer(&C.nox_win_unk5)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062452)),
	}
	old := make([]uint32, len(words))
	for i, p := range words {
		old[i] = *p
		*p = 0
	}
	return words, func() {
		for i, p := range words {
			*p = old[i]
		}
	}
}
