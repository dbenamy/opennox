//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// The client handles a draw through the generic flag-result message (87, team
// 65535, time-limit marker 1). The flagball winner message requires an actual team.
func TestMatchRosterFlagballDraw(t *testing.T) {
	for _, tc := range []struct {
		name   string
		scores []int
		winner uint16
	}{
		{"empty", nil, 0}, {"tied", []int{3, 3}, 0}, {"winner", []int{3, 1}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := newQuestRuntimeOwner(t)
			t.Cleanup(o.s.PortTestObjectiveTypes(nil, nil, nil))
			t.Cleanup(o.s.PortTestMapDrawableTeamMessages())
			t.Cleanup(noxflags.PortTestGameFlags(64))
			for i, score := range tc.scores {
				tm := o.s.Teams.Create(server.TeamID(i + 1))
				if tm == nil {
					t.Fatal("team allocation")
				}
				tm.Lessons = score
			}
			o.reset()
			legacy.PortTestMatchRosterFlagWinner()
			state := o.state()
			if len(state.Nodes) != 1 {
				t.Fatalf("result count %d, want 1", len(state.Nodes))
			}
			want := []byte{87, 255, 255, 1, 0, 0, 0, 0}
			if tc.winner != 0 {
				want[0] = 86
				want[1] = byte(tc.winner)
				want[2] = byte(tc.winner >> 8)
				want[3] = 0
			}
			binary.LittleEndian.PutUint32(want[4:], o.s.Frame())
			if !bytes.Equal(state.Nodes[0].Data, want) {
				t.Fatalf("result %x, want %x", state.Nodes[0].Data, want)
			}
			if !noxflags.HasGame(8) {
				t.Fatal("match did not enter its ended state")
			}
		})
	}
}
