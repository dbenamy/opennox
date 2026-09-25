//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// The caller owns the real drawable pool/index, alias table and outgoing queue.
// Packet storage is C-owned and padded; this exercises valid legacy records,
// including early returns, rather than assigning new malformed-input semantics.
func PortTestDrawableStream(stream int, input []byte, pos [2]int32) (int32, [2]int32) {
	packet, free := alloc.Make([]byte{}, len(input)+16)
	defer free()
	copy(packet, input)
	coords, freeCoords := alloc.Make([]int32{pos[0], pos[1], 0x12345678}, 3)
	defer freeCoords()
	var ret int32
	if stream == 1 {
		ret = int32(nox_xxx_netCliProcUpdateStream_494A60((*uint8)(unsafe.Pointer(&packet[0])), int32(31), (*uint32)(unsafe.Pointer(&coords[0]))))
	} else {
		ret = int32(uintptr(unsafe.Pointer(nox_xxx_netCliUpdateStream2_494C30((*uint8)(unsafe.Pointer(&packet[0])), int32(31), (*int32)(unsafe.Pointer(&coords[0]))))))
	}
	if coords[2] != 0x12345678 {
		panic("stream wrote beyond coordinate pair")
	}
	return ret, [2]int32{coords[0], coords[1]}
}
func PortTestDrawableState(op int, dr *client.Drawable, value int) uintptr {
	p := (*C.nox_drawable)(dr.C())
	switch op {
	case 0:
		return uintptr(drawableStatePredicate(dr))
	case 1:
		return uintptr(nox_xxx_spriteSetActiveMB_45A990_drawable(C.int(uintptr(dr.C()))))
	case 2:
		return uintptr(nox_xxx_spriteSetFrameMB_45AB80(C.int(uintptr(dr.C())), C.int(value)))
	case 3:
		return uintptr(unsafe.Pointer(C.nox_drawable_next_45A070(p)))
	case 4:
		return uintptr(unsafe.Pointer(C.sub_45A010(p)))
	}
	panic("unknown drawable state operation")
}
func PortTestDrawableMotion(op int, vp *noxrender.Viewport, dr *client.Drawable) int {
	if op == 0 {
		return drawableMotionDamped(vp, dr)
	}
	return drawableMotionTarget(dr)
}
