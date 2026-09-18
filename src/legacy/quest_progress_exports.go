package legacy

/*
#include "defs.h"
*/
import "C"

//export sub_5007E0
func sub_5007E0(name *C.char) *C.char { questProgressReset(GoString(name)); return nil }

//export sub_500A60
func sub_500A60() C.int { return C.int(questProgressWrite()) }

//export sub_500B70
func sub_500B70() C.int { return C.int(questProgressRead()) }
