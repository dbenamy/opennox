//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerTradeTransactions(t *testing.T) {
	gameMessageServerTradeTransactions(t, false)
}

func TestGameMessageServerTradeWithoutSession(t *testing.T) {
	gameMessageServerTradeTransactions(t, true)
}

func gameMessageServerTradeTransactions(t *testing.T, noSession bool) {
	t.Helper()
	var direct, messages []legacy.PortTestRoamSpec
	for _, subtype := range []byte{14, 16, 17, 18, 22, 23, 24, 25, 26, 28, 30} {
		for _, found := range []bool{false, true} {
			for _, gold := range []uint32{0, 1000} {
				for _, count := range []byte{0, 1, 3, 255} {
					makeCase := func(dispatch, noSession bool) legacy.PortTestRoamSpec {
						s := tradeEngineBase()
						p := s.Callbacks.Shop
						p.Gold = [2]uint32{gold, gold}
						p.Engine.Cache[10], p.Engine.Cache[11] = 14, 14
						if subtype == 14 || subtype == 16 || subtype == 17 {
							p.Peer = true
							p.Sequence[0].Value = 0
						}
						for i := 0; i < 3; i++ {
							code := uint32(7 + i)
							p.Items = append(p.Items, legacy.PortTestShopItem{NetCode: &code, Type: 15, Class: 8, Worth: 100, Health: true, HP: 25, MaxHP: 100, Pickup: true})
							switch subtype {
							case 22, 23:
								p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopAdd, Item: i})
							case 24, 25, 26, 28, 30:
								p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopInventory, Item: i, Side: 1})
							case 14, 16:
								p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestShopOffer, Item: i, Side: 1, Value: 100})
							}
						}
						p.Stock = []legacy.PortTestShopStock{{Type: 15, Count: 3}}
						code := uint16(999)
						if found {
							code = 7
							if subtype == 23 || subtype == 25 {
								code = 15
							}
						}
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeMessage, Side: 1, ServerMessage: &legacy.PortTestServerTradeSpec{Dispatch: dispatch, NoSession: noSession, Subtype: subtype, Code: code, Count: count}})
						return s
					}
					direct = append(direct, makeCase(false, noSession))
					messages = append(messages, makeCase(true, noSession))
				}
			}
		}
	}
	want, got := legacy.PortTestRoam(direct), legacy.PortTestRoam(messages)
	for i := range got {
		if !got[i].Intact || !want[i].Intact || !got[i].Callbacks.Intact || !want[i].Callbacks.Intact || !got[i].Callbacks.Shop.Intact || !want[i].Callbacks.Shop.Intact {
			t.Fatalf("trade fixture guard: case %d", i)
		}
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("trade message differs from qualified transaction: case %d", i)
		}
	}
	name := "game-server-trade-transactions"
	if noSession {
		name = "game-server-trade-without-session"
	}
	interactionCapture(t, name, got)
}

func TestGameMessageServerTradeOffers(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, subtype := range []byte{15, 16} {
		for _, code := range []uint16{7, 0x8017, 999} {
			for _, noSession := range []bool{false, true} {
				for _, class := range []uint32{8, 0x1000, 0x1000000} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
						s := tradeEngineBase()
						p := s.Callbacks.Shop
						p.Peer, p.Sequence[0].Value = true, 0
						netCode := uint32(7)
						p.Items = []legacy.PortTestShopItem{{NetCode: &netCode, Type: 15, Class: class, Worth: 100}}
						op := legacy.PortTestShopInventory
						if subtype == 16 {
							op = legacy.PortTestShopOffer
						}
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: op, Side: 1, Value: 100}, legacy.PortTestShopAction{Op: legacy.PortTestTradeMessage, Side: 1, ServerMessage: &legacy.PortTestServerTradeSpec{Dispatch: dispatch, NoSession: noSession, Found: code != 999, Subtype: subtype, Code: code}})
						return s
					}
					direct = append(direct, makeCase(false))
					messages = append(messages, makeCase(true))
				}
			}
		}
	}
	want, got := legacy.PortTestRoam(direct), legacy.PortTestRoam(messages)
	for i := range got {
		if !got[i].Intact || !want[i].Intact || !got[i].Callbacks.Intact || !want[i].Callbacks.Intact {
			t.Fatalf("offer fixture guard: case %d", i)
		}
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("trade offer route differs from qualified operation: case %d", i)
		}
	}
	interactionCapture(t, "game-server-trade-offers", got)
}

func TestGameMessageServerTradeOpening(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, code := range []uint16{7, 0x8017, 999} {
		for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
			for _, saving := range []bool{false, true} {
				for _, vendor := range []bool{false, true} {
					for _, game := range []uint32{0, 2048, 4096} {
						makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
							s := tradeEngineBase()
							p := s.Callbacks.Shop
							s.Lifecycle.GameFlags = game
							p.Balance["ShopAnkhCutoffStage"] = 0
							p.Load = &legacy.PortTestShopLoadSpec{}
							p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestTradeOpeningMessage, ServerMessage: &legacy.PortTestServerTradeSpec{Dispatch: dispatch, Found: code != 999, Flags: flags, Saving: saving, Vendor: vendor, Code: code}}}
							return s
						}
						direct = append(direct, makeCase(false))
						messages = append(messages, makeCase(true))
					}
				}
			}
		}
	}
	want, got := legacy.PortTestRoam(direct), legacy.PortTestRoam(messages)
	for i := range got {
		if !got[i].Intact || !want[i].Intact || !got[i].Callbacks.Intact || !want[i].Callbacks.Intact {
			t.Fatalf("trade opening fixture guard: case %d", i)
		}
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("trade opening differs from qualified operation: case %d", i)
		}
	}
	interactionCapture(t, "game-server-trade-opening", got)
}
