//go:build porttest && !server

package legacy

/*
#include "GAME1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
*/
import "C"

import (
	"image"
	"slices"
	"unsafe"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestFloorInput struct {
	X, Y        int32
	Tile        int32
	Flags       uint32
	NilViewport bool
}
type PortTestFloorResult struct {
	Values              []int
	Unchanged, Restored bool
}

func PortTestFloorEligibility(inputs []PortTestFloorInput) (out PortTestFloorResult) {
	rows, freeRows := alloc.Make([]uint32{}, 130)
	defer freeRows()
	cells, freeCells := alloc.Make([]uint32{}, 128*11+4)
	defer freeCells()
	oldGrid, oldFlags := C.ptr_5D4594_2650668, noxflags.GetEngine()
	defer func() {
		C.ptr_5D4594_2650668 = oldGrid
		noxflags.ResetEngine()
		noxflags.SetEngine(oldFlags)
		out.Restored = C.ptr_5D4594_2650668 == oldGrid && noxflags.GetEngine() == oldFlags
	}()
	for i := range rows {
		rows[i] = 0x31313131
	}
	for i := 1; i <= 128; i++ {
		rows[i] = uint32(uintptr(unsafe.Pointer(&cells[2])))
	}
	wantRows := append([]uint32(nil), rows...)
	out.Unchanged = true
	for _, in := range inputs {
		for i := range cells {
			cells[i] = 0x35353535
		}
		for y := 0; y < 128; y++ {
			off := 2 + y*11
			cells[off+1], cells[off+6] = uint32(in.Tile), uint32(in.Tile)
			cells[off+5], cells[off+10] = 0, 0
		}
		wantCells := append([]uint32(nil), cells...)
		installed := (**C.obj_5D4594_2650668_t)(unsafe.Pointer(&rows[1]))
		C.ptr_5D4594_2650668 = installed
		noxflags.ResetEngine()
		noxflags.SetEngine(noxflags.EngineFlag(in.Flags))
		// Invalid values in other position fields detect using the wrong viewport corner.
		vp := noxrender.Viewport{Screen: image.Rect(-500, -600, -700, -800), World: image.Rectangle{Min: image.Pt(-1000, -2000), Max: image.Pt(int(in.X), int(in.Y))}, Size: image.Pt(-300, -400), Field10: 123, Field11: 456, Jiggle12: 789}
		before := vp
		arg := &vp
		if in.NilViewport {
			arg = nil
			C.ptr_5D4594_2650668 = nil
			installed = nil
		}
		out.Values = append(out.Values, Nox_xxx_drawAllMB_475810_draw_B(arg))
		out.Unchanged = out.Unchanged && vp == before && noxflags.GetEngine() == noxflags.EngineFlag(in.Flags) && C.ptr_5D4594_2650668 == installed && slices.Equal(rows, wantRows) && slices.Equal(cells, wantCells)
	}
	return out
}
