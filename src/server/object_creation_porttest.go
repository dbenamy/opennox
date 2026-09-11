//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestCreationEnvironment extends the callback fixture with the lookup
// state used by 54C0C0..54CBB0. It requires both PortTestAICallbackTypes and
// PortTestAICallbackModifiers. IDs are the real uint16 TypeInd values returned
// by the name lookup used by the unchanged C code.
//
// configure restores this environment's previous case before applying its
// arguments. A true entry in disabled removes that synthetic type from both
// lookup tables. weaponType and armorType are Modifier.TypeInd values, not
// indices into either modifier list. effectsEnabled controls the three extra
// ModifierEff descriptors used by the Oblivion weapon cases.
func (s *Server) PortTestCreationEnvironment() (ids map[string]uint16, configure func(disabled map[string]bool, weaponType, armorType, durability uint32, balanceVals map[string]float64, effectsEnabled bool), restore func()) {
	if len(s.Types.byInd) != 23 || s.Types.ByID("Sword") == nil || s.Types.ByID("ToxicCloud") == nil {
		panic("PortTestCreationEnvironment requires PortTestAICallbackTypes and PortTestAICallbackModifiers")
	}

	oldTypes := s.Types
	byInd := append([]*ObjectType(nil), oldTypes.byInd...)
	byID := make(map[string]*ObjectType, len(oldTypes.byID)+12)
	for k, v := range oldTypes.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 12)...)
	s.Types.byID = byID

	// The existing main-fixture ToxicCloud is simple-only. C 54CB10 creates it
	// and writes the lifetime through Object.UpdateData (+748), so its type must
	// supply a real four-byte source buffer for NewObject to clone.
	oldCloud := oldTypes.byInd[10]
	cloud := *oldCloud
	cloudData, freeCloudData := alloc.New(uint32(0))
	cloud.UpdateData, cloud.UpdateDataSize = unsafe.Pointer(cloudData), 4
	s.Types.byInd[10], s.Types.byID[cloud.id] = &cloud, &cloud

	ids = make(map[string]uint16, 13)
	ids["ToxicCloud"] = 10
	names := []string{
		"UrchinShaman", "Wizard", "WizardWhite", "Beholder", "Lich", "LichLord", "Demon", "WizardGreen", "WillOWisp",
		"OblivionHeart", "OblivionWierdling", "OblivionOrb",
	}
	created := make(map[string]*ObjectType, len(names)+1)
	created[cloud.id] = &cloud
	for i, id := range names {
		name := strings.ToLower(id)
		ind := uint16(23 + i)
		t := &ObjectType{s: &s.Types, ind: ind, ind2: ind, id: name, class: object.ClassSimple, allowed: true, Mass: 1}
		s.Types.byInd[ind], s.Types.byID[name], created[name] = t, t, t
		ids[id] = ind
	}

	// Keep the callback fixture's modifier descriptors/lists intact as the base.
	// The three effects are appended as the third effect list, avoiding mutation
	// of its WeaponPower/Material chains.
	oldModif := s.Modif
	mkEff := func(name string, ind uint32) (*ModifierEff, func()) {
		eff, freeEff := alloc.New(ModifierEff{})
		p, freeName := alloc.CString(name)
		eff.name0, eff.ind4 = p, ind
		return eff, func() { freeName(); freeEff() }
	}
	lightning4, freeLightning4 := mkEff("Lightning4", 4)
	vampirism2, freeVampirism2 := mkEff("Vampirism2", 5)
	lightning3, freeLightning3 := mkEff("Lightning3", 6)
	lightning4.next136, vampirism2.prev140 = vampirism2, lightning4
	vampirism2.next136, lightning3.prev140 = lightning3, vampirism2
	weapon, freeWeapon := alloc.New(Modifier{})
	armor, freeArmor := alloc.New(Modifier{})

	oldBalance := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: oldBalance}
	s.Balance.file = overlay

	restoreTypes := func() {
		for name, t := range created {
			s.Types.byInd[int(t.ind)] = t
			s.Types.byID[name] = t
		}
		s.Types.byInd[10], s.Types.byID[cloud.id] = &cloud, &cloud
	}
	configure = func(disabled map[string]bool, weaponType, armorType, durability uint32, balanceVals map[string]float64, effectsEnabled bool) {
		restoreTypes()
		for name, off := range disabled {
			if !off {
				continue
			}
			name = strings.ToLower(name)
			if t := created[name]; t != nil {
				delete(s.Types.byID, name)
				s.Types.byInd[int(t.ind)] = nil
			}
		}

		m := oldModif
		weapon.TypeInd, weapon.Durability52, weapon.Next80, weapon.Prev84 = weaponType, durability, oldModif.Dword_5d4594_251600, nil
		armor.TypeInd, armor.Durability52, armor.Next80, armor.Prev84 = armorType, durability, oldModif.Dword_5d4594_251608, nil
		m.Dword_5d4594_251600, m.Dword_5d4594_251608 = weapon, armor
		if effectsEnabled {
			m.types[2], m.cnt = lightning4, 7
		} else {
			m.types[2] = nil
		}
		s.Modif = m

		overlay.Global = balance.Config{}
		for k, v := range balanceVals {
			overlay.Global[strings.ToLower(k)] = balance.Array{v}
		}
	}

	return ids, configure, func() {
		s.Types = oldTypes
		s.Modif = oldModif
		s.Balance.file = oldBalance
		freeArmor()
		freeWeapon()
		freeLightning3()
		freeVampirism2()
		freeLightning4()
		freeCloudData()
	}
}
