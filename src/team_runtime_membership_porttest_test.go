//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"slices"
	"testing"
)

// Uses real create/unlink/move entry points. Pointer identities are represented by
// fixture node numbers; list order, raw counters, masks and queue data remain exact.
func TestTeamRuntimeMembership(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name    string
		Members []int
		Count   uint32
		IDs     []byte
		Masks   [][4]uint32
		State   legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "team-runtime-membership", rows, "474835eaeaa59237441269061c77b3b4027e8ac34828bf01ab14aa844318ab7e")
	}()
	for _, remove := range []int{-1, 0, 1, 2} {
		for _, duplicate := range []bool{false, true} {
			label := fmt.Sprintf("remove%d/duplicate%t", remove, duplicate)
			t.Run(label, func(t *testing.T) {
				defer noxflags.PortTestGameFlags(1)()
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				tm := o.s.Teams.Create(1)
				for i := range o.units {
					o.units[i].TeamVal = server.ObjectTeam{}
					o.units[i].UpdateDataPlayer().Player.NetCodeVal = o.units[i].NetCode
				}
				ids := map[uint32]int{}
				for i := range o.units {
					u := &o.units[i]
					ids[uint32(uintptr(u.TeamPtr().C()))] = i
					legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
				}
				if duplicate {
					u := &o.units[1]
					legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
				}
				if got := objectXferGetWord(tm.C(), 48); got != uint32(len(o.units)) {
					t.Fatalf("create count %d", got)
				}
				o.reset()
				for i := range o.units {
					o.units[i].Field38 = 0
					o.units[i].Field37 = 0xffffffff
				}
				if remove >= 0 {
					legacy.PortTestTeamRuntimeUnlink(tm, o.units[remove].TeamPtr())
				}
				var members []int
				for p := objectXferGetWord(tm.C(), 44); p != 0; {
					i, ok := ids[p]
					if !ok || len(members) >= len(o.units) {
						t.Fatal("membership chain invalid")
					}
					members = append(members, i)
					p = o.units[i].TeamVal.Field0
				}
				var want []int
				for i := len(o.units) - 1; i >= 0; i-- {
					if i != remove {
						want = append(want, i)
					}
				}
				if !slices.Equal(members, want) {
					t.Fatalf("member order %v want %v", members, want)
				}
				count := objectXferGetWord(tm.C(), 48)
				if count != uint32(len(o.units)) {
					t.Fatal("raw unlink must leave count to caller")
				}
				r := row{Name: label, Members: members, Count: count, State: o.state()}
				for i := range o.units {
					u := &o.units[i]
					wantID := server.TeamID(1)
					if i == remove {
						wantID = 0
						if u.TeamVal.Field0 != 0 {
							t.Fatal("removed next link")
						}
					}
					if u.TeamVal.ID != wantID {
						t.Fatal("member ID")
					}
					r.IDs = append(r.IDs, byte(u.TeamVal.ID))
					r.Masks = append(r.Masks, [4]uint32{u.Field35, u.Field36, u.Field37, u.Field38})
				}
				rows = append(rows, r)
			})
		}
	}
}
