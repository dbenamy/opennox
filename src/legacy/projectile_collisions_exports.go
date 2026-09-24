package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4E9A30
func sub_4E9A30(a, b *nox_object_t) C.int {
	return C.int(bool2int(projectileTrapEligible(asObjectS(a), asObjectS(b))))
}

//export sub_4EB250
func sub_4EB250(a C.int) C.int { return C.int(inventoryInt(projectileChakramSelect(objectFromInt(a)))) }

//export sub_4EB340
func sub_4EB340(a *C.float, b C.int) {
	projectileChakramCandidate((*server.Object)(unsafe.Pointer(a)), (*types.Pointf)(unsafe.Pointer(uintptr(b))))
}

//export sub_4EB3E0
func sub_4EB3E0(a C.int) { projectileChakramFallback(objectFromInt(a)) }
