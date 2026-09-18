//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestMonsterControlDeath(t *testing.T) {
	o := newCollisionCoreOwner(t)
	oldSend := o.s.NetSendPacketXxx
	o.s.NetSendPacketXxx = legacy.Nox_xxx_netSendPacket_4E5030
	t.Cleanup(func() { o.s.NetSendPacketXxx = oldSend })
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Mimic", "Zombie", "VileZombie"}, nil, true, 0, 0))
	types := server.PortTestMonsterStateTypeIDs{Mimic: o.s.Types.IndByID("Mimic"), Zombie: o.s.Types.IndByID("Zombie"), VileZombie: o.s.Types.IndByID("VileZombie")}
	lists, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	a := newCreatureXferObject(t, o.s, "Monster")
	sa := *a
	ud := a.UpdateDataMonster()
	saved := *ud
	h := a.HealthData
	oldHealth := *h
	t.Cleanup(func() { *a = sa; *ud = saved; *h = oldHealth })
	child := newObjectXferSimple(t, o.s)
	savedChild := *child
	t.Cleanup(func() { *child = savedChild })
	ids := collisionCoreIDs(a, child, &o.units[0], &o.units[1], &o.units[2])
	o.objects = ids
	oldUnits := append([]server.Object(nil), o.units...)
	oldPlayers := make([]server.Player, len(o.units))
	for i := range o.units {
		oldPlayers[i] = *o.units[i].ControllingPlayer()
	}
	t.Cleanup(func() {
		for i := range o.units {
			*o.units[i].ControllingPlayer() = oldPlayers[i]
			o.units[i] = oldUnits[i]
		}
	})
	type row struct {
		Zombie, Owner, Observer, Quest                                    bool
		Status                                                            uint32
		Stack                                                             int8
		Return, Flags, Subclass, Decay, RNG, Head, OwnerAfter, ChildOwner uint32
		Actions                                                           [24]uint32
		Players                                                           [][2]uint32
		Packets                                                           [][]byte
		Reports                                                           legacy.PortTestReliableReportState
	}
	var rows []row
	for _, zombie := range []bool{false, true} {
		for _, owner := range []bool{false, true} {
			for _, observer := range []bool{false, true} {
				for _, quest := range []bool{false, true} {
					for _, status := range []uint32{0, 0x80, 0x100080} {
						for _, stack := range []int8{-1, 0, 23} {
							legacy.PortTestWorldMotionList("decay-clear", nil, 0)
							o.resetQueues()
							o.reset()
							*a = sa
							*child = savedChild
							*ud = server.MonsterUpdateData{}
							*h = oldHealth
							for i := range o.units {
								*o.units[i].ControllingPlayer() = oldPlayers[i]
								o.units[i] = oldUnits[i]
							}
							a.ObjClass = object.ClassMonster
							a.ObjFlags = 0x84
							a.ObjSubClass = 0x2180
							a.ObjOwner = nil
							a.Field128 = nil
							a.Field129 = nil
							a.Obj130 = nil
							a.Field115 = 0
							a.Field116 = 0
							a.Field117 = 0
							a.Buffs = 1 << server.ENCHANT_BLINDED
							a.TypeInd = uint16(types.Mimic)
							if zombie {
								a.TypeInd = uint16(types.Zombie)
							}
							ud.AIStackInd = stack
							ud.StatusFlags = object.MonsterStatus(status)
							for i := 0; i <= int(stack); i++ {
								ud.AIStack[i] = server.AIStackItem{Action: 1}
							}
							if owner {
								o.s.ObjSetOwner(&o.units[0], a)
							}
							child.ObjFlags = 4
							child.ObjOwner = nil
							child.Field128 = nil
							child.Field129 = nil
							o.s.ObjSetOwner(a, child)
							if observer {
								pl := o.units[1].ControllingPlayer()
								pl.Field3680 |= 2
								pl.CameraFollowObj = a
							}
							game := noxflags.GameFlag(0)
							if quest {
								game = 4096
							}
							resetFlags := noxflags.PortTestGameFlags(game)
							o.s.SetFrame(0xfffffff0)
							o.s.Rand.Logic = prand.New(23)
							rv := legacy.PortTestMonsterControl("death", a, nil, nil, 0)
							if ud.AIStackInd != 1 || ud.AIStack[0].Action != 31 || ud.AIStack[1].Action != 30 {
								t.Fatal("death stack")
							}
							if observer && o.units[1].ControllingPlayer().ObserveTarget() != nil {
								t.Fatal("death must release observer even for zombie")
							}
							if !zombie && (a.ObjOwner != nil || a.Buffs != 0 || a.ObjFlags&0x80 != 0) {
								t.Fatal("nonzombie death cleanup")
							}
							if zombie && (*lists["decay"] != 0 || a.Buffs == 0) {
								t.Fatal("zombie must preserve buffs and skip decay")
							}
							r := row{Zombie: zombie, Owner: owner, Observer: observer, Quest: quest, Status: status, Stack: stack, Return: rv, Flags: uint32(a.ObjFlags), Subclass: uint32(a.ObjSubClass), Decay: a.Field34, RNG: uint32(o.s.Rand.Logic.Index()), Head: collisionCoreID(t, ids, *lists["decay"]), OwnerAfter: collisionCoreID(t, ids, uint32(uintptr(a.ObjOwner.CObj()))), ChildOwner: collisionCoreID(t, ids, uint32(uintptr(child.ObjOwner.CObj()))), Packets: visibilityEffectsPackets(o.s), Reports: o.state()}
							for i, s := range ud.AIStack {
								r.Actions[i] = s.Action
							}
							for i := range o.units {
								p := o.units[i].ControllingPlayer()
								r.Players = append(r.Players, [2]uint32{p.Field3680, collisionCoreID(t, ids, uint32(uintptr(p.CameraFollowObj.CObj())))})
							}
							rows = append(rows, r)
							resetFlags()
						}
					}
				}
			}
		}
	}
	legacy.PortTestWorldMotionList("decay-clear", nil, 0)
	spellbookCapture(t, "monster-control-death", rows, "9677386c2227c0ac76a50ea9b26a31577d260dd332948a59e6d45e269343aa54")
}

func TestMonsterControlRevive(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Mimic", "Zombie", "VileZombie"}, nil, true, 0, 0))
	types := server.PortTestMonsterStateTypeIDs{Mimic: o.s.Types.IndByID("Mimic"), Zombie: o.s.Types.IndByID("Zombie"), VileZombie: o.s.Types.IndByID("VileZombie")}
	a := newCreatureXferObject(t, o.s, "Monster")
	sa := *a
	ud := a.UpdateDataMonster()
	saved := *ud
	h := a.HealthData
	sh := *h
	t.Cleanup(func() { *a = sa; *ud = saved; *h = sh })
	type row struct {
		Type, Class                                            int
		Flags, Action, Status, Return, FinalFlags, FinalStatus uint32
		Stack                                                  int8
		Actions                                                [24]uint32
		Health                                                 uint16
	}
	var rows []row
	for typ := 0; typ < 3; typ++ {
		for cls := 0; cls < 3; cls++ {
			for _, flags := range []uint32{4, 0x8004, 0x805c} {
				for _, act := range []uint32{0, 1, 30, 31} {
					for _, status := range []uint32{0, 0x100000, 0xffffffff} {
						*a = sa
						*ud = server.MonsterUpdateData{}
						*h = sh
						a.TypeInd = uint16([]int{types.Mimic, types.Zombie, types.VileZombie}[typ])
						a.ObjClass = object.ClassMonster
						a.ObjFlags = object.Flags(flags)
						if cls == 1 {
							a.ObjClass = object.ClassSimple
						}
						ud.StatusFlags = object.MonsterStatus(status)
						ud.AIStackInd = 0
						ud.AIStack[0].Action = act
						h.Cur = 1
						h.Max = 123
						target := a
						if cls == 2 {
							target = nil
						}
						rv := legacy.PortTestMonsterControl("revive", target, nil, nil, 0)
						if rv == uint32(uintptr(a.CObj())) {
							rv = 1001
						}
						if cls != 0 || flags&0x8000 == 0 {
							if ud.StatusFlags != object.MonsterStatus(status) || h.Cur != 1 {
								t.Fatal("revive admission mutated unit")
							}
						}
						if cls == 0 && flags&0x8000 != 0 && ud.StatusFlags&0x100000 != 0 {
							t.Fatal("revive status bit not cleared")
						}
						if cls == 0 && flags&0x8000 != 0 && typ != 0 && act == 31 && (h.Cur != 123 || ud.AIStackInd != 1 || ud.AIStack[0].Action != 61 || ud.AIStack[1].Action != 35) {
							t.Fatal("zombie get-up stack/health", ud.AIStack[0].Action, ud.AIStack[1].Action)
						}
						r := row{Type: typ, Class: cls, Flags: flags, Action: act, Status: status, Return: rv, FinalFlags: uint32(a.ObjFlags), FinalStatus: uint32(ud.StatusFlags), Stack: ud.AIStackInd, Health: h.Cur}
						for i, s := range ud.AIStack {
							r.Actions[i] = s.Action
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	spellbookCapture(t, "monster-control-revive", rows, "0fa65ceb0e83cb8729fb96f5989894dbad6acc464226c18d1926315d6f6f8857")
}

func TestMonsterControlChapterEnd(t *testing.T) {
	o := newReliableReportsOwner(t)
	configure := func(count int) {
		for i := range o.units {
			p := o.units[i].ControllingPlayer()
			p.Active = 0
			if i < count {
				p.Active = 1
			}
		}
	}
	disable, restore := legacy.PortTestMonsterChapterOwner()
	t.Cleanup(restore)
	serverConfigOwnBytes(t, 0x5D4594, 2386828, 8)
	type row struct {
		Count           int
		Chapter         byte
		Done            uint32
		Return, Disable uint32
		Reports         legacy.PortTestReliableReportState
	}
	var rows []row
	for count := 0; count <= 3; count++ {
		for _, chapter := range []byte{0, 1, 127, 128, 255} {
			for _, done := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				configure(count)
				o.reset()
				*disable = 0x12345678
				*memmap.PtrUint8(0x5D4594, 2386828) = chapter
				*memmap.PtrUint32(0x5D4594, 2386832) = done
				rv := legacy.PortTestMonsterControl("chapter", nil, nil, nil, 0)
				if *disable != 1 {
					t.Fatal("chapter must stop map drawing")
				}
				state := o.state()
				want := -1
				if count > 0 {
					want = []int{1, 7, 31}[count-1]
				}
				if want < 0 {
					if len(state.Nodes) != 0 {
						t.Fatal("chapter with no player enqueued")
					}
				} else {
					if len(state.Nodes) != 1 || int(state.Nodes[0].To) != want {
						t.Fatal("chapter recipient", count, want, state.Nodes)
					}
					flag := byte(0)
					if done == 1 {
						flag = 1
					}
					p := state.Nodes[0].Data
					if len(p) != 3 || p[0] != 214 || p[1] != chapter || p[2] != flag {
						t.Fatal("chapter payload", p)
					}
				}
				rows = append(rows, row{count, chapter, done, rv, *disable, state})
			}
		}
	}
	spellbookCapture(t, "monster-control-chapter", rows, "08eb38dbef32836b7c2fa13a9c9433ab86ffddb5db8c8a4ff494370d6501a4c2")
}
