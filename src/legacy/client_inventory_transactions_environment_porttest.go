//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_1062484;
extern uint32_t dword_5d4594_1062492;
extern uint32_t dword_5d4594_1062496;
extern uint32_t dword_5d4594_1062512;
extern uint32_t dword_5d4594_1062516;
extern uint32_t dword_5d4594_1062556;
extern uint32_t dword_5d4594_1062560;
extern uint32_t dword_5d4594_1062564;
extern uint32_t dword_5d4594_1049796_inventory_click_column_index;
extern uint32_t dword_5d4594_1049800_inventory_click_row_index;
extern uint32_t dword_5d4594_1049856;
extern uint32_t dword_5d4594_825736;
*/
import "C"
import "unsafe"

// PortTestInventoryTransactionWords borrows real inventory transaction globals.
func PortTestInventoryTransactionWords() ([]*uint32, func()) {
	words := []*uint32{
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062484)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062492)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062496)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062512)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062516)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062556)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062560)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1062564)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1049796_inventory_click_column_index)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1049800_inventory_click_row_index)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_1049856)),
		(*uint32)(unsafe.Pointer(&C.dword_5d4594_825736)),
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
