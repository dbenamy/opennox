package legacy

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

//export sub_4AD840
func sub_4AD840() C.int { return C.int(serverPanelsGeneralRefresh()) }

//export sub_4535E0
func sub_4535E0(p *C.int) *C.int {
	return (*C.int)(unsafe.Pointer(serverPanelsWeaponStore((*uint32)(unsafe.Pointer(p)))))
}

//export sub_4535F0
func sub_4535F0(v C.int) C.int { return C.int(serverPanelsArmorStore(uint32(v))) }

//export sub_4536B0
func sub_4536B0(p *C.uint32_t) C.int {
	return C.int(serverPanelsWeaponSnapshot((*uint32)(unsafe.Pointer(p))))
}

//export sub_453710
func sub_453710() C.int { return C.int(serverPanelsArmorSnapshot()) }

//export sub_453F70
func sub_453F70(p unsafe.Pointer) { serverPanelsSpellStore((*uint32)(p)) }

//export sub_454040
func sub_454040(p *C.uint32_t) C.int {
	return C.int(serverPanelsSpellSnapshot((*uint32)(unsafe.Pointer(p))))
}

//export sub_4540E0
func sub_4540E0(p C.int) C.int {
	return C.int(serverPanelsSpellApply((*uint32)(unsafe.Pointer(uintptr(uint32(p))))))
}

//export sub_455800
func sub_455800() *C.int { return (*C.int)(unsafe.Pointer(serverPanelsAccessRefresh())) }
