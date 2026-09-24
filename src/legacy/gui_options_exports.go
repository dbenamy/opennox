package legacy

/*
#include "defs.h"
*/
import "C"

//export nox_game_showOptions_4AA6B0
func nox_game_showOptions_4AA6B0() C.int { return C.int(optionsMenu.construct()) }

//export sub_4AB0C0
func sub_4AB0C0() C.int { return C.int(optionsMenuDone()) }
