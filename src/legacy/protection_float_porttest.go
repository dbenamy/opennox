//go:build porttest

package legacy

/*
#include <stdint.h>
#include <string.h>
#include "GAME5_2.h"
static unsigned short portTestFloatCW(void) { unsigned short cw; __asm__("fnstcw %0":"=m"(cw)); return cw; }
static uint32_t portTestProtectionFloat(int mode, int id, uint32_t bits) {
 float value;
 memcpy(&value, &bits, sizeof(value));
 if (mode == 0) return (uintptr_t)sub_56F8C0(id, value);
 return (uintptr_t)sub_56FA40(id, value);
}
*/
import "C"

func PortTestProtectionFloat(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		return uint32(C.portTestProtectionFloat(C.int(mode), C.int(id), C.uint32_t(value)))
	})
}

func PortTestProtectionFloatCW() uint16 { return uint16(C.portTestFloatCW()) }
