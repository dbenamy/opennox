//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestUpdateRegistryMotionMover(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	oldWPs := o.s.WPs
	var wps []*server.Waypoint
	t.Cleanup(func() {
		o.s.WPs = oldWPs
		for _, w := range wps {
			alloc.FreePtr(unsafe.Pointer(w))
		}
	})
	o.s.WPs.List = nil
	o.s.WPs.Pending = nil
	for _, p := range []types.Pointf{{110, 110}, {160, 100}, {100, 160}} {
		wps = append(wps, o.s.NewWaypoint(p))
	}
	u, v := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	savedU, savedV := *u, *v
	data := collisionCoreGuarded(t, o, 40)
	oldList, oldUp, oldRNG := o.s.Objs.List, o.s.Objs.UpdatableList, o.s.Rand.Logic
	t.Cleanup(func() {
		*u = savedU
		*v = savedV
		o.s.Objs.List = oldList
		o.s.Objs.UpdatableList = oldUp
		o.s.Rand.Logic = oldRNG
	})
	ids := collisionCoreIDs(u, v)
	for i, w := range wps {
		ids[unsafe.Pointer(w)] = uint32(2001 + i)
	}
	type row struct {
		State, Target, Links, Motion       int
		Powered, MissingWaypoint           bool
		Data                               [10]uint32
		Position, TargetPosition, Velocity [2]uint32
		Speed, SpeedBase, Updatable, Head  uint32
		RNG                                int
	}
	var rows []row
	for state := 0; state < 5; state++ {
		for target := 0; target < 6; target++ {
			for _, powered := range []bool{false, true} {
				for _, missing := range []bool{false, true} {
					for links := 0; links < 3; links++ {
						for motion := 0; motion < 12; motion++ {
							// Extra initialization cases check signed speed conversion without
							// multiplying every unrelated state combination.
							if motion >= 5 && (state != 0 || target != 5 || !powered || missing || links != 0) {
								continue
							}
							// State 1 requires a valid current waypoint in ordinary saved state.
							if state == 1 && missing {
								continue
							}
							o.s.PortTestAIEmptyMap()
							o.resetQueues()
							o.s.Rand.Logic = prand.New(23)
							*u = savedU
							*v = savedV
							worldGeometryResetObject(u, 1001, 100, 100, false)
							worldGeometryResetObject(v, 1002, 150, 150, false)
							u.ObjFlags = 4
							v.ObjFlags = 4
							u.Collide = nil
							v.Collide = nil
							u.Update = nil
							v.Update = nil
							u.UpdateData = data
							u.SpeedCur = 2
							u.SpeedBase = 3
							u.VelVec = types.Pointf{1, 1}
							if motion == 1 {
								u.VelVec = types.Pointf{1, 0}
							} else if motion == 2 {
								u.PosVec = types.Pointf{120, 120}
								u.NewPos = u.PosVec
							}
							if motion >= 3 {
								u.PosVec = types.Pointf{100, 110}
								if motion == 4 {
									u.PosVec.X = 110
								}
								u.NewPos = u.PosVec
								u.VelVec = types.Pointf{1, 0}
							}
							if powered {
								u.ObjFlags |= 0x1000000
							}
							v.Extent = 9001
							v.ObjNext = nil
							o.s.Objs.List = v
							u.IsUpdatable = 0
							u.UpdatableNext = nil
							u.UpdatablePrev = nil
							o.s.Objs.UpdatableList = nil
							o.s.Objs.AddToUpdatable(u)
							wps[0].PointsCnt = byte(links)
							wps[0].Points[0].Waypoint = wps[1]
							wps[0].Points[1].Waypoint = wps[2]
							words := (*[10]uint32)(data)
							*words = [10]uint32{}
							words[0] = uint32(state)
							words[1] = 8
							if motion >= 5 {
								words[1] = uint32([]int32{-2147483648, -9, -1, 0, 1, 9, 2147483647}[motion-5])
							}
							words[2] = wps[0].Index
							words[4] = wps[0].Index
							words[6] = wps[1].Index
							words[8] = v.Extent
							if missing {
								words[2] = 9999
								words[4] = 9999
							}
							switch target {
							case 0:
								words[8] = 0
							case 1:
								words[8] = 9999
							case 2:
								v.ObjFlags = 0
							case 3:
								v.ObjFlags = object.FlagDestroyed
							case 5:
								words[7] = uint32(uintptr(v.CObj()))
							}
							legacy.PortTestRegisteredUpdate(u, "MoverUpdate")
							got := *words
							for _, i := range []int{3, 5, 7} {
								got[i] = collisionCoreID(t, ids, got[i])
							}
							if target < 4 || state == 3 {
								if u.IsUpdatable != 0 || o.s.Objs.UpdatableList != nil {
									t.Fatal("invalid/stopped mover remained updatable", state, target)
								}
							}
							if target >= 4 && state == 0 && powered && !missing && (got[0] != 1 || u.PosVec != v.PosVec || u.SpeedCur != float32(float64(int32(words[1]))*.25) || u.SpeedBase != float32(float64(int32(words[1]))*.25)) {
								t.Fatal("mover initialization contract", got, u.PosVec, v.PosVec)
							}
							if target >= 4 && state == 1 && !powered && got[0] != 2 {
								t.Fatal("mover pause transition")
							}
							if target >= 4 && state == 2 && powered && (got[0] != 1 || u.PosVec != v.PosVec) {
								t.Fatal("mover resume transition")
							}
							if target >= 4 && state == 1 && powered && motion == 1 && math.Float32bits(u.VelVec.X) == 0 {
								t.Fatal("axis motion lost initialized velocity")
							}
							if target >= 4 && state == 1 && powered && motion >= 3 {
								if math.Float32bits(u.VelVec.Y) != 0x00800000 || (motion == 4 && math.Float32bits(u.VelVec.X) != 0x00800000) {
									t.Fatal("zero-axis mover velocity must use FLT_MIN", motion, u.VelVec)
								}
							}
							head := uint32(0)
							if o.s.Objs.UpdatableList != nil {
								head = collisionCoreID(t, ids, uint32(uintptr(o.s.Objs.UpdatableList.CObj())))
							}
							rows = append(rows, row{state, target, links, motion, powered, missing, got, motionBits(u.PosVec), motionBits(v.PosVec), motionBits(u.VelVec), math.Float32bits(u.SpeedCur), math.Float32bits(u.SpeedBase), u.IsUpdatable, head, o.s.Rand.Logic.Index()})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "update-registry-motion-mover", rows, updateRegistryHashes["update-registry-motion-mover"])
	updateRegistryReportCounts(t)
}
