//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "GAME3_2.h"
#include "client__gui__guiinv.h"
*/
import "C"

import (
	"math"
	"unsafe"
)

// PortTestInventoryDisplayModifierFunctions exposes identities of existing
// production callbacks; tests never execute a substitute modifier algorithm.
func PortTestInventoryDisplayModifierFunctions() [5]unsafe.Pointer {
	return [5]unsafe.Pointer{C.nox_xxx_lightngEffect_4E06F0, C.nox_xxx_fireEffect_4E0550, C.sub_4E0380, C.sub_4E0370, C.nox_xxx_effectDamageMultiplier_4E04C0}
}

// PortTestInventoryDisplay invokes the actual inventory display routines.
// Floating returns are captured as float64 bits; other returns use 386 words.
func PortTestInventoryDisplay(op int, a, b, c uintptr) uint64 {
	switch op {
	case 0:
		return uint64(uint32(C.sub_4625D0((*C.uint32_t)(unsafe.Pointer(a)))))
	case 1:
		return math.Float64bits(float64(C.sub_4626C0(C.int(a))))
	case 2:
		return math.Float64bits(float64(C.sub_462700(C.int(a))))
	case 3:
		return uint64(uint32(C.sub_4627F0((*C.uint32_t)(unsafe.Pointer(a)))))
	case 4:
		return uint64(uint32(C.sub_463370((*C.uint32_t)(unsafe.Pointer(a)), (*C.nox_point)(unsafe.Pointer(b)), (*C.uint32_t)(unsafe.Pointer(c)))))
	case 5:
		return uint64(uint32(C.sub_4633B0(C.int(a), (*C.float)(unsafe.Pointer(b)), (*C.float)(unsafe.Pointer(c)))))
	case 6:
		return uint64(uint32(C.sub_463420(C.int(a))))
	case 7:
		C.nox_client_makePlayerStatsDlg_463880((*C.int)(unsafe.Pointer(a)))
		return 0
	case 8:
		return uint64(uint32(C.nox_xxx_guiDrawInventoryTray_4643B0(C.int(a), C.int(b))))
	case 9:
		return uint64(uint32(C.sub_465D50_draw(C.int(a))))
	case 10:
		return uint64(uint32(C.nox_xxx_inventoryDrawProc_466580((*C.uint32_t)(unsafe.Pointer(a)))))
	case 11:
		return uint64(uintptr(unsafe.Pointer(C.sub_466660(C.int(a), (*C.int2)(unsafe.Pointer(b))))))
	case 12:
		return uint64(uint32(C.sub_466E20((*C.uint32_t)(unsafe.Pointer(a)))))
	case 13:
		return uint64(uint32(C.sub_466F50((*C.uint32_t)(unsafe.Pointer(a)), (*C.int)(unsafe.Pointer(b)))))
	case 14:
		return uint64(uint32(C.nox_xxx_inventoryNameSignInit_4671E0()))
	case 15:
		return uint64(uint32(C.sub_467750(C.int(a), C.char(b))))
	default:
		panic("inventory display operation")
	}
}
