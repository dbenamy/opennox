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

//export nox_xxx_math_509EA0
func nox_xxx_math_509EA0(index C.int) C.int { return C.int(geometryDirection4Index(int32(index))) }
