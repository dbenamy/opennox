package legacy

/*
#include "GAME1.h"
*/
import "C"

func Sub_4AD840() {
	serverPanelsGeneralRefresh()
}
func Sub_409E70(a1 int) {
	C.sub_409E70(C.int(a1))
}
func Sub_415960(a1 string) uint32 {
	return uint32(C.sub_415960(internWStr(a1)))
}
func Sub_415DA0(a1 string) uint32 {
	return uint32(C.sub_415DA0(internWStr(a1)))
}
