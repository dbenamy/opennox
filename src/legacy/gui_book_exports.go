package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

func nox_xxx_netSpellRewardCli_45CFE0(id, rank, notify, autoAdd C.int) {
	bookSpellReward(int(id), int(rank), int(notify), int(autoAdd))
}

func nox_xxx_netGuideRewardCli_45D140(id, notify C.int) { bookGuideReward(int(id), int(notify)) }

func sub_45D320(id C.int) C.int { return C.int(bookRemoveSpell(int(id))) }

func sub_45D400(id C.int) C.int { return C.int(bookRemoveGuide(int(id))) }

func nox_xxx_clientQuestDisableAbility_45D4A0(id C.int) *C.char {
	return (*C.char)(unsafe.Pointer(bookRemoveAbility(int(id))))
}

//export nox_xxx_bookClickSpell_45B1F0
func nox_xxx_bookClickSpell_45B1F0() C.int { return C.int(bookPageComplete(false)) }

//export nox_xxx_bookClickCreature_45B200
func nox_xxx_bookClickCreature_45B200() C.int { return C.int(bookPageComplete(true)) }

//export nox_xxx_book_45CF00
func nox_xxx_book_45CF00(w *C.uint32_t) C.int {
	return C.int(bookIconTooltip(AsWindowP(unsafe.Pointer(w))))
}
