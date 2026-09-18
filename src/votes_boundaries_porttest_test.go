//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPortVotesCapacityAndReuse(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	var nodes []unsafe.Pointer
	for i := 0; i < 64; i++ {
		p := legacy.PortTestVoteCreate(0, o.units[0].CObj())
		if p == nil {
			t.Fatalf("capacity %d", i)
		}
		nodes = append(nodes, p)
	}
	if legacy.PortTestVoteCreate(0, o.units[0].CObj()) != nil {
		t.Fatal("exceeded fixed pool capacity")
	}
	legacy.PortTestVoteDelete(nodes[20])
	p := legacy.PortTestVoteCreate(3, o.units[0].CObj())
	if p == nil {
		t.Fatal("freed slot not reused")
	}
	if len(legacy.PortTestVoteList()) != 64 {
		t.Fatal("reused list size")
	}
	for _, off := range []int{4, 8, 20, 28, 32, 36, 40} {
		if *(*uint32)(unsafe.Add(p, off)) != 0 {
			t.Fatalf("reused field %d", off)
		}
	}
}
func TestPortVotesTeamThreshold(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	p := legacy.PortTestVoteCreate(0, o.units[0].CObj())
	*(*byte)(unsafe.Add(p, 12)) = 5
	*(*uint32)(unsafe.Add(p, 20)) = 1
	for i := range o.units {
		o.units[i].TeamVal.ID = 1
	}
	for _, active := range []int{0, 1, 2, 3} {
		for i := range o.units {
			u := &o.units[i]
			u.UpdateDataPlayer().Player.Active = 0
			if i < active {
				u.UpdateDataPlayer().Player.Active = 1
			}
		}
		for _, team := range []byte{0, 1, 2} {
			o.units[0].TeamVal.ID = 1
			o.units[1].TeamVal.ID = server.TeamID(team)
			o.units[2].TeamVal.ID = server.TeamID(team)
			eligible := 0
			if active > 0 {
				eligible = 1
			}
			if team == 1 {
				eligible = active
			}
			for _, count := range []byte{0, 1, 2, 4, 5, 255} {
				*(*byte)(unsafe.Add(p, 4)) = count
				want := 0
				if count >= 5 || eligible > 0 && int(count) >= eligible-1 && count >= 2 {
					want = 1
				}
				if got := legacy.PortTestVoteThreshold(p); got != want {
					t.Fatalf("active%d team%d count%d got%d want%d", active, team, count, got, want)
				}
			}
		}
	}
}
func TestPortVotesQuestWithdrawal(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	settings, restore := legacy.PortTestVoteSettings()
	t.Cleanup(restore)
	*settings[0], *settings[1] = 6, 5
	for i := range o.units {
		o.units[i].UpdateDataPlayer().Player.Field4792 = 1
	}
	o.units[1].UpdateDataPlayer().Player.SetName("target")
	name := []uint16{'t', 'a', 'r', 'g', 'e', 't', 0}
	for _, kind := range []int{2, 3} {
		for _, i := range []int{0, 2, 0} {
			legacy.PortTestVoteCast(kind, o.units[i].CObj(), &name[0], false)
		}
		list := legacy.PortTestVoteList()
		if len(list) != 1 || *(*byte)(unsafe.Add(list[0], 4)) != 2 {
			t.Fatal("quest duplicate")
		}
		for _, i := range []int{0, 0, 2, 2} {
			legacy.PortTestVoteCast(kind, o.units[i].CObj(), &name[0], true)
		}
		if len(legacy.PortTestVoteList()) != 0 {
			t.Fatal("quest last withdrawal")
		}
	}
}
func TestPortVotesQuestTargetGuards(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	t.Cleanup(noxflags.PortTestGameFlags(0))
	target := &o.units[1]
	for _, state := range []int{0, 1, 2, 3} {
		target.ObjFlags = 0
		target.UpdateDataPlayer().Player.Field4792 = 1
		p := legacy.PortTestVoteCreate(3, o.units[0].CObj())
		*(*unsafe.Pointer)(unsafe.Add(p, 28)) = target.CObj()
		switch state {
		case 0:
			*(*unsafe.Pointer)(unsafe.Add(p, 28)) = nil
		case 1:
			target.ObjFlags = object.FlagDestroyed
		case 2:
			target.UpdateDataPlayer().Player.Field4792 = 0
		}
		legacy.Nox_xxx_voteUptate_506F30()
		// With only the target participating, the vote is impossible and removed.
		if len(legacy.PortTestVoteList()) != 0 {
			t.Fatalf("guard%d retained vote", state)
		}
	}
}

func TestPortVotesTeamAdmissionAndByteCount(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	old := noxflags.GetGamePlay()
	noxflags.UnsetGamePlay(old)
	noxflags.SetGamePlay(noxflags.GameplayFlag4)
	t.Cleanup(func() { noxflags.UnsetGamePlay(noxflags.GetGamePlay()); noxflags.SetGamePlay(old) })
	o.units[1].UpdateDataPlayer().Player.SetName("target")
	name := []uint16{'t', 'a', 'r', 'g', 'e', 't', 0}
	var rows [][13]uint32
	for _, ids := range [][2]byte{{0, 0}, {0, 1}, {1, 0}, {1, 2}, {1, 1}, {255, 255}} {
		o.units[0].TeamVal.ID = server.TeamID(ids[0])
		o.units[1].TeamVal.ID = server.TeamID(ids[1])
		legacy.PortTestVoteCast(0, o.units[0].CObj(), &name[0], false)
		list := legacy.PortTestVoteList()
		want := ids[0] != 0 && ids[0] == ids[1]
		if (len(list) != 0) != want {
			t.Fatalf("team admission %v", ids)
		}
		if !want {
			continue
		}
		p := list[0]
		if *(*uint32)(unsafe.Add(p, 20)) != 1 {
			t.Fatal("team filter not stored")
		}
		rows = append(rows, voteRecordWords(p, o.objects))
		// Count arithmetic is byte-sized even when an existing record is unusual.
		*(*byte)(unsafe.Add(p, 4)) = 255
		*(*uint32)(unsafe.Add(p, 8)) = 0
		legacy.PortTestVoteCast(0, o.units[0].CObj(), &name[0], false)
		if *(*byte)(unsafe.Add(p, 4)) != 0 || *(*uint32)(unsafe.Add(p, 8)) != 2 {
			t.Fatal("byte count wrap")
		}
		rows = append(rows, voteRecordWords(p, o.objects))
		legacy.PortTestVoteDelete(p)
	}
	spellbookCapture(t, "votes-team-admission", rows, "8bbb494632e2213b1ee0b222acb1babeb4d48fadcd2c6531fcb805316f8331c3")
}
