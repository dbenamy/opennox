package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

//export nox_thing_harpoon_draw
func nox_thing_harpoon_draw(vp *C.int, dr *C.nox_drawable) C.int {
	return C.int(spriteSlaveDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_harpoon_rope_draw
func nox_thing_harpoon_rope_draw(vp *C.int, dr *C.nox_drawable) C.int {
	return C.int(screenHarpoonRope((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_undead_killer_draw
func nox_thing_undead_killer_draw(vp *C.int, dr *C.nox_drawable) C.int {
	return C.int(screenUndead((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_maiden_draw
func nox_thing_maiden_draw(vp *C.uint32_t, dr *C.nox_drawable) C.int {
	return C.int(screenMaiden((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_player_waypoint_draw
func nox_thing_player_waypoint_draw(vp C.int, dr *C.nox_drawable) C.int {
	return C.int(screenWaypoint((*noxrender.Viewport)(unsafe.Pointer(uintptr(vp))), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_client_screenParticleDraw_489700
func nox_client_screenParticleDraw_489700(vp unsafe.Pointer, p *C.nox_screenParticle) C.int {
	return C.int(screenParticleDraw((*noxrender.Viewport)(vp), (*Nox_screenParticle)(unsafe.Pointer(p))))
}

//export nox_client_newScreenParticle_431540
func nox_client_newScreenParticle_431540(kind, x, y, vx, vy, gravity C.int, size, timer, phase, mode C.char) *C.nox_screenParticle {
	return (*C.nox_screenParticle)(unsafe.Pointer(screenParticleCreate(int(kind), int(x), int(y), int(vx), int(vy), int(gravity), byte(size), byte(timer), byte(phase), byte(mode))))
}

//export sub_48C6B0
func sub_48C6B0(x, y C.int) C.uint { return C.uint(screenDistance(int32(x), int32(y))) }
