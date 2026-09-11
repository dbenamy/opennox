//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME1.h"
#include "GAME4_1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern uint32_t nox_tile_def_cnt;
extern nox_tileDef_t nox_tile_defs_arr[176];
*/
import "C"

import (
	"bytes"
	"fmt"
	"slices"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

const (
	portTestGeneratorWaterNoTeleport = 1 + iota
	portTestGeneratorWaterDeepNoTeleport
	portTestGeneratorWaterShallowNoTeleport
	portTestGeneratorWaterSwampDeepNoTeleport
	portTestGeneratorWaterSwampShallowNoTeleport
)

// portTestGeneratorTileEnvironment supplies the real grid ABI consumed by
// nox_xxx_mapTileAllowTeleport_411A90. The owner fixture initializes the map
// separately with PortTestAIEmptyMap; this helper owns only legacy grid/tile
// state. configure installs one tile ID in both grid result fields. cold
// forces the C five-name lazy scan; hot preloads its result cache.
func portTestGeneratorTileEnvironment() (configure func(tileID int, cold bool), unchanged func() bool, restore func()) {
	const (
		rowStride  = 128*11 + 4
		cellWords  = 128 * rowStride
		guardWords = 2
	)
	rows, freeRows := alloc.Make([]uint32{}, 130) // rows[0] and rows[129] guard the dispatch range [1,128].
	cellMem, freeCells := alloc.Make([]uint32{}, cellWords+2*guardWords)
	cells := cellMem[guardWords : guardWords+cellWords]
	cellPrefix := cellMem[:guardWords]
	cellSuffix := cellMem[guardWords+cellWords:]

	tiles := unsafe.Slice((*server.TileDef)(unsafe.Pointer(&C.nox_tile_defs_arr[0])), 176)
	tileBytes := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(tiles))), len(tiles)*int(unsafe.Sizeof(server.TileDef{})))
	oldTiles := bytes.Clone(tileBytes)
	oldCount := C.nox_tile_def_cnt
	oldGrid := C.ptr_5D4594_2650668
	cacheOff := [...]uintptr{26516, 26520, 26524, 26528, 26532}
	oldCache := make([]uint32, len(cacheOff))
	for i, off := range cacheOff {
		oldCache[i] = *memmap.PtrUint32(0x587000, off)
	}

	edges := portTestEdgeTable()
	oldEdges := bytes.Clone(edges)
	left := unsafe.Slice(memmap.PtrUint8(0x85B3FC, portTestEdgeBase-8), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x85B3FC, uintptr(portTestEdgeBase+len(edges))), 8)
	oldLeft, oldRight := bytes.Clone(left), bytes.Clone(right)

	for i := range rows {
		rows[i] = 0x51510000 + uint32(i)
	}
	for i := range cellMem {
		cellMem[i] = 0x61610000 + uint32(i)
	}
	for x := 0; x < 128; x++ {
		rows[x+1] = uint32(uintptr(unsafe.Pointer(&cells[x*rowStride+2])))
	}
	baseRows := slices.Clone(rows)
	baseCells := slices.Clone(cells)
	for i := range left {
		left[i], right[i] = byte(0x31+i), byte(0xd1+i)
	}
	basePrefix, baseSuffix := slices.Clone(cellPrefix), slices.Clone(cellSuffix)
	baseLeft, baseRight := bytes.Clone(left), bytes.Clone(right)

	var wantRows, wantCells []uint32
	var wantTiles, wantEdges, wantPrefix, wantSuffix []byte
	var wantCache, initialCache [len(cacheOff)]uint32
	var wantGrid **C.obj_5D4594_2650668_t

	configure = func(tileID int, cold bool) {
		if tileID < 0 || tileID >= len(tiles) {
			panic(fmt.Sprintf("invalid generator tile ID %d", tileID))
		}
		copy(rows, baseRows)
		copy(cells, baseCells)
		copy(cellPrefix, basePrefix)
		copy(cellSuffix, baseSuffix)
		copy(left, baseLeft)
		copy(right, baseRight)
		copy(edges, oldEdges)

		for i := range tiles {
			tiles[i] = server.TileDef{}
		}
		for i, name := range []string{
			"WaterNoTeleport", "WaterDeepNoTeleport", "WaterShallowNoTeleport", "WaterSwampDeepNoTeleport", "WaterSwampShallowNoTeleport",
		} {
			copy(tiles[i+1].NameBuf[:], name)
		}
		C.nox_tile_def_cnt = 176
		for i, off := range cacheOff {
			if cold {
				*memmap.PtrUint32(0x587000, off) = 0xffffffff
			} else {
				*memmap.PtrUint32(0x587000, off) = uint32(i + 1)
			}
			initialCache[i] = *memmap.PtrUint32(0x587000, off)
			wantCache[i] = uint32(i + 1)
		}
		edges[52], edges[53] = 3, 3
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				off := x*rowStride + 2 + y*11
				// 411160 returns Field1 or Field6 according to the sub-tile.
				cells[off+1], cells[off+6] = uint32(tileID), uint32(tileID)
				cells[off+5], cells[off+10] = 0, 0
			}
		}
		wantGrid = (**C.obj_5D4594_2650668_t)(unsafe.Pointer(&rows[1]))
		C.ptr_5D4594_2650668 = wantGrid

		wantRows = slices.Clone(rows)
		wantCells = slices.Clone(cells)
		wantTiles = bytes.Clone(tileBytes)
		wantEdges = bytes.Clone(edges)
		wantPrefix = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(cellPrefix))), len(cellPrefix)*4))
		wantSuffix = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(cellSuffix))), len(cellSuffix)*4))
	}
	unchanged = func() bool {
		if wantRows == nil {
			return false
		}
		if C.ptr_5D4594_2650668 != wantGrid || C.nox_tile_def_cnt != 176 ||
			!slices.Equal(rows, wantRows) || !slices.Equal(cells, wantCells) ||
			!bytes.Equal(tileBytes, wantTiles) || !bytes.Equal(edges, wantEdges) ||
			!bytes.Equal(left, baseLeft) || !bytes.Equal(right, baseRight) {
			return false
		}
		prefix := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(cellPrefix))), len(cellPrefix)*4)
		suffix := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(cellSuffix))), len(cellSuffix)*4)
		if !bytes.Equal(prefix, wantPrefix) || !bytes.Equal(suffix, wantSuffix) {
			return false
		}
		warm, original := true, true
		for i, off := range cacheOff {
			v := *memmap.PtrUint32(0x587000, off)
			warm = warm && v == wantCache[i]
			original = original && v == initialCache[i]
		}
		return warm || original
	}
	restore = func() {
		C.ptr_5D4594_2650668 = oldGrid
		copy(tileBytes, oldTiles)
		C.nox_tile_def_cnt = oldCount
		for i, off := range cacheOff {
			*memmap.PtrUint32(0x587000, off) = oldCache[i]
		}
		copy(edges, oldEdges)
		copy(left, oldLeft)
		copy(right, oldRight)
		freeCells()
		freeRows()
	}
	return configure, unchanged, restore
}

// PortTestGeneratorTileContracts checks the retained teleport-tile service
// against known tile IDs before any generator algorithm is converted.
func PortTestGeneratorTileContracts() (cases int, intact bool) {
	configure, unchanged, restore := portTestGeneratorTileEnvironment()
	defer restore()
	intact = true
	for _, cold := range []bool{false, true} {
		for _, tile := range []int{0, 1, 2, 3, 4, 5, 6, 175} {
			configure(tile, cold)
			for _, xy := range [][2]float32{{100, 100}, {55, 55}, {145, 100}, {100, 145}, {128.5, 193.25}, {0, 0}, {500, 700}, {2000, 2000}} {
				p, free := alloc.New([2]float32{})
				*p = xy
				got := C.nox_xxx_mapTileAllowTeleport_411A90((*C.float2)(unsafe.Pointer(p))) != 0
				free()
				want := tile >= 1 && tile <= 5 && xy != [2]float32{55, 55} && xy != [2]float32{0, 0}
				if got != want {
					panic(fmt.Sprintf("tile %d cold %v point %v result %v want %v", tile, cold, xy, got, want))
				}
				cases++
			}
			intact = intact && unchanged()
		}
	}
	return
}
