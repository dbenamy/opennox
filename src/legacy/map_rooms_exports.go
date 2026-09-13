package legacy

/*
#include <stdint.h>
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"unsafe"
)

//export nox_xxx_mapGenRoundFloatToPtr_520DF0
func nox_xxx_mapGenRoundFloatToPtr_520DF0(a1 *C.float2, a2 *C.uint32_t) C.longlong {
	return (C.longlong)(mapRoomRound((*types.Pointf)(unsafe.Pointer(a1)), (*[2]int32)(unsafe.Pointer(a2))))
}

//export sub_520EA0
func sub_520EA0(a1 C.int) C.int { return (C.int)(mapRoomGridInit(mapRoomPointer(uint32(a1)))) }

//export sub_520F80
func sub_520F80() { mapRoomGridFree() }

//export sub_521200
func sub_521200(a1 C.int) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomOverlap((*mapRoom)(mapRoomPointer(uint32(a1)))))))
}

//export sub_521290
func sub_521290(a1 *C.int2) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomAt((*[2]int32)(unsafe.Pointer(a1))))))
}

//export sub_5212B0
func sub_5212B0(a1 C.int, a2 *C.uint32_t) C.int {
	return (C.int)(mapRoomResolveOverlap((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(unsafe.Pointer(a2))))
}

//export nox_xxx_mapgenAllocBuffer_5213E0
func nox_xxx_mapgenAllocBuffer_5213E0() C.int { return (C.int)(mapRoomScratchAlloc()) }

//export nox_xxx_mapgenFreeBuffer_521400
func nox_xxx_mapgenFreeBuffer_521400() { mapRoomScratchFree() }

//export nox_xxx_mapGenGetTopRoom_521710
func nox_xxx_mapGenGetTopRoom_521710() unsafe.Pointer { return unsafe.Pointer(mapRoomHead()) }

//export sub_521720
func sub_521720(a1 C.int) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomNext((*mapRoom)(mapRoomPointer(uint32(a1)))))))
}

//export nox_xxx_mapGenAddNewRoom_521730
func nox_xxx_mapGenAddNewRoom_521730(a1 *C.uint32_t) C.int {
	return (C.int)(mapRoomAdd((*mapRoom)(unsafe.Pointer(a1))))
}

//export sub_521760
func sub_521760(a1 C.int) C.int {
	return (C.int)(mapRoomRemove((*mapRoom)(mapRoomPointer(uint32(a1)))))
}

//export sub_5217A0
func sub_5217A0(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomWithinBounds(mapRoomPointer(uint32(a1)), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export sub_521820
func sub_521820(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomCanPlace(mapRoomPointer(uint32(a1)), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export nox_xxx_mapGenSetRoomPos_521880
func nox_xxx_mapGenSetRoomPos_521880(a1 *C.uint32_t, a2 *C.float2) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomSetPos((*mapRoom)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(a2))))))
}

//export sub_5218B0
func sub_5218B0(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomHasRoomNeighbor((*mapRoom)(mapRoomPointer(uint32(a1))), int32(a2)))
}

//export sub_521900
func sub_521900(a1 C.int, a2 C.int, a3 C.int) C.int {
	return (C.int)(mapRoomConnect((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2))), int32(a3)))
}

//export nox_xxx_mapGenMakeRoomStruct_521940
func nox_xxx_mapGenMakeRoomStruct_521940(a1 C.int, a2 C.int) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomNew(int32(a1), int32(a2))))
}

//export nox_xxx_mapGenPrepareRoom_521990
func nox_xxx_mapGenPrepareRoom_521990(a1 C.int) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomPrepare(mapRoomPointer(uint32(a1)))))
}

//export sub_521A10
func sub_521A10(lpMem unsafe.Pointer) { mapRoomFree((*mapRoom)(unsafe.Pointer(lpMem))) }

//export nox_xxx_mapGenFreeTopRoom_521A40
func nox_xxx_mapGenFreeTopRoom_521A40() *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(mapRoomFreeAll()))
}

//export sub_521A70
func sub_521A70(a1 C.int, a2 C.int, a3 C.int) C.int {
	return (C.int)(mapRoomConnectBoth((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2))), int32(a3)))
}

//export sub_521AA0
func sub_521AA0(a1 *C.uint32_t, a2 C.int) C.int {
	return (C.int)(mapRoomIsEntranceSide((*mapRoom)(unsafe.Pointer(a1)), int32(a2)))
}

//export sub_521B00
func sub_521B00(a1 C.int, a2 C.int) C.double {
	return (C.double)(mapRoomRandomX((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export sub_521B30
func sub_521B30(a1 C.int, a2 C.int) C.double {
	return (C.double)(mapRoomRandomY((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export sub_521B60
func sub_521B60(a1 C.int, a2 C.int) C.double {
	return (C.double)(mapRoomRandomReverseX((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export sub_521B90
func sub_521B90(a1 C.int, a2 C.int) C.double {
	return (C.double)(mapRoomRandomReverseY((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(mapRoomPointer(uint32(a2)))))
}

//export sub_521BC0
func sub_521BC0(a1 C.int, a2 *C.float2, a3 C.float, a4 C.float) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomAddExclusion((*mapRoom)(mapRoomPointer(uint32(a1))), (*types.Pointf)(unsafe.Pointer(a2)), float32(a3), float32(a4))))
}

//export sub_521C10
func sub_521C10(a1 C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(mapRoomRemoveGeneratedExclusions((*mapRoom)(mapRoomPointer(uint32(a1))))))
}

//export sub_521EB0
func sub_521EB0(a1 *C.float, a2 *C.float) C.int {
	return (C.int)(mapRoomContainsRect((*mapRoom)(unsafe.Pointer(a1)), (*mapRoomExclusion)(unsafe.Pointer(a2))))
}

//export sub_521F10
func sub_521F10(a1 C.int, a2 *C.float) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomFindExclusionOverlap((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoomExclusion)(unsafe.Pointer(a2))))))
}

//export sub_5226D0
func sub_5226D0(a1 C.int, a2 C.float, a3 C.int) C.int {
	return (C.int)(mapRoomRandomPoint((*mapRoom)(mapRoomPointer(uint32(a1))), float32(a2), (*types.Pointf)(mapRoomPointer(uint32(a3)))))
}

//export sub_522CA0
func sub_522CA0(a1 C.int, a2 *C.float) C.char {
	return (C.char)(mapRoomRememberPoint((*mapRoom)(mapRoomPointer(uint32(a1))), (*types.Pointf)(unsafe.Pointer(a2))))
}

//export nox_xxx_mapGenCheckRoomType_5238F0
func nox_xxx_mapGenCheckRoomType_5238F0(a1 *C.int) C.int {
	return (C.int)(mapRoomIsHall((*mapRoom)(unsafe.Pointer(a1))))
}

//export sub_523920
func sub_523920(a1 C.int) C.int { return (C.int)(mapRoomHallDirection(int32(a1))) }

//export sub_523960
func sub_523960(a1 C.int) C.int { return (C.int)(mapRoomOpposite(int32(a1))) }

//export sub_523970
func sub_523970(a1 C.int) C.int { return (C.int)(mapRoomHallSide(int32(a1))) }

//export sub_5239B0
func sub_5239B0(a1 C.int) C.int { return (C.int)(mapRoomHallOppositeSide(int32(a1))) }

//export sub_523A10
func sub_523A10(a1 C.int, a2 *C.float) C.int {
	return (C.int)(mapRoomTrimHall((*mapRoom)(mapRoomPointer(uint32(a1))), (*mapRoom)(unsafe.Pointer(a2))))
}

//export sub_523C30
func sub_523C30(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomHallExit((*mapRoom)(mapRoomPointer(uint32(a1))), (*types.Pointf)(mapRoomPointer(uint32(a2)))))
}

//export sub_523CB0
func sub_523CB0(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomHallEntry((*mapRoom)(mapRoomPointer(uint32(a1))), (*types.Pointf)(mapRoomPointer(uint32(a2)))))
}

//export sub_523D30
func sub_523D30(a1 *C.float, a2 *C.float) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomHallCenter((*mapRoom)(unsafe.Pointer(a1)), (*types.Pointf)(unsafe.Pointer(a2)))))
}

//export sub_523E30
func sub_523E30(a1 C.int, a2 C.int, a3 C.int) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomNewHall(int32(a1), int32(a2), int32(a3))))
}

//export nox_xxx_mapGenMakeHall_523EC0
func nox_xxx_mapGenMakeHall_523EC0(a1 C.int, a2 C.int, a3 C.int) *C.float {
	return (*C.float)(unsafe.Pointer(mapRoomPrepareHall(mapRoomPointer(uint32(a1)), int32(a2), int32(a3))))
}

//export sub_524070
func sub_524070(a1 C.int, a2 C.int) C.int {
	return (C.int)(uintptr(unsafe.Pointer(mapRoomAssignDecoration(mapRoomPointer(uint32(a1)), (*mapRoom)(mapRoomPointer(uint32(a2)))))))
}

//export nox_xxx_mapGenMakeRooms_524310
func nox_xxx_mapGenMakeRooms_524310(a1 C.int) C.int {
	return (C.int)(mapRoomAssignRequiredDecorations(mapRoomPointer(uint32(a1))))
}

//export nox_xxx_mapGenSetRngSeed_526AB0
func nox_xxx_mapGenSetRngSeed_526AB0(a1 C.uint) { mapRoomSeed(uint32(a1)) }

//export nox_xxx_mapGenRandFunc_526AC0
func nox_xxx_mapGenRandFunc_526AC0(a1 C.int, a2 C.int) C.int {
	return (C.int)(mapRoomRandomInt(int32(a1), int32(a2)))
}

//export sub_526BC0
func sub_526BC0(a1 C.float, a2 C.float) C.double {
	return (C.double)(mapRoomRandomFloat(float32(a1), float32(a2)))
}
