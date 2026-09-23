//go:build porttest

package legacy

import "unsafe"

// PortTestBookAbilityCursorWords exposes the actual quickbar cursor state.
func PortTestBookAbilityCursorWords() ([2]*uint32, func()) {
	p := [2]*uint32{(*uint32)(unsafe.Pointer(&dword_5d4594_1047932)), (*uint32)(unsafe.Pointer(&dword_5d4594_1047936))}
	old := [2]uint32{*p[0], *p[1]}
	return p, func() { *p[0], *p[1] = old[0], old[1] }
}
