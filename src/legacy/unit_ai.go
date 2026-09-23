package legacy

/*
#include "defs.h"
#include "GAME1_1.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "server__script__script.h"
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_gameSetAudioFadeoutMb_501AC0 func(v int)
	Nox_xxx_unitUpdateMonster_50A5C0     func(a1 *server.Object)
)

type Nox_player_polygon_check_data struct {
	Field_0 [35]uint32
}

//export sub_545E60
func sub_545E60(a1c *nox_object_t) int { return asObjectS(a1c).Sub_545E60() }

//export nox_xxx_gameSetAudioFadeoutMb_501AC0
func nox_xxx_gameSetAudioFadeoutMb_501AC0(v int) { Nox_xxx_gameSetAudioFadeoutMb_501AC0(v) }

//export nox_xxx_monsterPopAction_50A160
func nox_xxx_monsterPopAction_50A160(a1 *nox_object_t) int {
	return asObjectS(a1).MonsterPopAction()
}

//export nox_xxx_monsterPushAction_50A260_impl
func nox_xxx_monsterPushAction_50A260_impl(u *nox_object_t, act int, file *C.char, line int) unsafe.Pointer {
	return asObjectS(u).MonsterPushActionImpl(ai.ActionType(act), GoString(file), line).C()
}

//export nox_xxx_unitUpdateMonster_50A5C0
func nox_xxx_unitUpdateMonster_50A5C0(a1 *nox_object_t) {
	Nox_xxx_unitUpdateMonster_50A5C0(asObjectS(a1))
}

//export nox_xxx_monsterClearActionStack_50A3A0
func nox_xxx_monsterClearActionStack_50A3A0(a1 *nox_object_t) {
	asObjectS(a1).ClearActionStack()
}

//export sub_50B810
func sub_50B810(obj *nox_object_t, p *C.float2) int {
	return bool2int(GetServer().Sub_50B810(asObjectS(obj), (*types.Pointf)(unsafe.Pointer(p))))
}

//export sub_50B500
func sub_50B500() {
	GetServer().S().AI.Paths.Sub_50B500()
}

func sub_50B510() {
	GetServer().S().AI.Paths.Sub_50B510()
}

func Nox_xxx_mobSearchEdible_544A00(a1 *server.Object, a2 float32) int {
	if u := lifecycleFoodSearch(a1, a2, false); u != nil {
		return int(uintptr(u.CObj()))
	}
	return 0
}
func Nox_xxx_weaponGetStaminaByType_4F7E80(a1 int) int {
	return int(nox_xxx_weaponGetStaminaByType_4F7E80(C.int(a1)))
}
func Nox_xxx_mobGetMoveAttemptTime_534810(a1 *server.Object) int {
	return bool2int(monsterMoveAttempt(a1))
}
func Nox_xxx_unitIsDangerous_547120(a1 *server.Object, a2 *server.Object) {
	monsterDangerous(a1, a2)
}
func Nox_xxx_checkIsKillable_528190(a1 *server.Object) int {
	return int(visibilityKillable(a1))
}
func Nox_xxx_polygonIsPlayerInPolygon_4217B0(a1 unsafe.Pointer, a2 int) *Nox_player_polygon_check_data {
	return (*Nox_player_polygon_check_data)(unsafe.Pointer(mapPolygonFind((*[2]int32)(a1), uint32(a2), false)))
}
func Nox_xxx_mobAction_50A910(a1 *server.Object) {
	monsterControlRefresh(a1)
}
func Nox_xxx_monsterGetSoundSet_424300(a1 *server.Object) unsafe.Pointer {
	return resourceMonsterSound(a1)
}
func Nox_xxx_monsterPlayHurtSound_532800(a1 *server.Object) {
	unitHurtSound(a1)
}
func Nox_xxx_mobAction_5469B0(a1 *server.Object) {
	monsterIdleAudio(a1)
}
func Nox_xxx_unitUpdateSightMB_5281F0(a1 *server.Object) {
	visibilityUpdateSight(a1)
}
func Nox_xxx_monsterMainAIFn_547210(a1 *server.Object) {
	monsterMainAI(a1)
}
func Nox_xxx_updateNPCAnimData_50A850(a1 *server.Object) {
	monsterControlAnimation(a1)
}
func Nox_xxx_monsterPolygonEnter_421FF0(a1 *server.Object) {
	mapPolygonMonster(a1)
}
func Nox_xxx_monsterMimicCheckMorph_534950(a1 *server.Object) {
	monsterMimicMorph(a1)
}
func Sub_5466F0(a1 *server.Object) int {
	return investigateHeardSound(a1)
}
func Nox_xxx_mobHealSomeone_5411A0(a1 *server.Object) {
	monsterHealSomeone(a1)
}
func Nox_xxx_mobActionCast_5413B0(a1 *server.Object, a2 int) {
	monsterActionCast(a1, a2)
}
