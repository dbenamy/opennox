//go:build porttest

package server

import "github.com/opennox/libs/object"

func (s *Server) PortTestWorldTypes() func() {
	old := s.Types
	s.Types.byInd = append([]*ObjectType(nil), old.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+2)
	for k, v := range old.byID {
		s.Types.byID[k] = v
	}
	for _, name := range []string{"trigger", "pressureplate"} {
		id := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassSimple, allowed: true, Mass: 1}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[name] = t
	}
	return func() { s.Types = old }
}
