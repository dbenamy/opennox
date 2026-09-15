//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guimeter.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

// PortTestMeterCall enters the production C ABI; fixtures own all inputs.
func PortTestMeterCall(op int, w *gui.Window, a, b, c, d int) uint32 {
	switch op {
	case 0:
		return uint32(C.nox_xxx_playerGet_470A90())
	case 1:
		C.nox_xxx_cliShowHideTubes_470AA0(C.int(a))
		return 0
	case 2:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_guiHealthManaColorInit_470B00())))
	case 3:
		return uint32(C.sub_470C40(C.int(a)))
	case 4:
		return uint32(C.nox_xxx_cliSetTotalHealth_470C80(C.int(a), C.int(b)))
	case 5:
		return uint32(C.sub_470CB0(C.int(a)))
	case 6:
		return uint32(C.sub_470CC0())
	case 7:
		return uint32(C.sub_470CD0())
	case 8:
		return uint32(C.nox_xxx_cliSetManaAndMax_470CE0(C.int(a), C.int(b)))
	case 9:
		return uint32(C.nox_xxx_cliSetMana_470D10(C.int(a)))
	case 10:
		return uint32(C.sub_470D20(C.int(a), C.int(b)))
	case 11:
		C.sub_470D70()
		return 0
	case 12:
		return uint32(C.sub_470D90(C.int(a), C.int(b)))
	case 13:
		return uint32(C.nox_xxx_cliGetMana_470DD0())
	case 14:
		return uint32(C.sub_470DE0())
	case 15:
		return uint32(C.sub_470E90(C.int(uintptr(w.C())), C.int(a)))
	case 16:
		C.nox_win_init_cur_weapon((*C.nox_window)(unsafe.Pointer(w.C())), C.int(a), C.int(b), C.int(c), C.int(d))
		return 0
	case 17:
		return uint32(C.sub_470F40_draw((*C.nox_window)(unsafe.Pointer(w.C()))))
	case 18:
		return uint32(C.sub_471250((*C.uint32_t)(unsafe.Pointer(w.C()))))
	case 19:
		return uint32(C.sub_471450((*C.uint32_t)(unsafe.Pointer(w.C()))))
	case 20:
		return uint32(C.nox_xxx_guiBottleSlotDrawFn_471A80((*C.uint32_t)(unsafe.Pointer(w.C()))))
	case 21:
		return uint32(C.nox_xxx_guiBottleSlotProc_471B90(C.int(uintptr(w.C())), C.int(a)))
	case 22:
		return uint32(C.nox_xxx_drawHealthManaBar_471C00(C.int(uintptr(w.C()))))
	case 23:
		return uint32(C.sub_472080())
	case 24:
		return uint32(C.sub_4720C0(C.int(a), C.int(b)))
	case 25:
		return uint32(C.nox_xxx_guiHealthManaTubeProc_472100(C.int(uintptr(w.C())), C.int(a)))
	case 26:
		return uint32(C.sub_4721A0(C.int(a)))
	case 27:
		return uint32(C.nox_xxx_cliPrepareGameplay2_4721D0())
	case 28:
		C.nox_client_quickHealthPotion_472220()
		return 0
	case 29:
		C.nox_client_quickManaPotion_472240()
		return 0
	case 30:
		C.nox_client_quickCurePoisonPotion_472260()
		return 0
	case 31:
		return uint32(uintptr(unsafe.Pointer(C.sub_472280())))
	case 32:
		return uint32(uintptr(unsafe.Pointer(C.sub_472310())))
	case 33:
		return uint32(C.sub_4710B0())
	case 34:
		return uint32(C.sub_471160(C.int(uintptr(w.C())), C.int(a), C.int(b), C.int(c), C.int(d)))
	case 35:
		return uint32(C.nox_xxx_guiHealthManaInit_4714E0())
	case 36:
		return uint32(C.nox_xxx_guiHealthManaTubeDraw_471D10(C.int(uintptr(w.C()))))
	default:
		panic("unknown meter operation")
	}
}
func PortTestMeterCallback(op int) unsafe.Pointer {
	switch op {
	case 15:
		return C.sub_470E90
	case 17:
		return C.sub_470F40_draw
	case 18:
		return C.sub_471250
	case 19:
		return C.sub_471450
	case 20:
		return C.nox_xxx_guiBottleSlotDrawFn_471A80
	case 21:
		return C.nox_xxx_guiBottleSlotProc_471B90
	case 22:
		return C.nox_xxx_drawHealthManaBar_471C00
	case 25:
		return C.nox_xxx_guiHealthManaTubeProc_472100
	case 33:
		return C.sub_4710B0
	case 36:
		return C.nox_xxx_guiHealthManaTubeDraw_471D10
	default:
		return nil
	}
}
