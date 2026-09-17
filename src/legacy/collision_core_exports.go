package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"unsafe"
)

//export nox_xxx_allocHitArray_5486D0
func nox_xxx_allocHitArray_5486D0() { collisionResetHits() }

//export nox_xxx_collSysAddCollision_548630
func nox_xxx_collSysAddCollision_548630(a C.int, b C.uint, p *C.float2) {
	collisionAddHit(objectFromInt(a), uint32(b), (*types.Pointf)(unsafe.Pointer(p)))
}

//export sub_547DB0
func sub_547DB0(a C.int, p *C.float2) C.int {
	return C.int(collisionObjectContains(objectFromInt(a), (*types.Pointf)(unsafe.Pointer(p))))
}

//export sub_5481C0
func sub_5481C0(a C.int) { collisionScan(objectFromInt(a)) }

//export sub_54A990
func sub_54A990(p *C.float2, r C.float, a C.int, n *C.float2) C.double {
	return C.double(collisionBoxDistance((*types.Pointf)(unsafe.Pointer(p)), float32(r), objectFromInt(a), (*types.Pointf)(unsafe.Pointer(n))))
}

//export nox_xxx_unitHasCollideOrUpdateFn_537610
func nox_xxx_unitHasCollideOrUpdateFn_537610(a *C.nox_object_t) C.char {
	return C.char(collisionActivate(asObjectS(a)))
}

//export sub_537580
func sub_537580(a C.int) C.int { return C.int(objectFromInt(a).Field116 & 1) }

//export sub_5375A0
func sub_5375A0(a C.int) { collisionRemoveActive(objectFromInt(a)) }

//export sub_537740
func sub_537740() C.int { return C.int(collisionActiveHead) }

//export sub_537750
func sub_537750(a C.int) C.int { return C.int(collisionNextActive(objectFromInt(a))) }
