//go:build porttest

package legacy

import "unsafe"

// CString uses this same centralized allocation path in each profile.
func portTestStringMalloc(n uintptr) unsafe.Pointer { return legacyMalloc(n) }
