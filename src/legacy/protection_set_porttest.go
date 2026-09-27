//go:build porttest

package legacy

/*
#include <stdint.h>
*/
import "C"

func PortTestProtectionSet(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if mode == 4 {
			Nox_xxx_playerResetProtectionCRC_56F7D0(id, int(value))
			return 0
		}
		switch int32(mode) {
		case 0:
			return uint32(sub_56F780(C.int(id), C.int(value)))
		case 1:
			return uint32(nox_xxx_playerResetProtectionCRC_56F7D0(C.int(id), C.int(value)))
		case 2:
			return uint32(sub_56F820(C.int(id), C.uchar(value)))
		case 3:
			return uint32(nox_xxx_protectPlayerHPMana_56F870(C.int(id), C.ushort(value)))
		default:
			return 0
		}
	})
}
