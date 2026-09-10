package legacy

import "C"

import (
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
)

//export nox_xxx_playerCheckSpellClass_57AEA0
func nox_xxx_playerCheckSpellClass_57AEA0(class, ind C.int) C.int {
	flags := GetServer().S().Spells.Flags(spell.ID(ind))
	// Keep the raw int: narrowing to the byte-sized class enum accepts invalid inputs.
	var allowed things.SpellFlags
	switch class {
	case 1:
		allowed = things.SpellClassWizard
	case 2:
		allowed = things.SpellClassConjurer
	default:
		return 9
	}
	if flags&(things.SpellClassAny|allowed) != 0 {
		return 0
	}
	return 9
}
