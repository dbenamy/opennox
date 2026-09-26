//go:build porttest

package legacy

import (
	"image"
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

// PortTestInventoryDisplayModifierFunctions exposes identities of existing
// production callbacks; tests never execute a substitute modifier algorithm.
func PortTestInventoryDisplayModifierFunctions() [5]unsafe.Pointer {
	return [5]unsafe.Pointer{modifierKey(modifierIDLightningEffect), modifierKey(modifierIDFireEffect), modifierKey(modifierIDDurabilityMultiplierEffect), modifierKey(modifierIDArmorMultiplierEffect), modifierKey(modifierIDDamageMultiplierEffect)}
}

// PortTestInventoryDisplay invokes the actual inventory display routines.
// Floating returns are captured as float64 bits; other returns use 386 words.
func PortTestInventoryDisplay(op int, a, b, c uintptr) uint64 {
	switch op {
	case 0:
		return uint64(uint32(sub_4625D0((*gui.Window)(unsafe.Pointer(a)))))
	case 1:
		return math.Float64bits(uiInventoryElementValue(uiInventoryDrawable(uint32(a)), false))
	case 2:
		return math.Float64bits(uiInventoryElementValue(uiInventoryDrawable(uint32(a)), true))
	case 3:
		return uint64(uint32(sub_4627F0(unsafe.Pointer(a))))
	case 4:
		return uint64(uint32(sub_463370(unsafe.Pointer(a), unsafe.Pointer(b), unsafe.Pointer(c))))
	case 5:
		return uint64(uiInventoryScaledDurability(uiInventoryDrawable(uint32(a)), (*float32)(unsafe.Pointer(b)), (*float32)(unsafe.Pointer(c))))
	case 6:
		return uint64(uint32(sub_463420(int32(a))))
	case 7:
		nox_client_makePlayerStatsDlg_463880(unsafe.Pointer(a))
		return 0
	case 8:
		return uint64(uint32(nox_xxx_guiDrawInventoryTray_4643B0(int32(a), int32(b))))
	case 9:
		return uint64(uint32(uiInventoryCurrentWeaponDraw((*gui.Window)(unsafe.Pointer(uintptr(uint32(a)))))))
	case 10:
		return uint64(uint32(nox_xxx_inventoryDrawProc_466580((*gui.Window)(unsafe.Pointer(a)))))
	case 11:
		p := (*[2]int32)(unsafe.Pointer(b))
		return uint64(uintptr(unsafe.Pointer(uiInventoryHoverText(image.Pt(int(p[0]), int(p[1]))))))
	case 12:
		return uint64(uint32(sub_466E20((*uint32)(unsafe.Pointer(a)))))
	case 13:
		return uint64(uint32(sub_466F50((*gui.Window)(unsafe.Pointer(a)), (*gui.WindowData)(unsafe.Pointer(b)))))
	case 14:
		return uint64(uint32(nox_xxx_inventoryNameSignInit_4671E0()))
	case 15:
		return uint64(uint32(sub_467750(int32(a), int8(b))))
	default:
		panic("inventory display operation")
	}
}
