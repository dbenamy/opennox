//go:build porttest

package legacy

/*
#include "GAME2_1.h"
extern uint32_t nox_xxx_minimap_587000_149232;
*/
import "C"
import "unsafe"

// PortTestMinimap invokes the real C owner; no rendering/selection substitute.
func PortTestMinimap(op int, a, b, c, d, e int) uint32 {
	point := [2]int32{int32(a), int32(b)}
	p := (*C.int2)(unsafe.Pointer(&point[0]))
	switch op {
	case 0:
		C.nox_client_mapZoomIn_4724E0()
		return uint32(C.nox_xxx_minimap_587000_149232)
	case 1:
		C.nox_client_mapZoomOut_472500()
		return uint32(C.nox_xxx_minimap_587000_149232)
	case 2:
		return uint32(C.nox_xxx_cliSetMinimapZoom_472520(C.int(a)))
	case 3:
		return uint32(C.sub_4730D0(p, C.uchar(c), C.int(d)))
	case 4:
		return uint32(C.sub_473380(C.int(a), C.int(b), C.int(c), C.int(d), C.int(e)))
	case 5:
		return uint32(C.sub_4733B0((*C.uint32_t)(unsafe.Pointer(p))))
	case 6:
		return uint32(C.sub_473420((*C.uint32_t)(unsafe.Pointer(p))))
	case 7:
		C.nox_video_drawCircleRad3_4734F0((*C.int)(unsafe.Pointer(p)))
		return 0
	case 8:
		return uint32(C.nox_client_drawRectLines_473510(C.int(a), C.int(b), C.int(c), C.int(d)))
	case 9:
		C.nox_xxx_minimapDrawPoint_473570(C.int(a), C.int(b))
		return 0
	case 10:
		C.sub_4735C0(C.int(a), C.int(b))
		return 0
	case 11:
		return uint32(C.sub_472540(C.int(a)))
	case 12:
		C.nox_xxx_drawMinimap4Sprite_4725C0(C.int(a))
		return 0
	case 13:
		return uint32(C.nox_xxx_cliDrawMinimap_472600(C.int(a), C.int(b)))
	case 14:
		return uint32(C.nox_xxx_drawMinimapAndLines_4738E0())
	}
	panic("unknown minimap operation")
}
func PortTestMinimapZoom() (*uint32, func()) {
	p := (*uint32)(unsafe.Pointer(&C.nox_xxx_minimap_587000_149232))
	old := *p
	return p, func() { *p = old }
}
