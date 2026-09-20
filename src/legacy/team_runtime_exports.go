package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "common__system__team.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func teamRuntimeTeamInt(p C.int) *server.Team {
	return (*server.Team)(unsafe.Pointer(uintptr(uint32(p))))
}
func teamRuntimeMemberInt(p C.int) *server.ObjectTeam {
	return (*server.ObjectTeam)(unsafe.Pointer(uintptr(uint32(p))))
}
func teamRuntimeBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

//export sub_418800
func sub_418800(t, name *C.wchar2_t, flag C.int) {
	teamRuntimeSetName((*server.Team)(unsafe.Pointer(t)), (*uint16)(unsafe.Pointer(name)), uint32(flag))
}

//export sub_418830
func sub_418830(t, group C.int) C.int {
	teamRuntimeSetGroup(teamRuntimeTeamInt(t), uint32(group))
	return t
}

//export sub_418A40
func sub_418A40(name *C.wchar2_t) *C.char {
	return (*C.char)(teamRuntimeFind((*uint16)(unsafe.Pointer(name))).C())
}

//export sub_418BC0
func sub_418BC0(t C.int) C.int { return C.int(teamRuntimeCount(teamRuntimeTeamInt(t))) }

//export nox_xxx_objGetTeamByNetCode_418C80
func nox_xxx_objGetTeamByNetCode_418C80(code C.int) *C.uint32_t {
	return (*C.uint32_t)(teamRuntimeObject(int(code)).C())
}

//export nox_xxx_teamRenameMB_418CD0
func nox_xxx_teamRenameMB_418CD0(t, name *C.wchar2_t) {
	teamRuntimeRename((*server.Team)(unsafe.Pointer(t)), (*uint16)(unsafe.Pointer(name)))
}

//export sub_418D80
func sub_418D80(t C.int) { teamRuntimeClearPlayers(teamRuntimeTeamInt(t)) }

//export nox_xxx_netChangeTeamID_419090
func nox_xxx_netChangeTeamID_419090(t, score C.int) {
	teamRuntimeLessons(teamRuntimeTeamInt(t), int(score))
}

//export sub_4190F0
func sub_4190F0(name *C.wchar2_t) C.int {
	return teamRuntimeBool(teamRuntimeFind((*uint16)(unsafe.Pointer(name))) != nil)
}

//export nox_xxx_servObjectHasTeam_419130
func nox_xxx_servObjectHasTeam_419130(m C.int) C.int {
	return teamRuntimeBool(teamRuntimeMemberInt(m).Has())
}

//export nox_xxx_servCompareTeams_419150
func nox_xxx_servCompareTeams_419150(a, b C.int) C.int {
	return teamRuntimeBool(teamRuntimeMemberInt(a).SameAs(teamRuntimeMemberInt(b)))
}

func nox_xxx_teamCompare2_419180(m unsafe.Pointer, id C.uchar) C.int {
	return teamRuntimeBool(teamRuntimeContains(teamRuntimeMember(m), server.TeamID(id)))
}

//export nox_xxx_netChangeTeamMb_419570
func nox_xxx_netChangeTeamMb_419570(m unsafe.Pointer, code C.int) {
	teamRuntimeLeave(teamRuntimeMember(m), int(code))
}

//export sub_4196D0
func sub_4196D0(m, t unsafe.Pointer, code, relocate C.int) C.int {
	return C.int(teamRuntimeSwitch(teamRuntimeMember(m), asTeamP(t), int(code), int(relocate)))
}

//export sub_419900
func sub_419900(t, m C.int, code C.short) C.char {
	return C.char(teamRuntimeRequest(teamRuntimeTeamInt(t), teamRuntimeMemberInt(m), int16(code), 10))
}

//export sub_419960
func sub_419960(t, m C.int, code C.short) C.char {
	return C.char(teamRuntimeRequest(teamRuntimeTeamInt(t), teamRuntimeMemberInt(m), int16(code), 11))
}

//export sub_417DC0
func sub_417DC0() C.int { return C.int(teamRuntimeFlagCount) }

//export sub_417DE0
func sub_417DE0() C.char { return C.char(teamRuntimeGroupCount()) }

//export sub_417EC0
func sub_417EC0() C.bool { return C.bool(teamRuntimeScan()) }

//export sub_4181F0
func sub_4181F0(reset C.int) { teamRuntimeBalance(reset != 0) }

//export sub_418390
func sub_418390() C.int { return C.int(teamRuntimeEnable()) }

//export sub_4184D0
func sub_4184D0(t *C.nox_team_t) { teamRuntimeAnnounce(asTeam(t)) }

//export nox_xxx_wndGuiTeamCreate_4185B0
func nox_xxx_wndGuiTeamCreate_4185B0() C.int { return C.int(teamRuntimeCreateMap()) }

//export nox_xxx_toggleAllTeamFlags_418690
func nox_xxx_toggleAllTeamFlags_418690(on C.int) *C.char { teamRuntimeToggle(on != 0); return nil }

//export nox_xxx_createAtImpl_4191D0
func nox_xxx_createAtImpl_4191D0(id C.uchar, m unsafe.Pointer, notify, code, relocate C.int) {
	teamRuntimeJoin(server.TeamID(id), teamRuntimeMember(m), int(notify), int(code), int(relocate))
}
