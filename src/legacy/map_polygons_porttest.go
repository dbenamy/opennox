//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "GAME1_2.h"
extern uint32_t nox_xxx_polygonNextIdx_587000_60352;
extern uint32_t nox_xxx_polygonNextAngle_587000_60356;
extern uint32_t dword_5d4594_588068;
*/
import "C"
import (
	"math"
	"unsafe"
)

func PortTestMapPolygonWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"polygons": (*uint32)(unsafe.Pointer(&C.nox_xxx_polygonNextIdx_587000_60352)),
		"vertices": (*uint32)(unsafe.Pointer(&C.nox_xxx_polygonNextAngle_587000_60356)),
		"remap":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_588068)),
	}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestMapPolygonPointer(op string, p unsafe.Pointer, index int) unsafe.Pointer {
	switch op {
	case "reset":
		return unsafe.Pointer(C.sub_421B10())
	case "reset-records":
		return C.sub_421430()
	case "reset-vertices":
		return unsafe.Pointer(C.sub_421010())
	case "new":
		return unsafe.Pointer(C.sub_421230())
	case "first":
		return unsafe.Pointer(C.nox_xxx_polygonGetNext_4210A0())
	case "next":
		return unsafe.Pointer(C.sub_4210E0(C.int(uintptr(p))))
	case "get":
		return unsafe.Pointer(C.nox_xxx_polygonGetByIdx_4214A0(C.int(index)))
	case "vertex-first":
		return unsafe.Pointer(C.nox_xxx_polygon_420CA0())
	case "vertex-next":
		return unsafe.Pointer(C.nox_xxx_polygon_420CD0((*C.uint32_t)(p)))
	case "vertex-get":
		return unsafe.Pointer(C.nox_xxx_polygonGetAngle_421030(C.int(index)))
	case "remap-clear":
		return unsafe.Pointer(C.sub_420C70())
	default:
		panic(op)
	}
}
func PortTestMapPolygonScalar(op string, p unsafe.Pointer) int {
	switch op {
	case "free-vertex":
		return int(C.sub_420D10())
	case "free-polygon":
		return int(C.sub_421130())
	case "defaults-read":
		return int(C.sub_421160(C.int(uintptr(p))))
	case "defaults-write":
		return int(C.sub_4211D0(C.int(uintptr(p))))
	case "selected":
		return int(C.sub_421B40((*C.uint32_t)(p)))
	case "construct":
		C.sub_4214D0()
		return 0
	case "remap":
		C.sub_421040(C.int(uintptr(p)))
		return 0
	default:
		panic(op)
	}
}
func PortTestMapPolygonVertex(x, y float32, index, oldID int) unsafe.Pointer {
	return unsafe.Pointer(C.nox_xxx_polygonSetAngle_420D40(C.int(math.Float32bits(x)), C.int(math.Float32bits(y)), C.uint(index), C.int(oldID)))
}
func PortTestMapPolygonVertexSelect(x, y float32) unsafe.Pointer {
	return unsafe.Pointer(C.sub_420DA0(C.float(x), C.float(y)))
}
func PortTestMapPolygonVertexNearest(x, y, distance float32) unsafe.Pointer {
	return unsafe.Pointer(uintptr(uint32(C.sub_420E80(C.float(x), C.float(y), C.float(distance)))))
}
func PortTestMapPolygonRemapAdd(index, oldID int) unsafe.Pointer {
	return unsafe.Pointer(C.sub_420C40(C.int(index), C.int(oldID)))
}
func PortTestMapPolygonActor(op string, p unsafe.Pointer) {
	switch op {
	case "player":
		C.nox_xxx_questCheckSecretArea_421C70((*C.nox_object_t)(p))
	case "monster":
		C.nox_xxx_monsterPolygonEnter_421FF0((*C.nox_object_t)(p))
	default:
		panic(op)
	}
}

func PortTestMapPolygonContains(p unsafe.Pointer, point *[2]int32) int {
	return int(C.nox_xxx_polygon_421660((*C.int)(unsafe.Pointer(point)), C.int(uintptr(p))))
}
func PortTestMapPolygonIndex(point *[2]int32, cached int) int {
	return int(C.nox_xxx_polygonGetIdxA_421790((*C.int2)(unsafe.Pointer(point)), C.int(cached)))
}
func PortTestMapPolygonFind(op string, point *[2]int32, cached int, distance float32) unsafe.Pointer {
	switch op {
	case "point":
		return unsafe.Pointer(C.nox_xxx_polygonIsPlayerInPolygon_4217B0((*C.int2)(unsafe.Pointer(point)), C.int(cached)))
	case "script":
		return unsafe.Pointer(C.sub_421F10((*C.int)(unsafe.Pointer(point)), C.int(cached)))
	case "edge":
		return unsafe.Pointer(C.sub_421990((*C.int2)(unsafe.Pointer(point)), C.float(distance), C.int(cached)))
	default:
		panic(op)
	}
}
func PortTestMapPolygonEdge(p unsafe.Pointer, point *[2]int32, distance float32) int {
	return int(C.sub_421880(C.int(uintptr(unsafe.Pointer(point))), C.int(uintptr(p)), C.float(distance)))
}
func PortTestMapPolygonActorInit(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(uint32(C.sub_422140(C.int(uintptr(p))))))
}
func PortTestMapPolygonColor() { C.nox_xxx_polygonDrawColor_421B80() }
func PortTestMapPolygonSection(bypass int) int {
	return int(C.nox_server_mapRWPolygons_428CD0(C.int(bypass)))
}
