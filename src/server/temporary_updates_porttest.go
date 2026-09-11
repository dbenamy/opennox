//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"
)

// PortTestTemporaryTypes installs real allocator definitions for the creations
// observed by temporary-object updates. All buffers have the effects guard layout.
func (s *Server) PortTestTemporaryTypes(missing []string) func() {
	old := s.Types
	s.Types.byInd = append([]*ObjectType(nil), old.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+5)
	for k, v := range old.byID {
		s.Types.byID[k] = v
	}
	var frees []func()
	for _, name := range []string{"smallflame", "mediumflame", "meteor", "meteorexplode", "smallspider", "temporaryscorch0", "temporaryscorch1", "temporaryscorch2"} {
		id := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassSimple, allowed: true, Mass: 1}
		ud, free := alloc.Make([]byte{}, 80)
		frees = append(frees, free)
		for i := 64; i < 80; i++ {
			ud[i] = 0xa5
		}
		t.UpdateData = unsafe.Pointer(&ud[0])
		t.UpdateDataSize = 80
		cd, free := alloc.Make([]byte{}, 20)
		frees = append(frees, free)
		for i := 4; i < 20; i++ {
			cd[i] = 0x5a
		}
		t.CollideData = unsafe.Pointer(&cd[0])
		t.CollideDataSize = 20
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[name] = t
	}
	for _, name := range missing {
		delete(s.Types.byID, strings.ToLower(name))
	}
	return func() {
		s.Types = old
		for _, f := range frees {
			f()
		}
	}
}
