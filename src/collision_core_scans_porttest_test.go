//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestCollisionCoreObjectScan(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	ids := collisionCoreIDs(a, b)
	type row struct {
		Kinds         [2]uint32
		Flags         uint32
		Far, Callback bool
		A, B          worldGeometryObjectState
		Hits          [][5]uint32
	}
	var rows []row
	for _, ak := range []server.ShapeKind{server.ShapeKindCenter, server.ShapeKindCircle, server.ShapeKindBox} {
		for _, bk := range []server.ShapeKind{server.ShapeKindCircle, server.ShapeKindBox} {
			for _, flags := range []uint32{0, 8, 0x20, 0x40, 0x60} {
				for _, far := range []bool{false, true} {
					for _, callback := range []bool{false, true} {
						o.s.PortTestAIEmptyMap()
						o.resetHits()
						o.resetQueues()
						worldGeometryResetObject(a, 1001, 100, 100, ak == server.ShapeKindBox)
						worldGeometryResetObject(b, 1002, 105, 100, bk == server.ShapeKindBox)
						a.Shape.Kind = ak
						if far {
							b.NewPos = types.Pointf{500, 500}
							b.PosVec = b.NewPos
						}
						if ak == server.ShapeKindBox {
							legacy.PortTestCollisionCore("shape", a, nil, nil, 0)
						}
						if bk == server.ShapeKindBox {
							legacy.PortTestCollisionCore("shape", b, nil, nil, 0)
						}
						a.ObjFlags = object.Flags(flags)
						b.ObjFlags = object.FlagActive
						a.Collide = o.callback
						b.Collide = o.callback
						if !callback {
							b.Collide = nil
						}
						a.Pos24 = types.Pointf{17, 19}
						a.UpdateCollider(a.NewPos)
						o.s.Map.AddObjectToIndex(b)
						legacy.PortTestCollisionCore("scan", a, nil, nil, 0)
						if flags&0x60 != 0 || far || !callback || ak == server.ShapeKindCenter {
							if len(o.hits(ids)) != 0 || a.Pos24 != (types.Pointf{}) {
								t.Fatal("scan guard/reset", ak, bk, flags, far, callback, a.Pos24, o.hits(ids))
							}
						}
						if flags == 0 && !far && callback && ak != server.ShapeKindCenter && len(o.hits(ids)) != 1 {
							t.Fatal("real indexed pair not dispatched", ak, bk, o.hits(ids))
						}
						rows = append(rows, row{[2]uint32{uint32(ak), uint32(bk)}, flags, far, callback, worldGeometryState(a), worldGeometryState(b), o.hits(ids)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "collision-core-object-scan", rows, "6af4cb9c03228d91f15128f390a3fffccef60dcf2819174058cae3a270c3ce2f")
}
