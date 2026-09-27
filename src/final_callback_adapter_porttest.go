//go:build porttest

package opennox

import "unsafe"

// Keep the contract's inputs stable when the single production consumer becomes
// a typed Go callback. The observer itself remains a test-only foreign boundary.
func portTestFinalMapSection(arg unsafe.Pointer, name string, callback unsafe.Pointer) error {
	return nox_xxx_mapReadSectionSpecial_426F40(arg, name, callback)
}
