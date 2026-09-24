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

func sub_418800(t, name *C.wchar2_t, flag C.int) {
	teamRuntimeSetName((*server.Team)(unsafe.Pointer(t)), (*uint16)(unsafe.Pointer(name)), uint32(flag))
}

func sub_418830(t, group C.int) C.int {
	teamRuntimeSetGroup(teamRuntimeTeamInt(t), uint32(group))
	return t
}

//export nox_xxx_objGetTeamByNetCode_418C80
func nox_xxx_objGetTeamByNetCode_418C80(code C.int) *C.uint32_t {
	return (*C.uint32_t)(teamRuntimeObject(int(code)).C())
}

func nox_xxx_teamRenameMB_418CD0(t, name *C.wchar2_t) {
	teamRuntimeRename((*server.Team)(unsafe.Pointer(t)), (*uint16)(unsafe.Pointer(name)))
}

func sub_418D80(t C.int) { teamRuntimeClearPlayers(teamRuntimeTeamInt(t)) }

func nox_xxx_netChangeTeamID_419090(t, score C.int) {
	teamRuntimeLessons(teamRuntimeTeamInt(t), int(score))
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

func nox_xxx_netChangeTeamMb_419570(m unsafe.Pointer, code C.int) {
	teamRuntimeLeave(teamRuntimeMember(m), int(code))
}

func sub_4196D0(m, t unsafe.Pointer, code, relocate C.int) C.int {
	return C.int(teamRuntimeSwitch(teamRuntimeMember(m), asTeamP(t), int(code), int(relocate)))
}

func sub_419900(t, m C.int, code C.short) C.char {
	return C.char(teamRuntimeRequest(teamRuntimeTeamInt(t), teamRuntimeMemberInt(m), int16(code), 10))
}

func nox_xxx_createAtImpl_4191D0(id C.uchar, m unsafe.Pointer, notify, code, relocate C.int) {
	teamRuntimeJoin(server.TeamID(id), teamRuntimeMember(m), int(notify), int(code), int(relocate))
}
