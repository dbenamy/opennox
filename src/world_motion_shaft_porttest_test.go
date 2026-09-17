//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestWorldMotionShaftFall(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	ids := collisionCoreIDs(u)
	type row struct {
		Z, V                                          uint32
		Delta                                         [2]uint32
		Position, Proposed, Previous, Force, Velocity [2]uint32
		Height, Vertical, Flags, Sync                 uint32
		PlayerSync                                    [32]uint32
		Queues                                        [3]uint32
	}
	var rows []row
	for _, z := range []float32{-60, -50, -49, 0, 10, 90} {
		for _, v := range []float32{-3, -1, 0, 2} {
			for _, delta := range []types.Pointf{{0, 0}, {3, 4}, {-3, -4}, {0.1, -0.2}, {200, -100}} {
				o.s.PortTestAIEmptyMap()
				o.resetQueues()
				*u = saved
				worldGeometryResetObject(u, 1001, 200, 200, false)
				u.ObjFlags = 4 | 0x40000
				u.Collide = o.callback
				u.Update = nil
				u.Field115 = 0
				u.Field116 = 0x40
				u.Field38 = 0x1234
				clear(u.Field140[:])
				u.ZVal = z
				u.Field27 = v
				u.Pos39 = u.PosVec.Sub(delta)
				u.Field41 = math.Float32bits(500)
				u.Field42 = math.Float32bits(600)
				u.VelVec = types.Pointf{7, 8}
				u.ForceVec = types.Pointf{9, 10}
				legacy.PortTestWorldMotionPhysics("fall", u)
				returned := z+v < -50
				if u.VelVec != (types.Pointf{}) || u.Field27 != v-1 {
					t.Fatal("shaft fall must clear velocity and accelerate downward")
				}
				if returned {
					if u.ZVal != 90 || u.PosVec != (types.Pointf{500, 600}) || u.NewPos != u.PosVec || u.PrevPos != u.PosVec || u.ObjFlags&0x40000 != 0 || u.Field116&1 == 0 {
						t.Fatal("shaft return must relocate, reactivate and clear falling state")
					}
				} else if u.ZVal != z+v || u.PosVec != (types.Pointf{200, 200}) || u.ObjFlags&0x40000 == 0 {
					t.Fatal("shaft height/boundary contract")
				}
				if delta == (types.Pointf{}) && u.ForceVec != (types.Pointf{}) {
					t.Fatal("centered shaft pull must be zero")
				}
				if delta == (types.Pointf{3, 4}) && u.ForceVec != (types.Pointf{-1.8, -2.4}) {
					t.Fatal("shaft pull direction and magnitude", u.ForceVec)
				}
				rows = append(rows, row{math.Float32bits(z), math.Float32bits(v), motionBits(delta), motionBits(u.PosVec), motionBits(u.NewPos), motionBits(u.PrevPos), motionBits(u.ForceVec), motionBits(u.VelVec), math.Float32bits(u.ZVal), math.Float32bits(u.Field27), uint32(u.ObjFlags), u.Field38, u.Field140, o.queues(ids)})
			}
		}
	}
	spellbookCapture(t, "world-motion-shaft-fall", rows, "2109921813e9f353809b61f36257002cbcd2e13237942afd714fc681ff2f48c6")
}
