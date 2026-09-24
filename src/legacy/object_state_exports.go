package legacy

/*
#include "defs.h"
*/
import "C"

func nox_xxx_inventoryGetFirst_4E7980(a C.int) C.int {
	return C.int(inventoryInt(objectFromInt(a).InvFirstItem))
}

func nox_xxx_inventoryGetNext_4E7990(a C.int) C.int {
	if a == 0 {
		return 0
	}
	return C.int(inventoryInt(objectFromInt(a).InvNextItem))
}
