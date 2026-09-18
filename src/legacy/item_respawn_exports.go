package legacy

/*
#include "defs.h"
*/
import "C"

//export sub_4ED050
func sub_4ED050(u, stamp C.int) { itemDropCrowns(objectFromInt(u), uint32(stamp)) }
