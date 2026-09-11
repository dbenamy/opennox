package legacy

/*
#include <string.h>
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
*/
import "C"

import (
	"encoding/binary"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

func shopRepairQuote(u *server.Object, s *shopSession, code uint32) uint32 {
	var item *server.Object
	for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode == code {
			break
		}
	}
	if item == nil {
		return 0
	}
	charge := false
	if item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x47f0000 != 0 {
		data := unsafe.Slice((*byte)(item.UseData.Ptr), 110)
		charge = data[108] < data[109] && data[109] != 0
	}
	hp := item.HealthData
	if (hp != nil && hp.Cur != hp.Max && hp.Max != 0) || charge {
		b := [8]byte{0xc9, 31}
		binary.LittleEndian.PutUint16(b[2:], uint16(code))
		binary.LittleEndian.PutUint32(b[4:], uint32(shopPrice(2, s, item)))
		return shopSend(u, b[:], 0)
	}
	GetServer().S().Audio.EventObj(925, u, 0, 0)
	return uint32(uintptr(item.CObj()))
}
func shopRepair(u *server.Object, s *shopSession, code uint32) uint32 {
	gold := shopGetGold(u)
	if u.InvFirstItem == nil {
		return gold
	}
	var item *server.Object
	for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode == code {
			break
		}
	}
	if item == nil {
		return code
	}
	shopSubGold(u, uint32(shopPrice(2, s, item)))
	resourceSetHP(item, item.HealthData.Max)
	ind := C.int(uint8(u.UpdateDataPlayer().Player.PlayerInd))
	C.nox_xxx_itemReportHealth_4D87A0(ind, asObjectC(item))
	if item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x47f0000 != 0 {
		data := unsafe.Slice((*byte)(item.UseData.Ptr), 110)
		if C.nox_xxx_rechargeItem_53C520(C.int(uintptr(item.CObj())), 100) != 0 {
			C.nox_xxx_netReportCharges_4D82B0(ind, asObjectC(item), C.char(data[108]), C.char(data[109]))
		}
	}
	C.sub_4D8870(ind, C.int(uintptr(u.CObj())))
	GetServer().S().Audio.EventObj(803, u, 2, u.NetCode)
	return code
}
func shopSell(u *server.Object, s *shopSession, typ int32, count uint32) {
	shopGetGold(u)
	for i := uint32(0); i < count; i++ {
		var item *server.Object
		for item = u.InvFirstItem; item != nil; item = item.InvNextItem {
			if int32(item.TypeInd) == typ {
				break
			}
		}
		if item == nil {
			return
		}
		C.sub_4ED0C0(asObjectC(u), asObjectC(item))
		GetServer().DelayedDelete(item)
		shopAddGold(u, uint32(shopPrice(0, s, item)))
		C.sub_4D8870(C.int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), C.int(uintptr(u.CObj())))
		if i+1 == count {
			GetServer().S().Audio.EventObj(307, u, 2, u.NetCode)
			return
		}
	}
}
func shopLoad(s *shopSession) {
	var vendor *server.Object
	for _, u := range s.Units {
		if u != nil && u.ObjClass&2 != 0 && u.ObjSubClass&8 != 0 {
			vendor = u
			break
		}
	}
	if vendor == nil {
		return
	}
	core := GetServer().S()
	if noxflags.HasGame(noxflags.GameModeQuest) {
		cutoff := uint32(floatToInt32(float32(core.Balance.Float("ShopAnkhCutoffStage"))))
		if uint32(Nox_game_getQuestStage_4E3CC0()) < cutoff {
			if u := core.NewObjectByTypeID("AnkhTradable"); u != nil {
				shopAdd(s, u)
			}
		}
		for off := uintptr(234816); ; off += 4 {
			name := *memmap.PtrPtr(0x587000, off)
			if name == nil {
				break
			}
			u := asObjectS(C.nox_xxx_newObjectByTypeID_4E3810((*C.char)(name)))
			if u != nil {
				shopAdd(s, u)
			}
		}
		marker := core.NewObjectByTypeID("RewardMarker")
		if marker == nil {
			return
		}
		stage := uint32(Nox_game_getQuestStage_4E3CC0()) + 2
		for _, category := range [...]uint32{8, 8, 8, 8, 16, 16, 16, 16, 1, 1, 1, 1, 1, 1, 4, 4, 4} {
			*(*uint32)(marker.InitData) = category
			if u := (*server.Object)(unsafe.Pointer(C.nox_server_rewardgen_activateMarker_4F0720(C.int(uintptr(marker.CObj())), C.uint(stage)))); u != nil {
				shopAdd(s, u)
			}
		}
		if core.Rand.Logic.Int(0, 100) > 90 {
			*(*uint32)(marker.InitData) = 2
			if u := (*server.Object)(unsafe.Pointer(C.nox_server_rewardgen_activateMarker_4F0720(C.int(uintptr(marker.CObj())), C.uint(stage)))); u != nil {
				shopAdd(s, u)
			}
		}
		GetServer().DelayedDelete(marker)
		return
	}
	data := vendor.InitData
	if data == nil {
		return
	}
	for i := 0; i < int(*(*byte)(data)); i++ {
		e := (*shopStockEntry)(unsafe.Add(data, 4+28*i))
		for j := 0; j < int(e.Count); j++ {
			u := asObjectS(C.nox_xxx_newObjectWithTypeInd_4E3450(C.int(e.Type)))
			if u == nil {
				continue
			}
			if uint32(u.ObjClass)&0x13001000 != 0 {
				// C copied an uninitialized fifth word. Define only that former
				// indeterminate tail as zero; keep the four actual modifiers.
				var mods [5]uint32
				for k, p := range e.Modifiers {
					mods[k] = uint32(uintptr(unsafe.Pointer(p)))
				}
				C.nox_xxx_modifSetItemAttrs_4E4990(asObjectC(u), (*C.int)(unsafe.Pointer(&mods[0])))
			}
			if u.Xfer == C.nox_xxx_XFerSpellReward_4F5F30 {
				*(*byte)(u.UseData.Ptr) = byte(e.Reward)
			}
			if u.Xfer == C.nox_xxx_XFerAbilityReward_4F6240 {
				*(*byte)(u.UseData.Ptr) = byte(e.Reward)
			}
			if u.Xfer == C.nox_xxx_XFerFieldGuide_4F6390 {
				C.strcpy((*C.char)(u.UseData.Ptr), (*C.char)(C.nox_xxx_getUnitNameByThingType_4E3A80(C.int(e.Reward))))
			}
			shopAdd(s, u)
		}
	}
}
