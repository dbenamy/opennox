//go:build porttest

package legacy

/*
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
extern uint32_t dword_5d4594_2490508;
int nox_xxx_collidePentagram_4EAB20(int a1);
nox_object_t* sub_537700(void);
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
static uint32_t* coreRadialRun(nox_object_t* obj, float2* point, float radius, uint32_t code) {
 struct {float2* point; float radius; void (*callback)(int,uint32_t); uint32_t code;} args = {point,radius,coreRadialObserve,code};
 memset(coreRadialResult,0,sizeof(coreRadialResult));
 sub_5180B0((int)obj,(int)&args);
 return coreRadialResult;
}
static void* corePentagramPtr(void){return nox_xxx_collidePentagram_4EAB20;}
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestCollisionCore(op string, a, b *server.Object, p *types.Pointf, mode int32) int32 {
	ai, bi := C.int(uintptr(a.CObj())), C.int(uintptr(b.CObj()))
	switch op {
	case "angle-queue":
		worldAngleQueue(a.UpdateData)
	case "shape":
		geometryShapeBox(&a.Shape)
	case "contains":
		return int32(C.sub_547DB0(ai, (*C.float2)(unsafe.Pointer(p))))
	case "eligible":
		return int32(C.sub_548360(ai, bi))
	case "retained":
		return int32(C.sub_5485B0(ai, bi))
	case "pair":
		C.sub_548220((*C.int)(a.CObj()), (*C.float)(b.CObj()))
	case "scan":
		C.sub_5481C0(ai)
	case "circle-walls":
		C.sub_54FEF0(ai)
	case "circle-box":
		C.sub_54AD50(ai, bi, C.int(mode))
	case "gate-circle":
		C.sub_5488B0((*C.int)(a.CObj()), (*C.float)(b.CObj()), C.int(mode))
	case "dispatch":
		C.nox_xxx_collide_548740()
	case "angles":
		C.sub_548B60()
	case "elevator":
		C.sub_551AE0(ai, bi, C.int(mode))
	case "shaft":
		C.sub_551C40(ai, bi)
	case "types":
		C.sub_551BF0()
	case "active":
		return int32(C.sub_537580(ai))
	case "activate":
		return int32(C.nox_xxx_unitHasCollideOrUpdateFn_537610(asObjectC(a)))
	case "remove":
		C.sub_5375A0(ai)
	case "head":
		return int32(C.sub_537740())
	case "next":
		return int32(C.sub_537750(ai))
	case "pop":
		return int32(uintptr(unsafe.Pointer(C.sub_537700())))
	default:
		panic(op)
	}
	return 0
}
func PortTestCollisionCoreDistance(p *types.Pointf, radius float32, box *server.Object, out *types.Pointf) float64 {
	return float64(C.sub_54A990((*C.float2)(unsafe.Pointer(p)), C.float(radius), C.int(uintptr(box.CObj())), (*C.float2)(unsafe.Pointer(out))))
}
func PortTestCollisionCoreWallOpen(grid *[2]int32, u *server.Object) {
	C.sub_548100((*C.int2)(unsafe.Pointer(grid)), C.int(uintptr(u.CObj())))
}
func PortTestCollisionCoreAddHit(a, b *server.Object, sentinel uint32, normal *types.Pointf) {
	target := C.uint(sentinel)
	if b != nil {
		target = C.uint(uintptr(b.CObj()))
	}
	C.nox_xxx_collSysAddCollision_548630(C.int(uintptr(a.CObj())), target, (*C.float2)(unsafe.Pointer(normal)))
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
	words := map[string]*uint32{"trigger": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2490508))}
	for name, off := range map[string]uintptr{"powder": 2490512, "hand": 2490516, "small": 2491792, "medium": 2491796, "large": 2491800, "meteor": 2491804, "ready": 2491808} {
		words[name] = memmap.PtrUint32(0x5D4594, off)
	}
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
		p := *memmap.PtrUint32(0x5D4594, 2490520+uintptr(4*i))
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
	data := unsafe.Slice((*uint32)(unsafe.Pointer(C.coreRadialRun(asObjectC(u), (*C.float2)(unsafe.Pointer(p)), C.float(radius), C.uint32_t(code)))), 3)
	out := [3]uint32{data[0], data[1], data[2]}
	clear(data)
	return out
}
