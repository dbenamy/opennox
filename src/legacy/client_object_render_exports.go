package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_wndDraw_49F7F0
func nox_xxx_wndDraw_49F7F0() { objectRenderSaveClip() }

//export sub_49F860
func sub_49F860() C.int { return C.int(objectRenderRestoreClip()) }

func sub_4C5020(packet C.int) C.int {
	return C.int(objectRenderBeamAppend(unsafe.Pointer(uintptr(packet))))
}

//export sub_4C5050
func sub_4C5050() { objectRenderBeamReset() }
