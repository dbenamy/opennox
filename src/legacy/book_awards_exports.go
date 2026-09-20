package legacy

/*
#include "GAME1_1.h"
#include "GAME4_3.h"
#include "server__ability__ability.h"
#include "server__magic__plyrspel.h"
typedef const char book_award_const_char;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

//export nox_xxx_guide_427010
func nox_xxx_guide_427010(name *C.book_award_const_char) C.int {
	return C.int(bookGuideID(alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

//export nox_xxx_guiCreatureGetName_427240
func nox_xxx_guiCreatureGetName_427240(id C.int) C.int {
	return C.int(bookGuideCreatureName(int32(id)))
}

//export sub_53F930
func sub_53F930(a, b C.int) C.int { return C.int(bookUseGuide(objectFromInt(a), objectFromInt(b))) }

//export nox_xxx_useSpellReward_53F9E0
func nox_xxx_useSpellReward_53F9E0(a, b C.int) C.int {
	return C.int(bookUseSpell(objectFromInt(a), objectFromInt(b)))
}

//export nox_xxx_useAbilityReward_53FAE0
func nox_xxx_useAbilityReward_53FAE0(a, b C.int) C.int {
	return C.int(bookUseAbility(objectFromInt(a), objectFromInt(b)))
}
