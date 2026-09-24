//go:build porttest

package server

import (
	"github.com/opennox/libs/types"
	"runtime"
	"unsafe"
)

// Identical original raw-pointer cases now exercise the production pointer API.
func PortTestCollisionWith(u, target *Object, normal *types.Pointf) {
	u.CallCollideWith(target, normal)
}

// PortTestCollisionResult exercises a named native result owner without changing
// the actor's stored callback. Direct-helper fixtures must not alter owner state.
func PortTestCollisionResult(key unsafe.Pointer, obj, target *Object, normal *types.Pointf) uint32 {
	// Unlike the old C call, conversion to uintptr alone does not force Go
	// stack arguments to escape. Pin pointer inputs across arbitrary owner calls.
	var pins runtime.Pinner
	if obj != nil {
		pins.Pin(obj)
	}
	if target != nil {
		pins.Pin(target)
	}
	if normal != nil {
		pins.Pin(normal)
	}
	defer pins.Unpin()
	fnc := objectCollideGoFuncs[key].result
	if fnc == nil {
		panic("missing native collision result owner")
	}
	result := fnc(obj, uintptr(unsafe.Pointer(target)), uintptr(unsafe.Pointer(normal)))
	runtime.KeepAlive(obj)
	runtime.KeepAlive(target)
	runtime.KeepAlive(normal)
	return result
}

// PortTestRegisterNativeCollision confines synthetic registrations to one test.
func PortTestRegisterNativeCollision(name string, key unsafe.Pointer, fnc CollideResultFunc, size uintptr) func() {
	old, present := objectCollideGoFuncs[key]
	RegisterObjectCollideNative(name, key, fnc, size)
	return func() {
		delete(collideFuncs, name)
		if present {
			objectCollideGoFuncs[key] = old
		} else {
			delete(objectCollideGoFuncs, key)
		}
	}
}

// PortTestWorldChestDeath preserves the chest fixtures' deliberate reuse of the
// Pentagram owner as a one-argument death callback. Its native identity is data,
// so install the matching test-only typed death route before storing that key.
func PortTestWorldChestDeath() unsafe.Pointer {
	key := collideFuncs["PentagramCollide"].Func
	fnc := objectCollideGoFuncs[key].result
	if fnc == nil {
		panic("missing native Pentagram owner")
	}
	objDeath.Register(key, func(obj *Object) { fnc(obj, 0, 0) })
	return key
}
