//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabScriptsBuiltinDispatch(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	globals, restore := legacy.PortTestPrefabScriptsGlobals()
	defer restore()
	*globals[2] = 9
	oldMore, oldEven := legacy.Nox_script_shouldReadMoreXxx, legacy.Nox_script_shouldReadEvenMoreXxx
	defer func() {
		legacy.Nox_script_shouldReadMoreXxx = oldMore
		legacy.Nox_script_shouldReadEvenMoreXxx = oldEven
	}()
	for _, relocate := range []uint32{0, 1} {
		for _, more := range []bool{false, true} {
			for _, even := range []bool{false, true} {
				var calls []string
				legacy.Nox_script_shouldReadMoreXxx = func(id asm.Builtin) bool { calls = append(calls, fmt.Sprintf("first:%d", id)); return more }
				legacy.Nox_script_shouldReadEvenMoreXxx = func(id asm.Builtin) bool { calls = append(calls, fmt.Sprintf("second:%d", id)); return even }
				input := prefabScriptsWords(4, 10, 4, 20, 69, 211, 72)
				raw, f := prefabScriptsFiles(t, input, nil)
				if ret := legacy.PortTestPrefabScriptsCall(4, raw[0], raw[1], nil, relocate); ret != 1 {
					t.Fatal("dispatch copy return", ret)
				}
				a, b := uint32(10), uint32(20)
				var wantCalls []string
				if relocate != 0 {
					wantCalls = append(wantCalls, "first:211")
					if more {
						wantCalls = append(wantCalls, "second:211")
						b += 7
						if even {
							a += 7
						}
					}
				}
				if !reflect.DeepEqual(calls, wantCalls) {
					t.Fatalf("predicate calls %v want%v", calls, wantCalls)
				}
				if !bytes.Equal(prefabScriptsOutput(t, f[1]), prefabScriptsWords(4, a, 4, b, 69, 211, 72)) {
					t.Fatal("predicate-owned remapping differs")
				}
			}
		}
	}
}
