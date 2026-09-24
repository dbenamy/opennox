//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

func TestCollisionRegistryWorldSpellProjectile(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	owner := &o.units[2]
	a.ObjClass = object.ClassMissile
	a.ObjSubClass = 2
	data := o.record(t, 28)
	a.UpdateData = data
	t.Cleanup(func() { a.UpdateData = nil; o.s.ObjSetOwner(nil, a); o.s.Objs.DeletedList = nil; b.InvFirstItem = nil })
	t.Cleanup(o.s.PortTestSpellLifecycle([]server.PortTestSpellLifecycleDef{{Index: int(spell.SPELL_CURE_POISON), Valid: true, Enabled: true, Sounds: [3]int{0, 281, 0}}}, nil))
	o.s.Spells.DefByInd(spell.SPELL_CURE_POISON).Effect = spell.SPELL_CURE_POISON
	t.Cleanup(o.s.PortTestObjectReportAnimations(6, 1))
	item := newObjectXferSimple(t, o.s)
	item.ObjClass = object.ClassArmor
	item.ObjFlags = 0x100
	mods := o.record(t, 16)
	mod := (*server.ModifierEff)(o.record(t, int(unsafe.Sizeof(server.ModifierEff{}))))
	objectXferSetWord(mods, 8, uint32(uintptr(unsafe.Pointer(mod))))
	mod.DefendCollide88.Fnc = legacy.PortTestWorldInversionCallback()
	mod.DefendCollide88.Val = 1
	item.InitData = mods
	t.Cleanup(func() { item.InitData = nil })
	var rows []struct {
		Name                         string
		Velocity                     [2]uint32
		Direction                    uint16
		Owner, Target, Poison, Frame uint32
		Deleted                      bool
		Sounds                       []int
		RNG                          int
	}
	defer func() {
		collisionRegistryCapture(t, "collision-registry-world-spell-projectile", rows)
	}()
	for _, mode := range []string{"plain", "shield", "inversion", "buff", "sword", "walking-sword", "berserker-off", "berserker-on"} {
		for _, direction := range []uint16{0, 128} {
			for _, intended := range []bool{false, true} {
				for _, flags := range []object.Flags{0, 0x20, 0x8000} {
					name := fmt.Sprintf("%s/direction%d/intended%t/flags%d", mode, direction, intended, flags)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.s.PortTestCombatAudioReset()
						o.s.Rand.Logic = prand.New(7)
						o.s.Abils.Reset()
						*(*byte)(unsafe.Add(b.UpdateDataPlayer().Player.C(), 3)) = 0
						a.ObjFlags = 0
						a.DeletedNext = nil
						o.s.Objs.DeletedList = nil
						o.s.ObjSetOwner(nil, a)
						o.s.ObjSetOwner(owner, a)
						a.Direction1 = 0
						a.PosVec = types.Ptf(112, 101)
						a.PrevPos = types.Ptf(110, 100)
						a.NewPos = a.PosVec
						a.VelVec = types.Ptf(-3, 2)
						a.Field32 = 0
						b.PosVec = types.Ptf(100, 100)
						b.Direction1 = server.Dir16(direction)
						b.ObjFlags = flags
						b.Buffs = 0
						b.InvFirstItem = nil
						b.Poison540 = 3
						*(*byte)(unsafe.Add(b.UpdateData, 88)) = 0
						*(*byte)(unsafe.Add(b.UpdateData, 236)) = 0
						objectXferSetWord(b.UpdateDataPlayer().Player.C(), 4, 0)
						*o.globals["extensions"] = 0
						if mode == "shield" {
							*(*byte)(unsafe.Add(b.UpdateData, 88)) = 16
						}
						if mode == "sword" {
							*(*byte)(unsafe.Add(b.UpdateData, 88)) = 13
							objectXferSetWord(b.UpdateDataPlayer().Player.C(), 4, 0x400)
						}
						if mode == "walking-sword" {
							*o.globals["extensions"] = 4
							objectXferSetWord(b.UpdateDataPlayer().Player.C(), 4, 0x400)
						}
						if mode == "berserker-off" || mode == "berserker-on" {
							*(*byte)(unsafe.Add(b.UpdateData, 88)) = 1
							*(*byte)(unsafe.Add(b.UpdateDataPlayer().Player.C(), 3)) = 1
							o.s.Abils.GetFor(b).ExecList = &server.ExecAbilityClass{Abil: server.Ability(1), Active: 1}
							if mode == "berserker-on" {
								*o.globals["extensions"] = 16
							}
						}
						if mode == "inversion" {
							b.InvFirstItem = item
						}
						if mode == "buff" {
							b.Buffs = 1 << 27
						}
						objectXferSetWord(data, 0, uint32(uintptr(owner.CObj())))
						objectXferSetWord(data, 4, uint32(uintptr(owner.CObj())))
						if intended {
							objectXferSetWord(data, 4, uint32(uintptr(b.CObj())))
						}
						objectXferSetWord(data, 8, uint32(uintptr(owner.CObj())))
						objectXferSetWord(data, 12, uint32(spell.SPELL_CURE_POISON))
						objectXferSetWord(data, 16, 1)
						collisionRegistryWorld(8, a, b, nil)
						deleted := a.ObjFlags.Has(object.FlagDestroyed)
						if flags != 0 && (deleted || a.Direction1 != 0 || a.VelVec != types.Ptf(-3, 2) || a.Owner() != owner) {
							t.Fatal("projectile admission")
						}
						if deleted && (!intended || b.Poison540 != 2) {
							t.Fatal("spell delivered to intended target")
						}
						if !deleted && b.Poison540 != 3 {
							t.Fatal("poison changed without spell delivery")
						}
						if mode == "plain" && flags == 0 && deleted != intended {
							t.Fatal("plain target selection")
						}
						if mode == "inversion" && flags == 0 && deleted {
							t.Fatal("inversion stops delivery")
						}
						own := uint32(3)
						if a.Owner() == b {
							own = 2
						} else if a.Owner() != owner {
							t.Fatal("projectile owner")
						}
						target := uint32(3)
						if objectXferGetWord(data, 4) == uint32(uintptr(b.CObj())) {
							target = 2
						} else if objectXferGetWord(data, 4) != uint32(uintptr(owner.CObj())) {
							t.Fatal("projectile target")
						}
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							sounds = append(sounds, int(e.ID))
						}
						rows = append(rows, struct {
							Name                         string
							Velocity                     [2]uint32
							Direction                    uint16
							Owner, Target, Poison, Frame uint32
							Deleted                      bool
							Sounds                       []int
							RNG                          int
						}{name, [2]uint32{math.Float32bits(a.VelVec.X), math.Float32bits(a.VelVec.Y)}, uint16(a.Direction1), own, target, uint32(b.Poison540), a.Field32, deleted, sounds, o.s.Rand.Logic.Index()})
					})
				}
			}
		}
	}
}

func TestCollisionRegistryWorldSpellWall(t *testing.T) {
	o := newWorldCollisionOwner(t)
	a := newObjectXferSimple(t, o.s)
	var rows [][4]uint32
	for _, normal := range []types.Pointf{{1, 0}, {0, 1}, {1, 1}, {-1, 1}, {-1, -1}, {0, 0}} {
		a.VelVec = types.Ptf(-3, 2)
		collisionRegistryWorld(8, a, nil, &normal)
		want := types.Ptf(2, -3)
		if normal.X*normal.Y > 0 {
			want = types.Ptf(-2, 3)
		}
		if a.VelVec != want {
			t.Fatal("wall reflection")
		}
		rows = append(rows, [4]uint32{math.Float32bits(normal.X), math.Float32bits(normal.Y), math.Float32bits(a.VelVec.X), math.Float32bits(a.VelVec.Y)})
	}
	collisionRegistryWorld(8, a, nil, nil)
	collisionRegistryCapture(t, "collision-registry-world-spell-wall", rows)
}
