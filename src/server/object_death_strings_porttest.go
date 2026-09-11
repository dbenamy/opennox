//go:build porttest

package server

import (
	"encoding/json"
	"os"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestObjectDeathStrings installs the real StringManager, modifier lookup,
// and armor lookup used by ArmorDie/WeaponDie. It expects the callback fixture
// type table, where TypeInd 15 is an allocated synthetic item type.
//
// configure resets the previous case. pluralName is the modifier's UTF-16
// description; armor selects whether that description is exposed through the
// armor or weapon modifier lookup. The legacy fixture owns the guarded item
// InitData (+692); this helper deliberately allocates no item storage.
func (s *Server) PortTestObjectDeathStrings() (configure func(language int, pluralName string, armor bool), restore func()) {
	if len(s.Types.byInd) <= 15 || s.Types.byInd[15] == nil {
		panic("PortTestObjectDeathStrings requires callback fixture type 15")
	}
	oldSM, oldModif, oldArmor := s.sm, s.Modif, s.Armor

	file, err := os.CreateTemp("", "opennox-object-death-strings-*.json")
	if err != nil {
		panic(err)
	}
	path := file.Name()
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		panic(err)
	}

	armorDef, freeArmorDef := alloc.New(Modifier{})
	weaponDef, freeWeaponDef := alloc.New(Modifier{})
	var freeDesc func()

	reset := func() {
		s.Modif = oldModif
		s.Armor = oldArmor
		if freeDesc != nil {
			freeDesc()
			freeDesc = nil
		}
	}
	configure = func(language int, pluralName string, armor bool) {
		reset()

		entries := []strman.Entry{
			{ID: "Die.c:ArmorDieMetalPlural", Vals: []strman.Variant{{Str: "ARMOR_METAL_PLURAL:%s"}}},
			{ID: "Die.c:ArmorDieMetal", Vals: []strman.Variant{{Str: "ARMOR_METAL:%s"}}},
			{ID: "Die.c:ArmorDieWoodPlural", Vals: []strman.Variant{{Str: "ARMOR_WOOD_PLURAL:%s"}}},
			{ID: "Die.c:ArmorDieWood", Vals: []strman.Variant{{Str: "ARMOR_WOOD:%s"}}},
			{ID: "Die.c:ArmorDieHidePlural", Vals: []strman.Variant{{Str: "ARMOR_HIDE_PLURAL:%s"}}},
			{ID: "Die.c:ArmorDieHide", Vals: []strman.Variant{{Str: "ARMOR_HIDE:%s"}}},
			{ID: "Die.c:ArmorDieClothPlural", Vals: []strman.Variant{{Str: "ARMOR_CLOTH_PLURAL:%s"}}},
			{ID: "Die.c:ArmorDieCloth", Vals: []strman.Variant{{Str: "ARMOR_CLOTH:%s"}}},
			{ID: "Die.c:ArmorDieGeneric", Vals: []strman.Variant{{Str: "ARMOR_GENERIC:%s"}}},
			{ID: "Die.c:WeaponDieMetal", Vals: []strman.Variant{{Str: "WEAPON_METAL:%s"}}},
			{ID: "Die.c:WeaponDieWood", Vals: []strman.Variant{{Str: "WEAPON_WOOD:%s"}}},
			{ID: "Die.c:WeaponDieGeneric", Vals: []strman.Variant{{Str: "WEAPON_GENERIC:%s"}}},
			{ID: "ArmrLook.c:PortTestWeaponName", Vals: []strman.Variant{{Str: "PORTTEST_WEAPON"}}},
		}
		data, err := json.Marshal(struct {
			Lang    int            `json:"lang"`
			Entries []strman.Entry `json:"entries"`
		}{Lang: language, Entries: entries})
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			panic(err)
		}
		sm := strman.New()
		if err := sm.ReadJSON(path); err != nil {
			panic(err)
		}
		s.sm = sm

		desc, free := alloc.CString16(pluralName)
		freeDesc = free
		m := oldModif
		if armor {
			*armorDef = Modifier{TypeInd: 15, Desc8: desc, Next80: oldModif.Dword_5d4594_251608}
			m.Dword_5d4594_251608 = armorDef
		} else {
			*weaponDef = Modifier{TypeInd: 15, Desc8: desc, Next80: oldModif.Dword_5d4594_251600}
			m.Dword_5d4594_251600 = weaponDef
		}
		s.Modif = m

		s.Armor = oldArmor
		s.Armor.table = [PlayerArmorCnt]armorRecord{}
		s.Armor.table[0] = armorRecord{Name: "PortTestWeapon", NameStr: "PortTestWeaponName", TypeInd: 15, Bit: 0x8000}
		s.Armor.ready = true
		s.Armor.sm = sm
	}
	return configure, func() {
		reset()
		s.sm = oldSM
		freeWeaponDef()
		freeArmorDef()
		_ = os.Remove(path)
	}
}
