//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

var portTestAnimationCompleteIn func()

var portTestAnimationCompleteInKey byte

func init() {
	gui.RegisterAnimationCallbackGo(unsafe.Pointer(&portTestAnimationCompleteInKey), func() int {
		if portTestAnimationCompleteIn != nil {
			portTestAnimationCompleteIn()
		}
		return 0
	})
}

func PortTestAnimationCompleteIn(fn func()) (unsafe.Pointer, func()) {
	old := portTestAnimationCompleteIn
	portTestAnimationCompleteIn = fn
	return unsafe.Pointer(&portTestAnimationCompleteInKey), func() { portTestAnimationCompleteIn = old }
}
