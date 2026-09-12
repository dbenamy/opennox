package legacy

/*
#include "defs.h"
*/
import "C"
import "github.com/opennox/opennox/v1/server"

func Nox_xxx_cancelAllSpells_4FEE90(a1 *server.Object) {
	spellLifeCancelSelected(a1)
}
