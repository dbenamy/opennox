//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestMonsterControlDeathRewards(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Mimic", "Zombie", "VileZombie"}, nil, true, 0, 0))
	oldSend := o.s.NetSendPacketXxx
	o.s.NetSendPacketXxx = legacy.Nox_xxx_netSendPacket_4E5030
	t.Cleanup(func() { o.s.NetSendPacketXxx = oldSend })
	noxServer.abilities.s = noxServer
	o.s.Abils.Reset()
	a := newCreatureXferObject(t, o.s, "Monster")
	aSaved := *a
	ud := a.UpdateDataMonster()
	savedUD := *ud
	child := newObjectXferSimple(t, o.s)
	savedChild := *child
	pl := o.units[0].ControllingPlayer()
	savedPlayer := *pl
	lists, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	t.Cleanup(func() { *a = aSaved; *ud = savedUD; *child = savedChild; *pl = savedPlayer })
	ids := collisionCoreIDs(a, child, &o.units[0])
	ids[unsafe.Pointer(pl)] = 2001
	o.objects = ids
	type row struct {
		Flags    uint32
		Killer   int
		Special  bool
		Return   uint32
		Cooldown int
		Stats    [11]uint32
		Reports  legacy.PortTestReliableReportState
	}
	var rows []row
	for _, flags := range []uint32{0, 2048, 4096} {
		for killer := 0; killer < 3; killer++ {
			for _, special := range []bool{false, true} {
				legacy.PortTestWorldMotionList("decay-clear", nil, 0)
				*lists["decay"] = 0
				o.reset()
				o.resetQueues()
				*a = aSaved
				*ud = server.MonsterUpdateData{}
				*child = savedChild
				*pl = savedPlayer
				a.ObjClass = object.ClassMonster
				a.ObjSubClass = 0x2000
				a.ObjFlags = 4
				a.TypeInd = uint16(o.s.Types.IndByID("Mimic"))
				a.Obj130 = nil
				ud.AIStackInd = 0
				ud.AIStack[0].Action = 0
				if special {
					ud.Field546 = 2
					ud.Field547 = 2
				}
				if killer == 1 {
					a.Obj130 = &o.units[0]
				}
				if killer == 2 {
					child.ObjOwner = &o.units[0]
					a.Obj130 = child
				}
				o.s.Abils.GetFor(&o.units[0]).Cooldowns[server.AbilityBerserk] = 123
				o.s.Rand.Logic = prand.New(23)
				resetFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
				rv := legacy.PortTestMonsterControl("death", a, nil, nil, 0)
				if id, ok := ids[unsafe.Pointer(uintptr(rv))]; ok {
					rv = id
				}
				stats := questRuntimeStats(unsafe.Pointer(pl))
				cooldown := o.s.Abils.GetFor(&o.units[0]).Cooldowns[server.AbilityBerserk]
				if (cooldown == 0) != (special && killer != 0 && flags&2048 == 0) {
					t.Fatal("special kill ability reset", flags, killer, special, cooldown)
				}
				if stats[3] != uint32(bool2int(flags&4096 != 0 && killer != 0)) {
					t.Fatal("quest kill counter", flags, killer, stats)
				}
				rows = append(rows, row{flags, killer, special, rv, cooldown, stats, o.state()})
				resetFlags()
			}
		}
	}
	legacy.PortTestWorldMotionList("decay-clear", nil, 0)
	spellbookCapture(t, "monster-control-death-rewards", rows, "643af4ccfb4510d47675e7afa81d1befe06657fcaec9d1e3c92d604723245111")
}

func TestMonsterControlDeathInventory(t *testing.T) {
	o := newCollisionCoreOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Mimic", "Zombie", "VileZombie", "Glyph", "Torch", "Lantern"}, nil, true, 0, 0))
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	_, _, restore := o.s.PortTestPathWalls()
	t.Cleanup(restore)
	for _, off := range []uintptr{1568252, 1568256, 1568244} {
		serverConfigOwnBytes(t, 0x5D4594, off, 4)
	}
	lists, restore := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restore)
	oldRNG, oldPending := o.s.Rand.Logic, o.s.Objs.Pending
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG; o.s.Objs.Pending = oldPending })
	a := newCreatureXferObject(t, o.s, "Monster")
	sa := *a
	ud := a.UpdateDataMonster()
	savedUD := *ud
	it := newObjectXferSimple(t, o.s)
	si := *it
	t.Cleanup(func() { *a = sa; *ud = savedUD; *it = si })
	ids := collisionCoreIDs(a, it)
	o.objects = ids
	type row struct {
		Keep, Quest                                   bool
		Flags, Holder, Inventory, Pending, Decay, RNG uint32
		Pos                                           [2]float32
	}
	var rows []row
	for _, keep := range []bool{false, true} {
		for _, quest := range []bool{false, true} {
			legacy.PortTestWorldMotionList("decay-clear", nil, 0)
			o.resetQueues()
			o.reset()
			*a = sa
			*it = si
			*ud = server.MonsterUpdateData{}
			worldGeometryResetObject(a, 1001, 100, 100, false)
			a.TypeInd = uint16(o.s.Types.IndByID("Mimic"))
			a.ObjClass = object.ClassMonster
			a.ObjSubClass = 0
			if keep {
				a.ObjSubClass = 0x2000
			}
			ud.AIStackInd = 0
			ud.AIStack[0].Action = 0
			it.ObjClass = object.ClassSimple
			it.ObjFlags = 0 // Inventory objects are inactive until the real drop creates them.
			it.InvHolder = a
			it.InvNextItem = nil
			it.Drop.Ptr = nil
			a.InvFirstItem = it
			o.s.Objs.Pending = nil
			o.s.Rand.Logic = prand.New(23)
			flag := noxflags.GameFlag(0)
			if quest {
				flag = 4096
			}
			resetFlags := noxflags.PortTestGameFlags(flag)
			legacy.PortTestMonsterControl("death", a, nil, nil, 0)
			if (it.InvHolder == a) != keep || (a.InvFirstItem == it) != keep {
				t.Fatal("death inventory retention", keep, quest)
			}
			if !keep && o.s.Objs.Pending != it {
				t.Fatal("dropped item not staged in world")
			}
			rows = append(rows, row{keep, quest, uint32(it.ObjFlags), collisionCoreID(t, ids, uint32(uintptr(it.InvHolder.CObj()))), collisionCoreID(t, ids, uint32(uintptr(a.InvFirstItem.CObj()))), collisionCoreID(t, ids, uint32(uintptr(o.s.Objs.Pending.CObj()))), collisionCoreID(t, ids, *lists["decay"]), uint32(o.s.Rand.Logic.Index()), [2]float32{it.PosVec.X, it.PosVec.Y}})
			resetFlags()
		}
	}
	legacy.PortTestWorldMotionList("decay-clear", nil, 0)
	spellbookCapture(t, "monster-control-death-inventory", rows, "0c964e3465fab2d0e100ee9811bd8b4bed37d73466183797f57c2057416b26d4")
}
