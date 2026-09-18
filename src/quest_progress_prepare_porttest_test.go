//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestQuestProgressPreparation(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	pairs := questProgressTable(t, o)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"HecubahMarker", "NecromancerMarker"}, nil, true, 0, 0))
	generator := newObjectXferSimple(t, o.s)
	original := *generator
	data := o.record(t, 96)
	generator.UpdateData = data
	generator.ObjClass = object.ClassMonsterGenerator
	donor := newObjectXferSimple(t, o.s)
	donor.TypeInd = uint16(o.s.Types.IndByID(pairs[0][0]))
	var exits []*server.Object
	for i := 0; i < 3; i++ {
		u := newObjectXferSimple(t, o.s)
		u.ObjClass = object.ClassExit
		u.ObjSubClass = 1
		exits = append(exits, u)
	}
	marker := newObjectXferSimple(t, o.s)
	marker.TypeInd = uint16(o.s.Types.IndByID("NecromancerMarker"))
	oldList, oldDeleted, oldPending, oldRNG := o.s.Objs.First(), o.s.Objs.DeletedList, o.s.Objs.Pending, o.s.Rand.Logic
	t.Cleanup(func() {
		*generator = original
		o.s.Objs.SetObjects(oldList)
		o.s.Objs.DeletedList = oldDeleted
		o.s.Objs.Pending = oldPending
		o.s.Rand.Logic = oldRNG
	})
	all := append([]*server.Object{generator}, exits...)
	all = append(all, marker)
	saved := make([]server.Object, len(all))
	for i, u := range all {
		saved[i] = *u
	}
	ids := collisionCoreIDs(all...)
	type row struct {
		Name         string
		Words        [][4]uint32
		Deleted      []uint32
		RNG, Minions uint32
	}
	var rows []row
	for _, stage := range []uint32{0, 4, 5, 6, 7, 9, 10, 0x80000000, 0xffffffff} {
		for level := 0; level < 3; level++ {
			for rate := 0; rate < 5; rate++ {
				for _, present := range []bool{false, true} {
					for exitCount := 0; exitCount <= 3; exitCount++ {
						name := fmt.Sprintf("stage%d-level%d-rate%d-present%v-exits%d", stage, level, rate, present, exitCount)
						for i, u := range all {
							*u = saved[i]
							u.ObjNext = nil
							u.ObjPrev = nil
							u.ObjFlags = 0
						}
						clear(unsafe.Slice((*byte)(data), 96))
						if present {
							*(*unsafe.Pointer)(unsafe.Add(data, 16*level)) = donor.CObj()
						}
						*(*byte)(unsafe.Add(data, 83+level)) = byte(rate)
						*(*byte)(unsafe.Add(data, 87)) = 193
						list := append([]*server.Object{generator}, exits[:exitCount]...)
						list = append(list, marker)
						for i, u := range list {
							if i+1 < len(list) {
								u.ObjNext = list[i+1]
							}
							if i > 0 {
								u.ObjPrev = list[i-1]
							}
						}
						o.s.Objs.SetObjects(generator)
						o.s.Objs.DeletedList = nil
						o.s.Objs.Pending = nil
						o.s.Rand.Logic = prand.New(17 + exitCount)
						*memmap.PtrUint32(0x587000, 202028) = stage
						o.balance(map[string]float64{"QuestHardcoreStage": 9, "MinionsAlwaysStage": 10, "GeneratorMaxActiveCreaturesHigh": 129, "GeneratorMaxActiveCreaturesNormal": 3.75, "GeneratorMaxActiveCreaturesLow": 0, "GeneratorMaxActiveCreaturesSingular": 7})
						legacy.PortTestQuestProgressObject("prepare", nil, level)
						r := row{Name: name, RNG: uint32(o.s.Rand.Logic.Index()), Minions: memmap.Uint32(0x5D4594, 2388656)}
						for _, u := range all {
							r.Words = append(r.Words, [4]uint32{uint32(u.TypeInd), uint32(u.ObjFlags), u.DeletedAt, uint32(*(*byte)(unsafe.Add(data, 87)))})
						}
						for u := o.s.Objs.DeletedList; u != nil; u = u.DeletedNext {
							r.Deleted = append(r.Deleted, ids[u.CObj()])
							if len(r.Deleted) > len(all) {
								t.Fatal("delete cycle")
							}
						}
						if (generator.ObjFlags&object.FlagDestroyed != 0) == present {
							t.Fatal("generator admission", name)
						}
						if marker.ObjFlags&object.FlagDestroyed == 0 {
							t.Fatal("marker cleanup", name)
						}
						if present {
							want := []byte{129, 3, 0, 7, 193}[rate]
							if stage >= 9 && rate != 3 {
								want *= 2
							}
							if *(*byte)(unsafe.Add(data, 87)) != want {
								t.Fatal("generator cap", name)
							}
						}
						kept := 0
						for _, u := range exits[:exitCount] {
							if u.ObjFlags&object.FlagDestroyed == 0 {
								kept++
							}
						}
						if exitCount > 0 && kept != 1 {
							t.Fatal("exactly one exit", name, kept)
						}
						if stage < 5 && r.Minions != 0 {
							t.Fatal("early minions", name)
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	spellbookCapture(t, "quest-progress-preparation", rows, "6e6fe4144c7fbe314bc7c1728faa0cd691eedec9cc2c8e6f42a409b0b0418ec5")
}
