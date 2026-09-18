//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerOrchestrationRewards(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestOrchestrationRewardTypes())
	globals, restore := legacy.PortTestServerOrchestrationGlobals()
	t.Cleanup(restore)
	for _, off := range []uintptr{1568276, 1568284, 1568292, 1568296, 1568304} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
		*memmap.PtrUint32(0x5D4594, off) = 0
	}
	serverConfigOwnBytes(t, 0x587000, 207044, 4132)
	clear(unsafe.Slice(memmap.PtrUint8(0x587000, 207044), 4132))
	for i := 0; i < 8; i++ {
		*memmap.PtrUint32(0x587000, 207044+uintptr(i*8)) = 1
	}
	name := o.record(t, 32)
	copy(unsafe.Slice((*byte)(name), 32), "PortTestRewardWeapon")
	for off, v := range map[uintptr]uint32{208176: 1, 208184: uint32(o.s.Types.IndByID("PortTestRewardWeapon")), 208188: 1, 208192: 31} {
		*memmap.PtrUint32(0x587000, off) = v
	}
	*memmap.PtrPtr(0x587000, 208180) = name
	for i := 0; i < 5; i++ {
		*memmap.PtrUint32(0x587000, 211136+uintptr(8*i)) = uint32(100 + 100*i)
		*memmap.PtrUint32(0x587000, 211140+uintptr(8*i)) = uint32(100 + 100*i)
	}
	oldList, oldPending, oldDeleted, oldRNG := o.s.Objs.List, o.s.Objs.Pending, o.s.Objs.DeletedList, o.s.Rand.Logic
	t.Cleanup(func() {
		o.s.Objs.List = oldList
		o.s.Objs.Pending = oldPending
		o.s.Objs.DeletedList = oldDeleted
		o.s.Rand.Logic = oldRNG
	})
	type item struct {
		Name string
		Gold uint32
		Pos  types.Pointf
	}
	type row struct {
		Stage        uint32
		Seed, Kind   int
		Active       [3]bool
		World, Chest []item
		Deleted      []int
		RNG          int
	}
	var rows []row
	for _, stage := range []uint32{0, 1, 2, 8, 9, 10, 0xffffffff} {
		for _, seed := range []int{1, 23, 253} {
			for kind := 0; kind < 5; kind++ {
				t.Run(fmt.Sprintf("stage%d/seed%d/kind%d", stage, seed, kind), func(t *testing.T) {
					o.reset()
					if kind == 4 {
						defer o.s.PortTestRewardTypes(nil, []string{"Quiver"}, true, 0, 0)()
					}
					o.s.Rand.Logic = prand.New(seed)
					*o.quest["202028"] = stage
					for i := range o.units {
						objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4792, 1)
					}
					for _, k := range []string{"reward-marker", "ankh-marker", "selected-marker"} {
						*globals[k] = 0
					}
					var objects []*server.Object
					makeObj := func(name string) *server.Object {
						u := o.s.NewObjectByTypeID(name)
						if u == nil {
							t.Fatal(name)
						}
						u.ObjFlags = 0
						objects = append(objects, u)
						return u
					}
					a, b, c, chest, inside, plus, other := makeObj("RewardMarker"), makeObj("RewardMarkerPlus"), makeObj("RewardMarker"), makeObj("RewardChest"), makeObj("RewardMarker"), makeObj("RewardMarkerPlus"), makeObj("RewardOther")
					o.s.Objs.List = a
					o.s.Objs.Pending = nil
					o.s.Objs.DeletedList = nil
					a.ObjNext = b
					b.ObjNext = c
					c.ObjNext = chest
					chest.ObjNext = nil
					chest.Init = legacy.PortTestServerOrchestrationChestInit()
					chest.ObjClass = 0
					chest.InvFirstItem = inside
					inside.InvNextItem = other
					other.Field125 = inside
					other.InvNextItem = plus
					plus.Field125 = other
					for _, u := range []*server.Object{inside, other, plus} {
						u.InvHolder = chest
					}
					for i, u := range []*server.Object{a, b, c, inside, plus} {
						u.PosVec = types.Pointf{X: float32(100 + i*40), Y: 200}
						mask := uint32(32)
						if kind == 1 || kind == 4 {
							mask = 8
						} else if kind == 2 {
							mask = 0
						} else if kind == 3 {
							mask = 128
						}
						objectXferSetWord(u.InitData, 0, mask)
						objectXferSetWord(u.InitData, 216, 0)
					}
					objectXferSetWord(c.InitData, 216, 1)
					defer func() {
						generated := map[*server.Object]bool{}
						for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
							generated[u] = true
						}
						for u := chest.InvFirstItem; u != nil; u = u.InvNextItem {
							if u != other && u != inside && u != plus {
								generated[u] = true
							}
						}
						o.s.Objs.List = nil
						o.s.Objs.Pending = nil
						o.s.Objs.DeletedList = nil
						chest.InvFirstItem = nil
						for u := range generated {
							o.s.Objs.FreeObject(u)
						}
						for _, u := range objects {
							o.s.Objs.FreeObject(u)
						}
					}()
					legacy.PortTestServerOrchestration("rewards", nil, 0)
					r := row{Stage: stage, Seed: seed, Kind: kind, RNG: o.s.Rand.Logic.Index()}
					active := 0
					for i, u := range []*server.Object{a, b, c} {
						r.Active[i] = *(*byte)(unsafe.Add(u.InitData, 216))&128 != 0
						if r.Active[i] {
							active++
						}
						if !u.ObjFlags.Has(object.FlagDestroyed) {
							t.Fatal("world marker retained")
						}
					}
					for _, u := range []*server.Object{inside, plus} {
						if !u.ObjFlags.Has(object.FlagDestroyed) || u.InvHolder != nil {
							t.Fatal("chest marker retained")
						}
					}
					if other.InvHolder != chest || other.ObjFlags.Has(object.FlagDestroyed) {
						t.Fatal("ordinary inventory changed")
					}
					capture := func(u *server.Object) item {
						it := item{Name: o.s.Types.ByInd(int(u.TypeInd)).ID(), Pos: u.PosVec}
						if it.Name == "questgoldpile" || it.Name == "questgoldchest" {
							it.Gold = objectXferGetWord(u.InitData, 0)
							if it.Gold < 100 || it.Gold > 500 || it.Gold%100 != 0 {
								t.Fatal("gold tier", it)
							}
						}
						return it
					}
					for u := o.s.Objs.Pending; u != nil; u = u.ObjNext {
						if len(r.World) > 10 {
							t.Fatal("pending cycle")
						}
						r.World = append(r.World, capture(u))
					}
					for u := chest.InvFirstItem; u != nil; u = u.InvNextItem {
						if len(r.Chest) > 6 {
							t.Fatal("inventory cycle")
						}
						r.Chest = append(r.Chest, capture(u))
						if u.InvHolder != chest {
							t.Fatal("generated inventory owner")
						}
					}
					multiplier := 1
					if kind == 1 {
						multiplier = 2
					}
					wantWorld, wantChest := active*multiplier, 1+2*multiplier
					if kind == 3 {
						wantWorld++
					}
					if kind == 2 {
						wantWorld, wantChest = 0, 1
					}
					if len(r.World) != wantWorld || len(r.Chest) != wantChest {
						t.Fatal("generated item count", len(r.World), len(r.Chest), wantWorld, wantChest)
					}
					for u := o.s.Objs.DeletedList; u != nil; u = u.DeletedNext {
						idx := -1
						for i, v := range objects {
							if u == v {
								idx = i
							}
						}
						if idx < 0 || len(r.Deleted) > 5 {
							t.Fatal("deleted queue")
						}
						r.Deleted = append(r.Deleted, idx)
					}
					if len(r.Deleted) != 5 {
						t.Fatal("deleted marker count", r.Deleted)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "server-orchestration-rewards", rows, "17d07dde8a77dc78ee93b67c2c6823d9144128565e3870ebb02d88312f749ad7")
}
