//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerDeathKotrScoring(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown"}, nil, true, 0, 0))
	cache := serverConfigOwnBytes(t, 0x5D4594, 1567716, 4)
	clear(cache)
	var crowns [2]*server.Object
	for i := range crowns {
		crowns[i] = o.s.NewObjectByTypeID("Crown")
		if crowns[i] == nil {
			t.Fatal("crown allocation")
		}
		trackObjectXferTyped(t, o.s, crowns[i])
	}
	t.Cleanup(func() {
		for _, c := range crowns {
			o.s.ObjClearOwner(c)
		}
	})
	type record struct {
		Name           string
		Scores, Deaths [3]uint32
		Teams          [2]int
		Reports        legacy.PortTestReliableReportState
	}
	var rows []record
	for _, actor := range []int{-1, 0, 1} {
		for victimTeam := 0; victimTeam <= 2; victimTeam++ {
			for killerTeam := 0; killerTeam <= 2; killerTeam++ {
				if actor != 1 && killerTeam != 0 {
					continue
				}
				for mask := 0; mask < 4; mask++ {
					if actor != 1 && mask > 1 {
						continue
					}
					for variant, points := range [][3]float64{{2, 7, 11}, {1.75, 3.25, -2.75}} {
						name := fmt.Sprintf("actor=%d/teams=%d,%d/crowns=%d/points=%d", actor, victimTeam, killerTeam, mask, variant)
						t.Run(name, func(t *testing.T) {
							for _, c := range crowns {
								o.s.ObjClearOwner(c)
							}
							o.reset()
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameHost | noxflags.GameModeKOTR)
							noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
							noxflags.SetGamePlay(4)
							o.balance(map[string]float64{"KotRKingKillsPawnPoints": points[0], "KotRPawnKillsKingPoints": points[1], "KotRKingKillsKingPoints": points[2]})
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							teams := [2]*server.Team{o.s.Teams.Create(1), o.s.Teams.Create(2)}
							ids := [3]int{victimTeam, killerTeam, 0}
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								pl := u.UpdateDataPlayer().Player
								pl.PlayerUnit = u
								pl.NetCodeVal = u.NetCode
								objectXferSetWord(pl.C(), 2136, 7)
								objectXferSetWord(pl.C(), 2140, 11)
								if ids[i] != 0 {
									legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(ids[i]), u.TeamPtr(), 0, int(u.NetCode), 0)
								}
							}
							for i, c := range crowns {
								if mask&(1<<i) != 0 {
									o.s.ObjSetOwner(&o.units[i], c)
								}
							}
							for _, tm := range teams {
								tm.Lessons = 99
							}
							o.reset()
							var killer *server.Object
							if actor >= 0 {
								killer = &o.units[actor]
							}
							scores := [3]uint32{7, 7, 7}
							deaths := [3]uint32{11, 11, 11}
							ts := [2]int{99, 99}
							var messages [][]byte
							lesson := func(i int) { messages = append(messages, playerDeathLesson(o.units[i].NetCode, scores[i], deaths[i])) }
							teamScore := func(id, delta int) {
								if id == 0 {
									return
								}
								ts[id-1] += delta
								messages = append(messages, playerDeathTeamScore(uint32(id), uint32(ts[id-1])))
							}
							if actor >= 0 {
								kc := mask&2 != 0
								kt := killerTeam
								if actor == 0 {
									kc = mask&1 != 0
									kt = victimTeam
								}
								vc := mask&1 != 0
								friendly := actor == 0 || kt != 0 && kt == victimTeam
								if friendly {
									if kc {
										scores[actor]--
										lesson(actor)
										teamScore(kt, -1)
									}
								} else if kc || vc {
									index := 1
									if kc {
										index = 0
										if kt != 0 && vc {
											index = 2
										}
									}
									amount := int(int32(points[index]))
									scores[actor] += uint32(amount)
									teamScore(kt, amount)
									lesson(actor)
								}
								deaths[0]++
								lesson(0)
							}
							legacy.PortTestPlayerDeath("kotr", &o.units[0], killer, nil, nil)
							row := record{Name: name, Teams: [2]int{teams[0].Lessons, teams[1].Lessons}, Reports: o.state()}
							for i := range o.units {
								pl := o.units[i].UpdateDataPlayer().Player
								row.Scores[i] = objectXferGetWord(pl.C(), 2136)
								row.Deaths[i] = objectXferGetWord(pl.C(), 2140)
							}
							if row.Scores != scores || row.Deaths != deaths || row.Teams != ts {
								t.Fatalf("score/death/team %v/%v/%v want %v/%v/%v", row.Scores, row.Deaths, row.Teams, scores, deaths, ts)
							}
							playerDeathMessages(t, row.Reports, messages)
							for i, c := range crowns {
								var want *server.Object
								if mask&(1<<i) != 0 {
									want = &o.units[i]
								}
								if c.ObjOwner != want {
									t.Fatal("scoring-only mode changed crown owner")
								}
							}
							rows = append(rows, row)
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-death-kotr-scoring", rows, "")
}
