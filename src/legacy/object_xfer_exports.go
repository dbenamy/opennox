package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_xfer_4F3E30
func nox_xxx_xfer_4F3E30(version C.ushort, u *C.nox_object_t, count C.int) C.int {
	return C.int(objectXferInventory(uint16(version), asObjectS(u), int32(count)))
}

//export nox_xxx_servMapLoadPlaceObj_4F3F50
func nox_xxx_servMapLoadPlaceObj_4F3F50(u *C.nox_object_t, owner C.int, offset unsafe.Pointer) C.int {
	return C.int(objectXferPlace(asObjectS(u), unsafe.Pointer(uintptr(uint32(owner))), offset))
}

//export nox_xxx_mapReadWriteObjData_4F4530
func nox_xxx_mapReadWriteObjData_4F4530(u *C.nox_object_t, version C.int) C.int {
	return C.int(objectXferCommon(asObjectS(u), int(version)))
}

//export nox_xxx_XFerSpellPagePedistal_4F4A20
func nox_xxx_XFerSpellPagePedistal_4F4A20(u C.int) C.int {
	return C.int(objectXferPedestal(objectFromInt(u)))
}

//export nox_xxx_XFerReadable_4F4AB0
func nox_xxx_XFerReadable_4F4AB0(u C.int) C.int { return C.int(objectXferReadable(objectFromInt(u))) }

//export nox_xxx_XFerExit_4F4B90
func nox_xxx_XFerExit_4F4B90(u C.int) C.int { return C.int(objectXferExit(objectFromInt(u))) }

//export nox_xxx_XFerDoor_4F4CB0
func nox_xxx_XFerDoor_4F4CB0(u C.int) C.int { return C.int(objectXferDoor(objectFromInt(u))) }

// The registered callback always receives an object pointer. The old C float
// declaration reinterpreted its first stack slot as that pointer before reusing
// the parameter as numeric scratch. Keep the real 386 argument, without a float
// conversion or a fabricated second parameter.
//
//export nox_xxx_unitTriggerXfer_4F4E50
func nox_xxx_unitTriggerXfer_4F4E50(u *C.nox_object_t) C.int {
	return C.int(objectXferTrigger(asObjectS(u)))
}

//export nox_xxx_XFerHole_4F51D0
func nox_xxx_XFerHole_4F51D0(u C.int) C.int { return C.int(objectXferHole(objectFromInt(u))) }

//export nox_xxx_XFerTransporter_4F5300
func nox_xxx_XFerTransporter_4F5300(u C.int) C.int {
	return C.int(objectXferTransporter(objectFromInt(u)))
}

//export nox_xxx_XFerElevator_4F53D0
func nox_xxx_XFerElevator_4F53D0(u C.int) C.int { return C.int(objectXferElevator(objectFromInt(u))) }

//export nox_xxx_XFerElevatorShaft_4F54A0
func nox_xxx_XFerElevatorShaft_4F54A0(u C.int) C.int { return C.int(objectXferShaft(objectFromInt(u))) }

//export sub_4F5540
func sub_4F5540(p C.int) C.int {
	return C.int(objectXferLegacyScript(unsafe.Pointer(uintptr(uint32(p)))))
}

//export nox_xxx_XFerMover_4F5730
func nox_xxx_XFerMover_4F5730(u C.int) C.int { return C.int(objectXferMover(objectFromInt(u))) }

//export nox_xxx_XFerGlyph_4F5890
func nox_xxx_XFerGlyph_4F5890(u C.int) C.int { return C.int(objectXferGlyph(objectFromInt(u))) }

//export nox_xxx_XFerInvLight_4F5AA0
func nox_xxx_XFerInvLight_4F5AA0(u *C.int) C.int {
	return C.int(objectXferLight((*server.Object)(unsafe.Pointer(u))))
}

//export nox_xxx_XFerSentry_4F5E50
func nox_xxx_XFerSentry_4F5E50(u C.int) C.int { return C.int(objectXferSentry(objectFromInt(u))) }
