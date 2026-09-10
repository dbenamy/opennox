//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"
import "unsafe"

// PortTestProtectionChecksum exercises the actual C ABI, including its signed
// return value. No historical C implementation is needed for these ABI checks.
func PortTestProtectionChecksum(data []byte) uint32 {
	p := (*C.int)(unsafe.Pointer(unsafe.SliceData(data)))
	return uint32(C.nox_xxx_protectionStringCRCLen_56FAE0(p, C.uint(len(data))))
}

func PortTestProtectionNull(n uint32) uint32 {
	return uint32(C.nox_xxx_protectionStringCRCLen_56FAE0(nil, C.uint(n)))
}
