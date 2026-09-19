package legacy

/*
#include <stdint.h>
*/
import "C"

//export nox_gui_xxx_check_446360
func nox_gui_xxx_check_446360() C.uint { return C.uint(sessionQuitShown()) }

//export sub_445C40
func sub_445C40() { sessionQuitToggle() }

//export sub_446780
func sub_446780() C.int { return C.int(sessionMOTDClose()) }

//export sub_446950
func sub_446950() C.int { return C.int(sessionMOTDShown()) }
