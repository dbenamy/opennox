//go:build porttest

package server

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"unsafe"
)

// Original pointer-call baseline; switches to the new API only at conversion.
func PortTestCollisionWith(u, target *Object, normal *types.Pointf) {
	ccall.CallVoidPtr3(u.Collide, u.CObj(), unsafe.Pointer(target), unsafe.Pointer(normal))
}
