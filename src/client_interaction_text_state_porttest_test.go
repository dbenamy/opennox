//go:build porttest

package opennox

import (
	"encoding/binary"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionTextState(t *testing.T) {
	dst := serverConfigOwnBytes(t, 0x5D4594, 811376, 528)
	mode := serverConfigOwnBytes(t, 0x5D4594, 811060, 4)
	type row struct {
		Length int
		Mode   uint32
		Text   []uint16
	}
	var rows []row
	for _, length := range []int{0, 1, 255, 263} {
		text, free := alloc.Make([]uint16{}, length+1)
		for i := 0; i < length; i++ {
			text[i] = uint16(0x101 + i)
		}
		for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
			for i := range dst {
				dst[i] = 0x5a
			}
			ret := interactionCall("sub_435700", uintptr(unsafe.Pointer(&text[0])), uintptr(value))
			if uintptr(ret) != uintptr(unsafe.Pointer(&dst[0])) || binary.LittleEndian.Uint32(mode) != value {
				t.Fatal("text state result", length, value, ret)
			}
			got := make([]uint16, length+1)
			for i := range got {
				got[i] = binary.LittleEndian.Uint16(dst[2*i:])
			}
			if !slices.Equal(got, text) {
				t.Fatal("text units", length, value)
			}
			for _, v := range dst[2*(length+1):] {
				if v != 0x5a {
					t.Fatal("text tail modified", length)
				}
			}
			rows = append(rows, row{length, value, got})
		}
		free()
	}
	interactionCapture(t, "text-state", rows)
}
