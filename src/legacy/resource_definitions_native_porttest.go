//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestResourceMonsterSoundABI(u *server.Object) unsafe.Pointer {
	return nox_xxx_monsterGetSoundSet_424300(asObjectC(u))
}
