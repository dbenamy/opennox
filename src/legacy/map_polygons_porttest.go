//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestMapPolygonWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{"polygons": &mapPolygonNext, "vertices": &mapPolygonVertexNext, "remap": (*uint32)(unsafe.Pointer(&mapPolygonRemapHead))}
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
		return mapPolygonReset()
	case "reset-records":
		return mapPolygonResetRecords()
	case "reset-vertices":
		return mapPolygonResetVertices()
	case "new":
		return unsafe.Pointer(mapPolygonNew())
	case "first":
		return unsafe.Pointer(mapPolygonFirst())
	case "next":
		return unsafe.Pointer(mapPolygonAfter((*mapPolygon)(p).ID))
	case "get":
		return unsafe.Pointer(mapPolygonGet(uint32(index)))
	case "vertex-first":
		return unsafe.Pointer(mapPolygonVertexFirst())
	case "vertex-next":
		return unsafe.Pointer(mapPolygonVertexAfter((*mapPolygonVertex)(p).ID))
	case "vertex-get":
		return unsafe.Pointer(mapPolygonVertexAt(uint32(index)))
	case "remap-clear":
		return mapPolygonRemapClear()
	default:
		panic(op)
	}
}
func PortTestMapPolygonScalar(op string, p unsafe.Pointer) int {
	switch op {
	case "free-vertex":
		return int(mapPolygonFreeVertex())
	case "free-polygon":
		return int(mapPolygonFreeIndex())
	case "defaults-read":
		return mapPolygonDefaultsRead((*mapPolygon)(p))
	case "defaults-write":
		return mapPolygonDefaultsWrite((*mapPolygon)(p))
	case "selected":
		return bool2int(mapPolygonSelected((*mapPolygonVertex)(p)))
	case "construct":
		mapPolygonConstruct()
		return 0
	case "remap":
		mapPolygonRemapIDs((*mapPolygon)(p))
		return 0
	default:
		panic(op)
	}
}
func PortTestMapPolygonVertex(x, y float32, index, oldID int) unsafe.Pointer {
	return unsafe.Pointer(mapPolygonVertexSet(x, y, uint32(index), uint32(oldID)))
}
func PortTestMapPolygonVertexSelect(x, y float32) unsafe.Pointer {
	return unsafe.Pointer(mapPolygonSelect(x, y))
}
func PortTestMapPolygonVertexNearest(x, y, distance float32) unsafe.Pointer {
	return unsafe.Pointer(mapPolygonNearest(x, y, distance))
}
func PortTestMapPolygonRemapAdd(index, oldID int) unsafe.Pointer {
	return unsafe.Pointer(mapPolygonRemapAdd(uint32(index), uint32(oldID)))
}
func PortTestMapPolygonActor(op string, p unsafe.Pointer) {
	switch op {
	case "player":
		mapPolygonPlayer((*server.Object)(p))
	case "monster":
		mapPolygonMonster((*server.Object)(p))
	default:
		panic(op)
	}
}
func PortTestMapPolygonContains(p unsafe.Pointer, point *[2]int32) int {
	return bool2int(mapPolygonContains((*mapPolygon)(p), point))
}
func PortTestMapPolygonIndex(point *[2]int32, cached int) int {
	return int(mapPolygonIndex(point, uint32(cached)))
}
func PortTestMapPolygonFind(op string, point *[2]int32, cached int, distance float32) unsafe.Pointer {
	switch op {
	case "point":
		return unsafe.Pointer(mapPolygonFind(point, uint32(cached), false))
	case "script":
		return unsafe.Pointer(mapPolygonFind(point, uint32(cached), true))
	case "edge":
		return unsafe.Pointer(mapPolygonFindEdge(point, uint32(cached), distance))
	default:
		panic(op)
	}
}
func PortTestMapPolygonEdge(p unsafe.Pointer, point *[2]int32, distance float32) int {
	return bool2int(mapPolygonEdge((*mapPolygon)(p), point, distance))
}
func PortTestMapPolygonActorInit(p unsafe.Pointer) unsafe.Pointer { return mapPolygonActorInit(p) }
func PortTestMapPolygonColor()                                    { mapPolygonColor() }
func PortTestMapPolygonSection(bypass int) int                    { return mapPolygonSection(bypass) }
