//go:build porttest

package server

import (
	"encoding/json"
	"os"

	"github.com/opennox/libs/strman"
)

// PortTestMeterStrings installs a real language-bearing string manager.
func (s *Server) PortTestMeterStrings() (func(int), func()) {
	old := s.sm
	f, err := os.CreateTemp("", "opennox-meter-strings-*.json")
	if err != nil {
		panic(err)
	}
	path := f.Name()
	if err := f.Close(); err != nil {
		panic(err)
	}
	configure := func(language int) {
		data, err := json.Marshal(struct {
			Lang    int            `json:"lang"`
			Entries []strman.Entry `json:"entries"`
		}{language, []strman.Entry{
			{ID: "keybind:Q", Vals: []strman.Variant{{Str: "Q"}}},
			{ID: "keybind:F12", Vals: []strman.Variant{{Str: "F12"}}},
			{ID: "keybind:KpEnter", Vals: []strman.Variant{{Str: "Keypad Enter"}}},
			{ID: "guimeter.c:ToolTipCharges", Vals: []strman.Variant{{Str: "Charges"}}},
			{ID: "guimeter.c:CurePoisonSlotTT", Vals: []strman.Variant{{Str: "Cure poison"}}},
			{ID: "guimeter.c:HealthSlotTT", Vals: []strman.Variant{{Str: "Health potion"}}},
			{ID: "guimeter.c:ManaSlotTT", Vals: []strman.Variant{{Str: "Mana potion"}}},
			{ID: "guimeter.c:ToolTipMana", Vals: []strman.Variant{{Str: "Mana"}}},
			{ID: "guimeter.c:ToolTipHealth", Vals: []strman.Variant{{Str: "Health"}}},
			{ID: "guimeter.c:ToolTipCurWeapon", Vals: []strman.Variant{{Str: "Current weapon"}}},
			{ID: "ToolTip.c:BookOf", Vals: []strman.Variant{{Str: "Book of"}}},
			{ID: "ToolTip.c:LoreScroll", Vals: []strman.Variant{{Str: "Beast scroll"}}},
			{ID: "ToolTip.c:NoArmsInfo", Vals: []strman.Variant{{Str: "Missing: %S"}}},
		}})
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			panic(err)
		}
		sm := strman.New()
		if err := sm.ReadJSON(path); err != nil {
			panic(err)
		}
		s.sm = sm
	}
	return configure, func() { s.sm = old; _ = os.Remove(path) }
}
