//go:build porttest

package opennox

import (
	"bytes"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unicode/utf16"
)

func TestPortVotesNativeNameBound(t *testing.T) {
	o := newReliableReportsOwner(t)
	name := "abcdefghijklmnopqrstuvwxyz!"
	o.units[0].UpdateDataPlayer().Player.SetName(name)
	text := append(utf16.Encode([]rune(name)), 0)
	for _, withdraw := range []bool{false, true} {
		o.s.NetList.ResetAll()
		legacy.PortTestVoteNameMessage(&text[0], withdraw)
		var packets [][]byte
		o.s.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { packets = append(packets, bytes.Clone(b)); return false })
		want := make([]byte, 52)
		want[0] = 238
		if withdraw {
			want[1] = 2
		}
		for i, c := range text[:24] {
			want[2+2*i] = byte(c)
			want[3+2*i] = byte(c >> 8)
		}
		if len(packets) != 1 || !bytes.Equal(packets[0], want) {
			t.Fatalf("bounded name message %x", packets)
		}
	}
}

func TestPortVotesNativePoolLifetime(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(legacy.PortTestVoteOwner())
	for pass := 0; pass < 3; pass++ {
		if legacy.PortTestVoteCreate(0, o.units[0].CObj()) == nil {
			t.Fatal("live pool allocation")
		}
		legacy.Sub_506720()
		legacy.Sub_506720()
		if len(legacy.PortTestVoteList()) != 0 || legacy.PortTestVoteCreate(0, o.units[0].CObj()) != nil {
			t.Fatal("closed pool state")
		}
		if legacy.Nox_xxx_allocVoteArray_5066D0() != 1 {
			t.Fatal("pool restart")
		}
	}
}

func TestPortVotesNativeMissingLocalTeam(t *testing.T) {
	o, words, w := newVoteGUIOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	o.c.srv.Teams.Create(1)
	legacy.PortTestVoteGUI("show", 0)
	if *words["topic"] != 0 || len(teamUIRowNames(w.ChildByID(4320))) != 0 || !w.GetFlags().IsHidden() {
		t.Fatal("missing local team admission")
	}
}
