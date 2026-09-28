//go:build porttest

package legacy

import (
	"math"

	"github.com/opennox/opennox/v1/legacy/common/fpenv"
)

func PortTestProtectionFloat(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		f := float32(math.Float32frombits(value))
		switch mode {
		case 0:
			return uint32(portTestInvoke_sub_56F8C0(int32(id), f))
		default:
			return uint32(portTestInvoke_sub_56FA40(int32(id), f))
		}
	})
}

func PortTestProtectionFloatCW() uint16 { return uint16(fpenv.Control()) }

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_sub_56F8C0(id int32, value float32) uint32 {
	return uint32(updateProtectionFloat(int32(id), float32(value), false))
}

func portTestInvoke_sub_56FA40(id int32, value float32) uint32 {
	return uint32(updateProtectionFloat(int32(id), float32(value), true))
}
