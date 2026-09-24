package legacy

/*
#include <stdint.h>
#include <stdio.h>
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

//export nox_server_scriptExecuteFnForEachGroupObj_502670
func nox_server_scriptExecuteFnForEachGroupObj_502670(group *C.uchar, expected C.int, callback unsafe.Pointer, data C.int) {
	prefabGroupEach((*server.MapGroup)(unsafe.Pointer(group)), int32(expected), callback, uint32(data))
}
func sub_5029A0(a0 *C.char) C.int { return C.int(prefabFindName(mapRoomRaw(unsafe.Pointer(a0)))) }
func sub_5029F0(a0 C.int) C.int   { return C.int(prefabMetadataAt(int32(uint32(a0)))) }
func sub_502A20() C.int           { return C.int(*prefabGlobal(prefabCount)) }

//export sub_502A50
func sub_502A50(a0 *C.char) C.int { return C.int(prefabSetPath(mapRoomRaw(unsafe.Pointer(a0)), false)) }

//export sub_502AB0
func sub_502AB0(a0 *C.char) C.int { return C.int(prefabSetPath(mapRoomRaw(unsafe.Pointer(a0)), true)) }

//export sub_502B10
func sub_502B10() C.int                           { return C.int(prefabLibrary()) }
func sub_502D70(a0 C.int) C.int                   { return C.int(populationLoadPrefab(int32(a0))) }
func sub_502DF0() uint32                          { return prefabClose() }
func sub_502E70(a0 C.int) C.double                { return C.double(prefabDimension(int32(uint32(a0)), 64)) }
func sub_502EA0(a0 C.int) C.double                { return C.double(prefabDimension(int32(uint32(a0)), 68)) }
func nox_xxx_mapgenSaveMap_503830(a0 C.int) C.int { return C.int(prefabLoad(int32(uint32(a0)))) }
func sub_503B30(a0 *C.float2) C.int {
	return C.int(prefabInstantiate(prefabPoint(mapRoomRaw(unsafe.Pointer(a0)))))
}
func sub_503EC0(a0 C.int, a1 *C.float) C.int {
	return C.int(prefabRelativePosition(uint32(a0), prefabPoint(mapRoomRaw(unsafe.Pointer(a1)))))
}

//export nox_xxx_tileAllocTileInCoordList_5040A0
func nox_xxx_tileAllocTileInCoordList_5040A0(a0 C.int, a1 C.int, a2 C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabTileNew(int32(uint32(a0)), int32(uint32(a1)), math.Float32bits(float32(a2)))))
}
func nox_xxx_tileInit_504150(a0 C.int, a1 C.int) C.int {
	return C.int(prefabPlaceTiles(int32(uint32(a0)), int32(uint32(a1))))
}

//export sub_504290
func sub_504290(a0 C.char, a1 C.char) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWallNew(byte(uint32(a0)), byte(uint32(a1)))))
}

//export nox_xxx_cliWallGet_5042F0
func nox_xxx_cliWallGet_5042F0(a0 C.int, a1 C.int) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWallFind(int32(uint32(a0)), int32(uint32(a1)))))
}
func sub_504330(a0 C.int, a1 C.int) C.int {
	return C.int(prefabPlaceWalls(int32(uint32(a0)), int32(uint32(a1))))
}
func sub_5044B0(a0 C.int, a1 C.float, a2 C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabWaypointNew(uint32(a0), math.Float32frombits(math.Float32bits(float32(a1))), math.Float32frombits(math.Float32bits(float32(a2))))))
}
func sub_504560(a0 C.int, a1 C.int) C.int {
	return C.int(prefabPlaceWaypoints(int32(uint32(a0)), int32(uint32(a1))))
}
func nox_xxx_unitAddToList_5048A0(a0 C.int) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabObjectNew(uint32(a0))))
}
func sub_504910(a0 C.int, a1 C.int) C.int {
	return C.int(prefabPlaceObjects(int32(uint32(a0)), int32(uint32(a1))))
}
func sub_504980() C.int         { return C.int(prefabObjectHead()) }
func sub_5049C0(a0 C.int) C.int { return C.int(prefabObjectNext(uint32(a0))) }

//export sub_5049D0
func sub_5049D0() unsafe.Pointer { return mapRoomPointer(*prefabGlobal(prefabObjects)) }

//export sub_5049E0
func sub_5049E0(a0 C.int) C.int { return C.int(prefabNodeNext(uint32(a0))) }
func sub_504A10(a0 C.int) C.int { return C.int(prefabObjectRemove(uint32(a0))) }

//export sub_51D0E0
func sub_51D0E0()                { prefabResetWaypoint() }
func sub_51D0F0(a0 C.char) C.int { return C.int(prefabSetWaypointKind(byte(uint32(a0)))) }
func sub_51D120(a0 *C.float) *C.uint32_t {
	return (*C.uint32_t)(mapRoomPointer(prefabCreateWaypoint(prefabPoint(mapRoomRaw(unsafe.Pointer(a0))))))
}
func sub_51D1A0(a0 *C.float2) *C.float {
	return (*C.float)(mapRoomPointer(prefabFindWaypoint(prefabPoint(mapRoomRaw(unsafe.Pointer(a0))))))
}
func sub_51D3F0(a0 *C.float2, a1 *C.float2) uint32 {
	return prefabConnectWaypoint(prefabPoint(mapRoomRaw(unsafe.Pointer(a0))), prefabPoint(mapRoomRaw(unsafe.Pointer(a1))))
}
