package legacy

import "C"

import "github.com/opennox/opennox/v1/server"

var (
	Nox_xxx_updateSprings_5113A0 func()
)

//export nox_xxx_updateSprings_5113A0
func nox_xxx_updateSprings_5113A0() { Nox_xxx_updateSprings_5113A0() }
func Nox_xxx_unitHasCollideOrUpdateFn_537610(a1 *server.Object) {
	collisionActivate(a1)
}
