package legacy

/*
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME4_1.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_wallMath_427F30
func nox_xxx_wallMath_427F30(p *C.int2, r *C.int) C.int {
	return C.int(geometryWallPoint((*[2]int32)(unsafe.Pointer(p)), (*[8]int32)(unsafe.Pointer(r))))
}

//export sub_428170
func sub_428170(p unsafe.Pointer, r *C.int4) C.int {
	return C.int(geometryWallBounds((*[8]uint32)(p), (*[4]int32)(unsafe.Pointer(r))))
}

//export nox_xxx_pointInRect_4281F0
func nox_xxx_pointInRect_4281F0(p *C.int2, r *C.int4) C.int {
	return C.int(geometryRectInt((*[2]int32)(unsafe.Pointer(p)), (*[4]int32)(unsafe.Pointer(r))))
}

//export sub_4D3E30
func sub_4D3E30(p, out *C.float2) C.int {
	return C.int(geometryMapCoordinates((*types.Pointf)(unsafe.Pointer(p)), (*types.Pointf)(unsafe.Pointer(out))))
}

//export nox_xxx_mathDirection4ToAngle_509E90
func nox_xxx_mathDirection4ToAngle_509E90(index C.int) C.int {
	return C.int(geometryDirection4Angle(int32(index)))
}

//export nox_xxx_math_509EA0
func nox_xxx_math_509EA0(index C.int) C.int { return C.int(geometryDirection4Index(int32(index))) }

//export nox_xxx_math_509ED0
func nox_xxx_math_509ED0(p *C.float2) C.int {
	return C.int(geometryVectorAngle((*types.Pointf)(unsafe.Pointer(p))))
}

//export sub_54FFC0
func sub_54FFC0(grid *C.int2, u C.int) C.int {
	return C.int(geometryCircleWall((*[2]int32)(unsafe.Pointer(grid)), objectFromInt(u)))
}

//export sub_5504B0
func sub_5504B0(u C.int) { geometryBoxWalls(objectFromInt(u)) }

//export nox_xxx_collisionCheckCircleCircle_550D00
func nox_xxx_collisionCheckCircleCircle_550D00(a, b C.int) {
	geometryCircleCircle(objectFromInt(a), objectFromInt(b))
}

//export sub_550F80
func sub_550F80(a *C.float, b C.int) {
	geometryBoxBox((*server.Object)(unsafe.Pointer(a)), objectFromInt(b))
}

//export sub_551250
func sub_551250(a C.uint, b *C.float, mode C.int) {
	geometryGateBox(objectFromInt(C.int(a)), (*server.Object)(unsafe.Pointer(b)), int32(mode))
}
