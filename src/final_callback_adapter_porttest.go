//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"unsafe"
)

// Keep the contract's inputs stable when the single production consumer becomes
// a typed Go callback. The observer itself remains a test-only foreign boundary.
func portTestFinalMapSection(arg unsafe.Pointer, name string, callback unsafe.Pointer) error {
	var fn func(unsafe.Pointer) int
	if callback != nil {
		fn = func(p unsafe.Pointer) int { return ccall.CallIntPtr(callback, p) }
	}
	return nox_xxx_mapReadSectionSpecial_426F40(arg, name, fn)
}
