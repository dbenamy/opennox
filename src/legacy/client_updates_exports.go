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

//export nox_xxx_updDrawDBallCharge_4CE0C0
func nox_xxx_updDrawDBallCharge_4CE0C0(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateDeathBallCharge(dr))
}

//export sub_4CD690
func sub_4CD690(a1 *C.uint32_t, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	return C.int(updateHealDrain(vp, dr, false))
}

//export sub_4CD450
func sub_4CD450(a1 *C.uint32_t, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	return C.int(updateHealDrain(vp, dr, true))
}

//export nox_xxx_updDrawManabombCharge_4CCAC0
func nox_xxx_updDrawManabombCharge_4CCAC0(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateManaBomb(dr))
}

//export nox_xxx_updDrawMagicMissile_4CD9E0
func nox_xxx_updDrawMagicMissile_4CD9E0(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateMagicMissile(dr))
}

//export nox_xxx_updDrawMagic_4CDD80
func nox_xxx_updDrawMagic_4CDD80(a1 C.int, a2 *C.uint32_t) {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	updateMagicTrail(dr)
}

//export nox_xxx_updDrawSparkleTrail_4CDBF0
func nox_xxx_updDrawSparkleTrail_4CDBF0(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateTrailSparks(dr, false))
}

//export nox_xxx_updDrawTeleportWake_4CD8D0
func nox_xxx_updDrawTeleportWake_4CD8D0(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateTeleportWake(dr))
}

//export nox_xxx_updDrawVortexSource_4CC950
func nox_xxx_updDrawVortexSource_4CC950(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateVortex(dr))
}

//export nox_xxx_updDrawUndeadKiller_4CCCF0
func nox_xxx_updDrawUndeadKiller_4CCCF0() C.int {
	return 1
}

//export sub_4CCD00
func sub_4CCD00(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateHeight(dr, false))
}

//export nox_xxx_updDrawFist_4CCDB0
func nox_xxx_updDrawFist_4CCDB0(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateHeight(dr, true))
}

//export sub_4CCE70
func sub_4CCE70(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateFireballFrame(dr, 5))
}

//export sub_4CD090
func sub_4CD090(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateFireballFrame(dr, 4))
}

//export sub_4CD0C0
func sub_4CD0C0(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateFireballFrame(dr, 3))
}

//export sub_4CD0F0
func sub_4CD0F0(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateFireballFrame(dr, 2))
}

//export sub_4CD120
func sub_4CD120(a1 C.int, a2 *C.uint32_t) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(updateFireballFrame(dr, 1))
}

//export sub_4CD400
func sub_4CD400(a1 *C.uint32_t, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	return C.int(updateCharm(vp, dr))
}

//export nox_xxx_updDrawDBall_4CDF80
func nox_xxx_updDrawDBall_4CDF80(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	updateDeathBallSparks(dr, 3)
	return 1
}

//export sub_4CE0A0
func sub_4CE0A0(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	updateDeathBallSparks(dr, 1)
	return 1
}

//export nox_xxx_updDrawCloud_4CE1D0
func nox_xxx_updDrawCloud_4CE1D0(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateCloudFrame(dr, 75))
}

//export sub_4CE340
func sub_4CE340(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateCloudRise(dr))
}

//export sub_4CE360
func sub_4CE360(a1 C.int, a2 C.int) C.int {
	dr := (*client.Drawable)(unsafe.Pointer(uintptr(uint32(a2))))
	return C.int(updateCloudFrame(dr, 35))
}
