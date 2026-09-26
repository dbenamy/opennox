//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

// PortTestMeterCall invokes the native owners directly while fixtures still own
// the actual production state.
func PortTestMeterCall(op int, w *gui.Window, a, b, c, d int) uint32 {
	switch op {
	case 0:
		return uiMeterMode()
	case 1:
		nox_xxx_cliShowHideTubes_470AA0(int(a))
		return 0
	case 2:
		return uint32(uintptr(uiMeterInitColors()))
	case 3:
		return uint32(sub_470C40(int(a)))
	case 4:
		return uint32(nox_xxx_cliSetTotalHealth_470C80(a, b))
	case 5:
		return uint32(sub_470CB0(a))
	case 6:
		return uint32(sub_470CC0())
	case 7:
		return uint32(sub_470CD0())
	case 8:
		return uint32(nox_xxx_cliSetManaAndMax_470CE0(a, b))
	case 9:
		return uint32(nox_xxx_cliSetMana_470D10(a))
	case 10:
		return uint32(sub_470D20(a, b))
	case 11:
		sub_470D70()
		return 0
	case 12:
		return uint32(sub_470D90(int(int32(a)), int(int32(b))))
	case 13:
		return uint32(nox_xxx_cliGetMana_470DD0())
	case 14:
		return uint32(uiMeterHeartbeat())
	case 15:
		return uint32(sub_470E90(int(int32(uint32(uintptr(w.C())))), int(int32(a))))
	case 16:
		nox_win_init_cur_weapon(w, int(int32(a)), int(int32(b)), int(int32(c)), int(int32(d)))
		return 0
	case 17:
		return uint32(uiMeterWeaponDraw(w))
	case 18:
		return uint32(uiMeterChargeRaster(w))
	case 19:
		return uint32(uiMeterLabel(w))
	case 20:
		return uint32(uiMeterPotionDraw(w))
	case 21:
		return uint32(nox_xxx_guiBottleSlotProc_471B90(int(int32(uint32(uintptr(w.C())))), int(int32(a))))
	case 22:
		return uint32(uiMeterMiniBar(w))
	case 23:
		return uint32(uiMeterAdvanceCharge())
	case 24:
		return uint32(uiMeterCross(a, b))
	case 25:
		return uint32(nox_xxx_guiHealthManaTubeProc_472100(int(int32(uint32(uintptr(w.C())))), int(int32(a))))
	case 26:
		return uint32(sub_4721A0(int(int32(a))))
	case 27:
		return uint32(nox_xxx_cliPrepareGameplay2_4721D0())
	case 28:
		uiMeterQuickPotion(0)
		return 0
	case 29:
		uiMeterQuickPotion(1)
		return 0
	case 30:
		uiMeterQuickPotion(2)
		return 0
	case 31:
		return uint32(uintptr(uiMeterBindings()))
	case 32:
		return uint32(uintptr(unsafe.Pointer(sub_472310())))
	case 33:
		return uint32(uiMeterWeaponTooltip())
	case 34:
		return uint32(sub_471160(w, int(int32(a)), int(int32(b)), int(int32(c)), int(int32(d))))
	case 35:
		return uint32(uiMeterInit())
	case 36:
		return uint32(uiMeterTube(w))
	default:
		panic("unknown meter operation")
	}
}
func PortTestMeterCallback(op int) unsafe.Pointer {
	switch op {
	case 15:
		return uiMeterCallbackIdentity(15)
	case 17:
		return uiMeterCallbackIdentity(17)
	case 18:
		return uiMeterCallbackIdentity(18)
	case 19:
		return uiMeterCallbackIdentity(19)
	case 20:
		return uiMeterCallbackIdentity(20)
	case 21:
		return uiMeterCallbackIdentity(21)
	case 22:
		return uiMeterCallbackIdentity(22)
	case 25:
		return uiMeterCallbackIdentity(25)
	case 33:
		return uiMeterCallbackIdentity(33)
	case 36:
		return uiMeterCallbackIdentity(36)
	default:
		return nil
	}
}

// PortTestMeterInventoryWeaponDrawCallback preserves the callback identity used
// by inventory-window capture normalization after its Go owner is called directly.
func PortTestMeterInventoryWeaponDrawCallback() unsafe.Pointer {
	return uiMeterInventoryWeaponDrawIdentity()
}
