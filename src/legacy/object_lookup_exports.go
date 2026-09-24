package legacy

/*
#include "GAME3_3.h"
*/
import "C"

func nox_server_getObjectFromNetCode_4ECCB0(code C.int) *C.nox_object_t {
	return (*C.nox_object_t)(objectLookupByNetCode(uint32(code)).CObj())
}
