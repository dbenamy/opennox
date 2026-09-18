package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_xxx_serverHandleClientConsole_443E90
func nox_xxx_serverHandleClientConsole_443E90(pl *C.nox_playerInfo, action C.char, text *C.wchar2_t) C.int {
	return C.int(consoleCommandRemote(asPlayerS(pl), byte(action), GoWString(text)))
}
