//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerDeathArenaAssists(t *testing.T) {
	o := newMatchRosterOwner(t)
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
				for assistTeam := 0; assistTeam <= 2; assistTeam++ {
					for _, assistPresent := range []bool{false, true} {
						if !assistPresent && assistTeam != 0 {
							continue
						}
						for _, seed := range []uint32{0, 7, 0xffffffff} {
							name := fmt.Sprintf("actor=%d/teams=%d,%d,%d/assist=%t/seed=%x", actor, victimTeam, killerTeam, assistTeam, assistPresent, seed)
							t.Run(name, func(t *testing.T) {
								o.reset()
								noxflags.ResetGame()
								noxflags.SetGame(noxflags.GameHost | noxflags.GameModeArena)
								o.s.Teams.Reset()
								o.s.Teams.ActiveCnt = 0
								teams := [2]*server.Team{o.s.Teams.Create(1), o.s.Teams.Create(2)}
								ids := [3]int{victimTeam, killerTeam, assistTeam}
								for i := range o.units {
									u := &o.units[i]
									u.TeamVal = server.ObjectTeam{}
									pl := u.UpdateDataPlayer().Player
									pl.PlayerUnit = u
									pl.NetCodeVal = u.NetCode
									objectXferSetWord(pl.C(), 2136, seed)
									objectXferSetWord(pl.C(), 2140, seed)
									if ids[i] != 0 {
										legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(ids[i]), u.TeamPtr(), 0, int(u.NetCode), 0)
									}
								}
								for _, tm := range teams {
									tm.Lessons = int(seed)
								}
								o.reset()
								var killer, assist *server.Object
								var info *server.Player
								if actor >= 0 {
									killer = &o.units[actor]
								}
								if assistPresent {
									assist = &o.units[2]
									info = assist.UpdateDataPlayer().Player
								}
								scores := [3]uint32{seed, seed, seed}
								deaths := [3]uint32{seed, seed, seed}
								ts := [2]int{int(seed), int(seed)}
								var messages [][]byte
								lesson := func(i int) { messages = append(messages, playerDeathLesson(o.units[i].NetCode, scores[i], deaths[i])) }
								teamScore := func(id, delta int) {
									if id == 0 {
										return
									}
									ts[id-1] += delta
									messages = append(messages, playerDeathTeamScore(uint32(id), uint32(ts[id-1])))
								}
								if actor == 0 || actor < 0 && !assistPresent {
									// Historical arena policy: self/environment penalties do not add a death.
									scores[0]--
									lesson(0)
									teamScore(victimTeam, -1)
								} else {
									if actor == 1 {
										delta := 1
										if killerTeam != 0 && killerTeam == victimTeam {
											delta = -1
										}
										scores[1] += uint32(delta)
										lesson(1)
										teamScore(killerTeam, delta)
									}
									deaths[0]++
									lesson(0)
								}
								effectiveKillerTeam := killerTeam
								if actor == 0 {
									effectiveKillerTeam = victimTeam
								}
								if assistPresent && (assistTeam == 0 || assistTeam != victimTeam && assistTeam != effectiveKillerTeam) {
									scores[2]++
									lesson(2)
									teamScore(assistTeam, 1)
								}
								legacy.PortTestPlayerDeath("arena", &o.units[0], killer, assist, info)
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
								rows = append(rows, row)
							})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-death-arena-assists", rows, "93f3684188a03dc47d951c6dc1555942f1f178192853e955fe3ed04b8f77b1c7")
}
