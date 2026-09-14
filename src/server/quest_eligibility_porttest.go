//go:build porttest

package server

// Reuse the guarded type owner, then supply explicit equipment bits for the
// named eligibility subjects. Lookup itself still uses production tables.
func (s *Server) PortTestEligibilityTypes(names []string, armorBit, weaponBit uint32) func() {
	restore := s.PortTestRewardTypes(names, nil, true, armorBit, weaponBit)
	s.Armor.table = [PlayerArmorCnt]armorRecord{}
	s.Weapons.table = [PlayerWeaponCnt]weaponRecord{}
	for i, n := range []string{"EligibilityArmor", "StreetSneakers", "StreetPants", "StreetShirt"} {
		s.Armor.table[i] = armorRecord{TypeInd: int(s.Types.IndByID(n)), Bit: armorBit}
	}
	for i, n := range []string{"EligibilityWeapon", "SulphorousFlareWand", "InfinitePainWand"} {
		s.Weapons.table[i] = weaponRecord{TypeInd: int(s.Types.IndByID(n)), Bit: weaponBit}
	}
	return restore
}
