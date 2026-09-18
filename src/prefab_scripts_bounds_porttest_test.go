//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestPrefabScriptsBoundsOrdering(t *testing.T) {
	buf, free := alloc.Make([]uint32{}, 12)
	defer free()
	p := unsafe.Pointer(&buf[2])
	type row struct {
		Name   string
		Bounds []uint32
		Same   bool
	}
	var rows []row
	bases := [][8]uint32{{10, 0, 0, 10, 20, 10, 10, 20}, {1, 1, 1, 1, 1, 1, 1, 1}, {0, 0xffffffff, 0x80000000, 127, 0x7fffffff, 0x80000000, 0xffffffff, 0}, {0, 0, 1, 0, 0, 1, 1, 1}}
	for bi, base := range bases {
		for a := 0; a < 4; a++ {
			for b := 0; b < 4; b++ {
				if b == a {
					continue
				}
				for c := 0; c < 4; c++ {
					if c == a || c == b {
						continue
					}
					d := 6 - a - b - c
					order := [4]int{a, b, c, d}
					for i := range buf {
						buf[i] = 0xa5a5a5a5
					}
					for i, j := range order {
						copy(buf[2+2*i:2+2*i+2], base[2*j:2*j+2])
					}
					ret := legacy.PortTestPrefabScriptsCall(19, p, nil, nil, 0)
					same := ret == uint64(uintptr(p))
					if !same || buf[0] != 0xa5a5a5a5 || buf[1] != 0xa5a5a5a5 || buf[10] != 0xa5a5a5a5 || buf[11] != 0xa5a5a5a5 {
						t.Fatal("bounds return or guard changed")
					}
					if bi == 0 {
						for i, want := range base {
							if buf[2+i] != want {
								t.Fatal("normal extrema ordering", buf[2:10])
							}
						}
					}
					rows = append(rows, row{fmt.Sprintf("base%d/%v", bi, order), append([]uint32(nil), buf[2:10]...), same})
				}
			}
		}
	}
	// Arbitrary full-width words expose the C helper's mixed signed comparisons.
	state := uint32(0x13579bdf)
	for n := 0; n < 512; n++ {
		for i := 2; i < 10; i++ {
			state = state*1664525 + 1013904223
			buf[i] = state
		}
		ret := legacy.PortTestPrefabScriptsCall(19, p, nil, nil, 0)
		if ret != uint64(uintptr(p)) || buf[0] != 0xa5a5a5a5 || buf[11] != 0xa5a5a5a5 {
			t.Fatal("generated bounds guard")
		}
		rows = append(rows, row{fmt.Sprintf("generated%d", n), append([]uint32(nil), buf[2:10]...), true})
	}
	spellbookCapture(t, "prefab-scripts-bounds", rows, "79fd4a29af468edf4876dd3ab1ef00f1fdb8548e20fac8e6cf0125671e4a5249")
}
