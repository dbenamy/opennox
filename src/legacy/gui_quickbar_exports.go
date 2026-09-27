package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func sub_460380() unsafe.Pointer { return unsafe.Pointer(uintptr(quickbarClearAbilities())) }

//export nox_xxx_quickbarButtonBook_45F3F0
func nox_xxx_quickbarButtonBook_45F3F0() C.int { return C.int(quickbarBookTooltip()) }

//export sub_45F480
func sub_45F480(w C.int) C.int { return C.int(quickbarDirectionTooltip(bookWindow(uint32(w)))) }
