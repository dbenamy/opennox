//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestQuestProgressBossMarkers(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	template := newCreatureXferObject(t, o.s, "Monster")
	clear(unsafe.Slice((*byte)(template.UpdateData), 2200))
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Hecubah", "Necromancer", "HecubahMarker", "NecromancerMarker"}, nil, true, 0, 0))
	t.Cleanup(o.s.PortTestSpellEffectTypes(nil, []string{"Hecubah", "Necromancer"}, template))
	serverConfigOwnBytes(t, 0x587000, 202036, 4)
	*memmap.PtrFloat32(0x587000, 202036) = 1.5
	var markers []*server.Object
	for i := 0; i < 4; i++ {
		u := newObjectXferSimple(t, o.s)
		name := "HecubahMarker"
		if i >= 2 {
			name = "NecromancerMarker"
		}
		u.TypeInd = uint16(o.s.Types.IndByID(name))
		u.PosVec = types.Pointf{X: float32(100 + 30*i), Y: 100}
		markers = append(markers, u)
	}
	oldList, oldDeleted, oldPending, oldRNG := o.s.Objs.First(), o.s.Objs.DeletedList, o.s.Objs.Pending, o.s.Rand.Logic
	t.Cleanup(func() {
		o.s.Objs.SetObjects(oldList)
		o.s.Objs.DeletedList = oldDeleted
		o.s.Objs.Pending = oldPending
		o.s.Rand.Logic = oldRNG
	})
	ids := collisionCoreIDs(markers...)
	type spawn struct {
		Type   uint16
		Pos    types.Pointf
		Health [3]uint16
	}
	type row struct {
		Name         string
		Minions, RNG uint32
		Deleted      []uint32
		Spawns       []spawn
	}
	var rows []row
	for _, stage := range []uint32{0, 4, 5, 6, 7, 8, 9, 10, 11, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, seed := range []int{0, 1, 23, 83, 131} {
			for _, hecCount := range []int{0, 1, 2} {
				name := fmt.Sprintf("stage%d-seed%d-hec%d", stage, seed, hecCount)
				list := append([]*server.Object(nil), markers[:hecCount]...)
				list = append(list, markers[2:]...)
				for i, u := range list {
					u.ObjFlags = 0
					u.ObjNext = nil
					u.ObjPrev = nil
					u.DeletedNext = nil
					if i+1 < len(list) {
						u.ObjNext = list[i+1]
					}
					if i > 0 {
						u.ObjPrev = list[i-1]
					}
				}
				o.s.Objs.SetObjects(list[0])
				o.s.Objs.Pending = nil
				o.s.Objs.DeletedList = nil
				o.s.Rand.Logic = prand.New(seed)
				*memmap.PtrUint32(0x587000, 202028) = stage
				o.balance(map[string]float64{"QuestHardcoreStage": 9, "MinionsAlwaysStage": 10})
				legacy.PortTestQuestProgressObject("prepare", nil, 0)
				r := row{Name: name, Minions: memmap.Uint32(0x5D4594, 2388656), RNG: uint32(o.s.Rand.Logic.Index())}
				for u := o.s.Objs.DeletedList; u != nil; u = u.DeletedNext {
					r.Deleted = append(r.Deleted, ids[u.CObj()])
					if len(r.Deleted) > len(list) {
						t.Fatal("marker delete cycle")
					}
				}
				if len(r.Deleted) != len(list) {
					t.Fatal("every marker deleted", name)
				}
				for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
					r.Spawns = append(r.Spawns, spawn{u.TypeInd, u.PosVec, [3]uint16{u.HealthData.Cur, u.HealthData.Field2, u.HealthData.Max}})
					if len(r.Spawns) > 3 {
						t.Fatal("too many bosses")
					}
				}
				// Separate RNG contract predicts admission and which authored positions spawn.
				rng := prand.New(seed)
				admit := false
				if int32(stage) >= 5 {
					admit = stage == 5 || int32(stage) >= 10
					if !admit && stage&1 != 0 {
						admit = rng.IntClamp(1, 100) >= 50
					}
				}
				var want []types.Pointf
				if admit && hecCount > 0 {
					chosen := rng.IntClamp(1, hecCount) - 1
					want = append(want, markers[chosen].PosVec)
					for _, u := range markers[2:] {
						if rng.IntClamp(1, 100) >= 50 {
							want = append(want, u.PosVec)
						}
					}
				}
				if (r.Minions != 0) != admit || len(r.Spawns) != len(want) || r.RNG != uint32(rng.Index()) {
					t.Fatal("admission/spawn/RNG contract", name, r.Minions, len(r.Spawns), len(want), r.RNG, rng.Index())
				}
				for i, p := range want {
					if r.Spawns[len(want)-1-i].Pos != p {
						t.Fatal("selected marker", name)
					}
				}
				rows = append(rows, r)
				for u := o.s.Objs.Pending; u != nil; {
					next := u.ObjNext
					u.ObjNext = nil
					questProgressFree(o.s, u)
					u = next
				}
				o.s.Objs.Pending = nil
			}
		}
	}
	spellbookCapture(t, "quest-progress-markers", rows, "050cb42e87a2715242feba282562227677fcbf1a9adac8e1c66487a06585c116")
}
