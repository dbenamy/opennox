//go:build porttest

package legacy

import "unsafe"

// Stable, distinct identities used only by fixture snapshot normalization.
var portTestFixtureKeySlots [19]byte

func portTestFixtureKey(name string) unsafe.Pointer {
	switch name {
	case "nox_xxx_elevatorAud_53B490":
		return unsafe.Pointer(&portTestFixtureKeySlots[0])
	case "nox_xxx_elevatorFn_53B750":
		return unsafe.Pointer(&portTestFixtureKeySlots[1])
	case "nox_xxx_fnElevatorShaft_53B410":
		return unsafe.Pointer(&portTestFixtureKeySlots[2])
	case "nox_xxx_fnPentagramTeleport_53C060":
		return unsafe.Pointer(&portTestFixtureKeySlots[3])
	case "nox_xxx_pickupFlagCtf_4EA490":
		return unsafe.Pointer(&portTestFixtureKeySlots[4])
	case "sub_417F50":
		return unsafe.Pointer(&portTestFixtureKeySlots[5])
	case "sub_4E9A30":
		return unsafe.Pointer(&portTestFixtureKeySlots[6])
	case "sub_4EA7A0":
		return unsafe.Pointer(&portTestFixtureKeySlots[7])
	case "sub_4EA800":
		return unsafe.Pointer(&portTestFixtureKeySlots[8])
	case "sub_4EB250":
		return unsafe.Pointer(&portTestFixtureKeySlots[9])
	case "sub_4EB340":
		return unsafe.Pointer(&portTestFixtureKeySlots[10])
	case "sub_4EB3E0":
		return unsafe.Pointer(&portTestFixtureKeySlots[11])
	case "sub_4EB9B0":
		return unsafe.Pointer(&portTestFixtureKeySlots[12])
	case "sub_4ECBD0":
		return unsafe.Pointer(&portTestFixtureKeySlots[13])
	case "sub_4ECC00":
		return unsafe.Pointer(&portTestFixtureKeySlots[14])
	case "sub_53C140":
		return unsafe.Pointer(&portTestFixtureKeySlots[15])
	case "sub_53C240":
		return unsafe.Pointer(&portTestFixtureKeySlots[16])
	case "sub_548830":
		return unsafe.Pointer(&portTestFixtureKeySlots[17])
	case "sub_548860":
		return unsafe.Pointer(&portTestFixtureKeySlots[18])
	default:
		return nil
	}
}
