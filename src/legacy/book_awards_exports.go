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

func nox_xxx_guide_427010(name *C.char) C.int {
	return C.int(bookGuideID(alloc.GoString((*byte)(unsafe.Pointer(name)))))
}

func nox_xxx_guiCreatureGetName_427240(id C.int) C.int {
	return C.int(bookGuideCreatureName(int32(id)))
}
