//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerDeathArenaTeams(t *testing.T) {
	o := newMatchRosterOwner(t)
	type record struct {
		Name           string
		Scores, Deaths [3]uint32
		Teams          [2]int
		Reports        legacy.PortTestReliableReportState
	}
	var rows []record
	for victimTeam := 0; victimTeam <= 2; victimTeam++ {
		for killerTeam := 0; killerTeam <= 2; killerTeam++ {
			name := fmt.Sprintf("victim-team=%d/killer-team=%d", victimTeam, killerTeam)
			t.Run(name, func(t *testing.T) {
				o.reset()
				noxflags.ResetGame()
				noxflags.SetGame(noxflags.GameHost | noxflags.GameModeArena)
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
				for _, tm := range teams {
					tm.Lessons = 99
				}
				o.reset()
				wantScores := [3]uint32{7, 8, 7}
				wantTeams := [2]int{99, 99}
				if killerTeam != 0 {
					if killerTeam == victimTeam {
						wantScores[1] = 6
						wantTeams[killerTeam-1]--
					} else {
						wantTeams[killerTeam-1]++
					}
				}
				legacy.PortTestPlayerDeath("arena", &o.units[0], &o.units[1], nil, nil)
				row := record{Name: name, Teams: [2]int{teams[0].Lessons, teams[1].Lessons}, Reports: o.state()}
				for i := range o.units {
					pl := o.units[i].UpdateDataPlayer().Player
					row.Scores[i] = objectXferGetWord(pl.C(), 2136)
					row.Deaths[i] = objectXferGetWord(pl.C(), 2140)
				}
				if row.Scores != wantScores || row.Teams != wantTeams || row.Deaths != [3]uint32{12, 11, 11} {
					t.Fatalf("scores %v deaths %v teams %v want scores %v deaths [12 11 11] teams %v", row.Scores, row.Deaths, row.Teams, wantScores, wantTeams)
				}
				messages := [][]byte{playerDeathLesson(o.units[1].NetCode, wantScores[1], 11)}
				if killerTeam != 0 {
					messages = append(messages, playerDeathTeamScore(uint32(killerTeam), uint32(wantTeams[killerTeam-1])))
				}
				messages = append(messages, playerDeathLesson(o.units[0].NetCode, 7, 12))
				playerDeathMessages(t, row.Reports, messages)
				rows = append(rows, row)
			})
		}
	}
	spellbookCapture(t, "player-death-arena-teams", rows, "")
}
