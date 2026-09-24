//go:build porttest

package legacy

/*
#include <stdint.h>
static unsigned short portTestFloatCW(void) { unsigned short cw; __asm__("fnstcw %0":"=m"(cw)); return cw; }
*/
import "C"

import "math"

func PortTestProtectionFloat(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		f := C.float(math.Float32frombits(value))
		switch mode {
		case 0:
			return uint32(portTestInvoke_sub_56F8C0(C.int(id), f))
		default:
			return uint32(portTestInvoke_sub_56FA40(C.int(id), f))
		}
	})
}

func PortTestProtectionFloatCW() uint16 { return uint16(C.portTestFloatCW()) }

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_sub_56F8C0(id C.int, value C.float) C.uint32_t {
	return C.uint32_t(updateProtectionFloat(int32(id), float32(value), false))
}

func portTestInvoke_sub_56FA40(id C.int, value C.float) C.uint32_t {
	return C.uint32_t(updateProtectionFloat(int32(id), float32(value), true))
}
