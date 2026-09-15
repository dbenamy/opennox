//go:build porttest

package server

import (
	"encoding/json"
	"os"

	"github.com/opennox/libs/strman"
)

// PortTestTooltipStrings installs a real language-bearing string manager.
func (s *Server) PortTestTooltipStrings() (func(int), func()) {
	old := s.sm
	f, err := os.CreateTemp("", "opennox-tooltip-strings-*.json")
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
			{ID: "ToolTip.c:BookOf", Vals: []strman.Variant{{Str: "BookOf"}}},
			{ID: "ToolTip.c:LoreScroll", Vals: []strman.Variant{{Str: "LoreScroll"}}},
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
