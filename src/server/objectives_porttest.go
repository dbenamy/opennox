//go:build porttest

package server

import (
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// Stable low 16 address bits are needed for legacy short returns of team pointers.
// The actual team records and linked membership lists are still used.
func (s *Server) PortTestObjectiveTypes(ballUpdate, ballCollide unsafe.Pointer, missing []string) func() {
	oldTypes, oldTeams := s.Types, s.Teams
	storage, freeTeams := alloc.Make([]byte{}, 65536+17*int(unsafe.Sizeof(Team{})))
	base := uintptr(unsafe.Pointer(&storage[0]))
	aligned := (base + 65535) &^ uintptr(65535)
	teams := unsafe.Slice((*Team)(unsafe.Pointer(aligned)), 17)
	copy(teams, oldTeams.Arr)
	s.Teams.Arr = teams
	s.Teams.sm = s.sm
	s.Teams.pr = console.NewMultiPrinter()
	s.Types.byInd = append([]*ObjectType(nil), oldTypes.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(oldTypes.byID)+3)
	for k, v := range oldTypes.byID {
		s.Types.byID[k] = v
	}
	var frees []func()
	for _, name := range []string{"gameball", "gameballstart", "objectiveflag"} {
		id := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassSimple, allowed: true, Mass: 1}
		if name == "gameball" {
			t.class = object.ClassMissile
			t.Update = ballUpdate
			t.Collide = ballCollide
			ud, f := alloc.Make([]byte{}, 80)
			frees = append(frees, f)
			for i := 64; i < 80; i++ {
				ud[i] = 0xa5
			}
			t.UpdateData = unsafe.Pointer(&ud[0])
			t.UpdateDataSize = 80
			cd, f := alloc.Make([]byte{}, 20)
			frees = append(frees, f)
			for i := 4; i < 20; i++ {
				cd[i] = 0x5a
			}
			t.CollideData = unsafe.Pointer(&cd[0])
			t.CollideDataSize = 20
		}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[name] = t
	}
	for _, name := range missing {
		delete(s.Types.byID, name)
	}
	return func() {
		s.Types, s.Teams = oldTypes, oldTeams
		for _, f := range frees {
			f()
		}
		freeTeams()
	}
}
