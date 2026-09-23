//go:build porttest

package server

import "unsafe"

func PortTestDamageRegistry(name string) unsafe.Pointer {
	p, ok := damageFuncs[name]
	if !ok || p == nil {
		panic("unknown damage registry entry: " + name)
	}
	return p
}
