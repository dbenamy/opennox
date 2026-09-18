//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestQuestProgressBosses(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(legacy.PortTestQuestProgressOwner())
	template := newCreatureXferObject(t, o.s, "Monster")
	clear(unsafe.Slice((*byte)(template.UpdateData), 2200))
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Hecubah", "Necromancer", "RewardMarker", "RewardMarkerPlus", "QuestGoldPile", "QuestGoldChest", "RubyGem", "EmeraldGem", "DiamondGem"}, nil, true, 0, 0))
	t.Cleanup(o.s.PortTestSpellEffectTypes(nil, []string{"Hecubah", "Necromancer"}, template))
	// Controlled, nonzero reward category and gold ranges; real reward selection runs.
	raw := serverConfigOwnBytes(t, 0x587000, 207044, 64)
	clear(raw)
	*memmap.PtrUint32(0x587000, 207044+5*8) = 9
	serverConfigOwnBytes(t, 0x587000, 211136, 40)
	for i := uintptr(0); i < 5; i++ {
		*memmap.PtrUint32(0x587000, 211136+8*i) = uint32(10 + 10*i)
		*memmap.PtrUint32(0x587000, 211140+8*i) = uint32(20 + 10*i)
	}
	serverConfigOwnBytes(t, 0x5D4594, 1568276, 4)
	*memmap.PtrUint32(0x5D4594, 1568276) = 0
	serverConfigOwnBytes(t, 0x587000, 202036, 4)
	oldPending, oldRNG := o.s.Objs.Pending, o.s.Rand.Logic
	t.Cleanup(func() { o.s.Objs.Pending = oldPending; o.s.Rand.Logic = oldRNG })
	at := newObjectXferSimple(t, o.s)
	at.PosVec = types.Pointf{X: 123.25, Y: 177.5}
	definition := o.record(t, 248)
	objectXferSetWord(definition, 72, 93)
	markerType := o.s.Types.ByID("RewardMarker")
	type row struct {
		Name             string
		Present          bool
		Health           [3]uint16
		Words            [550]uint32
		Flags, Type, RNG uint32
		Position         types.Pointf
		Inventory        [][3]uint32
	}
	var rows []row
	for _, boss := range []string{"Hecubah", "Necromancer"} {
		typ := o.s.Types.ByID(boss)
		op := "hecubah"
		if boss == "Necromancer" {
			op = "necro"
		}
		for _, scale := range []float32{0, .5, 1, 1.3333334, 2, 655.36} {
			for _, hp := range []uint16{0, 100, 65535} {
				for _, override := range []int64{-1, 0, 93, 0xffffffff, 0x7fffffff} {
					for rewardMode := 0; rewardMode < 3; rewardMode++ {
						name := fmt.Sprintf("%s-scale%08x-hp%d-override%v-reward%d", boss, math.Float32bits(scale), hp, override, rewardMode)
						*typ.Health() = server.HealthData{Cur: hp, Field2: 0xabcd, Max: hp}
						*(*unsafe.Pointer)(unsafe.Add(typ.UpdateData, 484)) = nil
						if override >= 0 {
							objectXferSetWord(definition, 72, uint32(override))
							*(*unsafe.Pointer)(unsafe.Add(typ.UpdateData, 484)) = definition
						}
						*memmap.PtrFloat32(0x587000, 202036) = scale
						*memmap.PtrUint32(0x587000, 202028) = 7
						objectXferSetWord(markerType.InitData, 0, 0)
						if rewardMode == 2 {
							objectXferSetWord(markerType.InitData, 0, 32)
						}
						restore := func() {}
						if rewardMode == 0 {
							restore = o.s.PortTestRewardTypes(nil, []string{"RewardMarker"}, true, 0, 0)
						}
						o.s.Objs.Pending = nil
						o.s.Rand.Logic = prand.New(23)
						legacy.PortTestQuestProgressObject(op, at, 0)
						restore()
						u := o.s.Objs.Pending
						if u == nil || u.ObjNext != nil {
							t.Fatal("single staged boss", name)
						}
						r := row{Name: name, Present: true, Health: [3]uint16{u.HealthData.Cur, u.HealthData.Field2, u.HealthData.Max}, Words: *(*[550]uint32)(u.UpdateData), Flags: uint32(u.ObjFlags), Type: uint32(u.TypeInd), RNG: uint32(o.s.Rand.Logic.Index()), Position: u.PosVec}
						if r.Words[121] != 0 {
							r.Words[121] = 1
						}
						if r.Health[0] == 0 || r.Health[2] == 0 || r.Health[1] != 0xabcd || r.Position != at.PosVec {
							t.Fatal("spawn health/position", name, r.Health)
						}
						if r.Words[340] != 4 || r.Words[411] != 0x10000000 || r.Words[423] != 0x10000000 || r.Words[326] != 1062501089 || r.Words[410] != 0x8000000 || r.Words[444] != 0x20000000 || r.Words[415] != 0x40000000 {
							t.Fatal("boss attributes", name)
						}
						wantItems := 0
						if rewardMode == 2 {
							wantItems = 1
							if boss == "Hecubah" {
								wantItems = 4
							}
						}
						for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
							if it.InvHolder != u {
								t.Fatal("reward ownership")
							}
							r.Inventory = append(r.Inventory, [3]uint32{uint32(it.TypeInd), objectXferGetWord(it.InitData, 0), uint32(it.ObjFlags)})
						}
						if len(r.Inventory) != wantItems {
							t.Fatal("reward count", name, len(r.Inventory), wantItems)
						}
						rows = append(rows, r)
						// Detach lists, then release every owned allocation before returning records.
						for it := u.InvFirstItem; it != nil; {
							next := it.InvNextItem
							it.InvNextItem = nil
							it.InvHolder = nil
							questProgressFree(o.s, it)
							it = next
						}
						u.InvFirstItem = nil
						o.s.Objs.Pending = nil
						questProgressFree(o.s, u)
					}
				}
			}
		}
		restore := o.s.PortTestRewardTypes(nil, []string{boss}, true, 0, 0)
		o.s.Objs.Pending = nil
		legacy.PortTestQuestProgressObject(op, at, 0)
		restore()
		if o.s.Objs.Pending != nil {
			t.Fatal("missing boss type")
		}
		rows = append(rows, row{Name: boss + "-missing"})
	}
	spellbookCapture(t, "quest-progress-bosses", rows, "eb3bee27eeb2a93a273c9b45f64a98a176411a2656ec6ff4441512a06d4d2a33")
}
func questProgressFree(s *server.Server, u *server.Object) {
	if u.HealthData != nil {
		alloc.Free(u.HealthData)
		u.HealthData = nil
	}
	for _, p := range []*unsafe.Pointer{&u.InitData, &u.UpdateData, &u.CollideData, &u.Field189} {
		if *p != nil {
			alloc.FreePtr(*p)
			*p = nil
		}
	}
	s.Objs.FreeObject(u)
}
