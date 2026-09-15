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
	return C.nox_xxx_utilRect_49F930((*C.int4)(unsafe.Pointer(&rects[dst*4])), (*C.int4)(unsafe.Pointer(&rects[4])), (*C.int4)(unsafe.Pointer(&rects[8]))) != nil
}
func PortTestUIRenderOp(op int, a [6]int32) int32 {
	var p unsafe.Pointer
	switch op {
	case 0:
		return int32(C.sub_49F6D0(C.int(a[0])))
	case 1:
		p = unsafe.Pointer(C.nox_client_copyRect_49F6F0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case 2:
		p = unsafe.Pointer(C.sub_49F780(C.int(a[0]), C.int(a[1])))
	case 3:
		C.sub_49CD30(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]))
		return 0
	case 4:
		return int32(C.sub_430B50(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	default:
		panic("unknown UI render operation")
	}
	if p == nil {
		return 0
	}
	// noxCopyRect returns 1, despite the legacy pointer return declaration.
	return int32(uintptr(p))
}
func PortTestUIRenderFill(data []byte, offset int, color uint32, size int32, wrapper bool) int32 {
	p := unsafe.Pointer(&data[offset])
	if wrapper {
		return int32(C.sub_49D1C0(p, C.int(color), C.int(size)))
	}
	C.sub_49E3C0((*C.uint32_t)(p), C.int(color), C.uint(size))
	return 0
}
