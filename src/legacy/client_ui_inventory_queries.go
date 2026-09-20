package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
extern nox_inventory_cell_t nox_client_inventory_grid_1050020[NOX_INVENTORY_CELLS_MAX];
extern uint32_t array_5D4594_1049872[9];
extern uint32_t dword_5d4594_1062480;
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
)

type uiInventoryCell struct {
	Drawable                      *client.Drawable
	Codes                         [32]uint32
	Equipped, Alternate           uint32
	Count, Flags1, Flags2, Flags3 uint8
	Tail                          uint32
}
type uiInventoryLookup struct {
	Cell  *uiInventoryCell
	Index uint32
}

func uiInventoryGrid() []uiInventoryCell {
	return unsafe.Slice((*uiInventoryCell)(unsafe.Pointer(&C.nox_client_inventory_grid_1050020[0])), 84)
}
func uiInventoryEquipment() []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(&C.array_5D4594_1049872[0])), 9)
}
func uiInventoryDrawable(v uint32) *client.Drawable {
	return (*client.Drawable)(unsafe.Pointer(uintptr(v)))
}
func uiInventoryNext(dr *client.Drawable) *client.Drawable {
	return *(**client.Drawable)(unsafe.Add(dr.C(), 368))
}
func uiInventoryFindType(typ uint32) *uiInventoryCell {
	grid := uiInventoryGrid()
	for row := 0; row < 21; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			if cell.Count != 0 && cell.Drawable.TypeIDVal == typ {
				return cell
			}
		}
	}
	return nil
}
func uiInventoryFindCode(code uint32) *uiInventoryLookup {
	grid := uiInventoryGrid()
	for row := 0; row < 21; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			for index := 0; index < min(int(cell.Count), len(cell.Codes)); index++ {
				if cell.Codes[index] == code {
					out := (*uiInventoryLookup)(memmap.PtrOff(0x5D4594, 1049788))
					out.Cell = cell
					out.Index = uint32(index)
					return out
				}
			}
		}
	}
	return nil
}
func uiInventoryEquippedType(typ uint32) *client.Drawable {
	for _, v := range uiInventoryEquipment() {
		for dr := uiInventoryDrawable(v); dr != nil; dr = uiInventoryNext(dr) {
			if dr.TypeIDVal == typ {
				return dr
			}
		}
	}
	return nil
}
func uiInventoryCurrentWeapon() *client.Drawable {
	bow := memmap.PtrUint32(0x5D4594, 1063640)
	if *bow == 0 {
		*bow = uint32(GetClient().Cli().Things.IndByID("Bow"))
	}
	equipment := uiInventoryEquipment()
	for dr := uiInventoryDrawable(equipment[8]); dr != nil; dr = uiInventoryNext(dr) {
		if dr.TypeIDVal == *bow {
			return uiInventoryDrawable(equipment[8])
		}
	}
	return uiInventoryDrawable(equipment[7])
}
func uiInventoryDragged() *client.Drawable {
	return uiInventoryDrawable(memmap.Uint32(0x5D4594, 1049848))
}
func uiInventoryItem(code uint32) *client.Drawable {
	if found := uiInventoryFindCode(code); found != nil {
		return found.Cell.Drawable
	}
	if dr := uiInventoryDragged(); dr != nil && dr.NetCode32 == code {
		return dr
	}
	return nil
}
func uiInventoryTypeCount(typ uint32) int {
	if cell := uiInventoryFindType(typ); cell != nil {
		return int(cell.Count)
	}
	return 0
}
func uiInventorySelectedWeapon() *client.Drawable {
	player := uiMeterPlayer()
	if player == nil {
		return nil
	}
	mask := *(*uint32)(unsafe.Add(player, 4))
	for i := uint32(2); i < 27; i++ {
		if mask&(1<<i) != 0 {
			typ := GetServer().S().Weapons.Sub_415840(1 << i)
			if dr := uiInventoryEquippedType(uint32(typ)); dr != nil {
				if found := uiInventoryFindCode(dr.NetCode32); found != nil {
					return found.Cell.Drawable
				}
				return nil
			}
		}
	}
	return nil
}

//export sub_4615C0
func sub_4615C0() int { return int(uintptr(unsafe.Pointer(uiInventoryCurrentWeapon()))) }

//export sub_461600
func sub_461600(typ int) int {
	return int(uintptr(unsafe.Pointer(uiInventoryEquippedType(uint32(typ)))))
}

//export sub_461930
func sub_461930() int {
	for _, v := range uiInventoryEquipment() {
		for dr := uiInventoryDrawable(v); dr != nil; dr = uiInventoryNext(dr) {
			if uint32(dr.ObjClass)&0x1001000 != 0 {
				return 1
			}
		}
	}
	return 0
}

//export sub_461EF0
func sub_461EF0(code int) *C.char {
	return (*C.char)(unsafe.Pointer(uiInventoryFindCode(uint32(code))))
}

func sub_4676D0(code int) int { return int(uintptr(unsafe.Pointer(uiInventoryItem(uint32(code))))) }

//export sub_467700
func sub_467700(code int) int {
	if found := uiInventoryFindCode(uint32(code)); found != nil {
		return int(found.Cell.Count)
	}
	if dr := uiInventoryDragged(); dr != nil && dr.NetCode32 == uint32(code) {
		return 1
	}
	return 0
}

//export sub_467810
func sub_467810(col, row int) int {
	if col < 0 || col >= 4 || row < 0 || row >= 20 {
		return 0
	}
	return int(uiInventoryGrid()[row+21*col].Count)
}

//export sub_467870
func sub_467870(col, row int) *C.char {
	if col < 0 || col >= 4 || row < 0 || row >= 20 {
		return nil
	}
	return (*C.char)(unsafe.Pointer(&uiInventoryGrid()[row+21*col].Codes[0]))
}

//export sub_4678B0
func sub_4678B0() int {
	if cell := (*uiInventoryCell)(unsafe.Pointer(uintptr(C.dword_5d4594_1062480))); cell != nil {
		return int(cell.Codes[0])
	}
	return 0
}
