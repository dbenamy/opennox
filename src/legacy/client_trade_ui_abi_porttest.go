//go:build porttest

package legacy

import "unsafe"

// PortTestTradeUIBoundary preserves the former C boundary inputs through native owners.
func PortTestTradeUIBoundary(op int, args ...uintptr) uint32 {
	var a [10]uintptr
	copy(a[:], args)
	switch op {
	case 7:
		return uint32(uiAmountShow((*uint16)(unsafe.Pointer(a[0])), int(int32(a[1])), int(int32(a[2])), uint32(int32(a[3])), uint32(int32(a[4])), unsafe.Pointer(a[5]), uint32(int32(a[6])), uint32(int32(a[7])), unsafe.Pointer(a[8]), unsafe.Pointer(a[9])))
	case 27:
		return uiTradeAdd(unsafe.Pointer(uintptr(uint32(a[0]))))
	default:
		panic("trade UI native boundary")
	}
}
