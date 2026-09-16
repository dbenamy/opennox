//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1.h"
#include "GAME2_2.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern const int ptr_5D4594_2650668_cap;
extern uint32_t nox_xxx_waypointCounterMB_587000_154948;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// Use the production grid allocator; its row-free routine leaves the pointer
// table to its caller. Restore the previous owner after freeing this test grid.
func PortTestTileCompositionGrid() (*[128]*[128][11]uint32, *uint32, func()) {
	old := C.ptr_5D4594_2650668
	counter := C.nox_xxx_waypointCounterMB_587000_154948
	if C.ptr_5D4594_2650668_cap != 128 {
		panic("unexpected tile grid capacity")
	}
	if C.nox_xxx_tileAlloc_410F60_init() == 0 {
		panic("tile grid allocation")
	}
	rows := (*[128]*[128][11]uint32)(unsafe.Pointer(C.ptr_5D4594_2650668))
	return rows, (*uint32)(unsafe.Pointer(&C.nox_xxx_waypointCounterMB_587000_154948)), func() {
		C.nox_xxx_tileFree_410FC0_free()
		C.free(unsafe.Pointer(C.ptr_5D4594_2650668))
		C.ptr_5D4594_2650668 = old
		C.nox_xxx_waypointCounterMB_587000_154948 = counter
	}
}
func PortTestTileCompositionEdges(pos image.Point, edge unsafe.Pointer) {
	tileCompositionEdges(pos, (*[5]uint32)(edge))
}

func PortTestTileCompositionFull(vp *noxrender.Viewport) {
	tileCompositionFull(vp)
}
func PortTestTileCompositionHorizontal(vp *noxrender.Viewport, x int) {
	tileCompositionHorizontal(vp, x)
}
func PortTestTileCompositionVertical(vp *noxrender.Viewport, y int) {
	tileCompositionVertical(vp, y)
}
func PortTestTileCompositionRedraw(vp *noxrender.Viewport) int {
	return tileCompositionRedraw(vp)
}
