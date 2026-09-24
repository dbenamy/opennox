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

func TestModifierFireHooksRemainLate(t *testing.T) {
	mod, freeMod := alloc.New(server.ModifierEff{})
	defer freeMod()
	objects, freeObjects := alloc.Make([]server.Object{}, 3)
	defer freeObjects()
	data, freeData := alloc.New(uint32(0))
	defer freeData()
	oldRing, oldBlue := legacy.Nox_xxx_fireRingEffect_4E05B0, legacy.Nox_xxx_blueFREffect_4E05F0
	defer func() { legacy.Nox_xxx_fireRingEffect_4E05B0 = oldRing; legacy.Nox_xxx_blueFREffect_4E05F0 = oldBlue }()
	ring := server.PortTestModifierCallback("damage", "FireRingEffect")
	blue := server.PortTestModifierCallback("damage", "BlueFireRingEffect")
	if ring == nil || blue == nil || ring == blue {
		t.Fatal("fire hook keys", ring, blue)
	}
	for generation := 1; generation <= 3; generation++ {
		var captured [2][4]unsafe.Pointer
		var calls [2]int
		var versions [2]int
		makeHook := func(index int) func(unsafe.Pointer, *server.Object, *server.Object, *server.Object) {
			return func(m unsafe.Pointer, a, b, c *server.Object) {
				captured[index] = [4]unsafe.Pointer{m, a.CObj(), b.CObj(), c.CObj()}
				calls[index]++
				versions[index] = generation
			}
		}
		legacy.Nox_xxx_fireRingEffect_4E05B0 = makeHook(0)
		legacy.Nox_xxx_blueFREffect_4E05F0 = makeHook(1)
		for mask := 0; mask < 16; mask++ {
			args := [4]unsafe.Pointer{mod.C(), objects[0].CObj(), objects[1].CObj(), objects[2].CObj()}
			for i := range args {
				if mask&(1<<i) == 0 {
					args[i] = nil
				}
			}
			captured = [2][4]unsafe.Pointer{}
			calls = [2]int{}
			versions = [2]int{}
			*data = 0x89abcdef
			runtime.GC()
			for i, key := range []unsafe.Pointer{ring, blue} {
				legacy.PortTestModifierCall5(key, (*server.ModifierEff)(args[0]), (*server.Object)(args[1]), (*server.Object)(args[2]), (*server.Object)(args[3]), unsafe.Pointer(data))
				if captured[i] != args || calls[i] != 1 || versions[i] != generation {
					t.Fatal("late fire hook", generation, mask, i, captured[i], calls[i], versions[i])
				}
			}
			if *data != 0x89abcdef {
				t.Fatal("ignored fire data changed")
			}
		}
	}
}
