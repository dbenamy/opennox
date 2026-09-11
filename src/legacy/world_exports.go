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

//export nox_xxx_updateDoor_53AC50
func nox_xxx_updateDoor_53AC50(a1 C.int) C.char { return C.char(worldDoor(objectFromInt(a1))) }

//export nox_xxx_updatePush_53B030
func nox_xxx_updatePush_53B030(a1 C.int) { worldPush(objectFromInt(a1)) }

//export nox_xxx_updateToggle_53B060
func nox_xxx_updateToggle_53B060(a1 *C.uint32_t) C.char {
	return C.char(worldToggle((*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_updateTrigger_53B1B0
func nox_xxx_updateTrigger_53B1B0(a1 C.int) C.char { return C.char(worldTrigger(objectFromInt(a1))) }

//export sub_53B300
func sub_53B300(a1 C.int) C.char { return C.char(worldEnabledCollision(objectFromInt(a1))) }

//export nox_xxx_updateSwitch_53B320
func nox_xxx_updateSwitch_53B320(a1 *C.uint32_t) C.char {
	return C.char(worldSwitch((*server.Object)(unsafe.Pointer(a1))))
}

//export nox_xxx_updateElevatorShaft_53B380
func nox_xxx_updateElevatorShaft_53B380(a1 C.int) C.char {
	return C.char(worldShaft(objectFromInt(a1)))
}

//export nox_xxx_fnElevatorShaft_53B410
func nox_xxx_fnElevatorShaft_53B410(a1 C.int, a2 C.int) {
	worldShaftCandidate(objectFromInt(a1), objectFromInt(a2))
}

//export nox_xxx_elevatorAud_53B490
func nox_xxx_elevatorAud_53B490(a1 C.int, a2 C.int) { worldElevatorSound(objectFromInt(a1), a2 != 0) }

//export nox_xxx_updateElevator_53B5D0
func nox_xxx_updateElevator_53B5D0(a1 *C.uint32_t) {
	worldElevator((*server.Object)(unsafe.Pointer(a1)))
}

//export nox_xxx_elevatorFn_53B750
func nox_xxx_elevatorFn_53B750(a1 C.int, a2 C.int) {
	worldElevatorCandidate(objectFromInt(a1), objectFromInt(a2))
}

//export nox_xxx_updatePhantomPlayer_53B860
func nox_xxx_updatePhantomPlayer_53B860(a1 C.int) { worldPhantom(objectFromInt(a1)) }

//export nox_xxx_updateTeleportPentagram_53BEF0
func nox_xxx_updateTeleportPentagram_53BEF0(a1 C.int) C.int {
	return C.int(worldTeleport(objectFromInt(a1)))
}

//export nox_xxx_fnPentagramTeleport_53C060
func nox_xxx_fnPentagramTeleport_53C060(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), true)
}

//export nox_xxx_updateInvisiblePentagram_53C0C0
func nox_xxx_updateInvisiblePentagram_53C0C0(a1 C.int) C.int {
	return C.int(worldInvisibleTeleport(objectFromInt(a1)))
}

//export sub_53C140
func sub_53C140(a1 *C.float, a2 C.int) {
	worldTeleportCandidate((*server.Object)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(a2)))), false)
}

//export nox_xxx_updateBlow_53C160
func nox_xxx_updateBlow_53C160(a3 C.int) { worldBlow(objectFromInt(a3)) }

//export sub_53C240
func sub_53C240(a1 *C.float, arg4 C.int) {
	worldBlowCandidate((*server.Object)(unsafe.Pointer(a1)), objectFromInt(arg4))
}

//export nox_xxx_updateTrapDoor_53DE80
func nox_xxx_updateTrapDoor_53DE80(a1 *C.uint32_t) *C.int {
	return (*C.int)(unsafe.Pointer(uintptr(worldTrapDoor((*server.Object)(unsafe.Pointer(a1))))))
}

//export sub_548830
func sub_548830(a1 C.int) { worldAngleQueue(unsafe.Pointer(uintptr(uint32(a1)))) }

//export sub_548860
func sub_548860(a1 C.int, a2 C.short) { worldAngle(objectFromInt(a1), int16(a2)) }
