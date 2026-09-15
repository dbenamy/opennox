//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

// installListboxPalette uses the production initializer and restores every cell
// it owns, including the pointer table consumed by the C listbox callbacks.
func installListboxPalette(t *testing.T) {
	cells := []unsafe.Pointer{
		unsafe.Pointer(legacy.Get_nox_color_black_2650656_ptr()),
		memmap.PtrOff(0x852978, 4),
		memmap.PtrOff(0x85B3FC, 956),
		memmap.PtrOff(0x5D4594, 2597996),
		unsafe.Pointer(legacy.Get_nox_color_white_2523948_ptr()),
		unsafe.Pointer(legacy.Get_nox_color_violet_2598268_ptr()),
		memmap.PtrOff(0x85B3FC, 940),
		unsafe.Pointer(legacy.Get_nox_color_red_2589776_ptr()),
		memmap.PtrOff(0x85B3FC, 984),
		unsafe.Pointer(legacy.Get_dword_8531A0_2572_ptr()),
		unsafe.Pointer(legacy.Get_nox_color_green_2614268_ptr()),
		memmap.PtrOff(0x85B3FC, 944),
		unsafe.Pointer(legacy.Get_nox_color_cyan_2649820_ptr()),
		unsafe.Pointer(legacy.Get_nox_color_blue_2650684_ptr()),
		unsafe.Pointer(legacy.Get_nox_color_orange_2614256_ptr()),
		unsafe.Pointer(legacy.Get_nox_color_yellow_2589772_ptr()),
		memmap.PtrOff(0x852978, 0),
	}
	var oldValues [17]uint32
	var oldPointers [17]unsafe.Pointer
	for i, p := range cells {
		oldValues[i] = *(*uint32)(p)
		oldPointers[i] = *memmap.PtrPtr(0x85B3FC, 132+uintptr(i*4))
	}
	nox_xxx_loadDefColor_4A94A0()
	t.Cleanup(func() {
		for i, p := range cells {
			*(*uint32)(p) = oldValues[i]
			*memmap.PtrPtr(0x85B3FC, 132+uintptr(i*4)) = oldPointers[i]
		}
	})
}
