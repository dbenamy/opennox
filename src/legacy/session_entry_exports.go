package legacy

/*
#include "defs.h"
typedef const char session_entry_const_char;
*/
import "C"
import (
	"unsafe"
)

//export nox_server_currentMapGetFilename_409B30
func nox_server_currentMapGetFilename_409B30() *C.char {
	return (*C.char)(unsafe.Pointer(sessionMapFilename()))
}

//export nox_xxx_gameSetMapPath_409D70
func nox_xxx_gameSetMapPath_409D70(name *C.char) *C.char {
	return (*C.char)(unsafe.Pointer(sessionSetMapPath((*byte)(unsafe.Pointer(name)))))
}
