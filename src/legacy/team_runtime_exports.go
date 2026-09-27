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
