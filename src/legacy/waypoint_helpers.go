package legacy

/*
#include <stdint.h>
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

func waypointFromRaw(a1 C.int) *server.Waypoint {
	return (*server.Waypoint)(unsafe.Pointer(uintptr(uint32(a1))))
}

func waypointRaw(wp *server.Waypoint) C.int {
	return C.int(uint32(uintptr(unsafe.Pointer(wp))))
}

func waypointEnabledMask(wp *server.Waypoint, mask byte) bool {
	return wp.IsEnabled() && wp.HasFlag2Mask(mask)
}
