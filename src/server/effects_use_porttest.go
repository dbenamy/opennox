//go:build porttest

package server

import (
	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"
)

// PortTestEffectsUseEnvironment supplies actual balance arrays and optional
// projectile definitions to the unchanged production functions.
func (s *Server) PortTestEffectsUseEnvironment(values map[string][]float64, projectiles, disabled bool, speed float32) (ids [2]uint16, restore func()) {
	oldBalance := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: oldBalance}
	for k, v := range values {
		overlay.Global[strings.ToLower(k)] = balance.Array(append([]float64(nil), v...))
	}
	s.Balance.file = overlay
	oldTypes := s.Types
	var frees []func()
	s.Types.byID = make(map[string]*ObjectType, len(oldTypes.byID)+3)
	for k, v := range oldTypes.byID {
		s.Types.byID[k] = v
	}
	s.Types.byID["forcewand"] = oldTypes.byID["staffwooden"]
	if projectiles {
		s.Types.byInd = append([]*ObjectType(nil), oldTypes.byInd...)
		for i, name := range []string{"effectsbolt", "spark"} {
			id := uint16(len(s.Types.byInd))
			ids[i] = id
			t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassSimple, allowed: true, Mass: 1, Speed: speed, SpeedBase: speed}
			if i == 1 {
				ud, free := alloc.Make([]byte{}, 80)
				frees = append(frees, free)
				t.UpdateData = unsafe.Pointer(&ud[0])
				t.UpdateDataSize = 80
				for j := 64; j < 80; j++ {
					ud[j] = 0xa5
				}
				cd, free := alloc.Make([]byte{}, 20)
				frees = append(frees, free)
				t.CollideData = unsafe.Pointer(&cd[0])
				t.CollideDataSize = 20
				for j := 4; j < 20; j++ {
					cd[j] = 0x5a
				}
			}
			if disabled {
				s.Types.byInd = append(s.Types.byInd, nil)
			} else {
				s.Types.byInd = append(s.Types.byInd, t)
				s.Types.byID[name] = t
			}
		}
	}
	return ids, func() {
		s.Balance.file = oldBalance
		s.Types = oldTypes
		for _, f := range frees {
			f()
		}
	}
}
