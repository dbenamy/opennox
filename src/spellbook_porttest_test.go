//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func TestSpellbookInitialization(t *testing.T) {
	o := newSpellbookOwner(t)
	var rows []spellbookResult
	for missing := -1; missing < len(bookResourceNames); missing++ {
		o.resetBook(t)
		before := len(o.windows)
		if missing >= 0 {
			o.missing = bookResourceNames[missing]
		}
		result := o.bookCall("nox_xxx_bookInit_45B9D0")
		o.collect()
		if missing < 0 {
			if result != 1 || o.bookWindow() == nil || len(o.windows) != before+7 {
				t.Fatal("complete book window initialization", result, len(o.windows), before)
			}
			if *o.words["nox_xxx_aNox_cfg_0_587000_132132"] != 1 {
				t.Fatal("initial list mode")
			}
			if len(o.loads) != 15 {
				t.Fatal("book image loads", o.loads)
			} // arrows are reused for two buttons
		} else {
			if result != 0 || o.bookWindow() != nil || len(o.windows) != before {
				t.Fatal("partial image failure created windows")
			}
			if !reflect.DeepEqual(o.loads, bookResourceNames[:missing+1]) {
				t.Fatal("partial image load order", o.loads)
			}
		}
		rows = append(rows, o.bookSnapshot(o.missing, result))
	}
	spellbookCapture(t, "initialization", rows, "e814a31a3954dd373781f654266f8d2b4d55cb75ec3f89868a16afcf5acb6eeb")
}
func TestSpellbookPathBuffer(t *testing.T) {
	o := newSpellbookOwner(t)
	points, free := alloc.Make([]int32{}, 4)
	defer free()
	var rows []spellbookResult
	for _, count := range []uint32{0, 1, 18, 19, 20, 21} {
		for caseID, coords := range [][4]int32{{0, 0, 0, 0}, {1, -2, 3, -4}, {-2147483648, 2147483647, 16777217, -16777217}} {
			o.resetBook(t)
			copy(points, coords[:])
			for i := range o.region {
				o.region[i] = uint32(i)*0x12345 + 7
			}
			o.region[(1046680-1046612)/4] = count
			want := append([]uint32(nil), o.region...)
			expected := count
			if count < 20 {
				expected++
				offsets := []int{1046676, 1046680, 1046684, 1046688}
				for i, off := range offsets {
					want[(off+8*int(expected)-1046612)/4] = math.Float32bits(float32(coords[i]))
				}
				want[(1046680-1046612)/4] = expected
			}
			got := o.bookCall("sub_45D7D0", uint32(uintptr(unsafe.Pointer(&points[0]))), uint32(uintptr(unsafe.Pointer(&points[2]))))
			if got != expected || !reflect.DeepEqual(o.region, want) {
				t.Fatalf("path count %d case %d", count, caseID)
			}
			if !reflect.DeepEqual(points, coords[:]) {
				t.Fatal("path callback changed input")
			}
			rows = append(rows, o.bookSnapshot(fmt.Sprintf("path-%d-%d", count, caseID), got))
		}
	}
	spellbookCapture(t, "path", rows, "4babc9f66cb15e0efe1839e97ce8a6fbde59112d244c62701f77778ac1f0612e")
}
func TestSpellbookPalette(t *testing.T) {
	o := newSpellbookOwner(t)
	o.resetBook(t)
	if got := o.bookCall("nox_xxx_bookSetColor_45AC40"); got == 0 {
		t.Fatal("book palette")
	}
	spellbookCapture(t, "palette", []spellbookResult{o.bookSnapshot("palette", 0)}, "da8e6beca619d11d3d611a6edfbae1e88f6dc235c8cbcc9373f0f03b6fceb61a")
}
