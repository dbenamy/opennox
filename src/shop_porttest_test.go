//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func shopBase(op int) legacy.PortTestRoamSpec {
	s := callbackBase(69)
	s.Callbacks.Shop = &legacy.PortTestShopSpec{
		Op: op, Mode: 1, Session: 2,
		Item: legacy.PortTestShopItem{Type: 15, Class: 8, Worth: 101},
		Buy:  math.Float32bits(1.25), Sell: math.Float32bits(.5), GuideWorth: 137,
		SpellPrices: map[int]int{1: 201, 5: 65537, 9: -1},
		Balance:     map[string]float64{"QuestGuideWorthMultiplier": 1.5, "QuestModifierWorthMultiplier": 1.75, "DefaultAmmoAmount": 20, "DefaultAmmoAmountQuest": 30, "QuestSellMultiplier": .75, "RepairCoefficient": .3, "QuestSpellWorthMultiplier": 1.25},
	}
	return s
}

func TestShopPriceBase(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, worth := range []uint32{0, 1, 2, 3, 99, 100, 101, 16777215, 16777217, 0x7fffffff, 0x80000000, 0xffffffff} {
		for session := 0; session < 4; session++ {
			for _, mode := range []int{0, 1, 2, -1} {
				for _, quest := range []uint32{0, 4096} {
					s := shopBase(0)
					p := s.Callbacks.Shop
					p.Item.Worth, p.Session, p.Mode = worth, session, mode
					s.Lifecycle.GameFlags = quest
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		sp := specs[i].Callbacks.Shop
		if sp.Session == 0 && sp.Mode != 2 && sp.Item.Worth <= 101 {
			want := sp.Item.Worth
			if want == 0 {
				want = 1
			}
			if v.Callbacks.Return != want {
				t.Fatalf("case %d base cost=%d want%d", i, v.Callbacks.Return, want)
			}
		}
	}
	callbackHash(t, "shop-price-base", r, shopHashes["shop-price-base"])
}

func TestShopPriceModifiersHealth(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{8, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
		for mask := 0; mask < 16; mask++ {
			for mode := 0; mode < 3; mode++ {
				for variant := 0; variant < 4; variant++ {
					s := shopBase(0)
					p := s.Callbacks.Shop
					p.Mode, p.Session, p.Item.Class = mode, 2+variant%2, class
					p.Item.ModPrice = [4]int32{17, -13, 16777217, -16777215}
					for j := range p.Item.Mods {
						p.Item.Mods[j] = mask&(1<<j) != 0
					}
					p.Item.Health = true
					p.Item.HP = []uint16{0, 1, 99, 65535}[variant]
					p.Item.MaxHP = []uint16{100, 3, 100, 0}[variant]
					p.Item.QuestFlag = byte(variant % 2)
					if variant >= 2 {
						s.Lifecycle.GameFlags = 4096
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "shop-price-modifiers-health", legacy.PortTestRoam(specs), shopHashes["shop-price-modifiers-health"])
}

func TestShopPriceCharges(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, wand := range []bool{false, true} {
		for _, charges := range []byte{0, 1, 19, 20, 21, 29, 30, 31, 255} {
			for _, maximum := range []byte{0, 1, 20, 255} {
				for variant := 0; variant < 8; variant++ {
					s := shopBase(0)
					p := s.Callbacks.Shop
					p.Mode = variant % 3
					p.Item.Class, p.Item.Subclass = 0x1000000, 0x82
					p.Item.Use[0], p.Item.Use[1], p.Item.Use[2] = maximum, charges, byte(variant&1)
					if wand {
						p.Item.Class, p.Item.Subclass = 0x1000, 0x47f0000
						p.Item.Use[108], p.Item.Use[109] = charges, maximum
					}
					if variant&2 != 0 {
						s.Lifecycle.GameFlags = 4096
					}
					if variant&4 != 0 {
						p.Balance["DefaultAmmoAmount"], p.Balance["DefaultAmmoAmountQuest"] = 0, -1
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "shop-price-charges", legacy.PortTestRoam(specs), shopHashes["shop-price-charges"])
}

func TestShopPriceRewardsGems(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subclass := range []uint32{0, 1, 2, 3, 4} {
		for variant := 0; variant < 8; variant++ {
			for mode := 0; mode < 3; mode++ {
				s := shopBase(0)
				p := s.Callbacks.Shop
				p.Mode, p.Item.Class, p.Item.Subclass = mode, 0x100, subclass
				if subclass&1 != 0 {
					p.Item.Use[0] = []byte{0, 1, 5, 9}[variant%4]
				} else {
					copy(p.Item.Use[:], []string{"Zombie", "zombie", "missing", "VileZombie"}[variant%4])
				}
				if variant >= 4 {
					s.Lifecycle.GameFlags = 4096
				}
				if variant%4 == 3 {
					p.GuideWorth = -1
				}
				specs = append(specs, s)
			}
		}
	}
	for _, typ := range []uint16{15, 23, 24, 25, 0xffff} {
		for _, cold := range []bool{false, true} {
			for _, quest := range []uint32{0, 4096} {
				s := shopBase(0)
				p := s.Callbacks.Shop
				p.Mode, p.Item.Type = 0, typ
				if !cold {
					p.Cache = [3]uint32{23, 25, 24}
				}
				s.Lifecycle.GameFlags = quest
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-price-rewards-gems", legacy.PortTestRoam(specs), shopHashes["shop-price-rewards-gems"])
}

func TestShopStockMatching(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{8, 0x1000000, 0x100} {
		for _, subclass := range []uint32{0, 1, 2, 4} {
			for variant := 0; variant < 12; variant++ {
				for op := 1; op <= 2; op++ {
					s := shopBase(op)
					p := s.Callbacks.Shop
					p.Item.Class, p.Item.Subclass = class, subclass
					p.Item.Mods = [4]bool{true, false, true, false}
					entry := legacy.PortTestShopStock{Type: 15, Reward: 2, Mods: p.Item.Mods}
					if subclass&5 != 0 {
						p.Item.Use[0] = 2
					} else {
						copy(p.Item.Use[:], "zombie")
					}
					switch variant {
					case 1:
						entry.Type = 16
					case 2:
						entry.Type = 0x1000f
					case 3, 4, 5, 6:
						entry.Mods[variant-3] = !entry.Mods[variant-3]
					case 7:
						entry.Reward = 3
					case 8:
						p.Item.Use[0] = 0
					case 9:
						p.NilItem = op == 1
					case 10:
						p.NilEntry = op == 1
					case 11:
						p.Session = 3
					}
					p.Stock = []legacy.PortTestShopStock{{Type: 0}, entry, entry}
					p.Entry = 1
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "shop-stock-matching", legacy.PortTestRoam(specs), shopHashes["shop-stock-matching"])
}

var shopHashes = map[string]string{
	"shop-price-base":             "ff58691c15d75d5fae46e9b0a82dc0af4567367afc0246a5018647e2f549be03",
	"shop-price-modifiers-health": "66d46b7a12620f38a8070c8e3cfe9dd820d4f4ffdda870ef34b9caeca772b1e2",
	"shop-price-charges":          "46bc63ad25a97728b576ebca116ad4dac3854836a5d43c1001319df42c5d865e",
	"shop-price-rewards-gems":     "c4db2f46759f357a001dc51842a20dcb0bb3c9a04586cbd452b15a5f714b1ca1",
	"shop-stock-matching":         "50295d5d21fc1401d8db3bf87efa3e01249bd3082d13334d8dccabd5897143de",
	"shop-price-rounding":         "e72de64940723568b9e6ebad7b4e3ffaa0543e90ae45413357c81f6daa7b570d",
	"shop-stock-priority":         "7fed15c37ed9c260b64fe7555455040e634c80daa73a253c68a29a6bb7ba2aaa",
	"shop-session-sequences":      "5532a71d77bc1783e7f23aa05d3325893368f20c795fa83cafe53141e886efa2",
	"shop-stock-sequences":        "0b9dffef335e09b9809c9df55a83e49152120664e2cd8f4ef3524a862665b617",
	"shop-item-capacity":          "d000d582bc524f06ca1d28cd96afd7c5659bada0c363f24e9d3099ac030a5214",
	"shop-packets":                "051d65dfef6ab3b4349f59c959e2d2515c8176063943d9b19342264992063d41",
	"shop-gold-balance":           "9019e9f97a408eadff11f9722e7781d653c3cee26d17c16c6cb6431403825cd3",
	"shop-trade-completion":       "8cde6a61f8532c01768bd1d46c56b6cfc555f7cf88703ad90a1d347ded7c1b6e",
}

func TestShopPriceRounding(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	values := []float32{0, -1, .1, 1.0 / 3, 1, 1.01, float32(math.Inf(1)), float32(math.NaN())}
	for i := uint32(0); i < 384; i++ {
		s := shopBase(0)
		p := s.Callbacks.Shop
		p.Mode = int(i % 3)
		p.Buy, p.Sell = math.Float32bits(values[i%8]), math.Float32bits(values[(i/8)%8])
		p.Item.Worth = 16777000 + i*137
		p.Item.Health = true
		p.Item.HP, p.Item.MaxHP = uint16(i*53), uint16(1+i*97)
		p.Balance["RepairCoefficient"] = float64(values[(i/48)%8])
		specs = append(specs, s)
	}
	callbackHash(t, "shop-price-rounding", legacy.PortTestRoam(specs), shopHashes["shop-price-rounding"])
}

func TestShopStockPriority(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{0, 8, 0x1000, 0x1000000, 0x2000000, 0x10000000, 0x100, 0xffffffff} {
		for _, subclass := range []uint32{0, 1, 2, 4, 8, 0x82, 0x47f0000, 0xffffffff} {
			for _, value := range []uint32{0, 1, 101, 0xffffff, 0x1000000, 0x80000000, 0xffffffff} {
				s := shopBase(3)
				s.Callbacks.Shop.Item.Class, s.Callbacks.Shop.Item.Subclass, s.Callbacks.Shop.Priority = class, subclass, value
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-stock-priority", legacy.PortTestRoam(specs), shopHashes["shop-stock-priority"])
}

func TestShopSessionSequences(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, order := range [][]int{{0, 1, 2}, {2, 1, 0}, {1, 0, 2}, {1, 2, 0}} {
		for _, mode := range []uint32{0, 1, 2} {
			s := shopBase(4)
			p := s.Callbacks.Shop
			for i := 0; i < 3; i++ {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate, Value: mode})
			}
			for _, i := range order {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopDestroy, Session: i})
			}
			p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate}, legacy.PortTestShopAction{Op: legacy.PortTestShopReset}, legacy.PortTestShopAction{Op: legacy.PortTestShopReset})
			specs = append(specs, s)
		}
	}
	s := shopBase(4)
	p := s.Callbacks.Shop
	for i := 0; i < 66; i++ {
		p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate})
	}
	p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopDestroy, Session: 31}, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate}, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate}, legacy.PortTestShopAction{Op: legacy.PortTestShopReset}, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate})
	specs = append(specs, s)
	r := legacy.PortTestRoam(specs)
	for _, v := range r[:12] {
		steps := v.Callbacks.Shop.Sequence
		if steps[5].Alive != 0 || len(steps[5].Sessions) != 0 {
			t.Fatal("session destruction did not free all gold and sessions")
		}
		// The original bulk reset releases pool records, but retains the two
		// gold objects; the fixture owns and frees those orphaned allocations.
		if steps[7].Alive != 2 || steps[7].Head != 0 {
			t.Fatal("bulk reset lifetime changed")
		}
	}
	steps := r[len(r)-1].Callbacks.Shop.Sequence
	if steps[63].Return == 0 || steps[64].Return != 0 || steps[65].Return != 0 || steps[67].Return == 0 || steps[68].Return != 0 || steps[70].Return == 0 {
		t.Fatal("64-record session pool exhaustion/reuse")
	}
	callbackHash(t, "shop-session-sequences", r, shopHashes["shop-session-sequences"])
}

func TestShopStockSequences(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, order := range [][]int{{0, 1, 2, 3, 4, 5}, {5, 4, 3, 2, 1, 0}, {3, 1, 5, 0, 4, 2}} {
		for mode := uint32(0); mode < 3; mode++ {
			for cleanup := 0; cleanup < 3; cleanup++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				for _, worth := range []uint32{101, 0, 101, 55, 201, 1} {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: worth})
				}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate, Value: mode})
				for _, i := range order {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
				}
				for _, id := range []uint32{70000, 70001, 70005, 0, 0xffffffff} {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopFind, Value: id})
				}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: []int{legacy.PortTestShopDestroy, legacy.PortTestShopReset, legacy.PortTestShopFreeList}[cleanup]})
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		if len(steps[6].Nodes) != 6 {
			t.Fatalf("case %d stock count", i)
		}
		for j := 7; j < 10; j++ {
			if steps[j].Return == 0 {
				t.Fatalf("case %d missing stock lookup", i)
			}
		}
		if steps[10].Return != 0 || steps[11].Return != 0 {
			t.Fatalf("case %d unexpected stock lookup", i)
		}
		// All these objects have the same category, so node values must be
		// ascending. Equal prices retain the original insert-before behavior.
		for j := 1; j < 6; j++ {
			if steps[6].Nodes[j-1][1] > steps[6].Nodes[j][1] {
				t.Fatalf("case %d unsorted prices", i)
			}
		}
	}
	callbackHash(t, "shop-stock-sequences", r, shopHashes["shop-stock-sequences"])
}

func TestShopItemPoolCapacity(t *testing.T) {
	s := shopBase(4)
	p := s.Callbacks.Shop
	p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopCreate, Value: 1})
	for i := 0; i < 501; i++ {
		p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: uint32(100 + i%11)})
		p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
	}
	p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopFreeList}, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd}, legacy.PortTestShopAction{Op: legacy.PortTestShopDestroy})
	r := legacy.PortTestRoam([]legacy.PortTestRoamSpec{s})
	steps := r[0].Callbacks.Shop.Sequence
	if steps[500].Return == 0 || len(steps[500].Nodes) != 500 || steps[501].Return != 0 || steps[503].Return == 0 || len(steps[503].Nodes) != 1 {
		t.Fatal("500-record item pool exhaustion/reuse")
	}
	callbackHash(t, "shop-item-capacity", r, shopHashes["shop-item-capacity"])
}

func TestShopPackets(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, peer := range []bool{false, true} {
		for session := 2; session <= 3; session++ {
			for mask := uint32(0); mask < 16; mask++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				p.Session, p.Peer = session, peer
				p.Priority = 0xabcdef01
				p.Item.Class, p.Item.Health, p.Item.HP, p.Item.MaxHP = 0x1000000, true, 123, 200
				for j := range p.Item.Mods {
					p.Item.Mods[j] = mask&(1<<j) != 0
				}
				p.Sequence = []legacy.PortTestShopAction{
					{Op: legacy.PortTestShopCreate, Value: 1},
					{Op: legacy.PortTestShopGold, Value: 0x12345678}, {Op: legacy.PortTestShopGold, Side: 1, Value: 0xfedcba98},
					{Op: legacy.PortTestShopSet, Item: 6, Value: []uint32{0, 1, 2, 0xffffffff}[mask%4]},
					{Op: legacy.PortTestShopSet, Item: 7, Value: []uint32{0, 1, 2, 0xffffffff}[mask/4]},
					{Op: legacy.PortTestShopSet, Item: 10, Value: 7}, {Op: legacy.PortTestShopSet, Item: 11, Value: 3},
				}
				sides := []int{3 - session}
				if peer {
					sides = []int{0, 1}
				}
				for _, side := range sides {
					for packet := 0; packet < 8; packet++ {
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopPacket, Side: side, Item: packet})
					}
				}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		packets := steps[len(steps)-1].Packets
		want := 8
		if specs[i].Callbacks.Shop.Peer {
			want = 16
		}
		if len(packets) != want {
			t.Fatalf("case %d packets=%d want%d", i, len(packets), want)
		}
		for j, pkt := range packets {
			kind := 7 - j%8 // retained important queue is a head-insert list.
			if len(pkt.Data) != []int{4, 18, 2, 2, 2, 3, 14, 4}[kind] || pkt.Data[0] != 0xc9 || pkt.Data[1] != []byte{9, 8, 1, 2, 7, 3, 6, 5}[kind] {
				t.Fatalf("case %d packet %d shape %v", i, j, pkt)
			}
			if (pkt.Ordered != 0) != (kind != 2) || pkt.A4 != 0 || pkt.A5 != 1 {
				t.Fatalf("case %d packet %d metadata", i, j)
			}
		}
	}
	callbackHash(t, "shop-packets", r, shopHashes["shop-packets"])
}

func TestShopGoldBalance(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for session := 2; session <= 3; session++ {
		for _, gold := range []uint32{0, 1, 50, 1000, 0xffffffff} {
			for variant := 0; variant < 8; variant++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				p.Session, p.Gold, p.ProtectedGold = session, [2]uint32{gold, 100}, variant&4 != 0
				p.Peer = variant == 7
				p.Items = []legacy.PortTestShopItem{{Type: 15, Class: 8, Worth: 101}, {Type: 15, Class: 8, Worth: 203}}
				p.Sequence = []legacy.PortTestShopAction{
					{Op: legacy.PortTestShopCreate, Value: 1},
					{Op: legacy.PortTestShopGold, Value: 7}, {Op: legacy.PortTestShopGold, Side: 1, Value: 11},
					{Op: legacy.PortTestShopOffer, Item: 0, Value: []uint32{0, 101, 300, 0xffffffff}[variant%4]},
					{Op: legacy.PortTestShopOffer, Item: 1, Side: 1, Value: []uint32{0, 203, 99, 1}[variant%4]},
					{Op: legacy.PortTestShopTotal}, {Op: legacy.PortTestShopTotal, Side: 1},
					{Op: legacy.PortTestShopBalance}, {Op: legacy.PortTestShopBalance},
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-gold-balance", legacy.PortTestRoam(specs), shopHashes["shop-gold-balance"])
}

func TestShopTradeCompletion(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, peer := range []bool{false, true} {
		for session := 2; session <= 3; session++ {
			for variant := 0; variant < 6; variant++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				p.Peer, p.Session, p.Gold = peer, session, [2]uint32{100, 200}
				mode := uint32(1)
				if peer {
					mode = 0
				}
				p.Items = []legacy.PortTestShopItem{{Type: 15, Class: 8, Worth: 101}, {Type: 15, Class: 8, Worth: 203}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: mode}, {Op: legacy.PortTestShopGold, Value: 7}, {Op: legacy.PortTestShopGold, Side: 1, Value: 11}}
				if variant%3 > 0 {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Item: 0, Value: 101})
				}
				if variant%3 > 1 {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Item: 1, Side: 1, Value: 203})
				}
				if variant < 3 {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAccept}, legacy.PortTestShopAction{Op: legacy.PortTestShopAccept, Side: 1})
				} else {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopWithdraw, Value: 70000}, legacy.PortTestShopAction{Op: legacy.PortTestShopWithdraw, Value: 0xffffffff}, legacy.PortTestShopAction{Op: legacy.PortTestShopCancel})
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-trade-completion", legacy.PortTestRoam(specs), shopHashes["shop-trade-completion"])
}
