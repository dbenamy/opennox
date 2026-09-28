//go:build porttest

package server

import "unsafe"

// Bind observer identities without adding entries to the named parser tables.
func PortTestRegisterDamageCallback(key unsafe.Pointer, fn DamageValueFunc) {
	objDamageValue.Register(key, fn)
	objDamage.Register(key, func(obj, source, weapon *Object, amount, kind int32) bool {
		return fn(obj, source, weapon, amount, kind) != 0
	})
}
func PortTestRegisterUseCallback(key unsafe.Pointer, fn UseResultFunc) {
	objectUseNativeFuncs[key] = fn
	objUse.Register(key, func(obj, other *Object) bool { return fn(obj, other) != 0 })
}
func PortTestRegisterCollideCallback(key unsafe.Pointer, fn CollideResultFunc) {
	objectCollideGoFuncs[key] = objectCollideGoEntry{call: func(obj *Object, a, b uintptr) { fn(obj, a, b) }, result: fn}
}
