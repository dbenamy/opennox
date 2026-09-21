package legacy

/*
#include <stdint.h>
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_xxx_soundPlayerDamageSound_5328B0  func(obj, obj2 *server.Object) int
	Nox_xxx_soundDefaultDamageSound_532E20 func(obj, obj2 *server.Object) int
)

func Nox_xxx_clientPlaySoundSpecial_452D80(a1 sound.ID, a2 int) {
	audioEventPlay(int32(a1), int32(a2), 0, 0)
}

func Sub_4133D0(a1 *server.Object) int {
	return bool2int(runtimeMaterial(a1))
}

//export nox_xxx_soundPlayerDamageSound_5328B0
func nox_xxx_soundPlayerDamageSound_5328B0(a1, a2 *nox_object_t) int {
	return Nox_xxx_soundPlayerDamageSound_5328B0(asObjectS(a1), asObjectS(a2))
}

//export nox_xxx_soundDefaultDamageSound_532E20
func nox_xxx_soundDefaultDamageSound_532E20(a1, a2 *nox_object_t) int {
	return Nox_xxx_soundDefaultDamageSound_532E20(asObjectS(a1), asObjectS(a2))
}
