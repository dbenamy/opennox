package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
*/
import "C"
import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func inventoryPickup(u, it *server.Object, arg int) int {
	return bool2int(Nox_xxx_pickupDefault_4F31E0(u, it, arg, 0))
}
func inventoryFoodPickup(u, it *server.Object, arg int) int {
	if u == nil || it == nil {
		return 0
	}
	if C.sub_419E60(asObjectC(u)) == 0 && it.ObjSubClass&0x84 == 0 {
		ccall.CallVoidPtr2(it.Use.Ptr, u.CObj(), it.CObj())
	}
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		inventoryFoodSound(215640, u, it)
	}
	return rv
}
func inventoryUsePickup(u, it *server.Object, arg int) int {
	effectsUse(u, it)
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	return inventoryPickup(u, it, arg)
}
func inventoryTrapPickup(u, it *server.Object, arg int) int {
	if C.nox_xxx_unitHasThatParent_4EC4F0(asObjectC(it), asObjectC(u)) != 0 {
		rv := inventoryPickup(u, it, arg)
		if rv != 0 {
			inventorySound(824, u, 0, 0)
		}
		return rv
	}
	if u.ObjClass&4 != 0 {
		inventorySound(925, u, 2, int(u.NetCode))
	}
	return 0
}
func inventoryBookPickup(u, it *server.Object, arg int, ability bool) int {
	if noxflags.HasGame(6144) {
		effectsUse(u, it)
	}
	if it.ObjFlags&0x20 != 0 {
		return 1
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		sound := 826
		if !ability && it.ObjSubClass&1 == 0 {
			sound = 828
		}
		inventorySound(sound, u, 0, 0)
	}
	return rv
}
func inventoryAnkhPickup(u, it *server.Object) int {
	if u.ObjClass&4 == 0 {
		return 0
	}
	*(*uint32)(unsafe.Add(u.UpdateData, 320))++
	GetServer().DelayedDelete(it)
	inventorySound(1004, u, 0, 0)
	return 1
}
func inventoryCrownPickup(u, it *server.Object, arg int) int {
	data := it.UpdateData
	if u.ObjClass&4 == 0 {
		return 0
	}
	rv := inventoryPickup(u, it, arg)
	if rv != 0 {
		*(*uint32)(unsafe.Add(u.UpdateData, 264)) = GetServer().S().Frame()
		GetServer().S().ObjSetOwner(u, it)
		C.nox_xxx_buffApplyTo_4FF380(asObjectC(u), 30, 0, 5)
		inventorySound(313, u, 0, 0)
		inventoryMessage(10, u, uint32(u.TeamVal.ID))
		C.nox_xxx_netUnmarkMinimapSpec_417470(inventoryInt(it), 1)
	}
	*(*uint32)(unsafe.Add(data, 4)) = 0
	return rv
}
func inventoryTreasurePickup(u, it *server.Object, arg int) int {
	if inventoryPickup(u, it, arg) == 0 {
		return 0
	}
	if u.ObjClass&4 == 0 || !noxflags.HasGame(64) {
		return 1
	}
	pl := u.UpdateDataPlayer().Player
	inventorySound(307, u, 0, 0)
	pl.Field2152++
	pl.Field2156 = uint32(C.nox_xxx_scavengerTreasureMax_4D1600())
	C.nox_xxx_scavengerHuntReport_4D8CD0(inventoryInt(u))
	core := GetServer().S()
	if !u.TeamVal.Has() {
		if pl.Field2152 == uint32(C.nox_xxx_scavengerTreasureMax_4D1600()) {
			noxflags.SetGame(8)
			C.nox_xxx_changeScore_4D8E90(inventoryInt(u), 1)
			C.nox_xxx_netReportLesson_4D8EF0(asObjectC(u))
			for other := core.Players.FirstUnit(); other != nil; other = core.Players.NextUnit(other) {
				if other != u {
					C.nox_xxx_playerIncrementElimDeath_4D8D40(inventoryInt(other))
					C.nox_xxx_netReportLesson_4D8EF0(asObjectC(other))
				}
			}
		}
		return 1
	}
	if team := core.Teams.ByID(u.TeamVal.ID); team != nil {
		var total uint32
		for other := core.Players.FirstUnit(); other != nil; other = core.Players.NextUnit(other) {
			if Nox_xxx_teamCompare2_419180(&other.TeamVal, team.ID()) != 0 {
				total += other.UpdateDataPlayer().Player.Field2152
			}
		}
		if total == uint32(C.nox_xxx_scavengerTreasureMax_4D1600()) {
			noxflags.SetGame(8)
		}
	}
	return 1
}
