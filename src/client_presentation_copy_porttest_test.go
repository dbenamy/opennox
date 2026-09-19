//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientPresentationCopy(t *testing.T) {
	type record struct {
		Length, SourceOffset, DestinationOffset int
		Bytes                                   string
	}
	var rows []record
	source, freeSource := alloc.Make([]byte{}, 96)
	defer freeSource()
	dst, freeDst := alloc.Make([]byte{}, 96)
	defer freeDst()
	for i := range source {
		source[i] = byte(i*37 + 11)
	}
	for _, n := range []int{0, 2, 4, 6, 8, 30, 32, 34, 62, 64} {
		for _, so := range []int{0, 1, 2, 3} {
			for _, do := range []int{0, 1, 2, 3} {
				for i := range dst {
					dst[i] = 0xa5
				}
				want := append([]byte(nil), dst...)
				copy(want[8+do:8+do+n], source[8+so:8+so+n])
				legacy.PortTestPresentationCopy(unsafe.Pointer(&dst[8+do]), unsafe.Pointer(&source[8+so]), uint32(n))
				if !bytes.Equal(dst, want) {
					t.Fatalf("length=%d offsets=%d/%d copy/guard mismatch", n, so, do)
				}
				rows = append(rows, record{n, so, do, fmt.Sprintf("%x", dst)})
			}
		}
	}
	spellbookCapture(t, "client-presentation-copy", rows, "bd2af18d4698d259b904b7358337abf5cd841072f4a90faf1c08ba5f90294257")
}
