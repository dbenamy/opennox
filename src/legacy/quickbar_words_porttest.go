//go:build porttest

package legacy

import "unsafe"

func PortTestQuickbarWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_5d4594_1047548": (*uint32)(unsafe.Pointer(&dword_5d4594_1047548)),
		"dword_5d4594_1047552": (*uint32)(unsafe.Pointer(&dword_5d4594_1047552)),
		"dword_5d4594_1047932": (*uint32)(unsafe.Pointer(&dword_5d4594_1047932)),
		"dword_5d4594_1047936": (*uint32)(unsafe.Pointer(&dword_5d4594_1047936)),
		"dword_5d4594_1049484": (*uint32)(unsafe.Pointer(&dword_5d4594_1049484)),
		"dword_5d4594_1049496": (*uint32)(unsafe.Pointer(&dword_5d4594_1049496)),
		"dword_5d4594_1049500": (*uint32)(unsafe.Pointer(&dword_5d4594_1049500)),
		"dword_5d4594_1049504": (*uint32)(unsafe.Pointer(&dword_5d4594_1049504)),
		"dword_5d4594_1049508": (*uint32)(unsafe.Pointer(&dword_5d4594_1049508)),
		"dword_5d4594_1049512": (*uint32)(unsafe.Pointer(&dword_5d4594_1049512)),
		"dword_5d4594_1049516": (*uint32)(unsafe.Pointer(&dword_5d4594_1049516)),
		"dword_5d4594_1049520": (*uint32)(unsafe.Pointer(&dword_5d4594_1049520)),
		"dword_5d4594_1049524": (*uint32)(unsafe.Pointer(&dword_5d4594_1049524)),
		"dword_5d4594_1049532": (*uint32)(unsafe.Pointer(&dword_5d4594_1049532)),
		"dword_5d4594_1049536": (*uint32)(unsafe.Pointer(&dword_5d4594_1049536)),
		"dword_5d4594_1049692": (*uint32)(unsafe.Pointer(&dword_5d4594_1049692)),
		"dword_5d4594_1049696": (*uint32)(unsafe.Pointer(&dword_5d4594_1049696)),
	}
	old := make(map[string]uint32, len(words))
	for n, p := range words {
		old[n] = *p
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}
