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

func TestWorldMotionVelocitySprings(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	_, restore := legacy.PortTestWorldMotionVelocityGlobals()
	t.Cleanup(restore)
	oldSprings := noxServer.springs
	t.Cleanup(func() { noxServer.springs = oldSprings })
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	type row struct {
		Distance, Step             uint32
		Frozen, Buff, SecondActive bool
		A, B                       worldGeometryObjectState
		Queues                     [3]uint32
		Spring                     bool
		Length, Previous           uint32
	}
	var rows []row
	for _, distance := range []float32{20, 40, 400} {
		for _, step := range []float32{0, 0.01, 0.1, 1} {
			for _, frozen := range []bool{false, true} {
				for _, buff := range []bool{false, true} {
					for _, secondActive := range []bool{false, true} {
						o.s.PortTestAIEmptyMap()
						o.resetHits()
						o.resetQueues()
						noxServer.springs = serverSprings{}
						*a = sa
						*b = sb
						worldGeometryResetObject(a, 1001, 100, 100, false)
						worldGeometryResetObject(b, 1002, 120, 100, false)
						for _, u := range []*server.Object{a, b} {
							u.Shape.Kind = server.ShapeKindCenter
							u.ObjFlags = 4
							u.Update = o.callback
							u.Collide = o.callback
							u.Field115 = 0
							u.Field116 = 0
							u.HealthData = nil
							u.ForceVec = types.Pointf{}
							u.Float28 = 0
							u.Buffs = 0
						}
						if frozen {
							a.ObjFlags |= 2
						}
						if buff {
							a.Buffs = 1 << 5
						}
						noxServer.springs.Add(a, b)
						b.PosVec = types.Pointf{100 + distance, 100}
						b.NewPos = b.PosVec
						legacy.PortTestWorldMotionPhysics("activate", a)
						if secondActive {
							legacy.PortTestWorldMotionPhysics("activate", b)
						}
						legacy.PortTestWorldMotionVelocity(step)
						spring := noxServer.springs.head
						if (spring != nil) != (distance <= 256) {
							t.Fatal("spring distance removal")
						}
						if distance == 40 && step == 0.01 {
							if b.VelVec != (types.Pointf{-6, 0}) {
								t.Fatal("spring must activate and accelerate second object", b.VelVec)
							}
							want := types.Pointf{6, 0}
							if frozen {
								want = types.Pointf{}
							}
							if a.VelVec != want {
								t.Fatal("spring force survives immobilizing buff but frozen unit does not move", a.VelVec, want)
							}
						}
						r := row{Distance: math.Float32bits(distance), Step: math.Float32bits(step), Frozen: frozen, Buff: buff, SecondActive: secondActive, A: worldGeometryState(a), B: worldGeometryState(b), Queues: o.queues(ids), Spring: spring != nil}
						if spring != nil {
							r.Length = math.Float32bits(spring.curLen)
							r.Previous = math.Float32bits(spring.prevLen)
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-velocity-springs", rows, "9df1b1d145760f4769076502c458b3ed614cc08a56d4bd139704d5a025d4adee")
}

func TestWorldMotionVelocityFloorAndMonster(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, restore := legacy.PortTestWorldMotionTileGrid()
	t.Cleanup(restore)
	_, restore = legacy.PortTestWorldMotionVelocityGlobals()
	t.Cleanup(restore)
	oldSprings := noxServer.springs
	noxServer.springs = serverSprings{}
	t.Cleanup(func() { noxServer.springs = oldSprings })
	u := newCreatureXferObject(t, o.s, "Monster")
	saved := *u
	ud := u.UpdateDataMonster()
	oldUD := *ud
	t.Cleanup(func() { *u = saved; *ud = oldUD })
	ids := collisionCoreIDs(u)
	type row struct {
		Tile                    int32
		Flags                   uint32
		Monster, Action, Health bool
		Object                  worldGeometryObjectState
		Hits                    [][5]uint32
		Sync                    uint32
	}
	var rows []row
	for _, tile := range []int32{0, 5, 6, 7} {
		for _, flags := range []object.Flags{4, 4 | 2, 4 | 0x4000} {
			for _, monster := range []bool{false, true} {
				for _, action := range []bool{false, true} {
					for _, health := range []bool{false, true} {
						configure(tile)
						o.s.PortTestAIEmptyMap()
						o.resetHits()
						o.resetQueues()
						*u = saved
						*ud = oldUD
						worldGeometryResetObject(u, 1001, 100, 100, false)
						u.Shape.Kind = server.ShapeKindCenter
						u.ObjClass = object.ClassSimple
						if monster {
							u.ObjClass = object.ClassMonster
						}
						u.ObjFlags = flags
						u.HealthData = nil
						if health {
							u.HealthData = saved.HealthData
						}
						ud.AIStackInd = 0
						ud.AIStack[0].Action = 0
						if action {
							ud.AIStack[0].Action = 67
						}
						u.Field115 = 0
						u.Field116 = 0
						u.Update = o.callback
						u.Collide = o.callback
						u.VelVec = types.Pointf{1, 2}
						u.ForceVec = types.Pointf{3, 4}
						u.Float28 = 0
						u.Field38 = 0
						u.Buffs = 0
						legacy.PortTestWorldMotionPhysics("activate", u)
						legacy.PortTestWorldMotionVelocity(1)
						hits := o.hits(ids)
						frozen := flags&2 != 0 || monster && action
						wantHit := !frozen && flags&0x4000 == 0 && health && tile == 6
						if (len(hits) == 1) != wantHit || len(hits) > 1 {
							t.Fatal("hazard-floor eligibility", tile, flags, monster, action, health, hits)
						}
						if wantHit && hits[0] != ([5]uint32{1001, 6, 0, 0, 1001 % 256}) {
							t.Fatal("hazard-floor contact representation", hits)
						}
						if frozen && (u.VelVec != (types.Pointf{}) || u.NewPos != (types.Pointf{100, 100}) || u.Field38 != 0) {
							t.Fatal("monster action/freeze motion suppression")
						}
						if !frozen && (u.VelVec != (types.Pointf{4, 6}) || u.NewPos != (types.Pointf{104, 106})) {
							t.Fatal("floor test must exercise nonzero motion")
						}
						if !unchanged() {
							t.Fatal("velocity mutated tile grid")
						}
						rows = append(rows, row{tile, uint32(flags), monster, action, health, worldGeometryState(u), hits, u.Field38})
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-velocity-floor-monster", rows, "ae49812caa6a3b740c0f09b5321e9727fb67e934be678940f560c50c773d9041")
}

func TestWorldMotionVelocitySyncThreshold(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	_, restore := legacy.PortTestWorldMotionVelocityGlobals()
	t.Cleanup(restore)
	oldSprings := noxServer.springs
	noxServer.springs = serverSprings{}
	t.Cleanup(func() { noxServer.springs = oldSprings })
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	type row struct {
		X, Y        uint32
		Sync, Flags uint32
		Min, Max    [2]uint32
	}
	var rows []row
	for _, x := range []float32{100, 100.00999, 100.01, 100.01001} {
		for _, y := range []float32{100, 100.00999, 100.01, 100.01001} {
			o.s.PortTestAIEmptyMap()
			o.resetQueues()
			*u = saved
			worldGeometryResetObject(u, 1001, 100, 100, false)
			u.Shape.Kind = server.ShapeKindCenter
			u.NewPos = types.Pointf{x, y}
			u.ObjFlags = 4
			u.Collide = o.callback
			u.Update = o.callback
			u.HealthData = nil
			u.Field38 = 0
			u.Field115 = 0
			u.Field116 = 0
			u.Float28 = 0
			legacy.PortTestWorldMotionPhysics("activate", u)
			legacy.PortTestWorldMotionVelocity(0)
			want := uint32(0)
			if float64(x)-100 > 0.0099999998 || float64(y-100) > 0.0099999998 {
				want = 0xffffffff
			}
			if u.Field38 != want {
				t.Fatal("movement sync threshold", x, y, u.Field38, want)
			}
			rows = append(rows, row{math.Float32bits(x), math.Float32bits(y), u.Field38, uint32(u.ObjFlags), motionBits(u.CollideP1), motionBits(u.CollideP2)})
		}
	}
	spellbookCapture(t, "world-motion-velocity-sync-threshold", rows, "1e9bf7ba6e394e0a1e633d3000355ca5633a18fd8a665116b58f0791bea8b3d6")
}
