//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestCollisionRegistryWorldChestContents(t *testing.T) {
	o := newWorldCollisionOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	defer noxflags.PortTestGameFlags(noxflags.GameModeCoop)()
	for _, off := range []uintptr{1568252, 1568256, 1568244} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		t.Cleanup(func() { *p = old })
		*p = 0
	}
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	other := newObjectXferSimple(t, o.s)
	a.UpdateData = o.record(t, 8)
	a.PosVec = types.Ptf(100, 100)
	b.PosVec = types.Ptf(120, 110)
	a.Shape.Kind = server.ShapeKindCircle
	a.Shape.Circle.R = 10
	a.Shape.Circle.R2 = 100
	a.Death = server.PortTestWorldChestDeath()
	var items []*server.Object
	for i := 0; i < 4; i++ {
		it := newObjectXferSimple(t, o.s)
		it.NetCode = uint32(2001 + i)
		items = append(items, it)
	}
	t.Cleanup(func() {
		a.UpdateData = nil
		a.InvFirstItem = nil
		o.s.Objs.Pending = nil
		for _, it := range items {
			it.InvHolder = nil
			it.InvNextItem = nil
			it.Field125 = nil
			it.ObjNext = nil
			it.ObjPrev = nil
		}
	})
	var rows []struct {
		Name      string
		Opened    uint32
		Pending   []uint32
		Positions [][2]uint32
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-chest-contents", rows)
	}()
	for _, count := range []int{0, 1, 3, 4} {
		for _, subclass := range []object.SubClass{0, 0x100, 0x200, 0x400, 0x800} {
			for _, dead := range []bool{false, true} {
				for _, player := range []bool{false, true} {
					name := fmt.Sprintf("count%d/subclass%d/dead%t/player%t", count, subclass, dead, player)
					t.Run(name, func(t *testing.T) {
						o.reset()
						a.ObjSubClass = subclass
						a.ObjFlags = 0
						if dead {
							a.ObjFlags = object.FlagDead
						}
						objectXferSetWord(a.UpdateData, 4, 0)
						a.InvFirstItem = nil
						o.s.Objs.Pending = nil
						for i, it := range items {
							it.ObjFlags = 0
							it.PosVec = types.Pointf{}
							it.InvHolder = nil
							it.InvNextItem = nil
							it.Field125 = nil
							it.ObjNext = nil
							it.ObjPrev = nil
							if i < count {
								it.InvHolder = a
								if i == 0 {
									a.InvFirstItem = it
								} else {
									it.Field125 = items[i-1]
									items[i-1].InvNextItem = it
								}
							}
						}
						target := b
						if !player {
							target = other
						}
						collisionRegistryWorld(9, a, target, nil)
						opened := player && !dead
						if (objectXferGetWord(a.UpdateData, 4) == 1) != opened {
							t.Fatal("chest death callback")
						}
						var pending []uint32
						for it := o.s.Objs.Pending; it != nil; it = it.ObjNext {
							if len(pending) >= count {
								t.Fatal("pending list cycle")
							}
							pending = append(pending, it.NetCode)
						}
						wantCount := 0
						if opened {
							wantCount = count
						}
						if len(pending) != wantCount {
							t.Fatal("chest drop count")
						}
						var pos [][2]uint32
						for i := 0; i < count; i++ {
							it := items[i]
							if (it.InvHolder == nil) != opened {
								t.Fatal("chest inventory ownership")
							}
							pos = append(pos, [2]uint32{math.Float32bits(it.PosVec.X), math.Float32bits(it.PosVec.Y)})
						}
						rows = append(rows, struct {
							Name      string
							Opened    uint32
							Pending   []uint32
							Positions [][2]uint32
						}{name, objectXferGetWord(a.UpdateData, 4), pending, pos})
					})
				}
			}
		}
	}
}

func TestCollisionRegistryWorldChestKey(t *testing.T) {
	o := newWorldCollisionOwner(t)
	defer noxflags.PortTestGameFlags(noxflags.GameModeQuest)()
	names := []string{"SilverKey", "GoldKey"}
	t.Cleanup(o.s.PortTestAttackTypes(1, nil, names...))
	t.Cleanup(o.s.PortTestSpellEffectTypes(names, nil, nil))
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	a.UpdateData = o.record(t, 8)
	a.Death = server.PortTestWorldChestDeath()
	t.Cleanup(func() { a.UpdateData = nil; b.InvFirstItem = nil; o.s.Objs.DeletedList = nil })
	var rows []struct {
		Name            string
		Opened, Deleted bool
		Sounds          []int
	}
	for _, name := range names {
		for _, keyClass := range []bool{false, true} {
			for _, subclass := range []object.SubClass{0x100, 0x200, 0x400, 0x800} {
				label := fmt.Sprintf("%s/key%t/subclass%d", name, keyClass, subclass)
				t.Run(label, func(t *testing.T) {
					o.reset()
					o.s.PortTestCombatAudioReset()
					*o.throttle = 0
					objectXferSetWord(a.UpdateData, 4, 0)
					a.ObjSubClass = subclass
					a.ObjFlags = 0
					key := o.s.NewObjectByTypeID(name)
					if key == nil {
						t.Fatal("key factory")
					}
					key.ObjClass = object.ClassSimple
					if keyClass {
						key.ObjClass = object.ClassKey
					}
					key.ObjFlags = 0
					key.InvHolder = b
					b.InvFirstItem = key
					o.s.Objs.DeletedList = nil
					defer func() {
						b.InvFirstItem = nil
						key.InvHolder = nil
						o.s.Objs.DeletedList = nil
						o.s.Objs.FreeObject(key)
					}()
					collisionRegistryWorld(9, a, b, nil)
					opened := keyClass && name == "SilverKey"
					if (objectXferGetWord(a.UpdateData, 4) == 1) != opened || key.ObjFlags.Has(object.FlagDestroyed) != opened {
						t.Fatal("chest key admission/consumption")
					}
					var sounds []int
					for _, e := range o.s.PortTestCombatAudioSnapshot() {
						sounds = append(sounds, int(e.ID))
					}
					want := 1012
					if opened {
						want = 234
					}
					if len(sounds) != 1 || sounds[0] != want {
						t.Fatal("chest key sound")
					}
					rows = append(rows, struct {
						Name            string
						Opened, Deleted bool
						Sounds          []int
					}{label, objectXferGetWord(a.UpdateData, 4) == 1, key.ObjFlags.Has(object.FlagDestroyed), sounds})
				})
			}
		}
	}
	collisionRegistryCapture(t, "collision-registry-world-chest-key", rows)
}
