//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"strings"
	"testing"
)

func TestGameplayReportsJournalBoundaries(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, op := range []int{53, 54, 55} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, length := range []int{0, 1, 62, 63} {
				for _, flags := range []uint16{0, 1, 255, 256, 32767, 65535} {
					name := strings.Repeat("j", length)
					s := gameplayReportsBase(op, reportValue(recipient), legacy.PortTestGameplayReportArg{Kind: "record", Ref: 1})
					s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.Reports.Records = []legacy.PortTestGameplayReportRecord{{Text: name, Flags: flags}}
					want := make([]byte, 68)
					want[0] = 213
					want[1] = byte(op - 52)
					copy(want[2:], name)
					if op != 54 {
						binary.LittleEndian.PutUint16(want[66:], flags)
					}
					cases = append(cases, s)
					checks = append(checks, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		if len(ps) != 1 {
			t.Fatalf("journal boundary case%d count%d", i, len(ps))
		}
		p := ps[0]
		if !bytes.Equal(p.Data, checks[i]) || p.Recipient != []byte{1, 7, 31}[(i/24)%3] || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("journal boundary case%d bytes/routing", i)
		}
	}
	gameplayReportsCapture(t, "journal-boundaries", out)
}

func TestGameplayReportsNonPlayers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{1, 6, 21, 23, 24, 25, 26, 31, 35, 66} {
		s := gameplayReportsBase(op, reportValue(7), reportObject(4), reportValue(255))
		switch op {
		case 1, 6, 31, 35, 66:
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.Reports.Args = [5]legacy.PortTestGameplayReportArg{reportObject(4), reportValue(255)}
		}
		cases = append(cases, s)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		if len(r.Callbacks.Shop.Sequence[0].Packets) != 0 {
			t.Fatalf("non-player report case%d", i)
		}
	}
	gameplayReportsCapture(t, "non-players", out)
}

func TestGameplayReportsAggregateFloatCache(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		bits          uint32
		first, second int
	}
	var checks []check
	values := []uint32{0, 0x80000000, 0x3f800000, 0x7f800000, 0xff800000, 0x7fc12345, 0x7f812345}
	for _, v := range values {
		for _, old := range values {
			s := gameplayReportsBase(66, reportObject(1))
			s.Lifecycle.GameFlags = 0
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			a.ActorWords = map[int]uint32{440: 0}
			a.UpdateWords = map[int]uint32{4: 50 | 50<<16, 8: 100 | 50<<16, 228: v, 232: old}
			o.PlayerDataWords[0] = map[int]uint32{2164: 0, 2168: 0, 2172: 0, 2184: 0}
			p.Sequence = []legacy.PortTestShopAction{{Op: 1866}, {Op: 1866}}
			first := 0
			if math.Float32frombits(v) != math.Float32frombits(old) {
				first = 1
			}
			second := first
			if math.IsNaN(float64(math.Float32frombits(v))) {
				second++
			}
			cases = append(cases, s)
			checks = append(checks, check{v, first, second})
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		c := checks[i]
		for j, step := range r.Callbacks.Shop.Sequence {
			n := c.first
			if j == 1 {
				n = c.second
			}
			if len(step.Packets) != n {
				t.Fatalf("float cache case%d step%d count%d want%d", i, j, len(step.Packets), n)
			}
			for _, p := range step.Packets {
				if len(p.Data) != 5 || p.Data[0] != 73 || binary.LittleEndian.Uint32(p.Data[1:]) != c.bits || p.Recipient != 1 || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
					t.Fatalf("float cache case%d payload", i)
				}
			}
		}
	}
	gameplayReportsCapture(t, "aggregate-float-cache", out)
}

func TestGameplayReportsTeamBaseCache(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var eligible []bool
	for _, warm := range []bool{false, true} {
		for _, matching := range []bool{false, true} {
			for _, equipment := range []bool{false, true} {
				s := gameplayReportsBase(50, reportValue(7), reportObject(3))
				p := s.Callbacks.Shop
				sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.Reports
				p.Items[0].Class = 0x100000
				if equipment {
					p.Items[0].Class = 0x2000000
				}
				p.Items[0].Mods = [4]bool{}
				if warm {
					sp.Caches[0] = 23
					p.Items[0].Type = 24
					if matching {
						p.Items[0].Type = 23
					}
				} else {
					name := "SilverKey"
					if matching {
						name = "TeamBase"
					}
					p.Inventory.ItemTypes = []string{name}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: 1850}, {Op: 1850}}
				cases = append(cases, s)
				eligible = append(eligible, matching || equipment)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		steps := r.Callbacks.Shop.Sequence
		cache := steps[0].GameplayReports.Caches[0]
		if cache == 0 || cache != steps[1].GameplayReports.Caches[0] {
			t.Fatalf("team base case%d cache", i)
		}
		if i >= 4 && cache != 23 {
			t.Fatalf("team base case%d warm lookup", i)
		}
		for j, step := range steps {
			n := 0
			if eligible[i] {
				n = j + 1
			}
			if len(step.Packets) != n {
				t.Fatalf("team base case%d step%d count", i, j)
			}
			for _, p := range step.Packets {
				if len(p.Data) != 7 || p.Data[0] != 103 || !bytes.Equal(p.Data[3:], []byte{255, 255, 255, 255}) || p.Recipient != 7 || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
					t.Fatalf("team base case%d fields", i)
				}
			}
		}
	}
	gameplayReportsCapture(t, "team-base-cache", out)
}

func TestGameplayReportsZeroMaxHealth(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, op := range []int{15, 16, 20} {
		for _, recipient := range []uint32{1, 7, 31} {
			for _, hp := range []uint16{0, 1, 65535} {
				s := gameplayReportsBase(op, reportValue(recipient), reportObject(1))
				p := s.Callbacks.Shop
				p.Resources.NoHealth = false
				p.Resources.HP = hp
				p.Resources.MaxHP = 0
				p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: 123}
				var want []byte
				switch op {
				case 15:
					want = []byte{221, 123, 0, byte(hp), byte(hp >> 8), 0, 0}
				case 16:
					want = []byte{65, 123, 0, byte(hp >> 1)}
				}
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		step := r.Callbacks.Shop.Sequence[0]
		want := checks[i]
		if want == nil {
			if len(step.Packets) != 0 || step.Return == 0 {
				t.Fatalf("zero max item case%d", i)
			}
			continue
		}
		if len(step.Packets) != 1 {
			t.Fatalf("zero max health case%d count", i)
		}
		p := step.Packets[0]
		if !bytes.Equal(p.Data, want) || p.Recipient != []byte{1, 7, 31}[(i/3)%3] || p.Ordered != 1 || p.A4 != 0 || p.A5 != 1 {
			t.Fatalf("zero max health case%d fields", i)
		}
	}
	gameplayReportsCapture(t, "zero-max-health", out)
}

func TestGameplayReportsEliminationFrameWrap(t *testing.T) {
	inputs := []struct {
		frame, start uint32
		elapsed      bool
	}{{0, 1, true}, {0, 0xfffffe00, false}, {0, 0xfffffda8, false}, {0, 0xfffffda7, true}}
	var cases []legacy.PortTestRoamSpec
	for _, v := range inputs {
		s := gameplayReportsBase(36, reportObject(1))
		s.Owner.Frame = v.frame
		s.Owner.FPS = 30
		s.Lifecycle.GameFlags = 1024
		p := s.Callbacks.Shop
		p.EffectsUse.Balance["SuddenDeathCountdown"] = []float64{7}
		p.TemporaryUpdates.World.Objectives.Attack.Controls.Reports.Rules = &legacy.PortTestGameplayReportRules{StartFrame: v.start, Limit: 8}
		cases = append(cases, s)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		rules := r.Callbacks.Shop.Sequence[0].GameplayReports.Rules
		n := 0
		if inputs[i].elapsed {
			n = 1
		}
		if rules.Notified != uint32(n) || len(rules.Countdown) != n {
			t.Fatalf("elimination frame wrap case%d", i)
		}
	}
	gameplayReportsCapture(t, "elimination-frame-wrap", out)
}

func TestGameplayReportsInventoryClassMasks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][]byte
	for _, op := range []int{11, 12, 28, 50} {
		for _, class := range []uint32{0x1000000, 0x10000000, 0x13001000} {
			for _, modified := range []bool{false, true} {
				s := gameplayReportsBase(op, reportValue(7), reportObject(3))
				p := s.Callbacks.Shop
				p.Equipment.HolderOnly = true
				p.Inventory.WeaponBits[23] = 0xf0
				p.Inventory.ArmorBits = map[uint16]uint32{23: 0x800}
				p.Items[0].Type = 23
				p.Items[0].Class = class
				p.Items[0].Health = false
				p.Items[0].Mods = [4]bool{modified, modified, modified, modified}
				p.TemporaryUpdates.ItemWords[0][36] = 321
				p.TemporaryUpdates.World.Objectives.Attack.ActorWords = map[int]uint32{36: 123}
				mods := []byte{255, 255, 255, 255}
				if modified {
					mods = []byte{40, 41, 42, 43}
				}
				var want []byte
				switch op {
				case 11:
					want = []byte{80, 123, 128, 240, 0, 0, 0}
					if modified {
						want[0] = 81
						want = append(want, mods...)
					}
				case 12:
					want = []byte{84, 123, 0, 240, 0, 0, 0}
				case 28:
					want = append([]byte{76, 65, 1, 23, 0}, mods...)
				case 50:
					want = append([]byte{103, 65, 1}, mods...)
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
			t.Fatalf("inventory class case%d count", i)
		}
		p := ps[0]
		ordered, priority := byte(1), uint32(0)
		if checks[i][0] == 103 {
			ordered = 0
			priority = 1
		}
		if !bytes.Equal(p.Data, checks[i]) || p.Recipient != 7 || p.Ordered != ordered || p.A4 != 0 || p.A5 != priority {
			t.Fatalf("inventory class case%d mask precedence", i)
		}
	}
	gameplayReportsCapture(t, "inventory-class-masks", out)
}
