package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"unsafe"
)

//export sub_547DB0
func sub_547DB0(a C.int, p *C.float2) C.int {
	return C.int(collisionObjectContains(objectFromInt(a), (*types.Pointf)(unsafe.Pointer(p))))
}
