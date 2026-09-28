//go:build porttest

package legacy

import "unsafe"

// PortTestProtectionChecksum exercises the actual C ABI, including its signed
// return value. No historical C implementation is needed for these ABI checks.
func PortTestProtectionChecksum(data []byte) uint32 {
	p := (*int32)(unsafe.Pointer(unsafe.SliceData(data)))
	return uint32(nox_xxx_protectionStringCRCLen_56FAE0(p, uint32(len(data))))
}

func PortTestProtectionNull(n uint32) uint32 {
	return uint32(nox_xxx_protectionStringCRCLen_56FAE0(nil, uint32(n)))
}
