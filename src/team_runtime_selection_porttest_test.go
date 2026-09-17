//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestTeamRuntimeSelection(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	type row struct {
		Assignment               int
		Counts                   [3]int
		Least, Available, Groups int
	}
	var rows []row
	for assignment := 0; assignment < 27; assignment++ {
		t.Run(fmt.Sprint(assignment), func(t *testing.T) {
			o.s.Teams.Reset()
			o.s.Teams.ActiveCnt = 0
			var teams [3]*server.Team
			for i := range teams {
				teams[i] = o.s.Teams.Create(server.TeamID(i + 1))
			}
			for i := range o.units {
				o.units[i].TeamVal = server.ObjectTeam{}
				o.units[i].UpdateDataPlayer().Player.NetCodeVal = o.units[i].NetCode
			}
			var counts [3]int
			v := assignment
			for i := range o.units {
				index := v % 3
				v /= 3
				counts[index]++
				u := &o.units[i]
				legacy.Nox_xxx_createAtImpl_4191D0(teams[index].ID(), u.TeamPtr(), 0, int(u.NetCode), 0)
			}
			least := 0
			for i, tm := range teams {
				if got := legacy.PortTestTeamRuntimeSelect("count", tm); got != counts[i] {
					t.Fatalf("player count %d want %d", got, counts[i])
				}
				if counts[i] < counts[least] {
					least = i
				}
				objectXferSetWord(tm.C(), 60, uint32(i&1))
			}
			gotLeast := legacy.PortTestTeamRuntimeSelect("least", nil)
			available := legacy.PortTestTeamRuntimeSelect("available", nil)
			groups := legacy.PortTestTeamRuntimeSelect("group-count", nil)
			if gotLeast != least+1 || available != 1 || groups != 1 {
				t.Fatalf("team selection %d/%d/%d", gotLeast, available, groups)
			}
			rows = append(rows, row{assignment, counts, gotLeast, available, groups})
		})
	}
	spellbookCapture(t, "team-runtime-selection", rows, "9edb675f51b5fc22d66858b4010560c3731d49e0df022b43a1359ddc88a8fb58")
}

func TestTeamRuntimeEmptyFlagMaps(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	for _, op := range []string{"capflag", "flagball"} {
		t.Run(op, func(t *testing.T) {
			*o.roster["team-cap"] = 99
			if got := legacy.PortTestTeamRuntimeSelect(op, nil); got != 0 {
				t.Fatalf("empty flag map result %d, want zero", got)
			}
			if got := legacy.PortTestTeamRuntimeSelect("flag-count", nil); got != 0 {
				t.Fatalf("empty flag map count %d", got)
			}
		})
	}
}
