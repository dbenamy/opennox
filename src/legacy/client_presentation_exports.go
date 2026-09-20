package legacy

/*
#include <stdint.h>
*/
import "C"
import "unsafe"

//export nox_xxx_clientAddRayEffect_49C160
func nox_xxx_clientAddRayEffect_49C160(event C.int) {
	presentationRayAdd((*[7]byte)(unsafe.Pointer(uintptr(uint32(event)))))
}

//export nox_xxx_clientRemoveRayEffect_49C450
func nox_xxx_clientRemoveRayEffect_49C450(event C.int) {
	presentationRayRemove((*[7]byte)(unsafe.Pointer(uintptr(uint32(event)))))
}

//export nox_xxx_fxDrawTurnUndead_499880
func nox_xxx_fxDrawTurnUndead_499880(pos *C.short) {
	presentationTurnUndead((*[2]int16)(unsafe.Pointer(pos)))
}
