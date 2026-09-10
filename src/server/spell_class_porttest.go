//go:build porttest

package server

import (
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
)

// PortTestSpellClassDef supplies flags directly to the existing Spells.Flags
// lookup. Valid is deliberately independent: the C eligibility helper reads
// flags only and accepts an invalid definition when its class flags match.
type PortTestSpellClassDef struct {
	Index uint32
	Flags uint32
	Valid bool
}

func PortTestSpellClassServer(defs []PortTestSpellClassDef) *Server {
	s := new(Server)
	s.Spells.byID = make(map[spell.ID]*SpellDef, len(defs))
	for _, d := range defs {
		ind := spell.ID(int32(d.Index))
		if ind <= 0 {
			continue // production DefByInd treats nonpositive spell IDs as missing.
		}
		s.Spells.byID[ind] = &SpellDef{
			ID:    ind,
			Valid: d.Valid,
			Def: things.Spell{
				Flags: things.SpellFlags(d.Flags),
			},
		}
	}
	return s
}
