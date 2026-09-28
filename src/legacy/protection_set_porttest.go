//go:build porttest

package legacy

func PortTestProtectionSet(initial [][2]uint32, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed uint32, seed int, mode int, id, value uint32) PortTestRekeySnapshot {
	return portTestRekeyOperation(initial, key, sum, sequence, swapCount, rekeyCount, frame, floatSeed, seed, func() uint32 {
		if mode == 4 {
			Nox_xxx_playerResetProtectionCRC_56F7D0(id, int(value))
			return 0
		}
		switch int32(mode) {
		case 0:
			return uint32(sub_56F780(int32(id), int32(value)))
		case 1:
			return uint32(nox_xxx_playerResetProtectionCRC_56F7D0(int32(id), int32(value)))
		case 2:
			return uint32(sub_56F820(int32(id), uint8(value)))
		case 3:
			return uint32(nox_xxx_protectPlayerHPMana_56F870(int32(id), uint16(value)))
		default:
			return 0
		}
	})
}
