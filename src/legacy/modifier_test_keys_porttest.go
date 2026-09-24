//go:build porttest

package legacy

import "unsafe"

var modifierTestIdentitySlots [9]byte

func modifierTestKey(id int) unsafe.Pointer { return unsafe.Pointer(&modifierTestIdentitySlots[id]) }
