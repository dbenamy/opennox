//go:build porttest

package server

import "unsafe"

func PortTestDeathRegistry(name string) (unsafe.Pointer, uintptr) {
	entry, ok := deathFuncs[name]
	if !ok || entry.Func == nil {
		panic("unknown death registry entry: " + name)
	}
	return entry.Func, entry.DataSize
}
