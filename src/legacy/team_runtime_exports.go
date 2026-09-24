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
