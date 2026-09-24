package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_53BD10
func sub_53BD10(a1 C.int, a2 C.int) { temporaryAntiCandidate(objectFromInt(a1), objectFromInt(a2)) }

//export nox_xxx_waterBarrel_53CC30
func nox_xxx_waterBarrel_53CC30(a1 *C.float, a2 C.int) {
	temporaryWaterCandidate((*server.Object)(unsafe.Pointer(a1)), *(*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))))
}

//export nox_xxx_updateFlameCleanse_53D510
func nox_xxx_updateFlameCleanse_53D510(a1 C.int) { temporaryFlameCleanse(objectFromInt(a1)) }

//export sub_53D8C0
func sub_53D8C0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(a1), objectFromInt(a2), false)
}

//export nox_xxx_toxicCloudPoison_53D9D0
func nox_xxx_toxicCloudPoison_53D9D0(a1 C.int, a2 C.int) {
	temporaryCloudCandidate(objectFromInt(a1), objectFromInt(a2), true)
}

//export nox_xxx_createSpark_54FD80
func nox_xxx_createSpark_54FD80(a1 C.float, a2 C.float, a3 C.int, a4 C.int, a5 C.float, a6 C.float, a7 C.float, a8 C.int) *C.float {
	return (*C.float)(temporarySpark(types.Ptf(float32(a1), float32(a2)), types.Ptf(float32(a5), float32(a6)), int32(a3), int32(a4), float32(a7), objectFromInt(a8)).CObj())
}
