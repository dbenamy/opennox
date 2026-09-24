package legacy

/*
#include "defs.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "server__magic__plyrspel.h"


void nox_xxx_updateProjectile_53AC10(nox_object_t* a1);
void nox_xxx_updateDeathBall_53D080(nox_object_t* a1);
void nox_xxx___mkgmtime_538280(nox_object_t* a1);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/player"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_updatePlayer_4F8100         func(up *server.Object)
	Nox_xxx_updatePixie_53CD20          func(cobj *server.Object)
	Nox_xxx_updatePlayerObserver_4E62F0 func(a1p *server.Object)
	Nox_xxx_updateProjectile_53AC10     func(a1 *server.Object)
	Nox_xxx_updateDeathBall_53D080      func(a1 *server.Object)
	Nox_xxx___mkgmtime_538280           func(a1 *server.Object)
)

var _ = [1]struct{}{}[2200-unsafe.Sizeof(server.MonsterUpdateData{})]
var _ = [1]struct{}{}[556-unsafe.Sizeof(server.PlayerUpdateData{})]

func init() {
	_ = nox_xxx_updatePlayer_4F8100
	server.RegisterObjectUpdateGo("PlayerUpdate", C.nox_xxx_updatePlayer_4F8100, func(u *server.Object) { nox_xxx_updatePlayer_4F8100(asObjectC(u)) }, unsafe.Sizeof(server.PlayerUpdateData{}))
	_ = nox_xxx_updateProjectile_53AC10
	server.RegisterObjectUpdateGo("ProjectileUpdate", C.nox_xxx_updateProjectile_53AC10, func(u *server.Object) { nox_xxx_updateProjectile_53AC10(asObjectC(u)) }, 0)
	server.RegisterObjectUpdateGo("SpellProjectileUpdate", C.nox_xxx_spellFlyUpdate_53B940, func(u *server.Object) { nox_xxx_spellFlyUpdate_53B940(C.int(uintptr(u.CObj()))) }, unsafe.Sizeof(server.SpellProjectileUpdateData{}))
	server.RegisterObjectUpdateGo("AntiSpellProjectileUpdate", C.nox_xxx_updateAntiSpellProj_53BB00, func(u *server.Object) { nox_xxx_updateAntiSpellProj_53BB00(C.int(uintptr(u.CObj()))) }, 28)
	server.RegisterObjectUpdateGo("DoorUpdate", C.nox_xxx_updateDoor_53AC50, func(u *server.Object) { nox_xxx_updateDoor_53AC50(C.int(uintptr(u.CObj()))) }, 52)
	server.RegisterObjectUpdateGo("SparkUpdate", C.nox_xxx_updateSpark_53ADC0, func(u *server.Object) { nox_xxx_updateSpark_53ADC0(C.int(uintptr(u.CObj()))) }, 16)
	server.RegisterObjectUpdateGo("ProjectileTrailUpdate", C.nox_xxx_updateProjTrail_53AEC0, func(u *server.Object) { nox_xxx_updateProjTrail_53AEC0(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("PushUpdate", C.nox_xxx_updatePush_53B030, func(u *server.Object) { nox_xxx_updatePush_53B030(C.int(uintptr(u.CObj()))) }, 12)
	server.RegisterObjectUpdateGo("TriggerUpdate", C.nox_xxx_updateTrigger_53B1B0, func(u *server.Object) { nox_xxx_updateTrigger_53B1B0(C.int(uintptr(u.CObj()))) }, 60)
	server.RegisterObjectUpdateGo("ToggleUpdate", C.nox_xxx_updateToggle_53B060, func(u *server.Object) { nox_xxx_updateToggle_53B060((*C.uint32_t)(u.CObj())) }, 60)
	server.RegisterObjectUpdateGo("MonsterUpdate", C.nox_xxx_unitUpdateMonster_50A5C0, func(u *server.Object) { nox_xxx_unitUpdateMonster_50A5C0(asObjectC(u)) }, unsafe.Sizeof(server.MonsterUpdateData{}))
	server.RegisterObjectUpdateGo("LoopAndDamageUpdate", C.sub_53B300, func(u *server.Object) { sub_53B300(C.int(uintptr(u.CObj()))) }, 16)
	server.RegisterObjectUpdateGo("ElevatorUpdate", C.nox_xxx_updateElevator_53B5D0, func(u *server.Object) { nox_xxx_updateElevator_53B5D0((*C.uint32_t)(u.CObj())) }, 20)
	server.RegisterObjectUpdateGo("ElevatorShaftUpdate", C.nox_xxx_updateElevatorShaft_53B380, func(u *server.Object) { nox_xxx_updateElevatorShaft_53B380(C.int(uintptr(u.CObj()))) }, 16)
	server.RegisterObjectUpdateGo("PhantomPlayerUpdate", C.nox_xxx_updatePhantomPlayer_53B860, func(u *server.Object) { nox_xxx_updatePhantomPlayer_53B860(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("ObeliskUpdate", C.nox_xxx_updateObelisk_53C580, func(u *server.Object) { nox_xxx_updateObelisk_53C580(C.int(uintptr(u.CObj()))) }, unsafe.Sizeof(server.ObeliskUpdateData{}))
	server.RegisterObjectUpdateGo("LifetimeUpdate", C.nox_xxx_updateLifetime_53B8F0, func(u *server.Object) { nox_xxx_updateLifetime_53B8F0(C.int(uintptr(u.CObj()))) }, 4)
	server.RegisterObjectUpdateGo("MagicMissileUpdate", C.nox_xxx_updateMagicMissile_53BDA0, func(u *server.Object) { nox_xxx_updateMagicMissile_53BDA0(C.int(uintptr(u.CObj()))) }, 28)
	server.RegisterObjectUpdateGo("PixieUpdate", C.nox_xxx_updatePixie_53CD20, func(u *server.Object) { nox_xxx_updatePixie_53CD20(asObjectC(u)) }, 28)
	server.RegisterObjectUpdateGo("SkullUpdate", C.nox_xxx_updateShootingTrap_54F9A0, func(u *server.Object) { nox_xxx_updateShootingTrap_54F9A0(C.int(uintptr(u.CObj()))) }, 52)
	server.RegisterObjectUpdateGo("PentagramUpdate", C.nox_xxx_updateTeleportPentagram_53BEF0, func(u *server.Object) { nox_xxx_updateTeleportPentagram_53BEF0(C.int(uintptr(u.CObj()))) }, 24)
	server.RegisterObjectUpdateGo("InvisiblePentagramUpdate", C.nox_xxx_updateInvisiblePentagram_53C0C0, func(u *server.Object) { nox_xxx_updateInvisiblePentagram_53C0C0(C.int(uintptr(u.CObj()))) }, 24)
	server.RegisterObjectUpdateGo("SwitchUpdate", C.nox_xxx_updateSwitch_53B320, func(u *server.Object) { nox_xxx_updateSwitch_53B320((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("BlowUpdate", C.nox_xxx_updateBlow_53C160, func(u *server.Object) { nox_xxx_updateBlow_53C160(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("MoverUpdate", C.nox_xxx_unitUpdateMover_54F740, func(u *server.Object) { nox_xxx_unitUpdateMover_54F740(C.int(uintptr(u.CObj()))) }, 36)
	server.RegisterObjectUpdateGo("BlackPowderBarrelUpdate", C.nox_xxx_updateBlackPowderBarrel_53C9A0, func(u *server.Object) { nox_xxx_updateBlackPowderBarrel_53C9A0((*C.float)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("OneSecondDieUpdate", C.nox_xxx_updateOneSecondDie_53CB60, func(u *server.Object) { nox_xxx_updateOneSecondDie_53CB60(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("WaterBarrelUpdate", C.nox_xxx_updateWaterBarrel_53CB90, func(u *server.Object) { nox_xxx_updateWaterBarrel_53CB90(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("SelfDestructUpdate", C.nox_xxx_updateSelfDestruct_53CC90, func(u *server.Object) { nox_xxx_updateSelfDestruct_53CC90(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("BlackPowderBurnUpdate", C.nox_xxx_updateBlackPowderBurn_53CCB0, func(u *server.Object) { nox_xxx_updateBlackPowderBurn_53CCB0(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("DeathBallUpdate", C.nox_xxx_updateDeathBall_53D080, func(u *server.Object) { nox_xxx_updateDeathBall_53D080(asObjectC(u)) }, 0)
	server.RegisterObjectUpdateGo("DeathBallFragmentUpdate", C.nox_xxx_updateDeathBallFragment_53D220, func(u *server.Object) { nox_xxx_updateDeathBallFragment_53D220(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("MoonglowUpdate", C.nox_xxx_updateMoonglow_53D270, func(u *server.Object) { nox_xxx_updateMoonglow_53D270(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("SentryGlobeUpdate", C.nox_xxx_updateSentryGlobe_510E60, func(u *server.Object) { nox_xxx_updateSentryGlobe_510E60(C.int(uintptr(u.CObj()))) }, 12)
	server.RegisterObjectUpdateGo("TelekinesisUpdate", C.nox_xxx_updateTelekinesis_53D330, func(u *server.Object) { nox_xxx_updateTelekinesis_53D330(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("FistUpdate", C.nox_xxx_updateFist_53D400, func(u *server.Object) { nox_xxx_updateFist_53D400(C.int(uintptr(u.CObj()))) }, 4)
	server.RegisterObjectUpdateGo("MeteorShowerUpdate", C.nox_xxx_updateMeteorShower_53D5A0, func(u *server.Object) { nox_xxx_updateMeteorShower_53D5A0((*C.float)(u.CObj())) }, 4)
	server.RegisterObjectUpdateGo("MeteorUpdate", C.nox_xxx_meteorExplode_53D6E0, func(u *server.Object) { nox_xxx_meteorExplode_53D6E0(C.int(uintptr(u.CObj()))) }, 4)
	server.RegisterObjectUpdateGo("ToxicCloudUpdate", C.nox_xxx_updateToxicCloud_53D850, func(u *server.Object) { nox_xxx_updateToxicCloud_53D850(C.int(uintptr(u.CObj()))) }, 4)
	server.RegisterObjectUpdateGo("SmallToxicCloudUpdate", C.nox_xxx_updateSmallToxicCloud_53D960, func(u *server.Object) { nox_xxx_updateSmallToxicCloud_53D960(C.int(uintptr(u.CObj()))) }, 4)
	server.RegisterObjectUpdateGo("ArachnaphobiaUpdate", C.nox_xxx_updateArachnaphobia_53DA60, func(u *server.Object) { nox_xxx_updateArachnaphobia_53DA60((*C.int)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("ExpireUpdate", C.nox_xxx_updateExpire_53DB00, func(u *server.Object) { nox_xxx_updateExpire_53DB00(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("BreakUpdate", C.nox_xxx_updateBreak_53DB30, func(u *server.Object) { nox_xxx_updateBreak_53DB30((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("OpenUpdate", C.nox_xxx_updateOpen_53DBB0, func(u *server.Object) { nox_xxx_updateOpen_53DBB0((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("BreakAndRemoveUpdate", C.nox_xxx_updateBreakAndRemove_53DC30, func(u *server.Object) { nox_xxx_updateBreakAndRemove_53DC30((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("ChakramInMotionUpdate", C.nox_xxx_updateChakramInMotion_53DCC0, func(u *server.Object) { nox_xxx_updateChakramInMotion_53DCC0(C.int(uintptr(u.CObj()))) }, 28)
	server.RegisterObjectUpdateGo("FlagUpdate", C.nox_xxx_updateFlag_53DDF0, func(u *server.Object) { nox_xxx_updateFlag_53DDF0(C.int(uintptr(u.CObj()))) }, 12)
	server.RegisterObjectUpdateGo("TrapDoorUpdate", C.nox_xxx_updateTrapDoor_53DE80, func(u *server.Object) { nox_xxx_updateTrapDoor_53DE80((*C.uint32_t)(u.CObj())) }, 0)
	server.RegisterObjectUpdateGo("BallUpdate", C.nox_xxx_updateGameBall_53DF40, func(u *server.Object) { nox_xxx_updateGameBall_53DF40(C.int(uintptr(u.CObj()))) }, 32)
	server.RegisterObjectUpdateGo("CrownUpdate", C.nox_xxx_updateCrown_53E1D0, func(u *server.Object) { nox_xxx_updateCrown_53E1D0(C.int(uintptr(u.CObj()))) }, 12)
	server.RegisterObjectUpdateGo("UndeadKillerUpdate", C.nox_xxx_updateUndeadKiller_53E190, func(u *server.Object) { nox_xxx_updateUndeadKiller_53E190(C.int(uintptr(u.CObj()))) }, 0)
	server.RegisterObjectUpdateGo("HarpoonUpdate", C.nox_xxx_updateHarpoon_54F380, func(u *server.Object) { nox_xxx_updateHarpoon_54F380(asObjectC(u)) }, 4)
	server.RegisterObjectUpdateGo("MonsterGeneratorUpdate", C.nox_xxx_updateMonsterGenerator_54E930, func(u *server.Object) { nox_xxx_updateMonsterGenerator_54E930((*C.uint32_t)(u.CObj())) }, 164)

	server.RegisterObjectUpdateParse("PushUpdate", resourceObjectParser("update", "push"))
	server.RegisterObjectUpdateParse("TriggerUpdate", resourceObjectParser("update", "trigger"))
	server.RegisterObjectUpdateParse("ToggleUpdate", resourceObjectParser("update", "trigger"))
	server.RegisterObjectUpdateParse("LoopAndDamageUpdate", resourceObjectParser("update", "triple"))
	server.RegisterObjectUpdateParse("LifetimeUpdate", resourceObjectParser("update", "lifetime"))
	server.RegisterObjectUpdateParse("SkullUpdate", resourceObjectParser("update", "skull"))
}

//export nox_xxx_updatePlayer_4F8100
func nox_xxx_updatePlayer_4F8100(up *nox_object_t) { Nox_xxx_updatePlayer_4F8100(asObjectS(up)) }

//export nox_xxx_objectApplyForce_52DF80
func nox_xxx_objectApplyForce_52DF80(vec *C.float, obj *nox_object_t, force C.float) {
	GetServer().ApplyForce(asObjectS(obj), AsPointf(unsafe.Pointer(vec)), float64(force))
}

//export nox_xxx_updatePixie_53CD20
func nox_xxx_updatePixie_53CD20(cobj *nox_object_t) { Nox_xxx_updatePixie_53CD20(asObjectS(cobj)) }

//export nox_xxx_updatePlayerObserver_4E62F0
func nox_xxx_updatePlayerObserver_4E62F0(a1p *nox_object_t) {
	Nox_xxx_updatePlayerObserver_4E62F0(asObjectS(a1p))
}

//export nox_xxx_updateProjectile_53AC10
func nox_xxx_updateProjectile_53AC10(a1 *nox_object_t) {
	Nox_xxx_updateProjectile_53AC10(asObjectS(a1))
}

//export nox_xxx_updateDeathBall_53D080
func nox_xxx_updateDeathBall_53D080(a1 *nox_object_t) {
	Nox_xxx_updateDeathBall_53D080(asObjectS(a1))
}

//export nox_xxx___mkgmtime_538280
func nox_xxx___mkgmtime_538280(a1 *nox_object_t) {
	Nox_xxx___mkgmtime_538280(asObjectS(a1))
}

func Get_nox_xxx___mkgmtime_538280() unsafe.Pointer {
	return C.nox_xxx___mkgmtime_538280
}

func Nox_server_doPlayersAutoRespawn_40A5F0() int {
	return int(C.int(serverConfigRespawnGet()))
}
func Sub_4E4100() uint32 {
	return uint32(bool2int(questRuntimeRoom()))
}

func Nox_xxx_questCheckSecretArea_421C70(a1 *server.Object) {
	mapPolygonPlayer(a1)
}
func Nox_xxx_playerCanMove_4F9BC0(a1 *server.Object) int {
	return bool2int(controlCanMove(a1))
}
func Sub_4F9AB0(a1 *server.Object) int {
	return int(controlWalkWaypoint(a1))
}
func Nox_xxx_playerConfusedGetDirection_4F7A40(a1 *server.Object) server.Dir16 {
	return server.Dir16(controlConfusedDirection(a1))
}
func Nox_xxx_playerAttack_538960(a1 *server.Object) int {
	return int(nox_xxx_playerAttack_538960(asObjectC(a1)))
}
func Nox_xxx_playerRespawn_4F7EF0(a1 *server.Object) {
	controlRespawn(a1)
}
func Sub_4F9E10(a1 *server.Object) int {
	return int(controlFollowEnemy(a1))
}
func Sub_4F9A80(a1 *server.Object) int {
	return bool2int(controlHasWaypoint(a1))
}
func Nox_xxx_monsterTestBlockShield_533E70(a1 *server.Object) int {
	return int(uintptr(unsafe.Pointer(monsterShieldThreat(a1))))
}
func Nox_common_mapPlrActionToStateId_4FA2B0(a1 *server.Object) int {
	return int(controlActionState(a1))
}
func Nox_xxx_playerCanAttack_4F9C40(a1 *server.Object) int {
	return bool2int(controlCanAttack(a1))
}
func Nox_xxx_checkWinkFlags_4F7DF0(a1 *server.Object) int {
	return int(controlDropBall(a1))
}
func Nox_xxx_playerInputAttack_4F9C70(a1 *server.Object) {
	controlInputAttack(a1)
}
func Nox_xxx_playerSubStamina_4F7D30(a1 *server.Object, a2 int) int {
	return int(controlSubStamina(a1, int32(a2)))
}
func Nox_xxx_playerDoSchedSpell_4FB0E0(a1 *server.Object, a2 *server.Object) {
	controlScheduledSpell(a1, a2, false)
}
func Nox_xxx_playerDoSchedSpellQueue_4FB1D0(a1 *server.Object, a2 *server.Object) {
	controlScheduledSpell(a1, a2, true)
}
func Sub_4E7540(a1 *server.Object, a2 *server.Object) {
	sub_4E7540(asObjectC(a1), asObjectC(a2))
}
func Nox_xxx_playerCheckStrength_4F3180(a1 *server.Object, a2 *server.Object) bool {
	return equipmentCheckStrength(a1, a2)
}
func Nox_xxx_unitDamageClear_4EE5E0(a1 *server.Object, a2 int) {
	resourceDamage(a1, int32(a2))
}
func Nox_xxx_unitAdjustHP_4EE460(a1 *server.Object, a2 int) {
	resourceAdjustHP(a1, int32(a2))
}
func Sub_509CF0(a1 *byte, a2 player.Class, a3 uint32) int {
	return bool2int(matchRosterIdentityAllowed(GoStringP(unsafe.Pointer(a1)), byte(a2), a3))
}
func Sub_4D79C0(a1 *server.Object) {
	questRuntimeReconnect(a1)
}
func Sub_4D7480(a1 *server.Object) {
	questRuntimeGateReturn(a1)
}
