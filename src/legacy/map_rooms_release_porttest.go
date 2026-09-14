//go:build porttest

package legacy

import "unsafe"

// Observe allocation lifetime at the shared native allocator boundary. This
// also covers remaining C callers of the room-free exports.
var mapRoomTestRelease func(unsafe.Pointer)

func mapRoomBeforeRelease(p unsafe.Pointer) {
	if mapRoomTestRelease != nil {
		mapRoomTestRelease(p)
	}
}
