package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
extern uint32_t dword_5d4594_1062552, dword_5d4594_1049864, dword_5d4594_1062488;
extern nox_window* nox_win_unk5;
extern nox_window* dword_5d4594_1062452;
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
)

//export sub_4673F0
func sub_4673F0(a, b int) int {
	*memmap.PtrUint32(0x5D4594, 1062580) = uint32(a)
	*memmap.PtrUint32(0x5D4594, 1062584) = uint32(b)
	return a
}

func sub_467410(v int) int { *memmap.PtrUint32(0x5D4594, 1062540) = uint32(v); return v }

func sub_467420(v C.char) C.char { *memmap.PtrUint8(0x5D4594, 1062536) = byte(v); return v }

func sub_467430() C.uchar { return C.uchar(memmap.Uint8(0x5D4594, 1062536)) }

func sub_467440(v int) int { *memmap.PtrUint32(0x5D4594, 1062544) = uint32(v); return v }

func sub_467450(v int) int { *memmap.PtrUint32(0x5D4594, 1062548) = uint32(v); return v }

func sub_467470(index int, v float32) int {
	i := uint8(index)
	*memmap.PtrFloat32(0x5D4594, 1063100+uintptr(i)*4) = v
	return int(i)
}

func sub_467490(v int) int { C.dword_5d4594_1062552 = C.uint32_t(v); return v }

//export sub_4674A0
func sub_4674A0() int { return int(C.dword_5d4594_1062552) }

//export nox_window_set_visible_unk5
func nox_window_set_visible_unk5(v int) {
	uiMeterHide((*gui.Window)(unsafe.Pointer(C.nox_win_unk5)), v == 0)
}
func uiInventoryUsePotion(typ uint32) {
	if memmap.Uint32(0x5D4594, 1096672) == 1 || noxflags.HasGame(noxflags.GamePause) {
		return
	}
	if cell := uiInventoryFindType(typ); cell != nil {
		cell.Drawable.NetCode32 = cell.Codes[0]
		uiInventoryUse(cell.Drawable)
	}
}

//export sub_467590
func sub_467590() int {
	if p := uiMeterPlayer(); p != nil {
		return int(*(*int8)(unsafe.Add(p, 3684)))
	}
	return 1
}
func uiInventoryMode() int { return int(C.dword_5d4594_1049864) }

func uiInventoryItemHealth(code int, current, maximum int16) int16 {
	if found := uiInventoryFindCode(uint32(code)); found != nil {
		dr := found.Cell.Drawable
		*(*uint16)(unsafe.Add(dr.C(), 292)) = uint16(current)
		*(*uint16)(unsafe.Add(dr.C(), 294)) = uint16(maximum)
		return int16(uintptr(dr.C()))
	}
	dr := uiInventoryDragged()
	if dr != nil && dr.NetCode32 == uint32(code) {
		*(*uint16)(unsafe.Add(dr.C(), 292)) = uint16(current)
		*(*uint16)(unsafe.Add(dr.C(), 294)) = uint16(maximum)
		return maximum
	}
	return int16(uintptr(unsafe.Pointer(dr)))
}

//export sub_467680
func sub_467680() {
	if C.dword_5d4594_1049864 == 6 {
		C.dword_5d4594_1049864 = 0
	}
}

//export nox_xxx_wndGetHandle_4676A0
func nox_xxx_wndGetHandle_4676A0() *C.nox_window { return C.dword_5d4594_1062452 }

//export sub_467740
func sub_467740(v int) int { C.dword_5d4594_1062488 = C.uint32_t(v); return v }

//export sub_4678C0
func sub_4678C0() int { return int(C.dword_5d4594_1062488) }

func sub_467930(code, current, maximum int) *C.char {
	if code == 0 {
		return nil
	}
	found := uiInventoryFindCode(uint32(code))
	if found == nil {
		return nil
	}
	cell := found.Cell
	dr := cell.Drawable
	*(*uint16)(unsafe.Add(dr.C(), 448)) = uint16(current)
	*(*uint16)(unsafe.Add(dr.C(), 450)) = uint16(maximum)
	if cell.Equipped == 1 {
		return (*C.char)(unsafe.Pointer(uintptr(uint32(sub_470D90(current, maximum)))))
	}
	return (*C.char)(unsafe.Pointer(cell))
}
