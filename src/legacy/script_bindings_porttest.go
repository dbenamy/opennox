//go:build porttest

package legacy

import "unsafe"

// Exercise the production callback-transfer adapter with fixture-owned records.
func PortTestScriptBindingCallback(record, name unsafe.Pointer) int {
	return objectXferScript(record, name)
}

func PortTestScriptBindingRegistryOwner() (*uint32, func()) {
	p := &scriptBindingStringCount
	old := *p
	return p, func() { *p = old }
}

func PortTestScriptBindingIntern(text unsafe.Pointer) int {
	return scriptBindingIntern(text)
}
