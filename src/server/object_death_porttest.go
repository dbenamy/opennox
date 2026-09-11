//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"strings"
)

func (s *Server) PortTestDeathTypes() func() {
	old := s.Types
	s.Types.byInd = append([]*ObjectType(nil), old.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+3)
	for k, v := range old.byID {
		s.Types.byID[k] = v
	}
	for _, name := range []string{"BarrelBreaking", "BarrelPortTest", "GameBallStart"} {
		id := strings.ToLower(name)
		ind := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: ind, ind2: ind, id: id, class: object.ClassSimple, allowed: true, Mass: 1}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[id] = t
	}
	return func() { s.Types = old }
}
