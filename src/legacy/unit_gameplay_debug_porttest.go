//go:build porttest

package legacy

/*
#include "GAME4_3.h"
static void port_unit_debug(int kind, int frame, const char* name, int code) {
 if (kind == 0) nox_ai_debug_printf_5341A0("%d: Lost sight of %s(#%d)\n",frame,name,code);
 else nox_ai_debug_printf_5341A0("%d: %s(#%d) FRUSTRATED\n",frame,name,code);
}
*/
import "C"

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestUnitDebug(kind int, frame, code uint32, name string) {
	p, free := alloc.CString(name)
	defer free()
	C.port_unit_debug(C.int(kind), C.int(frame), (*C.char)(unsafe.Pointer(p)), C.int(code))
}
