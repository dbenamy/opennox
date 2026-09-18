//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestPrefabRuntimePaths(t *testing.T) {
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	first, freeFirst := alloc.Make([]byte{}, 2048)
	defer freeFirst()
	second, freeSecond := alloc.Make([]byte{}, 2048)
	defer freeSecond()
	*words["path"] = uint32(uintptr(unsafe.Pointer(&first[0])))
	*words["alternate"] = uint32(uintptr(unsafe.Pointer(&second[0])))
	type row struct {
		Name          string
		Return        uint64
		First, Second []byte
	}
	var rows []row
	for _, length := range []int{0, 1, 63, 2046, 2047, 2048} {
		for _, op := range []int{5, 6} {
			for _, nilArg := range []bool{false, true} {
				for i := range first {
					first[i] = 0xa5
					second[i] = 0x5a
				}
				name := strings.Repeat("n", length)
				p, free := alloc.CString(name)
				arg := uint32(uintptr(unsafe.Pointer(p)))
				if nilArg {
					arg = 0
				}
				ret := legacy.PortTestPrefabCall(op, [6]uint32{arg})
				free()
				wantFirst := bytes.Repeat([]byte{0xa5}, 2048)
				wantSecond := bytes.Repeat([]byte{0x5a}, 2048)
				var out []byte
				if op == 5 {
					out = wantFirst
				} else {
					out = wantSecond
				}
				if nilArg {
					out[0] = 0
					if ret != 0 {
						t.Fatal("nil path return")
					}
				} else {
					for i := 0; i < 2047; i++ {
						if i < length {
							out[i] = 'n'
						} else {
							out[i] = 0
						}
					}
					if ret != 1 {
						t.Fatal("path return")
					}
				}
				if !bytes.Equal(first, wantFirst) || !bytes.Equal(second, wantSecond) {
					t.Fatalf("op%d/len%d/nil%t bounded copy or untouched path changed", op, length, nilArg)
				}
				rows = append(rows, row{fmt.Sprintf("op%d/len%d/nil%t", op, length, nilArg), ret, append([]byte(nil), first...), append([]byte(nil), second...)})
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-paths", rows, "75c22902d77b3dd4c3a9d4d71e7b159b68c12cd6510126343fd6c4e62d43fef5")
}
