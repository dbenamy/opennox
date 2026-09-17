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

func TestCollisionCoreShapeContainment(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	type row struct {
		Kind         uint32
		Size, Offset [2]float32
		Result       int32
	}
	var rows []row
	for _, kind := range []server.ShapeKind{server.ShapeKindCenter, server.ShapeKindCircle, server.ShapeKindBox} {
		for _, size := range [][2]float32{{0, 0}, {10, 10}, {20, 30.25}, {7.25, 12.5}} {
			worldGeometryResetObject(u, 1001, 100, 100, kind == server.ShapeKindBox)
			u.Shape.Kind = kind
			u.Shape.Circle.R = size[0]
			u.Shape.Circle.R2 = size[0] * size[0]
			u.Shape.Box.W = size[0]
			u.Shape.Box.H = size[1]
			if kind == server.ShapeKindBox {
				legacy.PortTestCollisionCore("shape", u, nil, nil, 0)
			}
			for _, x := range []float32{-30, -10, -9.999, -0.001, 0, 0.001, 9.999, 10, 30} {
				for _, y := range []float32{-30, -10, -0.001, 0, 0.001, 10, 30} {
					p := types.Pointf{100 + x, 100 + y}
					rv := legacy.PortTestCollisionCore("contains", u, nil, &p, 0)
					if kind == server.ShapeKindCenter && rv != 0 {
						t.Fatal("center shape contains point")
					}
					if kind == server.ShapeKindCircle {
						dx, dy := float64(p.X)-100, float64(p.Y)-100
						want := int32(0)
						if dx*dx+dy*dy <= float64(u.Shape.Circle.R2) {
							want = 1
						}
						if rv != want {
							t.Fatal("circle containment", p, rv, want)
						}
					}
					rows = append(rows, row{uint32(kind), size, [2]float32{x, y}, rv})
				}
			}
		}
	}
	spellbookCapture(t, "collision-core-shape-containment", rows, "93689cc9490586e38875a1988342b31d0dca786d0eee609bd13442ce18687d24")
}
func TestCollisionCoreBoxDistance(t *testing.T) {
	o := newCollisionCoreOwner(t)
	box := newObjectXferSimple(t, o.s)
	old := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = old })
	type row struct {
		Size, Offset [2]float32
		Radius       float32
		Seed, RNG    int
		Distance     uint64
		Normal       [2]uint32
	}
	var rows []row
	offsets := []float32{-40, -20, -14.142, -10, -0.001, 0, 0.001, 10, 14.142, 20, 40}
	for _, size := range [][2]float32{{0, 0}, {20, 20}, {30.25, 10}, {7.25, 12.5}} {
		worldGeometryResetObject(box, 1001, 100, 100, true)
		box.Shape.Box.W = size[0]
		box.Shape.Box.H = size[1]
		legacy.PortTestCollisionCore("shape", box, nil, nil, 0)
		for _, x := range offsets {
			for _, y := range offsets {
				for _, radius := range []float32{0, 0.1, 5, 20} {
					o.s.Rand.Logic = prand.New(23)
					p := types.Pointf{100 + x, 100 + y}
					out := [4]uint32{0xa5a5a5a5, 0x3f000000, 0xbf000000, 0x5a5a5a5a}
					d := legacy.PortTestCollisionCoreDistance(&p, radius, box, (*types.Pointf)(unsafe.Pointer(&out[1])))
					if out[0] != 0xa5a5a5a5 || out[3] != 0x5a5a5a5a {
						t.Fatal("distance output guard")
					}
					if x == 40 && y == 0 && radius == 0 && d < 0 && (out[1] != 0x3f000000 || out[2] != 0xbf000000) {
						t.Fatal("separated corner changed output normal")
					}
					rows = append(rows, row{size, [2]float32{x, y}, radius, 23, o.s.Rand.Logic.Index(), math.Float64bits(d), [2]uint32{out[1], out[2]}})
				}
			}
		}
	}
	// Coincident centers use exactly one RNG draw and all four shipped normals.
	worldGeometryResetObject(box, 1001, 100, 100, true)
	legacy.PortTestCollisionCore("shape", box, nil, nil, 0)
	normals := map[[2]uint32]bool{}
	for seed := 0; seed < 32; seed++ {
		o.s.Rand.Logic = prand.New(seed)
		var out types.Pointf
		p := box.NewPos
		d := legacy.PortTestCollisionCoreDistance(&p, 10, box, &out)
		bits := [2]uint32{math.Float32bits(out.X), math.Float32bits(out.Y)}
		normals[bits] = true
		if o.s.Rand.Logic.Index() != seed+1 {
			t.Fatal("coincident box RNG consumption")
		}
		rows = append(rows, row{[2]float32{20, 20}, [2]float32{}, 10, seed, o.s.Rand.Logic.Index(), math.Float64bits(d), bits})
	}
	if len(normals) != 4 {
		t.Fatal("coincident box normals", len(normals))
	}
	spellbookCapture(t, "collision-core-box-distance", rows, "906b2e46d12aac568c69a8221ede2e7453eaba76544265c19d3536c9f4064be7")
}

func TestCollisionCoreRadialCaller(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	old := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = old })
	type row struct {
		Kind   uint32
		Offset [2]float32
		Radius float32
		Code   uint32
		Calls  [3]uint32
		RNG    int
	}
	var rows []row
	for _, kind := range []server.ShapeKind{server.ShapeKindCenter, server.ShapeKindCircle, server.ShapeKindBox} {
		worldGeometryResetObject(u, 1001, 100, 100, kind == server.ShapeKindBox)
		u.Shape.Kind = kind
		u.Shape.Circle.R = 10
		u.Shape.Circle.R2 = 100
		if kind == server.ShapeKindBox {
			legacy.PortTestCollisionCore("shape", u, nil, nil, 0)
		}
		// Deliberately separate the coordinates: circles/points read PosVec,
		// whereas the retained C box caller delegates using NewPos.
		u.PosVec = types.Pointf{120, 100}
		for _, x := range []float32{-40, -20, -10, -0.001, 0, 0.001, 9.999, 10, 10.001, 19.999, 20, 20.001, 30, 40} {
			for _, y := range []float32{0, 5, 20} {
				for _, radius := range []float32{0, 5, 10} {
					for _, code := range []uint32{0, 0xffffffff} {
						o.s.Rand.Logic = prand.New(23)
						p := types.Pointf{100 + x, 100 + y}
						got := legacy.PortTestCollisionCoreRadial(u, &p, radius, code)
						if got[0] > 1 {
							t.Fatal("duplicate radial callback", got)
						}
						if got[0] != 0 {
							if got[1] != uint32(uintptr(u.CObj())) || got[2] != code {
								t.Fatal("radial callback arguments", got)
							}
							got[1] = 1001
						} else if got[1] != 0 || got[2] != 0 {
							t.Fatal("unexpected radial callback state", got)
						}
						if kind != server.ShapeKindBox {
							distance := math.Hypot(float64(p.X-u.PosVec.X), float64(p.Y-u.PosVec.Y))
							if kind == server.ShapeKindCircle {
								distance -= 10
							}
							want := uint32(0)
							if float64(radius) > distance {
								want = 1
							}
							if got[0] != want {
								t.Fatal("strict radial contact boundary", kind, p, radius, got, want)
							}
						}
						rows = append(rows, row{uint32(kind), [2]float32{x, y}, radius, code, got, o.s.Rand.Logic.Index()})
					}
				}
			}
		}
	}
	spellbookCapture(t, "collision-core-radial-caller", rows, "30fd15ff066a7f60e6d172611b6aaeb10c3c8d86fd72948d83f8975f552a926c")
}
