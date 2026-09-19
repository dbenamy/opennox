//go:build porttest

package legacy

/*
#include "GAME3.h"
#include "client__shell__selclass.h"
extern uint32_t dword_5d4594_1307736;
*/
import "C"
import "unsafe"

func PortTestCharacterClassWindow(root unsafe.Pointer) func() {
	old := C.dword_5d4594_1307736
	C.dword_5d4594_1307736 = C.uint32_t(uintptr(root))
	return func() { C.dword_5d4594_1307736 = old }
}
func PortTestCharacterClassEvent(root, child unsafe.Pointer, event int) int {
	return int(C.sub_4A4A20(C.int(uintptr(root)), C.int(event), (*C.int)(child), 0))
}
func PortTestCharacterClassDraw(w unsafe.Pointer) int { return int(C.sub_4A49D0(C.int(uintptr(w)), 0)) }
