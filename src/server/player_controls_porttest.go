//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"unsafe"
)

func (s *Server) PortTestControlsTypes(init unsafe.Pointer, disallowed uint8) func() {
	names := []string{"StreetShirt", "StreetPants", "StreetSneakers", "WizardRobe", "Longsword", "WoodenShield", "Bow", "CrossBow", "FireStormWand", "DeathRayWand", "ForceWand"}
	restore := s.PortTestRewardTypes(names, nil, true, 0, 0)

	for _, name := range names {
		s.Types.ByID(name).Init = init
	}
	armor := []struct {
		name string
		bit  uint32
		flag uint8
	}{{"StreetShirt", 0x400, 1}, {"StreetPants", 4, 2}, {"StreetSneakers", 1, 4}, {"WizardRobe", 0x4000, 16}, {"WoodenShield", 0x1000000, 128}}
	for i, row := range armor {
		t := s.Types.ByID(row.name)
		t.class = object.ClassArmor
		t.allowed = disallowed&row.flag == 0
		s.Armor.table[i] = armorRecord{TypeInd: int(t.ind), Bit: row.bit}
	}
	weapons := []struct {
		name string
		bit  uint32
		flag uint8
	}{{"Longsword", 0x8000, 8}, {"Bow", 0x100, 32}, {"CrossBow", 0x200, 64}, {"FireStormWand", 0x10000, 0}, {"DeathRayWand", 0x20000, 0}, {"ForceWand", 0x40000, 0}}
	for i, row := range weapons {
		t := s.Types.ByID(row.name)
		t.class = object.ClassWeapon
		t.allowed = disallowed&row.flag == 0
		s.Weapons.table[i] = weaponRecord{TypeInd: int(t.ind), Bit: row.bit}
	}

	return restore
}
func (s *Server) PortTestControlsModifiers(mods []*ModifierEff, names []*byte) func() {
	old := s.Modif.types
	for i, m := range mods {
		m.name0 = names[i]
		m.ind4 = uint32(i + 1)
		if i > 0 {
			m.prev140 = mods[i-1]
		}
		if i+1 < len(mods) {
			m.next136 = mods[i+1]
		}
	}
	s.Modif.types = [3]*ModifierEff{mods[0], nil, nil}
	return func() { s.Modif.types = old }
}
