//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldMotionScorch(t *testing.T) {
	type row struct {
		Missing, Quest     bool
		Size, Seed         int
		Frame, Tick        uint32
		Ready              uint32
		Types              [3]uint32
		Created            bool
		Position           [2]uint32
		Flags, Expiry, Age uint32
		RNG                int
	}
	var rows []row
	for _, missing := range []bool{false, true} {
		name := "available"
		if missing {
			name = "missing-types"
		}
		t.Run(name, func(t *testing.T) {
			o := newCollisionCoreOwner(t)
			names := []string{"ScorchMarkFloorSmallA", "ScorchMarkFloorMediumA", "ScorchMarkFloorLargeB"}
			var unavailable []string
			if missing {
				unavailable = names
			}
			t.Cleanup(o.s.PortTestRewardTypes(names, unavailable, true, 0, 0))
			globals, restore := legacy.PortTestWorldMotionGlobals()
			t.Cleanup(restore)
			table, restore := legacy.PortTestWorldMotionScorchNames()
			t.Cleanup(restore)
			queues, restore := legacy.PortTestWorldMotionListGlobals()
			t.Cleanup(restore)
			oldRNG, oldPending := o.s.Rand.Logic, o.s.Objs.Pending
			t.Cleanup(func() { o.s.Rand.Logic = oldRNG; o.s.Objs.Pending = oldPending })
			for _, quest := range []bool{false, true} {
				restoreFlags := noxflags.PortTestGameFlags(0)
				t.Cleanup(restoreFlags)
				if quest {
					noxflags.SetGame(noxflags.GameModeQuest)
				}
				for _, size := range []int{-1, 0, 1, 2, 3} {
					for _, seed := range []int{0, 23} {
						for _, frame := range []uint32{123, 0xfffffff0} {
							for _, tick := range []uint32{30, 60} {
								o.s.SetFrame(frame)
								o.s.SetTickRate(tick)
								o.s.Rand.Logic = prand.New(seed)
								o.s.Objs.Pending = nil
								*queues["decay"] = 0
								*globals["scorch-ready"] = 0
								table[1] = 0
								table[3] = 0
								table[5] = 0
								pos := types.Pointf{100.25, 200.75}
								legacy.PortTestWorldMotionScorch(&pos, int32(size))
								r := row{Missing: missing, Quest: quest, Size: size, Seed: seed, Frame: frame, Tick: tick, Ready: *globals["scorch-ready"], Types: [3]uint32{table[1], table[3], table[5]}, RNG: o.s.Rand.Logic.Index()}
								if r.Ready != 1 {
									t.Fatal("scorch initializes types even for invalid sizes")
								}
								for i, n := range names {
									want := uint32(o.s.Types.IndByID(n))
									if r.Types[i] != want || !missing && want == 0 {
										t.Fatal("scorch named type initialization")
									}
								}
								u := o.s.Objs.Pending
								r.Created = u != nil
								if r.Created != (size >= 0 && size <= 2 && !missing) {
									t.Fatal("scorch size/allocation contract", size, missing)
								}
								if u != nil {
									trackObjectXferTyped(t, o.s, u)
									if u.ObjNext != nil || u.ObjOwner != nil || u.PosVec != pos || uint32(u.TypeInd) != r.Types[size] || *queues["decay"] != uint32(uintptr(u.CObj())) || u.ObjFlags&0x400000 == 0 {
										t.Fatal("scorch actual creation/decay owner")
									}
									life := (u.Field34 - frame) / tick
									min, max := uint32(10), uint32(20)
									if quest {
										min, max = 5, 8
									}
									if life < min || life > max || u.Field34-frame != life*tick {
										t.Fatal("scorch lifetime range and frame wrap", life)
									}
									r.Position = motionBits(u.PosVec)
									r.Flags = uint32(u.ObjFlags)
									r.Expiry = u.Field34
									r.Age = u.Field32
								}
								wantRNG := seed
								if size >= 0 && size <= 2 {
									wantRNG++
									if !missing {
										wantRNG++
									}
								}
								if r.RNG != wantRNG {
									t.Fatal("scorch RNG consumption", r.RNG, wantRNG)
								}
								rows = append(rows, r)
							}
						}
					}
				}
				restoreFlags()
			}
			t.Cleanup(func() { o.s.Objs.Pending = nil; *queues["decay"] = 0 })
		})
	}
	spellbookCapture(t, "world-motion-scorch", rows, "")
}
