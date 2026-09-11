//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME1.h"
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

const portTestSubtileRows, portTestSubtileRowSize = 64, 60
const portTestSubtileBase = 28644

type PortTestSubtileRow struct {
	Index         int
	Width, Height byte
}
type PortTestSubtileNode struct{ Value, Row, Edge int32 }
type PortTestSubtileLookupSpec struct {
	// Mode 0 calls 4113A0. Mode 1 calls 411350 over Nodes.
	Mode                     byte
	X, Y, Category, Fallback int32
	Rows                     []PortTestSubtileRow
	Nodes                    []PortTestSubtileNode
	NilList, NilPoint        bool // NilPoint is defined only with NilList.
}
type PortTestSubtileLookupResult struct {
	Return                                                          int32
	TableUnchanged, PointUnchanged, NodesUnchanged, GuardsUnchanged bool
}
type PortTestSubtileLookupSnapshot struct {
	Results                        []PortTestSubtileLookupResult
	TableBefore, TableAfterRestore []byte
	GuardsRestored                 bool
}

func portTestSubtileTable() []byte {
	return unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestSubtileBase), portTestSubtileRows*portTestSubtileRowSize)
}

// PortTestSubtileLookup exercises original C 4113A0 and 411350. Nodes and the
// two-word point are C allocations, so C link pointers never point into Go.
func PortTestSubtileLookup(specs []PortTestSubtileLookupSpec) (out PortTestSubtileLookupSnapshot) {
	table := portTestSubtileTable()
	left := unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestSubtileBase-8), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x85B3FC, uintptr(portTestSubtileBase+len(table))), 8)
	oldTable, oldLeft, oldRight := append([]byte(nil), table...), append([]byte(nil), left...), append([]byte(nil), right...)
	out.TableBefore = append([]byte(nil), oldTable...)
	defer func() {
		copy(table, oldTable)
		copy(left, oldLeft)
		copy(right, oldRight)
		out.TableAfterRestore = append([]byte(nil), table...)
		out.GuardsRestored = bytes.Equal(left, oldLeft) && bytes.Equal(right, oldRight)
	}()
	for i := range table {
		table[i] = 0
	}
	for i := range left {
		left[i] = byte(0x31 + i)
	}
	for i := range right {
		right[i] = byte(0xc1 + i)
	}
	wantLeft, wantRight := append([]byte(nil), left...), append([]byte(nil), right...)
	for _, s := range specs {
		for _, row := range s.Rows {
			if row.Index < 0 || row.Index >= portTestSubtileRows {
				panic("invalid subtile row")
			}
			off := row.Index * portTestSubtileRowSize
			table[off+52], table[off+53] = row.Width, row.Height
		}
		beforeTable := append([]byte(nil), table...)
		point := C.calloc(4, C.size_t(unsafe.Sizeof(C.int(0))))
		if point == nil {
			panic("calloc point")
		}
		pw := unsafe.Slice((*C.int)(point), 4)
		pw[0], pw[1], pw[2], pw[3] = 0x12345678, C.int(s.X), C.int(s.Y), 0x76543210
		beforePoint := append([]C.int(nil), pw...)
		var nodes unsafe.Pointer
		var nw []C.int
		if !s.NilList {
			nodes = C.calloc(C.size_t(len(s.Nodes)+2), 5*C.size_t(unsafe.Sizeof(C.int(0))))
			if nodes == nil {
				panic("calloc nodes")
			}
			nw = unsafe.Slice((*C.int)(nodes), (len(s.Nodes)+2)*5)
			for i := range nw {
				nw[i] = C.int(0x5a5a5a5a)
			}
			for i, n := range s.Nodes {
				if n.Row < 0 || n.Row >= 64 {
					panic("invalid node row")
				}
				off := (i + 1) * 5
				nw[off], nw[off+1], nw[off+2], nw[off+3] = C.int(n.Value), C.int(0x11110000+i), C.int(n.Row), C.int(n.Edge)
				nw[off+4] = 0
				if i+1 < len(s.Nodes) {
					nw[off+4] = C.int(uintptr(unsafe.Pointer(&nw[off+5])))
				}
			}
		}
		beforeNodes := append([]C.int(nil), nw...)
		var ret C.int
		if s.Mode == 0 {
			ret = C.sub_4113A0(&pw[1], C.int(s.Category))
		} else if s.Mode == 1 {
			if s.NilPoint && !s.NilList {
				panic("nonnull C list with nil point faults")
			}
			var head, p *C.int
			if !s.NilList && len(s.Nodes) > 0 {
				head = &nw[5]
			}
			if !s.NilPoint {
				p = &pw[1]
			}
			ret = C.sub_411350(head, p, C.int(s.Fallback))
		} else {
			panic("invalid subtile mode")
		}
		out.Results = append(out.Results, PortTestSubtileLookupResult{Return: int32(ret), TableUnchanged: bytes.Equal(table, beforeTable), PointUnchanged: bytes.Equal(unsafe.Slice((*byte)(point), 4*int(unsafe.Sizeof(C.int(0)))), unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(beforePoint))), len(beforePoint)*int(unsafe.Sizeof(C.int(0))))), NodesUnchanged: bytes.Equal(unsafe.Slice((*byte)(nodes), len(nw)*int(unsafe.Sizeof(C.int(0)))), unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(beforeNodes))), len(beforeNodes)*int(unsafe.Sizeof(C.int(0))))), GuardsUnchanged: bytes.Equal(left, wantLeft) && bytes.Equal(right, wantRight)})
		C.free(point)
		if nodes != nil {
			C.free(nodes)
		}
	}
	return out
}
