//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"math"
	"unsafe"
)

// PortTestBriefingWindow exercises actual window lifecycle and transition owners.
func PortTestBriefingWindow(op int, a, b, c, d uintptr) uint64 {
	switch op {
	case 0:
		return uint64(uintptr(unsafe.Pointer(briefingCreateWindow())))
	case 1:
		return uint64(uint32(briefingPlayVoice()))
	case 2:
		return uint64(uint32(briefingInput((*gui.Window)(unsafe.Pointer(a)), int(b), c, d)))
	case 3:
		return uint64(uint32(briefingBackgroundEvent((*gui.Window)(unsafe.Pointer(a)), int(b), c, d)))
	case 4:
		return uint64(uint32(briefingDrawWindow((*gui.Window)(unsafe.Pointer(a)), (*gui.WindowData)(unsafe.Pointer(b)))))
	case 5:
		return math.Float64bits(briefingScrollSpeed())
	case 6:
		return uint64(uint32(briefingCompletedDraw(nil, nil)))
	case 7:
		return uint64(uint32(briefingShow(int(a), int(b), byte(c))))
	case 8:
		return uint64(uint32(briefingDestroy()))
	}

	return 0
}
func PortTestBriefingWindowCallbacks() []unsafe.Pointer { return nil }
