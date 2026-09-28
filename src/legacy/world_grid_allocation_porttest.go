//go:build porttest

package legacy

import (
	"runtime"
	"unsafe"
)

type PortTestWorldGridAllocation struct {
	FailAt, Return, Rows, RowsFreed, Remaining, FinalRemaining int
	Sizes                                                      []int
	Zero, OuterRetained, ValidFrees                            bool
}

func PortTestWorldGridAllocate(failAt int) (r PortTestWorldGridAllocation) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	oldGrid := worldTileGrid
	r.FailAt = failAt
	r.Zero = true
	allocationTestGridStart(int(failAt))
	defer func() { allocationTestGridStop(); worldTileGrid = oldGrid }()
	r.Return = int(worldGridAllocate())
	outer := unsafe.Pointer(worldTileGrid)
	if outer != nil {
		rows := unsafe.Slice((*unsafe.Pointer)(outer), 128)
		for _, p := range rows {
			if p != nil {
				r.Rows++
				for _, v := range unsafe.Slice((*uint32)(p), 128*11) {
					r.Zero = r.Zero && v == 0
				}
			}
		}
		before := int(allocationTestGridStat(-1))
		worldGridFreeRows()
		r.RowsFreed = int(allocationTestGridStat(-2))
		r.OuterRetained = allocationTestGridContains(outer) != 0
		r.OuterRetained = r.OuterRetained && unsafe.Pointer(worldTileGrid) == outer
		r.Remaining = before - r.RowsFreed
		legacyFree(outer)
	}
	count := int(allocationTestGridStat(-1))
	r.FinalRemaining = count - int(allocationTestGridStat(-2))
	r.ValidFrees = allocationTestGridStat(-3) != 0
	for i := 0; i < count; i++ {
		r.Sizes = append(r.Sizes, int(allocationTestGridStat(int(i))))
	}
	return r
}

// Shared by fixture-grid allocation contracts while the C observer remains.
func portTestGridAllocationObserve(failAt int) { allocationTestGridStart(int(failAt)) }
func portTestGridAllocationStop()              { allocationTestGridStop() }
func portTestGridAllocationStat(index int) int { return int(allocationTestGridStat(int(index))) }
