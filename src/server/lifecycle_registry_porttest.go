//go:build porttest

package server

import "unsafe"

func PortTestLifecycleRegistry(name string, create bool) (unsafe.Pointer, uintptr) {
	if create {
		p, ok := createFuncs[name]
		if !ok {
			panic("missing create registration: " + name)
		}
		return p, 0
	}
	d, ok := initFuncs[name]
	if !ok {
		panic("missing init registration: " + name)
	}
	return d.Func, d.DataSize
}

// The original baseline invoked the raw two-pointer call here. The same cases
// now qualify the public typed API and its raw callback fallback.
func PortTestLifecycleInitWithArg(u *Object, arg unsafe.Pointer) {
	u.CallInitWithArg(arg)
}
