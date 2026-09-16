//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_2.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
)

func PortTestColorLight(op int, vp *noxrender.Viewport, d *client.Drawable) int {
	p := C.int(uintptr(d.C()))
	switch op {
	case 0:
		C.nox_xxx_colorLightAlterB0ArrayColor_4CE440(p)
	case 1:
		C.nox_xxx_colorLightAlterIntensity_4CE610(p)
	case 2:
		return int(C.nox_xxx_colorLightAlterRadius_4CE760(p))
	case 3:
		C.sub_4CE8C0(p)
	case 4:
		C.sub_4CE960(p)
	case 5:
		return int(C.nox_xxx_updDrawColorlight_4CE390((*C.uint32_t)(vp.C()), p))
	default:
		panic(op)
	}
	return 0
}
