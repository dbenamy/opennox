package legacy

/*
#include "defs.h"
void nox_xxx_unitNeedSync_4E44F0(nox_object_t* a1);
nox_object_t* nox_xxx_findObjectAtCursor_54AF40(nox_object_t* a1);
*/
import "C"
import "github.com/opennox/opennox/v1/server"

func Nox_xxx_findObjectAtCursor_54AF40(a1 *server.Object) *server.Object {
	return asObjectS(C.nox_xxx_findObjectAtCursor_54AF40(asObjectC(a1)))
}
func Nox_xxx_updateFallLogic_51B870(a1 *server.Object) {
	motionFall(a1)
}
func Sub_51B810(a1 *server.Object) {
	motionProjectileStep(a1)
}
func Sub_537770(a1 *server.Object) {
	motionProjectileDispatch(a1)
}
