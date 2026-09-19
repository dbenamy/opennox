//go:build porttest

package server

import "unsafe"

// PortTestUnitGameplayRegistration reads the actual production callback registry.
func PortTestUnitGameplayRegistration(name string, use bool) (unsafe.Pointer, uintptr) {
	table := updateFuncs
	if use {
		table = useFuncs
	}
	def, ok := table[name]
	if !ok || def.Func == nil {
		panic("missing unit gameplay registration: " + name)
	}
	return def.Func, def.DataSize
}
