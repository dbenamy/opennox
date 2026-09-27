//go:build porttest

package legacy

/*
void nox_porttest_animation_complete_in(void);
*/
import "C"

import "unsafe"

var portTestAnimationCompleteIn func()

//export nox_porttest_animation_complete_in
func nox_porttest_animation_complete_in() {
	if portTestAnimationCompleteIn != nil {
		portTestAnimationCompleteIn()
	}
}

func PortTestAnimationCompleteIn(fn func()) (unsafe.Pointer, func()) {
	old := portTestAnimationCompleteIn
	portTestAnimationCompleteIn = fn
	return C.nox_porttest_animation_complete_in, func() { portTestAnimationCompleteIn = old }
}
