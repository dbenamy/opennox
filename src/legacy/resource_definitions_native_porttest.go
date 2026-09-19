//go:build porttest

package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestResourceMonsterSoundABI(u *server.Object) unsafe.Pointer {
	return C.nox_xxx_monsterGetSoundSet_424300(asObjectC(u))
}
