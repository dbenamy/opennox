package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func teamRuntimeTeamInt(p int32) *server.Team {
	return (*server.Team)(unsafe.Pointer(uintptr(uint32(p))))
}
func teamRuntimeMemberInt(p int32) *server.ObjectTeam {
	return (*server.ObjectTeam)(unsafe.Pointer(uintptr(uint32(p))))
}
func teamRuntimeBool(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
