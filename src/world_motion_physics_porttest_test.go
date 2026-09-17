//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func motionBits(p types.Pointf) [2]uint32 {
	return [2]uint32{math.Float32bits(p.X), math.Float32bits(p.Y)}
}

func TestWorldMotionProjectileIntegration(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	type row struct {
		Shape                                            int
		Input                                            [7]uint32
		Position, NewPosition, Velocity, Force, Min, Max [2]uint32
	}
	var rows []row
	for shape := 0; shape < 3; shape++ {
		for _, pos := range []types.Pointf{{100, 200}, {0, 0}, {16777216, -16777216}, {0.001, -0.001}} {
			for _, vel := range []types.Pointf{{0, 0}, {1, -2}, {0.1234567, -3.1415927}, {-16777216, 16777216}} {
				for _, force := range []types.Pointf{{0, 0}, {3, 4}, {-0.000123456, 0.000123456}, {16777216, -16777216}} {
					for _, drag := range []float32{0, 0.1, 0.5, 1, 1.0000001, -0.5} {
						*u = saved
						worldGeometryResetObject(u, 1001, 13, 17, shape == 2)
						if shape == 0 {
							u.Shape.Kind = server.ShapeKindCenter
						}
						u.NewPos = pos
						u.VelVec = vel
						u.ForceVec = force
						u.Float28 = drag
						legacy.PortTestWorldMotionPhysics("projectile", u)
						if u.PosVec != pos || u.ForceVec != force {
							t.Fatal("projectile must retain force and advance current position to previous proposed position")
						}
						// The simple integral cases give an independent kinematic contract.
						if vel == (types.Pointf{1, -2}) && force == (types.Pointf{3, 4}) && drag == 0.5 {
							if u.VelVec != (types.Pointf{2, 1}) || u.NewPos != (types.Pointf{pos.X + 2, pos.Y + 1}) {
								t.Fatal("projectile force/drag integration", u.VelVec, u.NewPos)
							}
						}
						if shape == 0 && (u.CollideP1 != pos || u.CollideP2 != pos) {
							t.Fatal("center collider position")
						}
						if shape == 1 && (u.CollideP1 != (types.Pointf{pos.X - 10, pos.Y - 10}) || u.CollideP2 != (types.Pointf{pos.X + 10, pos.Y + 10})) {
							t.Fatal("projectile collider uses current position")
						}
						rows = append(rows, row{shape, [7]uint32{math.Float32bits(pos.X), math.Float32bits(pos.Y), math.Float32bits(vel.X), math.Float32bits(vel.Y), math.Float32bits(force.X), math.Float32bits(force.Y), math.Float32bits(drag)}, motionBits(u.PosVec), motionBits(u.NewPos), motionBits(u.VelVec), motionBits(u.ForceVec), motionBits(u.CollideP1), motionBits(u.CollideP2)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-projectile-integration", rows, "0d2adb0fe5e032bfb917e37e4211d031278645c441885e73be6f28ea50c4cf47")
}

func TestWorldMotionFall(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	ids := collisionCoreIDs(u)
	type row struct {
		Class, Flags                          uint32
		Input                                 [3]uint32
		Z, Vertical, FlagsAfter, Sync, Active uint32
		PlayerSync                            [32]uint32
		Queues                                [3]uint32
		Sounds                                []int
	}
	var rows []row
	for _, cls := range []object.Class{object.ClassSimple, object.ClassPlayer, 0} {
		for _, flags := range []object.Flags{4, 4 | 0x20000, 4 | 0x100000, 4 | 0x800000, 4 | 0x820000, 4 | 0x900000} {
			for _, z := range []float32{-1, 0, 0.00001, 1, 10, 100} {
				for _, v := range []float32{-20, -10, -9, -2, -0.5, 0, 0.5, 3} {
					for _, bounce := range []float32{0, 0.5, 1, 10} {
						o.resetQueues()
						o.s.PortTestCombatAudioReset()
						*u = saved
						worldGeometryResetObject(u, 1001, 100, 100, false)
						u.ObjClass = cls
						u.ObjFlags = flags
						u.ZVal = z
						u.Field27 = v
						u.Field29 = math.Float32bits(bounce)
						u.Field38 = 0x1234
						clear(u.Field140[:])
						u.Field116 = 0x40
						u.Field115 = 0
						u.Collide = o.callback
						u.Update = nil
						legacy.PortTestWorldMotionPhysics("fall", u)
						if flags&0x900000 == 0x100000 && (u.ZVal != z || u.Field27 != v || u.Field38 != 0x1234) {
							t.Fatal("supported object fell")
						}
						if flags&0x900000 == 0 && z == 10 && v == 3 && (u.ZVal != 13 || u.Field27 != 2) {
							t.Fatal("fall gravity contract")
						}
						if flags&0x800000 != 0 && z == 10 && v == 3 && (u.ZVal != 13 || u.Field27 != 2.5) {
							t.Fatal("bounce gravity contract")
						}
						if flags&0x900000 == 0 && z == 1 && v == -20 {
							if u.ZVal != 0 || u.Field27 != 0 || u.ObjFlags&0x20000 != 0 {
								t.Fatal("landing contract")
							}
							if (u.Field116&1 != 0) != (cls&1 == 0) {
								t.Fatal("landing activation eligibility")
							}
						}
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							if e.Obj != u || e.Kind != 0 || e.Code != 0 || e.ByPos {
								t.Fatal("landing audio routing")
							}
							sounds = append(sounds, int(e.ID))
						}
						wantSound := cls == object.ClassPlayer && flags&0x900000 == 0 && (z <= 0 && v < -10 || z > 0 && z+v <= 0 && v-1 < -10)
						if (len(sounds) == 1) != wantSound || len(sounds) > 1 || len(sounds) == 1 && sounds[0] != 280 {
							t.Fatal("landing audio threshold", cls, flags, z, v, sounds)
						}
						rows = append(rows, row{uint32(cls), uint32(flags), [3]uint32{math.Float32bits(z), math.Float32bits(v), math.Float32bits(bounce)}, math.Float32bits(u.ZVal), math.Float32bits(u.Field27), uint32(u.ObjFlags), u.Field38, u.Field116, u.Field140, o.queues(ids), sounds})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-fall", rows, "9639947135da4e35ee81d9b0ef80fb518cecf2d79cb1e8710654e66e71e44903")
}

func TestWorldMotionActivationWrappers(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	ids := collisionCoreIDs(u)
	type row struct {
		Class, Flags uint32
		Callback     bool
		Op           string
		Return       int8
		Active, Link uint32
		Queues       [3]uint32
	}
	var rows []row
	for _, cls := range []object.Class{object.ClassSimple, 0, object.ClassMissile} {
		for _, flags := range []object.Flags{0, 4, 4 | 0x40, 4 | 8} {
			for _, callback := range []bool{false, true} {
				for _, entry := range []string{"activate", "activate-nonsimple"} {
					o.resetQueues()
					*u = saved
					u.ObjClass = cls
					u.ObjFlags = flags
					u.Field115 = 0
					u.Field116 = 0x40
					u.Collide = nil
					u.Update = nil
					if callback {
						u.Collide = o.callback
					}
					for _, op := range []string{entry, entry, "remove-nonsimple", "remove-nonsimple"} {
						rv := legacy.PortTestWorldMotionPhysics(op, u)
						// This wrapper's simple-class return is the address low byte, not a status.
						// Assert the actual value before replacing only that ASLR-dependent field.
						if op == "activate-nonsimple" && cls&1 != 0 {
							if rv != int8(uintptr(u.CObj())) {
								t.Fatal("simple activation return")
							}
							rv = 0
						}
						if op == "activate-nonsimple" && cls&1 != 0 && u.Field116&1 != 0 {
							t.Fatal("simple class was activated")
						}
						if op == "remove-nonsimple" && cls&1 == 0 && u.Field116&1 != 0 {
							t.Fatal("nonsimple removal retained membership")
						}
						rows = append(rows, row{uint32(cls), uint32(flags), callback, op, rv, u.Field116, collisionCoreID(t, ids, u.Field115), o.queues(ids)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-activation-wrappers", rows, "5f7e9e7038d1f9d883f1c976a0bac3ffb126a6baa7459eb09c131c76449694ff")
}
