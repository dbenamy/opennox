//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Fixtures dispatch directly to the production Go owners.
func PortTestTeamRuntimeMessage(op string, tm *server.Team, obj *server.ObjectTeam, code int, name string) int {
	switch op {
	case "rename":
		teamRuntimeRename(tm, (*uint16)(unsafe.Pointer(internWStr(name))))
	case "change":
		return teamRuntimeSwitch(obj, tm, code, 0)
	case "join-request":
		return int(teamRuntimeRequest(tm, obj, int16(code), 10))
	case "change-request":
		return int(teamRuntimeRequest(tm, obj, int16(code), 11))
	default:
		panic(op)
	}
	return 0
}
func PortTestTeamRuntimeName(tm *server.Team, name string, flag uint32) {
	teamRuntimeSetName(tm, (*uint16)(unsafe.Pointer(internWStr(name))), flag)
}
func PortTestTeamRuntimeUnlink(tm *server.Team, obj *server.ObjectTeam) { teamRuntimeUnlink(tm, obj) }
func PortTestTeamRuntimePredicates(a, b *server.ObjectTeam, id byte) (int, int, int) {
	return int(teamRuntimeBool(a.Has())), int(teamRuntimeBool(a.SameAs(b))), int(teamRuntimeBool(teamRuntimeContains(a, server.TeamID(id))))
}
func PortTestTeamRuntimeLookup(name string) (*server.Team, bool) {
	t := teamRuntimeFind((*uint16)(unsafe.Pointer(internWStr(name))))
	return t, t != nil
}
func PortTestTeamRuntimeSelect(op string, tm *server.Team) int {
	switch op {
	case "count":
		return teamRuntimeCount(tm)
	case "group-count":
		return int(teamRuntimeGroupCount())
	case "least":
		if t := teamRuntimeLeast(); t != nil {
			return int(t.ID())
		}
		return 0
	case "available":
		if t := teamRuntimeAvailable(); t != nil {
			return int(t.ID())
		}
		return 0
	case "flag-count":
		return int(teamRuntimeFlagCount)
	case "capflag":
		return teamUIMapCTF()
	case "flagball":
		return int(teamUIMapBall())
	default:
		panic(op)
	}
}
func PortTestTeamRuntimeMap(op string, arg int) int {
	switch op {
	case "scan":
		return int(teamRuntimeBool(teamRuntimeScan()))
	case "assign":
		return teamRuntimeAssignFlags()
	case "toggle":
		teamRuntimeToggle(arg != 0)
		return 0
	case "create":
		return teamRuntimeCreateMap()
	case "balance":
		teamRuntimeBalance(arg != 0)
		return 0
	case "nearest":
		return teamRuntimeNearest()
	case "enable":
		return teamRuntimeEnable()
	case "crown":
		return teamRuntimeCrown()
	default:
		panic(op)
	}
}
func PortTestTeamRuntimeTextIndex() *uint32 {
	return (*uint32)(unsafe.Pointer(&interactionMessageHead))
}
func PortTestTeamRuntimeOther(op string, tm *server.Team, member *server.ObjectTeam, a, b int) int {
	switch op {
	case "group":
		teamRuntimeSetGroup(tm, uint32(a))
		return int(uintptr(tm.C()))
	case "clear":
		teamRuntimeClearPlayers(tm)
	case "lessons":
		teamRuntimeLessons(tm, a)
	case "leave":
		teamRuntimeLeave(member, a)
	case "announce":
		teamRuntimeAnnounce(tm)
	case "describe":
		teamRuntimeDescribe(tm, a)
	case "member-report":
		teamRuntimeMemberReport(member, a, b)
	case "reset":
		return GetServer().TeamsRemoveActive(false)
	default:
		panic(op)
	}
	return 0
}
func PortTestTeamRuntimeLinks(tm *server.Team, member *server.ObjectTeam) (unsafe.Pointer, unsafe.Pointer) {
	return teamRuntimeFirst(tm).C(), teamRuntimeNext(member).C()
}
func PortTestTeamRuntimeObject(code int) *server.ObjectTeam { return teamRuntimeObject(code) }

// Shared C respawn ownership, not a second implementation of its list logic.
func PortTestTeamRuntimeRespawns(objects []*server.Object) func() {
	oldPool, oldHead, oldAllow := itemRespawnPool, itemRespawnHead, itemRespawnEnabled
	itemRespawnInit()
	itemRespawnReset()
	for _, u := range objects {
		itemRespawnAdd(u)
	}
	return func() {
		itemRespawnFree()
		itemRespawnPool, itemRespawnHead, itemRespawnEnabled = oldPool, oldHead, oldAllow
	}
}
