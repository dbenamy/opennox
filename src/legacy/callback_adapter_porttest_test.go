//go:build porttest

package legacy_test

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestLegacyCallbackAdapters(t *testing.T) {
	vp, freeVP := alloc.New(noxrender.Viewport{})
	defer freeVP()
	dr, freeDrawable := alloc.New(client.Drawable{})
	defer freeDrawable()

	for _, which := range []int{1, 2} {
		dr.DrawFuncPtr = legacy.PortTestAdapterDrawCallback(which)
		for _, vpArg := range []*noxrender.Viewport{vp, nil} {
			for _, result := range []int{-1 << 31, -1, 0, 1, 1<<31 - 1} {
				legacy.PortTestAdapterDrawReset(result)
				got := legacy.CallDrawFunc(dr, vpArg)
				if got != result {
					t.Fatalf("draw callback %d return: got %d want %d", which, got, result)
				}
				if got := legacy.PortTestAdapterDrawValue(0); got != uintptr(unsafe.Pointer(vpArg)) {
					t.Fatalf("draw callback %d viewport: got %#x want %#x", which, got, uintptr(unsafe.Pointer(vpArg)))
				}
				if got := legacy.PortTestAdapterDrawValue(1); got != uintptr(unsafe.Pointer(dr)) {
					t.Fatalf("draw callback %d drawable: got %#x want %#x", which, got, uintptr(unsafe.Pointer(dr)))
				}
				if got := legacy.PortTestAdapterDrawValue(2); got != uintptr(which) {
					t.Fatalf("draw callback selection: got %d want %d", got, which)
				}
				if got := legacy.PortTestAdapterDrawValue(3); got != 1 {
					t.Fatalf("draw callback count: got %d want 1", got)
				}
			}
		}

	}

	obj, freeObject := alloc.New(server.Object{})
	defer freeObject()
	legacy.PortTestAdapterObjectReset()
	for i, objArg := range []*server.Object{obj, nil, obj, nil} {
		legacy.Nox_call_objectType_new_go(legacy.PortTestAdapterObjectCallback(), objArg)
		if got := legacy.PortTestAdapterObjectValue(0); got != uintptr(unsafe.Pointer(objArg)) {
			t.Fatalf("object callback pointer: got %#x want %#x", got, uintptr(unsafe.Pointer(objArg)))
		}
		if got := legacy.PortTestAdapterObjectValue(1); got != uintptr(i+1) {
			t.Fatalf("object callback count: got %d want %d", got, i+1)
		}
	}
}
