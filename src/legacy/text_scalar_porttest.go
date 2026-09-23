//go:build porttest

package legacy

/*
#include "noxstring.h"
*/
import "C"

import "unsafe"

func PortTestTextCompareWide(a, b *uint16) int {
	return int(C._nox_wcsicmp((*C.wchar2_t)(unsafe.Pointer(a)), (*C.wchar2_t)(unsafe.Pointer(b))))
}
func PortTestTextCompareNarrow(a, b *byte) int {
	return int(C.nox_strcmpi((*C.char)(unsafe.Pointer(a)), (*C.char)(unsafe.Pointer(b))))
}
func PortTestTextDecimal(p *uint16) int32 {
	return int32(C.nox_wcstol((*C.wchar2_t)(unsafe.Pointer(p)), nil, 10))
}
func PortTestTextCopy(dst, src *uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(C.nox_wcscpy((*C.wchar2_t)(unsafe.Pointer(dst)), (*C.wchar2_t)(unsafe.Pointer(src)))))
}
