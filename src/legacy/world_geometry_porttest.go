//go:build porttest

package legacy

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"unsafe"
)

// Pointer offsets refer to one guarded C allocation, allowing explicit aliasing.
// Negative offsets supply nil only for functions whose interface admits it.
type PortTestWorldGeometrySpec struct {
	Op      string
	Words   []uint32
	Offsets [4]int
	Ints    [3]int32
	Floats  [3]uint32
}
type PortTestWorldGeometryResult struct {
	Return   int32
	Words    []uint32
	GuardsOK bool
}

func PortTestWorldGeometry(spec PortTestWorldGeometrySpec) PortTestWorldGeometryResult {
	all, data, guards, free := portTestCollisionWords(spec.Words)
	defer free()
	var ptr [4]unsafe.Pointer
	for i, off := range spec.Offsets {
		if off >= 0 {
			if off >= len(data) {
				panic("geometry offset")
			}
			ptr[i] = unsafe.Pointer(&data[off])
		}
	}
	f := func(i int) float32 { return float32(math.Float32frombits(spec.Floats[i])) }
	var rv int32
	switch spec.Op {
	case "segments":
		rv = int32(geometrySegments((*[4]int32)(unsafe.Pointer(ptr[0])), (*[4]int32)(unsafe.Pointer(ptr[1]))))
	case "edge-project":
		rv = int32(geometryProjectEdge((*[2]int32)(ptr[0]), (*[4]int32)(unsafe.Pointer(ptr[1])), float32(f(0))))
	case "wall-point":
		rv = int32(geometryWallPoint((*[2]int32)(unsafe.Pointer(ptr[0])), (*[8]int32)(unsafe.Pointer(ptr[1]))))
	case "wall-bounds":
		rv = int32(geometryWallBounds((*[8]uint32)(unsafe.Pointer(ptr[0])), (*[4]int32)(unsafe.Pointer(ptr[1]))))
	case "rect-int":
		rv = int32(geometryRectInt((*[2]int32)(unsafe.Pointer(ptr[0])), (*[4]int32)(unsafe.Pointer(ptr[1]))))
	case "rect-float":
		rv = int32(geometryRectFloat((*types.Pointf)(unsafe.Pointer(ptr[0])), (*[4]float32)(unsafe.Pointer(ptr[1]))))
	case "box-calc":
		geometryShapeBox((*server.Shape)(unsafe.Pointer(ptr[0])))
	case "map-coordinates":
		rv = int32(geometryMapCoordinates((*types.Pointf)(unsafe.Pointer(ptr[0])), (*types.Pointf)(unsafe.Pointer(ptr[1]))))
	case "direction-angle":
		rv = int32(geometryDirectionAngle((*[2]uint32)(unsafe.Pointer(ptr[0]))))
	case "indexed-direction":
		rv = int32(geometryIndexedDirection(int32(spec.Ints[0]), (*[2]int32)(unsafe.Pointer(ptr[0]))))
	case "direction4-angle":
		rv = int32(geometryDirection4Angle(int32(spec.Ints[0])))
	case "direction4-index":
		rv = int32(geometryDirection4Index(int32(spec.Ints[0])))
	case "vector-angle":
		rv = int32(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(ptr[0]))))
	case "normalize":
		geometryNormalize((*types.Pointf)(unsafe.Pointer(ptr[0])))
	case "project-positive":
		rv = int32(geometryProjectPositive((*types.Pointf)(ptr[0]), float32(f(0)), float32(f(1)), int32(spec.Ints[0]), int32(spec.Ints[1]), (*types.Pointf)(ptr[1]), (*types.Pointf)(ptr[2])))
	case "project-negative":
		rv = int32(geometryProjectNegative((*types.Pointf)(unsafe.Pointer(ptr[0])), float32(f(0)), float32(f(1)), int32(spec.Ints[0]), int32(spec.Ints[1]), (*types.Pointf)(unsafe.Pointer(ptr[1])), (*types.Pointf)(unsafe.Pointer(ptr[2]))))
	case "quadrant":
		rv = int32(int32(geometryQuadrant((*types.Pointf)(unsafe.Pointer(ptr[0])), (*types.Pointf)(unsafe.Pointer(ptr[1])))))
	case "rectangle-crossings":
		rv = int32(geometryRectCrossings((*[4]float32)(unsafe.Pointer(ptr[0])), (*[4]float32)(unsafe.Pointer(ptr[1])), (*types.Pointf)(unsafe.Pointer(ptr[2])), int32(spec.Ints[0]), int32(spec.Ints[1])))
	case "horizontal-crossing":
		rv = int32(geometryCrossHorizontal((*[4]float32)(unsafe.Pointer(ptr[0])), float32(f(0)), float32(f(1)), float32(f(2)), (*types.Pointf)(unsafe.Pointer(ptr[1])), int32(spec.Ints[0])))
	case "vertical-crossing":
		rv = int32(geometryCrossVertical((*[4]float32)(unsafe.Pointer(ptr[0])), float32(f(0)), float32(f(1)), float32(f(2)), (*types.Pointf)(unsafe.Pointer(ptr[1])), int32(spec.Ints[0])))
	case "clipped-center":
		rv = int32(geometryClipCenter((*[4]float32)(unsafe.Pointer(ptr[0])), (*[4]float32)(unsafe.Pointer(ptr[1])), (*[4]float32)(unsafe.Pointer(ptr[2])), (*types.Pointf)(unsafe.Pointer(ptr[3]))))
	case "rect-float-alt":
		rv = int32(geometryRectFloat((*types.Pointf)(unsafe.Pointer(ptr[0])), (*[4]float32)(unsafe.Pointer(ptr[1]))))
	default:
		panic(spec.Op)
	}
	ok := true
	for i := 0; i < portTestCollisionGuardWords; i++ {
		ok = ok && all[i] == guards[i] && all[len(all)-portTestCollisionGuardWords+i] == guards[portTestCollisionGuardWords+i]
	}
	return PortTestWorldGeometryResult{int32(rv), append([]uint32(nil), data...), ok}
}

func PortTestWorldGeometryPhysics(op string, u, v *server.Object, words *[16]uint32, flags int32, distance float32) int {
	p := unsafe.Pointer(&words[0])
	p1 := unsafe.Pointer(&words[2])
	p2 := unsafe.Pointer(&words[4])
	p3 := unsafe.Pointer(&words[8])
	switch op {
	case "circle-wall":
		return int(geometryCircleWall((*[2]int32)(p), u))
	case "point-wall":
		return int(geometryPointWall(u, (*types.Pointf)(p)))
	case "game-ball":
		return int(geometryGameBall(u))
	case "box-walls":
		geometryBoxWalls(u)
	case "box-wall":
		return int(geometryBoxWall((*[2]int32)(p), u))
	case "horizontal-wall":
		return int(geometryWallHorizontal(u, (*types.Pointf)(p), (*types.Pointf)(p1), (*[4]float32)(p2), (*types.Pointf)(p3), distance))
	case "vertical-wall":
		return int(geometryWallVertical(u, (*types.Pointf)(p), (*types.Pointf)(p1), (*[4]float32)(p2), (*types.Pointf)(p3), distance))
	case "circle-circle":
		geometryCircleCircle(u, v)
	case "box-box":
		geometryBoxBox(u, v)
	case "gate-box":
		geometryGateBox(u, v, flags)
	default:
		panic(op)
	}
	return 0
}

func PortTestWorldGeometryGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"directionThreshold": (*uint32)(unsafe.Pointer(&geometryDirectionThreshold)),
		"objectForce":        (*uint32)(unsafe.Pointer(&collisionObjectForce)),
		"wallForce":          (*uint32)(unsafe.Pointer(&collisionWallForce)),
		"gameBall":           memmap.PtrUint32(0x5D4594, 2491788),
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

// Isolate the real Hit class and both indices, without replacing collision logic.
func PortTestWorldGeometryHitOwner() (reset func(), snapshot func(map[unsafe.Pointer]uint32) [][5]uint32, restore func()) {
	oldClass, oldHead := collisionHitClass, collisionHitHead
	buckets := collisionBuckets[:]
	oldBuckets := append([]uint32(nil), buckets...)
	collisionHitClass = nil
	collisionHitHead = 0
	clear(buckets)
	reset = func() { collisionResetHits() }
	reset()
	snapshot = func(ids map[unsafe.Pointer]uint32) [][5]uint32 {
		normalize := func(raw uint32) uint32 {
			if raw <= 6 {
				return raw
			}
			id, ok := ids[unsafe.Pointer(uintptr(raw))]
			if !ok {
				panic("collision outside fixture owner")
			}
			return id
		}
		var rows [][5]uint32
		for p := uint32(collisionHitHead); p != 0; {
			if len(rows) >= 1024 {
				panic("collision list cycle")
			}
			row := unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(p))), 7)
			rows = append(rows, [5]uint32{normalize(row[2]), normalize(row[3]), row[4], row[5], row[6]})
			p = row[1]
		}
		return rows
	}
	restore = func() {
		alloc.AsClass(collisionHitClass).Free()
		collisionHitClass = oldClass
		collisionHitHead = oldHead
		copy(buckets, oldBuckets)
	}
	return
}
