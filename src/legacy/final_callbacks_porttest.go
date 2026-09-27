//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
)

var portTestFinalPtr func(unsafe.Pointer) int32
var portTestFinalPtr2 func(unsafe.Pointer, unsafe.Pointer) int32

var portTestFinalKeys [2]byte

func init() {
	playerFileCallbacks[unsafe.Pointer(&portTestFinalKeys[0])] = func(p unsafe.Pointer) int { return int(portTestFinalPtr(p)) }
	screenParticleCallbacks[unsafe.Pointer(&portTestFinalKeys[1])] = func(vp *noxrender.Viewport, p *Nox_screenParticle) int {
		return int(portTestFinalPtr2(vp.C(), unsafe.Pointer(p)))
	}
}

// These observers preserve the callback arguments and signed results.
func PortTestObserveFinalPtr(fn func(unsafe.Pointer) int32) (unsafe.Pointer, func()) {
	old := portTestFinalPtr
	portTestFinalPtr = fn
	return unsafe.Pointer(&portTestFinalKeys[0]), func() { portTestFinalPtr = old }
}

func PortTestObserveFinalPtr2(fn func(unsafe.Pointer, unsafe.Pointer) int32) (unsafe.Pointer, func()) {
	old := portTestFinalPtr2
	portTestFinalPtr2 = fn
	return unsafe.Pointer(&portTestFinalKeys[1]), func() { portTestFinalPtr2 = old }
}

func PortTestFinalMenuCallback() gui.WindowFunc             { return MainMenuEvent }
func PortTestFinalConversationTooltip() unsafe.Pointer      { return finalTooltipKey(tooltipConversation) }
func PortTestFinalFlameCallback() unsafe.Pointer            { return flameCleanseCallbackKey() }
func PortTestFinalPlayerSectionCall(key unsafe.Pointer) int { return callPlayerFileSection(key) }

func PortTestFinalPointerCall(key, arg unsafe.Pointer) int { return callPlayerFileCallback(key, arg) }
