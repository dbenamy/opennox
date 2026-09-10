package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"

import "github.com/opennox/opennox/v1/internal/protection"

func addProtectionRecord(id int32, delta uint32) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	old := r.Value
	value := (old ^ key) + delta
	r.Value = value ^ key
	C.dword_5d4594_2516328 ^= C.uint32_t(old ^ r.Value)
	return uint32(nox_xxx_protectData_56F5C0())
}

//export sub_56F920
func sub_56F920(id, delta C.int) C.uint32_t {
	return C.uint32_t(addProtectionRecord(int32(id), uint32(delta)))
}

//export nox_xxx_protectMana_56F9E0
func nox_xxx_protectMana_56F9E0(id C.int, delta C.short) C.uint32_t {
	return C.uint32_t(addProtectionRecord(int32(id), uint32(int32(delta))))
}

//export sub_56F980
func sub_56F980(id C.int, delta C.uchar) C.uint32_t {
	return C.uint32_t(addProtectionRecord(int32(id), uint32(delta)))
}
