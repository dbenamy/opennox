package legacy

/*
#include "GAME1_1.h"
*/
import "C"
import "unsafe"

func nox_xxx_plrLoad_41A480(path *C.char) C.int { return C.int(playerFileClientLoad(GoString(path))) }

//export sub_41C280
func sub_41C280(_ unsafe.Pointer) C.int { return C.int(playerFileGUI()) }

//export nox_xxx_parseFileInfoData_41C3B0
func nox_xxx_parseFileInfoData_41C3B0(_ C.int) C.int { return C.int(playerFileMetadata()) }

//export sub_41C780
func sub_41C780(_ C.int) C.int { return C.int(playerFileMusic()) }

//export nox_xxx_netSavePlayer_41CE00
func nox_xxx_netSavePlayer_41CE00() C.int { return C.int(playerFileSaveRequest()) }

//export sub_41CEE0
func sub_41CEE0(info unsafe.Pointer, all C.int) C.int {
	return C.int(playerFileClientWrite(info, int(all)))
}
