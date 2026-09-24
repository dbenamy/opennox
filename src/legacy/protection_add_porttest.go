//go:build porttest

package legacy

func PortTestProtectionAdd(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if mode == 3 {
			Nox_xxx_protectMana_56F9E0(int(id), int16(value))
			return 0
		}
		switch mode {
		case 0:
			return uint32(portTestInvoke_sub_56F920(int32(id), int32(int32(value))))
		case 1:
			return uint32(portTestInvoke_nox_xxx_protectMana_56F9E0(int32(id), int16(value)))
		case 2:
			return uint32(portTestInvoke_sub_56F980(int32(id), uint8(value)))
		default:
			return 0
		}
	})
}

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_nox_xxx_protectMana_56F9E0(id int32, delta int16) uint32 {
	return uint32(addProtectionRecord(int32(id), uint32(int32(delta))))
}

func portTestInvoke_sub_56F920(id, delta int32) uint32 {
	return uint32(addProtectionRecord(int32(id), uint32(delta)))
}

func portTestInvoke_sub_56F980(id int32, delta uint8) uint32 {
	return uint32(addProtectionRecord(int32(id), uint32(delta)))
}
