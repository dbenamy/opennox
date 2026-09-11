package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "common__gamemech__pausefx.h"
extern uint32_t dword_5d4594_2488720, dword_5d4594_2488724;
*/
import "C"
import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func inventoryPriMessage(u *server.Object, key string) {
	C.nox_xxx_netPriMsgToPlayer_4DA2C0(asObjectC(u), internCStr(key), 0)
}
func inventoryWeaponPickup(u, it *server.Object, arg, equip int) int {
	if u.ObjClass&4 != 0 && noxflags.HasGame(4096) && it.ObjSubClass&0x200000 != 0 {
		limit := float32(C.nox_xxx_gamedataGetFloat_419D40(internCStr("ForceOfNatureStaffLimit")))
		if int32(C.nox_xxx_inventoryCountObjects_4E7D30(inventoryInt(u), C.int(it.ObjFlags))) >= floatToInt32(limit) {
			inventoryPriMessage(u, "pickup.c:MaxSameItem")
			inventorySound(925, u, 0, 0)
			return 0
		}
	}
	if !noxflags.HasGame(2048|4096) && C.sub_409F40(2) != 0 {
		duplicate := it.ObjSubClass&0x82 == 0 && C.sub_4E7EC0(inventoryInt(u), asObjectC(it)) != 0
		if it.ObjSubClass&0x40 != 0 {
			for owned := u.Field129; owned != nil; owned = owned.Field128 {
				if owned.ObjClass&0x1000000 != 0 && owned.ObjSubClass&0x40 != 0 {
					duplicate = true
					break
				}
			}
		}
		if duplicate {
			if u.ObjClass&4 != 0 {
				inventoryPriMessage(u, "weapon.c:CannotPickupDuplicateWeapon")
				inventorySound(925, u, 2, int(u.NetCode))
			}
			return 0
		}
	}
	if !noxflags.HasGame(2048|4096) && u.ObjClass&4 != 0 && !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		inventoryPriMessage(u, "weapon.c:WeaponEquipClassFail")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryPickup(u, it, arg) != 1 {
		return 0
	}
	if u.ObjClass&4 != 0 {
		ud := u.UpdateData
		equipped := C.int(0)
		if *(*uint32)(unsafe.Add(ud, 104)) == 0 && C.sub_419E60(asObjectC(u)) == 0 && C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) != 2 {
			equipped = C.nox_xxx_playerEquipWeapon_53A420((*C.uint32_t)(u.CObj()), asObjectC(it), C.int(equip), 0)
		}
		if C.sub_419E60(asObjectC(u)) == 0 && C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) == 2 {
			flags := *(*uint32)(unsafe.Add(unsafe.Pointer(u.UpdateDataPlayer().Player), 4))
			if flags&0xC != 0 && flags&2 == 0 {
				equipped = C.nox_xxx_playerEquipWeapon_53A420((*C.uint32_t)(u.CObj()), asObjectC(it), C.int(equip), 0)
			}
		}
		if equipped == 0 {
			offset := -1
			if it.ObjClass&0x1000 != 0 && it.ObjSubClass&0x47F0000 != 0 {
				offset = 108
			} else if it.ObjClass&0x1000000 != 0 && it.ObjSubClass&0x82 != 0 {
				offset = 1
			}
			if offset >= 0 {
				other := offset + 1
				if offset == 1 {
					other = 0
				}
				C.nox_xxx_netReportCharges_4D82B0(C.int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), asObjectC(it), C.char(*(*byte)(unsafe.Add(it.UseData.Ptr, offset))), C.char(*(*byte)(unsafe.Add(it.UseData.Ptr, other))))
			}
		}
	}
	C.sub_53A6C0(inventoryInt(u), asObjectC(it))
	C.nox_xxx_decay_5116F0(asObjectC(it))
	return 1
}
func inventoryAmmoPickup(u, it *server.Object, arg, equip int) int {
	bits := uint32(C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)))
	if u.ObjClass&4 != 0 && bits&0x82 != 0 {
		data := unsafe.Slice((*byte)(it.UseData.Ptr), 3)
		mods := unsafe.Slice((*uint32)(it.InitData), 4)
		for old := u.InvFirstItem; old != nil; old = old.InvNextItem {
			if old.TypeInd != it.TypeInd || old.ObjClass&0x1000000 == 0 || uint32(C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(old)))&bits == 0 {
				continue
			}
			prev := unsafe.Slice((*byte)(old.UseData.Ptr), 3)
			prevMods := unsafe.Slice((*uint32)(old.InitData), 4)
			same := true
			for i, v := range mods {
				if prevMods[i] != v {
					same = false
				}
			}
			if !same || prev[2] != 0 || int(data[0])+int(prev[0]) > 250 {
				continue
			}
			prev[1] += data[1]
			prev[0] += data[0]
			C.nox_xxx_netReportCharges_4D82B0(C.int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), asObjectC(old), C.char(prev[1]), C.char(prev[0]))
			GetServer().DelayedDelete(it)
			C.sub_53A6C0(inventoryInt(u), asObjectC(it))
			return 1
		}
	}
	return inventoryWeaponPickup(u, it, arg, equip)
}
func inventoryOblivionPickup(u, it *server.Object, arg, equip int) int {
	rv := inventoryWeaponPickup(u, it, arg, equip)
	if rv == 1 && u.ObjClass&4 != 0 && C.sub_419E60(asObjectC(u)) == 0 {
		for i, key := range []string{"weapon.c:PickupHalberdOblivion", "weapon.c:PickupHeartOblivion", "weapon.c:PickupWierdlingOblivion", "weapon.c:PickupOrbOblivion"} {
			if uint32(it.ObjSubClass)&(0x800000<<i) != 0 {
				inventoryPriMessage(u, key)
				inventorySound(914+i, u, 0, 0)
				break
			}
		}
		C.sub_57AF30(inventoryInt(u), 1)
		C.nox_xxx_playerTryEquip_4F2F70(asObjectC(u), asObjectC(it))
	}
	return rv
}
func inventoryPlain(it *server.Object) bool {
	for _, v := range unsafe.Slice((*uint32)(it.InitData), 4) {
		if v != 0 {
			return false
		}
	}
	return true
}
func inventoryArmorPickup(u, it *server.Object, arg, equip int) int {
	if *memmap.PtrUint32(0x5D4594, 2488712) == 0 {
		*memmap.PtrUint32(0x5D4594, 2488712) = uint32(GetServer().S().Types.IndByID("StreetSneakers"))
		*memmap.PtrUint32(0x5D4594, 2488716) = uint32(GetServer().S().Types.IndByID("WizardRobe"))
		C.dword_5d4594_2488720 = C.uint32_t(GetServer().S().Types.IndByID("WoodenShield"))
		C.dword_5d4594_2488724 = C.uint32_t(GetServer().S().Types.IndByID("SteelShield"))
	}
	if !noxflags.HasGame(2048|4096) && C.sub_409F40(2) != 0 && C.sub_4E7EC0(inventoryInt(u), asObjectC(it)) != 0 {
		inventoryPriMessage(u, "armor.c:CannotPickupDuplicateArmor")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if !noxflags.HasGame(2048|4096) && u.ObjClass&4 != 0 && !Nox_xxx_playerClassCanUseItem_57B3D0(it, u.UpdateDataPlayer().Player.PlayerClass()) {
		inventoryPriMessage(u, "armor.c:ArmorEquipClassFail")
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryPickup(u, it, arg) != 1 {
		return 0
	}
	if u.ObjClass&4 != 0 {
		old := (*server.Object)(unsafe.Pointer(C.nox_xxx_armorHaveSameSubclass_53E7B0(inventoryInt(u), inventoryInt(it))))
		if C.sub_419E60(asObjectC(u)) == 0 {
			wood, steel := uint32(C.dword_5d4594_2488720), uint32(C.dword_5d4594_2488724)
			sneakers, robe := *memmap.PtrUint32(0x5D4594, 2488712), *memmap.PtrUint32(0x5D4594, 2488716)
			replace := false
			switch {
			case old != nil && uint32(old.TypeInd) != wood && uint32(old.TypeInd) != steel && uint32(old.TypeInd) != sneakers && uint32(old.TypeInd) != robe:
			case old != nil && uint32(old.TypeInd) == robe:
				replace = inventoryPlain(old)
			case it.ObjSubClass&2 == 0:
				replace = true
			case old == nil:
				weapon := *(**server.Object)(unsafe.Add(u.UpdateData, 104))
				replace = weapon == nil || weapon.ObjSubClass&0x7FFE40C == 0
			case uint32(old.TypeInd) == wood:
				replace = inventoryPlain(old)
			case uint32(old.TypeInd) == steel:
				if uint32(it.TypeInd) == wood {
					replace = inventoryPlain(old) && !inventoryPlain(it)
				} else if uint32(it.TypeInd) == steel {
					replace = inventoryPlain(old)
				}
			}
			if replace {
				C.nox_xxx_playerEquipArmor_53E650((*C.uint32_t)(u.CObj()), asObjectC(it), C.int(equip), 0)
			}
		}
	}
	switch {
	case it.Material&0x10 != 0:
		inventorySound(804, u, 0, 0)
	case it.Material&8 != 0:
		inventorySound(810, u, 0, 0)
	case it.Material&4 != 0:
		inventorySound(807, u, 0, 0)
	case it.Material&2 != 0:
		sound := 813
		if it.ObjSubClass&0x20 != 0 {
			sound = 816
		}
		inventorySound(sound, u, 0, 0)
	}
	C.nox_xxx_decay_5116F0(asObjectC(it))
	return 1
}
