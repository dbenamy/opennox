//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"
import "unsafe"

func PortTestProtectionValidate(initial [][2]uint32, key, sum, id uint32, data []byte, size uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, 0x87654321, 0xffffffff, 0xffffffff, 1, 0x12345678, 123, func() uint32 {
		return uint32(C.sub_56FB00((*C.int)(unsafe.Pointer(unsafe.SliceData(data))), C.uint(size), C.int(id)))
	})
}
