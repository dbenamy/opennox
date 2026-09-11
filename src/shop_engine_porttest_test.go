//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

var tradeEngineHashes = map[string]string{
	"trade-engine-cache-counts":       "e11b0b40ec57f1cdc068b93444b4f96607d182ad4fecf9f8486420384a4aa73b",
	"trade-engine-cached-reopening":   "381f01d471de53269bce7ff819c17d7a1dab2e8a39118bf79b633e285869c141",
	"trade-engine-offer-admission":    "f21a5229a2d451afd9e82556c3ccc8d2de8486836333662a9e5f26995a56f9df",
	"trade-engine-offer-mutation":     "69918d948f63a9912330a85a8f3dfca2a83753543441299205e5b4c0ceb1bca8",
	"trade-engine-opening-rejections": "026b1f1c3e384346d021e2ca3112c0a93220e832bd673922277946f6d895024c",
	"trade-engine-opening":            "2e4100a183199e7f503ab56dd1777f67744e78956dc1666d61e746839d4d9007",
	"trade-engine-packets":            "b9cfac932a354ca3c3d761a07fe105e1beff5aab3e378b100570cc4f5f170210",
	"trade-engine-pool-exhaustion":    "6790d4e96855ba426ffec8faaf4cfe2a59c43002115e5e0aa92f80c2f8168be3",
	"trade-engine-purchase-limits":    "087ca9290c0d989afd862466a48b66cbfe36509acaaa8f0fa2ad107bfe31a3c6",
	"trade-engine-purchases":          "0cf4822d84144f0d73721668e0a6c144db7a126246411f7040f06e5c7e07fec4",
	"trade-engine-sales":              "c49cb55cfeeca7b9e95c45267bd4ddbc7907066f640046b5fd3c384a81def891",
	"trade-engine-stock-remove":       "4a71e2a77a6698cfa6fa7158e73cc9d68e23b18579b1651502e6a7918d51a514",
}

func tradeEngineBase() legacy.PortTestRoamSpec {
	s := shopBase(4)
	p := s.Callbacks.Shop
	p.Engine = &legacy.PortTestShopEngineSpec{Names: [2]string{"Ada", "Grace"}, VendorName: "merchant"}
	p.CaptureData = true
	p.Gold = [2]uint32{1000, 2000}
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestShopCreate, Value: 1}}
	return s
}
func TestShopEngineStockRemoval(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for orientation := 2; orientation <= 3; orientation++ {
		for count := 0; count <= 5; count++ {
			for selected := 0; selected < 6; selected++ {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				p.Session = orientation
				for i := 0; i < 6; i++ {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: uint32(10 + i)})
				}
				for i := 0; i < count; i++ {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
				}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeRemove, Item: selected}, legacy.PortTestShopAction{Op: legacy.PortTestTradeRemove, Item: selected})
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		if steps[len(steps)-1].Return != 0 {
			t.Fatalf("case %d repeated removal succeeded", i)
		}
	}
	callbackHash(t, "trade-engine-stock-remove", r, tradeEngineHashes["trade-engine-stock-remove"])
}
func TestShopEngineOfferAdmission(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for count := 0; count <= 4; count++ {
		for _, typ := range []uint16{0, 15, 16, 17, 18, 19, 65535} {
			for _, class := range []uint32{8, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				p.Peer = true
				p.Sequence[0].Value = 0
				for i := 0; i < count; i++ {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: uint16(15 + i), Class: 8, Worth: 10})
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Item: i, Value: 10})
				}
				p.Items = append(p.Items, legacy.PortTestShopItem{Type: typ, Class: class, Worth: 20})
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeOfferAllowed, Item: count}, legacy.PortTestShopAction{Op: legacy.PortTestTradeOfferAllowed, Item: count, Side: 1})
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		if steps[len(steps)-1].Return != 1 {
			t.Fatalf("case %d empty offer list rejected item", i)
		}
	}
	callbackHash(t, "trade-engine-offer-admission", r, tradeEngineHashes["trade-engine-offer-admission"])
}
func TestShopEnginePackets(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, peer := range []bool{false, true} {
		for variant := 0; variant < 16; variant++ {
			s := tradeEngineBase()
			p := s.Callbacks.Shop
			p.Peer = peer
			side := 1
			if peer {
				p.Sequence[0].Value = 0
				side = variant % 2
			}
			p.Items = []legacy.PortTestShopItem{{Type: 65535, Class: []uint32{8, 0x1000, 0x1000000, 0x2000000}[variant%4], Worth: []uint32{0, 100, 0x80000000, 0xffffffff}[variant%4]}}
			for i := range p.Items[0].Mods {
				p.Items[0].Mods[i] = variant&(1<<i) != 0
			}
			p.Engine.Names = [2]string{[]string{"", "Ada", "abcdefghijklmnopqrstuvwx", "李明"}[variant%4], "Grace"}
			p.Engine.VendorName = []string{"", "merchant", "abcdefghijklmnopqrstuvwxyzABCDE", "shop"}[variant%4]
			p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeSetPlayer, Side: side}, legacy.PortTestShopAction{Op: legacy.PortTestTradeShortfall, Side: side, Value: []uint32{0, 1, 65535, 65536}[variant%4]}, legacy.PortTestShopAction{Op: legacy.PortTestTradeOfferPacket, Side: side, Value: uint32(variant % 2)}, legacy.PortTestShopAction{Op: legacy.PortTestTradeSendStock, Side: side})
			if peer {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradePeerIntro, Side: side})
			} else {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeIntro, Side: side})
			}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "trade-engine-packets", legacy.PortTestRoam(specs), tradeEngineHashes["trade-engine-packets"])
}
func TestShopEngineCachesAndStockCounts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, quest := range []uint32{0, 4096} {
		for _, typ := range []uint16{0, 15, 23, 24, 25, 30, 65535} {
			for _, count := range []byte{0, 1, 2, 255} {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				s.Lifecycle.GameFlags = quest
				p.Items = []legacy.PortTestShopItem{{Type: typ, Class: 8, Worth: 10}}
				p.Stock = []legacy.PortTestShopStock{{Type: 16, Count: 3}, {Type: uint32(typ), Count: count}, {Type: 17, Count: 4}}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeStockRemovable}, legacy.PortTestShopAction{Op: legacy.PortTestTradeIsGem}, legacy.PortTestShopAction{Op: legacy.PortTestTradeStockDecrement}, legacy.PortTestShopAction{Op: legacy.PortTestTradeStockDecrement})
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "trade-engine-cache-counts", legacy.PortTestRoam(specs), tradeEngineHashes["trade-engine-cache-counts"])
}

func TestShopEngineOfferMutation(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, peer := range []bool{false, true} {
		for count := 0; count <= 4; count++ {
			for variant := 0; variant < 8; variant++ {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				p.Peer = peer
				if peer {
					p.Sequence[0].Value = 0
				}
				side := variant % 2
				p.Gold = [2]uint32{[]uint32{0, 20, 100, 0xffffffff}[variant%4], 200}
				for i := 0; i < count; i++ {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: uint16(15 + i), Class: 8, Worth: 20})
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Side: side, Item: i, Value: uint32(10 + i)})
				}
				typ := uint16(19)
				if variant&2 != 0 {
					typ = 15
				}
				class := uint32(8)
				if variant&4 != 0 {
					class = 0x1000
				}
				p.Items = append(p.Items, legacy.PortTestShopItem{Type: typ, Class: class, Worth: 100})
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeAddOffer, Side: side, Item: count})
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "trade-engine-offer-mutation", legacy.PortTestRoam(specs), tradeEngineHashes["trade-engine-offer-mutation"])
}
func TestShopEnginePurchases(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, quest := range []uint32{0, 4096} {
		for variant := 0; variant < 16; variant++ {
			for _, many := range []bool{false, true} {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				s.Lifecycle.GameFlags = quest
				p.Balance["MaxExtraLives"], p.Balance["ForceOfNatureStaffLimit"] = 3, 2
				p.Engine.ExtraLives = uint32(variant % 4)
				typ := []uint16{15, 23, 24, 30}[variant%4]
				p.Gold = [2]uint32{[]uint32{0, 124, 125, 1000}[variant/4], 2000}
				for i := 0; i < 3; i++ {
					p.Items = append(p.Items, legacy.PortTestShopItem{Type: typ, Class: 8, Worth: 100, Pickup: variant&4 != 0})
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
				}
				p.Stock = []legacy.PortTestShopStock{{Type: uint32(typ), Count: 3}}
				if many {
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeBuyMany, Side: 1, Value: uint32(variant % 4)})
				} else {
					code := uint32(70000)
					if variant&8 != 0 {
						code = 999
					}
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeBuy, Side: 1, Value: code})
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "trade-engine-purchases", legacy.PortTestRoam(specs), tradeEngineHashes["trade-engine-purchases"])
}
func TestShopEngineSales(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, quest := range []uint32{0, 4096} {
		for variant := 0; variant < 24; variant++ {
			s := tradeEngineBase()
			p := s.Callbacks.Shop
			s.Lifecycle.GameFlags = quest
			p.Items = []legacy.PortTestShopItem{{Type: 15, Class: 8, Worth: 100}, {Type: 15, Class: 8, Worth: 201}}
			p.Engine.Cache[10], p.Engine.Cache[11] = 14, 14
			if variant%3 == 1 {
				p.Engine.NoSellType = 15
			}
			if variant%3 == 2 {
				p.Engine.Cache[10], p.Engine.Cache[11] = 15, 15
			}
			if variant&4 == 0 {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopInventory, Side: 1}, legacy.PortTestShopAction{Op: legacy.PortTestShopInventory, Item: 1, Side: 1})
			}
			code := uint32(70000)
			if variant&8 != 0 {
				code = 999
			}
			p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeSellQuote, Side: 1, Value: code}, legacy.PortTestShopAction{Op: legacy.PortTestTradeSell, Side: 1, Value: code}, legacy.PortTestShopAction{Op: legacy.PortTestTradeSell, Side: 1, Value: code})
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		sp := specs[i].Callbacks.Shop
		steps := v.Callbacks.Shop.Sequence
		want := sp.Gold[0]
		stocked := false
		for _, a := range sp.Sequence {
			if a.Op == legacy.PortTestShopInventory {
				stocked = true
			}
		}
		code := sp.Sequence[len(sp.Sequence)-1].Value
		if stocked && code == 70000 && sp.Engine.NoSellType != 15 && sp.Engine.Cache[11] != 15 {
			if specs[i].Lifecycle.GameFlags&4096 != 0 {
				want += 75
			} else {
				want += 50
			}
		}
		for _, j := range []int{len(steps) - 2, len(steps) - 1} {
			if got := steps[j].Players[2][541]; got != want {
				t.Fatalf("case %d sale/repeated sale gold=%d want %d", i, got, want)
			}
		}
	}
	callbackHash(t, "trade-engine-sales", r, tradeEngineHashes["trade-engine-sales"])
}
func TestShopEngineOpening(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 2048, 4096, 6144} {
		for _, peer := range []bool{false, true} {
			for side := 0; side < 2; side++ {
				s := tradeEngineBase()
				p := s.Callbacks.Shop
				p.Peer = peer
				s.Lifecycle.GameFlags = flags
				p.Balance["ShopAnkhCutoffStage"] = 0
				p.Load = &legacy.PortTestShopLoadSpec{}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTradeStart, Side: side}, {Op: legacy.PortTestTradeStart, Side: side}}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		if steps[0].Return == 0 || steps[1].Return != 0 {
			t.Fatalf("case %d opening / already-trading result mismatch", i)
		}
	}
	callbackHash(t, "trade-engine-opening", r, tradeEngineHashes["trade-engine-opening"])
}

func TestShopEnginePurchaseLimits(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 2048, 4096, 6144} {
		for _, count := range []int{0, 1, 2, 3, 8, 9} {
			for _, many := range []bool{false, true} {
				for variant := 0; variant < 4; variant++ {
					s := tradeEngineBase()
					p := s.Callbacks.Shop
					s.Lifecycle.GameFlags = flags
					p.Gold[0] = 10000
					p.Balance["MaxExtraLives"], p.Balance["ForceOfNatureStaffLimit"] = 3, 2
					item := legacy.PortTestShopItem{Type: 15, Class: 0x10, Worth: 100, Pickup: variant&1 != 0}
					if variant&2 != 0 {
						item.Class = 0x1000
						item.Subclass = 0x200000
						item.Use[108], item.Use[109] = 1, 1
					}
					p.Items = append(p.Items, item)
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd})
					for i := 0; i < count; i++ {
						p.Items = append(p.Items, item)
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopInventory, Item: i + 1, Side: 1})
					}
					if many {
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeBuyMany, Side: 1, Value: 1})
					} else {
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeBuy, Side: 1, Value: 70000})
					}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "trade-engine-purchase-limits", legacy.PortTestRoam(specs), tradeEngineHashes["trade-engine-purchase-limits"])
}
func TestShopEngineCachedReopening(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 2048, 4096, 6144} {
		for side := 0; side < 2; side++ {
			s := tradeEngineBase()
			p := s.Callbacks.Shop
			s.Lifecycle.GameFlags = flags
			p.Balance["ShopAnkhCutoffStage"] = 0
			p.Load = &legacy.PortTestShopLoadSpec{}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTradeStart, Side: side}, {Op: legacy.PortTestShopCancel}, {Op: legacy.PortTestTradeStart, Side: side}, {Op: legacy.PortTestShopCancel, Session: 1}}
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		quest := specs[i].Lifecycle.GameFlags&4096 != 0
		if (steps[0].Return == steps[2].Return) != quest {
			t.Fatalf("case %d cached session identity", i)
		}
		if quest && steps[1].Cached[1] != steps[0].Return {
			t.Fatalf("case %d cache after exit", i)
		}
	}
	callbackHash(t, "trade-engine-cached-reopening", r, tradeEngineHashes["trade-engine-cached-reopening"])
}
func TestShopEnginePoolExhaustion(t *testing.T) {
	s := tradeEngineBase()
	p := s.Callbacks.Shop
	for i := 0; i < 501; i++ {
		p.Items = append(p.Items, legacy.PortTestShopItem{Type: 15, Class: 8, Worth: 10})
		if i < 500 {
			p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
		}
	}
	p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeAddOffer, Item: 500, Side: 1}, legacy.PortTestShopAction{Op: legacy.PortTestTradeRemove, Item: 0}, legacy.PortTestShopAction{Op: legacy.PortTestTradeAddOffer, Item: 500, Side: 1})
	r := legacy.PortTestRoam([]legacy.PortTestRoamSpec{s})
	steps := r[0].Callbacks.Shop.Sequence
	if steps[len(steps)-3].Return != 0 || steps[len(steps)-1].Return != 1 {
		t.Fatal("offer exhaustion/reuse mismatch")
	}
	callbackHash(t, "trade-engine-pool-exhaustion", r, tradeEngineHashes["trade-engine-pool-exhaustion"])
}

func TestShopEngineOpeningRejections(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, name := range []string{"", "Grace", "李明", "abcdefghijklmnopqrstuvwx"} {
		s := tradeEngineBase()
		p := s.Callbacks.Shop
		p.Peer = true
		p.Engine.Names[1] = name
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTradeStart}, {Op: legacy.PortTestTradeStart, Value: 1}, {Op: legacy.PortTestTradeStart, Value: 2}}
		specs = append(specs, s)
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		q := v.Callbacks.Shop.Sequence
		if q[0].Return == 0 || q[1].Return != 0 || q[2].Return != 0 {
			t.Fatalf("case %d trade rejection", i)
		}
	}
	callbackHash(t, "trade-engine-opening-rejections", r, tradeEngineHashes["trade-engine-opening-rejections"])
}
