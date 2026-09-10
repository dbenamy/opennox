//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
)

// PortTestRuleSpell describes one deliberately small spell definition installed
// by PortTestRuleServerSetup.
type PortTestRuleSpell struct {
	ID        string
	Index     int
	Title     string
	RuleFlags uint32
	Valid     bool
}

type PortTestRuleItem struct {
	Name string
	Bit  uint32
}

// PortTestRuleServer describes the lookup data installed in a fresh server.
// The tables are copies of the shipped weapon and armor tables; only their
// object type indices are synthetic and unique.
type PortTestRuleServer struct {
	Spells  []PortTestRuleSpell
	Weapons []PortTestRuleItem
	Armors  []PortTestRuleItem
}

// PortTestRuleServerSetup creates a minimal server whose existing production
// spell, weapon, and armor lookup methods are usable by rule-loader tests.
func PortTestRuleServerSetup() (*Server, PortTestRuleServer) {
	s := new(Server)
	s.Types.byInd = []*ObjectType{nil}
	s.Types.byID = make(map[string]*ObjectType)
	addType := func(id string) int {
		ind := len(s.Types.byInd)
		t := &ObjectType{s: &s.Types, ind: uint16(ind), id: id}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[strings.ToLower(id)] = t
		return ind
	}

	s.Weapons.table = weaponTable
	for i := range s.Weapons.table {
		if s.Weapons.table[i].Name != "" && s.Weapons.table[i].Bit != 0 {
			s.Weapons.table[i].TypeInd = addType(s.Weapons.table[i].Name)
		}
	}
	s.Weapons.ready = true

	s.Armor.table = armorTable
	for i := range s.Armor.table {
		if s.Armor.table[i].Name != "" && s.Armor.table[i].Bit != 0 {
			s.Armor.table[i].TypeInd = addType(s.Armor.table[i].Name)
		}
	}
	s.Armor.ready = true

	desc := PortTestRuleServer{
		Spells: []PortTestRuleSpell{
			{ID: spell.SPELL_FIREBALL.String(), Index: int(spell.SPELL_FIREBALL), Title: "Rule Fireball", RuleFlags: uint32(things.SpellClassAny), Valid: true},
			{ID: spell.SPELL_ANCHOR.String(), Index: int(spell.SPELL_ANCHOR), Title: "No Rule", RuleFlags: uint32(things.SpellTargeted), Valid: true},
			{ID: spell.SPELL_BLINK.String(), Index: int(spell.SPELL_BLINK), Title: "Fíre Burst", RuleFlags: uint32(things.SpellClassWizard), Valid: true},
			{ID: spell.SPELL_BURN.String(), Index: int(spell.SPELL_BURN), Title: "Invalid", RuleFlags: uint32(things.SpellClassAny), Valid: false},
		},
	}
	for _, it := range s.Weapons.table {
		if it.Name != "" && it.Bit != 0 {
			desc.Weapons = append(desc.Weapons, PortTestRuleItem{Name: it.Name, Bit: it.Bit})
		}
	}
	for _, it := range s.Armor.table {
		if it.Name != "" && it.Bit != 0 {
			desc.Armors = append(desc.Armors, PortTestRuleItem{Name: it.Name, Bit: it.Bit})
		}
	}
	s.Spells.byID = make(map[spell.ID]*SpellDef, len(desc.Spells))
	for _, it := range desc.Spells {
		ind := spell.ID(it.Index)
		s.Spells.byID[ind] = &SpellDef{
			ID:    ind,
			Valid: it.Valid,
			Title: it.Title,
			Def: things.Spell{
				ID:    it.ID,
				Flags: things.SpellFlags(it.RuleFlags),
			},
		}
	}
	return s, desc
}
