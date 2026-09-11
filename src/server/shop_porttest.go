//go:build porttest

package server

import (
	"strings"

	"github.com/opennox/libs/balance"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// PortTestShopEnvironment layers shop balance and definition inputs over the
// callback fixture. The penalty helper supplies the same three real gem names.
func (s *Server) PortTestShopEnvironment() (configure func(map[string]float64, int, map[int]int), restore func()) {
	_, _, freeGems := s.PortTestPenaltyEnvironment()
	// The shared combat fixture has eight objects. Shop capacity contracts
	// need 128 gold objects plus all 500 stock objects at the same time.
	oldObjectPool := s.Objs.alloc
	if !s.Objs.Init(768) {
		panic("shop object fixture capacity")
	}
	oldBalance, oldSpells := s.Balance.file, s.Spells.byID
	oldSpellServer := s.Spells.s
	s.Spells.s = s
	oldGuide := *s.Types.byInd[2]
	goldData, freeGold := alloc.Make([]byte{}, 4)
	gold := &ObjectType{s: &s.Types, ind: 26, ind2: 26, id: "gold", class: object.ClassSimple, allowed: true, InitData: unsafe.Pointer(&goldData[0]), InitDataSize: 4}
	s.Types.byInd = append(s.Types.byInd, gold)
	s.Types.byID[gold.id] = gold
	overlay := &balance.File{Global: balance.Config{}, Tags: make(map[balance.Tag]balance.Config), Parent: oldBalance}
	s.Balance.file = overlay
	configure = func(values map[string]float64, guideWorth int, prices map[int]int) {
		overlay.Global = balance.Config{}
		for key, value := range values {
			overlay.Global[strings.ToLower(key)] = balance.Array{value}
		}
		s.Types.byInd[2].Worth = guideWorth
		s.Spells.byID = make(map[spell.ID]*SpellDef, len(prices))
		for index, price := range prices {
			id := spell.ID(index)
			s.Spells.byID[id] = &SpellDef{ID: id, Valid: true, Def: things.Spell{Price: price}}
		}
	}
	return configure, func() {
		s.Objs.alloc.Free()
		s.Objs.alloc = oldObjectPool
		*s.Types.byInd[2] = oldGuide
		s.Balance.file, s.Spells.byID = oldBalance, oldSpells
		s.Spells.s = oldSpellServer
		freeGems()
		freeGold()
	}
}
