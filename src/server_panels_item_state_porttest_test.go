//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestServerPanelsItemMasks(t *testing.T) {
	o := newServerOptionsOwner(t)
	oldTypes, oldWeapons, oldArmor := o.c.srv.Types, o.c.srv.Weapons, o.c.srv.Armor
	t.Cleanup(func() { o.c.srv.Types = oldTypes; o.c.srv.Weapons = oldWeapons; o.c.srv.Armor = oldArmor })
	words, free := alloc.Make([]uint32{}, 3)
	defer free()
	type row struct {
		Pattern, Weapon, Armor uint32
		Result                 int
	}
	var rows []row
	for _, pattern := range []uint32{0, 1, 0x55555555, 0xaaaaaaaa, 0x7fffffff, 0x80000000, 0xffffffff} {
		t.Run(fmt.Sprintf("%x", pattern), func(t *testing.T) {
			s, desc := server.PortTestRuleServerSetup()
			// These production owners refer back to s.Types, which stays alive here.
			o.c.srv.Types = s.Types
			o.c.srv.Weapons = s.Weapons
			o.c.srv.Armor = s.Armor
			wantWeapon, wantArmor := uint32(0xffffffff), uint32(0xffffffff)
			for _, it := range desc.Weapons {
				allow := it.Bit&pattern != 0
				s.Types.ByID(it.Name).SetAllowed(allow)
				if !allow && it.Bit < 1<<27 {
					wantWeapon &^= it.Bit
				}
			}
			for _, it := range desc.Armors {
				allow := it.Bit&pattern != 0
				s.Types.ByID(it.Name).SetAllowed(allow)
				if !allow && it.Bit < 1<<26 {
					wantArmor &^= it.Bit
				}
			}
			words[0], words[1], words[2] = 0x12345678, 0x55667788, 0x87654321
			result := legacy.PortTestServerPanelsWeaponSnapshot(&words[1])
			armor := uint32(legacy.PortTestServerPanelsArmorSnapshot())
			if words[1] != wantWeapon || armor != wantArmor || words[0] != 0x12345678 || words[2] != 0x87654321 {
				t.Fatalf("masks %x/%x want %x/%x", words[1], armor, wantWeapon, wantArmor)
			}
			// Last weapon bit determines the incidental scalar return.
			wantResult := 0
			for _, it := range desc.Weapons {
				if it.Bit == 1<<26 {
					if pattern&it.Bit != 0 {
						wantResult = 1
					} else {
						wantResult = 251
					}
				}
			}
			if result != wantResult {
				t.Fatalf("result %d want %d", result, wantResult)
			}
			rows = append(rows, row{pattern, words[1], armor, result})
		})
	}
	spellbookCapture(t, "server-panels-item-masks", rows, "6e50814f5b70cad9d9f7ce72bfbd9ecf14a860b1ae523a86ef8601eb441558f5")
}
