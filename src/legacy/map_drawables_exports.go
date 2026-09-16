package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_client_mapSpecialRWObjectData_4AC610
func nox_client_mapSpecialRWObjectData_4AC610() C.int { return C.int(mapDrawableSection()) }
