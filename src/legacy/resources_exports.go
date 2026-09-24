package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_xxx_pickupGold_4F3A60_obj_pickup
func nox_xxx_pickupGold_4F3A60_obj_pickup(u, item, flags C.int) C.int {
	return C.int(bool2int(resourceGoldPickup(objectFromInt(u), objectFromInt(item), int(flags))))
}
