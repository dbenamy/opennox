package legacy

/*
#include "GAME5.h"
*/
import "C"

// The registered PlayerDie callback ignores the historical integer return value.
//
//export nox_xxx_diePlayer_54D2B0
func nox_xxx_diePlayer_54D2B0(u C.int) C.int {
	playerDeath(objectFromInt(u))
	return 0
}
