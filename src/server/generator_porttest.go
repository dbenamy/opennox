//go:build porttest

package server

import (
	"strings"
	"unsafe"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestGeneratorTypeIDs are the stable TypeInd values appended to the
// callback fixture's 23-entry registry. Beholder is resolved through the C
// generator cache; Monster is the ordinary source TypeInd used for spawning.
type PortTestGeneratorTypeIDs struct {
	DestroyedGenerator, Beholder, Monster, Glyph uint16
}

// PortTestGeneratorEnvironment extends PortTestAICallbackTypes with the
// minimum real object definitions used by the generator family. It does not
// create objects itself: legacy tests must free every allocated object before
// restore, because NewObject copies these C-owned update and health sources.
//
// configure replaces one overlay's values, rather than layering overlays. Its
// intended keys are QuestHardcoreStage, QuestHardcoreSpawnRateIncrease,
// QuestHardcoreSpawnCap, SpawnRate{High,Normal,Low,VeryLow,VeryVeryLow}Value,
// and MaxOnscreenMonsterCount; unrelated keys still fall through to the prior
// live balance file.
func (s *Server) PortTestGeneratorEnvironment() (ids PortTestGeneratorTypeIDs, configure func(map[string]float64), restore func()) {
	if len(s.Types.byInd) != 23 || s.Types.ByID("Beholder") != nil {
		panic("PortTestGeneratorEnvironment requires PortTestAICallbackTypes")
	}

	oldTypes := s.Types
	oldWeapons, oldArmor := s.Weapons, s.Armor
	s.Weapons.table = [PlayerWeaponCnt]weaponRecord{}
	s.Weapons.table[0] = weaponRecord{TypeInd: 15, Bit: 0x100}
	s.Weapons.ready = true
	s.Armor.table = [PlayerArmorCnt]armorRecord{}
	s.Armor.table[0] = armorRecord{TypeInd: 16, Bit: 0x1000000}
	s.Armor.ready = true
	byInd := append([]*ObjectType(nil), oldTypes.byInd...)
	byID := make(map[string]*ObjectType, len(oldTypes.byID)+4)
	for k, v := range oldTypes.byID {
		byID[k] = v
	}
	s.Types.byInd = append(byInd, make([]*ObjectType, 4)...)
	s.Types.byID = byID

	ids = PortTestGeneratorTypeIDs{DestroyedGenerator: 23, Beholder: 24, Monster: 25, Glyph: 26}
	destroyed := &ObjectType{
		s: &s.Types, ind: ids.DestroyedGenerator, ind2: ids.DestroyedGenerator,
		id: "destroyedgenerator", class: object.ClassSimple, allowed: true, Mass: 1,
	}

	makeMonster := func(ind uint16, id string) (*ObjectType, func()) {
		data, freeData := alloc.Make([]byte{}, 2200)
		health, freeHealth := alloc.New(HealthData{})
		*health = HealthData{Cur: 100, Max: 100}
		t := &ObjectType{
			s: &s.Types, ind: ind, ind2: ind, id: id,
			class: object.ClassMonster, allowed: true, Mass: 1,
			health: health, UpdateData: unsafe.Pointer(unsafe.SliceData(data)), UpdateDataSize: uintptr(len(data)),
		}
		return t, func() { freeHealth(); freeData() }
	}
	beholder, freeBeholder := makeMonster(ids.Beholder, "beholder")
	monster, freeMonster := makeMonster(ids.Monster, "porttestgeneratormonster")
	glyphInit, freeGlyph := alloc.Make([]byte{}, 64)
	glyph := &ObjectType{s: &s.Types, ind: ids.Glyph, ind2: ids.Glyph, id: "glyph", class: object.ClassSimple, allowed: true, Mass: 1, InitData: unsafe.Pointer(&glyphInit[0]), InitDataSize: 64}
	for _, typ := range []*ObjectType{destroyed, beholder, monster, glyph} {
		s.Types.byInd[int(typ.ind)] = typ
		s.Types.byID[typ.id] = typ
	}

	oldBalance := s.Balance.file
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: oldBalance}
	s.Balance.file = overlay
	configure = func(values map[string]float64) {
		overlay.Global = balance.Config{}
		for key, value := range values {
			overlay.Global[strings.ToLower(key)] = balance.Array{value}
		}
	}
	restore = func() {
		s.Types = oldTypes
		s.Weapons, s.Armor = oldWeapons, oldArmor
		s.Balance.file = oldBalance
		freeGlyph()
		freeMonster()
		freeBeholder()
	}
	return ids, configure, restore
}
