package legacy

/*
#include "GAME2_3.h"
*/
import "C"

//export sub_48E8E0
func sub_48E8E0(code C.int) { chatBubbleRemove(uint32(code)) }

//export sub_48E940
func sub_48E940() { chatBubbleClear() }
