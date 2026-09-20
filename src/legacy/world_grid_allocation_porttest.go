//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1.h"
extern obj_5D4594_2650668_t** ptr_5D4594_2650668;
void worldGridAllocObserve(int fail);
void worldGridAllocStop(void);
int worldGridAllocStat(int index);
int worldGridAllocContains(void* p);
*/
import "C"

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
	oldGrid := C.ptr_5D4594_2650668
	r.FailAt = failAt
	r.Zero = true
	C.worldGridAllocObserve(C.int(failAt))
	defer func() { C.worldGridAllocStop(); C.ptr_5D4594_2650668 = oldGrid }()
	r.Return = int(C.nox_xxx_tileAlloc_410F60_init())
	outer := unsafe.Pointer(C.ptr_5D4594_2650668)
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
		before := int(C.worldGridAllocStat(-1))
		C.nox_xxx_tileFree_410FC0_free()
		r.RowsFreed = int(C.worldGridAllocStat(-2))
		r.OuterRetained = C.worldGridAllocContains(outer) != 0
		r.OuterRetained = r.OuterRetained && unsafe.Pointer(C.ptr_5D4594_2650668) == outer
		r.Remaining = before - r.RowsFreed
		C.free(outer)
	}
	count := int(C.worldGridAllocStat(-1))
	r.FinalRemaining = count - int(C.worldGridAllocStat(-2))
	r.ValidFrees = C.worldGridAllocStat(-3) != 0
	for i := 0; i < count; i++ {
		r.Sizes = append(r.Sizes, int(C.worldGridAllocStat(C.int(i))))
	}
	return r
}
