//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

// Reuse real painting owners and guarded cache records. Only dispatch is new;
// floor/wall services, allocation tracking and state capture remain unchanged.
func PortTestPrefabPainting(cases []PortTestPaintSpec, owner func(*server.Server) (Server, func())) []PortTestPaintResult {
	words, restore := PortTestPrefabRuntimeGlobals()
	defer restore()
	ext := &paintTestExtension{globals: words}
	ext.globals["transparentFloor"] = memmap.PtrUint32(0x587000, 229704)
	for off := uintptr(1599484); off < 1599532; off += 4 {
		ext.globals[fmt.Sprintf("cacheBlob%d", off)] = memmap.PtrUint32(0x5D4594, off)
	}
	for _, off := range []uintptr{35920, 35924, 35928, 35932, 35936, 35940, 35944, 35960, 35964, 35968, 35980} {
		ext.globals[fmt.Sprintf("selection%d", off)] = memmap.PtrUint32(0x973F18, off)
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 { return uint32(PortTestPrefabCall(op, args)) }
	return portTestMapPainting(cases, owner, ext)
}
