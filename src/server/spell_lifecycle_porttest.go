//go:build porttest

package server

import (
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/sound"
)

// PortTestSpellLifecycleDef supplies real spell-definition fields, including
// costs and sounds that the earlier class-only fixture did not need.
type PortTestSpellLifecycleDef struct {
	Index          int
	Flags          uint32
	Valid, Enabled bool
	ManaCost       int
	Phonemes       []int
	Sounds         [3]int
}

func (s *Server) PortTestSpellLifecycle(defs []PortTestSpellLifecycleDef, tree *PhonemeLeaf) func() {
	oldDefs, oldTree, oldDur := s.Spells.byID, s.Spells.tree, s.Spells.Dur
	s.Spells.byID = make(map[spell.ID]*SpellDef, len(defs))
	s.Spells.tree = tree
	for _, d := range defs {
		if d.Index <= 0 {
			continue
		}
		var phon []spell.Phoneme
		for _, v := range d.Phonemes {
			phon = append(phon, spell.Phoneme(v))
		}
		id := spell.ID(d.Index)
		s.Spells.byID[id] = &SpellDef{ID: id, Valid: d.Valid, Enabled: d.Enabled, Def: things.Spell{Flags: things.SpellFlags(d.Flags), ManaCost: d.ManaCost, Phonemes: phon}, CastSound: sound.ID(d.Sounds[0]), OnSound: sound.ID(d.Sounds[1]), OffSound: sound.ID(d.Sounds[2])}
	}
	s.Spells.Dur = SpellsDuration{}
	s.Spells.Dur.init(s)
	s.Spells.Dur.Init()
	return func() { s.Spells.Dur.Free(); s.Spells.Dur = oldDur; s.Spells.tree = oldTree; s.Spells.byID = oldDefs }
}
