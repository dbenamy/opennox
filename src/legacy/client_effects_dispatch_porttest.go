//go:build porttest

package legacy

/*
#include "GAME3.h"
#include "GAME3_1.h"
#include "client__draw__fx.h"
#include "client__draw__drawrays.h"
#include "client__draw__lightning.h"
#include "client__draw__plasma.h"
#include "client__draw__glowdraw.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// PortTestClientEffects exercises retained production C entry points and native private helpers. The caller
// owns bounded input storage, drawable lists, renderer and RNG. Curve callback
// recording has a separate typed fixture to preserve its external ABI.
func PortTestClientEffects(op int, vp *noxrender.Viewport, dr *client.Drawable, a [8]int32, data unsafe.Pointer) uint32 {
	switch op {
	case 45:
		return uint32(effectPrepareLightning())
	case 0:
		sub_499490(C.int(a[0]), (*C.uint16_t)(data), C.int(a[1]), C.int(a[2]), C.char(a[3]), C.char(a[4]))
		return 0
	case 1:
		C.sub_499520(C.int(a[0]), (*C.short)(data), C.short(a[1]), C.char(a[2]), C.char(a[3]))
		return 0
	case 2:
		return uint32(nox_xxx_makePointFxCli_499610(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5])))
	case 3:
		return uint32(nox_xxx_drawEnergyBolt_499710(C.int(a[0]), C.int(a[1]), C.short(a[2]), C.int(a[3])))
	case 4:
		return uint32(C.sub_499950(C.int(a[0]), (*C.int2)(unsafe.Add(data, 0)), (*C.int2)(unsafe.Add(data, 8)), C.ushort(a[1]), C.char(a[2])))
	case 5:
		return uint32(nox_xxx_makeLightningParticles_4999D0(C.int(a[0]), (*C.int2)(unsafe.Add(data, 0)), (*C.int2)(unsafe.Add(data, 8))))
	case 6:
		return uint32(C.nox_xxx_draw_499E70(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]), C.int(a[6])))
	case 7:
		return uint32(sub_49A150((*C.int2)(unsafe.Add(data, 0)), C.int(a[0]), C.uchar(a[1])))
	case 8:
		return uint32(nox_xxx_netDrawRays_49BDD0((*C.uchar)(data)))
	case 9:
		return uint32(effectLightningStep(uint32(a[0]), uint32(a[1])))
	case 10:
		return uint32(effectLightningPasses(AsPoint(data), AsPoint(unsafe.Add(data, 8)), int(a[0]), (*[4]int16)(unsafe.Add(data, 16)), int(a[1]), int(a[2]), int(a[3])))
	case 11:
		return uint32(C.nox_thing_lightning_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 12:
		return uint32(C.nox_thing_chain_lightning_bolt_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 13:
		return uint32(C.nox_thing_energy_bolt_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 14:
		return uint32(C.nox_thing_green_bolt_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 15:
		return uint32(effectPlasma(int(a[0]), image.Pt(int(a[1]), int(a[2])), image.Pt(int(a[3]), int(a[4]))))
	case 16:
		return uint32(effectPlasmaSegment(AsPoint(data), AsPoint(unsafe.Add(data, 8)), int(*(*int32)(unsafe.Add(data, 16)))))
	case 17:
		return uint32(C.nox_thing_plasma_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 18:
		return uint32(effectSparkDraw(vp, dr, uint32(a[0]), uint32(a[1]), true))
	case 19:
		return uint32(C.nox_thing_magic_sparkle_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 20:
		return uint32(C.nox_thing_pixie_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 21:
		return uint32(C.nox_thing_pixie_dust_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 22:
		return uint32(C.nox_thing_blue_rain_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 23:
		return uint32(C.nox_thing_rain_orb_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 24:
		return uint32(C.nox_thing_red_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 25:
		return uint32(C.nox_thing_blue_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 26:
		return uint32(C.nox_thing_cyan_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 27:
		return uint32(C.nox_thing_green_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 28:
		return uint32(C.nox_thing_yellow_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 29:
		return uint32(C.nox_thing_violet_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 30:
		return uint32(C.nox_thing_death_ball_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 31:
		return uint32(C.nox_thing_white_spark_draw((*C.uint32_t)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 32:
		return uint32(C.nox_thing_particle_draw(C.int(uintptr(unsafe.Pointer(vp))), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 33:
		return uint32(C.nox_thing_glow_orb_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 34:
		return uint32(C.nox_thing_glow_orb_move_draw((*C.int)(unsafe.Pointer(vp)), (*C.nox_drawable)(unsafe.Pointer(dr))))
	case 35:
		effectPlasmaSetup(int(a[0]), image.Pt(int(a[1]), int(a[2])), image.Pt(int(a[3]), int(a[4])))
		return 0
	case 36:
		return uint32(effectOrb(vp, dr, a[0] != 0))
	case 37:
		return uint32(C.sub_4CA720(C.int(uintptr(unsafe.Pointer(vp))), C.int(uintptr(unsafe.Pointer(dr)))))
	case 38:
		return uint32(effectCurveColor(int(a[0])))
	case 39:
		return uint32(effectCurveGlow(int(a[0]), int(a[1]), int(a[2]), int8(a[3])))
	case 40:
		effectCurveRaster([4]image.Point{AsPoint(data), AsPoint(unsafe.Add(data, 8)), AsPoint(unsafe.Add(data, 16)), AsPoint(unsafe.Add(data, 24))}, int(a[0]), int(a[1]))
		return 0
	case 42:
		return uint32(effectSparkDraw(vp, dr, uint32(a[0]), uint32(a[1]), false))
	case 43:
		return uint32(effectMovingSpark(vp, dr, uint32(a[0]), uint32(a[1])))
	case 44:
		return uint32(effectSparkBounce(dr))
	default:
		panic("unknown client effect operation")
	}
}

func PortTestEffectsCallback(op int) unsafe.Pointer {
	switch op {
	case 11:
		return unsafe.Pointer(C.nox_thing_lightning_draw)
	case 12:
		return unsafe.Pointer(C.nox_thing_chain_lightning_bolt_draw)
	case 13:
		return unsafe.Pointer(C.nox_thing_energy_bolt_draw)
	case 14:
		return unsafe.Pointer(C.nox_thing_green_bolt_draw)
	case 17:
		return unsafe.Pointer(C.nox_thing_plasma_draw)
	case 19:
		return unsafe.Pointer(C.nox_thing_magic_sparkle_draw)
	case 20:
		return unsafe.Pointer(C.nox_thing_pixie_draw)
	case 21:
		return unsafe.Pointer(C.nox_thing_pixie_dust_draw)
	case 22:
		return unsafe.Pointer(C.nox_thing_blue_rain_spark_draw)
	case 23:
		return unsafe.Pointer(C.nox_thing_rain_orb_draw)
	case 24:
		return unsafe.Pointer(C.nox_thing_red_spark_draw)
	case 25:
		return unsafe.Pointer(C.nox_thing_blue_spark_draw)
	case 26:
		return unsafe.Pointer(C.nox_thing_cyan_spark_draw)
	case 27:
		return unsafe.Pointer(C.nox_thing_green_spark_draw)
	case 28:
		return unsafe.Pointer(C.nox_thing_yellow_spark_draw)
	case 29:
		return unsafe.Pointer(C.nox_thing_violet_spark_draw)
	case 30:
		return unsafe.Pointer(C.nox_thing_death_ball_spark_draw)
	case 31:
		return unsafe.Pointer(C.nox_thing_white_spark_draw)
	case 32:
		return unsafe.Pointer(C.nox_thing_particle_draw)
	case 33:
		return unsafe.Pointer(C.nox_thing_glow_orb_draw)
	case 34:
		return unsafe.Pointer(C.nox_thing_glow_orb_move_draw)
	case 37:
		return unsafe.Pointer(C.sub_4CA720)
	default:
		panic("unknown effect callback")
	}
}
