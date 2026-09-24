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

//export nox_xxx_collideMonsterEventProc_4E83B0
func nox_xxx_collideMonsterEventProc_4E83B0(a, b C.int) *C.uchar {
	return (*C.uchar)(stateMonsterCollision(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collideMimic_4E83D0
func nox_xxx_collideMimic_4E83D0(a, b C.int) *C.uchar {
	return (*C.uchar)(stateMimicCollision(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_collidePlayer_4E8460
func nox_xxx_collidePlayer_4E8460(a, b C.int) {
	statePlayerCollision(objectFromInt(a), objectFromInt(b))
}

//export nox_objectCollideDefault
func nox_objectCollideDefault(a, b C.int, c *C.float) C.int { return 0 }
