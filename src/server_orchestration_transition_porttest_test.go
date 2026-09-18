//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestServerOrchestrationPlayerTransition(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestControlsTypes(nil, 0))
	oldCoreAbilities := o.s.Abils.ByUnit
	o.s.Abils.Reset()
	t.Cleanup(func() { o.s.Abils.ByUnit = oldCoreAbilities })
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	serverConfigOwnBytes(t, 0x5D4594, 1523072, 4)
	*memmap.PtrUint32(0x5D4594, 1523072) = 2 // Preserve existing loadout.
	oldSolo := gameIsSwitchToSolo
	oldControl := o.s.Players.Control
	oldUpdates := o.s.Objs.UpdatableList
	t.Cleanup(func() {
		gameIsSwitchToSolo = oldSolo
		o.s.Players.Control = oldControl
		o.s.Objs.UpdatableList = oldUpdates
	})
	health := make([]*server.HealthData, 3)
	for i := range o.units {
		u := &o.units[i]
		old := u.HealthData
		t.Cleanup(func() { u.HealthData = old })
		health[i] = (*server.HealthData)(o.record(t, int(unsafe.Sizeof(server.HealthData{}))))
		u.HealthData = health[i]
	}
	type unit struct {
		State                                     uint8
		Previous, Current, Loadout, Frame, Status uint32
		HP, Mana                                  uint16
		EmptyControl, CameraCleared               bool
	}
	type row struct {
		Flags                uint32
		Dead, Observer, Solo bool
		Units                [3]unit
		Packets              [][]byte
	}
	var rows []row
	for _, flags := range []uint32{0, 512, 4096, 4608} {
		for _, dead := range []bool{false, true} {
			for _, observer := range []bool{false, true} {
				for _, solo := range []bool{false, true} {
					t.Run(fmt.Sprintf("flags%x/dead%t/observer%t/solo%t", flags, dead, observer, solo), func(t *testing.T) {
						o.reset()
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						gameIsSwitchToSolo = solo
						o.s.Objs.UpdatableList = nil
						for i := range o.units {
							u := &o.units[i]
							ud := u.UpdateDataPlayer()
							pl := ud.Player
							u.ObjFlags = 0
							if dead {
								u.ObjFlags = object.FlagDead
							}
							u.IsUpdatable = 0
							u.UpdatableNext = nil
							u.UpdatablePrev = nil
							u.Buffs = 0
							ad := o.s.Abils.GetFor(u)
							for j := range ad.Cooldowns {
								ad.Cooldowns[j] = 73
							}
							u.HealthData.Cur = 11
							u.HealthData.Max = 87
							ud.ManaCur = 7
							ud.ManaMax = 41
							ud.State = 0
							u.Field34 = 0
							objectXferSetWord(u.UpdateData, 280, 0)
							objectXferSetWord(pl.C(), 4676, uint32(50+i))
							objectXferSetWord(pl.C(), 4680, 123)
							objectXferSetWord(pl.C(), 4700, 77)
							pl.Field3680 = 0
							if observer {
								pl.Field3680 = 0x20
							}
							pl.CameraFollowObj = &o.units[(i+1)%3]
							o.s.Players.Control.Player(int(pl.PlayerInd)).Reset()
							o.s.Players.Control.Player(int(pl.PlayerInd)).Append([]server.PlayerCtrl{{Code: 1, Data: 123, Active: true}})
						}
						legacy.PortTestServerOrchestration("players", nil, 0)
						r := row{Flags: flags, Dead: dead, Observer: observer, Solo: solo, Packets: visibilityEffectsPackets(o.s)}
						for i := range o.units {
							u := &o.units[i]
							ud := u.UpdateDataPlayer()
							pl := ud.Player
							v := unit{State: uint8(ud.State), Previous: objectXferGetWord(pl.C(), 4680), Current: objectXferGetWord(pl.C(), 4676), Loadout: objectXferGetWord(pl.C(), 4700), Frame: u.Field34, Status: pl.Field3680, HP: u.HealthData.Cur, Mana: ud.ManaCur, EmptyControl: o.s.Players.Control.Player(int(pl.PlayerInd)).First() == nil, CameraCleared: pl.CameraFollowObj == nil}
							r.Units[i] = v
							defaults := dead || flags&512 == 0
							wantLoadout := uint32(77)
							wantHP, wantMana := uint16(11), uint16(7)
							if defaults {
								wantLoadout = 1
								if !solo {
									wantHP, wantMana = 87, 41
								}
							}
							if v.Loadout != wantLoadout || v.HP != wantHP || v.Mana != wantMana || !v.EmptyControl || v.Frame != o.s.Frame() || v.CameraCleared != observer {
								t.Fatal("player transition", i, v, wantLoadout, wantHP, wantMana)
							}
							if !dead && v.State != 13 {
								t.Fatal("player state", v.State)
							}
							if flags&4096 != 0 {
								if v.Current != 0 || v.Previous != uint32(50+i) {
									t.Fatal("quest stats", v)
								}
							} else if v.Current != uint32(50+i) || v.Previous != 123 {
								t.Fatal("nonquest stats", v)
							}
							for _, cooldown := range o.s.Abils.GetFor(u).Cooldowns {
								want := 73
								if defaults {
									want = 0
								}
								if cooldown != want {
									t.Fatal("ability cooldown", cooldown, want)
								}
							}
							if objectXferGetWord(u.UpdateData, 280) != 0 {
								t.Fatal("shop pointer")
							}
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-orchestration-player-transition", rows, "438fe05865744301dbb154d101fd3e36b93834d4ea8d68f6b75d1c47b9cd0062")
}

func TestServerOrchestrationFreshLoadout(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestControlsTypes(nil, 0))
	oldAbilities := noxServer.abilities
	noxServer.abilities.Init(noxServer)
	t.Cleanup(func() { noxServer.abilities = oldAbilities })
	oldPickup := o.s.Objs.DefaultPickup
	o.s.Objs.DefaultPickup = nox_xxx_pickupDefault_4F31E0
	t.Cleanup(func() { o.s.Objs.DefaultPickup = oldPickup })
	oldSolo := gameIsSwitchToSolo
	gameIsSwitchToSolo = true
	t.Cleanup(func() { gameIsSwitchToSolo = oldSolo })
	serverConfigOwnBytes(t, 0x5D4594, 1523072, 4)
	serverConfigOwnBytes(t, 0x5D4594, 1564960, 4)
	*memmap.PtrUint32(0x5D4594, 1564960) = 0
	var mods []*server.ModifierEff
	var names []*byte
	for _, name := range []string{"UserColor1", "ArmorQuality1", "Material1", "Replenishment1"} {
		mods = append(mods, (*server.ModifierEff)(o.record(t, int(unsafe.Sizeof(server.ModifierEff{})))))
		p := o.record(t, len(name)+1)
		copy(unsafe.Slice((*byte)(p), len(name)+1), name)
		names = append(names, (*byte)(p))
	}
	t.Cleanup(o.s.PortTestControlsModifiers(mods, names))
	for i := range o.units {
		u := &o.units[i]
		old := u.HealthData
		t.Cleanup(func() { u.HealthData = old })
		u.HealthData = (*server.HealthData)(o.record(t, int(unsafe.Sizeof(server.HealthData{}))))
		u.HealthData.Cur = 43
		u.HealthData.Max = 99
		u.CarryCapacity = 64
	}
	type row struct {
		Map       uint32
		Inventory [3][]string
		Loadout   [3]uint32
		Queue     legacy.PortTestReliableReportState
	}
	var rows []row
	for _, state := range []uint32{0, 1, 2, 3, 0x80000000, 0x80000001, 0xfffffffe, 0xffffffff} {
		t.Run(fmt.Sprintf("map%x", state), func(t *testing.T) {
			o.reset()
			defer noxflags.PortTestGameFlags(0)()
			*memmap.PtrUint32(0x5D4594, 1523072) = state
			for i := range o.units {
				u := &o.units[i]
				u.ObjFlags = 0
				u.InvFirstItem = nil
				u.Field129 = nil
				u.UpdateDataPlayer().State = 0
				objectXferSetWord(u.UpdateDataPlayer().Player.C(), 4700, 77)
			}
			defer func() {
				for i := range o.units {
					u := &o.units[i]
					for it := u.InvFirstItem; it != nil; {
						next := it.InvNextItem
						o.s.Objs.FreeObject(it)
						it = next
					}
					u.InvFirstItem = nil
					u.Field129 = nil
				}
			}()
			legacy.PortTestServerOrchestration("players", nil, 0)
			r := row{Map: state, Queue: o.state()}
			for i := range o.units {
				u := &o.units[i]
				r.Loadout[i] = objectXferGetWord(u.UpdateDataPlayer().Player.C(), 4700)
				for it := u.InvFirstItem; it != nil; it = it.InvNextItem {
					if len(r.Inventory[i]) >= 5 {
						t.Fatal("loadout cycle")
					}
					r.Inventory[i] = append(r.Inventory[i], o.s.Types.ByInd(int(it.TypeInd)).ID())
					if it.InvHolder != u || it.ObjOwner != u {
						t.Fatal("item owner")
					}
				}
				want := "[woodenshield longsword streetsneakers streetpants]"
				if state&2 != 0 {
					want = "[]"
				}
				if fmt.Sprint(r.Inventory[i]) != want || r.Loadout[i] != 1 {
					t.Fatal("fresh loadout", r.Inventory[i], want, r.Loadout[i])
				}
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "server-orchestration-fresh-loadout", rows, "70a384882d1bb5ea31f977c2b1b546a794ccbdfbba03a0627073cae743a8aabb")
}

func TestServerOrchestrationOpenShopTransition(t *testing.T) {
	o := newMatchRosterOwner(t)
	serverConfigOwnBytes(t, 0x5D4594, 2386364, 128)
	clear(unsafe.Slice(memmap.PtrUint8(0x5D4594, 2386364), 128))
	// Quest shops remain cached after exit; the real session record stays owned here.
	session := (*server.TradeSession)(o.record(t, int(unsafe.Sizeof(server.TradeSession{}))))
	session.Kind = 1
	session.Active = 1
	session.Units = [2]*server.Object{&o.units[0], &o.units[1]}
	for _, u := range session.Units {
		*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 280)) = unsafe.Pointer(session)
	}
	defer noxflags.PortTestGameFlags(4096 | 512)()
	legacy.PortTestServerOrchestration("players", nil, 0)
	for _, u := range session.Units {
		if *(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 280)) != nil {
			t.Fatal("shop participant still attached")
		}
	}
	if session.Active != 0 || *memmap.PtrPtr(0x5D4594, 2386364+4) != unsafe.Pointer(session) {
		t.Fatal("shop cache ownership")
	}
	spellbookCapture(t, "server-orchestration-open-shop-transition", struct {
		Active  uint32
		Cached  bool
		Packets [][]byte
	}{session.Active, true, visibilityEffectsPackets(o.s)}, "961f0ce56da23ec77b65914361581b1b1e55a750f55141a2c260a11ff3dfe9a7")
}
