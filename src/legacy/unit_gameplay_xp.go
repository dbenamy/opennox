package legacy

/*
#include "GAME3_2.h"
#include "common__gamemech__pausefx.h"
static void unit_xp_notice(nox_object_t* u, wchar2_t* text, unsigned int value) {
 nox_xxx_netSendLineMessage_4D9EB0((int)u,text,value);
}
*/
import "C"

import (
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func unitExperienceNotice(u *server.Object, key, file string, amount float32) {
	text := GetServer().S().Strings().GetStringInFile(strman.ID(key), file)
	C.unit_xp_notice(asObjectC(u), internWStr(text), C.uint32_t(uint32(int64(amount))))
}
func unitExperienceLevel(u *server.Object) {
	if Nox_xxx_gameGet_4DB1B0() && Sub_4DB1C0() != nil {
		return
	}
	pl := u.UpdateDataPlayer().Player
	level := (*byte)(unsafe.Add(unsafe.Pointer(pl), 3684))
	xp := *(*float32)(unsafe.Add(u.CObj(), 28))
	if !(GetServer().S().Balance.FloatInd("XPTable", int(int8(*level))+1) <= float64(xp)) {
		return
	}
	*level++
	addProtectionRecord(int32(*equipmentWord(unsafe.Pointer(pl), 4644)), 1)
	controlReadStats(u, 1)
	if noxflags.HasGame(noxflags.GameModeCoop) {
		runtimePauseStart(u, 0)
	} else {
		GetServer().S().Audio.EventObj(902, u, 2, u.NetCode)
		unitExperienceNotice(u, "LevelUP", "C:\\NoxPost\\src\\Server\\GameMech\\explevel.c", 0)
	}
}
func unitGiveExperience(u *server.Object, amount float32) {
	xp := (*float32)(unsafe.Add(u.CObj(), 28))
	*xp = amount + *xp
	updateProtectionFloat(int32(*equipmentWord(controlPlayer(u), 4604)), amount, true)
	gameplayReportExperience(u)
	unitExperienceNotice(u, "health.c:gainpoints", "C:\\NoxPost\\src\\Server\\GameMech\\explevel.c", amount)
	unitExperienceLevel(u)
}
func unitRewardExperience(u *server.Object, target float32) float64 {
	xp := (*float32)(unsafe.Add(u.CObj(), 28))
	before := *xp
	if before >= target {
		return 0
	}
	// The C call to the coefficient accessor spills subtraction to float32.
	delta := float32(target - before)
	gain := float64(delta)*float64(*memmap.PtrFloat32(0x587000, 206148)) + 1
	returned := float32(gain)
	*xp = float32(gain + float64(before))
	updateProtectionFloat(int32(*equipmentWord(controlPlayer(u), 4604)), returned, true)
	gameplayReportExperience(u)
	unitExperienceLevel(u)
	return float64(returned)
}
func unitMonsterReward(victim *server.Object) {
	if victim == nil || !noxflags.HasGame(noxflags.GameModeCoop) || victim.Obj130 == nil {
		return
	}
	attacker := victim.Obj130
	player := attacker.FindOwnerChainPlayer()
	if player.ObjClass&4 == 0 {
		return
	}
	for it := attacker; it != player; it = it.ObjOwner {
		if it.ObjClass&2 != 0 {
			if !server.Nox_xxx_creatureIsMonitored_500CC0(player, it) {
				return
			}
			break
		}
	}
	gain := unitRewardExperience(player, *(*float32)(unsafe.Add(victim.CObj(), 28)))
	if gain > 0 {
		unitExperienceNotice(player, "gainpoints", "C:\\NoxPost\\src\\Server\\Object\\health.c", float32(gain))
	}
}
