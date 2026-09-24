package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_xxx_XFerMonster_528DB0
func nox_xxx_XFerMonster_528DB0(u *nox_object_t) C.int {
	return C.int(creatureXferMonster(asObjectS(u)))
}

//export nox_xxx_XFerNPC_52ADE0
func nox_xxx_XFerNPC_52ADE0(u *nox_object_t) C.int {
	return C.int(creatureXferNPC(asObjectS(u)))
}
