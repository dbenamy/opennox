package legacy

/*
#include "GAME1.h"
#include "GAME1_1.h"
#include "common__random.h"
#include "GAME3_3.h"
#include "GAME4_1.h"
extern uint32_t dword_5d4594_1565628, dword_5d4594_1565632;
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"strings"
	"unsafe"
)

func stateDistance(u, t *server.Object) float64 {
	dx, dy := float64(u.PosVec.X)-float64(t.PosVec.X), float64(u.PosVec.Y)-float64(t.PosVec.Y)
	d := math.Sqrt(dy*dy + dx*dx)
	for _, o := range []*server.Object{u, t} {
		switch o.Shape.Kind {
		case 2:
			d -= float64(o.Shape.Circle.R)
		case 3:
			w := float64(*(*float32)(unsafe.Add(o.CObj(), 184))) * .5
			h := float32(float64(*(*float32)(unsafe.Add(o.CObj(), 188))) * .5)
			if w <= float64(h) {
				w = float64(h)
			}
			d -= w
		}
	}
	if d < 0.0099999998 {
		d = 0.0099999998
	}
	return d
}
func stateDirection(a, b *types.Pointf) int32 {
	dx, dy := float32(b.X-a.X), float32(b.Y-a.Y)
	C.dword_5d4594_1565628 = C.uint32_t(math.Float32bits(dx))
	C.dword_5d4594_1565632 = C.uint32_t(math.Float32bits(dy))
	p := float32(float64(dx)*.41304299 - float64(dy))
	q := float64(dx)*2.4210529 - float64(dy)
	r := float32(float64(dx)*-2.4210529 - float64(dy))
	s := float32(float64(dx)*-.41304299 - float64(dy))
	*memmap.PtrFloat32(0x5d4594, 1565652) = p
	*memmap.PtrFloat32(0x5d4594, 1565656) = float32(q)
	*memmap.PtrFloat32(0x5d4594, 1565636) = r
	*memmap.PtrFloat32(0x5d4594, 1567708) = s
	var bits uint32
	// The first three C comparisons use !(x < 0), including unordered values.
	if !(p < 0) {
		bits |= 8
	}
	if !(q < 0) {
		bits |= 4
	}
	if !(r < 0) {
		bits |= 2
	}
	if s >= 0 {
		bits |= 1
	}
	*memmap.PtrUint32(0x5d4594, 1565640) = bits
	return [16]int32{2, 0, 6, 4, 10, 0, 0, 0, 0, 0, 0, 5, 8, 9, 0, 1}[bits]
}
func stateFront(a *types.Pointf, dir int32, b *types.Pointf) int32 {
	v, free := alloc.New([2]int32{})
	defer free()
	C.nox_xxx_xferIndexedDirection_509E20(C.int(dir), (*C.int2)(unsafe.Pointer(v)))
	ind := stateDirection(a, b) + 16*(v[0]+3*v[1]+4)
	return int32(memmap.Uint32(0x587000, uintptr(202504+4*ind)))
}
func stateTeleport(u *server.Object, p *types.Pointf) {
	if u.Buffs&(1<<14) == 0 && u.ObjFlags&2 == 0 &&
		(!bool(C.nox_common_gameFlags_check_40A5C0(4096)) || u.ObjClass&2 == 0 || u.ObjSubClass&8 == 0) &&
		(bool(C.nox_common_gameFlags_check_40A5C0(2048)) || u.ObjClass&6 != 0) {
		C.nox_xxx_unitMove_4E7010(asObjectC(u), (*C.float2)(unsafe.Pointer(p)))
	}
}
func stateLoot(u *server.Object, p *types.Pointf) {
	off := uintptr(203080)
	if !strings.HasPrefix(alloc.GoString((*byte)(unsafe.Pointer(C.nox_xxx_getUnitName_4E39D0(asObjectC(u))))), "Barrel") {
		off = 203240
	}
	roll := uint32(C.nox_common_randomInt_415FA0(0, 99))
	for *memmap.PtrPtr(0x587000, off) != nil && memmap.Uint32(0x587000, off+8) <= roll {
		off += 12
	}
	if *memmap.PtrPtr(0x587000, off) == nil {
		return
	}
	pos, free := alloc.New(types.Pointf{})
	defer free()
	for i := int32(0); i < int32(memmap.Uint32(0x587000, off+4)); i++ {
		name := alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, off)))
		if it := GetServer().S().NewObjectByTypeID(name); it != nil {
			C.sub_4ED970(35, (*C.float2)(unsafe.Pointer(p)), (*C.float2)(unsafe.Pointer(pos)))
			GetServer().CreateObjectAt(it, nil, *pos)
		}
	}
}
