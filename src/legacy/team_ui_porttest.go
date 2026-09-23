//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Thin accessors call the native UI without duplicating its algorithms.
func PortTestTeamUI(op string, w *gui.Window, tm *server.Team, code, value int, text string) int {
	switch op {
	case "ctf-construct":
		return teamUICTFConstruct()
	case "ctf-show":
		return teamUIHUDShow(false, value)
	case "ctf-open":
		return int(teamUICTFOpen(byte(value)))
	case "ctf-hide":
		return teamUIHUDHide(false)
	case "ctf-select":
		return teamUICTFSelect(byte(code))
	case "ctf-tooltip":
		teamUICTFTooltip(byte(code), byte(value))
	case "ctf-destroy":
		return teamUIHUDDestroy(false)
	case "ball-construct":
		return teamUIBallConstruct()
	case "ball-show":
		return teamUIHUDShow(true, value)
	case "ball-open":
		return teamUIBallOpen()
	case "ball-hide":
		return teamUIHUDHide(true)
	case "ball-destroy":
		return teamUIHUDDestroy(true)
	case "players-construct":
		return teamUIPlayersConstruct((*gui.Window)(unsafe.Pointer(uintptr(code))))
	case "players-refresh":
		return teamUIPlayersRefresh()
	case "players-refresh-if-open":
		return teamUIPlayersRefreshOpen()
	case "players-destroy":
		return int(teamUIPlayersDestroy(value != 0))
	case "player-add":
		return teamUIPlayerAdd(code, text)
	case "player-remove":
		return teamUIPlayerRemove(code)
	case "player-find":
		return teamUIPlayerFind(code, value != 0)
	case "player-team":
		return teamUIPlayerTeam(code, value)
	case "team-add":
		return teamUITeamAdd(tm.Name())
	case "team-color":
		return int(teamUITeamColor(tm))
	case "team-rename":
		return teamUITeamRename(tm, text)
	case "team-remove":
		return teamUITeamRemove(text)
	case "team-find":
		return teamUITeamFind(text, value != 0)
	case "team-clear":
		return teamUITeamClear()
	case "team-join":
		teamUIJoin(tm)
	case "requests-reset":
		teamUIRequestsReset()
	case "map-ctf":
		return teamUIMapCTF()
	case "map-ball":
		return int(teamUIMapBall())
	default:
		panic(op)
	}
	return 0
}
func PortTestTeamUIDraw(op string, w *gui.Window) int {
	switch op {
	case "ctf":
		return teamUICTFDraw(w, w.DrawData())
	case "ball":
		return teamUIBallDraw(w, w.DrawData())
	case "players":
		return teamUIPlayersDraw(w, w.DrawData())
	default:
		panic(op)
	}
}
func PortTestTeamUIEvent(w, child *gui.Window, event, value int) int {
	return teamUIPlayersEvent(w, event, child, value)
}
func PortTestTeamUISelectedName(index int) string { return teamUISelectedName(index) }

// Extracted window pointers are independent of the old backing-blob slots.
func PortTestTeamUIWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"ball-start-type": (*uint32)(unsafe.Pointer(&dword_5d4594_527656)),
		"ctf":             teamUIWord(1045604),
		"ball":            (*uint32)(unsafe.Pointer(&dword_5d4594_1045636)),
		"ball-visible":    teamUIWord(1045640),
		"players":         teamUIWord(1045684),
		"join":            teamUIWord(1045688),
		"rename":          teamUIWord(1045692),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	region := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1045608)), 92)
	saved := append([]byte(nil), region...)
	clear(region)
	listClear((*legacyListNode)(memmap.PtrOff(0x5D4594, 1045652)))
	listClear((*legacyListNode)(memmap.PtrOff(0x5D4594, 1045668)))
	return words, func() {
		teamUIPlayersDestroy(true)
		teamUIHUDDestroy(false)
		teamUIHUDDestroy(true)
		copy(region, saved)
		for k, p := range words {
			*p = old[k]
		}
	}
}

// Snapshot the real intrusive row storage, without reproducing list UI behavior.
type PortTestTeamUIRow struct {
	Name            string
	Code, TeamColor uint32
	Palette         byte
}

func PortTestTeamUIRows(teams bool) []PortTestTeamUIRow {
	off := uintptr(1045652)
	if teams {
		off = 1045668
	}
	var rows []PortTestTeamUIRow
	for n := listNext((*legacyListNode)(memmap.PtrOff(0x5D4594, off))); n != nil; n = listNext(n) {
		if len(rows) >= 1024 {
			panic("team UI row cycle")
		}
		p := unsafe.Pointer(n)
		rows = append(rows, PortTestTeamUIRow{Name: alloc.GoString16((*uint16)(unsafe.Add(p, 12))), Code: *(*uint32)(unsafe.Add(p, 60)), Palette: *(*byte)(unsafe.Add(p, 64)), TeamColor: *(*uint32)(unsafe.Add(p, 68))})
	}
	return rows
}
