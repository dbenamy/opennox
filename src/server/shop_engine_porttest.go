//go:build porttest

package server

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/opennox/libs/strman"
)

var portTestTradeStrings struct {
	once sync.Once
	sm   *strman.StringManager
}

// Install actual localized strings so rejection messages exercise formatting,
// not the StringManager's missing-entry fallback. The immutable table is shared.
func (s *Server) PortTestTradeStrings() func() {
	portTestTradeStrings.once.Do(func() {
		entries := []strman.Entry{}
		for _, v := range [][2]string{
			{"NPC:PortTradeVendor", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"},
			{"Trade.c:StarterAlreadyTrading", "STARTER_ALREADY_TRADING"},
			{"Trade.c:OtherAlreadyTrading", "OTHER_ALREADY_TRADING:%s"},
			{"Trade.c:TradeMaxObjectsReached", "TRADE_POOL_FULL"},
			{"Trade.c:CantSellQuestItem", "CANNOT_SELL_QUEST"},
			{"Trade.c:CantSellItem", "CANNOT_SELL_ITEM"},
			{"pickup.c:MaxSameItem", "TOO_MANY_ITEMS"},
			{"pickup.c:GoldPickup", "GOLD:%d"},
		} {
			entries = append(entries, strman.Entry{ID: strman.ID(v[0]), Vals: []strman.Variant{{Str: v[1]}}})
		}
		f, err := os.CreateTemp("", "opennox-trade-strings-*.json")
		if err != nil {
			panic(err)
		}
		path := f.Name()
		defer os.Remove(path)
		err = json.NewEncoder(f).Encode(struct {
			Entries []strman.Entry `json:"entries"`
		}{entries})
		closeErr := f.Close()
		if err != nil {
			panic(err)
		}
		if closeErr != nil {
			panic(closeErr)
		}
		sm := strman.New()
		if err := sm.ReadJSON(path); err != nil {
			panic(err)
		}
		portTestTradeStrings.sm = sm
	})
	old := s.sm
	s.sm = portTestTradeStrings.sm
	return func() { s.sm = old }
}
