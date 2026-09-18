//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/server"

func PortTestPlayerStateQuery(op string, pl *server.Player, tm *server.Team) int {
	switch op {
	case "competitors":
		return playerStateCompetitors()
	case "team-players":
		return playerStateTeamPlayers(tm)
	case "multiple":
		return playerStateMultiple()
	case "admission":
		return playerStateAdmission(pl)
	case "elapsed":
		return playerStateElapsed()
	case "threshold":
		return int(playerStateThreshold())
	default:
		panic(op)
	}
}
func PortTestPlayerStateReset() { playerStateReset() }
func PortTestPlayerStateStatus(op string, pl *server.Player, flags uint32) {
	switch op {
	case "add":
		playerStateAddStatus(pl, flags)
	case "remove":
		playerStateRemoveStatus(pl, flags)
	case "report":
		playerStateReport(pl)
	case "all":
		playerStateAllStatus(pl)
	case "lessons":
		playerStateLessons(int32(flags))
	default:
		panic(op)
	}
}
func PortTestPlayerStateName(name *uint16) *server.Player { return playerStateByName(name) }
func PortTestPlayerStateMinimap(op string, obj *server.Object, index int, flags uint32) int {
	switch op {
	case "tracks":
		return playerStateTracks(index, obj)
	case "mark":
		playerStateMark(obj, flags)
	case "unmark":
		playerStateUnmark(obj, flags)
	default:
		panic(op)
	}
	return 0
}
func PortTestPlayerStateEquipment(op string, pl *server.Player, command byte, code, mask uint32, mods *[4]byte) {
	switch op {
	case "respawn":
		playerStateRespawn(pl, command)
	case "equip":
		playerStateEquip(command, code, mask, mods)
	case "unequip":
		playerStateUnequip(command, code, mask)
	default:
		panic(op)
	}
}
func PortTestPlayerStateReentry(value int32) int32 { return playerStateReentry(value) }
