//go:build porttest

package legacy

/*
#include "GAME3.h"
extern nox_gui_animation* nox_wnd_xxx_1307732;
extern nox_gui_animation* nox_wnd_xxx_1308092;
*/
import "C"
import "unsafe"

func PortTestCharacterAnimationOwner(color bool, anim unsafe.Pointer) func() {
	p := &C.nox_wnd_xxx_1307732
	if color {
		p = &C.nox_wnd_xxx_1308092
	}
	old := *p
	*p = (*C.nox_gui_animation)(anim)
	return func() { *p = old }
}
func PortTestCharacterAnimationDone(color bool) int {
	if color {
		return int(C.sub_4A6C90())
	}
	return int(C.sub_4A49A0())
}
func PortTestCharacterClassStart() int { return int(C.sub_4A4970()) }
