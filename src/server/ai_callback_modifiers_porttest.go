//go:build porttest

package server

import (
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestAICallbackModifiers installs the three exact name lookups used by
// the generated-gear callbacks and makes their synthetic item types exercise
// the real InitData (+692) and UseData (+736) copies. It requires the type
// table from PortTestAICallbackTypes.
//
// The returned pointers are the descriptors C writes into the first four
// InitData words. They are C-allocated so their raw 32-bit addresses are safe
// to retain in a C-created object. The fifth InitData word is deliberately not
// specified: sub_54A390 passes a 20-byte local while defining only four words.
type PortTestAICallbackModifiers struct {
	WeaponPower1 *ModifierEff
	Material1    *ModifierEff
	Material2    *ModifierEff
}

func (s *Server) PortTestAICallbackModifiers() (PortTestAICallbackModifiers, func()) {
	if len(s.Types.byInd) < 23 || s.Types.ByID("Sword") == nil || s.Types.ByID("FanChakram") == nil {
		panic("PortTestAICallbackModifiers requires PortTestAICallbackTypes")
	}

	mk := func(name string, ind uint32) (*ModifierEff, func()) {
		eff, freeEff := alloc.New(ModifierEff{})
		namep, freeName := alloc.CString(name)
		eff.name0, eff.ind4 = namep, ind
		return eff, func() {
			freeName()
			freeEff()
		}
	}
	power, freePower := mk("WeaponPower1", 1)
	material1, freeMaterial1 := mk("Material1", 2)
	material2, freeMaterial2 := mk("Material2", 3)
	material1.next136, material2.prev140 = material2, material1

	oldModif := s.Modif
	s.Modif = serverModifiers{sm: oldModif.sm, cnt: 4}
	s.Modif.types[0] = power
	s.Modif.types[1] = material1

	// Every source blob is C-owned. NewObject copies InitDataSize bytes into a
	// distinct C allocation, which is the exact storage C memcpy writes at +692.
	type savedType struct {
		t   *ObjectType
		old ObjectType
	}
	var saved []savedType
	var frees []func()
	for _, name := range []string{
		"Sword", "WoodenShield", "SteelShield", "StaffWooden",
		"Bow", "Quiver", "OgreAxe", "FanChakram",
	} {
		t := s.Types.ByID(name)
		saved = append(saved, savedType{t: t, old: *t})

		init, freeInit := alloc.Make([]byte{}, 20)
		frees = append(frees, freeInit)
		t.InitData, t.InitDataSize = unsafe.Pointer(&init[0]), 20
		t.class = object.ClassSimple
		switch name {
		case "WoodenShield", "SteelShield":
			t.class |= object.ClassArmor
		default:
			t.class |= object.ClassWeapon
		}
		if name == "FanChakram" {
			use, freeUse := alloc.Make([]byte{}, 2)
			frees = append(frees, freeUse)
			// Do not call SetPtr: this fixture owns and frees this raw source
			// allocation after restoring the type value.
			t.UseData.Ptr, t.UseDataSize = unsafe.Pointer(&use[0]), 2
			t.subclass = object.SubClass(0x82)
		}
	}

	defs := PortTestAICallbackModifiers{
		WeaponPower1: power,
		Material1:    material1,
		Material2:    material2,
	}
	return defs, func() {
		for _, v := range saved {
			*v.t = v.old
		}
		for i := len(frees) - 1; i >= 0; i-- {
			frees[i]()
		}
		s.Modif = oldModif
		freeMaterial2()
		freeMaterial1()
		freePower()
	}
}
