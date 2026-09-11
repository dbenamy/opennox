package legacy

/*
#include "GAME5_2.h"
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/opennox/libs/types"
)

// quietCollisionNaN preserves the payload/sign while modeling an x87 float load
// followed by a float store. Ordinary numbers and infinities keep their bits.
func quietCollisionNaN(bits uint32) uint32 {
	if bits&0x7fffffff > 0x7f800000 {
		bits |= 0x00400000
	}
	return bits
}

func collisionReflect(normal, velocity *types.Pointf) {
	x, y := math.Float32bits(velocity.X), math.Float32bits(velocity.Y)
	// The original x87 product does not underflow/overflow at float32 limits.
	if float64(normal.Y)*float64(normal.X) <= 0 {
		velocity.X = math.Float32frombits(y)
		velocity.Y = math.Float32frombits(quietCollisionNaN(x))
	} else {
		velocity.X = math.Float32frombits(quietCollisionNaN(y) ^ 0x80000000)
		velocity.Y = math.Float32frombits(quietCollisionNaN(x) ^ 0x80000000)
	}
}

func collisionContains(pos *types.Pointf, shape *[11]float32, point *types.Pointf) bool {
	// Keep the original operation order and PC53 intermediates. The C source's
	// float locals stayed in x87 registers; rounding them to float32 changes edges.
	const scale = 0.70709997
	x, y := float64(pos.X), float64(pos.Y)
	px, py := float64(point.X), float64(point.Y)
	v4 := float64(shape[5]) + x
	v5 := float64(shape[6]) + y
	if !((v5-v4+px-py)*scale < 0 &&
		(float64(shape[8])+y-(float64(shape[7])+x)+px-py)*scale > 0) {
		return false
	}
	v6 := float64(shape[9]) + x
	v7 := float64(shape[10]) + y
	return (v7+v6-px-py)*scale > 0 && (v5+v4-px-py)*scale < 0
}

//export nox_xxx_collideReflect_57B810
func nox_xxx_collideReflect_57B810(normal *C.float, velocity C.int) C.int {
	collisionReflect((*types.Pointf)(unsafe.Pointer(normal)), (*types.Pointf)(unsafe.Pointer(uintptr(uint32(velocity)))))
	return velocity
}

//export nox_xxx_map_57B850
func nox_xxx_map_57B850(pos *C.float2, shape *C.float, point *C.float2) C.int {
	return C.int(bool2int(collisionContains((*types.Pointf)(unsafe.Pointer(pos)), (*[11]float32)(unsafe.Pointer(shape)), (*types.Pointf)(unsafe.Pointer(point)))))
}
