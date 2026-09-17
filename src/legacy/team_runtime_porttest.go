//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "GAME1.h"
extern uint32_t dword_5d4594_825736;
#include "GAME3_3.h"
extern void* nox_alloc_respawn_1568020;
extern uint32_t dword_5d4594_1568024;
extern uint32_t nox_xxx_respawnAllow_587000_205200;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Call the actual team implementation; fixtures own teams and message queues.
func PortTestTeamRuntimeMessage(op string, tm *server.Team, obj *server.ObjectTeam, code int, name string) int {
	switch op {
	case "rename":
		C.nox_xxx_teamRenameMB_418CD0((*C.wchar2_t)(tm.C()), internWStr(name))
	case "change":
		return int(C.sub_4196D0(obj.C(), tm.C(), C.int(code), 0))
	case "join-request":
		return int(C.sub_419900(C.int(uintptr(tm.C())), C.int(uintptr(unsafe.Pointer(obj))), C.short(code)))
	case "change-request":
		return int(C.sub_419960(C.int(uintptr(tm.C())), C.int(uintptr(unsafe.Pointer(obj))), C.short(code)))
	default:
		panic(op)
	}
	return 0
}

func PortTestTeamRuntimeName(tm *server.Team, name string, flag uint32) {
	C.sub_418800((*C.wchar2_t)(tm.C()), internWStr(name), C.int(flag))
}
func PortTestTeamRuntimeUnlink(tm *server.Team, obj *server.ObjectTeam) {
	C.sub_418E40(tm.C(), obj.C())
}

func PortTestTeamRuntimePredicates(a, b *server.ObjectTeam, id byte) (int, int, int) {
	return int(C.nox_xxx_servObjectHasTeam_419130(C.int(uintptr(a.C())))), int(C.nox_xxx_servCompareTeams_419150(C.int(uintptr(a.C())), C.int(uintptr(b.C())))), int(C.nox_xxx_teamCompare2_419180(a.C(), C.uchar(id)))
}
func PortTestTeamRuntimeLookup(name string) (*server.Team, bool) {
	return (*server.Team)(unsafe.Pointer(C.sub_418A40(internWStr(name)))), C.sub_4190F0(internWStr(name)) != 0
}
func PortTestTeamRuntimeSelect(op string, tm *server.Team) int {
	switch op {
	case "count":
		return int(C.sub_418BC0(C.int(uintptr(tm.C()))))
	case "group-count":
		return int(C.sub_417DE0())
	case "least":
		t := (*server.Team)(unsafe.Pointer(C.sub_4189D0()))
		if t == nil {
			return 0
		}
		return int(t.ID())
	case "available":
		t := (*server.Team)(unsafe.Pointer(C.sub_418A10()))
		if t == nil {
			return 0
		}
		return int(t.ID())
	case "flag-count":
		return int(C.sub_417DC0())
	case "capflag":
		return int(C.nox_xxx_mapInfoSetCapflag_417EA0())
	case "flagball":
		return int(C.nox_xxx_mapInfoSetFlagball_417F30())
	default:
		panic(op)
	}
}

func PortTestTeamRuntimeMap(op string, arg int) int {
	switch op {
	case "scan":
		if C.sub_417EC0() {
			return 1
		}
		return 0
	case "assign":
		return int(C.nox_xxx_teamAssignFlags_418640())
	case "toggle":
		return int(uintptr(unsafe.Pointer(C.nox_xxx_toggleAllTeamFlags_418690(C.int(arg)))))
	case "create":
		return int(C.nox_xxx_wndGuiTeamCreate_4185B0())
	case "balance":
		C.sub_4181F0(C.int(arg))
		return 0
	case "nearest":
		return int(C.sub_4183C0())
	case "enable":
		return int(C.sub_418390())
	case "crown":
		return int(C.nox_xxx_mapInfoSetKotr_4180D0())
	default:
		panic(op)
	}
}

func PortTestTeamRuntimeTextIndex() *uint32 { return (*uint32)(unsafe.Pointer(&C.dword_5d4594_825736)) }
func PortTestTeamRuntimeOther(op string, tm *server.Team, member *server.ObjectTeam, a, b int) int {
	switch op {
	case "group":
		return int(C.sub_418830(C.int(uintptr(tm.C())), C.int(a)))
	case "clear":
		C.sub_418D80(C.int(uintptr(tm.C())))
	case "lessons":
		C.nox_xxx_netChangeTeamID_419090(C.int(uintptr(tm.C())), C.int(a))
	case "leave":
		C.nox_xxx_netChangeTeamMb_419570(member.C(), C.int(a))
	case "announce":
		C.sub_4184D0((*C.nox_team_t)(tm.C()))
	case "describe":
		C.sub_4197C0((*C.wchar2_t)(tm.C()), C.int(a))
	case "member-report":
		C.sub_4198A0(C.int(uintptr(member.C())), C.int(a), C.int(b))
	case "reset":
		return int(C.sub_417CF0())
	default:
		panic(op)
	}
	return 0
}
func PortTestTeamRuntimeLinks(tm *server.Team, member *server.ObjectTeam) (unsafe.Pointer, unsafe.Pointer) {
	return unsafe.Pointer(uintptr(C.nox_xxx_teamCheckSmth_418C60(C.int(uintptr(tm.C()))))), unsafe.Pointer(C.sub_418C70((*C.uint32_t)(member.C())))
}
func PortTestTeamRuntimeObject(code int) *server.ObjectTeam {
	return (*server.ObjectTeam)(unsafe.Pointer(C.nox_xxx_objGetTeamByNetCode_418C80(C.int(code))))
}

// Shared C respawn ownership, not a second implementation of its list logic.
func PortTestTeamRuntimeRespawns(objects []*server.Object) func() {
	oldPool, oldHead, oldAllow := C.nox_alloc_respawn_1568020, C.dword_5d4594_1568024, C.nox_xxx_respawnAllow_587000_205200
	C.nox_alloc_respawn_1568020 = nil
	C.dword_5d4594_1568024 = 0
	if C.nox_xxx_allocItemRespawnArray_4ECA60() == 0 {
		panic("respawn owner allocation")
	}
	C.sub_4EC5B0()
	for _, u := range objects {
		C.nox_xxx_respawnAdd_4EC5E0(asObjectC(u))
	}
	return func() {
		C.sub_4ECA90()
		C.nox_alloc_respawn_1568020, C.dword_5d4594_1568024, C.nox_xxx_respawnAllow_587000_205200 = oldPool, oldHead, oldAllow
	}
}
