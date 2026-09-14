package legacy

/*
#include "GAME3_2.h"
typedef const char gameplay_report_const_char;
typedef const nox_object_t gameplay_report_const_object;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4D7EA0
func sub_4D7EA0() C.int { return C.int(gameplayReportResetAll()) }

//export nox_xxx_netCreatureCmd_4D7EE0
func nox_xxx_netCreatureCmd_4D7EE0(a0 C.int, a1 C.char) C.int {
	return C.int(gameplayReportCreature(int(a0), byte(a1)))
}

//export nox_xxx_netNotifyRate_4D7F10
func nox_xxx_netNotifyRate_4D7F10(a0 C.int) C.int { return C.int(gameplayReportRate(int(a0))) }

//export sub_4D81A0
func sub_4D81A0(a0 C.int) {
	gameplayReportExperience((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))))
}

//export nox_xxx_netReportAnimFrame_4D81F0
func nox_xxx_netReportAnimFrame_4D81F0(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportAnimation(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netReportXStatus_4D8230
func nox_xxx_netReportXStatus_4D8230(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportXStatus(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export sub_4D82F0
func sub_4D82F0(a0 C.int, a1 *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(gameplayReportEquipment(int(a0), (*server.Object)(unsafe.Pointer(a1))))))
}

//export nox_xxx_netReportUnitCurrentHP_4D8620
func nox_xxx_netReportUnitCurrentHP_4D8620(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportCurrentHealth(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netSendTeam_4D8670
func nox_xxx_netSendTeam_4D8670(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportTeam(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netSendPlrHealthToTeam_4D86E0
func nox_xxx_netSendPlrHealthToTeam_4D86E0(a0 C.int) *C.char {
	return (*C.char)(unsafe.Pointer(uintptr(gameplayReportPlayerHealthToTeam(int(a0)))))
}

//export nox_xxx_netReportHealthDelta_4D8760
func nox_xxx_netReportHealthDelta_4D8760(a0 C.int, a1 C.short, a2 C.short) C.short {
	return C.short(gameplayReportHealthDelta(int(a0), uint16(a1), int16(a2)))
}

//export nox_xxx_netReportMana_4D8930
func nox_xxx_netReportMana_4D8930(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportMana(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))
}

//export nox_xxx_netSendDMWinner_4D8B90
func nox_xxx_netSendDMWinner_4D8B90(a0 C.int, a1 C.char) C.int {
	return C.int(gameplayReportDMWinner((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))), byte(a1)))
}

//export nox_xxx_netSendDMTeamWinner_4D8BF0
func nox_xxx_netSendDMTeamWinner_4D8BF0(a0 C.int, a1 C.char) C.int {
	return C.int(gameplayReportDMTeamWinner((*server.Team)(unsafe.Pointer(uintptr(uint32(a0)))), byte(a1)))
}

//export nox_xxx_netFlagballWinner_4D8C40
func nox_xxx_netFlagballWinner_4D8C40(a0 C.int) C.int {
	return C.int(gameplayReportFlagballWinner((*server.Team)(unsafe.Pointer(uintptr(uint32(a0))))))
}

//export nox_xxx_netFlagWinner_4D8C40_4D8C80
func nox_xxx_netFlagWinner_4D8C40_4D8C80(a0 C.int, a1 C.char) C.int {
	return C.int(gameplayReportFlagWinner((*server.Team)(unsafe.Pointer(uintptr(uint32(a0)))), byte(a1)))
}

//export nox_xxx_playerIncrementElimDeath_4D8D40
func nox_xxx_playerIncrementElimDeath_4D8D40(a0 C.int) {
	gameplayReportEliminationDeath((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))))
}

//export nox_xxx_changeScore_4D8E90
func nox_xxx_changeScore_4D8E90(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportChangeScore((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))), int(a1)))
}

//export nox_xxx_playerSubLessons_4D8EC0
func nox_xxx_playerSubLessons_4D8EC0(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportSubtractLessons((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))), int(a1)))
}

//export nox_xxx_netReportLesson_4D8EF0
func nox_xxx_netReportLesson_4D8EF0(a0 *C.nox_object_t) C.int {
	return C.int(gameplayReportLesson((*server.Object)(unsafe.Pointer(a0))))
}

//export nox_xxx_netTimerStatus_4D8F50
func nox_xxx_netTimerStatus_4D8F50(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportTimer(int(a0), uint32(a1)))
}

//export nox_xxx_netReportEnchant_4D8F90
func nox_xxx_netReportEnchant_4D8F90(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportEnchant(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netReportObjHidden_4D8FD0
func nox_xxx_netReportObjHidden_4D8FD0(a0 C.int, a1 *C.uint32_t) {
	gameplayReportHidden(int(a0), (*server.Object)(unsafe.Pointer(a1)))
}

//export nox_xxx_netReportUnitHeight_4D9020
func nox_xxx_netReportUnitHeight_4D9020(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportHeight(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netReportAcquireCreature_4D91A0
func nox_xxx_netReportAcquireCreature_4D91A0(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportAcquireCreature(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netFxShield_0_4D9200
func nox_xxx_netFxShield_0_4D9200(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportShield(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))
}

//export nox_xxx_netMonitorCreature_4D9250
func nox_xxx_netMonitorCreature_4D9250(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportMonitor(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))
}

//export nox_xxx_netReportTeamBase_4D92D0
func nox_xxx_netReportTeamBase_4D92D0(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportTeamBase(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))
}

//export nox_xxx_netSendReportNPC_4D93A0
func nox_xxx_netSendReportNPC_4D93A0(a0 C.int, a1 C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(gameplayReportNPC(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))))
}

//export nox_xxx_netSendJournalAdd_4D9440
func nox_xxx_netSendJournalAdd_4D9440(a0 C.int, a1 *C.nox_playerInfo_journal) C.int {
	return C.int(gameplayReportJournal(int(a0), unsafe.Pointer(a1), 1))
}

//export nox_xxx_netSendJournalRemove_4D94A0
func nox_xxx_netSendJournalRemove_4D94A0(a0 C.int, a1 *C.gameplay_report_const_char) C.int {
	return C.int(gameplayReportJournal(int(a0), unsafe.Pointer(a1), 2))
}

//export nox_xxx_netSendJournalUpdate_4D9500
func nox_xxx_netSendJournalUpdate_4D9500(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportJournal(int(a0), unsafe.Pointer(uintptr(uint32(a1))), 3))
}

//export nox_xxx_netSendChapterEnd_4D9560
func nox_xxx_netSendChapterEnd_4D9560(a0 C.int, a1 C.char, a2 C.int) C.int {
	return C.int(gameplayReportChapter(int(a0), byte(a1), int(a2)))
}

//export nox_xxx_netSendFlagStatus_4D95A0
func nox_xxx_netSendFlagStatus_4D95A0(a0 C.int, a1 C.char, a2 C.char, a3 C.char, a4 C.short) C.int {
	return C.int(gameplayReportFlag(int(a0), byte(a1), byte(a2), byte(a3), uint16(a4)))
}

//export nox_xxx_netSendBallStatus_4D95F0
func nox_xxx_netSendBallStatus_4D95F0(a0 C.int, a1 C.char, a2 C.short) C.int {
	return C.int(gameplayReportBall(int(a0), byte(a1), uint16(a2)))
}

//export nox_xxx_netReportSpellStat_4D9630
func nox_xxx_netReportSpellStat_4D9630(a0 C.int, a1 C.int, a2 C.char) C.int {
	return C.int(gameplayReportSpellStat(int(a0), uint32(a1), byte(a2)))
}

//export nox_xxx_netSendSecondaryWeapon_4D9670
func nox_xxx_netSendSecondaryWeapon_4D9670(a0 C.int, a1 *C.uint32_t, a2 C.char) C.int {
	return C.int(gameplayReportSecondary(int(a0), (*server.Object)(unsafe.Pointer(a1)), byte(a2)))
}

//export nox_xxx_netMsgLastQuiver_4D96B0
func nox_xxx_netMsgLastQuiver_4D96B0(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportQuiver(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netMsgInventoryLoaded_4D96E0
func nox_xxx_netMsgInventoryLoaded_4D96E0(a0 C.int) C.int {
	return C.int(gameplayReportInventoryLoaded(int(a0)))
}

//export nox_xxx_netFriendAddRemove_4D97A0
func nox_xxx_netFriendAddRemove_4D97A0(a0 C.int, a1 *C.uint32_t, a2 C.int) C.int {
	return C.int(gameplayReportFriend(int(a0), (*server.Object)(unsafe.Pointer(a1)), int(a2)))
}

//export sub_4D97E0
func sub_4D97E0(a0 C.int) C.int { return C.int(gameplayReportFriendReset(int(a0))) }

//export nox_xxx_playerReportAnything_4D9900
func nox_xxx_playerReportAnything_4D9900(a0 C.int) {
	gameplayReportAnything((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))))
}

//export sub_4D9CF0
func sub_4D9CF0(a0 C.int) C.int { return C.int(gameplayReportQuestStart(int(a0))) }

//export sub_4D9D20
func sub_4D9D20(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportQuestObject(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_netGauntlet_4D9E70
func nox_xxx_netGauntlet_4D9E70(a0 C.int) C.int { return C.int(gameplayReportGauntlet(int(a0))) }
