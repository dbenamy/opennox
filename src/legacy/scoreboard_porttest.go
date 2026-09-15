//go:build porttest

package legacy

/*
#include "GAME2_1.h"
#include "client__gui__guirank.h"
*/
import "C"
import "unsafe"

// PortTestScoreboard invokes the original owners. The variadic formatter is
// exercised through its actual row-rendering callers, not a C test substitute.
func PortTestScoreboard(op int, a, b, c uintptr) uint32 {
	switch op {
	case 0:
		return uint32(uintptr(unsafe.Pointer(C.sub_46DC60(C.int(a), C.uchar(b), C.int(c)))))
	case 1:
		return uint32(C.nox_xxx_guiDrawRank_46E870())
	case 2:
		return uint32(uintptr(unsafe.Pointer(C.sub_46F030())))
	case 3:
		return uint32(C.sub_46F080(C.int(a), C.int(b)))
	case 4:
		return uint32(uintptr(unsafe.Pointer(C.sub_46F8F0(C.int(a), C.int(b)))))
	case 5:
		return uint32(uintptr(unsafe.Pointer(C.sub_46FB50(C.int(a), (*C.uint8_t)(unsafe.Pointer(b))))))
	case 6:
		return uint32(C.sub_46FC50())
	case 7:
		return uint32(C.sub_46FD80())
	case 8:
		return uint32(C.sub_46DB80())
	case 9:
		return uint32(C.sub_46DC00(C.int(a), C.uchar(b), C.int(c)))
	case 11:
		return uint32(uintptr(unsafe.Pointer(C.sub_46DCC0())))
	case 12:
		return uint32(C.sub_46E080(C.int(a)))
	case 13:
		return uint32(C.sub_46E130(C.int(a)))
	case 14:
		return uint32(uintptr(unsafe.Pointer(C.sub_46E170((*C.wchar2_t)(unsafe.Pointer(a))))))
	case 15:
		return uint32(C.sub_46E1E0(C.int(a)))
	case 16:
		return uint32(uintptr(unsafe.Pointer(C.sub_46E4E0())))
	case 17:
		return uint32(C.sub_46F060())
	case 18:
		return uint32(C.nox_xxx_Proc_46F070())
	case 19:
		C.sub_46FAE0()
		return 0
	case 20:
		return uint32(C.sub_46FE60(C.int(a)))
	case 21:
		return uint32(C.sub_46FEB0(C.uchar(a)))
	case 22:
		return uint32(C.sub_46FEE0())
	case 23:
		return uint32(C.sub_46FF70(C.int(a)))
	case 24:
		return uint32(C.sub_46FFD0())
	case 25:
		return uint32(C.sub_470580())
	case 26:
		C.sub_4705B0()
		return 0
	case 27:
		return uint32(C.sub_4705F0(C.char(a), C.char(b), C.short(c)))
	case 28:
		return uint32(C.sub_470650(C.char(a), C.short(b)))
	case 29:
		return uint32(C.sub_470680())
	case 30:
		return uint32(C.sub_4706A0())
	}
	panic("unknown scoreboard operation")
}
func PortTestScoreboardCallbacks() []unsafe.Pointer {
	return []unsafe.Pointer{C.sub_46F060, C.nox_xxx_Proc_46F070, C.sub_46F080}
}
