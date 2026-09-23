//go:build porttest

package legacy

/*
#include "GAME2_2.h"


*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/server"
	"image"
	"unsafe"
)

func PortTestTileRasterOwner() (map[string]*uint32, []server.TileDef, func()) {
	words := map[string]*uint32{
		"width":    (*uint32)(unsafe.Pointer(&dword_5d4594_3798800)),
		"stride":   (*uint32)(unsafe.Pointer(&dword_5d4594_3798804)),
		"height":   (*uint32)(unsafe.Pointer(&dword_5d4594_3798808)),
		"columns":  (*uint32)(unsafe.Pointer(&dword_5d4594_3798812)),
		"rows":     (*uint32)(unsafe.Pointer(&dword_5d4594_3798816)),
		"originX":  (*uint32)(unsafe.Pointer(&dword_5d4594_3798820)),
		"originY":  (*uint32)(unsafe.Pointer(&dword_5d4594_3798824)),
		"tileX":    (*uint32)(unsafe.Pointer(&dword_5d4594_3798828)),
		"tileY":    (*uint32)(unsafe.Pointer(&dword_5d4594_3798832)),
		"scrollX":  (*uint32)(unsafe.Pointer(&dword_5d4594_3798836)),
		"scrollY":  (*uint32)(unsafe.Pointer(&dword_5d4594_3798840)),
		"flatFlag": (*uint32)(unsafe.Pointer(&dword_5d4594_1193156)),
		"dirty":    (*uint32)(unsafe.Pointer(&dword_5d4594_1193188)),
	}
	old := map[string]uint32{}
	for n, p := range words {
		old[n] = *p
	}
	defs, _, _, _, _ := portTestTilePtrs()
	savedDefs := append([]byte(nil), tileBytes(defs)...)
	begin, end := legacyGlobals.nox_video_tileBuf_ptr_3798796, legacyGlobals.nox_video_tileBuf_end_3798844
	draw, edges := tileDrawCallback, tileEdgeCallback
	return words, defs, func() {
		for n, p := range words {
			*p = old[n]
		}
		copy(tileBytes(defs), savedDefs)
		legacyGlobals.nox_video_tileBuf_ptr_3798796, legacyGlobals.nox_video_tileBuf_end_3798844 = begin, end
		tileDrawCallback, tileEdgeCallback = draw, edges
	}
}
func PortTestTileRasterPrimitive(fill bool, dst, source unsafe.Pointer, color uint32) {
	if fill {
		tileRasterFillPacked(dst, color)
	} else {
		tileRasterCopyPacked(dst, source)
	}
}

func PortTestTileRasterDispatch(pos image.Point, img noxrender.ImageHandle, tile uint32) {
	tileDrawCallback(pos, img, uint16(tile))
}
