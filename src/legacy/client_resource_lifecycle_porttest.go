//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestClientResourceWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"frameCount":  (*uint32)(unsafe.Pointer(&dword_5d4594_815748)),
		"modalWindow": (*uint32)(unsafe.Pointer(&dword_5d4594_816412)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestClientFrameAverage()      { clientFrameAverage() }
func PortTestClientShellState() int32  { return int32(memmap.Uint32(0x5D4594, 815092)) }
func PortTestClientBindingCount() byte { return memmap.Uint8(0x5D4594, 1193128) }
