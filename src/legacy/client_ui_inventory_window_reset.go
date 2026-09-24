package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
)

func uiInventoryResetPanelControls() int {
	*memmap.PtrUint8(0x5D4594, 1049869) = 0
	dword_5d4594_1062512 = dword_5d4594_1062516
	slider := uiInventoryWindowValue(uint32(dword_5d4594_1062508))
	uiInventorySliderValue(slider, 16395, 0, 850)
	uiInventorySliderValue(slider, 16394, (*gui.SliderData)(slider.WidgetData).Max-uint32(dword_5d4594_1062512), 0)
	nox_xxx_wndSetIcon_46AE60(C.int(dword_5d4594_1062528), 0)
	sub_46AEC0(C.int(dword_5d4594_1062528), C.int(dword_5d4594_1049976))
	uiInventoryWindowValue(uint32(dword_5d4594_1062528)).SetID(9105)
	*memmap.PtrUint8(0x5D4594, 1049870) = 0
	nox_xxx_wndSetIcon_46AE60(C.int(dword_5d4594_1062524), C.int(dword_5d4594_1049992))
	sub_46AEC0(C.int(dword_5d4594_1062524), C.int(dword_5d4594_1049996))
	uiInventoryWindowValue(uint32(dword_5d4594_1062524)).SetID(9107)
	return nox_window_set_hidden((*nox_window)(uiInventoryWindowValue(uint32(dword_5d4594_1062468)).C()), 0)
}
func uiInventoryResetWindow() int {
	grid := uiInventoryGrid()
	for row := 0; row < 21; row++ {
		for col := 0; col < 4; col++ {
			cell := &grid[col*21+row]
			if cell.Drawable != nil {
				GetClient().Nox_xxx_spriteDelete_45A4B0(cell.Drawable)
				cell.Drawable = nil
			}
			cell.Count = 0
			cell.Equipped = 0
			cell.Alternate = 0
		}
	}
	uiInventoryCloseIdentify()
	dword_5d4594_1049864 = 0
	uiInventorySetAlternate(nil)
	dword_5d4594_1062488 = 0
	clear(uiInventoryEquipment())
	dword_5d4594_1062492, dword_5d4594_1062496 = 0, 0
	*memmap.PtrUint8(0x5D4594, 1062536) = 0
	for _, off := range []uintptr{1062540, 1062544, 1062548} {
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	dword_5d4594_1062552 = 0
	uiMeterRefreshPotions()
	dword_587000_136184 = uint32(^uint32(224))
	*memmap.PtrUint8(0x5D4594, 1049868) = 0
	dword_5d4594_1062516, dword_5d4594_1062520, dword_5d4594_1062512 = 0, 0, 0
	return uiInventoryResetPanelControls()
}
