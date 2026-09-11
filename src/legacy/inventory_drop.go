package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
*/
import "C"
import (
	"encoding/binary"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/ccall"

	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func inventoryCache(off uintptr, name string) uint32 {
	p := memmap.PtrUint32(0x5D4594, off)
	if *p == 0 {
		*p = uint32(GetServer().S().Types.IndByID(name))
	}
	return *p
}
func inventorySound(id int, u *server.Object, a, b int) {
	C.nox_xxx_aud_501960(C.int(id), asObjectC(u), C.int(a), C.int(b))
}
func inventoryMessage(kind int, u *server.Object, value uint32) {
	var data [10]byte
	binary.LittleEndian.PutUint32(data[2:], u.NetCode)
	binary.LittleEndian.PutUint32(data[6:], value)
	C.nox_xxx_netInformTextMsg2_4DA180(C.int(kind), (*C.uint8_t)(unsafe.Pointer(&data[0])))
}
func inventoryDefaultDrop(u, it *server.Object, pos *types.Pointf) int {
	if it.InvHolder != u {
		return 0
	}
	if u.ObjClass&4 != 0 && inventoryDroppable(it) && C.sub_53EC80(inventoryInt(it), 1) != 0 {
		if u.ObjFlags&0x8020 == 0 {
			C.nox_xxx_netPriMsgToPlayer_4DA2C0(asObjectC(u), internCStr("drop.c:CantDropThat"), 0)
			inventorySound(925, u, 2, int(u.NetCode))
		}
		return 0
	}
	inventoryRemove(u, it)
	GetServer().CreateObjectAt(it, nil, *pos)
	if it.ObjClass&0x1000000 != 0 && C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(it)) == 4 {
		for ammo := u.InvFirstItem; ammo != nil; ammo = ammo.InvNextItem {
			if ammo.ObjClass&0x1000000 != 0 && C.nox_xxx_weaponInventoryEquipFlags_415820(asObjectC(ammo)) == 2 && ammo.ObjFlags&0x100 == 0 && *(*byte)(unsafe.Add(ammo.UseData.Ptr, 2)) != 0 {
				inventoryRemove(u, ammo)
				GetServer().DelayedDelete(ammo)
				break
			}
		}
	}
	if it.ObjClass&0x10000000 != 0 {
		team := it.TeamVal.ID
		value := C.sub_4ECBD0(inventoryInt(it))
		inventoryMessage(7, u, uint32(value))
		C.nox_xxx_netMarkMinimapForAll_4174B0(inventoryInt(it), 1)
		*(*uint32)(unsafe.Add(it.UpdateData, 8)) = GetServer().S().Frame()
		C.sub_4E82C0(C.uchar(team), 2, C.char(value), 0)
	}
	glyph := inventoryCache(1568252, "Glyph")
	if !noxflags.HasGame(2048|4096) && it.ObjFlags&0x80000 == 0 && it.ObjClass&0x10000000 == 0 && uint32(it.TypeInd) != glyph {
		Nox_xxx_unitSetDecayTime_511660(it, int(10*GetServer().S().TickRate()))
	}
	Nox_xxx_unitRaise_4E46F0(it, 0)
	if it.ObjClass&64 != 0 {
		inventorySound(821, it, 0, 0)
	}
	if it.ObjClass&2 != 0 {
		*(*uint32)(unsafe.Add(it.UpdateData, 1360)) = 15
		*(*uint32)(unsafe.Add(it.UpdateData, 1440)) |= 0x100
	}
	if *memmap.PtrUint32(0x5D4594, 1568256) == 0 {
		*memmap.PtrUint32(0x5D4594, 1568256) = uint32(GetServer().S().Types.IndByID("Torch"))
		*memmap.PtrUint32(0x5D4594, 1568244) = uint32(GetServer().S().Types.IndByID("Lantern"))
	}
	if uint32(it.TypeInd) == *memmap.PtrUint32(0x5D4594, 1568256) || uint32(it.TypeInd) == *memmap.PtrUint32(0x5D4594, 1568244) {
		Nox_xxx_spellBuffOff_4FF5B0(u, 15)
	}
	return 1
}
func inventoryTrapDrop(u, it *server.Object, pos *types.Pointf) int {
	if C.nox_xxx_mapTileAllowTeleport_411A90((*C.float2)(unsafe.Pointer(pos))) != 0 {
		inventorySound(925, u, 2, int(u.NetCode))
		return 0
	}
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	GetServer().S().ObjSetOwner(u, it)
	return 1
}
func inventoryGlyphDrop(u, it *server.Object, pos *types.Pointf) int {
	if inventoryTrapDrop(u, it, pos) == 0 {
		return 0
	}
	*(*types.Pointf)(unsafe.Add(it.InitData, 28)) = *pos
	dir := types.Pointf{X: u.PosVec.X - pos.X, Y: u.PosVec.Y - pos.Y}
	angle := uint16(C.nox_xxx_math_509ED0((*C.float2)(unsafe.Pointer(&dir))))
	it.Direction1, it.Direction2 = server.Dir16(angle), server.Dir16(angle)
	inventorySound(825, it, 0, 0)
	return 1
}
func inventoryDrop(u, it *server.Object, pos *types.Pointf) int {
	if it == nil {
		return 0
	}
	if noxflags.HasGame(0x2000) && !noxflags.HasGame(4096) && it.ObjClass&0x3001010 != 0 {
		it.ObjFlags |= 0x40
		C.nox_xxx_unit_511810(asObjectC(it))
	}
	if it.Drop.Ptr != nil {
		if fn := inventoryNativeDrops[it.Drop.Ptr]; fn != nil {
			return fn(u, it, pos)
		}
		return ccall.CallIntPtr3(it.Drop.Ptr, u.CObj(), it.CObj(), unsafe.Pointer(pos))
	}
	return inventoryDefaultDrop(u, it, pos)
}
func inventoryPotionDrop(u, it *server.Object, pos *types.Pointf) int {
	if inventoryDefaultDrop(u, it, pos) == 0 {
		return 0
	}
	inventorySound(833, it, 0, 0)
	if !noxflags.HasGame(2048 | 4096) {
		Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
	}
	return 1
}
func inventoryFoodSound(off uintptr, u, it *server.Object) {
	for ; *memmap.PtrUint16(0x587000, off+6) != 0; off += 8 {
		if uint32(it.ObjSubClass)&*memmap.PtrUint32(0x587000, off) != 0 || it.Material&*memmap.PtrUint16(0x587000, off+4) != 0 {
			inventorySound(int(*memmap.PtrUint16(0x587000, off+6)), u, 0, 0)
			return
		}
	}
}
func inventoryFoodDrop(u, it *server.Object, pos *types.Pointf) int {
	if u == nil || it == nil || pos == nil {
		return 0
	}
	rv := inventoryDefaultDrop(u, it, pos)
	if rv != 0 {
		if !noxflags.HasGame(2048) {
			Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
		}
		inventoryFoodSound(205704, u, it)
	}
	return rv
}
func inventoryEquipmentDrop(u, it *server.Object, pos *types.Pointf, armor bool) int {
	if inventoryDefaultDrop(u, it, pos) != 1 {
		return 0
	}
	if armor {
		C.sub_53EAE0(inventoryInt(it))
	} else {
		C.sub_53AAB0(inventoryInt(it))
	}
	if !noxflags.HasGame(2048|4096) && C.sub_409F40(2) != 0 {
		Nox_xxx_unitSetDecayTime_511660(it, int(25*GetServer().S().TickRate()))
	}
	return 1
}
