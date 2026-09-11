//go:build porttest

package server

import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestShopLoadTypes supplies the retained stock/reward factories with
// ordinary C-owned definitions. No reward-generation algorithm is replaced.
func (s *Server) PortTestShopLoadTypes() (configure func(bool, uint32), restore func()) {
	if len(s.Types.byInd) != 27 {
		panic("shop loader fixture type order")
	}
	old := s.Types
	s.Types.byInd = append(append([]*ObjectType(nil), old.byInd...), make([]*ObjectType, 6)...)
	s.Types.byID = make(map[string]*ObjectType, len(old.byID)+6)
	for key, t := range old.byID {
		s.Types.byID[key] = t
	}
	init, freeInit := alloc.Make([]byte{}, 216)
	use, freeUse := alloc.Make([]byte{}, 128)
	for i, name := range []string{"commonspellbook", "abilitybook", "fieldguide", "ankhtradable", "rewardmarker", "porttestshopweapon"} {
		id := uint16(27 + i)
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: name, class: object.ClassSimple, allowed: true, Worth: 101}
		if i < 3 {
			t.class, t.subclass = object.Class(0x100), object.SubClass([]uint32{1, 4, 2}[i])
			t.Xfer = xferFuncs[[]string{"SpellRewardXfer", "AbilityRewardXfer", "FieldGuideXfer"}[i]]
			t.UseData.Ptr, t.UseDataSize = unsafe.Pointer(&use[0]), 128
		}
		if i >= 4 {
			t.InitData, t.InitDataSize = unsafe.Pointer(&init[0]), 216
			if i == 5 {
				t.class, t.InitDataSize = object.ClassSimple|object.ClassWeapon, 20
			}
		}
		s.Types.byInd[id], s.Types.byID[name] = t, t
	}
	marker := s.Types.byInd[31]
	configure = func(enabled bool, chance uint32) {
		clear(init)
		binary.LittleEndian.PutUint32(init[212:], chance)
		if enabled {
			s.Types.byID[marker.id] = marker
		} else {
			delete(s.Types.byID, marker.id)
		}
	}
	return configure, func() { s.Types = old; freeUse(); freeInit() }
}
