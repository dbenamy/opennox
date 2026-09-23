//go:build porttest

package server

import (
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"unsafe"
)

// Isolate the existing boolean registry without mutating its underlying map.
// Tests are serial because object callback registrations are process globals.
func PortTestDamageBoolRegistration(ptr unsafe.Pointer, fn DamageFunc) func() {
	const name = "PortTestDamageBoolOnly"
	if _, ok := damageFuncs[name]; ok {
		panic("nested boolean registration fixture")
	}
	old := objDamage
	objDamage = ccall.NewFuncs(func(p unsafe.Pointer) DamageFunc { return old.Get(p) })
	RegisterObjectDamageGo(name, ptr, fn)
	return func() { objDamage = old; delete(damageFuncs, name) }
}
