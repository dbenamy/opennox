//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"math"
	"strings"
	"testing"
	"unsafe"
)

func TestPrefabRuntimeMetadata(t *testing.T) {
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	data, free := alloc.Make([]byte{}, 4*76)
	defer free()
	base := uint32(uintptr(unsafe.Pointer(&data[0])))
	*words["metadata"] = base
	type row struct {
		Case   string
		Return uint64
	}
	var rows []row
	for _, width := range []uint32{0, 0x80000000, 1, 0x3f800000, 0x43a2a1bc, 0xbf800000, 0x7f800000, 0xff800000, 0x7fc01234} {
		names := []string{"Alpha", "ALPHA", strings.Repeat("n", 63)}
		clear(data)
		for i, name := range names {
			copy(data[i*76:], name)
			binary.LittleEndian.PutUint32(data[i*76+64:], width)
			binary.LittleEndian.PutUint32(data[i*76+68:], width^0x80000000)
			binary.LittleEndian.PutUint32(data[i*76+72:], uint32(100+100*i))
		}
		for _, count := range []int32{-1, 0, 1, 3} {
			*words["count"] = uint32(count)
			got := legacy.PortTestPrefabCall(4, [6]uint32{})
			if got != uint64(uint32(count)) {
				t.Fatal("count changed")
			}
			rows = append(rows, row{fmt.Sprintf("count/%08x/%d", width, count), got})
			for _, index := range []int32{math.MinInt32, -1, 0, 1, 2, 3, 4, math.MaxInt32} {
				ptr := uint32(legacy.PortTestPrefabCall(3, [6]uint32{uint32(index)}))
				norm := uint64(0)
				if index >= 0 && index <= count {
					want := base + 76*uint32(index)
					if ptr != want {
						t.Fatal("metadata getter boundary")
					}
					norm = uint64(index + 1)
				} else if ptr != 0 {
					t.Fatal("invalid metadata pointer")
				}
				rows = append(rows, row{fmt.Sprintf("pointer/%08x/%d/%d", width, count, index), norm})
				for _, op := range []int{12, 13} {
					got := legacy.PortTestPrefabCall(op, [6]uint32{uint32(index)})
					want := -1.0
					if index >= 0 && index < count {
						bits := width
						if op == 13 {
							bits ^= 0x80000000
						}
						want = float64(math.Float32frombits(bits))
					}
					if math.IsNaN(want) {
						if !math.IsNaN(math.Float64frombits(got)) {
							t.Fatal("NaN changed")
						}
					} else if got != math.Float64bits(want) {
						t.Fatal("dimension value/boundary changed")
					}
					rows = append(rows, row{fmt.Sprintf("dimension/%08x/%d/%d/%d", width, count, index, op), got})
				}
			}
			for _, name := range []string{"", "alpha", "AlPhA", "missing", names[2], strings.ToUpper(names[2])} {
				text, free := alloc.CString(name)
				got := legacy.PortTestPrefabCall(2, [6]uint32{uint32(uintptr(unsafe.Pointer(text)))})
				free()
				want := uint64(math.MaxUint32)
				for i := int32(0); i < count; i++ {
					if strings.EqualFold(name, names[i]) {
						want = uint64(i)
						break
					}
				}
				if got != want {
					t.Fatal("case-insensitive first-match lookup changed")
				}
				rows = append(rows, row{fmt.Sprintf("name/%08x/%d/%s", width, count, name), got})
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-metadata", rows, "9342a1ef9d41f510d01a9d3714ade6d0ff6395746379d6fa499bbecff125a133")
}
