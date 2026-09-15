//go:build porttest

package server

import (
	"fmt"
	"github.com/opennox/libs/spell"
	"unsafe"
)

// PortTestInventoryWindowSpells installs owned definitions in the real spell registry.
func (s *Server) PortTestInventoryWindowSpells(icons []unsafe.Pointer) func() {
	old := s.Spells.byID
	s.Spells.byID = make(map[spell.ID]*SpellDef)
	for i := 1; i < SpellsMax; i++ {
		id := spell.ID(i)
		s.Spells.byID[id] = &SpellDef{ID: id, Effect: id, Enabled: true, Valid: true, Title: fmt.Sprintf("Spell %d", i), Icon: icons[i%len(icons)]}
	}
	return func() { s.Spells.byID = old }
}
