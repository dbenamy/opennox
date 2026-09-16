//go:build porttest

package legacy

/*
#include "GAME2_2.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// The only production caller discards the old scratch/pointer return.
func PortTestWallEdge(img noxrender.ImageHandle, pos image.Point, first, second *[3]uint32, bottom, left, right, crop, flags int) {
	C.nox_xxx_edgeDraw_480EF0(C.int(uintptr(img)), C.int(pos.X), C.int(pos.Y), (*C.int)(unsafe.Pointer(first)), (*C.int)(unsafe.Pointer(second)), C.int(bottom), C.int(left), C.int(right), C.int(crop), C.int(flags))
}
