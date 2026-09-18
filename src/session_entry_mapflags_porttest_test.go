//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestSessionEntryMapFlags(t *testing.T) {
	record, free := alloc.New([8]uint32{})
	defer free()
	for i := 0; i < 7; i++ {
		record[i] = 0xa5a50000 + uint32(i)
	}
	type row struct {
		Input  uint32
		Output int32
	}
	var rows []row
	// File flags have a distinct encoding from the game-mode bits.
	modes := [7]uint32{512, 4096, 256, 32, 16, 1024, 64}
	for mask := uint32(0); mask < 128; mask++ {
		for _, extra := range []uint32{0, 0x7fffff80, 0x80000000, 0xffffff80} {
			input := mask | extra
			record[7] = input
			got := legacy.PortTestSessionEntryMapFlags(unsafe.Pointer(record))
			want := uint32(0)
			for bit, mode := range modes {
				if mask&(1<<uint(bit)) != 0 {
					want |= mode
				}
			}
			if input&0x80000000 != 0 {
				want |= 128
			}
			if got != int32(want) {
				t.Fatalf("map flags%x got%x want%x", input, got, want)
			}
			for i := 0; i < 7; i++ {
				if record[i] != 0xa5a50000+uint32(i) {
					t.Fatal("map record mutated")
				}
			}
			if record[7] != input {
				t.Fatal("map flags mutated")
			}
			rows = append(rows, row{input, got})
		}
	}
	spellbookCapture(t, "session-entry-map-flags", rows, "ba4aceabb3dde39ad66ef1124a15ad7f84e225252a3d222faa8b54dc5edfeee6")
}
