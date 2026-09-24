//go:build porttest

package opennox

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestModifierNativeCallbackForwarding(t *testing.T) {
	mod, freeMod := alloc.New(server.ModifierEff{})
	defer freeMod()
	objects, freeObjects := alloc.Make([]server.Object{}, 4)
	defer freeObjects()
	buffer, freeBuffer := alloc.New([3]uint32{})
	defer freeBuffer()
	mod.Price20 = -1234567
	mod.AllowWeapons28 = 0xfedcba98
	before := *mod
	slots := new([3]byte)
	keys := [3]unsafe.Pointer{unsafe.Pointer(&slots[0]), unsafe.Pointer(&slots[1]), unsafe.Pointer(&slots[2])}
	var value int32
	var received [6]unsafe.Pointer
	defer func() { received = [6]unsafe.Pointer{} }()
	calls, kind := 0, 0
	record := func(which int, args [6]unsafe.Pointer) {
		runtime.GC()
		calls++
		kind = which
		received = args
	}
	server.RegisterModifierEffect3(keys[0], func(m *server.ModifierEff, a, b *server.Object) int32 {
		record(3, [6]unsafe.Pointer{m.C(), a.CObj(), b.CObj()})
		return value
	})
	server.RegisterModifierEffect5(keys[1], func(m *server.ModifierEff, a, b, c *server.Object, out unsafe.Pointer) {
		record(5, [6]unsafe.Pointer{m.C(), a.CObj(), b.CObj(), c.CObj(), out})
		if out != nil {
			*(*uint32)(out) = ^uint32(value)
		}
	})
	server.RegisterModifierEffect6(keys[2], func(m *server.ModifierEff, a, b, c, d *server.Object, out unsafe.Pointer) {
		record(6, [6]unsafe.Pointer{m.C(), a.CObj(), b.CObj(), c.CObj(), d.CObj(), out})
		if out != nil {
			*(*uint32)(out) = ^uint32(value)
		}
	})
	for _, mode := range []int{0, 1, 2, 3} {
		arity := 3
		if mode == 2 {
			arity = 5
		}
		if mode == 3 {
			arity = 6
		}
		for mask := 0; mask < 1<<arity; mask++ {
			for _, word := range []int32{0, 1, 255, 256, 0x7fffffff, -0x80000000, -1} {
				value = word
				received = [6]unsafe.Pointer{}
				calls = 0
				kind = 0
				*buffer = [3]uint32{0xaabbccdd, 0x12345678, 0x98765432}
				args := [6]unsafe.Pointer{mod.C(), objects[0].CObj(), objects[1].CObj(), objects[2].CObj(), objects[3].CObj(), unsafe.Pointer(&buffer[1])}
				if arity == 5 {
					args[4] = unsafe.Pointer(&buffer[1])
				}
				for i := range args {
					if i >= arity || mask&(1<<i) == 0 {
						args[i] = nil
					}
				}
				m := (*server.ModifierEff)(args[0])
				a := (*server.Object)(args[1])
				b := (*server.Object)(args[2])
				switch mode {
				case 0:
					if got := server.CallModifierEffect3Result(keys[0], m, a, b); got != word {
						t.Fatal("native result", got, word)
					}
				case 1:
					server.CallModifierEffect3Discard(keys[0], m, a, b)
				case 2:
					server.CallModifierEffect5(keys[1], m, a, b, (*server.Object)(args[3]), args[4])
				case 3:
					server.CallModifierEffect6(keys[2], m, a, b, (*server.Object)(args[3]), (*server.Object)(args[4]), args[5])
				}
				if received != args || calls != 1 || kind != arity {
					t.Fatalf("mode%d mask%x: args%v count%d kind%d want %v/1/%d", mode, mask, received, calls, kind, args, arity)
				}
				expected := [3]uint32{0xaabbccdd, 0x12345678, 0x98765432}
				if arity > 3 && args[arity-1] != nil {
					expected[1] = ^uint32(word)
				}
				if *buffer != expected || *mod != before {
					t.Fatal("native callback mutation", mode, mask, *buffer, expected)
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
