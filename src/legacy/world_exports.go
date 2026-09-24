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

//export nox_xxx_fnElevatorShaft_53B410
func nox_xxx_fnElevatorShaft_53B410(a1 C.int, a2 C.int) {
	worldShaftCandidate(objectFromInt(a1), objectFromInt(a2))
}

//export nox_xxx_elevatorAud_53B490
func nox_xxx_elevatorAud_53B490(a1 C.int, a2 C.int) { worldElevatorSound(objectFromInt(a1), a2 != 0) }

//export nox_xxx_elevatorFn_53B750
func nox_xxx_elevatorFn_53B750(a1 C.int, a2 C.int) {
	worldElevatorCandidate(objectFromInt(a1), objectFromInt(a2))
}

//export nox_xxx_fnPentagramTeleport_53C060
func nox_xxx_fnPentagramTeleport_53C060(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), true)
}

//export sub_53C140
func sub_53C140(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), false)
}

//export sub_53C240
func sub_53C240(a1 *C.float, arg4 C.int) {
	worldBlowCandidate((*server.Object)(unsafe.Pointer(a1)), objectFromInt(arg4))
}

//export sub_548830
func sub_548830(a1 C.int) { worldAngleQueue(unsafe.Pointer(uintptr(uint32(a1)))) }

//export sub_548860
func sub_548860(a1 C.int, a2 C.short) { worldAngle(objectFromInt(a1), int16(a2)) }
