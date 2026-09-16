//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3.h"
extern uint32_t dword_5d4594_527660;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"unsafe"
)

// Calls the real installed map reader; no replacement stream or drawable factory.
func PortTestMapDrawableBase(typ int, outer, inner int16, old bool, consumed *uint32) *client.Drawable {
	var p C.int
	if old {
		p = C.sub_4ABDA0(C.int(typ), C.short(inner), C.short(outer), (*C.uint32_t)(unsafe.Pointer(consumed)))
	} else {
		p = C.nox_xxx_spriteLoadFromMap_4AC020(C.int(typ), C.short(outer), (*C.uint32_t)(unsafe.Pointer(consumed)))
	}
	return (*client.Drawable)(unsafe.Pointer(uintptr(uint32(p))))
}
func PortTestMapDrawableRecord(op, typ int) int {
	switch op {
	case 0:
		return int(C.nox_xxx_clientLoadSomeObject_4AC6E0(C.ushort(typ)))
	case 1:
		return int(C.sub_4AC7B0(C.int(typ)))
	case 2:
		return int(C.nox_xxx_colorLightClientLoad_4AC980(C.int(typ)))
	case 3:
		return int(C.nox_xxx_cliLoadTeamBase_4ACE00(C.int(typ)))
	case 4:
		return int(C.sub_4ACEF0(C.int(typ)))
	case 5:
		return int(C.sub_4AD040(C.int(typ)))
	case 6:
		return int(C.nox_client_mapSpecialRWObjectData_4AC610())
	}
	panic(op)
}

func PortTestMapDrawableTeamWord() *uint32 { return (*uint32)(unsafe.Pointer(&C.dword_5d4594_527660)) }
