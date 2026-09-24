//go:build porttest

package server

import "unsafe"

// Read the actual resource callback registration, independently of legacy fixtures.
func PortTestItemIdentity(family, name string) (unsafe.Pointer, uintptr) {
	switch family {
	case "use":
		e, ok := useFuncs[name]
		if !ok {
			panic("unknown use name")
		}
		return e.Func, e.DataSize
	case "drop":
		p, ok := dropFuncs[name]
		if !ok {
			panic("unknown drop name")
		}
		return p, 0
	case "pickup":
		p, ok := pickupFuncs[name]
		if !ok {
			panic("unknown pickup name")
		}
		return p, 0
	default:
		panic("unknown item callback family")
	}
}
