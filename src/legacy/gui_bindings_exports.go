package legacy

/*
#include "defs.h"
*/
import "C"

//export sub_4CB880
func sub_4CB880() C.int { return C.int(bindingMenu.construct()) }

//export sub_4CBB70
func sub_4CBB70() C.int { return C.int(bindingMenuBack()) }

//export sub_4CBBB0
func sub_4CBBB0() C.int { return C.int(bindingMenuDone()) }
