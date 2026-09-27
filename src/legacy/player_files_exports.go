package legacy

/*
#include "GAME1_1.h"
*/
import "C"

func nox_xxx_plrLoad_41A480(path *C.char) C.int { return C.int(playerFileClientLoad(GoString(path))) }
