//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
#include "server__gamemech__explevel.h"
static uint32_t portTestProtectionAdd(int mode, int id, uint32_t value) {
	switch (mode) {
	case 0: return (uintptr_t)sub_56F920(id, (int)value);
	case 1: return (uintptr_t)nox_xxx_protectMana_56F9E0(id, (short)value);
	case 2: return (uintptr_t)sub_56F980(id, (unsigned char)value);
	default: return 0;
	}
}
*/
import "C"

func PortTestProtectionAdd(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if mode == 3 {
			Nox_xxx_protectMana_56F9E0(int(id), int16(value))
			return 0
		}
		return uint32(C.portTestProtectionAdd(C.int(mode), C.int(id), C.uint32_t(value)))
	})
}
