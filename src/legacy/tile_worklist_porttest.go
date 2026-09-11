//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include "GAME4_2.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
extern uint32_t dword_5d4594_2487248;

static obj_5D4594_2650668_t** portTestWorklistGridNew(void) {
	obj_5D4594_2650668_t** p = calloc(128, sizeof(*p));
	if (!p) return NULL;
	for (int i = 0; i < 128; i++) {
		p[i] = calloc(128, sizeof(*p[i]));
		if (!p[i]) { while (i) free(p[--i]); free(p); return NULL; }
	}
	return p;
}
static void portTestWorklistGridFree(obj_5D4594_2650668_t** p) {
	if (!p) return;
	for (int i = 0; i < 128; i++) free(p[i]);
	free(p);
}
static obj_5D4594_2650668_t** portTestWorklistGridGet(void) { return ptr_5D4594_2650668; }
static void portTestWorklistGridSet(obj_5D4594_2650668_t** p) { ptr_5D4594_2650668 = p; }
static void portTestWorklistCellSet(obj_5D4594_2650668_t** p, int x, int y, uint32_t one, uint32_t two) {
	p[x][y].field_1 = one; p[x][y].field_6 = two;
}
static int portTestWorklistGridEqual(obj_5D4594_2650668_t** a, obj_5D4594_2650668_t** b) {
	for (int i = 0; i < 128; i++) if (memcmp(a[i], b[i], 128*sizeof(*a[i]))) return 0;
	return 1;
}
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

const portTestWorklistWords = 1500 // 500 records × (x, y, flags)
const portTestWorklistNilRef = ^uint16(0)

type PortTestTileWorklistSpec struct {
	// Op 0 enqueues X/Y/Flags/Key after installing Field1/Field2 at X/Y.
	// Op 1 pops into output references.
	Op                  byte
	X, Y, Flags         int32
	Key, Field1, Field2 uint32
	// NilGrid is valid only on an enqueue path which C proves does not read the
	// grid: out-of-range coordinates or Flags&3 == 0.
	NilGrid bool
	// Pop references: 0..2 are separate guarded C words; 3=count; 4=overflow;
	// 5+n=queue word n (0..1499); portTestWorklistNilRef is a NULL output and
	// is valid only for a signed-nonpositive count.
	PopX, PopY, PopZ uint16
}

type PortTestTileWorklistResult struct {
	Return                               int
	Count, Overflow                      uint32
	Outputs                              [3]uint32
	OutputReadable                       [3]bool
	Queue                                []uint32
	GridUnchanged, GridPointerUnchanged  bool
	QueueGuardsUnchanged, OutputGuardsOK bool
}

type PortTestTileWorklistState struct {
	Count, Overflow uint32
	Queue           []uint32
	Left, Right     []byte
}

type PortTestTileWorklistSnapshot struct {
	Results              []PortTestTileWorklistResult
	Before, AfterRestore PortTestTileWorklistState
	GridPointerRestored  bool
}

func portTestWorklistQueue() []uint32 {
	return unsafe.Slice(memmap.PtrUint32(0x973F18, 16200), portTestWorklistWords)
}
func portTestWorklistState(count, overflow *uint32, queue []uint32, left, right []byte) PortTestTileWorklistState {
	return PortTestTileWorklistState{Count: *count, Overflow: *overflow, Queue: append([]uint32(nil), queue...), Left: append([]byte(nil), left...), Right: append([]byte(nil), right...)}
}
func portTestWorklistOut(ref uint16, words []C.uint32_t, count, overflow *uint32, queue []uint32) *C.uint32_t {
	switch {
	case ref < 3:
		return &words[1+ref] // words[0]/[4] are guards
	case ref == 3:
		return (*C.uint32_t)(unsafe.Pointer(count))
	case ref == 4:
		return (*C.uint32_t)(unsafe.Pointer(overflow))
	case ref == portTestWorklistNilRef:
		return nil
	case int(ref)-5 < len(queue):
		return (*C.uint32_t)(unsafe.Pointer(&queue[int(ref)-5]))
	default:
		panic("invalid worklist output reference")
	}
}
func portTestWorklistOutValue(ref uint16, words []C.uint32_t, count, overflow *uint32, queue []uint32) (uint32, bool) {
	if ref == portTestWorklistNilRef {
		return 0, false
	}
	return uint32(*portTestWorklistOut(ref, words, count, overflow, queue)), true
}

// PortTestTileWorklist invokes the live C ABI producer/consumer against an
// isolated C-owned 128x128 grid. It retains every physical queue word after
// each call, including inactive records, so a port cannot merely model a slice.
func PortTestTileWorklist(initialCount, initialOverflow uint32, initialQueue []uint32, specs []PortTestTileWorklistSpec) (snap PortTestTileWorklistSnapshot) {
	if len(initialQueue) != portTestWorklistWords {
		panic("worklist initial queue must contain 1500 words")
	}
	count := (*uint32)(unsafe.Pointer(&C.dword_5d4594_2487248))
	overflow := memmap.PtrUint32(0x973F18, 22200)
	queue := portTestWorklistQueue()
	left := unsafe.Slice(memmap.PtrUint8(0x973F18, 16192), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x973F18, 22204), 8)
	oldGrid := C.portTestWorklistGridGet()
	snap.Before = portTestWorklistState(count, overflow, queue, left, right)
	grid, expected := C.portTestWorklistGridNew(), C.portTestWorklistGridNew()
	if grid == nil || expected == nil {
		C.portTestWorklistGridFree(grid)
		C.portTestWorklistGridFree(expected)
		panic("calloc grid failed")
	}
	defer func() {
		copy(queue, snap.Before.Queue)
		*count, *overflow = snap.Before.Count, snap.Before.Overflow
		copy(left, snap.Before.Left)
		copy(right, snap.Before.Right)
		C.portTestWorklistGridSet(oldGrid)
		C.portTestWorklistGridFree(grid)
		C.portTestWorklistGridFree(expected)
		snap.AfterRestore = portTestWorklistState(count, overflow, queue, left, right)
		snap.GridPointerRestored = C.portTestWorklistGridGet() == oldGrid
	}()
	C.portTestWorklistGridSet(grid)
	copy(queue, initialQueue)
	*count, *overflow = initialCount, initialOverflow
	for i := range left {
		left[i] = byte(0x31 + i)
	}
	for i := range right {
		right[i] = byte(0xc1 + i)
	}
	wantLeft, wantRight := append([]byte(nil), left...), append([]byte(nil), right...)

	for _, s := range specs {
		mem := C.calloc(5, C.size_t(unsafe.Sizeof(C.uint32_t(0))))
		if mem == nil {
			panic("calloc outputs failed")
		}
		words := unsafe.Slice((*C.uint32_t)(mem), 5)
		words[0], words[1], words[2], words[3], words[4] = 0xa0a0a0a0, 0x01020304, 0x11223344, 0x55667788, 0xb0b0b0b0
		var ret C.int
		if s.Op == 0 {
			if s.X > 0 && s.X < 127 && s.Y > 0 && s.Y < 127 {
				C.portTestWorklistCellSet(grid, C.int(s.X), C.int(s.Y), C.uint32_t(s.Field1), C.uint32_t(s.Field2))
				C.portTestWorklistCellSet(expected, C.int(s.X), C.int(s.Y), C.uint32_t(s.Field1), C.uint32_t(s.Field2))
			}
			wantGrid := grid
			if s.NilGrid {
				wantGrid = nil
				C.portTestWorklistGridSet(nil)
			}
			C.sub_51DD50(C.int(s.X), C.int(s.Y), C.int(s.Flags), C.int(s.Key))
			pointerOK := C.portTestWorklistGridGet() == wantGrid
			if s.NilGrid {
				C.portTestWorklistGridSet(grid)
			}
			v0, ok0 := portTestWorklistOutValue(s.PopX, words, count, overflow, queue)
			v1, ok1 := portTestWorklistOutValue(s.PopY, words, count, overflow, queue)
			v2, ok2 := portTestWorklistOutValue(s.PopZ, words, count, overflow, queue)
			snap.Results = append(snap.Results, PortTestTileWorklistResult{Return: int(ret), Count: *count, Overflow: *overflow, Outputs: [3]uint32{v0, v1, v2}, OutputReadable: [3]bool{ok0, ok1, ok2}, Queue: append([]uint32(nil), queue...), GridUnchanged: C.portTestWorklistGridEqual(grid, expected) != 0, GridPointerUnchanged: pointerOK, QueueGuardsUnchanged: string(left) == string(wantLeft) && string(right) == string(wantRight), OutputGuardsOK: words[0] == 0xa0a0a0a0 && words[4] == 0xb0b0b0b0})
		} else if s.Op == 1 {
			ret = C.sub_51DE30(portTestWorklistOut(s.PopX, words, count, overflow, queue), portTestWorklistOut(s.PopY, words, count, overflow, queue), portTestWorklistOut(s.PopZ, words, count, overflow, queue))
			v0, ok0 := portTestWorklistOutValue(s.PopX, words, count, overflow, queue)
			v1, ok1 := portTestWorklistOutValue(s.PopY, words, count, overflow, queue)
			v2, ok2 := portTestWorklistOutValue(s.PopZ, words, count, overflow, queue)
			snap.Results = append(snap.Results, PortTestTileWorklistResult{Return: int(ret), Count: *count, Overflow: *overflow, Outputs: [3]uint32{v0, v1, v2}, OutputReadable: [3]bool{ok0, ok1, ok2}, Queue: append([]uint32(nil), queue...), GridUnchanged: C.portTestWorklistGridEqual(grid, expected) != 0, GridPointerUnchanged: C.portTestWorklistGridGet() == grid, QueueGuardsUnchanged: string(left) == string(wantLeft) && string(right) == string(wantRight), OutputGuardsOK: words[0] == 0xa0a0a0a0 && words[4] == 0xb0b0b0b0})
		} else {
			C.free(mem)
			panic("invalid worklist operation")
		}
		C.free(mem)
	}
	return snap
}

// portTestMainGrid provides real C/Go tile lookups for dodge decisions.
func portTestMainGrid() (configure func(uint32), intact func() bool, free func()) {
	old := C.portTestWorklistGridGet()
	grid, expected := C.portTestWorklistGridNew(), C.portTestWorklistGridNew()
	if grid == nil || expected == nil {
		panic("main fixture grid allocation")
	}
	C.portTestWorklistGridSet(grid)
	last := ^uint32(0)
	configure = func(tile uint32) {
		if tile == last {
			return
		}
		last = tile
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				C.portTestWorklistCellSet(grid, C.int(x), C.int(y), C.uint32_t(tile), C.uint32_t(tile))
				C.portTestWorklistCellSet(expected, C.int(x), C.int(y), C.uint32_t(tile), C.uint32_t(tile))
			}
		}
	}
	intact = func() bool {
		return C.portTestWorklistGridGet() == grid && C.portTestWorklistGridEqual(grid, expected) != 0
	}
	free = func() {
		C.portTestWorklistGridSet(old)
		C.portTestWorklistGridFree(grid)
		C.portTestWorklistGridFree(expected)
	}
	return
}
