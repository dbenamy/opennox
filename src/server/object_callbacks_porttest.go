//go:build porttest

package server

import (
	"github.com/opennox/libs/types"
	"unsafe"
)

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

func PortTestRegisterCreateCallback(key unsafe.Pointer, fn ObjectCreateFunc) {
	objectCreateFuncs.Register(key, fn)
}
func PortTestRegisterInitCallback(key unsafe.Pointer, fn ObjectInitFunc) { objectInitGoFuncs[key] = fn }
func PortTestRegisterInitArgCallback(key unsafe.Pointer, fn func(*Object, unsafe.Pointer)) {
	objectInitArgGoFuncs[key] = fn
}
func PortTestRegisterDeathCallback(key unsafe.Pointer, fn DeathFunc) { objDeath.Register(key, fn) }
func PortTestRegisterDropCallback(key unsafe.Pointer, fn func(*Object, *Object, types.Pointf) bool) {
	objDrop.Register(key, fn)
}
func PortTestRegisterPickupCallback(key unsafe.Pointer, fn PickupFunc) { objPickup.Register(key, fn) }
func PortTestRegisterXferCallback(key unsafe.Pointer, fn XferFunc)     { objectXferGoFuncs[key] = fn }
func PortTestRegisterDamageSoundCallback(key unsafe.Pointer, fn DamageSoundFunc) {
	objectDamageSoundGoFuncs[key] = fn
}
