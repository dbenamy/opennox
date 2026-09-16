//go:build porttest

package legacy

/*
#include <stdlib.h>
#include <stdint.h>
extern uint32_t dword_5d4594_251572;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Arrays allocated by the original readers use raw calloc/free. Fixture cleanup
// must not pass them to the separately tracked alloc.FreePtr owner.
func PortTestFloorAssetsRelease(defs []server.TileDef, edges []byte) {
	for i := range defs {
		if p := defs[i].Data32; p != nil {
			C.free(p)
			defs[i].Data32 = nil
		}
	}
	for i := 0; i < 64; i++ {
		p := (*unsafe.Pointer)(unsafe.Pointer(&edges[60*i+32]))
		if *p != nil {
			C.free(*p)
			*p = nil
		}
	}
}
func PortTestFloorAssetsOwner() ([]server.TileDef, []byte, *uint32, *uint32, func()) {
	defs, count, _, _, _ := portTestTilePtrs()
	edges := portTestEdgeTable()
	edgeCount := (*uint32)(unsafe.Pointer(&C.dword_5d4594_251572))
	oldDefs, oldEdges := append([]byte(nil), tileBytes(defs)...), append([]byte(nil), edges...)
	dc, ec := *count, *edgeCount
	clear(defs)
	clear(edges)
	*count = 0
	*edgeCount = 0
	return defs, edges, count, edgeCount, func() {
		PortTestFloorAssetsRelease(defs, edges)
		copy(tileBytes(defs), oldDefs)
		copy(edges, oldEdges)
		*count = dc
		*edgeCount = ec
	}
}
func PortTestFloorAssetsAlloc(n int) unsafe.Pointer { return C.calloc(C.size_t(n), 1) }
func PortTestFloorAssetsFacade(index int) int {
	if floorAssetFacade(&tileDefinitionsAll()[index]) {
		return 1
	}
	return 0
}
