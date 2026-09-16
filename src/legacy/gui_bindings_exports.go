package legacy

/*
#include "defs.h"
*/
import "C"

//export sub_4C35B0
func sub_4C35B0(cancel C.int) C.int { return C.int(bindingClose(int(cancel))) }

//export sub_4C4260
func sub_4C4260() { bindingShow() }

//export sub_4C4280
func sub_4C4280() C.int { return C.int(bindingVisible()) }

//export sub_4CB880
func sub_4CB880() C.int { return C.int(bindingMenu.construct()) }

//export sub_4CBB70
func sub_4CBB70() C.int { return C.int(bindingMenuBack()) }

//export sub_4CBBB0
func sub_4CBBB0() C.int { return C.int(bindingMenuDone()) }
