package legacy

/*
extern int nox_cheat_allowall;
*/
import "C"

func CheatEquipAll(v bool) { C.nox_cheat_allowall = C.int(bool2int(v)) }

var spellCharmAll bool

func CheatCharmAll(v bool) { spellCharmAll = v }
