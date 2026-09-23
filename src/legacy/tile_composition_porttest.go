//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1.h"
#include "GAME2_2.h"
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
	old := worldTileGrid
	counter := nox_xxx_waypointCounterMB_587000_154948
	if worldTileGridCapacity != 128 {
		panic("unexpected tile grid capacity")
	}
	if worldGridAllocate() == 0 {
		panic("tile grid allocation")
	}
	rows := (*[128]*[128][11]uint32)(unsafe.Pointer(worldTileGrid))
	return rows, (*uint32)(unsafe.Pointer(&nox_xxx_waypointCounterMB_587000_154948)), func() {
		worldGridFreeRows()
		C.free(unsafe.Pointer(worldTileGrid))
		worldTileGrid = old
		nox_xxx_waypointCounterMB_587000_154948 = counter
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
