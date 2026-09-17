//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestTeamRuntimeClearAndReset(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	tm := o.s.Teams.Create(1)
	o.s.Teams.Create(2)
	for i := range o.units {
		u := &o.units[i]
		u.TeamVal = server.ObjectTeam{}
		u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
		legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
	}
	o.reset()
	legacy.PortTestTeamRuntimeOther("clear", tm, nil, 0, 0)
	if objectXferGetWord(tm.C(), 44) != 0 || objectXferGetWord(tm.C(), 48) != 0 {
		t.Fatal("clear must update the caller-owned counter")
	}
	for i := range o.units {
		if o.units[i].TeamVal.ID != 0 || o.units[i].TeamVal.Field0 != 0 {
			t.Fatal("clear player membership")
		}
	}
	first := o.state()
	o.reset()
	legacy.PortTestTeamRuntimeOther("clear", tm, nil, 0, 0)
	if len(o.state().Nodes) != 0 {
		t.Fatal("repeated clear sent data")
	}
	noxflags.SetGamePlay(4)
	rv := legacy.PortTestTeamRuntimeOther("reset", nil, nil, 0, 0)
	if rv != 0 || o.s.Teams.Count() != 0 || noxflags.HasGamePlay(4) {
		t.Fatal("reset team lifetime/gameplay flag")
	}
	spellbookCapture(t, "team-runtime-clear-reset", struct {
		Clear, Reset legacy.PortTestReliableReportState
	}{first, o.state()}, "eaae58a3fe8929bb1467125754e3016838375eded10687a03a065df5cd82dd22")
}

func TestTeamRuntimeCapacity(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	type row struct{ Count, Available, After int }
	var rows []row
	for count := 0; count <= 16; count++ {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			o.s.Teams.Reset()
			o.s.Teams.ActiveCnt = 0
			for i := 0; i < count; i++ {
				tm := o.s.Teams.Create(server.TeamID(i + 1))
				if tm == nil {
					t.Fatal("team capacity owner")
				}
				legacy.PortTestTeamRuntimeOther("group", tm, nil, int(uint32(0x80000000)+uint32(i)), 0)
				if uint32(tm.Ind60()) != 0x80000000+uint32(i) {
					t.Fatal("group width")
				}
			}
			if got := legacy.PortTestTeamRuntimeSelect("group-count", nil); got != count {
				t.Fatal("group count", got, count)
			}
			got := legacy.PortTestTeamRuntimeSelect("available", nil)
			want := count + 1
			after := count + 1
			if count == 16 {
				want = 0
				after = 16
			}
			if got != want || o.s.Teams.Count() != after {
				t.Fatal("available-team allocation", got, want, o.s.Teams.Count(), after)
			}
			rows = append(rows, row{count, got, o.s.Teams.Count()})
		})
	}
	if legacy.PortTestTeamRuntimeSelect("count", nil) != 0 {
		t.Fatal("nil member count")
	}
	a, b := legacy.PortTestTeamRuntimeLinks(nil, nil)
	if a != nil || b != nil {
		t.Fatal("nil list traversal")
	}
	spellbookCapture(t, "team-runtime-capacity", rows, "f1080c8b5c1a93338f851c68f7231bf7e5b749db0b25c66dfa71d043a9c50cf6")
}

func TestTeamRuntimeClearDetachedMember(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	tm := o.s.Teams.Create(1)
	for i := range o.units {
		u := &o.units[i]
		u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
		if i < 2 {
			legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
		} else {
			u.TeamVal.ID = 1
		}
	}
	legacy.PortTestTeamRuntimeOther("clear", tm, nil, 0, 0)
	if objectXferGetWord(tm.C(), 48) != 0 || objectXferGetWord(tm.C(), 44) != 0 {
		t.Fatal("clear decremented for a detached member")
	}
	if o.units[0].TeamVal.ID != 0 || o.units[1].TeamVal.ID != 0 || o.units[2].TeamVal.ID != 1 {
		t.Fatal("clear membership predicate")
	}
}
