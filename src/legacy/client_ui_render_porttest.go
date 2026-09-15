//go:build porttest

package legacy

/*
#include "GAME1_2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
extern nox_render_data_t* nox_draw_curDrawData_3799572;
*/
import "C"
import "unsafe"

func PortTestUIRenderIntersect(rects *[12]int32, dst int) bool {
	return uiRenderIntersect((*[4]int32)(unsafe.Pointer(&rects[dst*4])), (*[4]int32)(unsafe.Pointer(&rects[4])), (*[4]int32)(unsafe.Pointer(&rects[8]))) != nil
}
func PortTestUIRenderOp(op int, a [6]int32) int32 {
	switch op {
	case 0:
		return int32(uiRenderFlag(uint32(a[0])))
	case 1:
		return int32(C.nox_client_copyRect_49F6F0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case 2:
		return int32(uiRenderNarrowClip(int(a[0]), int(a[1])))
	case 3:
		C.sub_49CD30(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]))
		return 0
	case 4:
		return int32(C.sub_430B50(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	default:
		panic("unknown UI render operation")
	}
}
func PortTestUIRenderFill(data []byte, offset int, color uint32, size int32, wrapper bool) int32 {
	p := unsafe.Pointer(&data[offset])
	if wrapper {
		return int32(C.sub_49D1C0(p, C.int(color), C.int(size)))
	}
	uiRenderFill(p, color, size)
	return 0
}
