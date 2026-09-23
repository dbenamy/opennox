//go:build porttest

package legacy

/*
#include "defs.h"


*/
import "C"

import "unsafe"

// PortTestMeterInventory borrows the actual client inventory storage. The extra
// twenty-first row is part of the C search ABI, even though UI slots use 20 rows.
func PortTestMeterInventory() (grid []byte, equipment []uint32, restore func()) {
	if C.sizeof_nox_inventory_cell_t != 148 || C.NOX_INVENTORY_CELLS_MAX != 84 {
		panic("client inventory ABI")
	}
	grid = unsafe.Slice((*byte)(unsafe.Pointer(&legacyGlobals.nox_client_inventory_grid_1050020[0])), 148*84)
	equipment = unsafe.Slice((*uint32)(unsafe.Pointer(&legacyGlobals.array_5D4594_1049872[0])), 9)
	oldGrid, oldEquipment := append([]byte(nil), grid...), append([]uint32(nil), equipment...)
	clear(grid)
	clear(equipment)
	return grid, equipment, func() { copy(grid, oldGrid); copy(equipment, oldEquipment) }
}
