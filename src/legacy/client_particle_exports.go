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

//export nox_thing_magic_draw
func nox_thing_magic_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleMagicDraw(vp, dr, false))
}

//export nox_thing_magic_missle_draw
func nox_thing_magic_missle_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleMagicDraw(vp, dr, true))
}

//export nox_thing_magic_missle_tail_link_draw
func nox_thing_magic_missle_tail_link_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleTailDraw(vp, dr, true))
}

//export nox_thing_magic_tail_link_draw
func nox_thing_magic_tail_link_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleTailDraw(vp, dr, false))
}

//export nox_thing_drain_mana_draw
func nox_thing_drain_mana_draw() C.int { return 1 }

//export nox_thing_bubble_draw
func nox_thing_bubble_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleBubbleDraw(vp, dr))
}

//export nox_thing_blue_rain_draw
func nox_thing_blue_rain_draw(a1 C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(uintptr(uint32(a1))))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleBlueRain(vp, dr))
}

//export nox_thing_levelup_draw
func nox_thing_levelup_draw(a1 C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(uintptr(uint32(a1))))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleLevelUp(vp, dr, false))
}

//export nox_thing_oblivion_up_draw
func nox_thing_oblivion_up_draw(a1 C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(uintptr(uint32(a1))))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleLevelUp(vp, dr, true))
}

//export nox_thing_spider_spit_draw
func nox_thing_spider_spit_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleSpiderSpit(vp, dr))
}

//export nox_thing_vortex_draw
func nox_thing_vortex_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(particleVortexDraw(vp, dr))
}

//export nox_xxx_spriteChangeLightColor_484BE0
func nox_xxx_spriteChangeLightColor_484BE0(p *C.uint32_t, r, g, b C.int) *C.uint32_t {
	return (*C.uint32_t)(particleLightColor(unsafe.Pointer(p), int(r), int(g), int(b)))
}

//export sub_484C00
func sub_484C00(p, angle C.int) C.longlong {
	return C.longlong(particleLightAngle(unsafe.Pointer(uintptr(uint32(p))), int(angle), false))
}

//export nox_xxx_spriteChangeLightSize_484C30
func nox_xxx_spriteChangeLightSize_484C30(p, angle C.int) C.longlong {
	return C.longlong(particleLightAngle(unsafe.Pointer(uintptr(uint32(p))), int(angle), true))
}

//export sub_484CE0
func sub_484CE0(p C.int, value C.float) C.int {
	return C.int(particleLightIntensity(unsafe.Pointer(uintptr(uint32(p))), float32(value), false))
}

//export nox_xxx_spriteChangeIntensity_484D70_light_intensity
func nox_xxx_spriteChangeIntensity_484D70_light_intensity(p C.int, value C.float) C.int {
	return C.int(particleLightIntensity(unsafe.Pointer(uintptr(uint32(p))), float32(value), true))
}
