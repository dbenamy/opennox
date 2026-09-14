package legacy

/*
#include "GAME3_3.h"
*/
import "C"

//export nox_server_getObjectFromNetCode_4ECCB0
func nox_server_getObjectFromNetCode_4ECCB0(code C.int) *C.nox_object_t {
	return (*C.nox_object_t)(objectLookupByNetCode(uint32(code)).CObj())
}

//export sub_4ECF10
func sub_4ECF10(id C.int) C.int {
	return C.int(uintptr(objectLookupByScriptID(uint32(id)).CObj()))
}
