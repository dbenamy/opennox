//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestWorldMotionVelocity(t *testing.T) {
	o := newCollisionCoreOwner(t)
	names := []string{"Trigger", "BlackPowder", "TelekinesisHand", "SmallFist", "MediumFist", "LargeFist", "Meteor", "Spike", "PeriodicSpike", "SmallFlameCleanse", "MediumFlameCleanse", "FlameCleanse", "LargeFlameCleanse", "SmallBlueFlameCleanse", "MediumBlueFlameCleanse", "BlueFlameCleanse", "LargeBlueFlameCleanse"}
	t.Cleanup(o.s.PortTestRewardTypes(names, nil, true, 0, 0))
	caches, restore := legacy.PortTestWorldMotionVelocityGlobals()
	t.Cleanup(restore)
	oldSprings := noxServer.springs
	noxServer.springs = serverSprings{}
	t.Cleanup(func() { noxServer.springs = oldSprings })
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	ids := collisionCoreIDs(u)
	type row struct {
		Wall, Mode                                                  int
		Step, Drag                                                  uint32
		Position, Proposed, Velocity, Acceleration, Force, Min, Max [2]uint32
		Flags, Sync                                                 uint32
		Cache                                                       [10]uint32
		Hits                                                        [][5]uint32
		Queues                                                      [3]uint32
	}
	var rows []row
	for wall := 0; wall < 3; wall++ {
		for mode := 0; mode < 8; mode++ {
			for _, step := range []float32{0, 0.001, 0.5, 1} {
				for _, drag := range []float32{0, 0.1, 1} {
					o.s.PortTestAIEmptyMap()
					configure(wall)
					o.resetHits()
					o.resetQueues()
					*caches = [10]uint32{}
					*u = saved
					worldGeometryResetObject(u, 1001, 100, 100, false)
					u.Shape.Kind = server.ShapeKindCenter
					u.HealthData = nil
					u.ObjFlags = 4
					u.Update = o.callback
					u.Collide = o.callback
					u.Field115 = 0
					u.Field116 = 0x40
					u.VelVec = types.Pointf{100, 0}
					u.ForceVec = types.Pointf{3, 4}
					u.Pos24 = types.Pointf{900, 901}
					u.Float28 = drag
					u.Field38 = 0
					u.Buffs = 0
					switch mode {
					case 1:
						u.ObjFlags |= 2
					case 2:
						u.Buffs = 1 << 5
					case 3:
						u.Buffs = 1 << 25
					case 4:
						u.Buffs = 1 << 28
					case 5:
						u.ObjFlags |= 0x4000
					case 6:
						u.TypeInd = uint16(o.s.Types.IndByID("SmallFlameCleanse"))
					case 7:
						u.TypeInd = uint16(o.s.Types.IndByID("LargeBlueFlameCleanse"))
					}
					legacy.PortTestWorldMotionPhysics("activate", u)
					if rv := legacy.PortTestWorldMotionVelocity(step); rv != 0 {
						t.Fatal("velocity traversal return", rv)
					}
					if !unchanged() {
						t.Fatal("velocity changed wall definition")
					}
					if u.Pos24 != (types.Pointf{}) {
						t.Fatal("collision scan must clear stale acceleration")
					}
					if mode == 1 && (u.VelVec != (types.Pointf{}) || u.NewPos != (types.Pointf{100, 100})) {
						t.Fatal("frozen object moved")
					}
					if wall == 0 && mode == 0 && step == 1 && drag == 0 && (u.VelVec != (types.Pointf{103, 4}) || u.NewPos != (types.Pointf{203, 104})) {
						t.Fatal("velocity force integration")
					}
					if mode >= 2 && mode <= 4 && drag == 0 && u.VelVec != (types.Pointf{100, 0}) {
						t.Fatal("immobilizing buff did not suppress applied force", mode, u.VelVec)
					}
					var want [10]uint32
					for i, name := range []string{"SmallFlameCleanse", "SmallFlameCleanse", "MediumFlameCleanse", "FlameCleanse", "LargeFlameCleanse", "SmallBlueFlameCleanse", "SmallBlueFlameCleanse", "MediumBlueFlameCleanse", "BlueFlameCleanse", "LargeBlueFlameCleanse"} {
						want[i] = uint32(o.s.Types.IndByID(name))
						if want[i] == 0 {
							t.Fatal("missing fixture type", name)
						}
					}
					if *caches != want {
						t.Fatal("velocity type cache must use actual named definitions", *caches)
					}
					rows = append(rows, row{wall, mode, math.Float32bits(step), math.Float32bits(drag), motionBits(u.PosVec), motionBits(u.NewPos), motionBits(u.VelVec), motionBits(u.Pos24), motionBits(u.ForceVec), motionBits(u.CollideP1), motionBits(u.CollideP2), uint32(u.ObjFlags), u.Field38, *caches, o.hits(ids), o.queues(ids)})
				}
			}
		}
	}
	spellbookCapture(t, "world-motion-velocity", rows, "4c75c39ed3e8685d10a1f6d95fa3e5363e77fff909ec3cd0a951dc1920482284")
}
