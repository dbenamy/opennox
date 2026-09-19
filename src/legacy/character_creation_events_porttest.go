//go:build porttest

package legacy

/*
#include "GAME3.h"
*/
import "C"
import "unsafe"

func PortTestCharacterColorEvent(root, child unsafe.Pointer, event int, point uint32) int {
	return int(C.sub_4A7330(C.int(uintptr(root)), C.int(event), (*C.int)(child), C.uint(point)))
}
func PortTestCharacterPaletteOutside(menu unsafe.Pointer, event int, point uint32) int {
	return int(C.sub_4A7270(C.int(uintptr(menu)), C.int(event), C.uint(point), 0))
}
