//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_1.h"
*/
import "C"
import "unsafe"

// PortTestTradeUIBoundary checks the retained C call boundaries, not another
// implementation. Ordinary fixture dispatch calls the Go owners directly.
func PortTestTradeUIBoundary(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	switch op {
	case 7:
		return uint32(C.nox_gui_itemAmountDialog_4C0430((*C.wchar2_t)(unsafe.Pointer(a[0])), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), unsafe.Pointer(a[5]), C.int(a[6]), C.int(a[7]), unsafe.Pointer(a[8]), unsafe.Pointer(a[9])))
	case 27:
		return uint32(nox_xxx_tradeClientAddItem_4C1790(C.int(a[0])))
	default:
		panic("trade UI C boundary")
	}
}
