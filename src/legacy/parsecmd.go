package legacy

/*
 */
import "C"

func CheatEquipAll(v bool) { nox_cheat_allowall = C.int(bool2int(v)) }

var spellCharmAll bool

func CheatCharmAll(v bool) { spellCharmAll = v }
