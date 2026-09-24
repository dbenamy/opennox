//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
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

// PortTestEffectsCurve records the ordered native segments and userdata.
// Input storage remains owned and guarded here.
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
	nativePoints := [4]image.Point{}
	for i, point := range *points {
		nativePoints[i] = image.Pt(int(point[0]), int(point[1]))
	}
	effectCurveSegments(nativePoints, sp.Steps, sp.Shift, func(from, to image.Point) {
		out = append(out, PortTestEffectsCurveSegment{
			From:  [2]int32{int32(from.X), int32(from.Y)},
			To:    [2]int32{int32(to.X), int32(to.Y)},
			Token: sp.Token,
		})
	})
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
