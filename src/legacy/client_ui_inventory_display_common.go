package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME2_1.h"
#include "GAME3_2.h"
#include "noxstring.h"
extern uint32_t dword_5d4594_1049844, dword_5d4594_1062456;
extern uint32_t dword_5d4594_1062476, dword_5d4594_1062480, dword_5d4594_1062484;
extern uint32_t dword_5d4594_1062512, dword_5d4594_1063116, dword_5d4594_1063120;
extern uint32_t dword_5d4594_1063636;
extern uint32_t dword_5d4594_1049796_inventory_click_column_index;
extern uint32_t dword_5d4594_1049800_inventory_click_row_index;
// ABI adapters to the existing production variadic formatter.
static int inventoryFormatInts(wchar2_t* dst, wchar2_t* fmt, int a, int b) {return nox_swprintf(dst,fmt,a,b);}
static int inventoryFormatFloat(wchar2_t* dst, wchar2_t* fmt, double a) {return nox_swprintf(dst,fmt,a);}
static int inventoryFormatText(wchar2_t* dst, wchar2_t* fmt, wchar2_t* a, wchar2_t* b) {return nox_swprintf(dst,fmt,a,b);}
static int inventoryFormatLiteral(wchar2_t* dst, wchar2_t* fmt) {return nox_swprintf(dst,fmt);}
*/
import "C"

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"golang.org/x/image/font"
)

func uiInventoryText(id string) string {
	return GetServer().S().Strings().GetStringInFile(strman.ID(id), `C:\NoxPost\src\Client\Gui\guiinv.c`)
}
func uiInventoryFormatInts(id string, a, b int) string {
	var dst [256]uint16
	C.inventoryFormatInts((*C.wchar2_t)(unsafe.Pointer(&dst[0])), (*C.wchar2_t)(unsafe.Pointer(internWStr(uiInventoryText(id)))), C.int(a), C.int(b))
	return alloc.GoString16(&dst[0])
}
func uiInventoryFormatFloat(id string, a float64) string {
	var dst [256]uint16
	C.inventoryFormatFloat((*C.wchar2_t)(unsafe.Pointer(&dst[0])), (*C.wchar2_t)(unsafe.Pointer(internWStr(uiInventoryText(id)))), C.double(a))
	return alloc.GoString16(&dst[0])
}
func uiInventoryFormatLiteral(text string) string {
	var dst [256]uint16
	C.inventoryFormatLiteral((*C.wchar2_t)(unsafe.Pointer(&dst[0])), (*C.wchar2_t)(unsafe.Pointer(internWStr(text))))
	return alloc.GoString16(&dst[0])
}
func uiInventorySmallFont() font.Face {
	return GetClient().R2().GetFonts().AsFont(unsafe.Pointer(uintptr(C.dword_5d4594_1063636)))
}
func uiInventoryMainWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1062456)))
}
func uiInventoryIdentifyWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(C.dword_5d4594_1062476)))
}
func uiInventoryAlternate() *uiInventoryCell {
	return (*uiInventoryCell)(unsafe.Pointer(uintptr(C.dword_5d4594_1062480)))
}
func uiInventorySelectedItem() *client.Drawable {
	return uiInventoryDrawable(uint32(C.dword_5d4594_1063116))
}
func uiInventoryItemModifier(dr *client.Drawable, index int) *server.ModifierEff {
	return *(**server.ModifierEff)(unsafe.Add(dr.C(), 432+uintptr(index*4)))
}
func uiInventoryElementValue(dr *client.Drawable, fire bool) float64 {
	if dr == nil || uint32(dr.Class())&0x13001000 == 0 {
		return 0
	}
	var fn unsafe.Pointer = C.nox_xxx_lightngEffect_4E06F0
	if fire {
		fn = C.nox_xxx_fireEffect_4E0550
	}
	for i := 2; i < 4; i++ {
		if m := uiInventoryItemModifier(dr, i); m != nil && m.AttackPreHit52.Fnc == fn {
			return float64(m.AttackPreHit52.Valf)
		}
	}
	return 0
}
func uiInventoryScaledDurability(dr *client.Drawable, current, maximum *float32) uint32 {
	*current = float32(*(*uint16)(unsafe.Add(dr.C(), 292)))
	*maximum = float32(*(*uint16)(unsafe.Add(dr.C(), 294)))
	ret := uiInventoryPointer(dr.C())
	if uint32(dr.Class())&0x13001000 != 0 {
		if m := uiInventoryItemModifier(dr, 1); m != nil && m.Defend76.Fnc == C.sub_4E0380 {
			*current = float32(int32(*current * m.Defend76.Valf))
			*maximum = float32(int32(*maximum * m.Defend76.Valf))
			ret = uint32(int32(*maximum))
		}
	}
	return ret
}

//export sub_463370
func sub_463370(w *C.uint32_t, pos *C.nox_point, out *C.uint32_t) C.int {
	p := image.Pt(int(pos.x), int(pos.y)).Sub(uiWindowPosition((*gui.Window)(unsafe.Pointer(w))))
	v := unsafe.Slice((*int32)(unsafe.Pointer(out)), 2)
	v[0], v[1] = int32(p.X), int32(p.Y)
	return C.int(p.Y)
}

//export sub_463420
func sub_463420(v C.int) C.int { *memmap.PtrUint32(0x5D4594, 1050012) = uint32(v); return v }

func uiInventoryHitRect(p image.Point, offset uintptr) bool {
	r := unsafe.Slice((*int32)(memmap.PtrOff(0x587000, offset)), 4)
	return p.X >= int(r[0]) && p.X <= int(r[2]) && p.Y >= int(r[1]) && p.Y <= int(r[3])
}
func uiInventoryHoverText(p image.Point) *uint16 {
	if uiInventoryHitRect(p, 136336) {
		if memmap.Uint8(0x5D4594, 1049870) == 1 {
			return nil
		}
		slot := uiInventoryEquipmentAt(p)
		if slot == -1 {
			return alloc.InternCString16(uiInventoryText("DollRegionError"))
		}
		if dr := uiInventoryDrawable(uiInventoryEquipment()[slot]); dr != nil {
			return uiItemTooltip(dr)
		}
		return alloc.InternCString16(uiInventoryText("ToolTipDrag"))
	}
	if memmap.Uint8(0x5D4594, 1049869) != 0 {
		return nil
	}
	var col, row int
	if uiInventoryHitRect(p, 136368) {
		col, row = (p.Y-13)/50, 20
	} else if uiInventoryHitRect(p, 136352) {
		col = (p.X - 314) / 50
		row = (p.Y + int(C.dword_5d4594_1062512) - 13) / 50
	} else {
		return nil
	}
	C.dword_5d4594_1049796_inventory_click_column_index = C.uint32_t(col)
	C.dword_5d4594_1049800_inventory_click_row_index = C.uint32_t(row)
	if col >= 0 && col < 4 && row >= 0 && row < 21 {
		cell := &uiInventoryGrid()[col*21+row]
		if cell.Count != 0 {
			cell.Drawable.NetCode32 = cell.Codes[0]
			return uiItemTooltip(cell.Drawable)
		}
	}
	return nil
}

//export sub_466660
func sub_466660(_ C.int, p *C.int2) *C.wchar2_t {
	return (*C.wchar2_t)(unsafe.Pointer(uiInventoryHoverText(image.Pt(int(p.field_0), int(p.field_4)))))
}

//export sub_466E20
func sub_466E20(w *C.uint32_t) C.int {
	var key string
	switch uint32(*w) {
	case 9105:
		key = "JournalModeTT"
	case 9106:
		key = "InventoryModeTT"
	case 9107:
		key = "StatsModeTT"
	case 9108:
		key = "PaperDollModeTT"
	case 9111:
		key = "CloseInventoryTT"
	default:
		return 0
	}
	uiCursorTooltip(alloc.InternCString16(uiInventoryText(key)))
	return 1
}

//export nox_xxx_inventoryNameSignInit_4671E0
func nox_xxx_inventoryNameSignInit_4671E0() C.int {
	dst := (*C.wchar2_t)(memmap.PtrOff(0x5D4594, 1062588))
	C.nox_wcscpy(dst, (*C.wchar2_t)(memmap.PtrOff(0x5D4594, 1063676)))
	p := uiMeterPlayer()
	level := 0
	if noxflags.HasGame(4096) || C.nox_xxx_isQuest_4D6F50() != 0 || C.sub_4D6F70() != 0 {
		level = int(min(uint32(C.dword_5d4594_1049844), 10))
	} else if p != nil {
		level = int(*(*int8)(unsafe.Add(p, 3684)))
	}
	if p == nil {
		return C.int(level)
	}
	class := *(*byte)(unsafe.Add(p, 2251))
	className := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, 29456+uintptr(class)*4)))
	title := uiInventoryText(fmt.Sprintf("experience:%s%d", className, level))
	return C.int(C.inventoryFormatText(dst, (*C.wchar2_t)(unsafe.Pointer(internWStr(uiInventoryText("ElaborateNameFormat")))), (*C.wchar2_t)(unsafe.Add(p, 4704)), (*C.wchar2_t)(unsafe.Pointer(internWStr(title)))))
}

//export sub_467750
func sub_467750(code C.int, status C.char) C.int {
	if code != 0 {
		if found := uiInventoryFindCode(uint32(code)); found != nil {
			if old := uiInventoryAlternate(); old != nil {
				old.Alternate = 0
			}
			C.dword_5d4594_1062480 = C.uint32_t(uiInventoryPointer(unsafe.Pointer(found.Cell)))
			found.Cell.Alternate = 1
			return 1
		}
	} else if old := uiInventoryAlternate(); old != nil {
		old.Alternate = 0
		C.dword_5d4594_1062480 = 0
	}
	if status != 0 {
		if status != 1 {
			return 0
		}
		Nox_xxx_printCentered_445490(uiInventoryText("Weapon2CantUse"))
		if C.dword_5d4594_1062484 == 0 {
			return 0
		}
		if found := uiInventoryFindCode(uint32(C.dword_5d4594_1062484)); found != nil {
			uiInventorySetAlternate(found.Cell)
			return 0
		}
	}
	C.dword_5d4594_1062484 = 0
	return 0
}
