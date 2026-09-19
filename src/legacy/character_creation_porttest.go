//go:build porttest

package legacy

/*
#include "GAME3.h"
#include "client__shell__selcolor.h"
extern uint32_t dword_5d4594_1308084;
extern uint32_t dword_5d4594_1308088;
extern uint32_t dword_5d4594_1307792;
*/
import "C"
import "unsafe"

func PortTestCharacterName(p unsafe.Pointer) int { return int(C.sub_4A6B50((*C.wchar2_t)(p))) }

func PortTestCharacterPaletteOwner(root, menu unsafe.Pointer) func() {
	a, b, c := C.dword_5d4594_1308084, C.dword_5d4594_1308088, C.dword_5d4594_1307792
	C.dword_5d4594_1308084 = C.uint32_t(uintptr(root))
	C.dword_5d4594_1308088 = C.uint32_t(uintptr(menu))
	return func() { C.dword_5d4594_1308084 = a; C.dword_5d4594_1308088 = b; C.dword_5d4594_1307792 = c }
}
func PortTestCharacterPaletteMatch(w unsafe.Pointer, bank int, rgb unsafe.Pointer) uintptr {
	return uintptr(unsafe.Pointer(C.sub_4A61E0((*C.uint32_t)(w), C.int(bank), (*C.uchar)(rgb))))
}
func PortTestCharacterPaletteFill(bank uint16) uintptr {
	return uintptr(unsafe.Pointer(C.sub_4A7530(C.ushort(bank))))
}
func PortTestCharacterPaletteClose(id int, index uint16) uintptr {
	C.dword_5d4594_1307792 = C.uint32_t(id)
	return uintptr(unsafe.Pointer(C.sub_4A72D0(C.ushort(index))))
}
func PortTestCharacterPaletteDraw(w unsafe.Pointer) int {
	return int(C.sub_4A6D20(C.int(uintptr(w)), 0))
}
