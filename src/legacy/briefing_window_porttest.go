//go:build porttest

package legacy

/*
#include "GAME2.h"
*/
import "C"
import (
	"math"
	"unsafe"
)

// PortTestBriefingWindow exercises actual window lifecycle and transition owners.
func PortTestBriefingWindow(op int, a, b, c, d uintptr) uint64 {
	switch op {
	case 0:
		return uint64(uintptr(unsafe.Pointer(C.sub_44E560())))
	case 1:
		return uint64(uint32(C.nox_xxx_playGMCAPsmth_44E3E0()))
	case 2:
		return uint64(uint32(C.nox_client_wndQuestBriefProc_44E630(C.int(a), C.int(b), C.int(c), C.int(d))))
	case 3:
		return uint64(uint32(C.nox_xxx_wndProc_44E6E0(C.int(a), C.int(b), C.int(c), C.int(d))))
	case 4:
		return uint64(uint32(C.sub_44E6F0((*C.uint32_t)(unsafe.Pointer(a)), C.int(b))))
	case 5:
		return math.Float64bits(float64(C.sub_44E8B0()))
	case 6:
		return uint64(uint32(C.sub_44E8D0()))
	case 7:
		return uint64(uint32(C.nox_client_lockScreenBriefing_450160(C.int(a), C.int(b), C.char(c))))
	case 8:
		return uint64(uint32(C.sub_4505E0()))
	}
	return 0
}
func PortTestBriefingWindowCallbacks() []unsafe.Pointer { return []unsafe.Pointer{C.sub_44E8D0} }
