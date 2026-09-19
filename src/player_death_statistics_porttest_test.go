//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
)

func TestPlayerDeathStatistics(t *testing.T) {
	o := newMatchRosterOwner(t)
	words, restore := legacy.PortTestStatisticsGlobals()
	t.Cleanup(restore)
	events := serverConfigOwnBytes(t, 0x5D4594, 608320, 128)
	type row struct {
		Name  string
		Pairs [][2]byte
	}
	var rows []row
	for _, op := range []string{"arena", "elimination"} {
		for _, actor := range []int{-1, 0, 1} {
			for _, friendly := range []bool{false, true} {
				for _, assist := range []bool{false, true} {
					for _, enabled := range []bool{false, true} {
						name := fmt.Sprintf("%s/actor=%d/friendly=%t/assist=%t/logging=%t", op, actor, friendly, assist, enabled)
						t.Run(name, func(t *testing.T) {
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameHost | noxflags.GameOnline)
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							o.s.Teams.Create(1)
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								objectXferSetWord(u.UpdateDataPlayer().Player.C(), 4648, uint32(i))
								objectXferSetWord(u.UpdateDataPlayer().Player.C(), 2136, 7)
								objectXferSetWord(u.UpdateDataPlayer().Player.C(), 2140, 11)
							}
							if friendly {
								for i := 0; i < 2; i++ {
									u := &o.units[i]
									legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
								}
							}
							o.reset()
							*words["players"], *words["events"] = 3, 0
							clear(events)
							if enabled {
								*o.words["rateMode"] = 1
							}
							var killer, assistant *server.Object
							var info *server.Player
							if actor >= 0 {
								killer = &o.units[actor]
							}
							if assist {
								assistant = &o.units[2]
								info = assistant.UpdateDataPlayer().Player
							}
							legacy.PortTestPlayerDeath(op, &o.units[0], killer, assistant, info)
							r := row{Name: name}
							for i := uint32(0); i < *words["events"]; i++ {
								if i >= 8 {
									t.Fatal("unexpected event count")
								}
								r.Pairs = append(r.Pairs, [2]byte{events[2*i], events[2*i+1]})
							}
							var want [][2]byte
							if enabled {
								if actor >= 0 {
									target := byte(0)
									if actor == 0 || friendly {
										target = byte(actor)
									}
									want = append(want, [2]byte{byte(actor), target})
								} else if op == "elimination" {
									want = append(want, [2]byte{0, 0})
								}
								if op == "arena" && assist {
									want = append(want, [2]byte{2, 0})
								}
							}
							if !reflect.DeepEqual(r.Pairs, want) {
								t.Fatalf("statistics pairs %v want %v", r.Pairs, want)
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-death-statistics", rows, "c6f280903ae925e7c3b4dc37b0ba6f0516abf6a8b43a119c79955dfbfdbb1218")
}
