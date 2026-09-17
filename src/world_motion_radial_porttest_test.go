//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestWorldMotionRadialScan(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	a, b, c := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	units := []*server.Object{a, b, c}
	saved := []server.Object{*a, *b, *c}
	t.Cleanup(func() {
		for i, u := range units {
			*u = saved[i]
		}
	})
	ids := collisionCoreIDs(units...)
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	type row struct {
		Radius uint32
		Shape  int
		Nil    bool
		Return uint32
		Calls  [][2]uint32
		RNG    int
	}
	var rows []row
	for shape := 0; shape < 3; shape++ {
		for _, radius := range []float32{0, 0.001, 9.999, 10, 10.001, 19.999, 20, 20.001, 40} {
			for _, nilPoint := range []bool{false, true} {
				o.s.PortTestAIEmptyMap()
				o.s.Rand.Logic = prand.New(23)
				for i, u := range units {
					*u = saved[i]
					worldGeometryResetObject(u, uint32(1001+i), 100+float32(10*i), 100, shape == 2)
					u.ObjFlags = 4
					if shape == 0 {
						u.Shape.Kind = server.ShapeKindCenter
					}
					if shape == 1 {
						u.Shape.Circle.R = 1
						u.Shape.Circle.R2 = 1
					}
					if shape == 2 {
						u.Shape.Box.W = 2
						u.Shape.Box.H = 2
						u.Shape.Box.Calc()
					}
					o.s.Map.AddObjectToIndex(u)
				}
				p := types.Pointf{100, 100}
				ptr := &p
				if nilPoint {
					ptr = nil
				}
				rv, calls := legacy.PortTestWorldMotionRadial(ptr, radius, 0xabcdef01)
				if rv != uint32(uintptr(unsafe.Pointer(ptr))) {
					t.Fatal("radial scan returns center pointer")
				}
				if rv != 0 {
					rv = 1
				}
				for i := range calls {
					calls[i][0] = collisionCoreID(t, ids, calls[i][0])
					if calls[i][1] != 0xabcdef01 {
						t.Fatal("radial callback context")
					}
				}
				if nilPoint && len(calls) != 0 {
					t.Fatal("nil radial center invoked callbacks")
				}
				if !nilPoint && shape == 0 {
					seen := map[uint32]bool{}
					for _, c := range calls {
						if seen[c[0]] {
							t.Fatal("duplicate radial candidate")
						}
						seen[c[0]] = true
					}
					for i := 0; i < 3; i++ {
						if seen[uint32(1001+i)] != (float32(10*i) < radius) {
							t.Fatal("radial strict point radius", radius, calls)
						}
					}
				}
				rows = append(rows, row{math.Float32bits(radius), shape, nilPoint, rv, calls, o.s.Rand.Logic.Index()})
			}
		}
	}
	spellbookCapture(t, "world-motion-radial-scan", rows, "90863b691e91f392e4a7c194785655a39c041f9a391b1049024bf0d4f5853577")
}
