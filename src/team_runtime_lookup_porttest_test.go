//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestTeamRuntimeMembershipPredicates(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	a := o.s.Teams.Create(1)
	b := o.s.Teams.Create(2)
	first := (*server.ObjectTeam)(o.record(t, 8))
	second := (*server.ObjectTeam)(o.record(t, 8))
	detached := (*server.ObjectTeam)(o.record(t, 8))
	zero := (*server.ObjectTeam)(o.record(t, 8))
	legacy.Nox_xxx_createAtImpl_4191D0(a.ID(), first, 0, 10001, 0)
	legacy.Nox_xxx_createAtImpl_4191D0(b.ID(), second, 0, 10002, 0)
	detached.ID = 1
	nodes := []*server.ObjectTeam{nil, zero, first, second, detached}
	type row struct{ I, J, ID, Has, Same, Member int }
	var rows []row
	for i, x := range nodes {
		for j, y := range nodes {
			for _, id := range []int{0, 1, 2, 3, 255} {
				t.Run(fmt.Sprintf("%d/%d/%d", i, j, id), func(t *testing.T) {
					has, same, member := legacy.PortTestTeamRuntimePredicates(x, y, byte(id))
					wantHas := 0
					if x != nil && x.ID != 0 {
						wantHas = 1
					}
					wantSame := 0
					if x != nil && y != nil && x.ID != 0 && x.ID == y.ID {
						wantSame = 1
					}
					wantMember := 0
					if x == first && id == 1 || x == second && id == 2 {
						wantMember = 1
					}
					if has != wantHas || same != wantSame || member != wantMember {
						t.Fatalf("predicates %d/%d/%d want %d/%d/%d", has, same, member, wantHas, wantSame, wantMember)
					}
					rows = append(rows, row{i, j, id, has, same, member})
				})
			}
		}
	}
	spellbookCapture(t, "team-runtime-predicates", rows, "f6e13e3551b803b158d7579bed1c72631dae4ff9af75339d815f1ce988f4f264")
}

func TestTeamRuntimeLookup(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name   string
		Found  int
		Exists bool
	}
	var rows []row
	names := []string{"Red", "BLUE", "Team Ω"}
	for i, name := range names {
		tm := o.s.Teams.Create(server.TeamID(i + 1))
		legacy.PortTestTeamRuntimeName(tm, name, 0)
	}
	for _, tc := range []struct {
		name string
		id   int
	}{{"Red", 1}, {"red", 1}, {"RED", 1}, {"Blue", 2}, {"BLUE", 2}, {"Team Ω", 3}, {"", 0}, {"missing", 0}, {"Red ", 0}} {
		t.Run(tc.name, func(t *testing.T) {
			tm, exists := legacy.PortTestTeamRuntimeLookup(tc.name)
			id := 0
			if tm != nil {
				id = int(tm.ID())
			}
			if id != tc.id || exists != (tc.id != 0) {
				t.Fatal("team lookup")
			}
			rows = append(rows, row{tc.name, id, exists})
		})
	}
	spellbookCapture(t, "team-runtime-lookup", rows, "bc06fb9af5f3459a82a937fbda38f4b0648450f49f46304d81101bd6b090e58a")
}
