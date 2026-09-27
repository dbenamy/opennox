//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

var portTestTooltipKey byte
var portTestTooltipWords [4]uint32

func init() {
	gui.RegisterTooltipCallbackGo(unsafe.Pointer(&portTestTooltipKey), func(window *gui.Window, draw *gui.WindowData, argument uintptr) {
		portTestTooltipWords[0]++
		portTestTooltipWords[1] = uint32(uintptr(window.C()))
		portTestTooltipWords[2] = uint32(uintptr(draw.C()))
		portTestTooltipWords[3] = uint32(argument)
	})
}

// PortTestGUITooltipProbe captures the callback words independently of the engine
// tooltip implementation. Cleanup restores the previous observer state.
func PortTestGUITooltipProbe() (unsafe.Pointer, *[4]uint32, func()) {
	saved := portTestTooltipWords
	portTestTooltipWords = [4]uint32{}
	return unsafe.Pointer(&portTestTooltipKey), &portTestTooltipWords, func() { portTestTooltipWords = saved }
}
