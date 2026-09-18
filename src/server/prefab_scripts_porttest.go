//go:build porttest

package server

import "unsafe"

// Return the actual registered transfer callback for identity-based contracts.
func PortTestPrefabScriptXfer(name string) unsafe.Pointer {
	p := xferFuncs[name]
	if p == nil {
		panic("missing real transfer callback: " + name)
	}
	return p
}
