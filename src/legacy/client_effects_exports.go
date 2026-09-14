package legacy

/*
#include "defs.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"image"
	"unsafe"
)

//export sub_499490
func sub_499490(typ C.int, p *C.uint16_t, x, y C.int, speed, period C.char) {
	effectCreateOrb(int(typ), (*[4]uint16)(unsafe.Pointer(p)), int(x), int(y), byte(speed), byte(period))
}

//export sub_499520
func sub_499520(typ C.int, p *C.short, angle C.short, direction, period C.char) {
	effectCreateOrbit(int(typ), (*[4]int16)(unsafe.Pointer(p)), int16(angle), byte(direction), byte(period))
}

//export nox_xxx_makePointFxCli_499610
func nox_xxx_makePointFxCli_499610(typ, count, speed, ttl, x, y C.int) C.int {
	return C.int(effectCreatePointSparks(int(typ), int(count), int(speed), int(ttl), int(x), int(y)))
}

//export nox_xxx_drawEnergyBolt_499710
func nox_xxx_drawEnergyBolt_499710(x, y C.int, z C.short, typ C.int) C.int {
	return C.int(effectCreateEnergySparks(int(x), int(y), int16(z), int(typ)))
}

//export sub_499950
func sub_499950(typ C.int, from, to *C.int2, z C.ushort, velocity C.char) C.int {
	dr := effectCreateRainOrb(int(typ), AsPoint(unsafe.Pointer(from)), AsPoint(unsafe.Pointer(to)), uint16(z), int8(velocity))
	return C.int(uintptr(unsafe.Pointer(dr)))
}

//export nox_xxx_makeLightningParticles_4999D0
func nox_xxx_makeLightningParticles_4999D0(typ C.int, from, to *C.int2) C.int {
	return C.int(effectLightningParticles(int(typ), AsPoint(unsafe.Pointer(from)), AsPoint(unsafe.Pointer(to))))
}

//export nox_xxx_draw_499E70
func nox_xxx_draw_499E70(kind, x, y, width, height, axis, direction C.int) C.int {
	return C.int(effectScreenParticles(int(kind), int(x), int(y), int(width), int(height), int(axis), int(direction)))
}

//export sub_49A150
func sub_49A150(pos *C.int2, typ C.int, amount C.uchar) C.int {
	return C.int(effectSparkBurst(AsPoint(unsafe.Pointer(pos)), int(typ), byte(amount)))
}

//export nox_xxx_netDrawRays_49BDD0
func nox_xxx_netDrawRays_49BDD0(packet *C.uchar) C.int {
	return C.int(effectDispatchRay((*[9]byte)(unsafe.Pointer(packet))))
}

//export sub_4CA720
func sub_4CA720(unused, drawable C.int) C.int {
	return C.int(effectOrbitUpdate((*client.Drawable)(unsafe.Pointer(uintptr(uint32(drawable))))))
}

//export sub_4BEDE0
func sub_4BEDE0(a, b, c, d *C.int2, steps C.int, shift C.float, callback, token C.int) {
	points := [4]image.Point{AsPoint(unsafe.Pointer(a)), AsPoint(unsafe.Pointer(b)), AsPoint(unsafe.Pointer(c)), AsPoint(unsafe.Pointer(d))}
	fn := unsafe.Pointer(uintptr(uint32(callback)))
	effectCurveSegments(points, int(steps), float32(shift), func(from, to image.Point) {
		// Typed pointers make the stack-copy/lifetime boundary visible to Go while
		// preserving the external three-word callback ABI and opaque token bits.
		ccall.CallVoidPtr3(fn, unsafe.Pointer(&from), unsafe.Pointer(&to), unsafe.Pointer(uintptr(uint32(token))))
	})
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
