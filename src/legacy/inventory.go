package legacy

/*
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
extern uint32_t dword_5d4594_2488728;
*/
import "C"

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func inventoryInt(u *server.Object) C.int { return C.int(uintptr(u.CObj())) }
func inventoryWeight(u *server.Object) {
	var weight uint32
	for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
		weight += uint32(it.Weight)
	}
	u.UpdateDataPlayer().Player.Field3656 = uint32(bool2int(weight > uint32(u.CarryCapacity)))
}
func inventoryRemove(u, it *server.Object) {
	if u == nil || it == nil {
		return
	}
	report := C.int(1)
	if u.ObjClass&4 != 0 {
		pl := u.UpdateDataPlayer().Player
		if !noxflags.HasGame(4096) && u.ObjFlags&0x8000 != 0 && it.ObjClass&0x13001000 != 0 {
			report = 0
		}
		if it.ObjClass&0x10000000 != 0 && noxflags.HasGame(32) {
			*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 4)) &^= 1
			if report == 1 {
				C.nox_xxx_netReportDequip_4D84C0(255, asObjectC(it))
			}
		}
		C.sub_53E430((*C.uint32_t)(u.CObj()), asObjectC(it), 0, report)
		C.nox_xxx_playerDequipWeapon_53A140((*C.uint32_t)(u.CObj()), asObjectC(it), 0, report)
		C.nox_xxx_netReportDrop_4D8B50(C.int(uint8(pl.PlayerInd)), asObjectC(it))
		toggleProtectionObject(int32(pl.Prot4632), it)
	} else if u.ObjClass&2 != 0 {
		if u.ObjSubClass&0x10 != 0 && it.ObjClass&0x10000000 != 0 && noxflags.HasGame(32) {
			C.nox_xxx_npcSetItemEquipFlags_4E4B20(inventoryInt(u), asObjectC(it), 0)
		}
		C.sub_53E430((*C.uint32_t)(u.CObj()), asObjectC(it), 1, 1)
		C.nox_xxx_playerDequipWeapon_53A140((*C.uint32_t)(u.CObj()), asObjectC(it), 1, 1)
	}
	if prev := it.Field125; prev != nil {
		prev.InvNextItem = it.InvNextItem
	} else {
		u.InvFirstItem = it.InvNextItem
	}
	if next := it.InvNextItem; next != nil {
		next.Field125 = it.Field125
	}
	it.InvHolder = nil
	GetServer().S().ObjClearOwner(it)
	if u.ObjClass&4 != 0 {
		inventoryWeight(u)
	}
}
func inventoryInsert(u, it *server.Object, report int) {
	if u == nil || it == nil || u.ObjFlags&0x20 != 0 || it.ObjFlags&0x20 != 0 {
		return
	}
	it.Field125 = nil
	it.InvNextItem = u.InvFirstItem
	if next := u.InvFirstItem; next != nil {
		next.Field125 = it
	}
	u.InvFirstItem = it
	it.InvHolder = u
	GetServer().S().ObjSetOwner(u, it)
	if u.ObjClass&4 != 0 {
		pl := u.UpdateDataPlayer().Player
		if report != 0 {
			C.nox_xxx_netReportPickup_4D8A60(C.int(uint8(pl.PlayerInd)), asObjectC(it))
		}
		toggleProtectionObject(int32(pl.Prot4632), it)
		inventoryWeight(u)
	}
	if it.ObjClass&64 != 0 {
		C.nox_xxx_aud_501960(820, asObjectC(u), 0, 0)
	}
}
func inventoryDropEligible(u, it *server.Object) bool {
	return it.ObjFlags&0x20 != 0 || u.ObjClass&6 == 0 || it.ObjFlags&0x10000000 == 0
}
func inventoryDroppable(it *server.Object) bool {
	if it == nil {
		return false
	}
	if C.dword_5d4594_2488728 == 0 {
		C.sub_53EC40()
	}
	for off := uintptr(279432); *memmap.PtrUint32(0x587000, off) != 0; off += 12 {
		if *memmap.PtrUint32(0x587000, off+4) == uint32(it.TypeInd) {
			return true
		}
	}
	return false
}
