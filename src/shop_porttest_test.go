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
	"shop-repairs":                "6ac553918b85a1f16e425cccad321c0b3c02312974dcc239221c1bc3bf204a0d",
	"shop-sales":                  "39cb08269918856425f14693d10a21e11818814bb0e8854fb33693d6ed19ffbf",
	"shop-cached-sessions":        "e55154a1f30a2e1e88e448822d47162e86752d6efd2bdaae0aed3145a44e8adf",
	"shop-stock-loading":          "bf808c7f64995d7fa8497e44a43ec367d96f03cd3c890ddf99faaef51cc3a582",
	"shop-quest-loading":          "dc0cfa85fac4eee21f50cb3f4cbddec8e21ec805e2465d485e5f081c2b7a38f5",
	"shop-offer-removal":          "6ed33429c83d82ab61d0c15638a026e07b07f84a90c890ead4cf625bd96325fd",
	"shop-quest-price-rounding":   "db88721cc49d89cc7e569e1b0e15b6a4e668d15eba5ef33d0f088f3257803861",
	"shop-stock-boundaries":       "2a9acc334312f6e7c768909a9e3cefbe7633090c53e6e7c52f03a13be78cd773",
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
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		if i%6 >= 3 {
			continue
		}
		sp := specs[i].Callbacks.Shop
		steps := v.Callbacks.Shop.Sequence
		last := steps[len(steps)-1]
		want := uint32(107)
		if sp.Session == 3 || sp.Peer {
			want = 111
		}
		if last.Players[2][541] != want {
			t.Fatalf("case%d completed gold=%d want%d", i, last.Players[2][541], want)
		}
		if sp.Peer {
			if last.Players[5][541] != 207 || last.Head != 0 {
				t.Fatalf("case%d peer completion state", i)
			}
			wantA, wantB := uint32(0), uint32(0)
			if i%6 >= 1 {
				wantB = 70000
			}
			if i%6 >= 2 {
				wantA = 70001
			}
			if last.Players[0][126] != wantA || last.Players[3][126] != wantB {
				t.Fatalf("case%d inventory delivery", i)
			}
		}
	}
	callbackHash(t, "shop-trade-completion", r, shopHashes["shop-trade-completion"])
}

func TestShopRepairs(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for session := 2; session <= 3; session++ {
		for _, hp := range [][2]uint16{{0, 100}, {1, 100}, {99, 100}, {100, 100}, {101, 100}, {0, 0}} {
			for variant := 0; variant < 4; variant++ {
				for _, gold := range []uint32{0, 1, 1000} {
					s := shopBase(4)
					p := s.Callbacks.Shop
					p.CaptureData, p.Session, p.Gold, p.ProtectedGold = true, session, [2]uint32{gold, 100}, variant&1 != 0
					item := legacy.PortTestShopItem{Type: 15, Class: 8, Worth: 100, Health: true, HP: hp[0], MaxHP: hp[1]}
					if variant >= 2 {
						item.Class, item.Subclass = 0x1000, 0x40000
						item.Use[108], item.Use[109] = byte(variant-2)*20, 20
						item.Use[112] = byte(variant-2) * 100
					}
					p.Items = []legacy.PortTestShopItem{item}
					side := 3 - session
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopRepairQuote, Side: side, Value: 70000}, {Op: legacy.PortTestShopRepair, Side: side, Value: 70000}, {Op: legacy.PortTestShopInventory, Side: side}, {Op: legacy.PortTestShopRepairQuote, Side: side, Value: 70000}, {Op: legacy.PortTestShopRepairQuote, Side: side, Value: 0xffffffff}, {Op: legacy.PortTestShopRepair, Side: side, Value: 0xffffffff}, {Op: legacy.PortTestShopRepair, Side: side, Value: 70000}, {Op: legacy.PortTestShopRepairQuote, Side: side, Value: 70000}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		for _, data := range steps[7].ObjectData {
			if data[0] != 70000 {
				continue
			}
			hp := data[len(data)-2:]
			if uint16(hp[0]) != uint16(hp[1]) {
				t.Fatalf("case%d repair did not set current to maximum: %v", i, hp)
			}
		}
	}
	callbackHash(t, "shop-repairs", r, shopHashes["shop-repairs"])
}

func TestShopSales(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for session := 2; session <= 3; session++ {
		for _, qty := range []uint32{0, 1, 2, 3, 4, 0xffffffff} {
			for _, typ := range []int{15, 16, 17, 0x1000f} {
				for _, protected := range []bool{false, true} {
					s := shopBase(4)
					p := s.Callbacks.Shop
					p.Session, p.Gold, p.ProtectedGold = session, [2]uint32{100, 0}, protected
					p.Items = []legacy.PortTestShopItem{{Type: 15, Class: 8, Worth: 100}, {Type: 16, Class: 8, Worth: 200}, {Type: 15, Class: 8, Worth: 100}}
					side := 3 - session
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopInventory, Side: side, Item: 0}, {Op: legacy.PortTestShopInventory, Side: side, Item: 1}, {Op: legacy.PortTestShopInventory, Side: side, Item: 2}, {Op: legacy.PortTestShopSell, Side: side, Item: typ, Value: qty}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		p := specs[i].Callbacks.Shop
		a := p.Sequence[4]
		want := uint32(100)
		if a.Item == 15 {
			want += 50 * min(a.Value, 2)
		} else if a.Item == 16 && a.Value != 0 {
			want += 100
		}
		if got := v.Callbacks.Shop.Sequence[4].Players[2][541]; got != want {
			t.Fatalf("case%d gold=%d want%d", i, got, want)
		}
	}
	callbackHash(t, "shop-sales", r, shopHashes["shop-sales"])
}

func TestShopCachedSessions(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for session := 2; session <= 3; session++ {
		for _, quest := range []uint32{0, 4096} {
			for count := 0; count < 4; count++ {
				s := shopBase(4)
				s.Lifecycle.GameFlags = quest
				p := s.Callbacks.Shop
				p.Session = session
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}}
				for i := 0; i < count; i++ {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: 100})
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
				}
				side := 3 - session
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopLookup, Side: side, Value: 70000}, legacy.PortTestShopAction{Op: legacy.PortTestShopLookup, Side: side, Value: 0xffffffff}, legacy.PortTestShopAction{Op: legacy.PortTestShopDetach, Side: side}, legacy.PortTestShopAction{Op: legacy.PortTestShopLookup, Side: side, Value: 70000}, legacy.PortTestShopAction{Op: legacy.PortTestShopCancel}, legacy.PortTestShopAction{Op: legacy.PortTestShopPlayerCleanup, Item: 1}, legacy.PortTestShopAction{Op: legacy.PortTestShopPlayerCleanup, Item: 1})
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-cached-sessions", legacy.PortTestRoam(specs), shopHashes["shop-cached-sessions"])
}

func TestShopStockLoading(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for session := 2; session <= 3; session++ {
		for count := byte(0); count < 4; count++ {
			for mask := 0; mask < 16; mask++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				p.Session, p.CaptureData = session, true
				mods := [4]bool{}
				for i := range mods {
					mods[i] = mask&(1<<i) != 0
				}
				p.Item.ModPrice = [4]int32{7, 11, 13, 17}
				p.Load = &legacy.PortTestShopLoadSpec{}
				p.Stock = []legacy.PortTestShopStock{{Type: 23, Count: count}, {Type: 27, Count: 1, Reward: 1}, {Type: 28, Count: 1, Reward: 5}, {Type: 29, Count: 1, Reward: 2}, {Type: 32, Count: count, Mods: mods}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopLoad}, {Op: legacy.PortTestShopDestroy}}
				specs = append(specs, s)
			}
		}
	}
	large := shopBase(4)
	large.Callbacks.Shop.CaptureData = true
	large.Callbacks.Shop.Load = &legacy.PortTestShopLoadSpec{}
	large.Callbacks.Shop.Stock = []legacy.PortTestShopStock{{Type: 23, Count: 255}}
	large.Callbacks.Shop.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopLoad}, {Op: legacy.PortTestShopDestroy}}
	specs = append(specs, large)
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		want := 0
		for _, entry := range specs[i].Callbacks.Shop.Stock {
			want += int(entry.Count)
		}
		if got := len(v.Callbacks.Shop.Sequence[1].Nodes); got != want {
			t.Fatalf("case%d stock=%d want%d", i, got, want)
		}
	}
	callbackHash(t, "shop-stock-loading", r, shopHashes["shop-stock-loading"])
}

func TestShopQuestLoading(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, stage := range []uint32{0, 1, 8, 9, 10, 11, 0xffffffff} {
		for _, seed := range []int{1, 2, 3, 7, 15, 31} {
			for variant := 0; variant < 4; variant++ {
				s := shopBase(4)
				s.Seed, s.Lifecycle.GameFlags = seed, 4096
				p := s.Callbacks.Shop
				p.CaptureData = true
				p.Balance["ShopAnkhCutoffStage"] = 10
				p.Load = &legacy.PortTestShopLoadSpec{Stage: stage, Marker: variant != 0, MarkerChance: []uint32{0, 0, 1, 4}[variant], Names: []string{"Diamond", "missing", "Gold"}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopLoad}, {Op: legacy.PortTestShopDestroy}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "shop-quest-loading", legacy.PortTestRoam(specs), shopHashes["shop-quest-loading"])
}

func TestShopOfferRemoval(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, peer := range []bool{false, true} {
		for session := 2; session <= 3; session++ {
			for remove := 0; remove < 6; remove++ {
				s := shopBase(4)
				p := s.Callbacks.Shop
				p.Peer, p.Session, p.Gold, p.ProtectedGold = peer, session, [2]uint32{1000, 2000}, true
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}, {Op: legacy.PortTestShopSet, Item: 6, Value: 1}, {Op: legacy.PortTestShopSet, Item: 7, Value: 1}}
				for i, value := range []uint32{101, 0x80000000, 0xffffffff, 201, 0} {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: uint32(100 + i)})
					side := 0
					if i >= 3 {
						side = 1
					}
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Item: i, Side: side, Value: value})
				}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopWithdraw, Value: uint32(70000 + remove)}, legacy.PortTestShopAction{Op: legacy.PortTestShopWithdraw, Value: uint32(70000 + remove)})
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		want := uint32(1)
		if i%6 == 5 {
			want = 0
		}
		if steps[8].Return != want || steps[9].Return != 0 {
			t.Fatalf("case%d repeated withdrawal", i)
		}
	}
	callbackHash(t, "shop-offer-removal", r, shopHashes["shop-offer-removal"])
}

func TestShopQuestPriceRounding(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for i := uint32(0); i < 256; i++ {
		s := shopBase(0)
		s.Lifecycle.GameFlags = 4096
		p := s.Callbacks.Shop
		p.Mode = int(i % 3)
		p.Item.Class = 0x1000000
		p.Item.Worth = 16777200 + i*113
		p.Item.Mods = [4]bool{true, true, i&1 != 0, i&2 != 0}
		p.Item.ModPrice = [4]int32{16777217, -16777215, 101, -37}
		values := []float64{.1, 1.00000003, 1.0 / 3, 1000000000001.25, -1.1, 0, .00000001, 1.01}
		p.Balance["QuestModifierWorthMultiplier"] = values[i%8]
		p.Balance["QuestSellMultiplier"] = values[(i/8)%8]
		p.Item.Health = true
		p.Item.HP, p.Item.MaxHP = uint16(1+i*37), uint16(1+i*97)
		specs = append(specs, s)
	}
	callbackHash(t, "shop-quest-price-rounding", legacy.PortTestRoam(specs), shopHashes["shop-quest-price-rounding"])
}

func TestShopStockBoundaries(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, n := range []int{0, 1, 59, 60} {
		for session := 2; session <= 3; session++ {
			for _, match := range []bool{false, true} {
				s := shopBase(2)
				p := s.Callbacks.Shop
				p.Session = session
				p.Stock = make([]legacy.PortTestShopStock, n)
				if match && n > 0 {
					p.Stock[n-1].Type = 15
				}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		want := uint32(0xffffffff)
		n := len(specs[i].Callbacks.Shop.Stock)
		if i%2 != 0 && n > 0 {
			want = uint32(n - 1)
		}
		if v.Callbacks.Return != want {
			t.Fatalf("case%d index=%d want%d", i, v.Callbacks.Return, want)
		}
	}
	callbackHash(t, "shop-stock-boundaries", r, shopHashes["shop-stock-boundaries"])
}
