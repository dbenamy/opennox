//go:build porttest

package legacy

import "unsafe"

func PortTestProtectionValidate(initial [][2]uint32, key, sum, id uint32, data []byte, size uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, 0x87654321, 0xffffffff, 0xffffffff, 1, 0x12345678, 123, func() uint32 {
		return uint32(portTestInvoke_sub_56FB00((*int32)(unsafe.Pointer(unsafe.SliceData(data))), uint32(size), int32(id)))
	})
}

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_sub_56FB00(data *int32, size uint32, id int32) int32 {
	return int32(protectionValidateString(unsafe.Pointer(data), uint32(size), int32(id)))
}
