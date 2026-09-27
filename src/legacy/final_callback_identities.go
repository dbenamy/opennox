package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

const (
	playerSectionGUI = iota
	playerSectionMetadata
	playerSectionMusic
)

var playerSectionKeys [3]byte

func playerSectionKey(id int) unsafe.Pointer { return unsafe.Pointer(&playerSectionKeys[id]) }

// The client section table remains mutable. Unknown callbacks keep its original
// foreign dispatch, including the nil argument and signed return value.
func callPlayerFileSection(key unsafe.Pointer) int {
	if ret, ok := callPlayerFileSectionGo(key); ok {
		return ret
	}
	return ccall.CallIntPtr(key, nil)
}

// Save metadata reads the same mutable table through a uintptr callback ABI.
func callPlayerFileMetadata(key unsafe.Pointer) int {
	if ret, ok := callPlayerFileSectionGo(key); ok {
		return ret
	}
	return ccall.CallIntUPtr(key, 0)
}

func callPlayerFileSectionGo(key unsafe.Pointer) (int, bool) {
	switch key {
	case playerSectionKey(playerSectionGUI):
		return playerFileGUI(), true
	case playerSectionKey(playerSectionMetadata):
		return playerFileMetadata(), true
	case playerSectionKey(playerSectionMusic):
		return playerFileMusic(), true
	default:
		return 0, false
	}
}

const (
	tooltipConversation = iota
	tooltipTeamAssign
	tooltipTeamDamage
)

var finalTooltipKeys [3]byte

func finalTooltipKey(id int) unsafe.Pointer { return unsafe.Pointer(&finalTooltipKeys[id]) }

var screenParticleKey, flameCleanseKey byte

func screenParticleCallbackKey() unsafe.Pointer { return unsafe.Pointer(&screenParticleKey) }
func flameCleanseCallbackKey() unsafe.Pointer   { return unsafe.Pointer(&flameCleanseKey) }

func init() {
	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipConversation), func(*gui.Window, *gui.WindowData, uintptr) {})
	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipTeamAssign), func(_ *gui.Window, data *gui.WindowData, _ uintptr) {
		serverOptionsTooltip(false, *(*byte)(unsafe.Pointer(data)))
	})
	gui.RegisterTooltipCallbackGo(finalTooltipKey(tooltipTeamDamage), func(_ *gui.Window, data *gui.WindowData, _ uintptr) {
		serverOptionsTooltip(true, *(*byte)(unsafe.Pointer(data)))
	})
	// FlameCleanse is assigned dynamically, not a named object-type update.
	server.RegisterObjectUpdateCallbackGo(flameCleanseCallbackKey(), temporaryFlameCleanse)
}

func callScreenParticleDraw(key unsafe.Pointer, vp *noxrender.Viewport, p *Nox_screenParticle) int {
	if key == screenParticleCallbackKey() {
		return screenParticleDraw(vp, p)
	}
	return ccall.CallIntPtr2(key, vp.C(), unsafe.Pointer(p))
}
