//go:build porttest

package legacy

import "unsafe"

func PortTestWorldWallAttach(drawable bool, value unsafe.Pointer, x, y int) unsafe.Pointer {
	if drawable {
		return worldDrawableAttach(asDrawable((*nox_drawable)(value)), x, y)
	}
	return worldDoorAttach(value, x, y).C()
}
