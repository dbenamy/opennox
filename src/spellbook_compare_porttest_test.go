//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSpellbookNameComparison(t *testing.T) {
	o := newSpellbookOwner(t)
	o.resetBook(t)
	o.bookCall("nox_xxx_bookInit_45B9D0")
	*o.words["dword_5d4594_1046868"] = 1
	names := [][]uint16{nil, {}, {'A'}, {'a'}, {'a', 'B'}, {'A', 'b'}, {'a', 0, 'z'}, {0x00c9}, {0x00e9}, {0x03a9}, {0x03c9}, {0xd800, 'x'}, {0xd801, 'x'}, {0xffff}}
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	old := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, old) })
	for i, name := range names {
		guides[7*(i+1)+1] = 1
		guides[7*(i+1)] = 0
		if name != nil {
			p := tooltipWide(t, name)
			guides[7*(i+1)] = uint32(uintptr(unsafe.Pointer(&p[0])))
		}
	}
	ids, free := alloc.New([2]uint32{})
	defer free()
	var rows []spellbookResult
	lower := func(v uint16) uint16 {
		if v >= 'A' && v <= 'Z' {
			return v + 32
		}
		return v
	}
	for i, a := range names {
		for j, b := range names {
			*ids = [2]uint32{uint32(i + 1), uint32(j + 1)}
			ret := o.bookCall("nox_xxx_guiSpellSortFn_45ABC0", uint32(uintptr(unsafe.Pointer(&ids[0]))), uint32(uintptr(unsafe.Pointer(&ids[1]))))
			if a == nil || b == nil {
				if ret != 0 {
					t.Fatal("missing names compare equal")
				}
			} else if i <= 6 && j <= 6 {
				want := int32(0)
				for k := 0; ; k++ {
					var x, y uint16
					if k < len(a) {
						x = lower(a[k])
					}
					if k < len(b) {
						y = lower(b[k])
					}
					if x != y || x == 0 {
						want = int32(x) - int32(y)
						break
					}
				}
				if int32(ret) != want {
					t.Fatalf("ASCII comparison names%d/%d got%d want%d", i, j, int32(ret), want)
				}
			}
			rows = append(rows, o.bookSnapshot(fmt.Sprintf("name%d-name%d", i, j), ret))
		}
	}
	spellbookCapture(t, "name-compare", rows, "d377183956fc746b21e47de689c0c8aab363edfd87bd1c78bd03e2f136c91693")
}
