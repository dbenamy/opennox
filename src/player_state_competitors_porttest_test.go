//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestPlayerStateCompetitors(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name                  string
		Competitors, Multiple int
		TeamPlayers           [3]int
	}
	var rows []row
	for _, teams := range []bool{false, true} {
		for _, ids := range [][3]byte{{1, 1, 1}, {1, 2, 2}, {1, 2, 3}, {0, 0, 0}} {
			for _, status := range [][3]uint32{{0, 0, 0}, {1, 0, 0}, {1, 1, 1}, {0x21, 1, 0x20}} {
				for _, headless := range []bool{false, true} {
					for _, missing := range []bool{false, true} {
						name := fmt.Sprintf("teams%t/ids%v/status%v/headless%t/missing%t", teams, ids, status, headless, missing)
						t.Run(name, func(t *testing.T) {
							defer noxflags.PortTestGameFlags(noxflags.GameHost)()
							o.reset()
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							var tm [3]*server.Team
							for i := range tm {
								tm[i] = o.s.Teams.Create(server.TeamID(i + 1))
							}
							noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
							if teams {
								noxflags.SetGamePlay(4)
							}
							noxflags.ResetEngine()
							if headless {
								noxflags.SetEngine(noxflags.EngineNoRendering)
							}
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								p := u.UpdateDataPlayer().Player
								p.PlayerUnit = u
								p.NetCodeVal = u.NetCode
								p.Field3680 = status[i]
								if ids[i] != 0 {
									legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(ids[i]), u.TeamPtr(), 0, int(u.NetCode), 0)
								}
							}
							if missing {
								o.units[1].UpdateDataPlayer().Player.PlayerUnit = nil
							}
							defer func() {
								for i := range o.units {
									o.units[i].UpdateDataPlayer().Player.PlayerUnit = &o.units[i]
								}
							}()
							eligible := 0
							activeUnits := 0
							occupied := map[byte]bool{}
							participatingTeams := map[byte]bool{}
							var wantTeam [3]int
							for i := range o.units {
								p := o.units[i].UpdateDataPlayer().Player
								if p.Field3680&1 == 0 && (!headless || p.PlayerInd != 31) {
									eligible++
									if p.PlayerUnit != nil && ids[i] != 0 {
										occupied[ids[i]] = true
										wantTeam[ids[i]-1]++
									}
								}
								if p.PlayerUnit != nil && (p.Field3680&1 == 0 || p.Field3680&0x20 != 0) {
									activeUnits++
									if ids[i] != 0 {
										participatingTeams[ids[i]] = true
									}
								}
							}
							// The shared sparse owner deliberately includes active slot3 with no unit.
							// Player-record counts include it; unit/team queries must not.
							extra := o.s.Players.ByInd(3)
							if extra == nil || extra.Active != 1 || extra.PlayerUnit != nil || extra.Field3680 != 0 {
								t.Fatal("unexpected nil-unit fixture player")
							}
							wantCount := eligible + 1
							wantMultiple := activeUnits > 1
							if teams {
								wantCount = len(occupied)
								wantMultiple = len(participatingTeams) > 1
							}
							r := row{Name: name, Competitors: legacy.PortTestPlayerStateQuery("competitors", nil, nil), Multiple: legacy.PortTestPlayerStateQuery("multiple", nil, nil)}
							for i := range tm {
								r.TeamPlayers[i] = legacy.PortTestPlayerStateQuery("team-players", nil, tm[i])
							}
							if r.Competitors != wantCount || r.Multiple != bool2int(wantMultiple) || r.TeamPlayers != wantTeam {
								t.Errorf("competitors=%d want%d multiple=%d want%t team=%v want%v", r.Competitors, wantCount, r.Multiple, wantMultiple, r.TeamPlayers, wantTeam)
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "player-state-competitors", rows, "8d2a7d504ddda0d77d9908dedeaf43a737bda956c1a2fc9e3f3174697cb45afe")
}
