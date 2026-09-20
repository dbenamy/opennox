//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageServerTradeTransactions(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, subtype := range []byte{14, 16, 17, 18, 22, 23, 24, 25, 26, 28, 30} {
		for _, found := range []bool{false, true} {
			for _, gold := range []uint32{0, 1000} {
				for _, count := range []byte{0, 1, 3, 255} {
					makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
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
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestTradeMessage, Side: 1, ServerMessage: &legacy.PortTestServerTradeSpec{Dispatch: dispatch, Subtype: subtype, Code: code, Count: count}})
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
		if !got[i].Intact || !want[i].Intact || !got[i].Callbacks.Intact || !want[i].Callbacks.Intact || !got[i].Callbacks.Shop.Intact || !want[i].Callbacks.Shop.Intact {
			t.Fatalf("trade fixture guard: case %d", i)
		}
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("trade message differs from qualified transaction: case %d", i)
		}
	}
	interactionCapture(t, "game-server-trade-transactions", got)
}
