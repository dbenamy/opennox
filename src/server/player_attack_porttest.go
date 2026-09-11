//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"
)

func (s *Server) PortTestAttackTypes(speed float32, missing []string) func() {
	old := s.Types
	s.Types.byInd = append([]*ObjectType(nil), old.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+5)
	for k, v := range old.byID {
		s.Types.byID[k] = v
	}
	var frees []func()
	region := func(payload int) unsafe.Pointer {
		b, f := alloc.Make([]byte{}, payload+16)
		frees = append(frees, f)
		for i := payload; i < len(b); i++ {
			b[i] = 0x5a
		}
		return unsafe.Pointer(&b[0])
	}
	for _, name := range []string{"archerarrow", "archerbolt", "weakarcherarrow", "fanchakraminmotion", "roundchakraminmotion"} {
		id := uint16(len(s.Types.byInd))
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassMissile, allowed: true, Mass: 1, Speed: speed, SpeedBase: speed}
		t.InitData = region(20)
		t.InitDataSize = 36
		t.CollideData = region(8)
		t.CollideDataSize = 24
		ud, f := alloc.Make([]byte{}, 80)
		frees = append(frees, f)
		for i := 64; i < 80; i++ {
			ud[i] = 0xa5
		}
		t.UpdateData = unsafe.Pointer(&ud[0])
		t.UpdateDataSize = 80
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
