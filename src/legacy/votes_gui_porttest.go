//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"

	"unsafe"
)

func PortTestVoteGUIOwner() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"topic": (*uint32)(unsafe.Pointer(&voteTopic)), "window": (*uint32)(unsafe.Pointer(&voteWindow)),
		"players": (*uint32)(unsafe.Pointer(&votePlayersWindow)), "topics": (*uint32)(unsafe.Pointer(&voteTopicsWindow)),
		"count": (*uint32)(unsafe.Pointer(&voteNameCount)), "previousCount": (*uint32)(unsafe.Pointer(&votePreviousNameCount)),
		"choice": (*uint32)(unsafe.Pointer(&voteChoice)), "previousChoice": (*uint32)(unsafe.Pointer(&votePreviousChoice)),
	}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	savedNames, savedPrevious := voteNames, votePreviousNames
	voteNames, votePreviousNames = [32][28]uint16{}, [32][28]uint16{}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
		voteNames, votePreviousNames = savedNames, savedPrevious
	}
}
func PortTestVoteGUI(op string, topic int) int {
	switch op {
	case "close":
		return voteGUIClose()
	case "reset-choice":
		return voteGUIReset()
	case "init":
		return voteGUIInit()
	case "show":
		voteGUIShow(uint32(topic))
	case "hide":
		return voteGUIHide()
	case "selection":
		return voteGUISelection()
	case "choice":
		return voteGUIResetChoice()
	case "topic":
		voteGUIConfirmTopic()
	default:
		panic(op)
	}
	return 0
}
func PortTestVoteGUIWindow() *gui.Window {
	return voteWindow
}
