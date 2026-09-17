//go:build porttest

package legacy

/*
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
int nox_xxx_collidePentagram_4EAB20(int a1);
static uint32_t coreContacts[1024];
static unsigned coreContactCount;
// This is an observation callback, not a collision implementation. Alignment
// makes the legacy activation helper's callback-address low byte deterministic.
static void __attribute__((aligned(256))) coreContact(void* a,void* b,float2* p) {
 if(coreContactCount>=256) return;
 uint32_t* r=coreContacts+4*coreContactCount++;
 r[0]=(uintptr_t)a;r[1]=(uintptr_t)b;
 memcpy(r+2,p,8);
}
static void* coreContactPtr(void){return coreContact;}
static uint32_t* coreContactData(void){return coreContacts;}
static unsigned* coreContactN(void){return &coreContactCount;}
static uint32_t coreRadialResult[3];
static void coreRadialObserve(int obj, uint32_t code) {
 coreRadialResult[0]++; coreRadialResult[1]=obj; coreRadialResult[2]=code;
}
static void* coreRadialCallback(void) {return coreRadialObserve;}
static uint32_t* coreRadialData(void) {return coreRadialResult;}
static void* corePentagramPtr(void){return nox_xxx_collidePentagram_4EAB20;}
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestCollisionCore(op string, a, b *server.Object, p *types.Pointf, mode int32) int32 {

	switch op {
	case "angle-queue":
		worldAngleQueue(a.UpdateData)
	case "shape":
		geometryShapeBox(&a.Shape)
	case "contains":
		return collisionObjectContains(a, p)
	case "eligible":
		return collisionEligible(a, b)
	case "retained":
		return collisionRetained(a, b)
	case "pair":
		collisionPair(a, b)
	case "scan":
		collisionScan(a)
	case "circle-walls":
		collisionCircleWalls(a)
	case "circle-box":
		collisionCircleBox(a, b, mode)
	case "gate-circle":
		collisionGateCircle(a, b, mode)
	case "dispatch":
		collisionDispatch()
	case "angles":
		collisionDrainAngles()
	case "elevator":
		collisionElevator(a, b, mode)
	case "shaft":
		collisionShaft(a, b)
	case "types":
		collisionInitTypes()
	case "active":
		return int32(a.Field116 & 1)
	case "activate":
		return int32(collisionActivate(a))
	case "remove":
		collisionRemoveActive(a)
	case "head":
		return int32(collisionActiveHead)
	case "next":
		return int32(collisionNextActive(a))
	case "pop":
		return int32(collisionObjectAddress(collisionPopActive()))
	default:
		panic(op)
	}
	return 0
}
func PortTestCollisionCoreDistance(p *types.Pointf, radius float32, box *server.Object, out *types.Pointf) float64 {
	return collisionBoxDistance(p, radius, box, out)
}
func PortTestCollisionCoreWallOpen(grid *[2]int32, u *server.Object) {
	collisionWallOpen(grid, u)
}
func PortTestCollisionCoreAddHit(a, b *server.Object, sentinel uint32, normal *types.Pointf) {
	target := C.uint(sentinel)
	if b != nil {
		target = C.uint(uintptr(b.CObj()))
	}
	collisionAddHit(a, uint32(target), normal)
}
func PortTestCollisionCoreObserver() (unsafe.Pointer, func(), func(map[unsafe.Pointer]uint32) [][4]uint32, func()) {
	data := unsafe.Slice((*uint32)(unsafe.Pointer(C.coreContactData())), 1024)
	count := (*uint32)(unsafe.Pointer(C.coreContactN()))
	old := append([]uint32(nil), data...)
	oldN := *count
	reset := func() { clear(data); *count = 0 }
	reset()
	snapshot := func(ids map[unsafe.Pointer]uint32) [][4]uint32 {
		out := make([][4]uint32, *count)
		for i := range out {
			copy(out[i][:], data[4*i:4*i+4])
			for j := 0; j < 2; j++ {
				if out[i][j] != 0 {
					id, ok := ids[unsafe.Pointer(uintptr(out[i][j]))]
					if !ok {
						panic("contact outside fixture")
					}
					out[i][j] = id
				}
			}
		}
		return out
	}
	return C.coreContactPtr(), reset, snapshot, func() { copy(data, old); *count = oldN }
}
func PortTestCollisionCorePentagram() unsafe.Pointer { return C.corePentagramPtr() }
func PortTestCollisionCoreGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"trigger": &collisionTrigger, "powder": &collisionPowder, "hand": &collisionHand, "small": &collisionSmallFist, "medium": &collisionMediumFist, "large": &collisionLargeFist, "meteor": &collisionMeteor, "ready": &collisionTypesReady}
	old := map[string]uint32{}
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

// Inspect every real bucket link and validate its bucket and object identities.
func PortTestCollisionCoreBuckets(ids map[unsafe.Pointer]uint32) map[uint32][][2]uint32 {
	out := map[uint32][][2]uint32{}
	for i := uint32(0); i < 256; i++ {
		p := collisionBuckets[i]
		for n := 0; p != 0; n++ {
			if n >= 1024 {
				panic("collision bucket cycle")
			}
			r := (*[7]uint32)(unsafe.Pointer(uintptr(p)))
			if r[6] != i {
				panic("wrong collision bucket")
			}
			pair := [2]uint32{r[2], r[3]}
			for j := range pair {
				if pair[j] > 6 {
					id, ok := ids[unsafe.Pointer(uintptr(pair[j]))]
					if !ok {
						panic("bucket outside fixture")
					}
					pair[j] = id
				}
			}
			out[i] = append(out[i], pair)
			p = r[0]
		}
	}
	return out
}

// Exercise the retained production C caller as well as its distance helper.
func PortTestCollisionCoreRadial(u *server.Object, p *types.Pointf, radius float32, code uint32) [3]uint32 {
	data := unsafe.Slice((*uint32)(unsafe.Pointer(C.coreRadialData())), 3)
	clear(data)
	motionRadialCandidate(u, p, radius, C.coreRadialCallback(), code)
	out := [3]uint32{data[0], data[1], data[2]}
	clear(data)
	return out
}
