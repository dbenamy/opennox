//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerDeathCrownTransfer(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "Glyph", "Torch", "Lantern"}, nil, true, 0, 0))
	for _, off := range []uintptr{1567716, 1568248, 1568252, 1568256, 1568244} {
		clear(serverConfigOwnBytes(t, 0x5D4594, off, 4))
	}
	cache := serverConfigOwnBytes(t, 0x5D4594, 2488728, 4)
	cache[0] = 1
	clear(serverConfigOwnBytes(t, 0x587000, 279432, 12))
	_, restoreMotion := legacy.PortTestWorldMotionListGlobals()
	t.Cleanup(restoreMotion)
	oldPending := o.s.Objs.Pending
	t.Cleanup(func() { o.s.Objs.Pending = oldPending })
	type record struct {
		Name           string
		Assigned       uint32
		Owned, Held    bool
		Pos            types.Pointf
		Scores, Deaths [3]uint32
		Reports        legacy.PortTestReliableReportState
	}
	var rows []record
	for _, teamed := range []bool{false, true} {
		for _, friendly := range []bool{false, true} {
			for _, disabled := range []bool{false, true} {
				for _, held := range []bool{false, true} {
					name := fmt.Sprintf("team=%t/friendly=%t/disabled=%t/held=%t", teamed, friendly, disabled, held)
					t.Run(name, func(t *testing.T) {
						o.reset()
						noxflags.ResetGame()
						noxflags.SetGame(noxflags.GameHost | noxflags.GameModeKOTR)
						noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
						if disabled {
							noxflags.SetGamePlay(4)
						}
						o.balance(map[string]float64{"KotRKingKillsPawnPoints": 2, "KotRPawnKillsKingPoints": 7, "KotRKingKillsKingPoints": 11})
						o.s.Teams.Reset()
						o.s.Teams.ActiveCnt = 0
						o.s.Teams.Create(1)
						o.s.Teams.Create(2)
						for i := range o.units {
							u := &o.units[i]
							u.TeamVal = server.ObjectTeam{}
							objectXferSetWord(u.UpdateDataPlayer().Player.C(), 2136, 0)
							objectXferSetWord(u.UpdateDataPlayer().Player.C(), 2140, 0)
						}
						victim, killer := &o.units[0], &o.units[1]
						if teamed {
							for i, u := range []*server.Object{victim, killer} {
								id := 1 + i
								if friendly {
									id = 1
								}
								legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(id), u.TeamPtr(), 0, int(u.NetCode), 0)
							}
						} else if friendly {
							killer = victim
						}
						victim.PosVec = types.Pointf{X: 81, Y: 123}
						c := o.s.NewObjectByTypeID("Crown")
						if c == nil {
							t.Fatal("crown")
						}
						trackObjectXferTyped(t, o.s, c)
						c.ObjClass = 0
						c.NetCode = 200
						o.s.ObjSetOwner(victim, c)
						if held {
							victim.InvFirstItem = c
							c.InvHolder = victim
						}
						o.s.Objs.Pending = nil
						defer func() {
							legacy.PortTestWorldMotionList("decay-remove", c, 0)
							o.s.ObjClearOwner(c)
							victim.InvFirstItem = nil
							c.InvHolder = nil
							c.InvNextItem = nil
							c.Field125 = nil
							o.s.Objs.Pending = nil
							c.ObjNext = nil
							c.ObjPrev = nil
						}()
						o.reset()
						legacy.PortTestPlayerDeath("kotr", victim, killer, nil, nil)
						r := record{Name: name, Owned: c.ObjOwner == victim, Held: c.InvHolder == victim, Pos: c.PosVec, Reports: o.state()}
						p := *(*unsafe.Pointer)(unsafe.Add(c.UpdateData, 4))
						if p != nil {
							if p != killer.CObj() {
								t.Fatal("unexpected pending crown owner")
							}
							r.Assigned = killer.NetCode
						}
						transfer := !teamed && !friendly && !disabled
						var want uint32
						if transfer {
							want = killer.NetCode
						}
						if r.Assigned != want || r.Owned != (!(transfer && held)) || r.Held != (held && !transfer) {
							t.Fatalf("crown assignment/drop %+v", r)
						}
						if transfer && held && r.Pos != victim.PosVec {
							t.Fatal("drop position")
						}
						for i := range o.units {
							pl := o.units[i].UpdateDataPlayer().Player
							r.Scores[i] = objectXferGetWord(pl.C(), 2136)
							r.Deaths[i] = objectXferGetWord(pl.C(), 2140)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
	spellbookCapture(t, "player-death-crown-transfer", rows, "521c860f04a5c8337632f1ff4ddf3e8dcf04072996b2e73d38261b73df28b08e")
}
