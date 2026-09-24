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

//export sub_4CA720
func sub_4CA720(unused, drawable C.int) C.int {
	return C.int(effectOrbitUpdate((*client.Drawable)(unsafe.Pointer(uintptr(uint32(drawable))))))
}

//export nox_thing_lightning_draw
func nox_thing_lightning_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectLightningDraw(vp, dr, 0))
}

//export nox_thing_chain_lightning_bolt_draw
func nox_thing_chain_lightning_bolt_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectLightningDraw(vp, dr, 0))
}

//export nox_thing_energy_bolt_draw
func nox_thing_energy_bolt_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectLightningDraw(vp, dr, 1))
}

//export nox_thing_green_bolt_draw
func nox_thing_green_bolt_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectLightningDraw(vp, dr, 2))
}

//export nox_thing_plasma_draw
func nox_thing_plasma_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectPlasmaDraw(vp, dr))
}

//export nox_thing_magic_sparkle_draw
func nox_thing_magic_sparkle_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectMagicSparkle(vp, dr))
}

//export nox_thing_pixie_draw
func nox_thing_pixie_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectPixie(vp, dr))
}

//export nox_thing_pixie_dust_draw
func nox_thing_pixie_dust_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectPixieDust(vp, dr))
}

//export nox_thing_blue_rain_spark_draw
func nox_thing_blue_rain_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectBlueRainSpark(vp, dr))
}

//export nox_thing_rain_orb_draw
func nox_thing_rain_orb_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectRainOrb(vp, dr))
}

//export nox_thing_red_spark_draw
func nox_thing_red_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 0))
}

//export nox_thing_blue_spark_draw
func nox_thing_blue_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 1))
}

//export nox_thing_cyan_spark_draw
func nox_thing_cyan_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 2))
}

//export nox_thing_green_spark_draw
func nox_thing_green_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 3))
}

//export nox_thing_yellow_spark_draw
func nox_thing_yellow_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 4))
}

//export nox_thing_violet_spark_draw
func nox_thing_violet_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 5))
}

//export nox_thing_death_ball_spark_draw
func nox_thing_death_ball_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 6))
}

//export nox_thing_white_spark_draw
func nox_thing_white_spark_draw(viewport *C.uint32_t, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectColoredSpark(vp, dr, 7))
}

//export nox_thing_particle_draw
func nox_thing_particle_draw(viewport C.int, drawable *C.nox_drawable) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectParticleUpdate(dr))
}

//export nox_thing_glow_orb_draw
func nox_thing_glow_orb_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectOrb(vp, dr, false))
}

//export nox_thing_glow_orb_move_draw
func nox_thing_glow_orb_move_draw(viewport *C.int, drawable *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(viewport))
	dr := (*client.Drawable)(unsafe.Pointer(drawable))
	return C.int(effectOrb(vp, dr, true))
}
