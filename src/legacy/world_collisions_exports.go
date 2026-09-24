package legacy

/*
#include <stdint.h>
#include <stdbool.h>
*/
import "C"

import (
	"unsafe"
)

//export sub_4E8E50
func sub_4E8E50() *C.uchar { return (*C.uchar)(unsafe.Pointer(worldQuestPending())) }

//export sub_4E8E60
func sub_4E8E60() C.int { return C.int(worldQuestCountdown()) }

//export nox_server_questMaybeWarp_4E8F60
func nox_server_questMaybeWarp_4E8F60() C.bool { return C.bool(worldQuestMaybeWarp()) }

//export sub_4E9010
func sub_4E9010() C.int { return C.int(bool2int(worldQuestExitReady())) }
