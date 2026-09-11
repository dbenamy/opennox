//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/object"
)

// PortTestMonsterStateTypeIDs extends the lifecycle synthetic registry with
// exactly the two names lazily resolved by 534840/534A10. Zombie IDs come from
// the installed lifecycle types and are returned so a fixture never assumes
// numeric object-type layout outside this fresh-server setup.
type PortTestMonsterStateTypeIDs struct{ Mimic, Plant, Zombie, VileZombie int }

func (s *Server) PortTestMonsterStateTypes() (PortTestMonsterStateTypeIDs, func()) {
	if len(s.Types.byInd) != 8 || s.Types.byInd[2] == nil || s.Types.byInd[3] == nil {
		panic("PortTestMonsterStateTypes requires PortTestLifecycleTypes on a fresh server")
	}
	old := s.Types
	s.Types.playerAnimFrames = make([][2]int, 64)
	for i := range s.Types.playerAnimFrames {
		s.Types.playerAnimFrames[i] = [2]int{3 + i%7, i % 3}
	}
	byInd := append([]*ObjectType(nil), old.byInd...)
	byID := make(map[string]*ObjectType, len(old.byID)+2)
	for k, v := range old.byID {
		byID[k] = v
	}
	s.Types.byInd, s.Types.byID = append(byInd, make([]*ObjectType, 2)...), byID
	ids := PortTestMonsterStateTypeIDs{Mimic: 8, Plant: 9, Zombie: 2, VileZombie: 3}
	for _, v := range []struct {
		id  string
		ind int
	}{{"Mimic", ids.Mimic}, {"CarnivorousPlant", ids.Plant}} {
		t := &ObjectType{s: &s.Types, ind: uint16(v.ind), ind2: uint16(v.ind), id: strings.ToLower(v.id), class: object.ClassSimple, allowed: true}
		s.Types.byInd[v.ind], s.Types.byID[t.id] = t, t
	}
	return ids, func() { s.Types = old }
}
