//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestCollisionCoreCircleBoxResponse(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	walls, intact, restore := o.s.PortTestPathWalls()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ids := collisionCoreIDs(a, b)
	old := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = old })
	type row struct {
		Mode, Wall, RNG int
		Offset, Size    [2]float32
		Flags           [3]uint32
		Mass            float32
		A, B            worldGeometryObjectState
		Hits            [][5]uint32
	}
	var rows []row
	for mode := 0; mode < 2; mode++ {
		for _, offset := range [][2]float32{{0, 0}, {5, 0}, {-5, 0}, {0, 5}, {0, -5}, {10, 10}, {20, 0}, {24.141, 0}, {24.142, 0}, {24.143, 0}, {40, 40}} {
			for _, size := range [][2]float32{{20, 20}, {30.25, 10}, {7.25, 12.5}} {
				for _, flags := range [][3]uint32{{8, 0, 0}, {8, 8, 0}, {8, 0, 8}, {4, 0, 0x2000}, {4, 0, 0}, {8, 0x8000000, 0x8000000}} {
					for _, mass := range []float32{0.5, 2, 7} {
						for wall := 0; wall < 3; wall++ {
							o.resetHits()
							o.resetQueues()
							walls(wall)
							o.s.Rand.Logic = prand.New(23)
							worldGeometryResetObject(a, 1001, 100+offset[0], 100+offset[1], false)
							worldGeometryResetObject(b, 1002, 100, 100, true)
							b.Shape.Box.W = size[0]
							b.Shape.Box.H = size[1]
							legacy.PortTestCollisionCore("shape", b, nil, nil, 0)
							a.ObjClass = object.Class(flags[0])
							b.ObjClass = a.ObjClass
							a.ObjFlags = object.Flags(flags[1])
							b.ObjFlags = object.Flags(flags[2])
							a.Mass = mass
							b.Mass = 3
							a.VelVec = types.Pointf{2, -3}
							b.VelVec = types.Pointf{-4, 5}
							a.PosVec = types.Pointf{100, 100}
							b.PosVec = types.Pointf{200, 100}
							legacy.PortTestCollisionCore("circle-box", a, b, nil, int32(mode))
							if mode == 0 && b.Pos24 != (types.Pointf{}) || mode == 1 && a.Pos24 != (types.Pointf{}) {
								t.Fatal("circle/box force applied to wrong object", mode)
							}
							if flags[1]&8 != 0 || flags[2]&8 != 0 || mode == 0 && flags[0]&6 != 0 && flags[2]&0x2000 != 0 {
								if a.Pos24 != (types.Pointf{}) || b.Pos24 != (types.Pointf{}) {
									t.Fatal("circle/box force suppression")
								}
							}
							if !intact() {
								t.Fatal("circle/box changed ray walls")
							}
							rows = append(rows, row{mode, wall, o.s.Rand.Logic.Index(), offset, size, flags, mass, worldGeometryState(a), worldGeometryState(b), o.hits(ids)})
						}
					}
				}
			}
		}
	}
	// Preserve the interior branch: offset five gives a five-unit response, not
	// a newly substituted textbook penetration calculation.
	walls(0)
	for mode := 0; mode < 2; mode++ {
		o.resetHits()
		worldGeometryResetObject(a, 1001, 105, 100, false)
		worldGeometryResetObject(b, 1002, 100, 100, true)
		legacy.PortTestCollisionCore("shape", b, nil, nil, 0)
		legacy.PortTestCollisionCore("circle-box", a, b, nil, int32(mode))
		wantA, wantB := types.Pointf{100, 0}, types.Pointf{}
		if mode == 1 {
			wantA, wantB = types.Pointf{}, types.Pointf{-100, 0}
		}
		if a.Pos24 != wantA || b.Pos24 != wantB || len(o.hits(ids)) != 1 {
			t.Fatal("circle/box interior response", mode, a.Pos24, b.Pos24, o.hits(ids))
		}
	}
	spellbookCapture(t, "collision-core-circle-box-response", rows, "eef97a727a2433c89ef2de248c25df84bbe04e5589be193fd3e63ffd0cff9067")
}
