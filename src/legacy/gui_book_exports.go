package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_client_toggleSpellbook_45AC70
func nox_client_toggleSpellbook_45AC70() { bookToggle() }

//export nox_xxx_bookHideMB_45ACA0
func nox_xxx_bookHideMB_45ACA0(reset C.int) C.int { return C.int(bookHide(int(reset))) }

//export sub_45CFC0
func sub_45CFC0() C.int { return C.int(bookShown()) }

func nox_xxx_netSpellRewardCli_45CFE0(id, rank, notify, autoAdd C.int) {
	bookSpellReward(int(id), int(rank), int(notify), int(autoAdd))
}

//export nox_xxx_netGuideRewardCli_45D140
func nox_xxx_netGuideRewardCli_45D140(id, notify C.int) { bookGuideReward(int(id), int(notify)) }

//export nox_xxx_bookSetForward_45D200
func nox_xxx_bookSetForward_45D200(kind *C.int, id C.int, pos *C.int2) *C.int {
	return (*C.int)(unsafe.Pointer(bookSetForward(uintptr(unsafe.Pointer(kind)), int(id), AsPoint(unsafe.Pointer(pos)))))
}

//export nox_xxx_abilityReward_45D290
func nox_xxx_abilityReward_45D290(id C.int, notify *C.char, autoAdd C.int) {
	bookAbilityReward(int(id), uintptr(unsafe.Pointer(notify)), int(autoAdd))
}

//export sub_45D320
func sub_45D320(id C.int) C.int { return C.int(bookRemoveSpell(int(id))) }

//export sub_45D400
func sub_45D400(id C.int) C.int { return C.int(bookRemoveGuide(int(id))) }

//export nox_xxx_clientQuestDisableAbility_45D4A0
func nox_xxx_clientQuestDisableAbility_45D4A0(id C.int) *C.char {
	return (*C.char)(unsafe.Pointer(bookRemoveAbility(int(id))))
}

//export sub_45D500
func sub_45D500(show C.int) C.int { return C.int(bookTemporaryShow(int(show))) }

//export nox_xxx_bookFillAll_45D570
func nox_xxx_bookFillAll_45D570(kind, id C.int) { bookAdd(int(kind), int(id)) }

//export sub_45D9B0
func sub_45D9B0() C.int { return C.int(*bookWord(1047520)) }

//export sub_45D870
func sub_45D870() { bookFinishAddition() }

//export nox_xxx_bookClickSpell_45B1F0
func nox_xxx_bookClickSpell_45B1F0() C.int { return C.int(bookPageComplete(false)) }

//export nox_xxx_bookClickCreature_45B200
func nox_xxx_bookClickCreature_45B200() C.int { return C.int(bookPageComplete(true)) }

//export nox_xxx_book_45CF00
func nox_xxx_book_45CF00(w *C.uint32_t) C.int {
	return C.int(bookIconTooltip(AsWindowP(unsafe.Pointer(w))))
}
