package legacy

/*
#include "defs.h"
const char** nox_xxx_getDefaultSoundSet_424350(const char* a1);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_monsterCreateFn_54C480 func(u *server.Object)
)

//export nox_xxx_monsterCreateFn_54C480
func nox_xxx_monsterCreateFn_54C480(u *nox_object_t) {
	Nox_xxx_monsterCreateFn_54C480(asObjectS(u))
}

func Nox_xxx_monsterDefByTT_517560(typ int) *server.MonsterDef {
	return monsterDefinitionByType(uint32(typ))
}

func Nox_xxx_monsterAutoSpells_54C0C0(u *server.Object) {
	monsterAutoSpells(u)
}

func Nox_xxx_getDefaultSoundSet_424350(name string) unsafe.Pointer {
	return resourceSoundByName(name)
}
