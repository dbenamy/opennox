//go:build porttest

package server

import (
	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
)

// PortTestAISpellDefs replaces only the live spell-definition map. It reuses
// PortTestSpellClassDef so valid/flags semantics match the existing spell-class
// fixture. Each configure call replaces the map instead of layering entries.
func (s *Server) PortTestAISpellDefs() (configure func([]PortTestSpellClassDef), restore func()) {
	old := s.Spells.byID
	configure = func(defs []PortTestSpellClassDef) {
		byID := make(map[spell.ID]*SpellDef, len(defs))
		for _, d := range defs {
			id := spell.ID(int32(d.Index))
			if id <= 0 {
				continue
			}
			byID[id] = &SpellDef{ID: id, Valid: d.Valid, Def: things.Spell{Flags: things.SpellFlags(d.Flags)}}
		}
		s.Spells.byID = byID
	}
	return configure, func() { s.Spells.byID = old }
}

// PortTestAIInversionRange overlays only InversionRange over the current
// balance file. set mutates the one overlay value in place, so repeated cases
// do not build a parent chain. restore reinstates the exact original file.
func (s *Server) PortTestAIInversionRange() (set func(float64), restore func()) {
	old := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: old}
	s.Balance.file = overlay
	set = func(v float64) { overlay.Global["inversionrange"] = balance.Array{v} }
	return set, func() { s.Balance.file = old }
}
