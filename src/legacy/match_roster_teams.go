package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "common__system__team.h"
*/
import "C"

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func matchRosterAssignTeam(pl *server.Player) {
	u := pl.PlayerUnit
	if u == nil {
		return
	}
	s := GetServer().S()
	if C.sub_40A740() == 0 && !noxflags.HasGame(0x8000) {
		if byte(s.Teams.Count()) != 0 {
			tm := (*server.Team)(unsafe.Pointer(C.sub_4189D0()))
			if tm != nil && !u.TeamVal.Has() {
				Nox_xxx_createAtImpl_4191D0(tm.IDVal, u.TeamPtr(), 1, int(u.NetCode), 1)
			}
		}
		return
	}
	if pl.Field2068 == 0 {
		return
	}
	tm := s.Teams.ByXxx(int(pl.Field2068))
	if tm == nil {
		capacity := int(memmap.Uint8(0x5D4594, 371516+52))
		if (noxflags.HasGame(96) || noxflags.HasGame(16) && noxflags.HasGamePlay(4)) && capacity > 2 {
			capacity = 2
		}
		count := func() int { return int(byte(C.sub_417DE0())) }
		if count() >= capacity {
			return
		}
		if noxflags.HasGame(96) && count() >= int(C.sub_417DC0()) {
			return
		}
		tm = (*server.Team)(unsafe.Pointer(C.sub_418A10()))
		if tm == nil {
			return
		}
		C.sub_418800((*C.wchar2_t)(tm.C()), (*C.wchar2_t)(unsafe.Pointer(&pl.Field2072[0])), 0)
		C.sub_418830(C.int(uintptr(tm.C())), C.int(pl.Field2068))
		C.sub_4184D0((*C.nox_team_t)(tm.C()))
	}
	if u.TeamVal.Has() {
		C.sub_4196D0(u.TeamPtr().C(), tm.C(), C.int(u.NetCode), 0)
	} else {
		Nox_xxx_createAtImpl_4191D0(tm.IDVal, u.TeamPtr(), 1, int(u.NetCode), 0)
	}
}
