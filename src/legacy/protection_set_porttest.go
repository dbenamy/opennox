//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
// Convert the decompiled pointer-typed return to scalar bits inside C. This
// adapter also compiles after the declarations are corrected to uint32_t.
static uint32_t portTestProtectionSet(int mode, int id, uint32_t value) {
	switch (mode) {
	case 0: return (uintptr_t)sub_56F780(id, (int)value);
	case 1: return (uintptr_t)nox_xxx_playerResetProtectionCRC_56F7D0(id, (int)value);
	case 2: return (uintptr_t)sub_56F820(id, (unsigned char)value);
	case 3: return (uintptr_t)nox_xxx_protectPlayerHPMana_56F870(id, (unsigned short)value);
	default: return 0;
	}
}
*/
import "C"

func PortTestProtectionSet(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if mode == 4 {
			Nox_xxx_playerResetProtectionCRC_56F7D0(id, int(value))
			return 0
		}
		return uint32(C.portTestProtectionSet(C.int(mode), C.int(id), C.uint32_t(value)))
	})
}
