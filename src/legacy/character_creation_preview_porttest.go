//go:build porttest

package legacy

/*
#include "GAME3.h"
extern void* dword_5d4594_1308156;
extern void* dword_5d4594_1308160;
extern void* dword_5d4594_1308164;
*/
import "C"
import "unsafe"

func PortTestCharacterPreviewOwner(pants, shirt, shoes unsafe.Pointer) func() {
	a, b, c := C.dword_5d4594_1308156, C.dword_5d4594_1308160, C.dword_5d4594_1308164
	C.dword_5d4594_1308156 = pants
	C.dword_5d4594_1308160 = shirt
	C.dword_5d4594_1308164 = shoes
	return func() { C.dword_5d4594_1308156 = a; C.dword_5d4594_1308160 = b; C.dword_5d4594_1308164 = c }
}
func PortTestCharacterPreview(w unsafe.Pointer) int { return int(C.sub_4A6DC0((*C.uint32_t)(w), 0)) }
