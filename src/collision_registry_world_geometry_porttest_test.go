//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestCollisionRegistryWorldTrapGeometry(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a, b := &o.units[0], &o.units[1]
	a.CollideData = o.record(t, 28)
	objectXferSetWord(a.CollideData, 8, uint32(123))
	objectXferSetWord(a.CollideData, 12, 0xffffffd3)
	a.PosVec = types.Pointf{100, 100}
	a.Shape = server.Shape{Kind: server.ShapeKindBox}
	a.Shape.Box.W = 20
	a.Shape.Box.H = 20
	a.Shape.Box.Calc()
	var rows []struct {
		Name        string
		Flags, X, Y uint32
		Origin      types.Pointf
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-trap-geometry", rows)
	}()
	for _, kind := range []server.ShapeKind{server.ShapeKindNone, server.ShapeKindCenter, server.ShapeKindCircle, server.ShapeKindBox} {
		for _, size := range []float32{0, 9, 10, 10.0001, 20, 21} {
			for _, pos := range []types.Pointf{{100, 100}, {101, 102}, {110, 100}, {120, 100}, {100, 120}} {
				for _, door := range []bool{false, true} {
					name := fmt.Sprintf("kind%d/size%g/pos%v/door%t", kind, size, pos, door)
					t.Run(name, func(t *testing.T) {
						a.ObjFlags = object.Flags(0x1000000)
						b.ObjFlags = 0
						b.ObjClass = object.ClassSimple
						if door {
							b.ObjClass = object.ClassDoor
						}
						b.Shape = server.Shape{Kind: kind}
						b.Shape.Circle.R = size
						b.Shape.Circle.R2 = size * size
						b.Shape.Box.W = size
						b.Shape.Box.H = size
						b.PosVec = pos
						b.Pos39 = types.Pointf{-1, -2}
						b.Field41 = 77
						b.Field42 = 88
						collisionRegistryWorld(14, a, b, nil)
						fits := true
						if kind == server.ShapeKindCircle {
							fits = size*2 <= 20
						}
						if kind == server.ShapeKindBox {
							fits = size <= 20
						}
						// The 20x20 rotated box contains these first three probe points, and excludes
						// the last two. Its strict boundary predicate is separately qualified.
						inside := pos.X < 120 && pos.Y < 120
						want := !door && fits && inside
						if (b.ObjFlags&0x60000 == 0x60000) != want {
							t.Fatalf("trap admission flags%x want%t", b.ObjFlags, want)
						}
						if want {
							if b.Field41 != math.Float32bits(123) || b.Field42 != math.Float32bits(-45) || b.Pos39 != a.PosVec {
								t.Fatal("trap destination")
							}
						} else if b.Field41 != 77 || b.Field42 != 88 || b.Pos39 != (types.Pointf{-1, -2}) {
							t.Fatal("rejected trap changed destination")
						}
						rows = append(rows, struct {
							Name        string
							Flags, X, Y uint32
							Origin      types.Pointf
						}{name, uint32(b.ObjFlags), b.Field41, b.Field42, b.Pos39})
					})
				}
			}
		}
	}
}
