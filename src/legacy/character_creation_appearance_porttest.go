//go:build porttest

package legacy

/*
#include "client__shell__selcolor.h"
extern uint32_t dword_5d4594_1307784;
extern uint32_t dword_5d4594_1308096;
extern uint32_t dword_5d4594_1308100;
extern uint32_t dword_5d4594_1308104;
extern uint32_t dword_5d4594_1308108;
extern uint32_t dword_5d4594_1308112;
extern uint32_t dword_5d4594_1308116;
extern uint32_t dword_5d4594_1308120;
extern uint32_t dword_5d4594_1308124;
extern uint32_t dword_5d4594_1308128;
extern uint32_t dword_5d4594_1308132;
extern uint32_t dword_5d4594_1308136;
extern uint32_t dword_5d4594_1308140;
extern uint32_t dword_5d4594_1308144;
extern uint32_t dword_5d4594_1308148;
extern uint32_t dword_5d4594_1308152;
*/
import "C"
import "unsafe"

func PortTestCharacterAppearanceOwner(player unsafe.Pointer, windows []unsafe.Pointer) func() {
	words := []*C.uint32_t{&C.dword_5d4594_1307784, &C.dword_5d4594_1308096, &C.dword_5d4594_1308100, &C.dword_5d4594_1308104, &C.dword_5d4594_1308108, &C.dword_5d4594_1308112, &C.dword_5d4594_1308116, &C.dword_5d4594_1308120, &C.dword_5d4594_1308124, &C.dword_5d4594_1308128, &C.dword_5d4594_1308132, &C.dword_5d4594_1308136, &C.dword_5d4594_1308140, &C.dword_5d4594_1308144, &C.dword_5d4594_1308148, &C.dword_5d4594_1308152}
	saved := make([]C.uint32_t, len(words))
	for i, p := range words {
		saved[i] = *p
		if i == 0 {
			*p = C.uint32_t(uintptr(player))
		} else {
			*p = C.uint32_t(uintptr(windows[i-1]))
		}
	}
	return func() {
		for i, p := range words {
			*p = saved[i]
		}
	}
}
func PortTestCharacterAppearance() unsafe.Pointer { return unsafe.Pointer(C.sub_4A68C0()) }

func PortTestCharacterCreateFile() int { return int(C.sub_4A75C0()) }
