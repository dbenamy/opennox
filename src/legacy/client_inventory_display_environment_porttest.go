//go:build porttest

package legacy

import "unsafe"

// PortTestInventoryDisplayWords borrows the real shared display globals.
func PortTestInventoryDisplayWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1049844":     (*uint32)(unsafe.Pointer(&dword_5d4594_1049844)),
		"dword_5d4594_1049864":     (*uint32)(unsafe.Pointer(&dword_5d4594_1049864)),
		"dword_5d4594_1062456":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062456)),
		"dword_5d4594_1062476":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062476)),
		"dword_5d4594_1062480":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062480)),
		"dword_5d4594_1062484":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062484)),
		"dword_5d4594_1062488":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062488)),
		"dword_5d4594_1062492":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062492)),
		"dword_5d4594_1062496":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062496)),
		"dword_5d4594_1062512":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062512)),
		"dword_5d4594_1062552":     (*uint32)(unsafe.Pointer(&dword_5d4594_1062552)),
		"dword_5d4594_1063116":     (*uint32)(unsafe.Pointer(&dword_5d4594_1063116)),
		"dword_5d4594_1063120":     (*uint32)(unsafe.Pointer(&dword_5d4594_1063120)),
		"dword_5d4594_1063636":     (*uint32)(unsafe.Pointer(&dword_5d4594_1063636)),
		"dword_8531A0_2576":        (*uint32)(unsafe.Pointer(&dword_8531A0_2576)),
		"nox_color_black_2650656":  (*uint32)(unsafe.Pointer(&nox_color_black_2650656)),
		"nox_color_blue_2650684":   (*uint32)(unsafe.Pointer(&nox_color_blue_2650684)),
		"nox_color_cyan_2649820":   (*uint32)(unsafe.Pointer(&nox_color_cyan_2649820)),
		"nox_color_orange_2614256": (*uint32)(unsafe.Pointer(&nox_color_orange_2614256)),
		"nox_color_red_2589776":    (*uint32)(unsafe.Pointer(&nox_color_red_2589776)),
		"nox_color_violet_2598268": (*uint32)(unsafe.Pointer(&nox_color_violet_2598268)),
		"nox_color_white_2523948":  (*uint32)(unsafe.Pointer(&nox_color_white_2523948)),
		"nox_color_yellow_2589772": (*uint32)(unsafe.Pointer(&nox_color_yellow_2589772)),
	}
	old := make(map[string]uint32, len(words))
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}
