//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
)

func TestPlayerStateAllStatus(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name  string
		Queue legacy.PortTestReliableReportState
	}
	var rows []row
	for _, to := range []int{1, 7, 31} {
		for _, mask := range []uint32{0, 0x423, 0xffffffff, 0x80000000} {
			o.reset()
			players := o.s.Players.List()
			for i, p := range players {
				p.NetCodeVal = 0x12340000 + uint32(17*i)
				p.Field3680 = mask ^ uint32(i)
			}
			legacy.PortTestPlayerStateStatus("all", o.s.Players.ByIndRaw(ntype.PlayerInd(to)), 0)
			state := o.state()
			if len(state.Nodes) != len(players) {
				t.Fatal("roster report count", len(state.Nodes), len(players))
			}
			for i, n := range state.Nodes {
				p := players[len(players)-1-i]
				want := make([]byte, 7)
				want[0] = 106
				binary.LittleEndian.PutUint16(want[1:], uint16(p.NetCodeVal))
				binary.LittleEndian.PutUint32(want[3:], p.Field3680&0x423)
				if !bytes.Equal(n.Data, want) || int(n.To) != to || n.Recipients != 0x80000082 {
					t.Fatal("roster report", i, n, want)
				}
			}
			rows = append(rows, row{fmt.Sprintf("to%d/mask%x", to, mask), state})
		}
	}
	spellbookCapture(t, "player-state-all-status", rows, "fd7515a95162a5b259e0277bf4c8cac080ca06719cd749a68fcb6adfbe237d9e")
}
func TestPlayerStateLessons(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Send  uint32
		Queue legacy.PortTestReliableReportState
	}
	var rows []row
	for _, send := range []uint32{0, 1, 0xffffffff} {
		o.reset()
		defer noxflags.PortTestGameFlags(1)()
		for _, p := range o.s.Players.List() {
			p.Lessons = -123
			p.Field2140 = 0xffffffff
		}
		legacy.PortTestPlayerStateStatus("lessons", nil, send)
		for _, p := range o.s.Players.List() {
			if p.Lessons != 0 || p.Field2140 != 0 {
				t.Fatal("score reset", p.PlayerInd)
			}
		}
		state := o.state()
		if len(state.Nodes) != 3*bool2int(send != 0) {
			t.Fatal("lesson report count", len(state.Nodes))
		}
		for i, n := range state.Nodes {
			code := o.units[len(o.units)-1-i].NetCode
			want := []byte{78, byte(code), byte(code >> 8), 0, 0, 0, 0, 0, 0, 0, 0}
			if !bytes.Equal(n.Data, want) || n.To != 255 || n.Priority != 1 || n.Ordered != 1 {
				t.Fatal("lesson report", n, want)
			}
		}
		rows = append(rows, row{send, state})
	}
	spellbookCapture(t, "player-state-lessons", rows, "ac639dad790ab66e5a0384e6dadcc48ead50096d8474a2a7dce2b1ae02f2cf77")
}
func TestPlayerStateName(t *testing.T) {
	o := newMatchRosterOwner(t)
	for _, p := range o.s.Players.List() {
		p.SetName("unused")
	}
	o.s.Players.ByInd(1).SetName("Alice")
	o.s.Players.ByInd(7).SetName("ALICE")
	o.s.Players.ByInd(31).SetName("Mage Ω")
	type row struct {
		Name  string
		Index int
	}
	rows := []row{}
	if legacy.PortTestPlayerStateName(nil) != nil {
		t.Fatal("nil name")
	}
	for _, c := range []struct {
		name  string
		index int
	}{{"Alice", 1}, {"aLiCe", 1}, {"ALICE", 1}, {"unused", 3}, {"Mage Ω", 31}, {"Mage ω", -1}, {"", -1}, {"missing", -1}} {
		name, free := alloc.CString16(c.name)
		pl := legacy.PortTestPlayerStateName(name)
		free()
		got := -1
		if pl != nil {
			got = int(pl.PlayerInd)
		}
		if got != c.index {
			t.Fatal(c.name, got, c.index)
		}
		rows = append(rows, row{c.name, got})
	}
	spellbookCapture(t, "player-state-name", rows, "c145961b9f3244dd556195174c4d8d69013c2c4d6f248c15c3139fa228869fc6")
}
