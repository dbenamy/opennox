//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guiinv.h"
*/
import "C"

import "unsafe"

type PortTestInventoryCell = uiInventoryCell
type PortTestInventoryLookup = uiInventoryLookup

func PortTestInventoryCells() []PortTestInventoryCell { return uiInventoryGrid() }

// PortTestInventoryTransaction enters the actual production inventory routines.
func PortTestInventoryTransaction(op int, a, b, c, d uintptr) uint32 {
	switch op {
	case 0:
		return uint32(nox_xxx_clientSetAltWeapon_461550(int32(a)))
	case 1:
		return uint32(nox_xxx_send2ServInvenFail_461630(int16(a)))
	case 2:
		return uint32(nox_xxx_spritePickup_461660(int32(a), int32(b), unsafe.Pointer(c)))
	case 3:
		return uint32(uiInventoryNewStack(uint32(a), uint32(b), unsafe.Pointer(c), (*[2]int32)(unsafe.Pointer(d))))
	case 4:
		return uint32(uintptr(unsafe.Pointer(uiInventoryAppend(uint32(a), uint32(b)))))
	case 5:
		return uint32(uiInventoryClearAlternateFlags())
	case 6:
		sub_461A80(int32(a))
		return 0
	case 7:
		return uint32(uintptr(unsafe.Pointer(sub_461B50())))
	case 8:
		return uint32(uiInventoryRemoveStack((*uiInventoryLookup)(unsafe.Pointer(a))))
	case 9:
		return uint32(uintptr(unsafe.Pointer(uiInventoryUnlink(uint32(a)))))
	case 10:
		sub_462040(int32(a))
		return 0
	case 11:
		return uint32(uiInventoryEquipmentSlot(uiInventoryDrawable(uint32(a))))
	case 12:
		return uint32(nox_xxx_clientEquip_4623B0(int32(a)))
	case 13:
		return uint32(uintptr(unsafe.Pointer(uiInventoryInsertEquipment(uiInventoryDrawable(uint32(a)), int(b)))))
	case 14:
		return uint32(sub_4624D0(int32(a)))
	case 15:
		return uint32(sub_4649B0(int32(a), int32(b), int32(c)))
	case 16:
		return uint32(sub_464B40(int32(a), int32(b)))
	case 17:
		return uint32(nox_xxx_clientDequip_464B70(int32(a)))
	case 18:
		return uint32(nox_xxx_trade_4657B0(int16(a)))
	case 19:
		return uint32(nox_xxx_clientTrade_465870(int16(a)))
	case 20:
		nox_xxx_cliInventorySpriteUpd_465A30()
		return 0
	case 21:
		return uint32(C.nox_xxx_clientDrop_465BE0((*C.int2)(unsafe.Pointer(a))))
	case 22:
		return uint32(C.nox_xxx_clientKeyEquip_465C30(C.int(a), C.int(b)))
	case 23:
		nox_xxx_clientUse_465C70(int32(a))
		return 0
	case 24:
		nox_client_invAlterWeapon_4672C0()
		return 0
	case 25:
		return uint32(uiInventoryCapacity(int32(a), int32(b)))
	default:
		panic("inventory transaction operation")
	}
}
