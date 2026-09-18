//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestItemRespawnOwnerTeams(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Owners [3]int
		Teams  [3]byte
		A, B   int
		Same   bool
	}
	var rows []row
	// All acyclic forests whose owners have a higher index, every pair including nil,
	// and unassigned/same/different team assignments. Identity is separate from team0.
	for own0 := 0; own0 <= 3; own0++ {
		if own0 == 1 {
			continue
		}
		for own1 := 0; own1 <= 3; own1++ {
			if own1 == 1 || own1 == 2 {
				continue
			}
			owners := [3]int{own0, own1, 0}
			for bits := 0; bits < 27; bits++ {
				n := bits
				var teams [3]byte
				for i := range o.units {
					teams[i] = byte(n % 3)
					n /= 3
					o.units[i].TeamVal = server.ObjectTeam{ID: server.TeamID(teams[i])}
					o.units[i].ObjOwner = nil
					if owners[i] != 0 {
						o.units[i].ObjOwner = &o.units[owners[i]-1]
					}
				}
				for a := 0; a <= 3; a++ {
					for b := 0; b <= 3; b++ {
						var ua, ub *server.Object
						if a != 0 {
							ua = &o.units[a-1]
						}
						if b != 0 {
							ub = &o.units[b-1]
						}
						want := false
						for x := a; x != 0; x = owners[x-1] {
							for y := b; y != 0; y = owners[y-1] {
								if x == y || teams[x-1] != 0 && teams[x-1] == teams[y-1] {
									want = true
								}
							}
						}
						got := legacy.PortTestItemRespawnSameTeam(ua, ub)
						if got != want {
							t.Fatalf("owners%v teams%v pair%d,%d got%v", owners, teams, a, b, got)
						}
						rows = append(rows, row{owners, teams, a, b, got})
					}
				}
			}
		}
	}
	for i := range o.units {
		o.units[i].ObjOwner = nil
	}
	spellbookCapture(t, "item-respawn-owner-teams", rows, "9c5ac1b7c160b8dac43293459b3f1bb828fa3b05098e86533156905f3a52bb3f")
}
