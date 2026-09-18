package legacy

/*
#include <stdint.h>
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func waypointFromRaw(a1 C.int) *server.Waypoint {
	return (*server.Waypoint)(unsafe.Pointer(uintptr(uint32(a1))))
}

func waypointRaw(wp *server.Waypoint) C.int {
	return C.int(uint32(uintptr(unsafe.Pointer(wp))))
}

//export nox_xxx_waypointNext_579870
func nox_xxx_waypointNext_579870(a1 C.int) C.int {
	if a1 == 0 {
		return 0
	}
	return waypointRaw(waypointFromRaw(a1).WpNext)
}

//export sub_5798A0
func sub_5798A0(a1 C.int) C.int {
	if a1 == 0 {
		return 0
	}
	return waypointRaw(waypointFromRaw(a1).WpNext)
}

//export sub_579E70
func sub_579E70() *C.uint32_t {
	wp, _ := alloc.New(server.Waypoint{})
	wp.Flags |= 0x1000000
	return (*C.uint32_t)(unsafe.Pointer(wp))
}

func waypointEnabledMask(wp *server.Waypoint, mask byte) bool {
	return wp.IsEnabled() && wp.HasFlag2Mask(mask)
}
