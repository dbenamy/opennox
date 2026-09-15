//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// PortTestScoreboard invokes the production owners directly. Row formatting is
// exercised through its real rendering callers, with no test-only algorithm.
func PortTestScoreboard(op int, a, b, c uintptr) uint32 {
	ptr := func(p *uint16) uint32 { return uint32(uintptr(unsafe.Pointer(p))) }
	signedByte := func(v byte) uint32 { return uint32(int32(int8(v))) }
	switch op {
	case 0:
		return uint32(scoreboardInsertSafe((*gui.Window)(unsafe.Pointer(a)), byte(b), (*uint16)(unsafe.Pointer(c))))
	case 1:
		return uint32(uintptr(unsafe.Pointer(scoreboardConstruct())))
	case 2:
		return ptr(scoreboardLoadClasses())
	case 3:
		return uint32(scoreboardDraw((*gui.Window)(unsafe.Pointer(a)), (*gui.WindowData)(unsafe.Pointer(b))))
	case 4:
		return uint32(scoreboardHeadings(int(a), int(b)))
	case 5:
		return ptr(scoreboardStatus(int(a), (*byte)(unsafe.Pointer(b))))
	case 6:
		return signedByte(scoreboardTimeHeading())
	case 7:
		return uint32(scoreboardLessonHeading())
	case 8:
		return uint32(scoreboardClearRows())
	case 9:
		return uint32(scoreboardInsert((*gui.Window)(unsafe.Pointer(a)), byte(b), (*uint16)(unsafe.Pointer(c))))
	case 11:
		scoreboardCollect()
		return 0
	case 12:
		return uint32(scoreboardObjective((*server.Player)(unsafe.Pointer(a))))
	case 13:
		return uint32(bool2int(scoreboardHasTeam(uint32(a))))
	case 14:
		return ptr(scoreboardClipName((*uint16)(unsafe.Pointer(a))))
	case 15:
		return uint32(bool2int(scoreboardHasPlayer(uint32(a))))
	case 16:
		scoreboardCollectOrdinary()
		return 0
	case 17, 18:
		return uint32(gui.EventRespInt(scoreboardEmptyEvent(nil, gui.AsWindowEvent(0, 0, 0))))
	case 19:
		scoreboardHighlight()
		return 0
	case 20:
		return uint32(scoreboardTeamIndex(uint32(a)))
	case 21:
		return uint32(scoreboardTeamColor(byte(a)))
	case 22:
		return signedByte(scoreboardLocalRank())
	case 23:
		return signedByte(scoreboardTeamRank(uint32(a)))
	case 24:
		return uint32(scoreboardDrawQuest())
	case 25:
		return uint32(bool2int(scoreboardVisible()))
	case 26:
		scoreboardOpen()
		return 0
	case 27:
		return signedByte(scoreboardSetFlag(byte(a), byte(b), uint16(c)))
	case 28:
		return signedByte(scoreboardSetBall(byte(a), uint16(b)))
	case 29:
		return uint32(scoreboardClearFlags())
	case 30:
		return uint32(bool2int(scoreboardEnabled()))
	}
	panic("unknown scoreboard operation")
}
func PortTestScoreboardCallbacks() []unsafe.Pointer { return []unsafe.Pointer{nil, nil, nil} }
