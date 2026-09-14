package legacy

/*
#include "GAME3_3.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4F2590
func sub_4F2590(ptr C.int) C.int {
	return C.int(bool2int(questEligibilityItem((*server.Object)(unsafe.Pointer(uintptr(uint32(ptr)))))))
}

//export sub_4F2C30
func sub_4F2C30(ptr C.int) C.int {
	return C.int(bool2int(questEligibilityInventoryLimit((*server.Object)(unsafe.Pointer(uintptr(uint32(ptr)))))))
}

//export nox_xxx_spell_4F2E70
func nox_xxx_spell_4F2E70(id C.int) C.int { return C.int(bool2int(questEligibilitySpell(uint32(id)))) }

//export sub_4F2EF0
func sub_4F2EF0(id C.int) C.int { return C.int(bool2int(questEligibilityBeast(uint32(id)))) }
