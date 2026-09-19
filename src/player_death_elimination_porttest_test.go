//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
)

func playerDeathLesson(code, score, deaths uint32) []byte {
	b := []byte{78, byte(code), byte(code >> 8)}
	b = binary.LittleEndian.AppendUint32(b, score)
	return binary.LittleEndian.AppendUint32(b, deaths)
}
func playerDeathTeamScore(id, score uint32) []byte {
	b := binary.LittleEndian.AppendUint32([]byte{196, 8}, id)
	return binary.LittleEndian.AppendUint32(b, score)
}
func playerDeathMessages(t *testing.T, state legacy.PortTestReliableReportState, want [][]byte) {
	t.Helper()
	var got [][]byte
	// The real queue stores newest first; compare in send/enqueue order.
	for i := len(state.Nodes) - 1; i >= 0; i-- {
		got = append(got, state.Nodes[i].Data)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("score messages got %x want %x", got, want)
	}
}
func TestPlayerDeathElimination(t *testing.T) {
	o := newMatchRosterOwner(t)
	type record struct {
		Name           string
		Scores, Deaths [3]uint32
		Teams          [2]int
		Reports        legacy.PortTestReliableReportState
	}
	var rows []record
	for _, actor := range []int{-1, 0, 1, 2} {
		for victimTeam := 0; victimTeam <= 2; victimTeam++ {
			for killerTeam := 0; killerTeam <= 2; killerTeam++ {
				for _, seed := range []uint32{0, 1, 0x7fffffff, 0xffffffff} {
					name := fmt.Sprintf("actor=%d/victim-team=%d/killer-team=%d/seed=%x", actor, victimTeam, killerTeam, seed)
					t.Run(name, func(t *testing.T) {
						o.reset()
						noxflags.ResetGame()
						noxflags.SetGame(noxflags.GameHost | noxflags.GameModeElimination)
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
						var killer *server.Object
						if actor >= 0 {
							killer = &o.units[actor]
						}
						oldClass := o.units[2].ObjClass
						if actor == 2 {
							o.units[2].ObjClass = 0
						}
						defer func() { o.units[2].ObjClass = oldClass }()
						scores := [3]uint32{seed, seed, seed}
						deaths := [3]uint32{seed + 1, seed, seed}
						ts := [2]int{int(seed), int(seed)}
						var messages [][]byte
						if actor == 0 {
							scores[0]--
						} else if actor == 1 {
							if killerTeam != 0 && killerTeam == victimTeam {
								scores[1]--
							} else {
								scores[1]++
							}
							messages = append(messages, playerDeathLesson(o.units[1].NetCode, scores[1], seed))
						}
						messages = append(messages, playerDeathLesson(o.units[0].NetCode, scores[0], deaths[0]))
						if victimTeam != 0 {
							ts[victimTeam-1] = int(seed + 1)
							messages = append(messages, playerDeathTeamScore(uint32(victimTeam), seed+1))
						}
						legacy.PortTestPlayerDeath("elimination", &o.units[0], killer, nil, nil)
						o.units[2].ObjClass = oldClass // Restore the fixture type before reading its player snapshot.
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
	spellbookCapture(t, "player-death-elimination", rows, "b1616b25cd55a28ade72129055ec77c14967018064a8baaf69cba74c661efeea")
}
