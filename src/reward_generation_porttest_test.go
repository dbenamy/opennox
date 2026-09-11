//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

var rewardHashes = map[string]string{
	"reward-category-masks":         "8b0513082443c8805aada0de7d7e9c6297c57ae8ddd9354b4dc51747bb2a14f0",
	"reward-empty-placement-19":     "03c34cc431a939f05e14e653b1c40a19a5f4ad1da8a9c77176da6cc34ae4efee",
	"reward-empty-placement-20":     "f74d2ee8bdff44826f36a62384ed426e2d8aa774c2530df6550939df16343174",
	"reward-explicit-10":            "2789c03af207e9c6803dec8c6c61cb4d804ba8e86d6de9ae1ddb32efd897cc78",
	"reward-explicit-12":            "06c5a6585766c97e1816dabc66186e3474ee23482053b961ad2532e07fed7891",
	"reward-explicit-13":            "10ffa1ca664c6d6a2e8a8272676ca25a731033a9cc7b6b2a5069fb050db255d2",
	"reward-factory-10":             "594ce9816283916e917cb9fb88ce8dab717742af8e83fce442b2086254dc3015",
	"reward-factory-12":             "e56e96a3470b84770efba753d99029dddce0a4642f2a8dafc5604817f52dc6b1",
	"reward-factory-13":             "330a99e86a029a885c205e3bc150fee085dc45fc0bb038687ef5b62428b6edf1",
	"reward-factory-14":             "b5dae720fe9a5e594efee7eb3a1294fe456a7f9014446cb89eae6bce85d94c6c",
	"reward-factory-15":             "a260c38f862cf7df89a7497322408424d002f9e5617da2aa4c5df1ab15367960",
	"reward-factory-16":             "d18facaf548e38639b3f33ded662a310e28a774b4b54cce7d67e6efb354e4ebd",
	"reward-factory-17":             "73e9e3183bf03e39fe8d53c3379d4f692239e08225d6b527e27b8cccb87bb081",
	"reward-factory-18":             "73e9e3183bf03e39fe8d53c3379d4f692239e08225d6b527e27b8cccb87bb081",
	"reward-filter-10":              "d359e19c7f9ba6556009f49d8983be9e73d875f11e7fdd54ce6d9315551324e2",
	"reward-filter-13":              "8c43fe4661d9217971196b23df634a14d18a9754d265384a509990590d92e4fa",
	"reward-filter-14":              "a24a3af8dda9dee40fd1eb9ac40c8b9bc526b50c66847ebe49fa116bcc1a9327",
	"reward-filter-15":              "f23c43b3120f75a19b310c6279eff9daf373335cfa982899fa982c63b9fda989",
	"reward-filter-16":              "de88aaf5b5891d90b29627da4cace6ef8a99405b831467c2a2270f3eae3cd6da",
	"reward-gold-players":           "faf94281941b7582fe700e82ebdf79b5da63d82278b2c6da4a9a35279befb8d7",
	"reward-init-00":                "96a3227b7d6bc956a3f49473813ecbe7069f67e2c8f147e2f0e46f2a17dd1ee2",
	"reward-init-01":                "15318534d0d16867097b21cd66c13f02122b51182e2e78ddad42039c15317c21",
	"reward-init-02":                "c639c9a477ccb0f77b523f4950b8db04b29c8f06ec5fd271fb96c447a0abbd44",
	"reward-init-03":                "644d0c9297b97e2648f20ad786865519b9f4593f0b535d236afff928cd5715ac",
	"reward-init-04":                "f30e754618322083ae7f33620a86bd453b7ad4c6a1b4b50231f7f992ea6858ee",
	"reward-init-05":                "f30e754618322083ae7f33620a86bd453b7ad4c6a1b4b50231f7f992ea6858ee",
	"reward-init-06":                "65857efac63cc3892459fc5824c221cefb4e5bd27e6f9503f99e14d0f016a78f",
	"reward-init-07":                "c639c9a477ccb0f77b523f4950b8db04b29c8f06ec5fd271fb96c447a0abbd44",
	"reward-init-08":                "63f06a763a529dd0656cebf60085b8ff1aca01875ac747e5b97fe717093b9e2f",
	"reward-marker":                 "c7b46d134f130601194dc02b12b68a1b427f88b430e49f864054f99a66b1a960",
	"reward-missing-10":             "408e82d123286309d57c92f1557ead9265a9a9893242115b7397afddb780dba3",
	"reward-missing-12":             "408e82d123286309d57c92f1557ead9265a9a9893242115b7397afddb780dba3",
	"reward-missing-13":             "408e82d123286309d57c92f1557ead9265a9a9893242115b7397afddb780dba3",
	"reward-missing-14":             "bff9f28964b239e4ff73d2aeebc6f1dc8818c7669adb87964f8d256196ad8f4a",
	"reward-missing-15":             "bff9f28964b239e4ff73d2aeebc6f1dc8818c7669adb87964f8d256196ad8f4a",
	"reward-missing-16":             "bff9f28964b239e4ff73d2aeebc6f1dc8818c7669adb87964f8d256196ad8f4a",
	"reward-missing-17":             "6b1d879cea97502d19e2f4e6e7f9d06c63be01da782f473b8cf36e90462e4293",
	"reward-modifier-boundary-14-6": "24563bfb5845a9f4bf4b0da69aecc60df2c8e635c63cddaa84a487ed55e1260f",
	"reward-modifier-boundary-14-7": "0296478de4551af289deb2cbbbb5c5c90a98305e93358d6d54b35336995b6b89",
	"reward-modifier-boundary-15-6": "5355c06f5b8acda226da0f9a957cd6f1e60245a7cb1b6ea7325304a716da9465",
	"reward-modifier-boundary-15-7": "a36d7a264d02cd4c79c730773b9b26d39f134c07449a102de687da8bf33ec0d4",
	"reward-named-initializer":      "3fa0460493b8e75bb2ee949c7a5db96946907b9cf6d0bfd843111321cd4af6ae",
	"reward-placement-19":           "42bb3e22949b030d709824a8c8861ee2a9821afef33ff410224657c57a98b0a1",
	"reward-placement-20":           "82897ceec6fe3140f5a7c981d951ea21784d674948aefa449c5aa696d654fe98",
	"reward-positive":               "b201f775007478cb389b34314da6617a9f91fa2818910fc1845978db858318a7",
	"reward-tiers":                  "d0a6b846542e0618889bd5b9595ee5fa521a7e857a98697e6a33c81303163f52",
	"reward-weapon-mode-4":          "5a565a9594dd2513c3562f5f05ac165ed6db64545757721795857f4f001fea4a",
	"reward-weapon-mode-5":          "47e3198eb8e5aa41b0a2b5c4783575c9abd0840ced1fc4ad06809c3fdeffa19d",
}

func rewardBase(op int) legacy.PortTestRoamSpec {
	s := attackBase()
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	a.Actor = 3
	a.Reward = &legacy.PortTestRewardSpec{Allowed: true, ArmorBit: 1, WeaponBit: 1, InitWords: map[int]uint32{}, UpdateWords: map[int]uint32{}, ActorType: "RewardMarker"}
	a.ActorWords = map[int]uint32{12: 1}
	p.Inventory.Linked = nil
	p.Inventory.Owned = nil
	p.Equipment.ActiveWeapon = 0
	p.Equipment.WeaponFlags = 0
	p.Sequence = []legacy.PortTestShopAction{{Op: 1300 + op}}
	if p.Balance == nil {
		p.Balance = make(map[string]float64)
	}
	for key, v := range map[string]float64{"GeneratorMaxActiveCreaturesHigh": 9, "GeneratorMaxActiveCreaturesNormal": 5, "GeneratorMaxActiveCreaturesLow": 2, "GeneratorMaxActiveCreaturesSingular": 1} {
		p.Balance[key] = v
	}
	return s
}
func rewardHash(t *testing.T, name string, cases []legacy.PortTestRoamSpec) {
	t.Helper()
	callbackHash(t, name, effectsTimedRun(t, cases), rewardHashes[name])
}
func TestRewardInitialization(t *testing.T) {
	for op := 0; op < 9; op++ {
		var cases []legacy.PortTestRoamSpec
		for i := 0; i < 32; i++ {
			s := rewardBase(op)
			s.Seed = i + 1
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			r := a.Reward
			r.InitWords[0] = uint32(i % 9)
			if op == 6 {
				r.InitWords[0] = 0
				if i%4 == 0 {
					r.InitWords[0] = 77
				}
			}
			a.ActorWords[20] = uint32(i % 16)
			a.ActorWords[12] = uint32(i % 16)
			r.UpdateWords[80] = uint32(i%5) << 24
			r.UpdateWords[84] = 0x03020100
			r.GeneratorStage = uint32(i % 4)
			r.UpdateWords[16] = 0 // empty type name for direction/name initializer.
			cases = append(cases, s)
		}
		rewardHash(t, fmt.Sprintf("reward-init-%02d", op), cases)
	}
}
func TestRewardTiers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, stage := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0x7fffffff, 0xffffffff} {
		for seed := uint32(1); seed <= 64; seed++ {
			s := rewardBase(11)
			s.Seed = int(seed)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward.Stage = stage
			cases = append(cases, s)
		}
	}
	rewardHash(t, "reward-tiers", cases)
}
func TestRewardFactories(t *testing.T) {
	for _, op := range []int{10, 12, 13, 14, 15, 16, 17, 18} {
		var cases []legacy.PortTestRoamSpec
		for stage := uint32(0); stage <= 11; stage++ {
			for seed := uint32(1); seed <= 24; seed++ {
				s := rewardBase(op)
				s.Seed = int(seed)
				r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
				r.Stage = stage
				cases = append(cases, s)
			}
		}
		rewardHash(t, fmt.Sprintf("reward-factory-%02d", op), cases)
	}
}
func TestRewardFilters(t *testing.T) {
	for _, op := range []int{10, 13, 14, 15, 16} {
		var cases []legacy.PortTestRoamSpec
		for mode := 1; mode <= 3; mode++ {
			for seed := uint32(1); seed <= 16; seed++ {
				s := rewardBase(op)
				s.Seed = int(seed)
				r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
				r.Stage = 9
				r.TableMode = mode
				cases = append(cases, s)
			}
		}
		rewardHash(t, fmt.Sprintf("reward-filter-%02d", op), cases)
	}
}
func TestRewardMarker(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for category := uint32(1); category <= 128; category <<= 1 {
		for chance := uint32(0); chance <= 4; chance++ {
			for seed := uint32(1); seed <= 12; seed++ {
				s := rewardBase(9)
				s.Seed = int(seed)
				r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
				r.Stage = 8
				r.InitWords[0] = category
				r.InitWords[212] = chance
				if seed%2 == 0 {
					r.ActorType = "RewardMarkerPlus"
				}
				cases = append(cases, s)
			}
		}
	}
	rewardHash(t, "reward-marker", cases)
}
func TestRewardGoldPlayers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for n := 1; n <= 3; n++ {
		for _, xp := range []float32{0, 1, 1.5, 100, 99999} {
			s := rewardBase(6)
			o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
			o.Players = n
			o.PlayerWords = []map[int]uint32{{28: math.Float32bits(xp)}, {28: math.Float32bits(xp + 1)}, {28: math.Float32bits(xp + 2)}}
			cases = append(cases, s)
		}
	}
	rewardHash(t, "reward-gold-players", cases)
}
func rewardByte(words map[int]uint32, off int, v byte) {
	base := off &^ 3
	shift := uint(off&3) * 8
	words[base] = (words[base] &^ (255 << shift)) | uint32(v)<<shift
}
func TestRewardExplicitBooks(t *testing.T) {
	for _, op := range []int{10, 12, 13} {
		var cases []legacy.PortTestRoamSpec
		count, start, flag := 137, 8, uint32(1)
		if op == 12 {
			count, start, flag = 6, 145, 2
		}
		if op == 13 {
			count, start, flag = 41, 151, 4
		}
		for choice := 0; choice < count; choice++ {
			for _, value := range []byte{0, 1, 2} {
				s := rewardBase(op)
				r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
				r.InitWords[4] = flag
				rewardByte(r.InitWords, start+choice, value)
				cases = append(cases, s)
			}
		}
		for seed := 1; seed <= 32; seed++ {
			s := rewardBase(op)
			s.Seed = seed
			r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
			r.InitWords[4] = flag
			for i := 0; i < count; i++ {
				rewardByte(r.InitWords, start+i, 1)
			}
			cases = append(cases, s)
		}
		rewardHash(t, fmt.Sprintf("reward-explicit-%02d", op), cases)
	}
}
func TestRewardMissingFactories(t *testing.T) {
	for _, op := range []int{10, 12, 13, 14, 15, 16, 17} {
		var cases []legacy.PortTestRoamSpec
		for seed := 1; seed <= 24; seed++ {
			s := rewardBase(op)
			s.Seed = seed
			r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
			r.Stage = 10
			r.Missing = []string{"ConjurerSpellBook", "WizardSpellBook", "CommonSpellBook", "AbilityBook", "FieldGuide", "QuestGoldChest", "QuestGoldPile", "RubyGem", "EmeraldGem", "DiamondGem"}
			if op == 14 || op == 15 || op == 16 {
				r.Allowed = false
			}
			cases = append(cases, s)
		}
		rewardHash(t, fmt.Sprintf("reward-missing-%02d", op), cases)
	}
}
func TestRewardPlacement(t *testing.T) {
	for _, op := range []int{19, 20} {
		var cases []legacy.PortTestRoamSpec
		for mode := 0; mode < 6; mode++ {
			for seed := 1; seed <= 24; seed++ {
				s := rewardBase(op)
				s.Seed = seed
				o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
				r := o.Attack.Reward
				r.Stage = uint32(seed % 3)
				o.Players = seed%3 + 1
				o.PlayerDataWords = []map[int]uint32{{4792: 1}, {4792: 1}, {4792: 1}}
				o.ObjectList = []int{3, 4, 5}
				r.Types = map[int]string{3: "RewardMarker", 4: "RewardMarker", 5: "RewardMarker"}
				if mode%3 == 1 {
					r.Types[4] = "RewardMarkerPlus"
				}
				if mode%3 == 2 {
					r.Types[5] = "RedPotion"
				}
				r.InitWords[0] = 0x80
				r.InitByRef = map[int]map[int]uint32{4: {0: 0x80}, 5: {0: 0x80}}
				if mode >= 3 {
					r.InitWords[216] = 1
				}
				cases = append(cases, s)
			}
		}
		rewardHash(t, fmt.Sprintf("reward-placement-%02d", op), cases)
	}
}

func TestRewardEquipmentSlots(t *testing.T) {
	for _, mode := range []int{4, 5} {
		var cases []legacy.PortTestRoamSpec
		for seed := 1; seed <= 256; seed++ {
			s := rewardBase(15)
			s.Seed = seed
			r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
			r.Stage = 10
			r.TableMode = mode
			cases = append(cases, s)
		}
		rewardHash(t, fmt.Sprintf("reward-weapon-mode-%d", mode), cases)
	}
}
func TestRewardNamedInitializer(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for dir := 0; dir < 9; dir++ {
		for _, name := range []string{"RewardMarker", "RewardMarkerPlus", "missing"} {
			s := rewardBase(4)
			r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
			r.InitWords[0] = uint32(dir)
			r.UpdateName = name
			cases = append(cases, s)
		}
	}
	rewardHash(t, "reward-named-initializer", cases)
}

func TestRewardEmptyPlacement(t *testing.T) {
	for _, op := range []int{19, 20} {
		s := rewardBase(op)
		o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
		o.ObjectList = nil
		rewardHash(t, fmt.Sprintf("reward-empty-placement-%d", op), []legacy.PortTestRoamSpec{s})
	}
}

func TestRewardPositive(t *testing.T) {
	weapon := rewardBase(15)
	weapon.Seed = 3
	weapon.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward.Stage = 10
	weapon.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward.TableMode = 5
	book := rewardBase(12)
	br := book.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
	br.InitWords[4] = 2
	rewardByte(br.InitWords, 146, 1)
	cases := []legacy.PortTestRoamSpec{weapon, book}
	results := effectsTimedRun(t, cases)
	for i, r := range results {
		if len(r.Lifecycle.Created) != 1 || r.Callbacks.Shop.Sequence[0].Return != 1001 {
			t.Fatalf("case %d did not return a new captured object", i)
		}
	}
	mods := results[0].Callbacks.Modifiers
	if len(mods) != 6 {
		t.Fatalf("weapon modifier capture %v", mods)
	}
	for slot, v := range mods[1:5] {
		if v == 0 {
			t.Fatalf("weapon modifier slot %d missing", slot)
		}
	}
	bookData := results[1].Callbacks.Modifiers
	if len(bookData) != 6 || bookData[5] != 1 {
		t.Fatalf("explicit ability book contents: %v", bookData)
	}
	callbackHash(t, "reward-positive", results, rewardHashes["reward-positive"])
}

func TestRewardModifierBoundaries(t *testing.T) {
	for _, op := range []int{14, 15} {
		for _, mode := range []int{6, 7} {
			var cases []legacy.PortTestRoamSpec
			for stage := uint32(0); stage <= 11; stage++ {
				for seed := 1; seed <= 16; seed++ {
					s := rewardBase(op)
					s.Seed = seed
					r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
					r.Stage = stage
					r.TableMode = mode
					cases = append(cases, s)
				}
			}
			rewardHash(t, fmt.Sprintf("reward-modifier-boundary-%d-%d", op, mode), cases)
		}
	}
}
func TestRewardCategoryMasks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mask := range []uint32{0, 3, 0x18, 0xff, 0xffff, 0xffffffff} {
		for _, chance := range []uint32{0, 5, 0xffffffff} {
			for seed := 1; seed <= 16; seed++ {
				s := rewardBase(9)
				s.Seed = seed
				r := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Reward
				r.Stage = 8
				r.InitWords[0] = mask
				r.InitWords[212] = chance
				cases = append(cases, s)
			}
		}
	}
	rewardHash(t, "reward-category-masks", cases)
}
