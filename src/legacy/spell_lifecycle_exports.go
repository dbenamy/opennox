package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4FD030
func sub_4FD030(a C.int, b C.short) C.ushort {
	return C.ushort(spellLifeRefundMana(objectFromInt(a), int16(b)))
}

//export nox_xxx_collide_4FDF90
func nox_xxx_collide_4FDF90(a, b C.int) { spellLifeCollide(objectFromInt(a), objectFromInt(b)) }

//export nox_xxx_spellGetPhoneme_4FE1C0
func nox_xxx_spellGetPhoneme_4FE1C0(a C.int, b C.char) C.int {
	return C.int(spellLifePhoneme(int32(a), int8(b)))
}

//export nox_xxx_spellByBookInsert_4FE340
func nox_xxx_spellByBookInsert_4FE340(a C.int, b *C.int, c, d, e C.int) C.int {
	return C.int(spellLifeInsertBook(objectFromInt(a), unsafe.Pointer(b), int32(c), int32(d), int32(e)))
}

//export sub_4FEA70
func sub_4FEA70(a C.int, b *C.float2) C.int {
	return C.int(spellLifeMoved(objectFromInt(a), (*types.Pointf)(unsafe.Pointer(b))))
}

//export nox_xxx_playerCancelSpells_4FEAE0
func nox_xxx_playerCancelSpells_4FEAE0(a *C.nox_object_t) C.int {
	return C.int(spellLifeCancelPlayer(asObjectS(a)))
}

//export nox_xxx_netStartDurationRaySpell_4FF130
func nox_xxx_netStartDurationRaySpell_4FF130(a C.int) *C.char {
	return (*C.char)(unsafe.Pointer(uintptr(spellLifeRayMessage((*server.DurSpell)(unsafe.Pointer(uintptr(a)))))))
}

//export sub_4FF2D0
func sub_4FF2D0(a, b C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(spellLifeFindDuration(int32(a), objectFromInt(b)))))
}

//export nox_xxx_testUnitBuffs_4FF350
func nox_xxx_testUnitBuffs_4FF350(a *C.nox_object_t, b C.char) C.int {
	return C.int(bool2int(spellLifeHasBuff(asObjectS(a), int32(int8(b)))))
}

//export nox_xxx_buffApplyTo_4FF380
func nox_xxx_buffApplyTo_4FF380(a *C.nox_object_t, b C.int, c C.short, d C.char) {
	spellLifeApplyBuff(asObjectS(a), int32(b), int16(c), int8(d))
}

//export nox_xxx_unitGetBuffTimer_4FF550
func nox_xxx_unitGetBuffTimer_4FF550(a *C.nox_object_t, b C.int) C.int {
	return C.int(spellLifeBuffTimer(asObjectS(a), int32(b)))
}

//export nox_xxx_buffGetPower_4FF570
func nox_xxx_buffGetPower_4FF570(a *C.nox_object_t, b C.int) C.char {
	return C.char(spellLifeBuffPower(asObjectS(a), int32(b)))
}

//export nox_xxx_unitClearBuffs_4FF580
func nox_xxx_unitClearBuffs_4FF580(a *C.nox_object_t) { spellLifeClearBuffs(asObjectS(a)) }

//export nox_xxx_spellBuffOff_4FF5B0
func nox_xxx_spellBuffOff_4FF5B0(a *C.nox_object_t, b C.int) C.int {
	return C.int(spellLifeBuffOff(asObjectS(a), int32(b)))
}
