//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

type worldGeometryObjectState struct {
	Class, Flags                               uint32
	Position, Previous, Velocity, Acceleration [2]uint32
}

func worldGeometryState(u *server.Object) worldGeometryObjectState {
	bits := func(p types.Pointf) [2]uint32 { return [2]uint32{math.Float32bits(p.X), math.Float32bits(p.Y)} }
	return worldGeometryObjectState{uint32(u.ObjClass), uint32(u.ObjFlags), bits(u.NewPos), bits(u.PrevPos), bits(u.VelVec), bits(u.Pos24)}
}
func worldGeometryResetObject(u *server.Object, id uint32, x, y float32, box bool) {
	u.ObjClass = object.ClassSimple
	u.ObjFlags = 0
	u.NetCode = id
	u.Mass = 2
	u.NewPos = types.Pointf{x, y}
	u.PosVec = u.NewPos
	u.PrevPos = u.NewPos
	u.VelVec = types.Pointf{}
	u.Pos24 = types.Pointf{}
	u.Shape = server.Shape{Kind: server.ShapeKindCircle}
	u.Shape.Circle.R = 10
	u.Shape.Circle.R2 = 100
	if box {
		u.Shape = server.Shape{Kind: server.ShapeKindBox}
		u.Shape.Box.W = 20
		u.Shape.Box.H = 20
		u.Shape.Box.Calc()
	}
}
func TestWorldGeometryObjectResponses(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	w.s.PortTestAIEmptyMap()
	configure, unchanged, free := w.s.PortTestPathWalls()
	t.Cleanup(free)
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, w.s), newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001, b.CObj(): 1002}
	oldRNG := w.s.Rand.Logic
	t.Cleanup(func() { w.s.Rand.Logic = oldRNG })
	type row struct {
		Op    string
		Delta [2]float32
		Flags [3]uint32
		Mass  float32
		Wall  int
		A, B  worldGeometryObjectState
		Hits  [][5]uint32
		RNG   int
	}
	var rows []row
	for _, op := range []string{"circle-circle", "box-box"} {
		for _, delta := range [][2]float32{{0, 0}, {8, 0}, {-8, 0}, {0, 8}, {0, -8}, {10, 10}, {19.999, 0}, {20, 0}, {20.001, 0}, {100, 0}} {
			for _, flags := range [][3]uint32{{8, 0, 0}, {8, 8, 0}, {8, 0, 8}, {4, 0, 0x2000}, {4, 0, 0}, {8, 0x8000000, 0x8000000}} {
				for _, mass := range []float32{0.5, 2, 7} {
					for wall := 0; wall < 3; wall++ {
						configure(wall)
						reset()
						w.s.Rand.Logic = prand.New(23)
						worldGeometryResetObject(a, 1001, 100, 100, op == "box-box")
						worldGeometryResetObject(b, 1002, 100+delta[0], 100+delta[1], op == "box-box")
						a.ObjClass = object.Class(flags[0])
						b.ObjClass = a.ObjClass
						a.ObjFlags = object.Flags(flags[1])
						b.ObjFlags = object.Flags(flags[2])
						a.Mass = mass
						b.Mass = 3
						a.VelVec = types.Pointf{2, -3}
						b.VelVec = types.Pointf{-4, 5}
						// Collision uses current positions for the actor visibility ray, independently of proposed positions.
						a.PosVec = types.Pointf{100, 100}
						b.PosVec = types.Pointf{200, 100}
						var scratch [16]uint32
						legacy.PortTestWorldGeometryPhysics(op, a, b, &scratch, 0, 0)
						got := hits(ids)
						wantRNG := 23
						if op == "circle-circle" && delta == ([2]float32{}) {
							wantRNG++
						}
						if w.s.Rand.Logic.Index() != wantRNG {
							t.Fatal("collision RNG", op, delta, w.s.Rand.Logic.Index(), wantRNG)
						}
						if !unchanged() {
							t.Fatal("collision mutated wall")
						}
						// A duplicate contact must not enqueue a second event, even though force accumulates.
						firstA, firstB := worldGeometryState(a), worldGeometryState(b)
						legacy.PortTestWorldGeometryPhysics(op, a, b, &scratch, 0, 0)
						if len(hits(ids)) != len(got) {
							t.Fatal("duplicate collision event", op)
						}
						rows = append(rows, row{op, delta, flags, mass, wall, firstA, firstB, got, wantRNG})
					}
				}
			}
		}
	}
	// Independent magnitude/direction contract: 5 units penetration * force 40 / mass 2.
	configure(0)
	reset()
	worldGeometryResetObject(a, 1001, 100, 100, false)
	worldGeometryResetObject(b, 1002, 115, 100, false)
	var scratch [16]uint32
	legacy.PortTestWorldGeometryPhysics("circle-circle", a, b, &scratch, 0, 0)
	if a.Pos24 != (types.Pointf{-100, 0}) || b.Pos24 != (types.Pointf{}) || len(hits(ids)) != 1 {
		t.Fatal("circle spring contract", a.Pos24, b.Pos24, hits(ids))
	}
	spellbookCapture(t, "world-geometry-object-responses", rows, "51f80077438849d3cd0512a824f046c6aca2553fd7a3881bc589b8ee49f135f3")
}
func TestWorldGeometryPointWall(t *testing.T) {
	w := newWorldCollisionOwner(t)
	globals := worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(w.s.PortTestObjectiveTypes(nil, nil, nil))
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a := newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001}
	typ := a.TypeInd
	ball := uint32(w.s.Types.ByID("GameBall").Ind())
	type row struct {
		Delta, Velocity [2]float32
		Radius, Mass    float32
		Ball            bool
		Return          int
		Object          worldGeometryObjectState
		Hits            [][5]uint32
		Cache           uint32
	}
	var rows []row
	for _, delta := range [][2]float32{{0, 0}, {3, 4}, {-3, 4}, {3, -4}, {-3, -4}, {10, 0}, {10.001, 0}, {9.999, 0}, {0, 20}} {
		for _, vel := range [][2]float32{{0, 0}, {-3, -4}, {3, 4}, {4, -3}} {
			for _, radius := range []float32{0, 5, 10, 20} {
				for _, mass := range []float32{0.5, 2, 7} {
					for _, isBall := range []bool{false, true} {
						reset()
						worldGeometryResetObject(a, 1001, 100+delta[0], 100+delta[1], false)
						a.TypeInd = typ
						if isBall {
							a.TypeInd = uint16(ball)
						}
						a.Mass = mass
						a.Shape.Circle.R = radius
						a.VelVec = types.Pointf{vel[0], vel[1]}
						*globals["gameBall"] = 0
						var scratch [16]uint32
						copy(scratch[:], worldGeometryWords(100, 100))
						rv := legacy.PortTestWorldGeometryPhysics("point-wall", a, nil, &scratch, 0, 0)
						rows = append(rows, row{delta, vel, radius, mass, isBall, rv, worldGeometryState(a), hits(ids), *globals["gameBall"]})
					}
				}
			}
		}
	}
	reset()
	worldGeometryResetObject(a, 1001, 105, 100, false)
	a.TypeInd = typ
	a.VelVec = types.Pointf{-3, 4}
	*globals["gameBall"] = 0
	var scratch [16]uint32
	copy(scratch[:], worldGeometryWords(100, 100))
	rv := legacy.PortTestWorldGeometryPhysics("point-wall", a, nil, &scratch, 0, 0)
	if rv != 1 || a.Pos24 != (types.Pointf{250, 0}) || a.VelVec != (types.Pointf{0, 4}) || len(hits(ids)) != 1 {
		t.Fatal("point-wall response contract", rv, a.Pos24, a.VelVec, hits(ids))
	}
	for _, isBall := range []bool{false, true} {
		a.TypeInd = typ
		if isBall {
			a.TypeInd = uint16(ball)
		}
		*globals["gameBall"] = 0
		for i := 0; i < 2; i++ {
			rv := legacy.PortTestWorldGeometryPhysics("game-ball", a, nil, &scratch, 0, 0)
			if (rv != 0) != isBall || *globals["gameBall"] != ball {
				t.Fatal("cold/hot GameBall lookup", isBall, rv, *globals["gameBall"])
			}
		}
	}
	a.TypeInd = typ
	spellbookCapture(t, "world-geometry-point-wall", rows, "7f3f5a637132f9603cd7883c7d2789756925c9468a0440ea24740701243eaa10")
}
func TestWorldGeometryWallResponseAxes(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a := newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001}
	type row struct {
		Op               string
		Center, Previous [2]float32
		Length, Mass     float32
		Return           int
		Words            [16]uint32
		Object           worldGeometryObjectState
		Hits             [][5]uint32
	}
	var rows []row
	for _, op := range []string{"horizontal-wall", "vertical-wall"} {
		for _, center := range [][2]float32{{0, 0}, {5, 5}, {-5, -5}, {20, 20}, {-20, 0}, {0, -20}, {10, 10}, {10.001, 10.001}} {
			for _, previous := range [][2]float32{center, {-center[0], -center[1]}, {0, 0}} {
				for _, length := range []float32{0, 5, 20, 40} {
					for _, mass := range []float32{0.5, 2, 7} {
						reset()
						worldGeometryResetObject(a, 1001, 100, 100, true)
						a.Mass = mass
						a.VelVec = types.Pointf{3, -4}
						var scratch [16]uint32
						copy(scratch[:], worldGeometryWords(center[0], center[1], previous[0], previous[1], center[0]-10, center[1]-10, center[0]+10, center[1]+10, 0, 0))
						rv := legacy.PortTestWorldGeometryPhysics(op, a, nil, &scratch, 0, length)
						if (rv != 0) != (len(hits(ids)) == 1) {
							t.Fatal("wall response return/event", op, rv, hits(ids))
						}
						rows = append(rows, row{op, center, previous, length, mass, rv, scratch, worldGeometryState(a), hits(ids)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-geometry-wall-response-axes", rows, "349fa76b36fd70f113c92803a9dfc52f7986c2c946a4fd012276df799aa72b09")
}

func TestWorldGeometryCoincidentDirections(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, w.s), newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001, b.CObj(): 1002}
	old := w.s.Rand.Logic
	t.Cleanup(func() { w.s.Rand.Logic = old })
	type row struct {
		Seed, Index int
		Object      worldGeometryObjectState
		Hits        [][5]uint32
	}
	var rows []row
	directions := map[[2]uint32]bool{}
	for seed := 0; seed < 32; seed++ {
		reset()
		w.s.Rand.Logic = prand.New(seed)
		worldGeometryResetObject(a, 1001, 100, 100, false)
		worldGeometryResetObject(b, 1002, 100, 100, false)
		var scratch [16]uint32
		legacy.PortTestWorldGeometryPhysics("circle-circle", a, b, &scratch, 0, 0)
		h := hits(ids)
		if len(h) != 1 || w.s.Rand.Logic.Index() != seed+1 || a.Pos24 == (types.Pointf{}) {
			t.Fatal("coincident response", seed, h, a.Pos24)
		}
		directions[[2]uint32{h[0][2], h[0][3]}] = true
		rows = append(rows, row{seed, w.s.Rand.Logic.Index(), worldGeometryState(a), h})
	}
	if len(directions) != 4 {
		t.Fatal("coincident normals did not cover all four directions", len(directions))
	}
	spellbookCapture(t, "world-geometry-coincident-directions", rows, "6f401ec9d40b8aa88f22efc43c33041fb06e46f11e1b74ea4d0d2c9148425fe9")
}
