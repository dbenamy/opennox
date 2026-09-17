//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_3484;
*/
import "C"

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestMatchRosterSendPlayers(to int) { matchRosterSendPlayers(to) }
func PortTestMatchRosterSendSettings()      { matchRosterSettings() }
func PortTestMatchRosterOwnCache() func() {
	old := matchRosterFlagType
	return func() { matchRosterFlagType = old }
}
func PortTestMatchRosterPlayerPacket(buf []byte, pl *server.Player) {
	server.EncodePlayerRoster(buf, pl)
}
func PortTestMatchRoster(op string, u *server.Object, pl *server.Player, data unsafe.Pointer, a [4]uint32) uint32 {
	switch op {
	case "gui-settings":
		return uint32(matchRosterGUISettings(byte(a[0]), data, int(a[1])))
	case "inventory":
		matchRosterInventory(int(a[0]), pl)
	case "player-ids":
		return uint32(matchRosterPlayerIDs(pl))
	case "objective-minimap":
		return uint32(matchRosterMinimap(int(a[0])))
	case "object-report-mask":
		matchRosterReportMask(int(a[0]))
	case "object-resync-mask":
		return uint32(matchRosterResyncMask(byte(a[0])))
	case "object-clear-mask":
		return uint32(matchRosterClearMask(byte(a[0])))
	case "wall-open":
		return uint32(matchRosterWall((*server.Wall)(data), 59))
	case "wall-close":
		return uint32(matchRosterWall((*server.Wall)(data), 60))
	case "team-roster":
		matchRosterTeamRoster(int(a[0]))
	case "simple-object":
		return uint32(matchRosterSimpleObject(int(a[0]), u))
	case "assign-team":
		matchRosterAssignTeam(pl)
	case "flag-state":
		return uint32(matchRosterFlagState(byte(a[0]), byte(a[1]), byte(a[2]), uint16(a[3])))
	case "winner-min":
		return uint32(matchRosterWinner(true))
	case "winner-max":
		return uint32(matchRosterWinner(false))
	case "winner-flag":
		return uint32(matchRosterFlagWinner())
	case "check-victory":
		matchRosterCheckVictory()
	case "check-limit":
		return uint32(matchRosterCheckLimit())
	case "remember":
		matchRosterRemember(pl)
	case "remembered":
		return uint32(bool2int(matchRosterHasIdentity(pl)))
	case "forget":
		return uint32(bool2int(matchRosterForget()))
	default:
		panic(op)
	}
	return 0
}
func PortTestMatchRosterQuery(name string, class int, group uint32) uint32 {
	return uint32(bool2int(matchRosterIdentityAllowed(name, byte(class), group)))
}
func PortTestMatchRosterFlagPointers(index byte) (unsafe.Pointer, unsafe.Pointer) {
	return matchRosterFlagBase(), matchRosterFlagRecord(index)
}
func PortTestMatchRosterGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"flag-type": &matchRosterFlagType, "remembered-initialized": &matchRosterRememberedInit, "team-cap": &teamRuntimeFlagCount, "server-subflags": (*uint32)(unsafe.Pointer(&C.dword_5d4594_3484))}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	remembered := matchRosterRemembered
	matchRosterRemembered = nil
	return words, func() {
		matchRosterRemembered = remembered
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestMatchRosterRemembered() [][]byte {
	var out [][]byte
	for _, id := range matchRosterRemembered {
		b := make([]byte, 20)
		copy(b[:12], id.Name)
		binary.LittleEndian.PutUint32(b[12:], id.Group)
		b[16] = id.Class
		out = append(out, b)
	}
	return out
}
func PortTestMatchRosterFlagWinner() uint32 { return uint32(matchRosterFlagWinner()) }
