//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func penaltyBase(op int) legacy.PortTestRoamSpec {
	s := callbackBase(44 + op)
	s.Lifecycle.GameFlags = 4096
	p := &legacy.PortTestPenaltySpec{Gold: 101, SpellTable: [][2]uint32{{1, 1}, {4, 1}, {9, 1}, {27, 1}, {34, 1}, {41, 1}, {136, 1}}, BeastTable: [][2]uint32{{1, 1}, {7, 1}, {40, 1}}, Armor: map[uint16]uint32{16: 2, 17: 4}}
	p.Spells[1] = 1
	p.Spells[4] = 2
	p.Spells[136] = 1
	p.Beasts[1] = 1
	p.Beasts[7] = 2
	p.Beasts[40] = 1
	s.Callbacks.Penalty = p
	return s
}
func TestQuestPenaltySmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 7; op++ {
		for class := byte(0); class < 3; class++ {
			s := penaltyBase(op)
			s.Callbacks.Penalty.Class = class
			specs = append(specs, s)
		}
	}
	callbackHash(t, "penalty-smoke", legacy.PortTestRoam(specs), penaltyHashes["penalty-smoke"])
}
func TestQuestPenaltyInventorySmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{0, 1, 2, 6} {
		s := penaltyBase(op)
		s.Callbacks.Penalty.Items = []legacy.PortTestPenaltyItem{
			{Type: 23, Class: 8, Worth: 101}, {Type: 23, Class: 8, Worth: 103}, {Type: 23, Class: 8, Worth: 105},
			{Type: 24, Class: 8, Worth: 55}, {Type: 25, Class: 8, Worth: 99},
			{Type: 16, Class: 0x2000000, Flags: 0x100}, {Type: 17, Class: 0x2000000, Flags: 0x100},
			{Type: 15, Class: 0x1000000, Subclass: 4, Flags: 0x100, Mods: [4]bool{true}},
			{Type: 21, Class: 0x1000000, Eligible: true},
		}
		specs = append(specs, s)
	}
	callbackHash(t, "penalty-inventory-smoke", legacy.PortTestRoam(specs), penaltyHashes["penalty-inventory-smoke"])
}

func TestQuestPenaltyKnowledge(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op <= 5; op++ {
		if op == 1 || op == 2 {
			continue
		}
		for class := byte(0); class < 4; class++ {
			for mode := 0; mode < 6; mode++ {
				for seed := 1; seed <= 4; seed++ {
					s := penaltyBase(op)
					s.Seed = seed
					p := s.Callbacks.Penalty
					p.Class = class
					p.Spells = [137]uint32{}
					p.Beasts = [41]uint32{}
					switch mode {
					case 1:
						p.Spells[0] = 1
						p.Beasts[0] = 1
					case 2:
						p.Spells[1] = 1
						p.Beasts[1] = 1
					case 3:
						for _, id := range []int{9, 27, 34, 41} {
							p.Spells[id] = 1
						}
						p.Beasts[7] = 2
					case 4:
						for _, id := range []int{1, 4, 5, 136} {
							p.Spells[id] = 3
						}
						p.Beasts[1] = 1
						p.Beasts[40] = 1
					case 5:
						p.Spells[136] = 1
						p.Beasts[40] = 1
						p.SpellTable = [][2]uint32{{136, 0}}
						p.BeastTable = [][2]uint32{{40, 0}}
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "penalty-knowledge", legacy.PortTestRoam(specs), penaltyHashes["penalty-knowledge"])
}
func TestQuestPenaltyEquipment(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, sub := range []uint32{0, 2, 4, 0x100, 0x104, 0x10000} {
		for mods := 0; mods < 5; mods++ {
			for order := 0; order < 3; order++ {
				for _, eligible := range []bool{false, true} {
					s := penaltyBase(1)
					p := s.Callbacks.Penalty
					selected := legacy.PortTestPenaltyItem{Type: 15, Class: 0x1000000, Flags: 0x100, Subclass: sub}
					if mods > 0 {
						selected.Mods[mods-1] = true
					}
					replacement := legacy.PortTestPenaltyItem{Type: 21, Class: 0x1000000, Eligible: eligible}
					switch order {
					case 0:
						p.Items = []legacy.PortTestPenaltyItem{selected}
					case 1:
						p.Items = []legacy.PortTestPenaltyItem{replacement, selected}
					case 2:
						p.Items = []legacy.PortTestPenaltyItem{selected, replacement}
					}
					specs = append(specs, s)
				}
			}
		}
	}
	for count := 0; count <= 4; count++ {
		for _, bit := range []uint32{0, 1, 2, 4, 0x400, 0x405} {
			for _, root := range []bool{false, true} {
				for seed := 1; seed <= 4; seed++ {
					op := 2
					if root {
						op = 0
					}
					s := penaltyBase(op)
					s.Seed = seed
					p := s.Callbacks.Penalty
					p.Armor = map[uint16]uint32{16: bit}
					for i := 0; i < count; i++ {
						p.Items = append(p.Items, legacy.PortTestPenaltyItem{Type: 16, Class: 0x2000000, Flags: 0x100})
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "penalty-equipment", legacy.PortTestRoam(specs), penaltyHashes["penalty-equipment"])
}
func TestQuestPenaltyGems(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for count := 0; count <= 5; count++ {
		for _, gold := range []uint32{0, 1, 2, 3, 0xffffffff} {
			for _, worth := range []uint32{0, 1, 2, 3, 101, 0x7fffffff, 0xffffffff} {
				for _, root := range []bool{false, true} {
					op := 6
					if root {
						op = 0
					}
					s := penaltyBase(op)
					p := s.Callbacks.Penalty
					p.Gold = gold
					for i := 0; i < count; i++ {
						for _, typ := range []uint16{23, 24, 25} {
							p.Items = append(p.Items, legacy.PortTestPenaltyItem{Type: typ, Class: 8, Worth: worth + uint32(i)})
						}
					}
					if count%2 != 0 {
						p.Cache = [6]uint32{23, 24, 25, 23, 25, 24}
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "penalty-gems", legacy.PortTestRoam(specs), penaltyHashes["penalty-gems"])
}
func TestQuestPenaltyCorpus(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	state := uint32(0x54cbd0)
	next := func() uint32 { state ^= state << 13; state ^= state >> 17; state ^= state << 5; return state }
	for n := 0; n < 1024; n++ {
		s := penaltyBase(int(next() % 7))
		s.Seed = int(next() % 1000)
		p := s.Callbacks.Penalty
		p.Class = byte(next() % 4)
		p.Gold = next()
		s.Lifecycle.GameFlags = []uint32{0, 2048, 4096}[next()%3]
		p.Spells = [137]uint32{}
		p.Beasts = [41]uint32{}
		for _, id := range []int{0, 1, 2, 3, 4, 5, 9, 27, 34, 41, 136} {
			p.Spells[id] = next() % 4
		}
		for _, id := range []int{0, 1, 7, 40} {
			p.Beasts[id] = next() % 4
		}
		p.Items = nil
		count := int(next() % 12)
		for i := 0; i < count; i++ {
			item := legacy.PortTestPenaltyItem{Type: uint16(23 + next()%3), Class: 8, Worth: next() % 100000, Eligible: next()%2 == 0}
			switch next() % 3 {
			case 0:
				item.Type = 15
				item.Class = 0x1000000
				item.Subclass = []uint32{0, 2, 4, 0x104, 0x10000}[next()%5]
				item.Flags = []uint32{0, 0x100, 0x120}[next()%3]
				item.Mods[next()%4] = next()%2 == 0
			case 1:
				item.Type = uint16(16 + next()%2)
				item.Class = 0x2000000
				item.Flags = []uint32{0, 0x100, 0x120}[next()%3]
			}
			p.Items = append(p.Items, item)
		}
		if next()%2 == 0 {
			p.Cache = [6]uint32{23, 24, 25, 23, 25, 24}
		}
		specs = append(specs, s)
	}
	callbackHash(t, "penalty-corpus", legacy.PortTestRoam(specs), penaltyHashes["penalty-corpus"])
}

func TestQuestPenaltyGoldProtection(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, gold := range []uint32{0, 1, 2, 3, 0xffffffff} {
		for _, count := range []int{0, 1, 2, 3} {
			s := penaltyBase(0)
			p := s.Callbacks.Penalty
			p.Gold = gold
			p.ProtectedGold = true
			for i := 0; i < count; i++ {
				for _, typ := range []uint16{23, 24, 25} {
					p.Items = append(p.Items, legacy.PortTestPenaltyItem{Type: typ, Class: 8, Worth: uint32(100 + i)})
				}
			}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "penalty-protection", legacy.PortTestRoam(specs), penaltyHashes["penalty-protection"])
}

var penaltyHashes = map[string]string{
	"penalty-smoke":           "935d5c69aae17f96ba002918c3ab251c162e3befb5b8993ea7b3b5f0e6cd97cb",
	"penalty-inventory-smoke": "284b6f32f7891acca86436e97b60e082df8832fa5fb775ee484d88afbba5149b",
	"penalty-knowledge":       "43b81f95b13274ec0a001f340bea3fce8252fb4255214f1d0621a833653ed91f",
	"penalty-equipment":       "7599536683bf26456102e04348984b21a5ab53839f7f1091cd899f57b72dc678",
	"penalty-gems":            "9678247c376349f62e21c332ecde30986f6b3d1859104d69e373be78988e2c73",
	"penalty-corpus":          "9cb2754bde579826b7486dedb7b222c6ce60bb172a8159db5d13818275942c15",
	"penalty-protection":      "682c54b2ea10f9335abf9d55cde4669ed874cb5cac1140ebf66774773e9cfd15",
}
