package legacy

/*
#include "GAME3_2.h"
typedef const nox_object_t gameplay_report_const_object;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func nox_xxx_netNotifyRate_4D7F10(a0 C.int) C.int { return C.int(gameplayReportRate(int(a0))) }

func sub_4D81A0(a0 C.int) {
	gameplayReportExperience((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))))
}

func sub_4D82F0(a0 C.int, a1 *C.uint32_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(gameplayReportEquipment(int(a0), (*server.Object)(unsafe.Pointer(a1))))))
}

func nox_xxx_netReportUnitCurrentHP_4D8620(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportCurrentHealth(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

func nox_xxx_netFlagballWinner_4D8C40(a0 C.int) C.int {
	return C.int(gameplayReportFlagballWinner((*server.Team)(unsafe.Pointer(uintptr(uint32(a0))))))
}

func nox_xxx_netFlagWinner_4D8C40_4D8C80(a0 C.int, a1 C.char) C.int {
	return C.int(gameplayReportFlagWinner((*server.Team)(unsafe.Pointer(uintptr(uint32(a0)))), byte(a1)))
}

func nox_xxx_playerIncrementElimDeath_4D8D40(a0 C.int) {
	gameplayReportEliminationDeath((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))))
}

func nox_xxx_changeScore_4D8E90(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportChangeScore((*server.Object)(unsafe.Pointer(uintptr(uint32(a0)))), int(a1)))
}

func nox_xxx_netReportLesson_4D8EF0(a0 *C.nox_object_t) C.int {
	return C.int(gameplayReportLesson((*server.Object)(unsafe.Pointer(a0))))
}

func nox_xxx_netReportEnchant_4D8F90(a0 C.int, a1 *C.uint32_t) C.int {
	return C.int(gameplayReportEnchant(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

func nox_xxx_netReportUnitHeight_4D9020(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportHeight(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

func nox_xxx_netReportAcquireCreature_4D91A0(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportAcquireCreature(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}

func nox_xxx_netMonitorCreature_4D9250(a0 C.int, a1 C.int) C.int {
	return C.int(gameplayReportMonitor(int(a0), (*server.Object)(unsafe.Pointer(uintptr(uint32(a1))))))
}

func nox_xxx_netReportSpellStat_4D9630(a0 C.int, a1 C.int, a2 C.char) C.int {
	return C.int(gameplayReportSpellStat(int(a0), uint32(a1), byte(a2)))
}

func nox_xxx_netSendSecondaryWeapon_4D9670(a0 C.int, a1 *C.uint32_t, a2 C.char) C.int {
	return C.int(gameplayReportSecondary(int(a0), (*server.Object)(unsafe.Pointer(a1)), byte(a2)))
}

func sub_4D9CF0(a0 C.int) C.int { return C.int(gameplayReportQuestStart(int(a0))) }

func sub_4D9D20(a0 C.int, a1 *C.nox_object_t) C.int {
	return C.int(gameplayReportQuestObject(int(a0), (*server.Object)(unsafe.Pointer(a1))))
}
