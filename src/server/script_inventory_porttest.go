//go:build porttest

package server

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func (s *Server) PortTestScriptInventoryTypes(init unsafe.Pointer) func() {
	names := []string{"OblivionHalberd", "OblivionHeart", "OblivionWierdling", "OblivionOrb"}
	restore := s.PortTestRewardTypes(names, nil, true, 0, 0)
	oldModifiers := s.Modif.Dword_5d4594_251600
	mods, free := alloc.Make([]Modifier{}, len(names))
	for i, name := range names {
		t := s.Types.ByID(name)
		t.Init = init
		t.class = object.ClassWeapon
		t.subclass = object.SubClass(uint32(1) << uint(23+i))
		mods[i].TypeInd = uint32(t.ind)
		mods[i].Next80 = s.Modif.Dword_5d4594_251600
		s.Modif.Dword_5d4594_251600 = &mods[i]
		s.Weapons.table[i] = weaponRecord{TypeInd: int(t.ind), Bit: uint32(1) << uint(23+i)}
	}
	return func() { s.Modif.Dword_5d4594_251600 = oldModifiers; free(); restore() }
}
