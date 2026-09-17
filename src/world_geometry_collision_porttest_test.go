//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestWorldGeometryBoxFlagContracts(t *testing.T) {
	world := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, world.s), newObjectXferSimple(t, world.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001, b.CObj(): 1002}
	for _, tc := range []struct {
		class, af, bf uint32
		force         bool
	}{{8, 0, 0, true}, {8, 8, 0, false}, {8, 0, 8, false}, {4, 0, 0x2000, false}, {8, 0x8000000, 0, true}, {8, 0, 0x8000000, true}} {
		t.Run(fmt.Sprintf("class%x/a%x/b%x", tc.class, tc.af, tc.bf), func(t *testing.T) {
			reset()
			for i, u := range []*server.Object{a, b} {
				u.ObjClass = object.ClassSimple
				u.ObjFlags = 0
				u.NetCode = uint32(1001 + i)
				u.Mass = 2
				u.NewPos = types.Pointf{100 + float32(i)*8, 100}
				u.PosVec = u.NewPos
				u.PrevPos = u.NewPos
				u.VelVec = types.Pointf{}
				u.Pos24 = types.Pointf{}
				u.Shape = server.Shape{Kind: server.ShapeKindBox}
				u.Shape.Box.W = 20
				u.Shape.Box.H = 20
				u.Shape.Box.Calc()
			}
			a.ObjClass = object.Class(tc.class)
			a.ObjFlags = object.Flags(tc.af)
			b.ObjFlags = object.Flags(tc.bf)
			var scratch [16]uint32
			legacy.PortTestWorldGeometryPhysics("box-box", a, b, &scratch, 0, 0)
			if got := hits(ids); len(got) != 1 || got[0][0] != 1002 || got[0][1] != 1001 {
				t.Fatal("real collision queue", got)
			}
			moved := a.Pos24 != (types.Pointf{})
			if moved != tc.force {
				t.Errorf("force applied %t want %t; acceleration %v", moved, tc.force, a.Pos24)
			}
			if a.ObjFlags&0x8000000 != 0 || b.ObjFlags&0x8000000 != 0 {
				t.Errorf("collision wake flags remain: %x %x", a.ObjFlags, b.ObjFlags)
			}
			// The selected response accumulates force on its first object only.
			if b.Pos24 != (types.Pointf{}) {
				t.Error("second object acceleration changed", b.Pos24)
			}
		})
	}
}
