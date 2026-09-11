//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"unsafe"
)

// PortTestRewardTypes supplies actual allocator-backed definitions, including
// guarded storage for the full reward marker and generated book/modifier data.
// Definition selection remains in the production reward routines.
func (s *Server) PortTestRewardTypes(names []string, unavailable []string, allowed bool, armorBit, weaponBit uint32) func() {
	oldTypes, oldArmor, oldWeapons := s.Types, s.Armor, s.Weapons
	s.Types.byInd = append([]*ObjectType(nil), oldTypes.byInd...)
	s.Types.byID = make(map[string]*ObjectType, len(oldTypes.byID)+len(names))
	for k, v := range oldTypes.byID {
		s.Types.byID[k] = v
	}
	s.Armor.table = [PlayerArmorCnt]armorRecord{}
	s.Armor.ready = true
	// Preserve unrelated weapon definitions; install the two explicit test rows.
	var frees []func()
	region := func(n int) unsafe.Pointer {
		b, f := alloc.Make([]byte{}, n+16)
		frees = append(frees, f)
		for i := n; i < len(b); i++ {
			b[i] = 0x5a
		}
		return unsafe.Pointer(&b[0])
	}
	for _, name := range names {
		id := uint16(len(s.Types.byInd))
		key := strings.ToLower(name)
		t := &ObjectType{s: &s.Types, ind: id, ind2: id, id: key, class: object.ClassSimple, allowed: allowed, Mass: 1}
		t.InitData = region(256)
		t.InitDataSize = 272
		t.UpdateData = region(64)
		t.UpdateDataSize = 80
		for i := 64; i < 80; i++ {
			*(*byte)(unsafe.Add(t.UpdateData, i)) = 0xa5
		}
		t.UseData.Ptr = region(64)
		t.UseDataSize = 80
		switch key {
		case "porttestrewardarmor":
			t.class = object.ClassArmor
			s.Armor.table[0] = armorRecord{TypeInd: int(id), Bit: armorBit}
		case "porttestrewardweapon", "porttestrewardwand":
			t.class = object.ClassWeapon
			if key == "porttestrewardwand" {
				t.class = object.ClassWand
				t.subclass = 0x10000
			}
			slot := 0
			if key == "porttestrewardwand" {
				slot = 1
			}
			s.Weapons.table[slot] = weaponRecord{TypeInd: int(id), Bit: weaponBit}
		}
		s.Types.byInd = append(s.Types.byInd, t)
		s.Types.byID[key] = t
	}
	for _, name := range unavailable {
		key := strings.ToLower(name)
		if t := s.Types.byID[key]; t != nil {
			delete(s.Types.byID, key)
			s.Types.byInd[t.ind] = nil
		}
	}
	return func() {
		s.Types = oldTypes
		s.Armor = oldArmor
		s.Weapons = oldWeapons
		for _, f := range frees {
			f()
		}
	}
}

func (s *Server) PortTestRewardModifier(mod *ModifierEff, name *byte) func() {
	old := s.Modif.types
	mod.name0 = name
	mod.ind4 = 42
	s.Modif.types = [3]*ModifierEff{nil, nil, mod}
	return func() { s.Modif.types = old }
}
