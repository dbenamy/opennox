package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import "unsafe"

//export sub_420DA0
func sub_420DA0(x, y C.float) *C.uint {
	return (*C.uint)(unsafe.Pointer(mapPolygonSelect(float32(x), float32(y))))
}

//export sub_4211D0
func sub_4211D0(p C.int) C.int {
	return C.int(mapPolygonDefaultsWrite((*mapPolygon)(unsafe.Pointer(uintptr(uint32(p))))))
}

//export sub_4214D0
func sub_4214D0() { mapPolygonConstruct() }

//export nox_xxx_polygonIsPlayerInPolygon_4217B0
func nox_xxx_polygonIsPlayerInPolygon_4217B0(point *C.int2, cached C.int) *C.nox_player_polygon_check_data {
	return (*C.nox_player_polygon_check_data)(unsafe.Pointer(mapPolygonFind((*[2]int32)(unsafe.Pointer(point)), uint32(cached), false)))
}

//export sub_421990
func sub_421990(point *C.int2, distance C.float, cached C.int) *C.int {
	return (*C.int)(unsafe.Pointer(mapPolygonFindEdge((*[2]int32)(unsafe.Pointer(point)), uint32(cached), float32(distance))))
}

//export sub_421B10
func sub_421B10() *C.uint32_t { return (*C.uint32_t)(mapPolygonReset()) }

//export sub_422140
func sub_422140(p C.int) C.int {
	return C.int(uintptr(mapPolygonActorInit(unsafe.Pointer(uintptr(uint32(p))))))
}
