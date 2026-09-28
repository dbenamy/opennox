//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"runtime"
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
func portTestWorklistOut(ref uint16, words []uint32, count, overflow *uint32, queue []uint32) *uint32 {
	switch {
	case ref < 3:
		return &words[1+ref] // words[0]/[4] are guards
	case ref == 3:
		return (*uint32)(unsafe.Pointer(count))
	case ref == 4:
		return (*uint32)(unsafe.Pointer(overflow))
	case ref == portTestWorklistNilRef:
		return nil
	case int(ref)-5 < len(queue):
		return (*uint32)(unsafe.Pointer(&queue[int(ref)-5]))
	default:
		panic("invalid worklist output reference")
	}
}
func portTestWorklistOutValue(ref uint16, words []uint32, count, overflow *uint32, queue []uint32) (uint32, bool) {
	if ref == portTestWorklistNilRef {
		return 0, false
	}
	return uint32(*portTestWorklistOut(ref, words, count, overflow, queue)), true
}

// PortTestTileWorklist invokes the native producer/consumer against an
// isolated allocator-owned 128x128 grid. It retains every physical queue word after
// each call, including inactive records, so a port cannot merely model a slice.
func PortTestTileWorklist(initialCount, initialOverflow uint32, initialQueue []uint32, specs []PortTestTileWorklistSpec) (snap PortTestTileWorklistSnapshot) {
	if len(initialQueue) != portTestWorklistWords {
		panic("worklist initial queue must contain 1500 words")
	}
	count := (*uint32)(unsafe.Pointer(&dword_5d4594_2487248))
	overflow := memmap.PtrUint32(0x973F18, 22200)
	queue := portTestWorklistQueue()
	left := unsafe.Slice(memmap.PtrUint8(0x973F18, 16192), 8)
	right := unsafe.Slice(memmap.PtrUint8(0x973F18, 22204), 8)
	oldGrid := worldTileGrid
	snap.Before = portTestWorklistState(count, overflow, queue, left, right)
	grid, expected := portTestWorklistGridNew(), portTestWorklistGridNew()
	if grid == nil || expected == nil {
		portTestWorklistGridFree(grid)
		portTestWorklistGridFree(expected)
		panic("calloc grid failed")
	}
	defer func() {
		copy(queue, snap.Before.Queue)
		*count, *overflow = snap.Before.Count, snap.Before.Overflow
		copy(left, snap.Before.Left)
		copy(right, snap.Before.Right)
		worldTileGrid = oldGrid
		portTestWorklistGridFree(grid)
		portTestWorklistGridFree(expected)
		snap.AfterRestore = portTestWorklistState(count, overflow, queue, left, right)
		snap.GridPointerRestored = worldTileGrid == oldGrid
	}()
	worldTileGrid = (**worldTileCell)(unsafe.Pointer(grid))
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
		mem := legacyCalloc(5, uintptr(unsafe.Sizeof(uint32(0))))
		if mem == nil {
			panic("calloc outputs failed")
		}
		words := unsafe.Slice((*uint32)(mem), 5)
		words[0], words[1], words[2], words[3], words[4] = 0xa0a0a0a0, 0x01020304, 0x11223344, 0x55667788, 0xb0b0b0b0
		var ret int32
		if s.Op == 0 {
			if s.X > 0 && s.X < 127 && s.Y > 0 && s.Y < 127 {
				portTestWorklistCellSet(grid, int32(s.X), int32(s.Y), uint32(s.Field1), uint32(s.Field2))
				portTestWorklistCellSet(expected, int32(s.X), int32(s.Y), uint32(s.Field1), uint32(s.Field2))
			}
			wantGrid := grid
			if s.NilGrid {
				wantGrid = nil
				worldTileGrid = nil
			}
			sub_51DD50(int32(s.X), int32(s.Y), int32(s.Flags), int32(s.Key))
			pointerOK := unsafe.Pointer(worldTileGrid) == unsafe.Pointer(wantGrid)
			if s.NilGrid {
				worldTileGrid = (**worldTileCell)(unsafe.Pointer(grid))
			}
			v0, ok0 := portTestWorklistOutValue(s.PopX, words, count, overflow, queue)
			v1, ok1 := portTestWorklistOutValue(s.PopY, words, count, overflow, queue)
			v2, ok2 := portTestWorklistOutValue(s.PopZ, words, count, overflow, queue)
			snap.Results = append(snap.Results, PortTestTileWorklistResult{Return: int(ret), Count: *count, Overflow: *overflow, Outputs: [3]uint32{v0, v1, v2}, OutputReadable: [3]bool{ok0, ok1, ok2}, Queue: append([]uint32(nil), queue...), GridUnchanged: portTestWorklistGridEqual(grid, expected) != 0, GridPointerUnchanged: pointerOK, QueueGuardsUnchanged: string(left) == string(wantLeft) && string(right) == string(wantRight), OutputGuardsOK: words[0] == 0xa0a0a0a0 && words[4] == 0xb0b0b0b0})
		} else if s.Op == 1 {
			ret = sub_51DE30(portTestWorklistOut(s.PopX, words, count, overflow, queue), portTestWorklistOut(s.PopY, words, count, overflow, queue), portTestWorklistOut(s.PopZ, words, count, overflow, queue))
			v0, ok0 := portTestWorklistOutValue(s.PopX, words, count, overflow, queue)
			v1, ok1 := portTestWorklistOutValue(s.PopY, words, count, overflow, queue)
			v2, ok2 := portTestWorklistOutValue(s.PopZ, words, count, overflow, queue)
			snap.Results = append(snap.Results, PortTestTileWorklistResult{Return: int(ret), Count: *count, Overflow: *overflow, Outputs: [3]uint32{v0, v1, v2}, OutputReadable: [3]bool{ok0, ok1, ok2}, Queue: append([]uint32(nil), queue...), GridUnchanged: portTestWorklistGridEqual(grid, expected) != 0, GridPointerUnchanged: unsafe.Pointer(worldTileGrid) == unsafe.Pointer(grid), QueueGuardsUnchanged: string(left) == string(wantLeft) && string(right) == string(wantRight), OutputGuardsOK: words[0] == 0xa0a0a0a0 && words[4] == 0xb0b0b0b0})
		} else {
			legacyFree(mem)
			panic("invalid worklist operation")
		}
		legacyFree(mem)
	}
	return snap
}

// portTestMainGrid provides real C/Go tile lookups for dodge decisions.
func portTestMainGrid() (configure func(uint32), intact func() bool, free func()) {
	old := worldTileGrid
	grid, expected := portTestWorklistGridNew(), portTestWorklistGridNew()
	if grid == nil || expected == nil {
		panic("main fixture grid allocation")
	}
	worldTileGrid = (**worldTileCell)(unsafe.Pointer(grid))
	last := ^uint32(0)
	configure = func(tile uint32) {
		if tile == last {
			return
		}
		last = tile
		for x := 0; x < 128; x++ {
			for y := 0; y < 128; y++ {
				portTestWorklistCellSet(grid, int32(x), int32(y), uint32(tile), uint32(tile))
				portTestWorklistCellSet(expected, int32(x), int32(y), uint32(tile), uint32(tile))
			}
		}
	}
	intact = func() bool {
		return unsafe.Pointer(worldTileGrid) == unsafe.Pointer(grid) && portTestWorklistGridEqual(grid, expected) != 0
	}
	free = func() {
		worldTileGrid = old
		portTestWorklistGridFree(grid)
		portTestWorklistGridFree(expected)
	}
	return
}

// Reuse the owned tile grid used by the worklist contracts. Both tile halves
// receive the requested floor kind; the production world-to-tile lookup stays live.
func PortTestWorldMotionTileGrid() (func(int32), func() bool, func()) {
	old := worldTileGrid
	grid, want := portTestWorklistGridNew(), portTestWorklistGridNew()
	if grid == nil || want == nil {
		portTestWorklistGridFree(grid)
		portTestWorklistGridFree(want)
		panic("world motion tile allocation")
	}
	worldTileGrid = (**worldTileCell)(unsafe.Pointer(grid))
	configure := func(kind int32) {
		for _, table := range []**worldTileCell{grid, want} {
			for _, p := range unsafe.Slice(table, 128) {
				for i := range unsafe.Slice(p, 128) {
					row := unsafe.Slice(p, 128)
					row[i][1] = uint32(kind)
					row[i][6] = uint32(kind)
				}
			}
		}
	}
	return configure, func() bool {
			return unsafe.Pointer(worldTileGrid) == unsafe.Pointer(grid) && portTestWorklistGridEqual(grid, want) != 0
		}, func() {
			worldTileGrid = old
			portTestWorklistGridFree(grid)
			portTestWorklistGridFree(want)
		}
}

// These fixture helpers retain the original row-by-row allocation and byte checks.
func portTestWorklistGridNew() **worldTileCell {
	p := (**worldTileCell)(legacyCalloc(128, unsafe.Sizeof((*worldTileCell)(nil))))
	if p == nil {
		return nil
	}
	rows := unsafe.Slice(p, 128)
	for i := range rows {
		rows[i] = (*worldTileCell)(legacyCalloc(128, unsafe.Sizeof(worldTileCell{})))
		if rows[i] == nil {
			for j := i; j > 0; {
				j--
				legacyFree(unsafe.Pointer(rows[j]))
			}
			legacyFree(unsafe.Pointer(p))
			return nil
		}
	}
	return p
}
func portTestWorklistGridFree(p **worldTileCell) {
	if p == nil {
		return
	}
	for _, row := range unsafe.Slice(p, 128) {
		legacyFree(unsafe.Pointer(row))
	}
	legacyFree(unsafe.Pointer(p))
}
func portTestWorklistCellSet(p **worldTileCell, x, y int32, one, two uint32) {
	row := unsafe.Slice(unsafe.Slice(p, 128)[x], 128)
	row[y][1], row[y][6] = one, two
}
func portTestWorklistGridEqual(a, b **worldTileCell) int32 {
	ar, br := unsafe.Slice(a, 128), unsafe.Slice(b, 128)
	for i := range ar {
		if !bytes.Equal(unsafe.Slice((*byte)(unsafe.Pointer(ar[i])), 128*44), unsafe.Slice((*byte)(unsafe.Pointer(br[i])), 128*44)) {
			return 0
		}
	}
	return 1
}

// PortTestWorklistAllocation observes fixture-grid allocation and partial cleanup.
type PortTestWorklistAllocationResult struct {
	Success, Zero, ValidFrees, NilTracked                        bool
	Sizes                                                        []int
	Freed, Remaining, TrackedBefore, TrackedDuring, TrackedAfter int
}

func PortTestWorklistAllocation(failAt int) (out PortTestWorklistAllocationResult) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if alloc.PortTestAllocationLive(nil) {
		panic("preexisting nil allocation marker")
	}
	out.TrackedBefore = alloc.PortTestAllocationCount()
	portTestGridAllocationObserve(failAt)
	defer portTestGridAllocationStop()
	p := portTestWorklistGridNew()
	out.Success = p != nil
	out.Zero = true
	if p != nil {
		for _, row := range unsafe.Slice((**worldTileCell)(unsafe.Pointer(p)), 128) {
			for _, cell := range unsafe.Slice(row, 128) {
				for _, word := range cell {
					out.Zero = out.Zero && word == 0
				}
			}
		}
	}
	out.TrackedDuring = alloc.PortTestAllocationCount()
	portTestWorklistGridFree(p)
	out.NilTracked = alloc.PortTestAllocationLive(nil)
	// Safe calloc records nil on failure. Preserve that observation, then remove
	// only this fixture's marker so subsequent cases begin with clean ownership.
	if out.NilTracked {
		alloc.FreePtr(nil)
	}
	out.TrackedAfter = alloc.PortTestAllocationCount()
	n := portTestGridAllocationStat(-1)
	out.Freed = portTestGridAllocationStat(-2)
	out.Remaining = n - out.Freed
	out.ValidFrees = portTestGridAllocationStat(-3) != 0
	for i := 0; i < n; i++ {
		out.Sizes = append(out.Sizes, portTestGridAllocationStat(i))
	}
	return
}
