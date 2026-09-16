package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "GAME1_2.h"
#include "GAME2.h"
#include "client__gui__guiinv.h"
extern uint32_t dword_5d4594_1049796_inventory_click_column_index, dword_5d4594_1049800_inventory_click_row_index;
extern uint32_t dword_5d4594_1062480, dword_5d4594_1062492;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func uiInventoryValidCell(col, row int) bool { return col >= 0 && col < 4 && row >= 0 && row < 21 }

//export sub_464B40
func sub_464B40(col, row C.int) C.int {
	return C.int(bool2int(uiInventoryValidCell(int(col), int(row))))
}
func uiInventoryPlace(dr *client.Drawable, col, row int) int {
	if !uiInventoryValidCell(col, row) {
		return 0
	}
	cell := &uiInventoryGrid()[row+21*col]
	if cell.Count != 0 && (uint32(dr.ObjClass)&0x4000000 != 0 || cell.Drawable.TypeIDVal != dr.TypeIDVal) {
		return 0
	}
	if cell.Count >= 32 {
		return 0
	}
	if cell.Count == 0 {
		cell.Drawable = GetClient().Nox_new_drawable_for_thing(int(dr.TypeIDVal))
		if cell.Drawable == nil {
			uiInventoryError("DrawablesExhausted")
			return 0
		}
		cell.Drawable.ObjFlags |= 0x40000000
		uiInventoryCopyItem(cell.Drawable, dr)
	}
	cell.Codes[cell.Count] = dr.NetCode32
	cell.Count++
	cell.Equipped = 0
	for _, v := range uiInventoryEquipment() {
		for eq := uiInventoryDrawable(v); eq != nil; eq = uiInventoryNext(eq) {
			if eq.NetCode32 == dr.NetCode32 {
				cell.Equipped = 1
				if cell.Alternate != 0 {
					uiInventorySetAlternate(nil)
					cell.Alternate = 0
				}
				return 1
			}
		}
	}
	return 1
}

//export sub_4649B0
func sub_4649B0(v, col, row C.int) C.int {
	return C.int(uiInventoryPlace(uiInventoryDrawable(uint32(v)), int(col), int(row)))
}
func uiInventorySetClick(col, row int) {
	C.dword_5d4594_1049796_inventory_click_column_index = C.uint32_t(col)
	C.dword_5d4594_1049800_inventory_click_row_index = C.uint32_t(row)
}
func uiInventoryDragCopy() {
	index := int(C.dword_5d4594_1049800_inventory_click_row_index) + 21*int(C.dword_5d4594_1049796_inventory_click_column_index)
	cell := &uiInventoryGrid()[index]
	if cell.Count == 0 {
		return
	}
	dr := GetClient().Nox_new_drawable_for_thing(int(cell.Drawable.TypeIDVal))
	uiInventorySetDragged(dr)
	if dr == nil {
		uiInventoryError("DrawablesExhausted")
		return
	}
	dr.ObjFlags |= 0x40000000
	dr.NetCode32 = cell.Codes[0]
	uiInventoryCopyItem(dr, cell.Drawable)
	uiInventoryRemoveStack(&uiInventoryLookup{Cell: cell})
}

//export nox_xxx_cliInventorySpriteUpd_465A30
func nox_xxx_cliInventorySpriteUpd_465A30() { uiInventoryDragCopy() }

//export sub_467B00
func sub_467B00(typ, quantity C.int) C.int {
	count := 0
	grid := uiInventoryGrid()
	for row := 0; row < 20; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[row+21*col]
			if cell.Count == 0 {
				count++
				continue
			}
			if cell.Drawable.TypeIDVal != uint32(typ) {
				continue
			}
			limit := int32(31)
			if uint32(cell.Drawable.ObjClass)&0x10 != 0 {
				limit = 3
				if noxflags.HasGame(6144) {
					limit = 9
				}
			}
			if uint32(cell.Drawable.ObjClass)&0x4000000 == 0 && int32(quantity)+int32(cell.Count) <= limit {
				count++
			}
		}
	}
	return C.int(count)
}
func uiInventoryAlterWeapon() {
	playerDr := uiInventoryDrawable(memmap.Uint32(0x852978, 8))
	if playerDr == nil || GetClient().Cli().Cursor != 0 || quickbarAbilityAvailable(1) != 0 {
		return
	}
	player := uiMeterPlayer()
	if player == nil || *(*uint32)(unsafe.Add(playerDr.C(), 276)) == 34 {
		return
	}
	if C.nox_xxx_pointInRect_4281F0((*C.int2)(memmap.PtrOff(0x5D4594, 1062572)), (*C.int4)(memmap.PtrOff(0x587000, 136336))) == 1 {
		Nox_xxx_cursorSetDraggedItem_477690(nil)
	}
	alt := uiInventoryCellRef(uint32(C.dword_5d4594_1062480))
	dequip := func(dr *client.Drawable) {
		C.dword_5d4594_1062492 = C.uint32_t(uiInventoryPointer(dr.C()))
		uiInventoryDequipRequest(dr)
		C.nox_xxx_clientPlaySoundSpecial_452D80(895, 100)
	}
	if alt != nil && GetServer().S().Weapons.Nox_xxx_ammoCheck_415880(int(alt.Drawable.TypeIDVal)) == 2 {
		typ := GetServer().S().Weapons.Sub_415840(2)
		if dr := uiInventoryEquippedType(uint32(typ)); dr != nil {
			dequip(dr)
		}
		return
	}
	mask := *(*uint32)(unsafe.Add(player, 4))
	for bit := uint32(1); bit < 27; bit++ {
		if bit == 1 || mask&(1<<bit) == 0 {
			continue
		}
		typ := GetServer().S().Weapons.Sub_415840(1 << bit)
		if dr := uiInventoryEquippedType(uint32(typ)); dr != nil {
			dequip(dr)
			return
		}
	}
	if alt != nil {
		alt.Drawable.NetCode32 = alt.Codes[0]
		uiInventoryEquipRequest(alt.Drawable)
		C.nox_xxx_clientPlaySoundSpecial_452D80(895, 100)
	}
}

//export nox_client_invAlterWeapon_4672C0
func nox_client_invAlterWeapon_4672C0() { uiInventoryAlterWeapon() }
