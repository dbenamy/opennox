//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"unsafe"
)

// Preserve the map contract's guards and arguments through its native observer.
func portTestFinalMapSection(arg unsafe.Pointer, name string, callback unsafe.Pointer) error {
	var fn func(unsafe.Pointer) int
	if callback != nil {
		fn = func(p unsafe.Pointer) int { return legacy.PortTestFinalPointerCall(callback, p) }
	}
	return nox_xxx_mapReadSectionSpecial_426F40(arg, name, fn)
}
