//go:build porttest

package legacy

/*
#include "GAME1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
static uint32_t porttest_grid_bench(float2* point, unsigned n) {
 uint32_t sum=0;
 for (unsigned i=0;i<n;i++) sum+=nox_xxx_tileNFromPoint_411160(point);
 return sum;
}
*/
import "C"

import (
	"bytes"
	"math"
	"slices"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestGridLookupResult struct {
	Results                                                                                  []int32
	Checksum                                                                                 uint32
	GridUnchanged, TableUnchanged, NodesUnchanged, InputGuardsOK, PointerUnchanged, Restored bool
}

// PortTestGridLookup supplies guarded C-owned rows, cells, nodes and point.
// The grid is immutable throughout each batch; compare every byte afterward.
func PortTestGridLookup(inputs [][2]uint32, listCount int, goWrapper bool) PortTestGridLookupResult {
	return portTestGridLookup(inputs, listCount, goWrapper, 0)
}
func PortTestGridLookupBenchmark(n int, goWrapper bool) PortTestGridLookupResult {
	return portTestGridLookup([][2]uint32{{math.Float32bits(128.5), math.Float32bits(193.25)}}, 0, goWrapper, n)
}
func portTestGridLookup(inputs [][2]uint32, listCount int, goWrapper bool, benchN int) (out PortTestGridLookupResult) {
	if listCount < 0 || listCount > 12 {
		panic("invalid grid list count")
	}
	const rowStride = 128*11 + 4
	rows, freeRows := alloc.Make([]uint32{}, 130)
	defer freeRows()
	cells, freeCells := alloc.Make([]uint32{}, 128*rowStride)
	defer freeCells()
	nodes, freeNodes := alloc.Make([]uint32{}, 2*(listCount+2)*5)
	defer freeNodes()
	point, freePoint := alloc.Make([]uint32{}, 4)
	defer freePoint()
	for i := range rows {
		rows[i] = 0x31313131
	}
	for i := range cells {
		cells[i] = 0x35350000 + uint32(i)
	}
	for i := range nodes {
		nodes[i] = 0x5a5a5a5a
	}
	for half := 0; half < 2; half++ {
		base := half * (listCount + 2) * 5
		for n := 0; n < listCount; n++ {
			off := base + (n+1)*5
			prefix := uint32(0x2a000000)
			if half == 1 {
				prefix = 0xca000000
			}
			nodes[off], nodes[off+1], nodes[off+2], nodes[off+3], nodes[off+4] = prefix+uint32(n), 0x17170000+uint32(n), 0, uint32(n), 0
			if n+1 < listCount {
				nodes[off+4] = uint32(uintptr(unsafe.Pointer(&nodes[off+5])))
			}
		}
	}
	for x := 0; x < 128; x++ {
		rows[x+1] = uint32(uintptr(unsafe.Pointer(&cells[x*rowStride+2])))
		for y := 0; y < 128; y++ {
			off := x*rowStride + 2 + y*11
			cells[off+1] = 0x1e000000 | uint32(x<<8|y)
			cells[off+6] = 0xbe000000 | uint32(x<<8|y)
			cells[off+5], cells[off+10] = 0, 0
			if listCount > 0 {
				cells[off+5] = uint32(uintptr(unsafe.Pointer(&nodes[5])))
				cells[off+10] = uint32(uintptr(unsafe.Pointer(&nodes[(listCount+2)*5+5])))
			}
		}
	}
	wantRows, wantCells, wantNodes := append([]uint32(nil), rows...), append([]uint32(nil), cells...), append([]uint32(nil), nodes...)
	table := portTestEdgeTable()
	oldTable := append([]byte(nil), table...)
	left := unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestEdgeBase-8), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x85B3FC, uintptr(portTestEdgeBase+len(table))), 8)
	oldLeft, oldRight := append([]byte(nil), left...), append([]byte(nil), right...)
	oldGrid := C.ptr_5D4594_2650668
	defer func() {
		C.ptr_5D4594_2650668 = oldGrid
		copy(table, oldTable)
		copy(left, oldLeft)
		copy(right, oldRight)
		out.Restored = C.ptr_5D4594_2650668 == oldGrid && bytes.Equal(table, oldTable) && bytes.Equal(left, oldLeft) && bytes.Equal(right, oldRight)
	}()
	table[52], table[53] = 3, 3
	for i := range left {
		left[i], right[i] = byte(0x51+i), byte(0xd1+i)
	}
	wantLeft, wantRight := append([]byte(nil), left...), append([]byte(nil), right...)
	wantTable := append([]byte(nil), table...)
	installed := (**C.obj_5D4594_2650668_t)(unsafe.Pointer(&rows[1]))
	C.ptr_5D4594_2650668 = installed
	out.Results = make([]int32, len(inputs))
	out.InputGuardsOK, out.PointerUnchanged = true, true
	for i, b := range inputs {
		point[0], point[1], point[2], point[3] = 0x12345678, b[0], b[1], 0x76543210
		if benchN > 0 {
			if goWrapper {
				p := types.Pointf{X: math.Float32frombits(b[0]), Y: math.Float32frombits(b[1])}
				for n := 0; n < benchN; n++ {
					out.Checksum += uint32(Nox_xxx_tileNFromPoint_411160(p))
				}
			} else {
				out.Checksum = uint32(C.porttest_grid_bench((*C.float2)(unsafe.Pointer(&point[1])), C.uint(benchN)))
			}
		} else if goWrapper {
			out.Results[i] = int32(Nox_xxx_tileNFromPoint_411160(types.Pointf{X: math.Float32frombits(b[0]), Y: math.Float32frombits(b[1])}))
		} else {
			out.Results[i] = int32(C.nox_xxx_tileNFromPoint_411160((*C.float2)(unsafe.Pointer(&point[1]))))
		}
		out.InputGuardsOK = out.InputGuardsOK && point[0] == 0x12345678 && point[1] == b[0] && point[2] == b[1] && point[3] == 0x76543210
		out.PointerUnchanged = out.PointerUnchanged && C.ptr_5D4594_2650668 == installed
	}
	out.GridUnchanged = slices.Equal(rows, wantRows) && slices.Equal(cells, wantCells)
	out.NodesUnchanged = slices.Equal(nodes, wantNodes)
	out.TableUnchanged = bytes.Equal(table, wantTable) && bytes.Equal(left, wantLeft) && bytes.Equal(right, wantRight)
	return out
}
