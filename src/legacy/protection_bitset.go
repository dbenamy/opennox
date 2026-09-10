package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
)

//export sub_56FCB0
func sub_56FCB0(index, enabled C.int) C.int {
	return C.int(protection.Bit(int32(index), int32(enabled)))
}

//export nox_xxx_playerAwardSpellProtectionCRC_56FCE0
func nox_xxx_playerAwardSpellProtectionCRC_56FCE0(id, index, enabled C.int) C.int {
	if id < 657757279 {
		return id
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return id
	}
	old := r.Value
	r.Value = ((old ^ key) | protection.Bit(int32(index), int32(enabled))) ^ key
	C.dword_5d4594_2516328 ^= C.uint32_t(old ^ r.Value)
	return C.int(r.Value)
}

//export nox_xxx_playerApplyProtectionCRC_56FD50
func nox_xxx_playerApplyProtectionCRC_56FD50(id C.int, data unsafe.Pointer, count C.int) C.int {
	if id < 657757279 {
		return 0
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	var flags uint32
	if count > 1 {
		flags = protection.Flags(unsafe.Slice((*int32)(data), int(count)))
	}
	if flags^key == r.Value {
		return 1
	}
	return 0
}
