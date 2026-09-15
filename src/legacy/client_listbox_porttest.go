//go:build porttest

package legacy

/*
#include "GAME3.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func PortTestListboxScrollIndex(d *gui.ScrollListBoxData) int {
	return int(C.sub_4A4800(C.int(uintptr(unsafe.Pointer(d)))))
}
