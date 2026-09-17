//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4.h"
#include "GAME4_1.h"
extern unsigned int dword_587000_230092;
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Primitive entry dispatch; all behavior belongs to the original production code.
func PortTestObjectReports(op int, a, b *server.Object, player, x, y int, pos *types.Pointf) int {
	ai, bi := C.int(uintptr(unsafe.Pointer(a))), C.int(uintptr(unsafe.Pointer(b)))
	switch op {
	case 0:
		return int(C.sub_518770())
	case 1:
		return int(C.nox_xxx_netSendPhantomPlrMb_5187E0(C.int(player), bi))
	case 2:
		return int(C.nox_xxx_netSendSimpleObj_5188A0(C.int(player), bi))
	case 3:
		return int(C.nox_xxx_netSendComplexObject_518960(C.int(player), (*C.uint)(unsafe.Pointer(b)), C.int(x)))
	case 4:
		return int(C.nox_xxx_netSpriteUpdate_518AE0(ai, C.int(player), (*C.uint)(unsafe.Pointer(b))))
	case 5:
		return int(C.nox_xxx_netPlayerObjSend_518C30(asObjectC(a), asObjectC(b), C.int(x), C.int(y)))
	case 6:
		return int(C.nox_xxx_netSendObjects2Plr_519410(asObjectC(a), asObjectC(b)))
	case 7:
		return int(C.sub_519710(a.UpdateData))
	case 8:
		return int(C.sub_501C00((*C.float)(unsafe.Pointer(pos)), asObjectC(b)))
	case 9:
		C.nox_xxx_netUpdateRemotePlr_501CA0(asObjectC(a))
		return 0
	default:
		panic("object report operation")
	}
}
func PortTestObjectReportsDirectionThreshold() *int32 {
	return (*int32)(unsafe.Pointer(&C.dword_587000_230092))
}
