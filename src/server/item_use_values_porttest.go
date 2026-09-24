//go:build porttest

package server

import (
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"unsafe"
)

// Isolate mutable registries for serial callback contracts and restore them.
func PortTestItemUseRegistration(key unsafe.Pointer, fn UseResultFunc) func() {
	const name = "PortTestItemUseResult"
	if _, ok := useFuncs[name]; ok {
		panic("nested item use registration")
	}
	oldBool, oldNative := objUse, objectUseNativeFuncs
	objUse = ccall.NewFuncs(func(p unsafe.Pointer) UseFunc { return oldBool.Get(p) })
	objectUseNativeFuncs = make(map[unsafe.Pointer]UseResultFunc, len(oldNative)+1)
	for p, f := range oldNative {
		objectUseNativeFuncs[p] = f
	}
	RegisterObjectUseNative(name, key, fn, 0)
	return func() { objUse = oldBool; objectUseNativeFuncs = oldNative; delete(useFuncs, name) }
}
