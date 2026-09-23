//go:build porttest

package server

import (
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"unsafe"
)

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

// Baseline for the new two-argument API: the original raw call used by live
// initialization owners. Switch this bridge to CallInitWithArg during conversion;
// identical argument cases and expectations then qualify the new API itself.
func PortTestLifecycleInitWithArg(u *Object, arg unsafe.Pointer) {
	ccall.CallVoidPtr2(u.Init, u.CObj(), arg)
}
