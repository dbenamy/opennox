//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Add real registered transfer callbacks and owned type templates to a fresh
// fixture server. Its existing object-pool owner remains responsible for objects.
func (s *Server) PortTestItemXferTypes() func() {
	oldInd := s.Types.byInd
	oldID := make(map[string]*ObjectType, len(s.Types.byID))
	for k, v := range s.Types.byID {
		oldID[k] = v
	}
	var owned []unsafe.Pointer
	allocate := func(n uintptr) unsafe.Pointer {
		if n == 0 {
			return nil
		}
		p, _ := alloc.Malloc(n)
		owned = append(owned, p)
		return p
	}
	for _, sp := range []struct {
		name              string
		class             object.Class
		init, use, update uintptr
	}{
		{"SpellReward", object.ClassSimple, 0, 4, 0},
		{"AbilityReward", object.ClassSimple, 0, 4, 0},
		{"FieldGuide", object.ClassSimple, 0, 64, 0},
		{"Weapon", object.ClassWeapon, 20, 116, 8},
		{"Armor", object.ClassArmor, 20, 0, 8},
		{"Ammo", object.ClassWeapon, 20, 4, 0},
		{"Team", object.ClassFlag, 20, 0, 8},
		{"Gold", object.ClassSimple, 4, 0, 0},
		{"Obelisk", object.ClassImmobile, 0, 0, 4},
		{"ToxicCloud", object.ClassSimple, 0, 0, 4},
		{"MonsterGenerator", object.ClassImmobile, 0, 0, 96},
		{"RewardMarker", object.ClassImmobile, 220, 0, 0},
	} {
		id := uint16(len(s.Types.byInd))
		typ := &ObjectType{s: &s.Types, ind: id, ind2: id, id: strings.ToLower("Port" + sp.name), class: sp.class, allowed: true, Mass: 1}
		typ.Xfer = xferFuncs[sp.name+"Xfer"]
		if typ.Xfer == nil {
			panic("missing real item xfer: " + sp.name)
		}
		typ.InitDataSize = sp.init
		typ.InitData = allocate(sp.init)
		typ.UseDataSize = sp.use
		if sp.use != 0 {
			typ.UseData.SetPtr(allocate(sp.use))
		}
		typ.UpdateDataSize = sp.update
		typ.UpdateData = allocate(sp.update)
		s.Types.byInd = append(s.Types.byInd, typ)
		s.Types.byID[typ.id] = typ
	}
	return func() {
		s.Types.byInd = oldInd
		s.Types.byID = oldID
		for _, p := range owned {
			alloc.FreePtr(p)
		}
	}
}

// Own the actual player bitset used by the quest item-health policy.
func (s *Server) PortTestItemXferQuestPlayers(mask uint32) func() {
	old := s.Players.playersXxx
	s.Players.playersXxx = mask
	return func() { s.Players.playersXxx = old }
}
