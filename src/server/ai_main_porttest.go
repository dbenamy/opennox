//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/object"
)

// PortTestMainTypeIDs identifies the two synthetic cloud types used by the
// main-AI dangerous-unit predicate. They extend the state fixture's ten-type
// fresh-server registry and therefore deliberately have stable indices 10/11.
type PortTestMainTypeIDs struct {
	ToxicCloud, SmallToxicCloud int
}

// PortTestMainTypes extends PortTestMonsterStateTypes with the exact two names
// lazily resolved by nox_xxx_unitIsDangerous_547120. The types are simple: the
// fixture supplies each object's class bits separately to exercise both C
// branches. Cleanup restores the complete prior type registry.
func (s *Server) PortTestMainTypes() (PortTestMainTypeIDs, func()) {
	if len(s.Types.byInd) != 10 || s.Types.byInd[8] == nil || s.Types.byInd[9] == nil {
		panic("PortTestMainTypes requires PortTestMonsterStateTypes on a fresh server")
	}
	old := s.Types
	byInd := append([]*ObjectType(nil), old.byInd...)
	byID := make(map[string]*ObjectType, len(old.byID)+2)
	for k, v := range old.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 2)...)
	s.Types.byID = byID
	ids := PortTestMainTypeIDs{ToxicCloud: 10, SmallToxicCloud: 11}
	for _, v := range []struct {
		id  string
		ind int
	}{
		{"ToxicCloud", ids.ToxicCloud},
		{"SmallToxicCloud", ids.SmallToxicCloud},
	} {
		t := &ObjectType{
			s: &s.Types, ind: uint16(v.ind), ind2: uint16(v.ind), id: strings.ToLower(v.id),
			class: object.ClassSimple, allowed: true,
		}
		s.Types.byInd[v.ind] = t
		s.Types.byID[t.id] = t
	}
	return ids, func() { s.Types = old }
}

// PortTestMainHost returns a setter for the real server host-unit slot and a
// restore function. It changes only serverPlayers' two host pointers: callers
// provide an existing C-owned fixture Object, so no Player/Object memory is
// changed after their read-only snapshot. Passing nil models an absent host.
func (s *Server) PortTestMainHost() (set func(*Object), restore func()) {
	oldPlayer, oldUnit := s.Players.hostPlayer, s.Players.hostUnit
	set = func(unit *Object) {
		// nox_getHostPlayerUnit reads hostUnit. A Player is intentionally not
		// synthesized or edited here; the main-AI cursor branch only needs unit.
		s.Players.SetHost(nil, unit)
	}
	restore = func() { s.Players.SetHost(oldPlayer, oldUnit) }
	return set, restore
}
