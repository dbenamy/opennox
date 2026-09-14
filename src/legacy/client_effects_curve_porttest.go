//go:build porttest

package legacy

/*
#include "GAME3_1.h"
extern void portTestEffectsCurveSegment(int2*, int2*, int);
*/
import "C"

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

type PortTestEffectsCurveSpec struct {
	Points [4][2]int32
	Steps  int
	Shift  float32
	Token  int32
}
type PortTestEffectsCurveSegment struct {
	From, To [2]int32
	Token    int32
}

var portTestEffectsCurveSegments *[]PortTestEffectsCurveSegment

//export portTestEffectsCurveSegment
func portTestEffectsCurveSegment(a, b *C.int2, token C.int) {
	*portTestEffectsCurveSegments = append(*portTestEffectsCurveSegments, PortTestEffectsCurveSegment{
		From: *(*[2]int32)(unsafe.Pointer(a)), To: *(*[2]int32)(unsafe.Pointer(b)), Token: int32(token),
	})
}

// PortTestEffectsCurve exercises the production C callback ABI, recording exact
// ordered segments and userdata. All input storage is owned and checked here.
func PortTestEffectsCurve(sp PortTestEffectsCurveSpec) []PortTestEffectsCurveSegment {
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x581450, 9872)), 64)
	oldTable := append([]byte(nil), table...)
	copy(table, blobdata.PortTestClientEffectsCurveTable())
	defer copy(table, oldTable)
	p, free := alloc.Make([]byte{}, 48)
	defer free()
	for i := range p {
		p[i] = 0xa5
	}
	points := (*[4][2]int32)(unsafe.Pointer(&p[8]))
	*points = sp.Points
	var out []PortTestEffectsCurveSegment
	old := portTestEffectsCurveSegments
	portTestEffectsCurveSegments = &out
	defer func() { portTestEffectsCurveSegments = old }()
	C.sub_4BEDE0((*C.int2)(unsafe.Pointer(&points[0])), (*C.int2)(unsafe.Pointer(&points[1])), (*C.int2)(unsafe.Pointer(&points[2])), (*C.int2)(unsafe.Pointer(&points[3])), C.int(sp.Steps), C.float(sp.Shift), C.int(uintptr(unsafe.Pointer(C.portTestEffectsCurveSegment))), C.int(sp.Token))
	if *points != sp.Points {
		panic("curve mutated input points")
	}
	for i := 0; i < 8; i++ {
		if p[i] != 0xa5 || p[40+i] != 0xa5 {
			panic("curve input guard changed")
		}
	}
	return out
}
