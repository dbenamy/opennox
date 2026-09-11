//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME1.h"
#include "GAME4_3.h"
*/
import "C"

import (
	"bytes"
	"slices"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
)

const portTestEdgeRows, portTestEdgeRowSize = 64, 60
const portTestEdgeBase = 28644
const portTestEdgeMapBase, portTestEdgeMapWords = 282736, 144

type PortTestEdgeDirectSpec struct {
	Width, Height byte
	Edge          int32
	Index         int32
	Normalize     bool
}
type PortTestEdgeMapSpec struct {
	Width, Height                        byte
	Index, Current, Category, Normalized int32
	Mapping                              uint32
}
type PortTestEdgeResult struct {
	Result                 int32
	LogicIndex, OtherIndex int
}
type PortTestEdgeMapResult struct {
	Return                 int
	Current                uint32
	LogicIndex, OtherIndex int
	RecordGuardsOK         bool
}
type PortTestEdgeSnapshot struct {
	Direct                                            []PortTestEdgeResult
	Mapped                                            []PortTestEdgeMapResult
	LogicBefore, LogicAfter, OtherBefore, OtherAfter  int
	TableUnchanged, MappingUnchanged, GuardsUnchanged bool
	Restored                                          bool
	TableBefore, TableAfterRestore                    []byte
	MappingBefore, MappingAfterRestore                []uint32
}

func portTestEdgeTable() []byte {
	return unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestEdgeBase), portTestEdgeRows*portTestEdgeRowSize)
}
func portTestEdgeMap() []uint32 {
	return unsafe.Slice(memmap.PtrUint32(0x587000, portTestEdgeMapBase), portTestEdgeMapWords)
}

// PortTestEdgeMapping batches direct 543EB0 and dependent 543E60 live C ABI calls.
func PortTestEdgeMapping(seed int, direct []PortTestEdgeDirectSpec, mapped []PortTestEdgeMapSpec) (out PortTestEdgeSnapshot) {
	oldGet := GetServer
	core := new(server.Server)
	core.Rand.Logic, core.Rand.Other = prand.New(seed), prand.New(seed+1)
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	table, mp := portTestEdgeTable(), portTestEdgeMap()
	left := unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestEdgeBase-8), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x85B3FC, uintptr(portTestEdgeBase+len(table))), 8)
	oldTable, oldMap := append([]byte(nil), table...), append([]uint32(nil), mp...)
	out.TableBefore, out.MappingBefore = append([]byte(nil), oldTable...), append([]uint32(nil), oldMap...)
	mapLeft := unsafe.Slice(memmap.PtrUint8(0x587000, portTestEdgeMapBase-8), 8)
	mapRight := unsafe.Slice(memmap.PtrUint8(0x587000, portTestEdgeMapBase+portTestEdgeMapWords*4), 8)
	oldMapLeft, oldMapRight := append([]byte(nil), mapLeft...), append([]byte(nil), mapRight...)
	oldLeft, oldRight := append([]byte(nil), left...), append([]byte(nil), right...)
	defer func() {
		copy(table, oldTable)
		copy(mp, oldMap)
		copy(left, oldLeft)
		copy(right, oldRight)
		copy(mapLeft, oldMapLeft)
		copy(mapRight, oldMapRight)
		GetServer = oldGet
		out.Restored = bytes.Equal(table, oldTable) && slices.Equal(mp, oldMap) && bytes.Equal(left, oldLeft) && bytes.Equal(right, oldRight) && bytes.Equal(mapLeft, oldMapLeft) && bytes.Equal(mapRight, oldMapRight)
		out.TableAfterRestore, out.MappingAfterRestore = append([]byte(nil), table...), append([]uint32(nil), mp...)
	}()
	for i := range table {
		table[i] = 0
	}
	for i := range mp {
		mp[i] = 0xfeedface
	}
	for i := range left {
		left[i] = byte(0x31 + i)
	}
	for i := range right {
		right[i] = byte(0xc1 + i)
	}
	for i := range mapLeft {
		mapLeft[i], mapRight[i] = byte(0x71+i), byte(0xe1+i)
	}
	wantMapLeft, wantMapRight := append([]byte(nil), mapLeft...), append([]byte(nil), mapRight...)
	wantTable, wantMap := append([]byte(nil), table...), append([]uint32(nil), mp...)
	wantLeft, wantRight := append([]byte(nil), left...), append([]byte(nil), right...)
	out.TableUnchanged, out.MappingUnchanged, out.GuardsUnchanged = true, true, true
	check := func() {
		out.TableUnchanged = out.TableUnchanged && bytes.Equal(table, wantTable)
		out.MappingUnchanged = out.MappingUnchanged && slices.Equal(mp, wantMap)
		out.GuardsUnchanged = out.GuardsUnchanged && bytes.Equal(left, wantLeft) && bytes.Equal(right, wantRight) && bytes.Equal(mapLeft, wantMapLeft) && bytes.Equal(mapRight, wantMapRight)
	}
	out.LogicBefore, out.OtherBefore = core.Rand.Logic.Index(), core.Rand.Other.Index()
	for _, s := range direct {
		if s.Index < 0 || s.Index >= portTestEdgeRows {
			panic("invalid edge row")
		}
		off := int(s.Index) * portTestEdgeRowSize
		table[off+52], table[off+53] = s.Width, s.Height
		wantTable[off+52], wantTable[off+53] = s.Width, s.Height
		var result C.int
		if s.Normalize {
			result = C.sub_411490(C.int(s.Index), C.int(s.Edge))
		} else {
			result = C.nox_xxx_mapGenEdge_543EB0(C.int(s.Index), C.int(s.Edge))
		}
		out.Direct = append(out.Direct, PortTestEdgeResult{Result: int32(result), LogicIndex: core.Rand.Logic.Index(), OtherIndex: core.Rand.Other.Index()})
		check()
	}
	for _, s := range mapped {
		if s.Index < 0 || s.Index >= portTestEdgeRows {
			panic("invalid edge row")
		}
		off := int(s.Index) * portTestEdgeRowSize
		table[off+52], table[off+53] = s.Width, s.Height
		wantTable[off+52], wantTable[off+53] = s.Width, s.Height
		class := int(s.Normalized)
		slot := int(s.Category) + 12*class
		if slot < 0 || slot >= portTestEdgeMapWords {
			panic("unnormalized edge mapping slot")
		}
		for i := range mp {
			mp[i], wantMap[i] = 255, 255
		}
		mp[slot], wantMap[slot] = s.Mapping, s.Mapping
		mem := C.calloc(6, C.size_t(unsafe.Sizeof(C.uint32_t(0))))
		if mem == nil {
			panic("calloc edge record")
		}
		words := unsafe.Slice((*C.uint32_t)(mem), 6)
		words[0], words[1], words[2], words[3], words[4], words[5] = 0xa0a0a0a0, 0x11111111, 0x22222222, C.uint32_t(s.Index), C.uint32_t(s.Current), 0xb0b0b0b0
		ret := C.sub_543E60(C.int(uintptr(unsafe.Pointer(&words[1]))), C.int(s.Category))
		out.Mapped = append(out.Mapped, PortTestEdgeMapResult{Return: int(ret), Current: uint32(words[4]), LogicIndex: core.Rand.Logic.Index(), OtherIndex: core.Rand.Other.Index(), RecordGuardsOK: words[0] == 0xa0a0a0a0 && words[1] == 0x11111111 && words[2] == 0x22222222 && words[3] == C.uint32_t(s.Index) && words[5] == 0xb0b0b0b0})
		C.free(mem)
		check()
	}
	out.LogicAfter, out.OtherAfter = core.Rand.Logic.Index(), core.Rand.Other.Index()
	return out
}
