package legacy

/*
#include "defs.h"
typedef const char session_entry_const_char;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4CFE00
func sub_4CFE00() C.int { return C.int(sessionMapState()) }

//export nox_xxx_action_4DA9F0
func nox_xxx_action_4DA9F0(u *C.nox_object_t) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(sessionShadowRemove((*server.Object)(unsafe.Pointer(u)))))
}

//export nox_client_countPlayerFiles02_4DC630
func nox_client_countPlayerFiles02_4DC630() C.int { return C.int(sessionCharacterCount()) }

//export nox_xxx_game_4DCCB0
func nox_xxx_game_4DCCB0() C.int { return C.int(sessionSaveAllowed()) }

//export nox_server_currentMapGetFilename_409B30
func nox_server_currentMapGetFilename_409B30() *C.char {
	return (*C.char)(unsafe.Pointer(sessionMapFilename()))
}

//export nox_xxx_mapGetMapName_409B40
func nox_xxx_mapGetMapName_409B40() *C.char { return (*C.char)(unsafe.Pointer(sessionMapName())) }

//export sub_409B50
func sub_409B50(name *C.session_entry_const_char) C.uint {
	return C.uint(sessionSetSelectedMap((*byte)(unsafe.Pointer(name))))
}

//export nox_xxx_gameSetMapPath_409D70
func nox_xxx_gameSetMapPath_409D70(name *C.char) *C.char {
	return (*C.char)(unsafe.Pointer(sessionSetMapPath((*byte)(unsafe.Pointer(name)))))
}
