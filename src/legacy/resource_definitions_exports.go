package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import "unsafe"

//export nox_xxx_monsterGetSoundSet_424300
func nox_xxx_monsterGetSoundSet_424300(u *nox_object_t) unsafe.Pointer {
	return resourceMonsterSound(asObjectS(u))
}
