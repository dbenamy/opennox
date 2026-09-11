//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"sort"
	"strings"
)

// Supply actual lookup tables and player control state for retained dependencies.
func (s *Server) PortTestInventoryEnvironment(blocked bool, weapons, armor map[uint16]uint32) func() {
	oldTypes, oldWeapons, oldArmor, oldBlocked := s.Types, s.Weapons, s.Armor, s.Players.playersXxx
	oldTeams := s.Teams
	teams, freeTeams := alloc.Make([]Team{}, 3)
	s.Teams = serverTeams{Arr: teams, ActiveCnt: 2}
	for i := 1; i < 3; i++ {
		teams[i] = Team{IDVal: TeamID(i), ind: byte(i), active: 1}
	}
	s.Types.byInd = append([]*ObjectType(nil), oldTypes.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(oldTypes.byID)+8)
	for k, v := range oldTypes.byID {
		s.Types.byID[k] = v
	}
	for _, name := range []string{"Glyph", "Torch", "Lantern", "Crown", "StreetSneakers", "WizardRobe", "WoodenShield", "SteelShield"} {
		id := strings.ToLower(name)
		if s.Types.byID[id] != nil {
			continue
		}
		ind := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: ind, ind2: ind, id: id, class: object.ClassSimple, allowed: true, Mass: 1}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[id] = t
	}
	keys := func(m map[uint16]uint32) []int {
		out := make([]int, 0, len(m))
		for k := range m {
			out = append(out, int(k))
		}
		sort.Ints(out)
		return out
	}
	s.Weapons = serverWeapons{ready: true, sm: s.Strings()}
	for i, k := range keys(weapons) {
		s.Weapons.table[i] = weaponRecord{TypeInd: k, Bit: weapons[uint16(k)]}
	}
	s.Armor = serverArmor{ready: true, sm: s.Strings()}
	for i, k := range keys(armor) {
		s.Armor.table[i] = armorRecord{TypeInd: k, Bit: armor[uint16(k)]}
	}
	s.Players.playersXxx = 0
	if blocked {
		s.Players.playersXxx = 1 << 1
	}
	return func() {
		s.Types, s.Weapons, s.Armor, s.Players.playersXxx = oldTypes, oldWeapons, oldArmor, oldBlocked
		s.Teams = oldTeams
		freeTeams()
	}
}
