//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/object"
)

// PortTestPenaltyEnvironment installs the three gem names consumed by
// sub_54D080 and a deliberately small armor lookup table consumed by
// sub_54CD30. It requires the callback fixture's 23-type server state.
//
// configure first resets this fixture's prior case. armorBits maps an item
// TypeInd to the raw armor bit returned by Sub_415D10. It does not require
// object definitions or asset initialization; both original C functions only
// use TypeInd/name lookup and the armor bit.
func (s *Server) PortTestPenaltyEnvironment() (gemIDs map[string]uint16, configure func(armorBits map[uint16]uint32), restore func()) {
	if len(s.Types.byInd) != 23 {
		panic("PortTestPenaltyEnvironment requires callback fixture types")
	}
	oldTypes := s.Types
	oldArmor := s.Armor

	byInd := append([]*ObjectType(nil), oldTypes.byInd...)
	byID := make(map[string]*ObjectType, len(oldTypes.byID)+3)
	for k, v := range oldTypes.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 3)...)
	s.Types.byID = byID

	gemIDs = map[string]uint16{
		"Diamond": 23,
		"Emerald": 24,
		"Ruby":    25,
	}
	gems := make(map[string]*ObjectType, len(gemIDs))
	for name, ind := range gemIDs {
		id := strings.ToLower(name)
		t := &ObjectType{
			s: &s.Types, ind: ind, ind2: ind, id: id,
			class: object.ClassSimple, allowed: true,
		}
		s.Types.byInd[ind] = t
		s.Types.byID[id] = t
		gems[id] = t
	}

	reset := func() {
		for id, t := range gems {
			s.Types.byInd[t.ind] = t
			s.Types.byID[id] = t
		}
		s.Armor = oldArmor
		s.Armor.table = [PlayerArmorCnt]armorRecord{}
		s.Armor.ready = true
	}
	configure = func(armorBits map[uint16]uint32) {
		reset()
		i := 0
		for typ, bit := range armorBits {
			if i == len(s.Armor.table) {
				panic("too many porttest armor entries")
			}
			s.Armor.table[i] = armorRecord{TypeInd: int(typ), Bit: bit}
			i++
		}
	}
	return gemIDs, configure, func() {
		s.Types = oldTypes
		s.Armor = oldArmor
	}
}
