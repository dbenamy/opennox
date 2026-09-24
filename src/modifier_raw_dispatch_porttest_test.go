//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestModifierRawCallbackForwarding(t *testing.T) {
	mod, freeMod := alloc.New(server.ModifierEff{})
	defer freeMod()
	objects, freeObjects := alloc.Make([]server.Object{}, 4)
	defer freeObjects()
	buffer, freeBuffer := alloc.New([3]uint32{})
	defer freeBuffer()
	mod.Price20 = -1234567
	mod.AllowWeapons28 = 0xfedcba98
	before := *mod
	data := unsafe.Pointer(&buffer[1])
	// result3, discarded-result3, void3, void5 and void6 have distinct contracts.
	for _, mode := range []int{0, 1, 2, 3, 4} {
		arity := 3
		if mode == 3 {
			arity = 5
		}
		if mode == 4 {
			arity = 6
		}
		for mask := 0; mask < 1<<arity; mask++ {
			for _, value := range []int32{0, 1, 255, 256, 0x7fffffff, -0x80000000, -1} {
				*buffer = [3]uint32{0xaabbccdd, 0x12345678, 0x98765432}
				all := [6]unsafe.Pointer{unsafe.Pointer(mod), unsafe.Pointer(&objects[0]), unsafe.Pointer(&objects[1]), unsafe.Pointer(&objects[2]), unsafe.Pointer(&objects[3]), data}
				if arity == 5 {
					all[4] = data
				}
				var want [6]uintptr
				for i := range all {
					if i >= arity || mask&(1<<i) == 0 {
						all[i] = nil
					}
					want[i] = uintptr(all[i])
				}
				keys := legacy.PortTestModifierObserverReset(value)
				seen := map[unsafe.Pointer]bool{}
				for _, key := range keys {
					if key == nil || seen[key] {
						t.Fatal("observer keys", keys)
					}
					seen[key] = true
				}
				m := (*server.ModifierEff)(all[0])
				a := (*server.Object)(all[1])
				b := (*server.Object)(all[2])
				runtime.GC()
				kind := 1
				switch mode {
				case 0:
					if got := legacy.PortTestModifierCall3Result(keys[0], m, a, b); got != value {
						t.Fatal("raw result", got, value)
					}
				case 1:
					legacy.PortTestModifierCall3Discard(keys[0], m, a, b)
				case 2:
					kind = 2
					legacy.PortTestModifierCall3Discard(keys[1], m, a, b)
				case 3:
					kind = 3
					legacy.PortTestModifierCall5(keys[2], m, a, b, (*server.Object)(all[3]), all[4])
				case 4:
					kind = 4
					legacy.PortTestModifierCall6(keys[3], m, a, b, (*server.Object)(all[3]), (*server.Object)(all[4]), all[5])
				}
				got, calls, gotKind := legacy.PortTestModifierObserverSnapshot()
				if got != want || calls != 1 || gotKind != kind {
					t.Fatalf("mode%d mask%x forwarding %v/%d/%d want %v/1/%d", mode, mask, got, calls, gotKind, want, kind)
				}
				expected := [3]uint32{0xaabbccdd, 0x12345678, 0x98765432}
				if arity > 3 && all[arity-1] != nil {
					expected[1] = ^uint32(value)
				}
				if *buffer != expected || *mod != before {
					t.Fatal("raw callback mutation", mode, mask, *buffer, expected)
				}
				for i := range objects {
					if objects[i] != (server.Object{}) {
						t.Fatal("object modified", i)
					}
				}
			}
		}
	}
}
