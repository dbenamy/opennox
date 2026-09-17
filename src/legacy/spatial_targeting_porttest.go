//go:build porttest

package legacy

/*
#include "GAME5_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2491592;
extern unsigned int nox_player_netCode_85319C;
static int spatialPredict(nox_object_t*a,nox_object_t*b,float speed,float2*out) {return nox_xxx_projAddVelocitySmth_533080((int)a,(int)b,speed,(int)out);}
static int spatialProbe(nox_object_t*a,float2*next,float2*prev) {return sub_54E810((int)a,next,(int)prev);}
static int spatialCandidate(nox_object_t*a,nox_object_t*b,float2*next,float2*prev) {int ctx[4]={(int)a,0,(int)prev,(int)next};sub_54E850((int)b,(int)ctx);return ctx[1];}
static void spatialCursorCandidate(nox_object_t*u,float2*p) {nox_xxx_playerCursorScanFn_54AFB0((int)u,(float*)p);}

int sub_57CDB0(int2*, float*, float2*);
char sub_57F1D0(float2*);
int sub_57F2A0(float2*, int, int);
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Invoke the production C helpers; ownership and expected geometry live in the tests.
func PortTestSpatialQuadrant(p *types.Pointf) int8 {
	return int8(C.sub_57F1D0((*C.float2)(unsafe.Pointer(p))))
}
func PortTestSpatialTriangle(p *types.Pointf, grid [2]int32) int32 {
	return int32(C.sub_57F2A0((*C.float2)(unsafe.Pointer(p)), C.int(grid[0]), C.int(grid[1])))
}
func PortTestSpatialNormal(grid *[2]int32, ray *[4]float32, normal *types.Pointf) int32 {
	return int32(C.sub_57CDB0((*C.int2)(unsafe.Pointer(grid)), (*C.float)(unsafe.Pointer(ray)), (*C.float2)(unsafe.Pointer(normal))))
}

func PortTestSpatialEligible(a, b *server.Object, teamFilter bool) int32 {
	if teamFilter {
		return int32(C.sub_54E6F0(C.int(uintptr(a.CObj())), C.int(uintptr(b.CObj()))))
	}
	return int32(C.sub_54E730(C.int(uintptr(a.CObj())), C.int(uintptr(b.CObj()))))
}
func PortTestSpatialPredict(a, b *server.Object, speed float32, out *types.Pointf) uint32 {
	return uint32(C.spatialPredict((*C.nox_object_t)(a.CObj()), (*C.nox_object_t)(b.CObj()), C.float(speed), (*C.float2)(unsafe.Pointer(out))))
}
func PortTestSpatialProbe(a *server.Object, next, prev *types.Pointf) uint32 {
	return uint32(C.spatialProbe((*C.nox_object_t)(a.CObj()), (*C.float2)(unsafe.Pointer(next)), (*C.float2)(unsafe.Pointer(prev))))
}
func PortTestSpatialCandidate(a, b *server.Object, next, prev *types.Pointf) uint32 {
	return uint32(C.spatialCandidate((*C.nox_object_t)(a.CObj()), (*C.nox_object_t)(b.CObj()), (*C.float2)(unsafe.Pointer(next)), (*C.float2)(unsafe.Pointer(prev))))
}
func PortTestSpatialRay(ray *[4]float32) int32 {
	return int32(C.nox_xxx_traceRay_5374B0((*C.float4)(unsafe.Pointer(ray))))
}
func PortTestSpatialCursorCandidate(u *server.Object, p *types.Pointf) {
	C.spatialCursorCandidate((*C.nox_object_t)(u.CObj()), (*C.float2)(unsafe.Pointer(p)))
}
func PortTestSpatialCursor(u *server.Object) uint32 {
	return uint32(uintptr(unsafe.Pointer(C.nox_xxx_findObjectAtCursor_54AF40((*C.nox_object_t)(u.CObj())))))
}
func PortTestSpatialGlobals() (map[string]*uint32, func()) {
	out := map[string]*uint32{"owner": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2491592)), "local": (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C))}
	for name, off := range map[string]uintptr{"chosen": 2491596, "score": 2491600, "polyp": 2491604} {
		out[name] = memmap.PtrUint32(0x5D4594, off)
	}
	old := map[string]uint32{}
	for n, p := range out {
		old[n] = *p
		*p = 0
	}
	return out, func() {
		for n, p := range out {
			*p = old[n]
		}
	}
}
