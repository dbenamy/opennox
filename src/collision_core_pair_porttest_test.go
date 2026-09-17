//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestCollisionCorePairDispatch(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b, link := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	units := []*server.Object{a, b}
	ids := collisionCoreIDs(a, b, link)
	savedA, savedB, savedLink := a.UpdateData, b.UpdateData, link.UpdateData
	a.UpdateData = collisionCoreGuarded(t, o, 64)
	b.UpdateData = collisionCoreGuarded(t, o, 64)
	link.UpdateData = collisionCoreGuarded(t, o, 64)
	t.Cleanup(func() {
		a.UpdateData = savedA
		b.UpdateData = savedB
		link.UpdateData = savedLink
		a.Update = nil
		b.Update = nil
		o.s.Objs.RemoveFromUpdatable(a)
		o.s.Objs.RemoveFromUpdatable(b)
	})
	ids[a.UpdateData] = 2001
	ids[b.UpdateData] = 2002
	data := [][]byte{unsafe.Slice((*byte)(a.UpdateData), 64), unsafe.Slice((*byte)(b.UpdateData), 64)}
	binary.LittleEndian.PutUint32(unsafe.Slice((*byte)(link.UpdateData), 64)[16:], 64)
	type row struct {
		AbsScratch       uint32
		Classes          [2]uint32
		Boxes            [2]bool
		A, B             worldGeometryObjectState
		Height, Vertical [2]uint32
		Data             [2][]byte
		Hits             [][5]uint32
		Queues           [3]uint32
		Updatable        [2]uint32
	}
	var rows []row
	for _, ac := range []uint32{8, 0x80, 0x4000, 0x8000} {
		for _, bc := range []uint32{8, 0x80, 0x4000, 0x8000} {
			for _, ab := range []bool{false, true} {
				for _, bb := range []bool{false, true} {
					o.resetHits()
					*o.absScratch = 0x3fc00000
					o.resetQueues()
					o.s.SetFrame(123)
					worldGeometryResetObject(a, 1001, 100, 100, ab)
					worldGeometryResetObject(b, 1002, 105, 100, bb)
					a.ObjClass = object.Class(ac)
					b.ObjClass = object.Class(bc)
					for i, u := range units {
						o.s.Objs.RemoveFromUpdatable(u)
						clear(data[i])
						u.ObjFlags = 4
						u.Collide = o.callback
						u.Update = nil
						u.Field115 = 0
						u.Field116 = 0
						u.ObjOwner = nil
						u.TeamVal = server.ObjectTeam{}
						u.ZVal = 0
						u.Field27 = 3
						u.VelVec = types.Pointf{2, -3}
						if u.Shape.Kind == server.ShapeKindBox {
							legacy.PortTestCollisionCore("shape", u, nil, nil, 0)
						}
						if u.ObjClass == object.ClassDoor {
							u.Update = legacy.PortTestGeometryDoorUpdate()
						}
						if uint32(u.ObjClass) == 0x8000 {
							binary.LittleEndian.PutUint32(data[i][4:], uint32(uintptr(link.CObj())))
						}
					}
					legacy.PortTestCollisionCore("pair", a, b, nil, 0)
					r := row{AbsScratch: *o.absScratch, Classes: [2]uint32{ac, bc}, Boxes: [2]bool{ab, bb}, A: worldGeometryState(a), B: worldGeometryState(b), Height: [2]uint32{math.Float32bits(a.ZVal), math.Float32bits(b.ZVal)}, Vertical: [2]uint32{math.Float32bits(a.Field27), math.Float32bits(b.Field27)}, Hits: o.hits(ids), Queues: o.queues(ids), Updatable: [2]uint32{uint32(a.IsUpdatable), uint32(b.IsUpdatable)}}
					for i, u := range units {
						r.Data[i] = append([]byte(nil), data[i]...)
						if uint32(u.ObjClass) == 0x8000 {
							binary.LittleEndian.PutUint32(r.Data[i][4:], 1003)
						}
						if raw := binary.LittleEndian.Uint32(r.Data[i][36:]); raw != 0 {
							binary.LittleEndian.PutUint32(r.Data[i][36:], collisionCoreID(t, ids, raw))
						}
					}
					if ac == 8 && bc == 8 && len(r.Hits) != 1 {
						t.Fatal("ordinary shape-pair route", ab, bb, r.Hits)
					}
					if ac == bc && ac != 8 && len(r.Hits) != 0 {
						t.Fatal("same special class pair accepted", ac)
					}
					rows = append(rows, r)
				}
			}
		}
	}
	o.resetQueues()
	for _, u := range units {
		u.Field115 = 0
		u.Field116 = 0
		u.ObjClass = object.ClassSimple
		u.ObjFlags = 0
	}
	spellbookCapture(t, "collision-core-pair-dispatch", rows, "158e1002c51e6e12346587ae8178eb0945a5ea15cb0d7ec58ecfeab66a6f15bc")
}
