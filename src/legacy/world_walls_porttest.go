//go:build porttest

package legacy

/*
#include "GAME2_1.h"
#include "GAME2_2.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/server"
	"image"
	"unsafe"
)

// Calls original production C; expected outputs are captured before conversion.
func PortTestWorldWalls(op int, vp *noxrender.Viewport, dr *client.Drawable, wl *server.Wall, p image.Point) (int, image.Point) {
	var out image.Point
	v := (*C.nox_draw_viewport_t)(unsafe.Pointer(vp))
	drawable := (*C.nox_drawable)(unsafe.Pointer(dr))
	switch op {
	case 0:
		ret := C.sub_4739E0((*C.uint32_t)(unsafe.Pointer(v)), (*C.int2)(unsafe.Pointer(&p)), (*C.int2)(unsafe.Pointer(&out)))
		return int(ret), out
	case 1:
		ret := C.sub_473A10((*C.uint32_t)(unsafe.Pointer(v)), (*C.int2)(unsafe.Pointer(&p)), (*C.uint32_t)(unsafe.Pointer(&out)))
		return int(ret), out
	case 2:
		C.nox_xxx_drawWalls_473C10(v, unsafe.Pointer(wl))
		return 0, out
	case 3:
		return int(C.sub_474B40(drawable)), out
	case 4:
		return int(C.nox_xxx_sprite_4756E0_drawable(drawable)), out
	case 5:
		return int(C.nox_xxx_sprite_475740_drawable(drawable)), out
	case 6:
		return int(C.nox_xxx_sprite_4757A0_drawable(drawable)), out
	case 7:
		return int(C.sub_4757D0_drawable(drawable)), out
	case 8:
		return int(C.sub_47D380(C.int(p.X), C.int(p.Y))), out
	}
	panic("unknown world-wall operation")
}
