//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestSessionEntryOperatorRoster(t *testing.T) {
	newMatchRosterOwner(t)
	t.Cleanup(legacy.PortTestSessionEntryRosterOwner())
	type row struct {
		Op      string
		Index   int
		Members []int
		Queries [7]int
	}
	var rows []row
	for _, c := range []struct {
		op    string
		index int
		want  []int
	}{
		{"init", 0, []int{}}, {"add", -1, []int{}}, {"add", 32, []int{}}, {"add", 0, []int{}},
		{"add", 1, []int{1}}, {"add", 7, []int{1, 7}}, {"add", 1, []int{1, 7}},
		{"add", 31, []int{1, 7, 31}}, {"add", 3, []int{1, 7, 31, 3}}, {"init", 0, []int{1, 7, 31, 3}},
		{"remove", 7, []int{1, 31, 3}}, {"remove", 1, []int{31, 3}}, {"remove", 31, []int{3}},
		{"remove", 3, []int{}}, {"remove", 7, []int{}}, {"add", 7, []int{7}},
		{"clear", 0, []int{}}, {"clear", 0, []int{}}, {"add", 1, []int{1}},
	} {
		legacy.PortTestSessionEntryRoster(c.op, c.index)
		r := row{Op: c.op, Index: c.index, Members: legacy.PortTestSessionEntryRosterMembers()}
		if !reflect.DeepEqual(r.Members, c.want) {
			t.Fatalf("%s %d members%v want%v", c.op, c.index, r.Members, c.want)
		}
		for i, index := range []int{-1, 0, 1, 3, 7, 31, 32} {
			want := 0
			for _, member := range c.want {
				if member == index {
					want = 1
				}
			}
			r.Queries[i] = legacy.PortTestSessionEntryRoster("contains", index)
			if r.Queries[i] != want {
				t.Fatal("roster membership", index, r.Queries[i], want)
			}
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "session-entry-operator-roster", rows, "037ca86e9884eec01892104fdacf08504d53c333b8d773ebb4c81edac5555438")
}
