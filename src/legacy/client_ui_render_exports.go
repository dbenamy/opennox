package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_client_copyRect_49F6F0
func nox_client_copyRect_49F6F0(x, y, w, h C.int) C.int {
	return C.int(uiRenderCopyRect(int(x), int(y), int(w), int(h)))
}
