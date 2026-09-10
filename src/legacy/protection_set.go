package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"

import "github.com/opennox/opennox/v1/internal/protection"

func setProtectionRecord(id int32, value uint32) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	old := r.Value
	r.Value = value ^ key
	C.dword_5d4594_2516328 ^= C.uint32_t(old ^ r.Value)
	return uint32(nox_xxx_protectData_56F5C0())
}

//export sub_56F780
func sub_56F780(id, value C.int) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

//export nox_xxx_playerResetProtectionCRC_56F7D0
func nox_xxx_playerResetProtectionCRC_56F7D0(id, value C.int) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

//export sub_56F820
func sub_56F820(id C.int, value C.uchar) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}

//export nox_xxx_protectPlayerHPMana_56F870
func nox_xxx_protectPlayerHPMana_56F870(id C.int, value C.ushort) C.uint32_t {
	return C.uint32_t(setProtectionRecord(int32(id), uint32(value)))
}
