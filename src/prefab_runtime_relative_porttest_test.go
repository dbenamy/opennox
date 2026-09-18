//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestPrefabRuntimeRelativeCoordinates(t *testing.T) {
	s := newObjectXferOwner(t)
	u := s.NewObjectByTypeInd(1)
	if u == nil {
		t.Fatal("object allocation")
	}
	defer s.Objs.FreeObject(u)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	output, free := alloc.Make([]uint32{}, 4)
	defer free()
	type row struct {
		Name             string
		Return           uint64
		Position, Output [2]uint32
		Guards           [2]uint32
	}
	var rows []row
	inverse := func(p types.Pointf) (types.Pointf, types.Pointf) {
		if p.X <= 80.5 {
			p.X = 82.5
		}
		if p.Y <= 80.5 {
			p.Y = 81.5
		}
		if p.X >= 5853.5 {
			p.X = 5851.5
		}
		if p.Y >= 5853.5 {
			p.Y = 5852.5
		}
		return p, types.Ptf(float32((float64(p.X)-1-float64(p.Y))*0.70710677), float32((float64(p.Y)+float64(p.X)-5912)*0.70710677))
	}
	for _, state := range [][3]uint32{{0, 0, 0}, {1, 1, 0}, {0, 1, 0}, {math.MaxUint32, math.MaxUint32, 0}, {1, 1, 1}, {1, 1, 2}} {
		for _, origin := range []int32{0, 81, 2957, 5854, math.MinInt32, math.MaxInt32} {
			for _, pos := range []types.Pointf{{0, 0}, {80.5, 80.5}, {82.5, 81.5}, {2957.25, 2956.75}, {5853.5, 5853.5}, {-10000, 10000}} {
				*words["loaded"], *words["selected"], *words["placed"] = state[0], state[1], state[2]
				*memmap.PtrInt32(0x5D4594, 1599508) = origin
				*memmap.PtrInt32(0x5D4594, 1599512) = -origin
				u.PosVec = pos
				output[0], output[1], output[2], output[3] = 0xa5a5a5a5, 0x11223344, 0x55667788, 0x5a5a5a5a
				ret := legacy.PortTestPrefabCall(16, [6]uint32{uint32(uintptr(u.CObj())), uint32(uintptr(unsafe.Pointer(&output[1])))})
				wantRet := uint64(0)
				wantPos := pos
				wantOut := [2]uint32{0x11223344, 0x55667788}
				if state[0] == state[1] && state[0] != math.MaxUint32 && state[2] != 1 {
					wantRet = 1
					var object, base types.Pointf
					wantPos, object = inverse(pos)
					_, base = inverse(types.Ptf(float32(origin), float32(-origin)))
					wantOut = [2]uint32{math.Float32bits(object.X - base.X), math.Float32bits(object.Y - base.Y)}
				}
				if ret != wantRet || u.PosVec != wantPos || output[1] != wantOut[0] || output[2] != wantOut[1] || output[0] != 0xa5a5a5a5 || output[3] != 0x5a5a5a5a {
					t.Fatalf("state%v/origin%d/pos%v returned%d output%x want%x", state, origin, pos, ret, output, wantOut)
				}
				rows = append(rows, row{fmt.Sprintf("state%v/origin%d/pos%v", state, origin, pos), ret, [2]uint32{math.Float32bits(u.PosVec.X), math.Float32bits(u.PosVec.Y)}, [2]uint32{output[1], output[2]}, [2]uint32{output[0], output[3]}})
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-relative-coordinates", rows, "805d23e46fc3ad50c13f455889d1a9c7d176dbc549b130a9cf7ad498ddc8fbe6")
}
