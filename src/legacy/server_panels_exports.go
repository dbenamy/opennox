package legacy

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

func sub_4540E0(p C.int) C.int {
	return C.int(serverPanelsSpellApply((*uint32)(unsafe.Pointer(uintptr(uint32(p))))))
}
