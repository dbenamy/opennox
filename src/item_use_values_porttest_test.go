//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

var itemUseContractKey uint32

func TestItemUseNativeResults(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	it, free2 := alloc.New(server.Object{})
	defer free2()
	var got [2]*server.Object
	var calls int
	var value int32
	key := unsafe.Pointer(&itemUseContractKey)
	restore := server.PortTestItemUseRegistration(key, func(a, b *server.Object) int32 { got = [2]*server.Object{a, b}; calls++; return value })
	defer restore()
	cb := server.UseFuncPtr{Ptr: key}
	collisionRegistryGrow(128)
	for _, v := range []int32{-2147483648, -1, 0, 1, 256, 2147483647} {
		value = v
		for _, a := range []*server.Object{nil, u} {
			for _, b := range []*server.Object{nil, it} {
				for mode := 0; mode < 3; mode++ {
					calls = 0
					got = [2]*server.Object{}
					switch mode {
					case 0:
						if x := cb.CallResult(a, b); x != v {
							t.Fatalf("native result: %d != %d", x, v)
						}
					case 1:
						if x := cb.Get()(a, b); x != (v != 0) {
							t.Fatal("native bool", v, x)
						}
					case 2:
						cb.CallDiscard(a, b)
					}
					if calls != 1 || got != ([2]*server.Object{a, b}) {
						t.Fatal("native arguments/count", mode, v, calls, got)
					}
				}
			}
		}
	}
}

func TestItemUseRawResults(t *testing.T) {
	u, free := alloc.New(server.Object{})
	defer free()
	it, free2 := alloc.New(server.Object{})
	defer free2()
	key, words, result, restore := legacy.PortTestXferSoundRawObserver()
	defer restore()
	cb := server.UseFuncPtr{Ptr: key}
	for _, v := range []int32{-2147483648, -1, 0, 1, 256, 2147483647} {
		*result = v
		for _, a := range []*server.Object{nil, u} {
			for _, b := range []*server.Object{nil, it} {
				for mode := 0; mode < 3; mode++ {
					*words = [3]uint32{}
					switch mode {
					case 0:
						if x := cb.CallResult(a, b); x != v {
							t.Fatalf("raw result: %d != %d", x, v)
						}
					case 1:
						if x := cb.Get()(a, b); x != (v != 0) {
							t.Fatal("raw bool", v, x)
						}
					case 2:
						cb.CallDiscard(a, b)
					}
					if *words != ([3]uint32{1, uint32(uintptr(unsafe.Pointer(a))), uint32(uintptr(unsafe.Pointer(b)))}) {
						t.Fatal("raw arguments/count", mode, v, *words)
					}
				}
			}
		}
	}
}
