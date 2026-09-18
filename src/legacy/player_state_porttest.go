//go:build porttest

package legacy

/*
#include "GAME1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestPlayerStateQuery(op string, pl *server.Player, tm *server.Team) int {
	switch op {
	case "competitors":
		return int(C.sub_40A770())
	case "team-players":
		return int(C.nox_xxx_countNonEliminatedPlayersInTeam_40A830((*C.nox_team_t)(unsafe.Pointer(tm))))
	case "multiple":
		return int(C.nox_xxx_gamePlayIsAnyPlayers_40A8A0())
	case "admission":
		return int(C.sub_40AA70((*C.nox_playerInfo)(unsafe.Pointer(pl))))
	case "elapsed":
		return int(C.sub_40AA00())
	case "threshold":
		return int(C.sub_40AA40())
	default:
		panic(op)
	}
}

func PortTestPlayerStateReset() { C.sub_40A970() }

func PortTestPlayerStateStatus(op string, pl *server.Player, flags uint32) {
	p := (*C.nox_playerInfo)(unsafe.Pointer(pl))
	switch op {
	case "add":
		C.nox_xxx_netNeedTimestampStatus_4174F0(p, C.int(flags))
	case "remove":
		C.nox_xxx_playerUnsetStatus_417530(p, C.int(flags))
	case "report":
		C.nox_xxx_netReportPlayerStatus_417630(p)
	case "all":
		C.nox_xxx_sendAllClientStatus_4175C0(C.int(uintptr(unsafe.Pointer(pl))))
	case "lessons":
		C.nox_xxx_playerForceSendLessons_416E50(C.int(flags))
	default:
		panic(op)
	}
}
func PortTestPlayerStateName(name *uint16) *server.Player {
	return (*server.Player)(unsafe.Pointer(C.nox_xxx_playerByName_4170D0((*C.wchar2_t)(unsafe.Pointer(name)))))
}
func PortTestPlayerStateMinimap(op string, obj *server.Object, index int, flags uint32) int {
	switch op {
	case "tracks":
		return int(C.nox_xxx_playerMapTracksObj_4173D0(C.int(index), (*C.nox_object_t)(unsafe.Pointer(obj))))
	case "mark":
		C.nox_xxx_netMarkMinimapForAll_4174B0(C.int(uintptr(unsafe.Pointer(obj))), C.int(flags))
	case "unmark":
		C.nox_xxx_netUnmarkMinimapSpec_417470(C.int(uintptr(unsafe.Pointer(obj))), C.int(flags))
	default:
		panic(op)
	}
	return 0
}

func PortTestPlayerStateEquipment(op string, pl *server.Player, command byte, code, mask uint32, mods *[4]byte) {
	switch op {
	case "respawn":
		C.nox_xxx_cliPlayerRespawn_417680(C.int(uintptr(unsafe.Pointer(pl))), C.char(command))
	case "equip":
		C.nox_xxx_clientEquipWeaponArmor_417AA0(C.char(command), C.int(code), C.int(mask), C.int(uintptr(unsafe.Pointer(mods))))
	case "unequip":
		C.sub_417B80(C.char(command), C.int(code), C.int(mask))
	default:
		panic(op)
	}
}

func PortTestPlayerStateReentry(value int32) int32 { return int32(C.sub_40AA60(C.int(value))) }
