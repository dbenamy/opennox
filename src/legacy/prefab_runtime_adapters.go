package legacy

/*
#include <stdint.h>
#include <stdio.h>
#include "defs.h"
*/
import "C"
import (
	"math"
)

func sub_502DF0() uint32 { return prefabClose() }

//export nox_xxx_tileAllocTileInCoordList_5040A0
func nox_xxx_tileAllocTileInCoordList_5040A0(a0 C.int, a1 C.int, a2 C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabTileNew(int32(uint32(a0)), int32(uint32(a1)), math.Float32bits(float32(a2)))))
}

//export sub_504290
func sub_504290(a0 C.char, a1 C.char) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWallNew(byte(uint32(a0)), byte(uint32(a1)))))
}

func sub_5044B0(a0 C.int, a1 C.float, a2 C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWaypointNew(uint32(a0), math.Float32frombits(math.Float32bits(float32(a1))), math.Float32frombits(math.Float32bits(float32(a2))))))
}

func nox_xxx_unitAddToList_5048A0(a0 C.int) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabObjectNew(uint32(a0))))
}

func sub_51D0E0() { prefabResetWaypoint() }
