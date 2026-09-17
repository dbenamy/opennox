//go:build porttest

package legacy

/*
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME3_2.h"
#include "GAME4_1.h"
#include "GAME5.h"
int sub_427C80(int4*, int4*);
extern unsigned int dword_587000_230092;
extern uint32_t dword_587000_292488, dword_587000_292492;
extern uint32_t dword_5d4594_2491544;
extern void* nox_alloc_hit_2491548;
*/
import "C"

import (
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
	f := func(i int) C.float { return C.float(math.Float32frombits(spec.Floats[i])) }
	var rv C.int
	switch spec.Op {
	case "segments":
		rv = C.sub_427C80((*C.int4)(ptr[0]), (*C.int4)(ptr[1]))
	case "edge-project":
		rv = C.sub_427DF0(C.int(uintptr(ptr[0])), (*C.int)(ptr[1]), f(0))
	case "wall-point":
		rv = C.nox_xxx_wallMath_427F30((*C.int2)(ptr[0]), (*C.int)(ptr[1]))
	case "wall-bounds":
		rv = C.sub_428170(ptr[0], (*C.int4)(ptr[1]))
	case "rect-int":
		rv = C.nox_xxx_pointInRect_4281F0((*C.int2)(ptr[0]), (*C.int4)(ptr[1]))
	case "rect-float":
		rv = C.sub_428220((*C.float2)(ptr[0]), (*C.float4)(ptr[1]))
	case "box-calc":
		C.nox_shape_box_calc((*C.nox_shape)(ptr[0]))
	case "map-coordinates":
		rv = C.sub_4D3E30((*C.float2)(ptr[0]), (*C.float2)(ptr[1]))
	case "direction-angle":
		rv = C.nox_xxx_xferDirectionToAngle_509E00((*C.uint32_t)(ptr[0]))
	case "indexed-direction":
		rv = C.nox_xxx_xferIndexedDirection_509E20(C.int(spec.Ints[0]), (*C.int2)(ptr[0]))
	case "direction4-angle":
		rv = C.nox_xxx_mathDirection4ToAngle_509E90(C.int(spec.Ints[0]))
	case "direction4-index":
		rv = C.nox_xxx_math_509EA0(C.int(spec.Ints[0]))
	case "vector-angle":
		rv = C.nox_xxx_math_509ED0((*C.float2)(ptr[0]))
	case "normalize":
		C.nox_xxx_utilNormalizeVector_509F20((*C.float2)(ptr[0]))
	case "project-positive":
		rv = C.sub_550280(C.int(uintptr(ptr[0])), f(0), f(1), C.int(spec.Ints[0]), C.int(spec.Ints[1]), C.int(uintptr(ptr[1])), C.int(uintptr(ptr[2])))
	case "project-negative":
		rv = C.sub_5502F0((*C.float2)(ptr[0]), f(0), f(1), C.int(spec.Ints[0]), C.int(spec.Ints[1]), (*C.float2)(ptr[1]), (*C.float2)(ptr[2]))
	case "quadrant":
		rv = C.int(C.sub_550CB0((*C.float2)(ptr[0]), (*C.float2)(ptr[1])))
	case "rectangle-crossings":
		rv = C.sub_5516A0((*C.float4)(ptr[0]), (*C.float4)(ptr[1]), (*C.float2)(ptr[2]), C.int(spec.Ints[0]), C.int(spec.Ints[1]))
	case "horizontal-crossing":
		rv = C.sub_551780((*C.float4)(ptr[0]), f(0), f(1), f(2), (*C.float2)(ptr[1]), C.int(spec.Ints[0]))
	case "vertical-crossing":
		rv = C.sub_551870((*C.float4)(ptr[0]), f(0), f(1), f(2), (*C.float2)(ptr[1]), C.int(spec.Ints[0]))
	case "clipped-center":
		rv = C.sub_551960((*C.float4)(ptr[0]), (*C.float4)(ptr[1]), (*C.float4)(ptr[2]), (*C.float2)(ptr[3]))
	case "rect-float-alt":
		rv = C.sub_551A90((*C.float2)(ptr[0]), (*C.float4)(ptr[1]))
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
	a, b := C.int(uintptr(u.CObj())), C.int(uintptr(v.CObj()))
	p := unsafe.Pointer(&words[0])
	p1 := unsafe.Pointer(&words[2])
	p2 := unsafe.Pointer(&words[4])
	p3 := unsafe.Pointer(&words[8])
	switch op {
	case "circle-wall":
		return int(C.sub_54FFC0((*C.int2)(p), a))
	case "point-wall":
		return int(C.sub_550380(C.int(flags), a, (*C.float2)(p)))
	case "game-ball":
		return int(C.sub_550480(a))
	case "box-walls":
		C.sub_5504B0(a)
	case "box-wall":
		return int(C.sub_550580((*C.int2)(p), (*C.float)(u.CObj())))
	case "horizontal-wall":
		return int(C.sub_550760(a, (*C.float2)(p), (*C.float2)(p1), (*C.float4)(p2), (*C.float2)(p3), C.float(distance)))
	case "vertical-wall":
		return int(C.sub_550A10(a, (*C.float2)(p), (*C.float2)(p1), (*C.float4)(p2), (*C.float2)(p3), C.float(distance)))
	case "circle-circle":
		C.nox_xxx_collisionCheckCircleCircle_550D00(a, b)
	case "box-box":
		C.sub_550F80((*C.float)(u.CObj()), b)
	case "gate-box":
		C.sub_551250(C.uint(a), (*C.float)(v.CObj()), C.int(flags))
	default:
		panic(op)
	}
	return 0
}
func PortTestWorldGeometryGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"directionThreshold": (*uint32)(unsafe.Pointer(&C.dword_587000_230092)),
		"objectForce":        (*uint32)(unsafe.Pointer(&C.dword_587000_292488)),
		"wallForce":          (*uint32)(unsafe.Pointer(&C.dword_587000_292492)),
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
	oldClass, oldHead := C.nox_alloc_hit_2491548, C.dword_5d4594_2491544
	buckets := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2490520), 256)
	oldBuckets := append([]uint32(nil), buckets...)
	C.nox_alloc_hit_2491548 = nil
	C.dword_5d4594_2491544 = 0
	clear(buckets)
	reset = func() { C.nox_xxx_allocHitArray_5486D0() }
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
		for p := uint32(C.dword_5d4594_2491544); p != 0; {
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
		alloc.AsClass(C.nox_alloc_hit_2491548).Free()
		C.nox_alloc_hit_2491548 = oldClass
		C.dword_5d4594_2491544 = oldHead
		copy(buckets, oldBuckets)
	}
	return
}
