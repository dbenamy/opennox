//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsRate(t *testing.T) {
	g := memmap.PtrUint32(0x587000, 4728)
	old := *g
	defer func() { *g = old }()
	var out []legacy.PortTestRoamResult
	for _, rate := range []uint32{0, 1, 127, 255, 256, 65535, 0xffffffff} {
		*g = rate
		var cases []legacy.PortTestRoamSpec
		for _, recipient := range []uint32{1, 7, 31} {
			cases = append(cases, gameplayReportsBase(4, reportValue(recipient)))
		}
		results := controlsRun(t, cases)
		for i, r := range results {
			ps := r.Callbacks.Shop.Sequence[0].Packets
			if len(ps) != 1 {
				t.Fatal("rate message count")
			}
			p := ps[0]
			if !bytes.Equal(p.Data, []byte{236, byte(rate)}) || p.Recipient != []byte{1, 7, 31}[i] || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
				t.Fatal("rate defined fields or routing")
			}
		}
		out = append(out, results...)
	}
	gameplayReportsCapture(t, "rate", out)
}

func TestGameplayReportsStatsAndScavenger(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, op := range []int{26, 35} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, value := range []uint32{0, 1, 255, 256, 32767, 65535} {
				s := gameplayReportsBase(op, reportValue(recipient), reportObject(1), reportValue(value))
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				w := value ^ 0x5555
				p.Resources.MaxHP = uint16(value)
				a.ActorWords = map[int]uint32{36: 123, 488: value << 16}
				a.UpdateWords = map[int]uint32{8: w}
				o.PlayerDataWords[0] = map[int]uint32{2232: value << 24, 2236: value>>8 | w<<24, 2240: w >> 8, 2152: value, 2156: w}
				var want []byte
				if op == 26 {
					want = []byte{72, 123, 0, byte(value), byte(value >> 8), byte(w), byte(w >> 8), byte(value), byte(value >> 8), byte(value), byte(value >> 8), byte(w), byte(w >> 8), byte(value)}
				} else {
					a.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportObject(1)}
					want = []byte{85, 123, 0, byte(value), byte(value >> 8), byte(w), byte(w >> 8)}
				}
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		if len(ps) != 1 {
			t.Fatalf("stats case%d count%d", i, len(ps))
		}
		p := ps[0]
		recipient := []byte{1, 7, 31}[(i/6)%3]
		if checks[i][0] == 85 {
			recipient = 255
		}
		if !bytes.Equal(p.Data, checks[i]) || p.Recipient != recipient || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("stats case%d fields or routing", i)
		}
	}
	gameplayReportsCapture(t, "stats-scavenger", out)
}

func TestGameplayReportsTeamHealth(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		hp           uint16
		active, team bool
	}
	var checks []check
	for players := 0; players <= 3; players++ {
		for _, recipient := range []uint32{0, 1, 7, 31} {
			for _, hp := range []uint16{0, 1, 255, 256, 65535} {
				for _, team := range []bool{false, true} {
					s := gameplayReportsBase(18, reportValue(recipient))
					s.Lifecycle.GameFlags = 0
					if team {
						s.Lifecycle.GameFlags = 4096
					}
					p := s.Callbacks.Shop
					p.Resources.HP = hp
					p.Resources.MaxHP = 65535
					p.TemporaryUpdates.World.Objectives.Players = players
					p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: 123}
					mask := uint32(0x80000082)
					p.TemporaryUpdates.World.Objectives.Attack.Controls.Reports.RecipientMask = &mask
					cases = append(cases, s)
					checks = append(checks, check{hp, players > 0 && recipient == 1, team})
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		c := checks[i]
		n := 0
		if c.active {
			n = 1
			if c.team {
				n = 2
			}
		}
		if len(ps) != n {
			t.Fatalf("team health case%d count%d want%d", i, len(ps), n)
		}
		for j, p := range ps {
			want := []byte{67, byte(c.hp), byte(c.hp >> 8)}
			recipient := byte(1)
			if c.team && j == 0 {
				want = []byte{196, 12, 123, 0, byte(uint32(c.hp) * 100 / 65535)}
				recipient = 129
			}
			if !bytes.Equal(p.Data, want) || p.Recipient != recipient || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
				t.Fatalf("team health case%d fields or routing", i)
			}
		}
	}
	gameplayReportsCapture(t, "team-health", out)
}

func TestGameplayReportsHidden(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, recipient := range []uint32{1, 7, 31} {
		for _, class := range []uint32{0x100000, 0x800000, 0x900000} {
			for _, flags := range []uint32{0, 1, 0x1000000, 0xffffffff} {
				s := gameplayReportsBase(42, reportValue(recipient), reportObject(4))
				p := s.Callbacks.Shop
				p.Items[1].Class = class
				p.TemporaryUpdates.ItemWords[1][36] = 123
				p.TemporaryUpdates.ItemWords[1][16] = flags
				var want []byte
				if class&0x800000 != 0 {
					opcode := byte(56)
					if flags&0x1000000 != 0 {
						opcode = 55
					}
					want = []byte{opcode, 123, 0}
				}
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		want := checks[i]
		if want == nil {
			if len(ps) != 0 {
				t.Fatalf("hidden ineligible case%d", i)
			}
			continue
		}
		if len(ps) != 1 {
			t.Fatalf("hidden case%d count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, want) || p.Recipient != []byte{1, 7, 31}[i/12] || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("hidden case%d fields or routing", i)
		}
	}
	gameplayReportsCapture(t, "hidden", out)
}
