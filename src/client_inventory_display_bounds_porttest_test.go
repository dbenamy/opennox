//go:build porttest

package opennox

import (
	"image"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Native contract: the old C path's unterminated 256-unit temporary and ensuing
// concatenation were out of bounds for these inputs, so they are not C oracles.
func TestClientInventoryDisplayLongIdentifyHeader(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	// Borrow the header plus two guard units of raw mapped storage. The old
	// address immediately after this buffer is now an unused hole: the font
	// global was relocated. Check its actual owner separately below.
	storage := tooltipMapped(t, 1063124, 258)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	for _, name := range []string{strings.Repeat("A", 246), strings.Repeat("A", 247), strings.Repeat("A", 248), strings.Repeat("A", 300), strings.Repeat("Ω", 300), strings.Repeat("😀", 200)} {
		o.reset(t)
		o.inventoryRects()
		o.identifyWindows(t)
		o.c.Mouse = image.Pt(16, 22)
		dr := o.item(t, "RedApple", 123)
		o.c.Things.TypeByID("RedApple").PrettyName = alloc.InternCString16(name)
		*o.displayWords["dword_5d4594_1063116"] = uint32(uintptr(dr.C()))
		header := storage[:256]
		old := [2]uint16{storage[256], storage[257]}
		storage[256], storage[257] = 0x5678, 0x1234
		font := *o.displayWords["dword_5d4594_1063636"]
		o.call(t, 0, 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
		gotGuard := [2]uint16{storage[256], storage[257]}
		storage[256], storage[257] = old[0], old[1]
		if *o.displayWords["dword_5d4594_1063636"] != font {
			t.Fatal("identify header changed the relocated font global")
		}
		if gotGuard != [2]uint16{0x5678, 0x1234} {
			t.Fatal("identify header changed adjacent storage")
		}
		want := utf16.Encode([]rune("Identify " + name))
		for i := 0; i < min(len(want), 255); i++ {
			if header[i] != want[i] {
				t.Fatalf("header unit%d got%x want%x", i, header[i], want[i])
			}
		}
		if header[min(len(want), 255)] != 0 {
			t.Fatal("identify header lacks terminator")
		}
	}
}
