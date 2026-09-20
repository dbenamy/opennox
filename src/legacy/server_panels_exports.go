package legacy

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

//export sub_4AD840
func sub_4AD840() C.int { return C.int(serverPanelsGeneralRefresh()) }

func sub_4540E0(p C.int) C.int {
	return C.int(serverPanelsSpellApply((*uint32)(unsafe.Pointer(uintptr(uint32(p))))))
}

//export sub_455800
func sub_455800() *C.int { return (*C.int)(unsafe.Pointer(serverPanelsAccessRefresh())) }
