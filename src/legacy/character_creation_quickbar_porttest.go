//go:build porttest

package legacy

/*
#include "client__shell__selclass.h"
extern uint32_t dword_5d4594_1307724;
*/
import "C"
import "unsafe"

func PortTestCharacterClassOwner(p unsafe.Pointer) func() {
	old := C.dword_5d4594_1307724
	C.dword_5d4594_1307724 = C.uint32_t(uintptr(p))
	return func() { C.dword_5d4594_1307724 = old }
}
func PortTestCharacterQuickbar(mode int) uintptr { return uintptr(C.sub_4A4B70(C.int(mode))) }
