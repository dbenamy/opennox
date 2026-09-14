//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameplayReportsAggregateLocal(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks []map[byte]int
	for _, class := range []byte{0, 1, 2} {
		for mask := uint32(0); mask < 64; mask++ {
			s := gameplayReportsBase(66, reportObject(1))
			s.Lifecycle.GameFlags = 0
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			p.Resources.PlayerClass = class
			p.Resources.HP = 23
			p.Resources.MaxHP = 100
			oldHP, oldMana := uint32(23), uint32(55)
			if mask&16 != 0 {
				oldHP = 0
			}
			if mask&32 != 0 {
				oldMana = 0
			}
			armor := uint32(0x3fa00000)
			oldArmor := armor
			if mask&1 != 0 {
				oldArmor = 0
			}
			stat, oldStat := uint32(27), uint32(27)
			if mask&2 != 0 {
				oldStat = 0
			}
			oldByte := uint32(17)
			if mask&4 != 0 {
				oldByte = 0
			}
			dirty := uint32(0)
			if mask&8 != 0 {
				dirty = 1
			}
			a.ActorWords = map[int]uint32{36: 123, 440: 17}
			a.UpdateWords = map[int]uint32{4: 55 | oldMana<<16, 8: 100 | oldHP<<16, 228: armor, 232: oldArmor}
			o.PlayerDataWords[0] = map[int]uint32{2164: stat, 2168: oldStat, 2172: oldByte, 2184: dirty, 2248: uint32(class) << 24}
			p.Sequence = []legacy.PortTestShopAction{{Op: 1866}, {Op: 1866}}
			want := map[byte]int{}
			if mask&1 != 0 {
				want[73]++
			}
			if mask&2 != 0 {
				want[74]++
			}
			if mask&4 != 0 {
				want[91]++
			}
			if dirty != 0 {
				want[221]++
				want[72]++
				if class != 0 {
					want[222]++
				}
			}
			if mask&16 != 0 {
				want[67]++
			}
			if mask&32 != 0 && class != 0 {
				want[69]++
			}
			cases = append(cases, s)
			checks = append(checks, want)
		}
	}
	for _, ref := range []int{0, 4} {
		s := gameplayReportsBase(66, reportObject(ref))
		s.Callbacks.Shop.Sequence = []legacy.PortTestShopAction{{Op: 1866}, {Op: 1866}}
		cases = append(cases, s)
		checks = append(checks, map[byte]int{})
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		steps := r.Callbacks.Shop.Sequence
		got := map[byte]int{}
		for _, p := range steps[0].Packets {
			if len(p.Data) == 0 || p.Recipient != 1 {
				t.Fatalf("aggregate local case%d route", i)
			}
			got[p.Data[0]]++
		}
		if !reflect.DeepEqual(got, checks[i]) {
			t.Fatalf("aggregate local case%d opcodes %v want%v", i, got, checks[i])
		}
		if !reflect.DeepEqual(steps[0].Packets, steps[1].Packets) {
			t.Fatalf("aggregate local case%d did not settle", i)
		}
	}
	gameplayReportsCapture(t, "aggregate-local", out)
}

func TestGameplayReportsAggregateCoop(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type check struct {
		players, keys int
		level         uint32
		prior         byte
	}
	var checks []check
	for players := 1; players <= 3; players++ {
		for keys := 0; keys < 4; keys++ {
			for _, level := range []uint32{0, 1, 255, 256} {
				for _, prior := range []byte{0, 1} {
					s := gameplayReportsBase(66, reportObject(1))
					s.Lifecycle.GameFlags = 4096
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					a := o.Attack
					o.Players = players
					a.ActorWords = map[int]uint32{36: 123, 440: 0}
					a.UpdateWords = map[int]uint32{4: 50 | 50<<16, 8: 100 | 50<<16, 228: 0, 232: 0, 320: level}
					// Every recipient's cached state is a byte; values above255 intentionally keep reporting.
					for off := 452; off < 548; off += 4 {
						a.UpdateWords[off] = uint32(prior) * 0x01010101
					}
					o.PlayerDataWords[0] = map[int]uint32{2164: 0, 2168: 0, 2172: 0, 2184: 0}
					p.Inventory.ItemTypes = []string{"SilverKey", "GoldKey"}
					p.Inventory.Linked = nil
					p.Inventory.Owned = nil
					for j := 0; j < 2; j++ {
						if keys&(1<<j) != 0 {
							p.Inventory.Linked = append(p.Inventory.Linked, j)
							p.Inventory.Owned = append(p.Inventory.Owned, j)
						}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 1866}, {Op: 1866}}
					cases = append(cases, s)
					checks = append(checks, check{players, keys, level, prior})
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		steps := r.Callbacks.Shop.Sequence
		c := checks[i]
		expected := 0
		if c.level != uint32(c.prior) {
			expected += c.players
		}
		for key := 0; key < 2; key++ {
			value := byte((c.keys >> key) & 1)
			if value != c.prior {
				expected += c.players + 1
			}
		}
		if len(steps[0].Packets) != expected {
			t.Fatalf("aggregate coop case%d count%d want%d", i, len(steps[0].Packets), expected)
		}
		extra := 0
		if c.level > 255 {
			extra = c.players
		}
		if len(steps[1].Packets) != expected+extra {
			t.Fatalf("aggregate coop case%d repeated state", i)
		}
		for _, p := range steps[0].Packets {
			if len(p.Data) != 5 || p.Data[0] != 240 || p.Data[3] != 123 || p.Data[4] != 0 || p.Ordered != 0 || p.A4 != 0 || p.A5 != 1 {
				t.Fatalf("aggregate coop case%d fields", i)
			}
			value := byte(c.level)
			switch p.Data[1] {
			case 4:
			case 22:
				value = byte(c.keys & 1)
			case 23:
				value = byte((c.keys >> 1) & 1)
			default:
				t.Fatalf("aggregate coop case%d subtype", i)
			}
			if p.Data[2] != value {
				t.Fatalf("aggregate coop case%d value", i)
			}
		}
		a, b := steps[0].GameplayReports.Caches, steps[1].GameplayReports.Caches
		if a[1] == 0 || a[2] == 0 || a[1] == a[2] || a != b {
			t.Fatalf("aggregate coop case%d lookup caches", i)
		}
	}
	gameplayReportsCapture(t, "aggregate-coop", out)
}
