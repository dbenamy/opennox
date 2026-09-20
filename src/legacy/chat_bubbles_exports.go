package legacy

/*
#include "GAME2_3.h"
*/
import "C"

func sub_48E8E0(code C.int) { chatBubbleRemove(uint32(code)) }

func sub_48E940() { chatBubbleClear() }
