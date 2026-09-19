//go:build porttest

package legacy

/*
#include "GAME4_3.h"
*/
import "C"
import "unsafe"

// PortTestResourceParser invokes the original registered parser on owned buffers.
func PortTestResourceParser(name string, input, data unsafe.Pointer) int {
	s := (*C.char)(input)
	d := C.int(uintptr(data))
	switch name {
	case "wandcast":
		return int(uintptr(unsafe.Pointer(C.sub_5361B0(s, d))))
	case "wand":
		return int(uintptr(unsafe.Pointer(C.sub_536260(s, d))))
	case "skull":
		return int(C.sub_5364E0(s, d))
	case "push":
		return int(C.sub_536550(s, (*C.uint32_t)(data)))
	case "triple":
		return int(C.sub_536580(s, d))
	case "trigger":
		return int(C.sub_5365B0(s, d))
	case "lifetime":
		return int(C.sub_536600(s, d))
	case "spawn":
		return int(C.sub_536B40(s, d))
	case "projectile":
		return int(C.sub_536D80(s, d))
	case "audio":
		return int(C.sub_536DA0(s, (*C.int)(data)))
	case "spark":
		return int(C.sub_536DE0(s, (*C.uchar)(data)))
	case "damage":
		return int(C.nox_xxx_collideDamageLoad_536E10(s, d))
	case "mana":
		return int(C.sub_536E50(s, (*C.uchar)(data)))
	case "arrow":
		return int(C.sub_536E80(s, (*C.int)(data)))
	}
	panic(name)
}
